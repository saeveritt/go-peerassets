package api

import (
	"github.com/gorilla/mux"
	"github.com/saeveritt/go-peerassets/storage"
	"net/http"
	"strconv"
)


func AgaveRouter() *mux.Router {
	// Create new router
	r := mux.NewRouter()
	// Set the subrouter path prefix
	api := r.PathPrefix("/v1").Subrouter()
	// Define the function handlers per route
	api.HandleFunc("/assets", getAssets).Methods(http.MethodGet)
	api.HandleFunc("/transactions",getAddress).Methods(http.MethodGet)
	return r
}

func getAddress( w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	// Create empty byte array which will store the JSON Response
	var j []byte
	// Assign address variable to the value passed in the GET request
	var address = r.URL.Query().Get("address")
	// If address is not empty
	if address != ""{
		// Check for address in storage. Each address has its own dedicated bucket.
		j, _ = storage.GetAddress(address)
		// if there was an error writing the JSON byte array,it will send empty array
		// else it sends a JSON byte array Response with the results
		w.Write(j)
	}else{
		// if address argument is empty, return empty byte array
		w.Write(j)
	}
}

func getAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var j []byte

	limit := r.URL.Query().Get("limit")
	// Because this is a string, convert l to an integer
	l, err := strconv.Atoi(limit)
	// TODO: return error in json
	if err != nil{w.Write(j)}
	page := r.URL.Query().Get("page")
	// Because this is a string, convert p to an integer
	p, err := strconv.Atoi(page)
	// TODO: return error in json
	if err != nil{w.Write(j)}

	if p == 0 && l == 0 {
		// If page and limit are both 0, return all decks
		j, _ = storage.GetDecks(0, 0)
	}
	if p > 0 && l == 0 {
		// If page is greater then zero, return page
		j, _ = storage.GetDecks(10, p)
	}
	if p == 0 && l > 0 {
		// If limit is greater then zero, send one page with l amount
		j, _ = storage.GetDecks(l, 1)
	}
	if p > 0 && l > 0{
		// if both page and l are greater than zero, return l amount on page p
		j, _ = storage.GetDecks(l, p)
	}
	// write the JSON to the Response Writer
	w.Write(j)
}