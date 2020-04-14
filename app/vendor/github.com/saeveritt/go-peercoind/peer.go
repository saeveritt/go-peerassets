package bitcoind

type Peer struct {

	ID	int32	`json:"id"`

	// The ip address and port of the peer
	Addr string `json:"addr"`

	// Local address
	Addrlocal string `json:"addrlocal"`

	AddrBind  string  `json:"addrbind"`

	// The services
	Services string `json:"services"`

	RelayTxes	bool `json:"relaytxes"`

	// The time in seconds since epoch (Jan 1 1970 GMT) of the last send
	Lastsend uint64 `json:"lastsend"`

	// The time in seconds since epoch (Jan 1 1970 GMT) of the last receive
	Lastrecv uint64 `json:"lastrecv"`

	// The total bytes sent
	Bytessent uint64 `json:"bytessent"`

	// The total bytes received
	Bytesrecv uint64 `json:"bytesrecv"`

	// The connection time in seconds since epoch (Jan 1 1970 GMT)
	Conntime uint64 `json:"conntime"`

	TimeOffset int64	`json:"timeoffset"`

	// Ping time
	Pingtime float64 `json:"pingtime"`

	// Ping Wait
	Pingwait float64 `json:"pingwait"`

	// The peer version, such as 7001
	Version uint32 `json:"version"`

	// The string version
	Subver string `json:"subver"`

	// Inbound (true) or Outbound (false)
	Inbound bool `json:"inbound"`

	AddNode	bool	`json:"addnode"`

	//  The starting height (block) of the peer
	Startingheight int32 `json:"startingheight"`

	// The ban score (stats.nMisbehavior)
	Banscore int32 `json:"banscore"`

	SyncedHeaders	int32 `json:"synced_headers"`

	SyncedBlocks	int32 `json:"synced_blocks"`

	WhiteListed		bool	`json:"whitelisted"`

	// If sync node
	Syncnode bool `json:"syncnode,omitempty"`
}
