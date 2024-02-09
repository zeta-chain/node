package keeper

import (
	"fmt"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := SetupKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}

func TestGenerateAddress(t *testing.T) {
	types.SetConfig(false)
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
