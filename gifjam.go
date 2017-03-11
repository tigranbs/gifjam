package main

import (
	"gifjam/api"
	"gifjam/gifGrabber"
	"gifjam/config"
)

func main() {
	config.ParseConfig()
	go gifGrabber.StartDownloader()
	api.StartApiServer()
}
