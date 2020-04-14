package bitcoind

// A ScriptSig represents a scriptsyg
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// Vin represent an IN value
type Vin struct {
	Coinbase  	string    	`json:"coinbase"`
	Txid      	string    	`json:"txid"`
	Vout      	int       	`json:"vout"`
	ScriptSig 	ScriptSig 	`json:"scriptSig"`
	Sequence  	uint32    	`json:"sequence"`
	TxinWitness	[]string	`json:"txinwitness,omitempty"`
}

type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int      `json:"reqSigs,omitempty"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses,omitempty"`
}

// Vout represent an OUT value
type Vout struct {
	Value        float64      `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

// RawTx represents a raw transaction
type RawTransaction struct {
	InActiveChain	bool	`json:"in_active_chain,omitempty"`
	Hex           	string 	`json:"hex"`
	Txid          	string 	`json:"txid"`
	Hash 			string	`json:"hash"`
	Size 			int64	`json:"size"`
	VSize			int64	`json:"vsize"`
	Version       	uint32 	`json:"version"`
	LockTime      	uint32 	`json:"locktime"`
	Vin           	[]Vin  	`json:"vin"`
	Vout          	[]Vout 	`json:"vout"`
	BlockHash     	string 	`json:"blockhash,omitempty"`
	Confirmations 	uint64 	`json:"confirmations,omitempty"`
	Time          	int64  	`json:"time,omitempty"`
	Blocktime     	int64  	`json:"blocktime,omitempty"`
}

// TransactionDetails represents details about a transaction
type TransactionDetails struct {
	Account  	string  `json:"account"`
	Address  	string  `json:"address,omitempty"`
	Category 	string  `json:"category"`
	Amount   	float64 `json:"amount"`
	Fee      	float64 `json:"fee,omitempty"`
	Label	 	string	`json:"label,omitempty"`
	Vout	 	int32	`json:"vout"`
	Abandoned	bool 	`json:"abandoned,omitempty"`
}

// Transaction represents a transaction
type Transaction struct {
	Amount          float64              `json:"amount"`
	Fee             float64              `json:"fee,omitempty"`
	Confirmations   int64                `json:"confirmations"`
	BlockHash       string               `json:"blockhash"`
	BlockIndex      int64                `json:"blockindex"`
	BlockTime       int64                `json:"blocktime"`
	TxID            string               `json:"txid"`
	Time			int64				 `json:"time"`
	TimeReceived	int64				 `json:"timereceived"`
	Details         []TransactionDetails `json:"details,omitempty"`
	Hex             string               `json:"hex,omitempty"`
}

// UTransactionOut represents a unspent transaction out (UTXO)
type UTransactionOut struct {
	Bestblock     string       `json:"bestblock"`
	Confirmations uint32       `json:"confirmations"`
	Value         float64      `json:"value"`
	ScriptPubKey  ScriptPubKey `json:"scriptPubKey"`
	Version       uint32       `json:"version"`
	Coinbase      bool         `json:"coinbase"`
}

// TransactionOutSet represents statistics about the unspent transaction output database
type TransactionOutSet struct {
	Height          uint32  `json:"height"`
	Bestblock       string  `json:"bestblock"`
	Transactions    float64 `json:"transactions"`
	TxOuts          float64 `json:"txouts"`
	Bogosize		uint64	`json:"bogosize"`
	HashSerialized2 string  `json:"hash_serialized_2"`
	DiskSize		uint64	`json:"disk_size"`
	TotalAmount     float64 `json:"total_amount"`
}
