package main

import (
	"github.com/saeveritt/go-peerassets/storage"
	//"github.com/saeveritt/go-peerassets/utils"
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

func init(){
	//utils.ImportRootP2TH()
	//utils.Scan(0)
	//storage.Connect()
	//storage.PutRootAsset()
	//storage.Close()
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
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/assets", getAssets).Methods(http.MethodGet)
	api.HandleFunc("/assets", postAssets).Methods(http.MethodPost)
	api.HandleFunc("/address", postAssets).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8089", r))
}



