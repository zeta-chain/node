package querytests

import (
	"fmt"
	"strconv"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/client/cli"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (s *CliTestSuite) TestShowInboundHashToCctx() {
	ctx := s.network.Validators[0].ClientCtx
	objs := s.crosschainState.InboundHashToCctxList
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc          string
		idInboundHash string

		args []string
		err  error
		obj  types.InboundHashToCctx
	}{
		{
			desc:          "found",
			idInboundHash: objs[0].InboundHash,

			args: common,
			obj:  objs[0],
		},
		{
			desc:          "not found",
			idInboundHash: strconv.Itoa(100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	} {
		s.Run(tc.desc, func() {
			args := []string{
				tc.idInboundHash,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowInboundHashToCctx(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				s.Require().True(ok)
				s.Require().ErrorIs(stat.Err(), tc.err)
			} else {
				s.Require().NoError(err)
				var resp types.QueryGetInboundHashToCctxResponse
				s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				s.Require().NotNil(resp.InboundHashToCctx)
				tc := tc
				s.Require().Equal(nullify.Fill(&tc.obj),
					nullify.Fill(&resp.InboundHashToCctx),
				)
			}
		})
	}
}

func (s *CliTestSuite) TestListInboundHashToCctx() {
	ctx := s.network.Validators[0].ClientCtx
	objs := s.crosschainState.InboundHashToCctxList
	cctxCount := len(s.crosschainState.CrossChainTxs)
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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListInboundHashToCctx(), args)
			s.Require().NoError(err)
			var resp types.QueryAllInboundHashToCctxResponse
			s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.InboundHashToCctx), step)
			s.Require().Subset(nullify.Fill(objs),
				nullify.Fill(resp.InboundHashToCctx),
			)
		}
	})
	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			// #nosec G115 always in range
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListInboundHashToCctx(), args)
			s.Require().NoError(err)
			var resp types.QueryAllInboundHashToCctxResponse
			s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.InboundHashToCctx), step)
			s.Require().Subset(nullify.Fill(objs),
				nullify.Fill(resp.InboundHashToCctx),
			)
			next = resp.Pagination.NextKey
		}
	})
	s.Run("Total", func() {
		// #nosec G115 always in range
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListInboundHashToCctx(), args)
		s.Require().NoError(err)
		var resp types.QueryAllInboundHashToCctxResponse
		s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().NoError(err)
		// saving CCTX also adds a new mapping
		// #nosec G115 always in range
		s.Require().Equal(len(objs)+cctxCount, int(resp.Pagination.Total))
		s.Require().ElementsMatch(nullify.Fill(objs),
			nullify.Fill(resp.InboundHashToCctx),
		)
	})
}
