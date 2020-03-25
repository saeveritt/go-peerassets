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

var db *bolt.DB
var err error

func must(err error){
	// Standard error handler. Will os.exit() on log.Fatal if error occurs.
	// Only use for db Connection
	if err != nil{
		log.Fatal(err)
	}
}

func Connect() (*bolt.DB,error){
	// Open the local db file and set it to db. This modifies the global variable so that
	// functions in this file can use it
	db, err = bolt.Open("storage/assets.db",0600,&bolt.Options{Timeout: 1 * time.Second})
	must(err)
	return db, err
}

func Close(){
	// if Connection to db is Open, Close it.
	if db.GoString() != "" {
		db.Close()
	}
}
func CreateBucket(bucket string,) {
	// Connect to local db
	Connect()
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
	Close()
}

func GetDecks(limit int, page int) ([]byte,error){
	// Connect to local db
	Connect()
	// Create a map, key < Deck ID >, Value < Deck Protobuf >
	res  := make(map[string]*protobuf.DeckSpawn)
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
				// Create an entry in the res map where key is string(<Deck ID>)
				// and the value is set to the Deck Object
				res[string(k)] = d
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
					// Create an entry in the res map where key is string(<Deck ID>)
					// and the value is set to the Deck Object
					res[string(k)] = d
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
	Close()
	// Will return nil if bucket is not found
	return j, nil
}



func GetAddress(address string, txtype string)([]byte,error){
	//make(map[string]string{"card":"Card-","deck":"Deck-"})
	if txtype == "deck"{
		res  := make(map[string]*protobuf.DeckSpawn)
	}
	else if txtype == "card"{
		res  := make(map[string]*protobuf.CardTransfer)
	}
	else{ return []byte, nil } 

	// Connect to local db
	Connect()

	db.View(func(tx *bolt.Tx) error{
		bucket := tx.Bucket([]byte(address))
		return bucket.ForEach( func(k []byte, v []byte) error{

			if string(k)[0:5] == "Deck-"{
				res[string(k[5:])] = protobuf.ParseDeck(v)
			}
			if string(k)[0:5] == "Card-"{
				res[string(k[5:])] = protobuf.ParseCard(v)
			}
			return nil
		})
		// Close local db connection
		Close()
		return nil
	})
	res := make(map[string]interface{})
	res["decks"] = resD
	res["cards"] = resC
	j,err := json.Marshal(res)
	if err != nil{
		return j, err
	}
	// Will return nil if bucket is not found
	return j, nil
}


func Put(bucket string,key string,value []byte) {
	var b *bolt.Bucket
	var err error
	db.Update(func(tx *bolt.Tx) error {
		b, err = tx.CreateBucketIfNotExists([]byte(bucket))
		must(err)
		must(b.Put([]byte(key), value))
		return nil
	})
}
func Get(bucket string, key string) []byte{

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
	return v
}

func PrefixScan(bucket string, keyPrefix string) map[string]string{

	M := make(map[string]string)
	count := 0

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket)).Cursor()

		if b != nil {

			prefix := []byte(keyPrefix)

			for k, v := b.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = b.Next() {
				M[string(k)] = string(v)
				if count > 20000 {
					break
				}
				count++
			}
		}
		return nil
	})
	return M
}