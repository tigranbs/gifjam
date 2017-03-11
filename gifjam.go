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

	fmt.Println("Only GIF Items Count ->", len(items))
}
