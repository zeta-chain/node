package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
)

func TestSendVoterMsgServerCreate(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	srv := NewMsgServerImpl(*keeper)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "A"
	for i := 0; i < 5; i++ {
		idx := fmt.Sprintf("%d", i)
		expected := &types.MsgCreateSendVoter{Creator: creator, TxHash: fmt.Sprintf("txahsh%s", idx), Index: idx}
		_, err := srv.CreateSendVoter(wctx, expected)
		require.NoError(t, err)
		rst, found := keeper.GetSendVoter(ctx, expected.Index)
		require.True(t, found)
		assert.Equal(t, expected.Creator, rst.Creator)
	}
}
