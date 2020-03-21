package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/saeveritt/go-peerassets/networks"
	"log"
)

func GetParams() *networks.NetParameters{

	return networks.Default()
}

func ToWIF(priv string) string{
	prefix := GetParams().WIFPrefix
	suffix := "01" // Default is compressed
	extended := prefix + priv + suffix
	hexBytes, err := hex.DecodeString(extended)
	if err != nil{
		log.Print(err)
	}
	checkSum := CheckSum(hexBytes)
	extended = extended + checkSum
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
