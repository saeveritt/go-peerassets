package utils

import (
	"log"
	"github.com/saeveritt/go-peerassets/networks"
	"encoding/hex"
	"crypto/sha256"
)

func GetParams() *networks.NetParameters{
	network := networks.Default()
	if network == "Peercoin"{
		return networks.Peercoin()
	}
	if network == "Peercoin-Testnet"{
		return networks.PeercoinTestnet()
	}
	return networks.PeercoinTestnet()
}

func ToWIF(priv string) string{
	prefix := GetParams().WIFPrefix
	suffix := "01" // Default is compressed
	extended := prefix + priv + suffix
	log.Print(extended)
	hexBytes, err := hex.DecodeString(extended)
	if err != nil{
		log.Print(err)
	}
	checkSum := CheckSum(hexBytes)
	extended = extended + checkSum
	log.Print(extended)
	hexBytes, err = hex.DecodeString(extended)
	if err != nil {
		log.Fatal(err)
	}
	return Encode(hexBytes)
}

func Hash256(hexBytes []byte) []byte{
	// Performs sha256 twice
	hash1 := sha256.New()
	hash1.Write(hexBytes)
	hash2 := sha256.New()
	hash2.Write(hash1.Sum(nil))
	return hash2.Sum(nil)
}

func CheckSum(hexBytes []byte) string{
	// applies Base58 Encoding and returns first four bytes
	h256 := Hash256(hexBytes)
	str := hex.EncodeToString(h256)
	return str[0:8]
}
