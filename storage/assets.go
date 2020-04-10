package storage

import (
	"fmt"
	"github.com/saeveritt/go-peerassets/config"
	"github.com/saeveritt/go-peerassets/protobuf"
	"github.com/saeveritt/go-peerassets/utils"
	ppcd "github.com/saeveritt/go-peercoind"
	"log"
)


func PutRootAsset(){
	// Loads all valid assets registered to main p2th address registry
	txs := utils.RootTransactions()
	rawtxs := utils.RawTransactions(txs)
	i := 0 // Deck counter
	for _, rawtx := range rawtxs{
		height := utils.GetBlockHeight(rawtx.Txid)
		sender := utils.GetSender(rawtx)
		receiver := utils.GetReceiver(rawtx)
		opReturn := utils.GetMetaData(rawtx)
		deck,err := utils.DeckParse(opReturn)
		if err != nil{
			continue
		}
		err = utils.ValidateDeckBasic(receiver, deck)
		if err != nil {
			continue
		}
		proto, err := deck.XXX_Marshal(nil, false)
		if err != nil{
			continue
		}
		if sender != "coinbase/coinstake"{
			PutDeck(sender, rawtx)
			PutDeckProto(proto, rawtx)
			PutDeckCreator(sender, rawtx, proto)
			PutDeckHeight(height, rawtx)
			i++
			fmt.Printf("\r%d Decks Validated", i)
		}
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
			fmt.Printf("Rescanning Blockchain from Height: %v", lastScanned)
			current := utils.RescanBlockchain(lastScanned)
			fmt.Printf("Rescanned to Height: %v", current)
		}else{
			lowest := GetLowestBlock()
			fmt.Printf("Rescanning Blockchain from Height: %v", lowest)
			current := utils.RescanBlockchain(lowest)
			fmt.Printf("Rescanned to Height: %v", current)

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
		}

	}else {
		decks := data.Subscribed.Decks
		if len(decks) > 0 {
			PutCards(decks)
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
	for _, deckid := range deckids {
		cards := utils.GetCards(deckid)
		for _, card := range cards {
			ProcessDeckCardKeys(card)
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