package evmostx

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func Test_CreateTx(t *testing.T) {
	pub := "036af105fb761129e0e215664d0aac495f374c47910be160b671ae02be4ad02543"
	pubbite, _ := hex.DecodeString(pub)
	txsdk := EvmosSdk{
		"evmos1cuvqlx3cvmqahnhgunjpexseu5h7qzsgxkhmya",
		"evmos1cuvqlx3cvmqahnhgunjpexseu5h7qzsgxkhmya",
		"",
		"aevmos",
		"1",
		18,
		10 + 100,
		25000,
		15000,
		0,
		100,
		pubbite,
	}
	aa := txsdk.CreateTx()
	fmt.Println(aa)
}
