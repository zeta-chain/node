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

func (s *CliTestSuite) TestListInTxTrackers() {
	ctx := s.network.Validators[0].ClientCtx
	objs := s.crosschainState.InTxTrackerList
	s.Run("List all trackers", func() {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListInTxTrackers(), args)
		s.Require().NoError(err)
		var resp types.QueryAllInTxTrackersResponse
		s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().Equal(len(objs), len(resp.InTxTracker))
		s.Require().ElementsMatch(nullify.Fill(objs), nullify.Fill(resp.InTxTracker))
	})
}

func (s *CliTestSuite) TestListInTxTrackersByChain() {
	ctx := s.network.Validators[0].ClientCtx
	objs := s.crosschainState.InTxTrackerList
	request := func(next []byte, offset, limit uint64, total bool, chainID int) []string {
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
		args = append(args, fmt.Sprintf("%d", chainID))
		return args
	}
	s.Run("ByOffset", func() {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false, 5)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListInTxTrackerByChain(), args)
			s.Require().NoError(err)
			var resp types.QueryAllInTxTrackerByChainResponse
			s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.InTxTracker), step)
			s.Require().Subset(nullify.Fill(objs),
				nullify.Fill(resp.InTxTracker),
			)
		}
	})
	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false, 5)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListInTxTrackerByChain(), args)
			s.Require().NoError(err)
			var resp types.QueryAllInTxTrackerByChainResponse
			s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.InTxTracker), step)
			s.Require().Subset(
				nullify.Fill(objs),
				nullify.Fill(resp.InTxTracker),
			)
			next = resp.Pagination.NextKey
		}
	})
	s.Run("Total", func() {
		args := request(nil, 0, uint64(len(objs)), true, 5)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListInTxTrackerByChain(), args)
		s.Require().NoError(err)
		var resp types.QueryAllInTxTrackerByChainResponse
		s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().NoError(err)
		s.Require().Equal(len(objs), int(resp.Pagination.Total))
		s.Require().ElementsMatch(nullify.Fill(objs),
			nullify.Fill(resp.InTxTracker),
		)
	})
	s.Run("Incorrect Chain ID ", func() {
		args := request(nil, 0, uint64(len(objs)), true, 15)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListInTxTrackerByChain(), args)
		s.Require().NoError(err)
		var resp types.QueryAllInTxTrackerByChainResponse
		s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().NoError(err)
		s.Require().Equal(0, int(resp.Pagination.Total))
	})
}
