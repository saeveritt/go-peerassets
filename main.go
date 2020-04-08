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
	startFlag := flag.Bool("server",false,"Start server")
	cardFlag := flag.Bool("cards",false,"Load All Cards")
	flag.Parse()
	if *loadFlag {
		utils.ImportRootP2TH()
		storage.PutRootAsset(false)
	}
	if *cardFlag{
		storage.PutAllCards()
	}
	StartServer(*startFlag, *portFlag)

}

func main() {
	//utils.GetCards("d460651e1d9147770ec9d4c254bcc68ff5d203a86b97c09d00955fb3f714cab3")
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



