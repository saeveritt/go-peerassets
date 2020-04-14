package main

import (
	"github.com/saeveritt/go-peerassets/app/api"
	"github.com/saeveritt/go-peerassets/app/rpc"
	"github.com/saeveritt/go-peerassets/app/storage"
	"github.com/saeveritt/go-peerassets/app/utils"
	"log"
	"net/http"
	"time"
)

func init(){
	for !Load(){time.Sleep( 3 * time.Second)}
}

func main() {
	StartServer(true, "8089")
}

func StartServer(server bool,port string){
	if server {
		server := "0.0.0.0"
		r := api.AgaveRouter()
		log.Println("Starting go-peerassets server...")
		log.Println("----Success! Running on " + server + ":" + port)
		log.Fatal(http.ListenAndServe(server+":"+port, r))
	}
}

func Load() (start bool){
	cli, _ := rpc.Connect("Peercoin-Testnet")
	info, _ := cli.GetInfo()
	if info.Blocks == info.Headers && info.Headers != 0 {
		log.Printf("%v of %v Blocks", info.Blocks, info.Headers)
		utils.ImportRootP2TH()
		storage.PutRootAsset()
		storage.ImportSubscribed()
		storage.ImportSubscribedCards()
		return true
	}
	if info.Headers == 0{
		log.Print("Connecting to local RPC node...")
	}else {
		log.Printf("%v of %v Blocks", info.Blocks, info.Headers)
	}
	return false
}


