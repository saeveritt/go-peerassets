package main

import (
	"flag"
	"fmt"
	"github.com/saeveritt/go-peerassets/app/api"
	"github.com/saeveritt/go-peerassets/app/storage"
	"github.com/saeveritt/go-peerassets/app/utils"
	"net/http"
	"log"
)

func init(){
	loadFlag := flag.Bool("load", false, "Import Root P2TH Assets")
	portFlag := flag.String("port", "8089", "Port to start server on")
	startFlag := flag.Bool("server",false,"Start server")
	cardFlag := flag.Bool("cards",false,"Load All Cards")
	flag.Parse()

	if *loadFlag {
		utils.ImportRootP2TH()
		storage.PutRootAsset()
		storage.ImportSubscribed()
	}
	if *cardFlag{
		storage.ImportSubscribedCards()
	}
	StartServer(*startFlag, *portFlag)

}

func main() {
}

func StartServer(server bool,port string){
	if server {
		server := "0.0.0.0"
		r := api.AgaveRouter()
		fmt.Println("\nStarting go-peerassets server...")
		fmt.Println("----Success! Running on " + server + ":" + port)
		log.Fatal(http.ListenAndServe(server+":"+port, r))
	}
}



