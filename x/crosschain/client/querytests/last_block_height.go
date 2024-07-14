package querytests

import (
	"fmt"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/zetacore/x/crosschain/client/cli"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (s *CliTestSuite) TestShowLastBlockHeight() {
	ctx := s.network.Validators[0].ClientCtx
	objs := s.crosschainState.LastBlockHeightList
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		id   string
		args []string
		err  error
		obj  *types.LastBlockHeight
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
		s.Run(tc.desc, func() {
			args := []string{tc.id}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowLastBlockHeight(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				s.Require().True(ok)
				s.Require().ErrorIs(stat.Err(), tc.err)
			} else {
				s.Require().NoError(err)
				var resp types.QueryGetLastBlockHeightResponse
				s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				s.Require().NotNil(resp.LastBlockHeight)
				s.Require().Equal(tc.obj, resp.LastBlockHeight)
			}
		})
	}
}

func (s *CliTestSuite) TestListLastBlockHeight() {
	ctx := s.network.Validators[0].ClientCtx
	objs := s.crosschainState.LastBlockHeightList
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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListLastBlockHeight(), args)
			s.Require().NoError(err)
			var resp types.QueryAllLastBlockHeightResponse
			s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			for j := i; j < len(objs) && j < i+step; j++ {
				s.Assert().Equal(objs[j], resp.LastBlockHeight[j-i])
			}
		}
	})
	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			// #nosec G115 always in range
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListLastBlockHeight(), args)
			s.Require().NoError(err)
			var resp types.QueryAllLastBlockHeightResponse
			s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			for j := i; j < len(objs) && j < i+step; j++ {
				s.Assert().Equal(objs[j], resp.LastBlockHeight[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	s.Run("Total", func() {
		// #nosec G115 always in range
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListLastBlockHeight(), args)
		s.Require().NoError(err)
		var resp types.QueryAllLastBlockHeightResponse
		s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().NoError(err)
		// #nosec G115 always in range
		s.Require().Equal(len(objs), int(resp.Pagination.Total))
		s.Require().Equal(objs, resp.LastBlockHeight)
	})
}
