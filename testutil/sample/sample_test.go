package sample

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/app"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestEthAddress(t *testing.T) {

	// just checking something

	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cfg.Seal()

	a := "zeta1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsxn0x73"
	addr, err := sdk.AccAddressFromBech32(a)
	require.NoError(t, err)

	foo := ethcommon.BytesToAddress(addr.Bytes()).String()
	_ = foo

	ethAddress := EthAddress()
	require.NotEqual(t, ethcommon.Address{}, ethAddress)

	// don't generate the same address
	ethAddress2 := EthAddress()
	require.NotEqual(t, ethAddress, ethAddress2)
}
