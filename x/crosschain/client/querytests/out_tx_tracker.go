package querytests

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/client/cli"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (s *CliTestSuite) TestListOutTxTracker() {
	ctx := s.network.Validators[0].ClientCtx
	objs := s.crossChainState.OutTxTrackerList
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
	s.Run("ByOffset", func() {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListOutTxTracker(), args)
			s.Require().NoError(err)
			var resp types.QueryAllOutTxTrackerResponse
			s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.OutTxTracker), step)
			s.Require().Subset(nullify.Fill(objs),
				nullify.Fill(resp.OutTxTracker),
			)
		}
	})
	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListOutTxTracker(), args)
			s.Require().NoError(err)
			var resp types.QueryAllOutTxTrackerResponse
			s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.OutTxTracker), step)
			s.Require().Subset(
				nullify.Fill(objs),
				nullify.Fill(resp.OutTxTracker),
			)
			next = resp.Pagination.NextKey
		}
	})
	s.Run("Total", func() {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListOutTxTracker(), args)
		s.Require().NoError(err)
		var resp types.QueryAllOutTxTrackerResponse
		s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().NoError(err)
		s.Require().Equal(len(objs), int(resp.Pagination.Total))
		s.Require().ElementsMatch(nullify.Fill(objs),
			nullify.Fill(resp.OutTxTracker),
		)
	})
}
