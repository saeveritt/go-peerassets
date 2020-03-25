package api

import (
	"github.com/gorilla/mux"
	"github.com/saeveritt/go-peerassets/storage"
	"net/http"
	"strconv"
)


func AgaveRouter() *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/v1").Subrouter()
	api.HandleFunc("/assets", getAssets).Methods(http.MethodGet)
	api.HandleFunc("/address",getAddress).Methods(http.MethodGet)
	return r
}

func getAddress( w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	var j []byte
	var address = r.URL.Query().Get("address")
	j,_ = storage.GetAddress(address)
	w.Write(j)
}



func getAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var j []byte

	limit := r.URL.Query().Get("limit")
	l, _ := strconv.Atoi(limit)
	page := r.URL.Query().Get("page")
	p, _ := strconv.Atoi(page)

	if page == "" && limit == "" {
		j, _ = storage.GetDecks()
	}
	if p > 0 && limit == "" {
		j, _ = storage.GetDecksPages(10, p)
	}
	if l > 0 && page == "" {
		j, _ = storage.GetDecksPages(l, 1)
	}
	if p > 0 && l > 0{
		j, _ = storage.GetDecksPages(l, p)
	}

	w.Write(j)
}