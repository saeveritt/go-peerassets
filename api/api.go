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
	//api.HandleFunc("/assets", postAssets).Methods(http.MethodPost)
	//api.HandleFunc("/address", postAssets).Methods(http.MethodPost)
	return r
}


func getAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var j []byte
	limit := r.Form.Get("limit")
	l, _ := strconv.Atoi(limit)
	page := r.Form.Get("page")
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

	w.Write(j)
}