package storage

import (
	"github.com/saeveritt/go-peerassets/utils"
)

func PutRootAsset(){
	txs := utils.RootTransactions()
	rawtxs := utils.RawTransactions(txs)
	for _, rawtx := range rawtxs{
		sender := utils.GetSender(rawtx)
		opReturn := utils.GetMetaData(rawtx)
		deck := utils.DeckParse(opReturn)
		proto, err := deck.XXX_Marshal(nil,false)
		must(err)
		if sender != "coinbase/coinstake" && len(proto) != 0{
			Put(rawtx.Txid,sender,proto)
		}
	}
	db.Close()
}
