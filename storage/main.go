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

func Connect() (*bolt.DB,error){
	db, err = bolt.Open("storage/assets.db",0600,&bolt.Options{Timeout: 1 * time.Second})
	must(err)
	return db, err
}

func Close(){
	if db.GoString() != "" {
		db.Close()
	}
}
func CreateBucket(bucket string,) {

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func GetDecks() ([]byte,error){
	// Connect to local db
	Connect()
	// Create a map, key < Deck ID >, Value < Deck Protobuf >
	res  := make(map[string]*protobuf.DeckSpawn)
	// Create a View to query the local database
	// Input to db.View is a Function that will iterate and grab Keys and Values
	db.View(func(tx *bolt.Tx) error{
		// DecksProto bucket contains Key: <Deck ID>, Value: < Deck Protobuf >
		bucket := tx.Bucket([]byte("DecksProto"))
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

func GetDecksPages(limit int, page int) ([]byte,error){
	// Connect to local db
	Connect()
	// Create a map, key < Deck ID >, Value < Deck Protobuf >
	res  := make(map[string]*protobuf.DeckSpawn)
	// Create a View to query the local database
	// Input to db.View is a Function that will iterate and grab Keys and Values
	db.View(func(tx *bolt.Tx) error{
		bucket := tx.Bucket([]byte("DecksProto"))
		n := -1 //counter
		return bucket.ForEach( func(k []byte, v []byte) error{
			n++
			if n >= page*limit-limit && n < page * limit {
				d := protobuf.ParseDeck(v)
				res[string(k)] = d
			}
			return nil
		})

		return nil
	})
	j,err := json.Marshal(res)
	if err != nil{
		return j, err
	}
	Close()
	// Will return nil if bucket is not found
	return j, nil
}

func GetAddress(address string)([]byte,error){
	Connect()
	resD  := make(map[string]*protobuf.DeckSpawn)
	resC  := make(map[string]*protobuf.CardTransfer)

	db.View(func(tx *bolt.Tx) error{
		bucket := tx.Bucket([]byte(address))
		return bucket.ForEach( func(k []byte, v []byte) error{

			if string(k)[0:5] == "Deck-"{
				resD[string(k[5:])] = protobuf.ParseDeck(v)
			}
			if string(k)[0:5] == "Card-"{
				resC[string(k[5:])] = protobuf.ParseCard(v)
			}
			return nil
		})
		return nil
	})
	res := make(map[string]interface{})
	res["decks"] = resD
	res["cards"] = resC
	j,err := json.Marshal(res)
	if err != nil{
		return j, err
	}
	Close()
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


func must(err error){
	if err != nil{
		log.Fatal(err)
	}
}