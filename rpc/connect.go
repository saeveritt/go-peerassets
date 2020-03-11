package rpc

import (
	ppcd "github.com/saeveritt/go-peercoind"
	"go-peerassets/networks"

	"log"
)

func Connect(name string) (*ppcd.Bitcoind, *networks.NetParameters){
	var net *networks.NetParameters
	if name == "Peercoin"{
		net = networks.Peercoin()
	}
	if name == "Peercoin-Testnet"{
		net = networks.PeercoinTestnet()
	}


	cli, err := ppcd.New(net.Host,net.Port,net.User,net.Password,false)
	must(err)
	return cli, net
}

func must(err error){
	if err !=nil{
		log.Fatal(err)
	}
}

