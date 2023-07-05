package cli_test

import (
	"fmt"
	"github.com/zeta-chain/zetacore/app"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/zetacore/testutil/network"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/client/cli"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func networkWithPermissionFlagsObjects(t *testing.T) (*network.Network, types.PermissionFlags) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	permissionFlags := &types.PermissionFlags{}
	nullify.Fill(&permissionFlags)
	state.PermissionFlags = permissionFlags
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	cfg.GenesisState = network.SetupZetaGenesisState(t, cfg.GenesisState, cfg.Codec)
	net, err := network.New(t, app.NodeDir, cfg)
	return net, *state.PermissionFlags
}

func TestShowPermissionFlags(t *testing.T) {
	net, obj := networkWithPermissionFlagsObjects(t)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		args []string
		err  error
		obj  types.PermissionFlags
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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowPermissionFlags(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetPermissionFlagsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.PermissionFlags)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.PermissionFlags),
				)
			}
		})
	}
}
