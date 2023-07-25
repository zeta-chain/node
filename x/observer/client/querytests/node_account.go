package querytests

import (
	"fmt"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/zeta-chain/zetacore/x/observer/client/cli"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CliTestSuite) TestShowNodeAccount() {
	ctx := s.network.Validators[0].ClientCtx
	objs := s.observerState.NodeAccountList
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		id   string
		args []string
		err  error
		obj  *zetaObserverTypes.NodeAccount
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
		s.Run(tc.desc, func() {
			args := []string{tc.id}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowNodeAccount(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				s.Require().True(ok)
				s.Require().ErrorIs(stat.Err(), tc.err)
			} else {
				s.Require().NoError(err)
				var resp zetaObserverTypes.QueryGetNodeAccountResponse
				s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				s.Require().NotNil(resp.NodeAccount)
				s.Require().Equal(tc.obj, resp.NodeAccount)
			}
		})
	}
}
