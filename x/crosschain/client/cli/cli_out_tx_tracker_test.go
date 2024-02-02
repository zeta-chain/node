package cli_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/testutil/network"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func networkWithOutTxTrackerObjects(t *testing.T, n int) (*network.Network, []types.OutTxTracker) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	assert.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		outTxTracker := types.OutTxTracker{
			Index: fmt.Sprintf("testchain-%d", i),
		}
		nullify.Fill(&outTxTracker)
		state.OutTxTrackerList = append(state.OutTxTrackerList, outTxTracker)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	assert.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	//cfg.GenesisState = network.SetupZetaGenesisState(t, cfg.GenesisState, cfg.Codec)
	net, err := network.New(t, app.NodeDir, cfg)
	return net, state.OutTxTrackerList
}

//
//func TestShowOutTxTracker(t *testing.T) {
//	net, objs := networkWithOutTxTrackerObjects(t, 2)
//
//	ctx := net.Validators[0].ClientCtx
//	common := []string{
//		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
//	}
//	for _, tc := range []struct {
//		desc    string
//		idIndex string
//
//		args []string
//		err  error
//		obj  types.OutTxTracker
//	}{
//		{
//			desc:    "found",
//			idIndex: objs[0].Index,
//
//			args: common,
//			obj:  objs[0],
//		},
//		{
//			desc:    "not found",
//			idIndex: fmt.Sprintf("testchain-%d", 100000),
//
//			args: common,
//			err:  status.Error(codes.NotFound, "not found"),
//		},
//	} {
//		tc := tc
//		t.Run(tc.desc, func(t *testing.T) {
//			args := []string{
//				tc.idIndex,
//			}
//			args = append(args, tc.args...)
//			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowOutTxTracker(), args)
//			if tc.err != nil {
//				stat, ok := status.FromError(tc.err)
//				assert.True(t, ok)
//				assert.ErrorIs(t, stat.Err(), tc.err)
//			} else {
//				assert.NoError(t, err)
//				var resp types.QueryGetOutTxTrackerResponse
//				assert.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
//				assert.NotNil(t, resp.OutTxTracker)
//				assert.Equal(t,
//					nullify.Fill(&tc.obj),
//					nullify.Fill(&resp.OutTxTracker),
//				)
//			}
//		})
//	}
//}
//
//func TestListOutTxTracker(t *testing.T) {
//	net, objs := networkWithOutTxTrackerObjects(t, 5)
//
//	ctx := net.Validators[0].ClientCtx
//	request := func(next []byte, offset, limit uint64, total bool) []string {
//		args := []string{
//			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
//		}
//		if next == nil {
//			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
//		} else {
//			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
//		}
//		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
//		if total {
//			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
//		}
//		return args
//	}
//	t.Run("ByOffset", func(t *testing.T) {
//		step := 2
//		for i := 0; i < len(objs); i += step {
//			args := request(nil, uint64(i), uint64(step), false)
//			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListOutTxTracker(), args)
//			assert.NoError(t, err)
//			var resp types.QueryAllOutTxTrackerResponse
//			assert.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
//			assert.LessOrEqual(t, len(resp.OutTxTracker), step)
//			assert.Subset(t,
//				nullify.Fill(objs),
//				nullify.Fill(resp.OutTxTracker),
//			)
//		}
//	})
//	t.Run("ByKey", func(t *testing.T) {
//		step := 2
//		var next []byte
//		for i := 0; i < len(objs); i += step {
//			args := request(next, 0, uint64(step), false)
//			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListOutTxTracker(), args)
//			assert.NoError(t, err)
//			var resp types.QueryAllOutTxTrackerResponse
//			assert.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
//			assert.LessOrEqual(t, len(resp.OutTxTracker), step)
//			assert.Subset(t,
//				nullify.Fill(objs),
//				nullify.Fill(resp.OutTxTracker),
//			)
//			next = resp.Pagination.NextKey
//		}
//	})
//	t.Run("Total", func(t *testing.T) {
//		args := request(nil, 0, uint64(len(objs)), true)
//		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListOutTxTracker(), args)
//		assert.NoError(t, err)
//		var resp types.QueryAllOutTxTrackerResponse
//		assert.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
//		assert.NoError(t, err)
//		assert.Equal(t, len(objs), int(resp.Pagination.Total))
//		assert.ElementsMatch(t,
//			nullify.Fill(objs),
//			nullify.Fill(resp.OutTxTracker),
//		)
//	})
//}
