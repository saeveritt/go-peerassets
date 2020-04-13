package storage

import (
	"github.com/boltdb/bolt"
	"github.com/saeveritt/go-peerassets/app/protobuf"
	"github.com/saeveritt/go-peerassets/app/utils"
	"sort"
)
type Card struct{
	sender		string
	receiver	string
	height		int32
	tx_index	int64
	card_index	int32
	amount		int64
}

func CalculateMulti(deckid string) map[string]int64{
	rawtx := utils.RawTransactions([]string{deckid})
	owner := utils.GetSender(rawtx[0])
	cards := GetCards(deckid)
	balances := make(map[string]int64)
	for _, card := range cards{
		if card.sender == owner{
			balances[card.receiver] += card.amount
		}
		if card.sender != owner{
			if balances[card.sender] >= card.amount{
				balances[card.receiver] += card.amount
				balances[card.sender] -= card.amount
			}
		}
	}
	return balances
}

func PutBalances(deckids []string){
	Connect()
	defer Close()
	for _,deckid := range deckids {
		var balances map[string]int64
		//deck, _ := utils.GetDeckInfo(deckid)
		//if protobuf.DeckSpawn_MODE_name[deck.IssueMode] == "MULTI" {
		balances = CalculateMulti(deckid)
		//}
		if len(balances) > 0 {
			for address, balance := range balances {
				bal := utils.Uint64Byte(uint64(balance))
				Put( "Balance-"+deckid, address, bal)
				}
		}
	}
}

func GetCards(deckid string) []Card {
	Connect()
	defer Close()
	var Cards []Card

	if err := db.View( func(tx *bolt.Tx) error{
		b := tx.Bucket( []byte(deckid))
		if b == nil { return nil}
		if err := b.ForEach( func(k []byte,v []byte) error{
			c := protobuf.ParseCard(v)
			Cards = append(Cards, Card{
				c.Sender,
				c.Receiver[0],
				c.BlockHeight[0],
				c.TxIndex[0],
				c.CardIndex[0],
				c.Amount[0],
			})
			return nil
		}); err != nil{
			return err
		}
		return nil
	}); err != nil{
		return []Card{}
	}
	result := SortCards(Cards)
	return result
}

func SortCards(Cards []Card) []Card {
	sort.Slice(Cards, func(p, q int) bool {
		if Cards[p].height == Cards[q].height{
			if Cards[p].tx_index == Cards[q].tx_index{
				return Cards[p].card_index < Cards[p].card_index
			}else{
				return Cards[p].tx_index < Cards[q].tx_index
			}
		}
		return Cards[p].height < Cards[q].height })
	return Cards
}
