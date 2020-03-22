package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/saeveritt/go-peerassets/storage"
	"net/http"
	"strconv"
)


func AgaveRouter() *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/v1").Subrouter()
	api.HandleFunc("/assets", getAssets).Methods(http.MethodGet)
	api.HandleFunc("/assets", postAssets).Methods(http.MethodPost)
	api.HandleFunc("/address", postAssets).Methods(http.MethodPost)
	return r
}


func getAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := storage.GetDecks()
	w.Write(j)

}

func postAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	r.ParseForm()
	limit := r.Form.Get("limit")
	l,_ := strconv.Atoi(limit)
	page := r.Form.Get("page")
	p,_ := strconv.Atoi(page)
	j,_ := storage.GetDecksPages(l,p)
	json.

	w.Write([]byte(`{"message": "post called"}`))
}