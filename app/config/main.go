package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config struct{
	Subscribed Subscribed `json:"subscribed"`
	Environment Environment
}

type Subscribed struct{
	All 		bool 		`json:"all"`
	Decks 		[]string	`json:"decks"`
}

type Environment struct{
	rpcUsername string
	rpcPassword string
	rpcHost	string
	rpcPort	string

}


func Open() (Config, error){
	var config Config

	env := Environment{
		os.Getenv("RPC_USERNAME"),
		os.Getenv("RPC_PASSWORD"),
		os.Getenv("RPC_HOST"),
		os.Getenv("RPC_PORT"),
	}

	config.Environment = env
	pwd, _ := os.Getwd()
	file, err := ioutil.ReadFile(pwd+"/config.json")
	_ = json.Unmarshal(file, &config)
	if err != nil{
		err = errors.New("config.json is missing from root directory.")
	}
	return config, err
}