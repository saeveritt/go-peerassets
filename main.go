package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/saeveritt/go-peerassets/storage"
	"github.com/saeveritt/go-peerassets/utils"
	"log"
	"net/http"
)

func init(){
	utils.ImportRootP2TH()
	storage.PutRootAsset()
}

var(
)

func getAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := storage.GetDecks()
	w.Write(j)

}

func postAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}


func main() {
	server := "localhost"
	port := "8089"
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/assets", getAssets).Methods(http.MethodGet)
	api.HandleFunc("/assets", postAssets).Methods(http.MethodPost)
	api.HandleFunc("/address", postAssets).Methods(http.MethodPost)
	fmt.Println("\nStarting go-peerassets server...")
	fmt.Println("----Success! Running on "+server + ":" + port)
	log.Fatal(http.ListenAndServe(server + ":" + port, r))
}



