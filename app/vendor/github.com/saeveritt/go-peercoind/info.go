package bitcoind

// Info response for getblockchaininfo
type Info struct {

	Chain string `json:"chain"`

	Blocks uint32 `json:"blocks"`

	Headers uint32 `json:"headers"`

	BestBlockHash string `json:"bestblockhash"`

	Difficulty float64 `json:"difficulty"`

	MedianTime uint64 `json:"mediantime"`

	VerificationProgress float64 `json:"verificationprogress"`

	InitialBlockDownload bool `json:"initialblockdownload"`

	ChainWork string `json:"chainwork"`

	SizeOnDisk uint64 `json:"size_on_disk"`

	Pruned	bool	`json:"pruned"`

	SoftForks []*SoftForks `json:"softforks"`

}

// WalletInfo - wallet state info
// https://bitcoincore.org/en/doc/0.16.0/rpc/wallet/getwalletinfo/
type WalletInfo struct {
	WalletName            string  `json:"walletname"`
	WalletVersion         float64 `json:"walletversion"`
	Balance               float64 `json:"balance"`
	UnconfirmedBalance    float64 `json:"unconfirmed_balance"`
	ImmatureBalance       float64 `json:"immature_balance"`
	TxCount               int64   `json:"txcount"`
	KeyPoolOldest         int64   `json:"keypoololdest"`
	KeyPoolSize           int64   `json:"keypoolsize"`
	KeyPoolSizeHdInternal int64   `json:"keypoolsize_hd_internal"`
	UnlockedUntil         *int64  `json:"unlocked_until"`
	PaytxFee              float64 `json:"paytxfee"`
	HdMasterKeyID         *string `json:"hdmasterkeyid"`
}

type SoftForks struct{
	ID string `json:"id"`
	Version uint32 `json:"version"`
	Reject	*Reject	`json:"reject"`
}

type Reject struct{
	Status	bool `json:"status"`
}