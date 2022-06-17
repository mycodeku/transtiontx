package evmostx

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	sdk "github.com/mycodeku/transtionhelper/types"
)

type EvmosAddr struct {
	EthAddr  string
	ComsAddr string
}

func (self *EvmosAddr) GetAddr(pubkey []byte) error {
	pub, _ := DecompressPubKey(pubkey)
	ethaddress := crypto.PubkeyToAddress(*pub).String()

	ublicKeyB, err := hex.DecodeString(ethaddress[2:])
	if err != nil {
		return err
	}
	cosmosaddress := fmt.Sprintf("%s", sdk.AccAddress(ublicKeyB))

	self.ComsAddr = cosmosaddress
	self.EthAddr = ethaddress

	return nil
}
func (self *EvmosAddr) EthAddrToComsAddr() error {
	ublicKeyB, err := hex.DecodeString(self.EthAddr[2:])
	if err != nil {
		return err
	}
	cosmosaddress := fmt.Sprintf("%s", sdk.AccAddress(ublicKeyB))

	self.ComsAddr = cosmosaddress
	return err
}
func (self *EvmosAddr) ComsAddrToEthAddr() error {
	addrbite, err := sdk.AccAddressFromBech32(self.ComsAddr)
	if err != nil {
		return err
	}

	self.ComsAddr = "0x" + hex.EncodeToString(addrbite.Bytes())
	return nil
}

func DecompressPubKey(pubKey []byte) (*ecdsa.PublicKey, error) {
	x, y := secp256k1.DecompressPubkey(pubKey)
	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return &ecdsa.PublicKey{X: x, Y: y, Curve: btcec.S256()}, nil
}
