package querytests

import (
	"fmt"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/client/cli"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (s *CliTestSuite) TestListInboundTrackers() {
	ctx := s.network.Validators[0].ClientCtx
	objs := s.crosschainState.InboundTrackerList
	s.Run("List all trackers", func() {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListInboundTrackers(), args)
		s.Require().NoError(err)
		var resp types.QueryAllInboundTrackersResponse
		s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().Equal(len(objs), len(resp.InboundTracker))
		s.Require().ElementsMatch(nullify.Fill(objs), nullify.Fill(resp.InboundTracker))
	})
}

func (s *CliTestSuite) TestListInboundTrackersByChain() {
	ctx := s.network.Validators[0].ClientCtx
	objs := s.crosschainState.InboundTrackerList
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
			// #nosec G115 always positive
			args := request(nil, uint64(i), uint64(step), false, 5)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListInboundTrackerByChain(), args)
			s.Require().NoError(err)
			var resp types.QueryAllInboundTrackerByChainResponse
			s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.InboundTracker), step)
			s.Require().Subset(nullify.Fill(objs),
				nullify.Fill(resp.InboundTracker),
			)
		}
	})
	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			// #nosec G115 always positive
			args := request(next, 0, uint64(step), false, 5)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListInboundTrackerByChain(), args)
			s.Require().NoError(err)
			var resp types.QueryAllInboundTrackerByChainResponse
			s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.InboundTracker), step)
			s.Require().Subset(
				nullify.Fill(objs),
				nullify.Fill(resp.InboundTracker),
			)
			next = resp.Pagination.NextKey
		}
	})
	s.Run("Total", func() {
		args := request(nil, 0, uint64(len(objs)), true, 5)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListInboundTrackerByChain(), args)
		s.Require().NoError(err)
		var resp types.QueryAllInboundTrackerByChainResponse
		s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().NoError(err)
		s.Require().Equal(uint64(len(objs)), resp.Pagination.Total)
		s.Require().ElementsMatch(nullify.Fill(objs),
			nullify.Fill(resp.InboundTracker),
		)
	})
	s.Run("Incorrect Chain ID ", func() {
		args := request(nil, 0, uint64(len(objs)), true, 15)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListInboundTrackerByChain(), args)
		s.Require().NoError(err)
		var resp types.QueryAllInboundTrackerByChainResponse
		s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().NoError(err)
		s.Require().Equal(uint64(0), resp.Pagination.Total)
	})
}
