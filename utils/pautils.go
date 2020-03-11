package utils

import (
	ppcd "github.com/saeveritt/go-peercoind"
	"log"
	"go-peerassets/protobuf"
	"go-peerassets/rpc"
	"go-peerassets/networks"
	"strings"
	"encoding/hex"
)
var cli *ppcd.Bitcoind
var net *networks.NetParameters

func init(){
	cli, net = rpc.Connect(networks.Default())
}

func must(err error){
	if err != nil{
		log.Fatal(err)
	}
}

func ImportRootP2TH() {
	resp, err := cli.ValidateAddress(net.Address)
	must(err)
	if resp.IsMine{return}
	must(cli.ImportPrivKey(net.WIF,net.Address,true))
}

func ImportDeck(txid string){
	_, err := cli.GetAddressesByAccount(txid)
	must(err)
	wif := ToWIF(txid)
	must(cli.ImportPrivKey(wif,txid,true))
	return
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

func GetMetaData(rawTx ppcd.RawTransaction) string{
	// Retrieves OP_RETURN Data from raw transaction
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