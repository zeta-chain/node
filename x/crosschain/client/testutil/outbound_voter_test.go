//go:build PRIVNET
// +build PRIVNET

package testutil

import (
	"fmt"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	crosschainCli "github.com/zeta-chain/zetacore/x/crosschain/client/cli"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (s *IntegrationTestSuite) TestCCTXOutBoundVoter() {
	type Vote struct {
		voterAddress string
		voteType     observerTypes.VoteType
		isFakeVote   bool
	}
	tt := []struct {
		name                  string
		votes                 []Vote
		correctBallotResult   observerTypes.BallotStatus
		cctxStatus            crosschaintypes.CctxStatus
		falseBallotIdentifier string
	}{
		{
			name: "All observers voted success",
			votes: []Vote{
				{voterAddress: "zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
			},
			correctBallotResult: observerTypes.BallotStatus_BallotFinalized_SuccessObservation,
			cctxStatus:          crosschaintypes.CctxStatus_PendingOutbound,
		},
		{
			name: "All observers voted success 2",
			votes: []Vote{
				{voterAddress: "zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t", voteType: observerTypes.VoteType_SuccessObservation, isFakeVote: false},
			},
			correctBallotResult: observerTypes.BallotStatus_BallotFinalized_SuccessObservation,
			cctxStatus:          crosschaintypes.CctxStatus_PendingOutbound,
		},
	}
	for _, test := range tt {
		test := test
		s.Run(test.name, func() {
			broadcaster := s.network.Validators[0]
			for _, val := range s.network.Validators {
				out, err := clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetAccountCmd(), []string{val.Address.String(), "--output", "json"})
				var account authtypes.AccountI
				s.NoError(val.ClientCtx.Codec.UnmarshalInterfaceJSON(out.Bytes(), &account))
				signedTx := BuildSignedGasPriceVote(s.T(), val, s.cfg.BondDenom, account)
				_, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetBroadcastCommand(), []string{signedTx.Name(), "--broadcast-mode", "sync"})
				s.Require().NoError(err)
			}
			s.Require().NoError(s.network.WaitForNBlocks(2))
			for _, val := range s.network.Validators {
				out, err := clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetAccountCmd(), []string{val.Address.String(), "--output", "json"})
				var account authtypes.AccountI
				s.NoError(val.ClientCtx.Codec.UnmarshalInterfaceJSON(out.Bytes(), &account))
				signedTx := BuildSignedTssVote(s.T(), val, s.cfg.BondDenom, account)
				out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetBroadcastCommand(), []string{signedTx.Name(), "--broadcast-mode", "sync"})
				s.Require().NoError(err)
			}
			s.Require().NoError(s.network.WaitForNBlocks(2))
			for _, val := range s.network.Validators {
				out, err := clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetAccountCmd(), []string{val.Address.String(), "--output", "json"})
				var account authtypes.AccountI
				s.NoError(val.ClientCtx.Codec.UnmarshalInterfaceJSON(out.Bytes(), &account))
				message := test.name
				signedTx := BuildSignedInboundVote(s.T(), val, s.cfg.BondDenom, account, message)
				out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetBroadcastCommand(), []string{signedTx.Name(), "--broadcast-mode", "sync"})
				s.Require().NoError(err)
			}

			s.Require().NoError(s.network.WaitForNBlocks(2))
			cctxIdentifier := GetBallotIdentifier(test.name)
			out, err := clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, crosschainCli.CmdShowSend(), []string{cctxIdentifier, "--output", "json"})
			cctx := crosschaintypes.QueryGetCctxResponse{}
			s.NoError(broadcaster.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &cctx))
			s.Assert().Equal(crosschaintypes.CctxStatus_PendingOutbound, cctx.CrossChainTx.CctxStatus.Status)

			for _, val := range s.network.Validators {
				valVote := Vote{}
				for _, vote := range test.votes {
					if vote.voterAddress == val.Address.String() {
						valVote = vote
					}
				}
				if valVote.voteType == observerTypes.VoteType_NotYetVoted {
					continue
				}
				out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetAccountCmd(), []string{val.Address.String(), "--output", "json"})
				var account authtypes.AccountI
				s.NoError(val.ClientCtx.Codec.UnmarshalInterfaceJSON(out.Bytes(), &account))
				message := test.name
				if valVote.isFakeVote {
					message = message + "falseVote"
				}
				signedTx := BuildSignedOutboundVote(s.T(), val, s.cfg.BondDenom, account, message, cctxIdentifier)
				out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetBroadcastCommand(), []string{signedTx.Name(), "--broadcast-mode", "sync"})
				s.Require().NoError(err)
			}
			s.Require().NoError(s.network.WaitForNBlocks(2))
			out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, crosschainCli.CmdShowSend(), []string{cctxIdentifier, "--output", "json"})
			cctx = crosschaintypes.QueryGetCctxResponse{}
			s.NoError(broadcaster.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &cctx))
			fmt.Println(cctx.CrossChainTx.CctxStatus.Status)
		})
	}

}
