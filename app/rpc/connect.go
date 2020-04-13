package rpc

import (
	"github.com/saeveritt/go-peerassets/app/networks"
	ppcd "github.com/saeveritt/go-peercoind"

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

