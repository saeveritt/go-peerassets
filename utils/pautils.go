package utils

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/saeveritt/go-peerassets/networks"
	"github.com/saeveritt/go-peerassets/protobuf"
	"github.com/saeveritt/go-peerassets/rpc"
	ppcd "github.com/saeveritt/go-peercoind"
	"log"
	"strings"
)
var cli *ppcd.Bitcoind
var net *networks.NetParameters

func init(){
	cli, net = rpc.Connect(networks.Default().Name)
}

func must(err error){
	if err != nil{
		log.Print(err)
	}
}

func ImportRootP2TH() {
	resp, err := cli.ValidateAddress(net.Address)
	must(err)
	if resp.IsMine{
		fmt.Println("P2TH previously imported. Scanning for assets...")
		return}
	// This will load the P2TH Main Registry with Account Name set to <Address>
	must(cli.ImportPrivKey(net.WIF,net.Address,false))
}

func ImportDeck(deckid string){
	_, err := cli.GetAddressesByAccount(deckid)
	must(err)
	wif := ToWIF(deckid)
	// Imports deck and sets the Account name to the deckid (txid)
	must(cli.ImportPrivKey(wif,deckid,true))
	return
}

func GetCards(deckid string) []*protobuf.CardTransfer{
	resp, err := cli.ListTransactions(deckid,99999999,0)
	deck := GetDeckInfo(deckid)
	txs := make([]string,len(resp))
	txIndex := make(map[string]int64)
	txBlock := make(map[string]int32)
	txSender := make(map[string]string)
	must(err)
	for _,tx := range resp{
		    block,_ := cli.GetBlockheader(tx.BlockHash)
			txIndex[tx.TxID] = tx.BlockIndex
			txBlock[tx.TxID] = int32(block.Height)
			txs = append(txs, tx.TxID)
	}
	rawtxs := RawTransactions(txs)
	cards := []*protobuf.CardTransfer{}

	for _,rawtx := range rawtxs{
		data := GetMetaData(rawtx)
		if len(data) < 1{continue}
		card := CardParse(data)
		//cardReceiver := []*protobuf.CardTransfer{}
		if card.NumberOfDecimals != deck.NumberOfDecimals {continue}
		txSender[rawtx.Txid] = GetSender(rawtx)
		for i, amount := range card.Amount{
			card.Amount = []int64{amount}
			card.Sender = txSender[rawtx.Txid]
			receiver := rawtx.Vout[2+i].ScriptPubKey.Addresses[0]
			card.DeckId = deckid
			card.CardId = rawtx.Txid
			card.Receiver = []string{receiver}
			card.BlockHeight = []int32{txBlock[rawtx.Txid]}
			card.TxIndex = []int64{txIndex[rawtx.Txid]}
			card.CardIndex = []int32{int32(i)}
			log.Print(card)
		}
		cards = append(cards,card)

	}
	//for _, card := range cards{
	//	index : = 2
	//
	//}
	return cards


}


func Scan(height uint64){
	log.Print("Rescanning Blockchain..")
	err := cli.RescanBlockchain(height)
	must(err)
}

func GetBlockHeight(txid string) uint64{
	rawtx, err := cli.GetRawTransaction(txid,true)
	block, err := cli.GetBlock(rawtx.BlockHash)
	must(err)
	return block.Height
}

func RootTransactions() []string{
	var txs []string
	resp, err := cli.ListTransactions(net.Address,99999999,0)
	must(err)
	for _, tx := range resp{
		txs = append(txs, tx.TxID)
	}
	return txs
}

func RawTransactions(transactions []string) []ppcd.RawTransaction{
	var rawTxs []ppcd.RawTransaction
	for _, txid := range transactions{
		rawTx, err := cli.GetRawTransaction(txid, true)
		must(err)
		rawTxs = append(rawTxs, rawTx)
	}
	return rawTxs
}

func GetSender(childTx ppcd.RawTransaction ) (sender string){
	// Define that vin[0] is from the targeted transaction sender
	childVin := childTx.Vin[0]
	// Acquire the vout index of the childTx vin to trace back to sender
	index := childVin.Vout
	// Acquire the Raw Parent transaction
	parentTx, err := cli.GetRawTransaction(childTx.Vin[0].Txid,true)
	if err != nil{
		return "coinbase/coinstake"
	}
	// Retrieve the first address in the scriptpubkey
	sender = parentTx.Vout[index].ScriptPubKey.Addresses[0]
	return
}
func GetCardReceiver(rawtx ppcd.RawTransaction, n int) string{
	// n represents the Vout index number
	receiver := rawtx.Vout[n].ScriptPubKey.Addresses[0]
	return receiver
}


func GetReceiver(rawtx ppcd.RawTransaction ) (receiver string){
	// Define that vin[0] is from the targeted transaction sender
	receiver = ""
	if len(rawtx.Vout[0].ScriptPubKey.Addresses) > 0 {
		receiver = rawtx.Vout[0].ScriptPubKey.Addresses[0]
	}
	return
}

func GetMetaData(rawTx ppcd.RawTransaction) string{
	// Retrieves OP_RETURN Data from raw transaction
	if len(rawTx.Vout) < 1 {return ""}
	asm := rawTx.Vout[1].ScriptPubKey.Asm
	// Seperates into array of strings by spaces
	s := strings.Fields(asm)
	if s[0] != "OP_RETURN"{return ""}
	if len(s) <= 1{return ""}
	return s[1]
}

func DeckParse(opReturn string) (Deck *protobuf.DeckSpawn){
	// convert hex string to bytes
	hexBytes, err := hex.DecodeString(opReturn)
	must(err)
	// Returns Unmarshalled bytes as Deck
	Deck = protobuf.ParseDeck(hexBytes)
	return
}

func CardParse(opReturn string) (Card *protobuf.CardTransfer){
	// convert hex string to bytes
	hexBytes, err := hex.DecodeString(opReturn)
	must(err)
	// Returns Unmarshalled bytes as Deck
	Card = protobuf.ParseCard(hexBytes)
	return
}

func ValidateDeckBasic(receiver string, deck *protobuf.DeckSpawn) error{
	if deck.Name == ""{
		return errors.New("deck name cannot be empty")
	}
	if _,ok := protobuf.DeckSpawn_MODE_name[deck.IssueMode]; !ok{
		return errors.New("issue mode not valid/supported")
	}
	if receiver != net.Address{
		return errors.New("deckspawn must be sent to main p2th address to be valid")
	}
	return nil
}

func GetDeckInfo(deckid string) *protobuf.DeckSpawn{
	rawtx,err := cli.GetRawTransaction(deckid,true)
	must(err)
	meta := GetMetaData(rawtx)
	deck := DeckParse(meta)
	return deck



}

//func ValidateCardBasic(deckid string, rawtx ppcd.RawTransaction, card *protobuf.CardTransfer){
//	sender:= GetSender(rawtx)
//	receiver := GetReceiver(rawtx)
//	owner := string(storage.Get("Decks",deckid))
//	if sender != owner {
//		// Check sender balance
//		balance := storage.Get(sender,"Balance-"+deckid)
//		for _, val := range card.Amount {
//			if byteUint64(balance) < val {
//
//			}
//		}
//	}
//	if sender == owner{
//
//		storage.Put(receiver,"Balance-"+deckid,uint64Byte(card.Amount))
//	}
//
//}


func uint64Byte( value uint64) (b []byte){
	b = make([]byte,8)
	binary.BigEndian.PutUint64(b, value)
	return
}

func byteUint64( value []byte) (val uint64){
	val = binary.BigEndian.Uint64(value)
	return
}