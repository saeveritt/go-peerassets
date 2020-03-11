package main

import (
	"github.com/saeveritt/go-peerassets/utils"
	"github.com/saeveritt/go-peerassets/storage"
)

func init(){
	utils.ImportRootP2TH()
}

func main() {
	storage.Connect()
	storage.PutRootAsset()
	storage.Close()
}



