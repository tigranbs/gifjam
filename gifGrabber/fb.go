package gifGrabber

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	fb "github.com/huandu/facebook"
	"log"
	"net/http"
	"strings"
	"time"
)

type FeedItem struct {
	ID          string `json:"id"`
	Message     string `json:"message"`
	CreatedTime string `json:"created_time"`
	Link        string `json:"link"`
	Picture     string `json:"full_picture"`
	Type        string `json:"type"`
	// This field used by Storage metadata
	Visible bool `json:"-"`
}

func GetFeed(page_id string) ([]FeedItem, error) {
	res, err := fb.Get("/"+page_id+"/feed", fb.Params{
		"access_token": fbToken,
		"fields":       "message,created_time,link,full_picture,type",
		"limit":        "100",
	})
	if err != nil {
		return nil, err
	}

	data_json, err := json.Marshal(res.Get("data"))
	if err != nil {
		return nil, err
	}

	items := []FeedItem{}

	err = json.Unmarshal(data_json, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func FilterGifs(items []FeedItem) ([]FeedItem, error) {
	gif_items := []FeedItem{}
	for _, item := range items {
		if item.Type != "link" {
			continue
		}

		retry_count := 0

		for {
			if retry_count > 20 {
				log.Println("Http Request Retry count existed for url", item.Link, " | Moving to the next element")
				break
			}

			resp, err := http.Get(item.Link)
			if err != nil {
				log.Println("Unable to make a request to url", item.Link, "-> ", err.Error(), " | Trying again")
				time.Sleep(time.Second * 1)
				retry_count++
				continue
			}

			content_type := resp.Header.Get("Content-Type")

			// if this is a direct GIF link then adding it to our list
			if strings.Contains(content_type, "image/gif") {
				gif_items = append(gif_items, item)
			} else if strings.Contains(content_type, "text/html") { // if this is an html content then probably GIF is inside meta "og:image" tag
				document, err := goquery.NewDocumentFromReader(resp.Body)
				resp.Body.Close()
				resp.Close = true
				if err != nil {
					log.Println("Unable to read HTML content from url", item.Link, " | Trying again")
					retry_count++
					continue
				}

				attr, exists := document.Find(`meta[property="og:image"]`).Attr("content")
				if !exists {
					// trying again with og:url
					attr, exists = document.Find(`meta[property="og:url"]`).Attr("content")
				}

				if !exists {
					log.Println("There is No Gif for URl ->", item.Link)
				} else {
					// changing gif link to real image link
					item.Link = attr
					gif_items = append(gif_items, item)
				}
			}

			break
		}
	}

	return gif_items, nil
}
