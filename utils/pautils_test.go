package utils

import (
	"encoding/hex"
	"testing"
)


func TestDeckParse(t *testing.T) {
	bytes := hex.EncodeToString([]byte{8,1,18, 8, 116, 101, 115, 116, 49, 49, 49, 52, 24, 2, 32, 7})
	opReturn := string(bytes)
	d := DeckParse(opReturn)
	if d.Version == 1 && d.Name == "test1114" && d.NumberOfDecimals == 2 && d.IssueMode == 7{
		t.Logf("%v",d)
	}else{
		t.Error("DeckParse Failed")
	}
}