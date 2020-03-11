# go-peerassets!
This is a go implementation of the PeerAssets protocol defined by https://github.com/peerassets

##### What is this?
go-peerassets communicates with your local node using standard RPC commands. It aggregates
all relevant PeerAssets transaction on the chosen network and stores them in a local database.
This allows for asset tracking, state calculations, and much more! You can think of this as a 
type of blockexplorer for PeerAssets related transactions. 
##### What are the goals of this project?
This project is being worked on to ultimately create a backend blockchain parser that stores
data to be served via a REST API in the near future. This will allow for web-wallets to integrate
PeerAssets tracking and creation.
 
# Requirements
```
go get github.com/saeveritt/go-peercoind
go get github.com/golang/protobuf/proto
go get github.com/boltdb/bolt
```
# Configure
#####  networks/config.go
Set the default network type in the **Default()** function.
Currently supports "Peercoin" and "Peercoin-Testnet".
Make sure to modify your rpcuser and rpcpassword to match your local node 
inside of the NetParameters structure for the Network you've chosen.

```
func Default() string{
	// Set Default Network here
	// ex. "Peercoin-Testnet", "Peercoin"
	return "Peercoin-Testnet"
}
```


# Language
This section defines the language used in the following sections.

**deckID** = *Deck (asset) spawning transaction ID*\
**Address** = Relevant address to be queried\
**cAddress** = *Deck creator address*\
**dProto** = *PeerAssets deck protobuf*\
**cProto** = *PeerAsset card protobuf*

# Storage

**go-peerassets** stores data in *storage/assets.db*, which means all your files are automatically saved locally and are accessible offline if needed.

### BoltDB
[Bolt](https://github.com/boltdb/bolt) is a pure Go key/value store. The structure is as follows,


    Bucket [ deckID ]
	    Key [ cAddress ]
		    Value [ dProto ]

    Bucket [ Address ]
    	Key [ deckID-cardID]
    		Value [ cProto ]
