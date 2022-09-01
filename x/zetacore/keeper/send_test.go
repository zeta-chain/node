package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"math/rand"
	"testing"
)

func createNSend(keeper *Keeper, ctx sdk.Context, n int) []types.Send {
	items := make([]types.Send, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetSend(ctx, items[i])
	}
	return items
}
func createNSendWithStatus(keeper *Keeper, ctx sdk.Context, n int, status types.SendStatus) []types.Send {
	items := make([]types.Send, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		items[i].Status = status
		keeper.SetSend(ctx, items[i])
	}
	return items
}

func TestSends(t *testing.T) {
	sendsTest := []struct {
		TestName        string
		PendingInbound  int
		PendingOutbound int
		OutboundMined   int
		Confirmed       int
		PendingRevert   int
		Reverted        int
		Aborted         int
	}{
		{
			TestName:        "test pending",
			PendingInbound:  10,
			PendingOutbound: 10,
			Confirmed:       10,
			PendingRevert:   10,
			Aborted:         10,
			OutboundMined:   10,
			Reverted:        10,
		},
		{
			TestName:        "test pending",
			PendingInbound:  rand.Intn(300-10) + 10,
			PendingOutbound: rand.Intn(300-10) + 10,
			Confirmed:       rand.Intn(300-10) + 10,
			PendingRevert:   rand.Intn(300-10) + 10,
			Aborted:         rand.Intn(300-10) + 10,
			OutboundMined:   rand.Intn(300-10) + 10,
			Reverted:        rand.Intn(300-10) + 10,
		},
	}
	for _, tt := range sendsTest {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			keeper, ctx := setupKeeper(t)
			sends := []types.Send{}
			sends = append(sends, createNSendWithStatus(keeper, ctx, tt.PendingInbound, types.SendStatus_PendingInbound)...)
			sends = append(sends, createNSendWithStatus(keeper, ctx, tt.PendingOutbound, types.SendStatus_PendingOutbound)...)
			sends = append(sends, createNSendWithStatus(keeper, ctx, tt.Confirmed, types.SendStatus_Confirmed)...)
			sends = append(sends, createNSendWithStatus(keeper, ctx, tt.PendingRevert, types.SendStatus_PendingRevert)...)
			sends = append(sends, createNSendWithStatus(keeper, ctx, tt.Aborted, types.SendStatus_Aborted)...)
			sends = append(sends, createNSendWithStatus(keeper, ctx, tt.OutboundMined, types.SendStatus_OutboundMined)...)
			sends = append(sends, createNSendWithStatus(keeper, ctx, tt.Reverted, types.SendStatus_Reverted)...)
			assert.Equal(t, tt.PendingOutbound, len(keeper.GetAllPendingOutBoundSend(ctx)))
			assert.Equal(t, tt.PendingInbound, len(keeper.GetAllPendingInBoundSend(ctx)))
			assert.Equal(t, len(sends), len(keeper.GetAllSend(ctx)))
			for _, s := range sends {
				send, found := keeper.GetSend(ctx, s.Index, s.Status)
				assert.True(t, found)
				assert.Equal(t, s, send)
			}

		})
	}
}
