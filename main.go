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
	portFlag := flag.String("port", "8089", "Port to start server on")
	flag.Parse()
	if *loadFlag {
		utils.ImportRootP2TH()
		storage.PutRootAsset()
	}
	StartServer(*portFlag)

}

func main() {
}

func StartServer(port string){
	server := "0.0.0.0"
	r := api.AgaveRouter()
	fmt.Println("\nStarting go-peerassets server...")
	fmt.Println("----Success! Running on "+server + ":" + port)
	log.Fatal(http.ListenAndServe(server + ":" + port, r))

}



