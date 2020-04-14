package storage

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/saeveritt/go-peerassets/app/config"
	"github.com/saeveritt/go-peerassets/app/protobuf"
	"github.com/saeveritt/go-peerassets/app/utils"
	ppcd "github.com/saeveritt/go-peercoind"
	"log"
)


func PutRootAsset(){
	// Loads all valid assets registered to main p2th address registry
	txs := utils.RootTransactions()
	rawtxs := utils.RawTransactions(txs)
	i := 0 // Deck counter
	for _, rawtx := range rawtxs{
			err := ImportDeck(rawtx.Txid)
			if err != nil{ continue }
			i++
			fmt.Printf("\r%d Decks Validated", i)
		}
}

func ImportSubscribed() error{
	data, err := config.Open()
	if err != nil{
		log.Fatal(err)
	}
	if data.Subscribed.All {
		decks := GetAllDecks()
		if len(decks)  > 0 {
			utils.ImportDecks(decks)
		}
	}else{
		utils.ImportDecks(data.Subscribed.Decks)
		}
		lastScanned := GetScanHeight()
		if lastScanned != 0 {
			log.Printf("Rescanning Blockchain from Height: %v", lastScanned)
			current := utils.RescanBlockchain(lastScanned)
			log.Printf("Rescanned to Height: %v", current)
		}else{
			lowest := GetLowestBlock()
			log.Printf("Rescanning Blockchain from Height: %v", lowest)
			current := utils.RescanBlockchain(lowest)
			log.Printf("Rescanned to Height: %v", current)

		}
	return nil
}

func ImportSubscribedCards(){
	data, err := config.Open()
	if err != nil{
		log.Fatal(err)
	}
	if data.Subscribed.All {
		decks := GetAllDecks()
		if len(decks) > 0 {
			PutCards(decks)
			PutBalances(decks)

		}

	}else {
		decks := data.Subscribed.Decks
		if len(decks) > 0 {
			PutCards(decks)
			PutBalances(decks)
		}
	}
}

func PutScanHeight(height uint64){
	bHeight := utils.Uint64Byte(height)
	Put("DecksHeight","LastScanned",bHeight)
}

func GetScanHeight() uint64{
	scanHeight := Get("DecksHeight","LastScanned")
	if scanHeight == nil{
		return uint64(0)
	}
	return utils.ByteUint64(scanHeight)
}

func ImportDeck(txid string) error{
	rawtx := utils.RawTransactions([]string{txid})[0]
	sender := utils.GetSender(rawtx)
	if sender != "coinbase/coinstake" {
		height := utils.GetBlockHeight(rawtx.Txid)
		receiver := utils.GetReceiver(rawtx)
		opReturn := utils.GetMetaData(rawtx)
		deck, err := utils.DeckParse(opReturn)

		if err != nil {
			return err
		}
		err = utils.ValidateDeckBasic(receiver, deck)
		if err != nil{return err}
		proto, err := deck.XXX_Marshal(nil, false)
		if err != nil {
			return err
		}
		PutDeck(sender, rawtx)
		PutDeckProto(proto, rawtx)
		PutDeckCreator(sender, rawtx, proto)
		PutDeckHeight(height, rawtx)
	}else{ return errors.New("Coinbase Transaction")}
	return nil
}

func PutDeck(sender string, rawtx ppcd.RawTransaction){
	//Import deck information into local db
	// Bucket: Decks, Key: DeckSpawn ID, Value: Deck Owner
	Put("Decks",rawtx.Txid,[]byte(sender))
}
func PutDeckHeight(height uint64, rawtx ppcd.RawTransaction){
	//Import deck information into local db
	bHeight := utils.Uint64Byte(height)
	// Bucket: Decks, Key: DeckSpawn ID, Value: Deck Owner
	Put("DecksHeight",rawtx.Txid,bHeight)
}

func PutDeckProto(proto []byte, rawtx ppcd.RawTransaction){
	//Import deck information into local db
	// Bucket: Decks, Key: DeckSpawn ID, Value: Deck Owner
	Put( "DecksProto",rawtx.Txid, proto)
}
func PutDeckCreator(sender string, rawtx ppcd.RawTransaction,proto []byte){
	// Bucket: <sender address>, Key: "Deck-" + <Deckspawn ID>, Value: <proto>
	deck := protobuf.AddressCardKey{Type: 0x01, DeckId: rawtx.Txid}
	deckKey,_ := deck.XXX_Marshal(nil,false)
	PutByte(sender,deckKey,proto)
}


func PutCards(deckids []string){
	// Loads all valid assets registered to main p2th address registry
	accounts := make(map[string]map[string]bool)
	for _, deckid := range deckids {
		cards := utils.GetCards(deckid)
		for _, card := range cards {
			ProcessDeckCardKeys(card)
			if accounts[card.Receiver[0]] == nil{
				accounts[card.Receiver[0]] = make(map[string]bool)
			}
			accounts[card.Receiver[0]][deckid] = true
		}
	}
	for address, deckids := range accounts{
		for deckid, _ := range deckids {
			baseKey := protobuf.AddressCardKey{
				Type: 0x01,
				DeckId: deckid,
			}
			accountKey,_ := baseKey.XXX_Marshal(nil,false)
			PutByte(address,accountKey,[]byte("true"))
		}
	}
}

func ProcessDeckCardKeys(card *protobuf.CardTransfer){
	baseKey := protobuf.AddressCardKey{
		Type:                 0x02,
		CardType: 			  0x01,
		DeckId:               card.DeckId,
		BlockHeight:          card.BlockHeight[0],
		TxIndex:              card.TxIndex[0],
		CardIndex:            card.CardIndex[0],
	}
	sendKey,_ := baseKey.XXX_Marshal(nil,false)
	baseKey.CardType = 0x02
	receiveKey,_ := baseKey.XXX_Marshal(nil,false)

	baseKey.Type = 0x00
	baseKey.CardType = 0x00
	baseKey.DeckId = ""
	deckKey,_ := baseKey.XXX_Marshal(nil,false)
	proto,_ := card.XXX_Marshal(nil,false)
	PutByte(card.Sender,sendKey, proto)
	PutByte(card.Receiver[0],receiveKey, proto)
	PutByte(card.DeckId, deckKey,proto)
}


func GetUserBalances(address string) map[string]int64{
	Connect()
	defer Close()
	balances := make( map[string]int64 )
	deckType := protobuf.AddressCardKey{Type: 0x01}
	db.View( func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte(address))
		c := b.Cursor()
		prefix,_ := deckType.XXX_Marshal(nil,false)
		for k,_ := c.Seek(prefix); bytes.HasPrefix(k, prefix);k,_ = c.Next(){
			key := protobuf.ParseKey(k)
			balances[key.DeckId] = 0
		}
		for k, _ := range balances {
			b = tx.Bucket([]byte("Balance-" + k))
			if b != nil {
				byteAmount := b.Get([]byte( address ))
				if byteAmount != nil {
					amount := utils.ByteUint64(byteAmount)
					balances[k] = int64(amount)
				}
			}
		}
	return nil
	})
	return balances
}