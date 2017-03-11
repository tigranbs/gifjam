package config

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	Host    string   `json:"host"`
	FbPages []string `json:"fb_pages"`
	Mongo   string   `json:"mongo"`
	FbToken string   `json:"fb_token"`
}

var (
	GlobalConfig Config
)

func ParseConfig() {
	if len(os.Args) != 2 {
		fmt.Println("Application should get only configuration file as an argument")
	}
	// Parsing config file
	conf_data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Unable to read given configuration file", os.Args[1], "->", err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(conf_data, &GlobalConfig)
	if err != nil {
		fmt.Println("Unable to parse given configuration file", os.Args[1], "->", err.Error())
		os.Exit(1)
	}
}