package main

import (
	"github.com/saeveritt/go-peerassets/storage"
	"github.com/saeveritt/go-peerassets/utils"
)

func init(){
	utils.ImportRootP2TH()
	utils.Scan(0)
}

func main() {
	storage.Connect()
	storage.PutRootAsset()
	storage.Close()
}



