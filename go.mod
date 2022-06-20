module github.com/mycodeku/transtiontx

go 1.16

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

require (
	github.com/btcsuite/btcd v0.22.1
	github.com/ethereum/go-ethereum v1.10.16
	github.com/mycodeku/transtionhelper v0.0.0-20220620005849-a315113d9162
	github.com/shopspring/decimal v1.3.1
	github.com/tharsis/ethermint v0.16.1
)

replace github.com/tendermint/tendermint => github.com/tendermint/tendermint v0.34.20-0.20220517115723-e6f071164839

replace github.com/cosmos/cosmos-sdk => github.com/cosmos/cosmos-sdk v0.43.0

require (
	github.com/cosmos/cosmos-sdk v0.45.5-0.20220523154235-2921a1c3c918
	google.golang.org/grpc v1.47.0 // indirect
)
