package main

import (
	"flag"
	"fmt"
	"github.com/saeveritt/go-peerassets/api"
	"github.com/saeveritt/go-peerassets/storage"
	"github.com/saeveritt/go-peerassets/utils"
	"log"
	"net/http"
)

func init(){
	loadFlag := flag.Bool("load", false, "Import Root P2TH Assets")
	flag.Parse()

	if *loadFlag {
		utils.ImportRootP2TH()
		storage.PutRootAsset()
	}
}

func main() {



	server := "0.0.0.0"
	port := "8089"
	r := api.AgaveRouter()
	fmt.Println("\nStarting go-peerassets server...")
	fmt.Println("----Success! Running on "+server + ":" + port)
	log.Fatal(http.ListenAndServe(server + ":" + port, r))
}



