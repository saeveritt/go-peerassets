package protobuf

import (
	"log"
)

func ParseDeck(buf []byte) *DeckSpawn {
	Deck := &DeckSpawn{}
	err := Deck.XXX_Unmarshal(buf)
	if err != nil{
		log.Print(err)
	}
	return Deck
}

func ParseCard(buf []byte) *CardTransfer {
	Card := &CardTransfer{}
	err := Card.XXX_Unmarshal(buf)
	if err != nil{
		log.Print(err)
	}
	return Card
}

func ParseKey(buf []byte) *AddressCardKey {
	Key := &AddressCardKey{}
	err := Key.XXX_Unmarshal(buf)
	if err != nil{
		log.Print(err)
	}
	return Key
}