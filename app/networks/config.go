package networks

import "github.com/saeveritt/go-peerassets/app/config"

var conf config.Config

func Default() *NetParameters {

	conf,_ = config.Open()
	// Set Default Network here
	// ex. "Peercoin-Testnet", "Peercoin"
	net := "Peercoin-Test"
	if net == "Peercoin-Test"{return PeercoinTestnet()}
	return Peercoin()
}

type NetParameters struct{
	Name			string
	Network 		string
	Address 		string
	WIF 			string
	WIFPrefix		string
	Fee 			float32
	Host 			string
	Port 			int
	User 			string
	Password		string
}


func Peercoin() *NetParameters {
	p := NetParameters{}
	p.Name = "Peercoin"
	p.Network = "PPC"
	p.Address = "PAprodbYvZqf4vjhef49aThB9rSZRxXsM6"
	p.WIF = "U624wXL6iT7XZ9qeHsrtPGEiU78V1YxDfwq75Mymd61Ch56w47KE"
	p.WIFPrefix = "b7" //0xb7
	p.Fee = 0.001
	p.Host = "localhost"
	p.Port = 9902
	p.User = "pothos"
	p.Password = "pothos"
	return &p
}

func PeercoinTestnet() *NetParameters {
	p := NetParameters{}
	p.Name = "Peercoin-Testnet"
	p.Network = "tPPC"
	p.Address = "miHhMLaMWubq4Wx6SdTEqZcUHEGp8RKMZt"
	p.WIF = "cTJVuFKuupqVjaQCFLtsJfG8NyEyHZ3vjCdistzitsD2ZapvwYZH"
	p.WIFPrefix = "ef" //0xef
	p.Fee = 0.001
	p.Host = "peercoind"
	p.Port = 19904
	p.User = "peercoind"
	p.Password = "peercoindrpc"
	return &p
}

func BitcoinCash() *NetParameters {
	// https://developer.bitcoin.com/mastering-bitcoin-cash/
	p := NetParameters{}
	p.Name = "BitcoinCash"
	p.Network = "BCH"
	// Leads with a p or q -- This needs to be defined at a later date
	p.Address = ""
	p.WIF = ""
	p.WIFPrefix = "80" //0x80
	p.Fee = 0.001 // Same as PeerAssets Fee
	p.Host = "localhost"
	p.Port = 8332
	p.User = "pothos"
	p.Password = "pothos"
	return &p
}

func BitcoinCashTestnet() *NetParameters {
	// https://developer.bitcoin.com/mastering-bitcoin-cash/
	p := NetParameters{}
	p.Name = "BitcoinCash-Testnet"
	p.Network = "tBCH"
	// Leads with a p or q -- This needs to be defined at a later date
	p.Address = ""
	p.WIF = ""
	p.WIFPrefix = "ef" //0xef
	p.Fee = 0.001 // Same as PeerAssets Fee
	p.Host = "localhost"
	p.Port = 18332
	p.User = "pothos"
	p.Password = "pothos"
	return &p
}

func Litecoin() *NetParameters {
	p := NetParameters{}
	p.Name = "Litecoin"
	p.Network = "LTC"
	p.Address = ""
	p.WIF = ""
	p.WIFPrefix = "32" //0x32
	p.Fee = 0.001 // Same as PeerAssets Fee but also same as LTC base fee
	p.Host = "localhost"
	p.Port = 9333
	p.User = "pothos"
	p.Password = "pothos"
	return &p
}

func LitecoinTestnet() *NetParameters {
	p := NetParameters{}
	p.Name = "Litecoin-Testnet"
	p.Network = "tLTC"
	p.Address = ""
	p.WIF = ""
	p.WIFPrefix = "3a" // 0x3a
	p.Fee = 0.001 // Same as PeerAssets Fee but also same as LTC base fee
	p.Host = "localhost"
	p.Port = 19335
	p.User = "pothos"
	p.Password = "pothos"
	return &p
}