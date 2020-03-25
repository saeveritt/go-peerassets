package main

import (
	"fmt"
	"github.com/saeveritt/go-peerassets/api"
	"log"
	"net/http"
)

func init(){
	//utils.ImportRootP2TH()
	//storage.PutRootAsset()
}

func main() {
	server := "0.0.0.0"
	port := "8089"
	r := api.AgaveRouter()
	fmt.Println("\nStarting go-peerassets server...")
	fmt.Println("----Success! Running on "+server + ":" + port)
	log.Fatal(http.ListenAndServe(server + ":" + port, r))
}



