package main

import (
	_ "gifjam/gifGrabber"
	"gifjam/api"
)

func main() {
	//go gifGrabber.StartDownloader()
	api.StartApiServer()
}
