package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/saeveritt/go-peerassets/app/rpc"
	"github.com/saeveritt/go-peerassets/app/storage"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var apiError = make(map[string]interface{})

func AgaveRouter() *mux.Router {
	// Create new router
	r := mux.NewRouter()
	// Set the subrouter path prefix
	api := r.PathPrefix("/v1").Subrouter()
	// Define the function handlers per route
	api.HandleFunc("/assets", getAssets).Methods(http.MethodGet)
	api.HandleFunc("/transactions", getTransactions).Methods(http.MethodGet)
	api.HandleFunc("/balances", getBalances).Methods(http.MethodGet)
	api.HandleFunc("/alert", walletNotify).Methods(http.MethodPost)
	api.HandleFunc("/send", sendTransaction).Methods(http.MethodPost)
	return r
}

func sendTransaction(w http.ResponseWriter, r *http.Request){
	logClient(r)
	rawtx,err:= ioutil.ReadAll( r.Body)
	if err != nil{
		return
	}
	_, err = strconv.ParseUint(string(rawtx), 16, 64)
	if err != nil {
		return
	}
	cli,_ := rpc.Connect("Peercoin-Testnet")
	resp, err := cli.SendRawTransaction( string(rawtx) )
	if err !=nil{
		log.Print(err)
	}
	log.Print(resp)
}


func walletNotify(w http.ResponseWriter, r *http.Request){
	logClient(r)
	r.ParseForm()
	txid := r.Form.Get("txid")
	cli,_ := rpc.Connect("Peercoin-Testnet")
	info,_ := cli.GetInfo()
	if info.Blocks == info.Headers {
		log.Printf("walletnotify: %v", txid)
	}
}


func getBalances( w http.ResponseWriter, r *http.Request ){
	logClient(r)
	w.Header().Set( "Content-Type","application/json")

	var j []byte

	address := r.URL.Query().Get("address")
	//deck := r.URL.Query().Get("deck")

	if len(address) == 34{
		j, _ = json.Marshal( storage.GetUserBalances(address) )
		w.Write(j)
	}
}


func getTransactions( w http.ResponseWriter, r *http.Request){
	logClient(r)
	w.Header().Set("Content-Type","application/json")
	// Create empty byte array which will store the JSON Response
	var j []byte
	// Assign variables to the values passed in the GET request
	var address = r.URL.Query().Get("address")
	var deck = r.URL.Query().Get("deck")
	var txType = r.URL.Query().Get("type")
	var limit = r.URL.Query().Get("limit")
	var page = r.URL.Query().Get("page")
	l,p,_ := pageLimit( limit, page)
	// If address is not empty
	if deck != "" && len(deck) == 64{
		j,_ = storage.GetDeckTransactions(deck,l,p)
		w.Write(j)
		return
	} else if len(address) == 34 && txType != ""{
		// Check for address in storage. Each address has its own dedicated bucket.
		j, _ = storage.GetAddress(address,txType,l,p)
		// if there was an error writing the JSON byte array,it will send empty array
		// else it sends a JSON byte array Response with the results
		//w.WriteHeader(200)
		w.Write(j)
		return
	}else{
		// if address argument is empty, return empty byte array
		apiError["error"] = "Invalid arguments"
		j, _ = json.Marshal(apiError)
		//w.WriteHeader(400)
		w.Write( j )
		return
	}
}

func getAssets(w http.ResponseWriter, r *http.Request) {
	logClient(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	// Assign variables to the values passed in the GET request
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	l,p,_ := pageLimit(limit, page)
	// Get Asset
	j,_ := storage.GetDecks(l, p)
	// write the JSON to the Response Writer
	w.Write(j)
}


////////////////////////////////////////////
///////////// Utilities ///////////////////
//////////////////////////////////////////

func pageLimit(limit string, page string) (int,int, error){
	if !isDigit(page) && !isDigit(limit){
		return 0,0,nil
	}
	l, _ := strconv.Atoi(limit)
	p, _ := strconv.Atoi(page)

	if  (p > 0 && l > 0){
		return l, p, nil
	}
	if p > 0 && l == 0 {
		return 10, p, nil
	}
	if p == 0 && l > 0 {
		return l, 1, nil
	}
	if p==0 && l ==0 {
		return 10,1,nil
	}
	return 0,0, nil
}

func isDigit(s string) bool{
	b := true
	if s == ""{ return b}
	for _, c := range s {
		if c < '0' || c > '9' {
			b = false
			break
		}
	}
	return b
}

func logClient(r *http.Request){
	log.Print("{'Client IP': '" + r.RemoteAddr + "', 'URI': '" + r.RequestURI + "'}")
}