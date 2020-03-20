package networks

func Default() *NetParameters{
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
	p.Host = "localhost"
	p.Port = 9904
	p.User = "pothos"
	p.Password = "pothos"
	return &p
}