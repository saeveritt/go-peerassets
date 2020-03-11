package main

import (
	"go-peerassets/utils"
	"go-peerassets/storage"
)

func init(){
	utils.ImportRootP2TH()
}

func main() {
	storage.Connect()
	storage.PutRootAsset()
	storage.Close()
}



