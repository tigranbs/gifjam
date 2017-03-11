package main

import (
	"fmt"
	"gifjam/gifGrabber"
)

func main() {
	items, err := gifGrabber.GetFeed("gifporn1")
	if err != nil {
		fmt.Println("Unable to get Feed -> ", err.Error())
		return
	}

	fmt.Println("Items Count ->", len(items))

	items, err = gifGrabber.FilterGifs(items)
	if err != nil {
		fmt.Println("Unable to get Gifs -> ", err.Error())
		return
	}

	for _, item :=range items {
		new_file, err := gifGrabber.SaveItem(&item)
		if err != nil {
			fmt.Println("Unable to save file -> ", err.Error())
			continue
		}

		if !new_file {
			fmt.Println("Item Exists already -> ", item.Link)
		}
	}
}
