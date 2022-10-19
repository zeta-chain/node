package cli_test

import (
	"fmt"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/zetacore/testutil/network"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/fungible/client/cli"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func networkWithZetaDepositAndCallContractObjects(t *testing.T) (*network.Network, types.ZetaDepositAndCallContract) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	zetaDepositAndCallContract := &types.ZetaDepositAndCallContract{}
	nullify.Fill(&zetaDepositAndCallContract)
	state.ZetaDepositAndCallContract = zetaDepositAndCallContract
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), *state.ZetaDepositAndCallContract
}

func TestShowZetaDepositAndCallContract(t *testing.T) {
	net, obj := networkWithZetaDepositAndCallContractObjects(t)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		args []string
		err  error
		obj  types.ZetaDepositAndCallContract
	}{
		{
			desc: "get",
			args: common,
			obj:  obj,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			var args []string
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowZetaDepositAndCallContract(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetZetaDepositAndCallContractResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.ZetaDepositAndCallContract)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.ZetaDepositAndCallContract),
				)
			}
		})
	}
}
