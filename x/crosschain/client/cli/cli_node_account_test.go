package cli_test

import (
	"fmt"
	"github.com/zeta-chain/zetacore/app"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/zetacore/testutil/network"
	"github.com/zeta-chain/zetacore/x/crosschain/client/cli"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func networkWithNodeAccountObjects(t *testing.T, n int) (*network.Network, []*types.NodeAccount) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		state.NodeAccountList = append(state.NodeAccountList, &types.NodeAccount{Operator: strconv.Itoa(i), GranteeAddress: "signer"})
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	cfg.GenesisState = network.SetupZetaGenesisState(t, cfg.GenesisState, cfg.Codec)
	net, err := network.New(t, app.NodeDir, cfg)
	return net, state.NodeAccountList
}

func TestShowNodeAccount(t *testing.T) {
	net, objs := networkWithNodeAccountObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		id   string
		args []string
		err  error
		obj  *types.NodeAccount
	}{
		{
			desc: "found",
			id:   objs[0].Operator,
			args: common,
			obj:  objs[0],
		},
		{
			desc: "not found",
			id:   "not_found",
			args: common,
			err:  status.Error(codes.InvalidArgument, "not found"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{tc.id}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowNodeAccount(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetNodeAccountResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.NodeAccount)
				require.Equal(t, tc.obj, resp.NodeAccount)
			}
		})
	}
}

func TestListNodeAccount(t *testing.T) {
	net, objs := networkWithNodeAccountObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListNodeAccount(), args)
			require.NoError(t, err)
			var resp types.QueryAllNodeAccountResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			for j := i; j < len(objs) && j < i+step; j++ {
				assert.Equal(t, objs[j], resp.NodeAccount[j-i])
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListNodeAccount(), args)
			require.NoError(t, err)
			var resp types.QueryAllNodeAccountResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			for j := i; j < len(objs) && j < i+step; j++ {
				assert.Equal(t, objs[j], resp.NodeAccount[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListNodeAccount(), args)
		require.NoError(t, err)
		var resp types.QueryAllNodeAccountResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.Equal(t, objs, resp.NodeAccount)
	})
}
