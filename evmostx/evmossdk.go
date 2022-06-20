package evmostx

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	comostxtypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
)

type EvmosSdk struct {
	FromAdd, ToAdd, Memo, Denom, Value, Fee          string
	Precision                                        uint8
	TimeoutHeight, Sequence, GasLimit, AccountNumber uint64
	PubkeyByte                                       []byte
}

func (self *EvmosSdk) precisionTransaction(value string, Precision uint8) types.Int {
	p := decimal.New(1, int32(Precision))
	v, err := decimal.NewFromString(value)
	if err != nil {
		panic(fmt.Errorf("precisionTransaction get value err %v", err))
	}
	return types.NewIntFromBigInt(v.Mul(p).BigInt())
}

func (self *EvmosSdk) CreateTx() client.TxBuilder {
	initClientCtx := client.Context{}.WithChainID("evmos_9001-2")
	initClientCtx = initClientCtx.WithTxConfig(simapp.MakeTestEncodingConfig().TxConfig)

	amount := self.precisionTransaction(self.Value, self.Precision)
	amounttosend, ok := types.NewIntFromString(amount.String())
	if !ok {
		panic(fmt.Errorf("create amount err"))
	}
	sendMsg := &comostxtypes.MsgSend{FromAddress: self.FromAdd, ToAddress: self.ToAdd, Amount: types.Coins{types.NewCoin("aevmos", amounttosend)}}
	builder := initClientCtx.TxConfig.NewTxBuilder()
	builder.SetMsgs(sendMsg)
	builder.SetGasLimit(self.GasLimit)
	builder.SetTimeoutHeight(self.TimeoutHeight)

	feeamount := self.precisionTransaction(self.Fee, self.Precision)
	feeamounttosend, ok := types.NewIntFromString(feeamount.String())
	builder.SetFeeAmount(types.Coins{types.NewCoin("aevmos", feeamounttosend)})
	aa := &ethsecp256k1.PubKey{Key: self.PubkeyByte}
	sig := signing.SignatureV2{
		PubKey: aa,
		Data: &signing.SingleSignatureData{
			SignMode: signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
		},
		Sequence: self.Sequence,
	}

	builder.SetSignatures(sig)

	return builder
}
func (self *EvmosSdk) CreatTxByte(builder client.TxBuilder) (string, string) {
	initClientCtx := client.Context{}.WithChainID("evmos_9001-2")
	initClientCtx = initClientCtx.WithTxConfig(simapp.MakeTestEncodingConfig().TxConfig)

	signerData := xauthsigning.SignerData{
		ChainID:       "evmos_9001-2",
		AccountNumber: self.AccountNumber,
		Sequence:      self.Sequence,
	}
	signBytes, err := initClientCtx.TxConfig.SignModeHandler().GetSignBytes(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON, signerData, builder.GetTx())
	if err != nil {
		fmt.Println(err)
	}
	return hex.EncodeToString(crypto.Keccak256Hash(signBytes).Bytes()), hex.EncodeToString(signBytes)
}

func (self *EvmosSdk) GetTxWithSignature(sign string, builder client.TxBuilder) string {
	signBytes, err := hex.DecodeString(sign)
	if err != nil {
		fmt.Println(err)
	}
	signBytes = signBytes[:64]

	sigData := signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
		Signature: signBytes,
	}
	aa := &ethsecp256k1.PubKey{Key: self.PubkeyByte}
	sig := signing.SignatureV2{
		PubKey:   aa,
		Data:     &sigData,
		Sequence: self.Sequence,
	}
	builder.SetSignatures(sig)

	initClientCtx := client.Context{}.WithChainID("evmos_9001-2")
	initClientCtx = initClientCtx.WithTxConfig(simapp.MakeTestEncodingConfig().TxConfig)
	txBytes, err := initClientCtx.TxConfig.TxEncoder()(builder.GetTx())

	return base64.StdEncoding.EncodeToString(txBytes)
}

func (self *EvmosSdk) PubKey() cryptotypes.PubKey {
	return &ethsecp256k1.PubKey{
		Key: self.PubkeyByte,
	}
}
