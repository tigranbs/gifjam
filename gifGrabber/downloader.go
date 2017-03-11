package gifGrabber

import (
	"log"
	"time"
)

// Calling Facebook API every 10 minutes
const API_CALL_INTERVAL = 600

func StartDownloader() {
	for {
		// getting Feed Items
		items, err := GetFeed("gifporn1")
		if err != nil {
			log.Println("Unable to get Feed -> ", err.Error())
			return
		}

		items, err = FilterGifs(items)
		if err != nil {
			log.Println("Unable to Filter Gifs -> ", err.Error())
			return
		}

		go SaveItems(items)

		time.Sleep(time.Second * API_CALL_INTERVAL)
	}
}

func SaveItems(items []FeedItem) {
	for _, item :=range items {
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