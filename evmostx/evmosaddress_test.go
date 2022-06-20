package evmostx

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func Test_GetAddr(t *testing.T) {
	var addr_test EvmosAddr
	pub := "036af105fb761129e0e215664d0aac495f374c47910be160b671ae02be4ad02543"
	pubbite, _ := hex.DecodeString(pub)
	addr_test.GetAddr(pubbite)
	//eth_address = 0xC7180f9a3866c1DBcee8e4E41c9A19e52fe00A08
	//coms_addr = evmos1cuvqlx3cvmqahnhgunjpexseu5h7qzsgxkhmya
	fmt.Println(addr_test.EthAddr)
	fmt.Println(addr_test.ComsAddr)
}

func Test_EthAddrChange(t *testing.T) {
	var addr_test EvmosAddr
	addr_test.EthAddr = "0xC7180f9a3866c1DBcee8e4E41c9A19e52fe00A08"
	addr_test.EthAddrToComsAddr()
	fmt.Println(addr_test.ComsAddr)
}

func Test_ComosAddrChange(t *testing.T) {
	var addr_test EvmosAddr
	addr_test.ComsAddr = "evmos1e8ayvt0p8xm7hx9nvs32232ghlgaaq6sqlk3f8"
	addr_test.ComsAddrToEthAddr()
	fmt.Println(addr_test.EthAddr)
}
