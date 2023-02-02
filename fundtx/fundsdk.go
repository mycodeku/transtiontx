package fundtx

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	comostxtypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

type FundSdk struct {
	FromAdd, ToAdd, Memo, Denom, Value, Fee          string
	Precision                                        uint8
	TimeoutHeight, Sequence, GasLimit, AccountNumber uint64
	PubkeyByte                                       []byte
}

func (self *FundSdk) precisionTransaction(value string, Precision uint8) types.Int {
	p := decimal.New(1, int32(Precision))
	v, err := decimal.NewFromString(value)
	if err != nil {
		panic(fmt.Errorf("precisionTransaction get value err %v", err))
	}
	return types.NewIntFromBigInt(v.Mul(p).BigInt())
}
func (self *FundSdk) CreateTx(ChainId string, CoinType string) client.TxBuilder {
	initClientCtx := client.Context{}.WithChainID(ChainId)
	initClientCtx = initClientCtx.WithTxConfig(simapp.MakeTestEncodingConfig().TxConfig)

	amount := self.precisionTransaction(self.Value, self.Precision)
	amounttosend, ok := types.NewIntFromString(amount.String())
	if !ok {
		panic(fmt.Errorf("create amount err"))
	}
	sendMsg := &comostxtypes.MsgSend{FromAddress: self.FromAdd, ToAddress: self.ToAdd, Amount: types.Coins{types.NewCoin(CoinType, amounttosend)}}
	builder := initClientCtx.TxConfig.NewTxBuilder()
	builder.SetMsgs(sendMsg)
	builder.SetGasLimit(self.GasLimit)
	builder.SetTimeoutHeight(self.TimeoutHeight)
	builder.SetMemo(self.Memo)

	feeamount := self.precisionTransaction(self.Fee, self.Precision)
	feeamounttosend, ok := types.NewIntFromString(feeamount.String())
	builder.SetFeeAmount(types.Coins{types.NewCoin(CoinType, feeamounttosend)})
	aa := &secp256k1.PubKey{Key: self.PubkeyByte}
	aa.Type()
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

func (self *FundSdk) CreatTxByte(builder client.TxBuilder, ChainId string) (string, string) {
	initClientCtx := client.Context{}.WithChainID(ChainId)
	initClientCtx = initClientCtx.WithTxConfig(simapp.MakeTestEncodingConfig().TxConfig)

	signerData := xauthsigning.SignerData{
		ChainID:       ChainId,
		AccountNumber: self.AccountNumber,
		Sequence:      self.Sequence,
	}
	signBytes, err := initClientCtx.TxConfig.SignModeHandler().GetSignBytes(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON, signerData, builder.GetTx())
	if err != nil {
		fmt.Println(err)
	}
	return hex.EncodeToString(crypto.Keccak256Hash(signBytes).Bytes()), hex.EncodeToString(signBytes)
}

func (self *FundSdk) GetTxWithSignature(sign string, builder client.TxBuilder, ChainId string) (string, string) {
	signBytes, err := hex.DecodeString(sign)
	if err != nil {
		fmt.Println(err)
	}
	signBytes = signBytes[:64]

	sigData := signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
		Signature: signBytes,
	}
	aa := &secp256k1.PubKey{Key: self.PubkeyByte}
	sig := signing.SignatureV2{
		PubKey:   aa,
		Data:     &sigData,
		Sequence: self.Sequence,
	}
	builder.SetSignatures(sig)

	initClientCtx := client.Context{}.WithChainID(ChainId)
	initClientCtx = initClientCtx.WithTxConfig(simapp.MakeTestEncodingConfig().TxConfig)
	txBytes, err := initClientCtx.TxConfig.TxEncoder()(builder.GetTx())
	txmg, _ := json.Marshal(builder.GetTx())
	return base64.StdEncoding.EncodeToString(txBytes), hex.EncodeToString(txmg)
}

func (self *FundSdk) PubKey() cryptotypes.PubKey {
	return &secp256k1.PubKey{
		Key: self.PubkeyByte,
	}
}
