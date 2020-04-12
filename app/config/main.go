package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config struct{
	Subscribed Subscribed `json:"subscribed"`
}

type Subscribed struct{
	All 		bool 		`json:"all"`
	Decks 		[]string	`json:"decks"`
}

func Open() (Config, error){
	var config Config
	pwd, _ := os.Getwd()
	file, err := ioutil.ReadFile(pwd+"/config.json")
	_ = json.Unmarshal(file, &config)
	if err != nil{
		err = errors.New("config.json is missing from root directory.")
	}
	return config, err
}