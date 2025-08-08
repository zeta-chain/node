package keeper_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_CheckPausedZRC20(t *testing.T) {
	addrUnpausedZRC20A, addrUnpausedZRC20B, addrUnpausedZRC20C, addrPausedZRC20 :=
		sample.EthAddress(),
		sample.EthAddress(),
		sample.EthAddress(),
		sample.EthAddress()

	tt := []struct {
		name    string
		receipt *ethtypes.Receipt
		wantErr bool
	}{
		{
			name:    "should pass if receipt is nil",
			receipt: nil,
			wantErr: false,
		},
		{
			name: "should pass if receipt is empty",
			receipt: &ethtypes.Receipt{
				Logs: []*ethtypes.Log{},
			},
			wantErr: false,
		},
		{
			name: "should pass if receipt contains unpaused ZRC20 and non ZRC20 addresses",
			receipt: &ethtypes.Receipt{
				Logs: []*ethtypes.Log{
					{
						Address: sample.EthAddress(),
					},
					{
						Address: addrUnpausedZRC20A,
					},
					{
						Address: addrUnpausedZRC20B,
					},
					{
						Address: addrUnpausedZRC20C,
					},
					{
						Address: addrUnpausedZRC20A,
					},
					{
						Address: addrUnpausedZRC20A,
					},
					nil,
					{
						Address: sample.EthAddress(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should fail if receipt contains paused ZRC20 and non ZRC20 addresses",
			receipt: &ethtypes.Receipt{
				Logs: []*ethtypes.Log{
					{
						Address: sample.EthAddress(),
					},
					{
						Address: addrUnpausedZRC20A,
					},
					{
						Address: addrUnpausedZRC20B,
					},
					{
						Address: addrUnpausedZRC20C,
					},
					{
						Address: addrPausedZRC20,
					},
					{
						Address: sample.EthAddress(),
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			k, ctx, _, _ := keepertest.FungibleKeeper(t)

			assertUnpaused := func(zrc20 string) {
				fc, found := k.GetForeignCoins(ctx, zrc20)
				require.True(t, found)
				require.False(t, fc.Paused)
			}
			assertPaused := func(zrc20 string) {
				fc, found := k.GetForeignCoins(ctx, zrc20)
				require.True(t, found)
				require.True(t, fc.Paused)
			}

			// setup ZRC20
			k.SetForeignCoins(ctx, sample.ForeignCoins(t, addrUnpausedZRC20A.Hex()))
			k.SetForeignCoins(ctx, sample.ForeignCoins(t, addrUnpausedZRC20B.Hex()))
			k.SetForeignCoins(ctx, sample.ForeignCoins(t, addrUnpausedZRC20C.Hex()))
			pausedZRC20 := sample.ForeignCoins(t, addrPausedZRC20.Hex())
			pausedZRC20.Paused = true
			k.SetForeignCoins(ctx, pausedZRC20)

			// check paused status
			assertUnpaused(addrUnpausedZRC20A.Hex())
			assertUnpaused(addrUnpausedZRC20B.Hex())
			assertUnpaused(addrUnpausedZRC20C.Hex())
			assertPaused(addrPausedZRC20.Hex())

			// process test
			h := k.EVMHooks()
			err := h.PostTxProcessing(ctx, common.Address{}, core.Message{}, tc.receipt)
			if tc.wantErr {
				require.ErrorIs(t, err, types.ErrPausedZRC20)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
