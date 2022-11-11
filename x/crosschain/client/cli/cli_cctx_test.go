package cli_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	fungibleModuleTypes "github.com/zeta-chain/zetacore/x/fungible/types"
	"strconv"
	"testing"
	"time"

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

func networkWithSendObjects(t *testing.T, n int) (*network.Network, []*types.CrossChainTx) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	bankState := banktypes.GenesisState{}
	//evmState := evmtypes.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[banktypes.ModuleName], &bankState))
	//require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[evmtypes.ModuleName], &evmState))
	amount, _ := sdk.NewIntFromString("100000000000000000000000000000000")
	bankState.Balances = append(bankState.Balances, banktypes.Balance{
		Address: fungibleModuleTypes.ModuleAddress.String(),
		Coins: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, amount),
			sdk.NewCoin("azeta", amount),
			sdk.NewCoin("zeta", amount),
			sdk.NewCoin("gETH", amount),
			sdk.NewCoin("geth", amount),
			sdk.NewCoin("eth", amount),
			sdk.NewCoin("ETH", amount),
		),
	})
	//evmState.Accounts = append(evmState.Accounts, evmtypes.GenesisAccount{
	//	Address: fungibleModuleTypes.ModuleAddressEVM.String(),
	//	Code:    "",
	//	Storage: evmtypes.Storage{
	//		{Key: common.BytesToHash([]byte("key")).String(), Value: common.BytesToHash([]byte("value")).String()},
	//	},
	//})
	for i := 0; i < n; i++ {
		state.CrossChainTxs = append(state.CrossChainTxs, &types.CrossChainTx{
			Creator: "ANY",
			Index:   strconv.Itoa(i),
			CctxStatus: &types.Status{
				Status:              types.CctxStatus_PendingInbound,
				StatusMessage:       "",
				LastUpdateTimestamp: 0,
			},
			ZetaMint:  sdk.OneUint(),
			ZetaBurnt: sdk.OneUint(),
			ZetaFees:  sdk.OneUint()},
		)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	buf, err = cfg.Codec.MarshalJSON(&bankState)
	require.NoError(t, err)
	cfg.GenesisState[banktypes.ModuleName] = buf
	//buf, err = cfg.Codec.MarshalJSON(&evmState)
	//require.NoError(t, err)
	//cfg.GenesisState[evmtypes.ModuleName] = buf

	net := network.New(t, cfg)
	_, err = net.WaitForHeight(1)
	return net, state.CrossChainTxs
}

func TestShowSend(t *testing.T) {
	net, objs := networkWithSendObjects(t, 2)
	h, err := net.WaitForHeightWithTimeout(2, time.Minute*2)
	//err := net.WaitForNextBlock()
	//assert.NoError(t, err)
	//err = net.WaitForNextBlock()
	//assert.NoError(t, err)
	//h, err := net.LatestHeight()
	assert.NoError(t, err)
	fmt.Println(h)
	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		id   string
		args []string
		err  error
		obj  *types.CrossChainTx
	}{
		{
			desc: "found",
			id:   objs[0].Index,
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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowSend(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetCctxResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.CrossChainTx)
				require.Equal(t, tc.obj, resp.CrossChainTx)
			}
		})
	}
}

func TestListSend(t *testing.T) {
	net, objs := networkWithSendObjects(t, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListSend(), args)
			require.NoError(t, err)
			var resp types.QueryAllCctxResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			for j := i; j < len(objs) && j < i+step; j++ {
				assert.Equal(t, objs[j], resp.CrossChainTx[j-i])
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListSend(), args)
			require.NoError(t, err)
			var resp types.QueryAllCctxResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			for j := i; j < len(objs) && j < i+step; j++ {
				assert.Equal(t, objs[j], resp.CrossChainTx[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListSend(), args)
		require.NoError(t, err)
		var resp types.QueryAllCctxResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.Equal(t, objs, resp.CrossChainTx)
	})
}
