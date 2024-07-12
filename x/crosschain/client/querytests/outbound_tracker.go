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

func (s *CliTestSuite) TestListOutboundTracker() {
	ctx := s.network.Validators[0].ClientCtx
	objs := s.crosschainState.OutboundTrackerList
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
			// #nosec G115 always in range
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListOutboundTracker(), args)
			s.Require().NoError(err)
			var resp types.QueryAllOutboundTrackerResponse
			s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.OutboundTracker), step)
			s.Require().Subset(nullify.Fill(objs),
				nullify.Fill(resp.OutboundTracker),
			)
		}
	})
	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			// #nosec G115 always in range
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListOutboundTracker(), args)
			s.Require().NoError(err)
			var resp types.QueryAllOutboundTrackerResponse
			s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.OutboundTracker), step)
			s.Require().Subset(
				nullify.Fill(objs),
				nullify.Fill(resp.OutboundTracker),
			)
			next = resp.Pagination.NextKey
		}
	})
	s.Run("Total", func() {
		// #nosec G115 always in range
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListOutboundTracker(), args)
		s.Require().NoError(err)
		var resp types.QueryAllOutboundTrackerResponse
		s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().NoError(err)
		// #nosec G115 always in range
		s.Require().Equal(len(objs), int(resp.Pagination.Total))
		s.Require().ElementsMatch(nullify.Fill(objs),
			nullify.Fill(resp.OutboundTracker),
		)
	})
}
