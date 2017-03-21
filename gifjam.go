package main

import (
	"gifjam/api"
	"gifjam/gifGrabber"
	"gifjam/config"
	"os"
	"fmt"
)

const HelpText =
	`GifJam Application should receive only 2 parameters to execute different tasks
	gifjam <parameter> <configuration file>

	Parameters
		grab	Using this parameters it will start one time grabber action,
				probably this should be included in Cron Task
		server	Starting API server

	Configuration file
		This field should contain path to configuration file
		In JSON format
`

func main() {
	config.ParseConfig()

	if len(os.Args) == 1 {
		fmt.Println(HelpText)
		return
	}

	switch os.Args[1] {
	case "grab":
		gifGrabber.StartDownloader()
	case "server":
		api.StartApiServer()
	default:
		fmt.Println(HelpText)
		return
	}
}
