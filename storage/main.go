package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/saeveritt/go-peerassets/protobuf"
	"log"
	"time"
)

//var db *bolt.DB
//var err error

func must(err error){
	// Standard error handler. Will os.exit() on log.Fatal if error occurs.
	// Only use for db Connection
	if err != nil{
		log.Print(err)
	}
}

func Connect() (*bolt.DB,error){
	// Open the local db file and set it to db. This modifies the global variable so that
	// functions in this file can use it
	db, err := bolt.Open("storage/assets.db",0600,&bolt.Options{Timeout: 1 * time.Second})
	must(err)
	return db, err
}

func CreateBucket(bucket string,) {
	// Connect to local db
	db,_ := Connect()
	// Use db.Update and pass a function with a bolt.TX
	db.Update(func(tx *bolt.Tx) error {
		// Create the Bucket if it does not exist
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		// Error Handling, Non-Standard
		if err != nil {
			// Return Bucket Exists
			return fmt.Errorf("Bucket Exists: %s", err)
		}
		// Return nil for Error
		return nil
	})
	// Close local db connection
	db.Close()
}

func GetDecks(limit int, page int) ([]byte,error){
	// Connect to local db
	db,_ :=Connect()
	// Create a map, key < Deck ID >, Value < Deck Protobuf >
	var res []interface{}
	bucketName := "DecksProto"
	// Create a View to query the local database
	// Input to db.View is a Function that will iterate and grab Keys and Values
	if limit == 0 && page == 0 {
		db.View(func(tx *bolt.Tx) error{
			// DecksProto bucket contains Key: <Deck ID>, Value: < Deck Protobuf >
			bucket := tx.Bucket([]byte(bucketName))
			// Setup a return which will iterate through each Key, Value in Bucket and append to map
			return bucket.ForEach( func(k []byte, v []byte) error{
				// Grab the <Deck Protobuf> and Parse it into a Deck Object
				d := protobuf.ParseDeck(v)
				j := FormatDeckResponse(string(k), d)
				// Create an entry in the res map where key is string(<Deck ID>)
				// and the value is set to the Deck Object
				res = append(res,j)
				// Return nil because there were no errors iterating through the bucket
				return nil
			})
			// Return nil after apending to map
			return nil
		})
	} else{
		// This View is for processing requests with page and limit arguments
		db.View(func(tx *bolt.Tx) error{
			// DecksProto bucket contains Key: <Deck ID>, Value: < Deck Protobuf >
			bucket := tx.Bucket([]byte(bucketName))
			// Create a counter to restrict which entries from the db we append to Response
			n := -1
			return bucket.ForEach( func(k []byte, v []byte) error{
				// Increase counter
				n++
				// Add entries to Response map that are within the terms defined by limit and page
				if n >= page*limit-limit && n < page * limit {
					// Grab the <Deck Protobuf> and Parse it into a Deck Object
					d := protobuf.ParseDeck(v)
					j := FormatDeckResponse(string(k), d)
					// Create an entry in the res map where key is string(<Deck ID>)
					// and the value is set to the Deck Object
					res = append(res,j)
				}
				// Return nil because there were no errors iterating through the bucket
				return nil
			})
			// Return nil after apending to map
			return nil
		})
	}
	// Create a variable to store the JSON byte array which will be used to write the Response
	j,err := json.Marshal(res)

	if err != nil{
		return j, err
	}
	// Close local db connection
	db.Close()
	// Will return nil if bucket is not found
	return j, nil
}

func GetAddress(address string, txType string, limit int, page int)([]byte,error){
	var j []byte
	// Make sure that type is either deck or card
	if (txType != "deck" && txType != "card") || len(address) != 34{
		return j, nil
	}
	// Create an empty map for the Response
	var res []interface{}
	// Connect to local db
	db, _ := Connect()
	prefix := map[string][]byte{"card":[]byte("Card-"),"deck":[]byte("Deck-")}
	log.Print(prefix)
	// This View is for processing requests with page and limit arguments
	db.View(func(tx *bolt.Tx) error{
		fmt.Print(address)
		c := tx.Bucket([]byte(address)).Cursor()
		n := -1
		// Use Seek to iterate through the bucket based on specified prefix
		for k,v := c.Seek(prefix[txType]); k != nil && bytes.HasPrefix(k,prefix[txType]); k,v = c.Next(){
			n++
			// Limit the output results based on GET arguments passed in request
			if n >= page*limit-limit && n < page * limit {
				// Handle the parsing based on what "type" argument was passed
				switch txType {
					case "card":
						c := protobuf.ParseCard(v)
						send := address == c.Sender
						j := FormatCardResponse(c,send)
						res = append(res,j)
					case "deck":
						d := protobuf.ParseDeck(v)
						j := FormatDeckResponse(string(k[5:]), d)
						res = append(res, j)
					}
				}
			}
		return nil
		})
		// Close local db connection
		db.Close()

	j,err := json.Marshal(res)
	if err != nil{
		return j, err
	}
	// Will return nil if bucket is not found
	return j, nil
}

func FormatCardResponse( c *protobuf.CardTransfer, send bool) map[string]interface{}{
	j := make(map[string]interface{})
	if send { j["type"] = "send"}else{ j["type"] = "receive"}
	j["deck_id"] = c.DeckId
	j["card_id"] = c.CardId
	j["block_height"] = c.BlockHeight[0]
	j["tx_index"] = c.TxIndex[0]
	j["card_index"] = c.CardIndex[0]
	j["sender"] = c.Sender
	j["receiver"] = c.Receiver[0]
	j["amount"] = c.Amount[0]
	j["data"] = c.AssetSpecificData
	return j
}

func FormatDeckResponse(deckid string, d *protobuf.DeckSpawn) map[string]interface{}{
	j := make(map[string]interface{})
	j["txid"] = deckid
	j["name"] = d.Name
	j["version"] = d.Version
	j["mode"] = protobuf.DeckSpawn_MODE_name[d.IssueMode]
	j["decimals"] = d.NumberOfDecimals
	j["data"] = d.AssetSpecificData
	j["fee"] = d.Fee
	return j

}


func Put(bucket string,key string,value []byte) {
	db, _ := Connect()
	var b *bolt.Bucket
	var err error
	db.Update(func(tx *bolt.Tx) error {
		b, err = tx.CreateBucketIfNotExists([]byte(bucket))
		must(err)
		if err != nil { return nil}
		must(b.Put([]byte(key), value))
		return nil
	})
	db.Close()
}
func Get(bucket string, key string) []byte{
	db,_ := Connect()
	var v []byte
	db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucket))

		if b != nil {

			v = b.Get([]byte(key))

		} else{

			v = nil
		}
		return nil

	})
	db.Close()
	return v
}