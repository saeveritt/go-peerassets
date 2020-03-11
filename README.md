# go-peerassets!
# Language
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
