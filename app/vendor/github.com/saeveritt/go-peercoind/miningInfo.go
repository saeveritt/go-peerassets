package bitcoind

// A MiningInfo represents a mininginfo response
type MiningInfo struct {
	// The current block
	Blocks uint64 `json:"blocks"`

	// The last block weight
	CurrentBlockWeight uint64 `json:"currentblockweight"`

	// The last block transaction
	CurrentBlockTx uint64 `json:"currentblocktx"`

	// The current difficulty
	Difficulty float64 `json:"difficulty"`

	NetworkHashps float64 `json:"networkhashps"`

	NetworkGHashps float64 `json:"networkghps"`

	PooledTX	uint32	`json:"pooledtx"`

	Chain	string	`json:"chain"`

	Warnings	string	`json:"warnings"`
}
