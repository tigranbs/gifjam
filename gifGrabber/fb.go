package gifGrabber

import (
	"fmt"
	"encoding/json"
	fb "github.com/huandu/facebook"
	"os"
)

var (
	accessToken = os.Getenv("FACEBOOK_TOKEN")
)

type FeedItem struct {
	ID string `json:"id"`
	Message string `json:"message"`
	CreatedTime string `json:"created_time"`
	Link string `json:"link"`
	Picture string `json:"full_picture"`
	Type string `json:"type"`
}

func GetFeed(page_id string) ([]FeedItem, error) {
	res, err := fb.Get("/" + page_id + "/feed", fb.Params{
		"access_token": accessToken,
		"fields": "message,created_time,link,full_picture,type",
		"limit": "100",
	})
	if err != nil {
		return nil, err
	}

	data_json, err := json.Marshal(res.Get("data"))
	if err != nil {
		fmt.Println("Unabel to unmarshal data -> ", err.Error())
		return
	}

	items := []FeedItem{}

	err = json.Unmarshal(data_json, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}