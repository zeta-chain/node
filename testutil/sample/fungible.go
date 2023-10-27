package sample

import (
	"testing"

	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func ForeignCoins(t *testing.T, address string) types.ForeignCoins {
	r := newRandFromStringSeed(t, address)

	return types.ForeignCoins{
		Zrc20ContractAddress: address,
		Asset:                EthAddress().String(),
		ForeignChainId:       r.Int63(),
		Decimals:             uint32(r.Uint64()),
		Name:                 StringRandom(r, 32),
		Symbol:               StringRandom(r, 32),
		CoinType:             common.CoinType_ERC20,
		GasLimit:             r.Uint64(),
	}
}

func SystemContract() *types.SystemContract {
	return &types.SystemContract{
		SystemContract: EthAddress().String(),
		ConnectorZevm:  EthAddress().String(),
	}
}
