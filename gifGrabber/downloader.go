package gifGrabber

import (
	"log"
	"time"
	"os"
	"gifjam/config"
)

// Calling Facebook API every 10 minutes
const API_CALL_INTERVAL = 600
var (
	fbToken string
)

func StartDownloader() {
	// Init MongoDB
	initDB()

	if len(config.GlobalConfig.FbToken) == 0 {
		fbToken = os.Getenv("GIFJAM_FACEBOOK_TOKEN")
	} else {
		fbToken = config.GlobalConfig.FbToken
	}

	for {
		items := []FeedItem{}

		for _, page :=range config.GlobalConfig.FbPages {
			// getting Feed Items
			it, err := GetFeed(page)
			if err != nil {
				log.Println("Unable to get Feed -> ", err.Error())
				continue
			}

			items = append(items, it...)
		}

		if len(items) > 0 {
			items, err := FilterGifs(items)
			if err != nil {
				log.Println("Unable to Filter Gifs -> ", err.Error())
				time.Sleep(time.Second * API_CALL_INTERVAL)
				continue
			}

			go SaveItems(items)
		}

		time.Sleep(time.Second * API_CALL_INTERVAL)
	}
}

func SaveItems(items []FeedItem) {
	for _, item := range items {
		new_file, err := SaveItem(&item)
		if err != nil {
			log.Println("Unable to save file -> ", err.Error())
			continue
		}

		if !new_file {
			log.Println("Item Exists already -> ", item.Link)
		}
	}
}
