package storage

import (
	"fmt"
	"log"
	"github.com/saeveritt/go-peerassets/utils"
	ppcd "github.com/saeveritt/go-peercoind"
)

var(
	//add to subscribed map the list of deck id's you wish to import
	subscribed = map[string]bool{
		"*": true,
	}
)

func PutRootAsset(){
	// Loads all valid assets registered to main p2th address registry
	Connect()
	txs := utils.RootTransactions()
	rawtxs := utils.RawTransactions(txs)
	i := 0 // Deck counter
	for _, rawtx := range rawtxs{
		if _,ok := subscribed["*"];!ok{continue}
		if _,ok := subscribed[rawtx.Txid];!ok {
			sender := utils.GetSender(rawtx)
			receiver := utils.GetReceiver(rawtx)
			opReturn := utils.GetMetaData(rawtx)
			deck := utils.DeckParse(opReturn)
			err := utils.ValidateDeckBasic(receiver, deck)
			if err != nil {
				//log.Print(err)
				continue
			}
			proto, err := deck.XXX_Marshal(nil, false)
			must(err)
			if sender != "coinbase/coinstake" && len(proto) != 0 {
				PutDeck(sender, rawtx)
				PutDeckProto(proto, rawtx)
				PutDeckCreator(sender, rawtx, proto)
				i++
				fmt.Printf("\r%d Decks Validated", i)
			}
		}
	}
	Close()
}

func PutDeck(sender string, rawtx ppcd.RawTransaction){
	//Import deck information into local db
	utils.ImportDeck(rawtx.Txid)
	// Bucket: Decks, Key: DeckSpawn ID, Value: Deck Owner
	Put("Decks",rawtx.Txid,[]byte(sender))
}
func PutDeckProto(proto []byte, rawtx ppcd.RawTransaction){
	//Import deck information into local db
	utils.ImportDeck(rawtx.Txid)
	// Bucket: Decks, Key: DeckSpawn ID, Value: Deck Owner
	Put("DecksProto",rawtx.Txid, proto)
}
func PutDeckCreator(sender string, rawtx ppcd.RawTransaction,proto []byte){
	// Bucket: <sender address>, Key: "Deck-" + <Deckspawn ID>, Value: <proto>
	Put(sender,"Deck-" + rawtx.Txid,proto)
}

func RescanBlockchain(txid string) {
	height := utils.GetBlockHeight(txid)
	log.Print("Scanning Transactions For: " + txid)
	utils.Scan(height)
}

func PutCards(deckid string){
	// Loads all valid assets registered to main p2th address registry
	txs := utils.RootTransactions()
	rawtxs := utils.RawTransactions(txs)
	for _, rawtx := range rawtxs{
		if _,ok := subscribed["*"];!ok{continue}
		if _,ok := subscribed[rawtx.Txid];!ok {
			sender := utils.GetSender(rawtx)
			receiver := utils.GetReceiver(rawtx)
			opReturn := utils.GetMetaData(rawtx)
			deck := utils.DeckParse(opReturn)
			err := utils.ValidateDeckBasic(receiver, deck)
			if err != nil {
				log.Print(err)
				continue
			}
			proto, err := deck.XXX_Marshal(nil, false)
			must(err)
			if sender != "coinbase/coinstake" && len(proto) != 0 {
				PutDeck(sender, rawtx)
				PutDeckCreator(sender, rawtx, proto)
			}
		}
	}
	db.Close()
}