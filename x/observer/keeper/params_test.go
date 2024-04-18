package keeper_test

import (
	"fmt"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestGetParams(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	params := types.DefaultParams()

	err := k.SetParams(ctx, params)
	require.NoError(t, err)

	p, found := k.GetParams(ctx)
	require.True(t, found)
	require.EqualValues(t, params, p)
}

func TestGenerateAddress(t *testing.T) {
	addr := sdk.AccAddress(crypto.AddressHash([]byte("Output1" + strconv.Itoa(1))))
	addrString := addr.String()
	fmt.Println(addrString)
	addbech32, _ := sdk.AccAddressFromBech32(addrString)
	valAddress := sdk.ValAddress(addbech32)
	v, _ := sdk.ValAddressFromBech32(valAddress.String())
	fmt.Println(v.String())
	accAddress := sdk.AccAddress(v)
	a, _ := sdk.AccAddressFromBech32(accAddress.String())
	fmt.Println(a.String())
}
