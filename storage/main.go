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
	Connect()
	res  := make(map[string]*protobuf.DeckSpawn)
	db.View(func(tx *bolt.Tx) error{
		bucket := tx.Bucket([]byte("DecksProto"))
		return bucket.ForEach( func(k []byte, v []byte) error{
			d := protobuf.ParseDeck(v)
			res[string(k)] = d
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