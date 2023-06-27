//go:build TESTNET
// +build TESTNET

package testutil

import (
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	observerCli "github.com/zeta-chain/zetacore/x/observer/client/cli"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (s *IntegrationTestSuite) TestCCTXInboundVoter() {
	tt := []struct {
		name        string
		votes       map[string]observerTypes.VoteType
		finalresult observerTypes.BallotStatus
	}{
		{
			name: "All observers voted success",
			votes: map[string]observerTypes.VoteType{
				"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax": observerTypes.VoteType_SuccessObservation,
				"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2": observerTypes.VoteType_SuccessObservation,
				"zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4": observerTypes.VoteType_SuccessObservation,
				"zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c": observerTypes.VoteType_SuccessObservation,
				"zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca": observerTypes.VoteType_SuccessObservation,
				"zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt": observerTypes.VoteType_SuccessObservation,
				"zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4": observerTypes.VoteType_SuccessObservation,
				"zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy": observerTypes.VoteType_SuccessObservation,
				"zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav": observerTypes.VoteType_SuccessObservation,
				"zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t": observerTypes.VoteType_SuccessObservation,
			},
			finalresult: observerTypes.BallotStatus_BallotFinalized_SuccessObservation,
		},
		{
			name: "8 votes only",
			votes: map[string]observerTypes.VoteType{
				"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax": observerTypes.VoteType_SuccessObservation,
				"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2": observerTypes.VoteType_SuccessObservation,
				"zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4": observerTypes.VoteType_SuccessObservation,
				"zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c": observerTypes.VoteType_SuccessObservation,
				"zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca": observerTypes.VoteType_SuccessObservation,
				"zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt": observerTypes.VoteType_NotYetVoted,
				"zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4": observerTypes.VoteType_SuccessObservation,
				"zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy": observerTypes.VoteType_SuccessObservation,
				"zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav": observerTypes.VoteType_NotYetVoted,
				"zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t": observerTypes.VoteType_SuccessObservation,
			},
			finalresult: observerTypes.BallotStatus_BallotFinalized_SuccessObservation,
		},
		{
			name: "5 votes only",
			votes: map[string]observerTypes.VoteType{
				"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax": observerTypes.VoteType_SuccessObservation,
				"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2": observerTypes.VoteType_SuccessObservation,
				"zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4": observerTypes.VoteType_SuccessObservation,
				"zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c": observerTypes.VoteType_SuccessObservation,
				"zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca": observerTypes.VoteType_SuccessObservation,
				"zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt": observerTypes.VoteType_NotYetVoted,
				"zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4": observerTypes.VoteType_NotYetVoted,
				"zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy": observerTypes.VoteType_NotYetVoted,
				"zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav": observerTypes.VoteType_NotYetVoted,
				"zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t": observerTypes.VoteType_NotYetVoted,
			},
			finalresult: observerTypes.BallotStatus_BallotInProgress,
		},
		{
			name: "7 votes only,Just crossed threshold",
			votes: map[string]observerTypes.VoteType{
				"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax": observerTypes.VoteType_SuccessObservation,
				"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2": observerTypes.VoteType_SuccessObservation,
				"zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4": observerTypes.VoteType_SuccessObservation,
				"zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c": observerTypes.VoteType_SuccessObservation,
				"zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca": observerTypes.VoteType_SuccessObservation,
				"zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt": observerTypes.VoteType_NotYetVoted,
				"zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4": observerTypes.VoteType_SuccessObservation,
				"zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy": observerTypes.VoteType_NotYetVoted,
				"zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav": observerTypes.VoteType_NotYetVoted,
				"zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t": observerTypes.VoteType_SuccessObservation,
			},
			finalresult: observerTypes.BallotStatus_BallotFinalized_SuccessObservation,
		},
	}
	for _, test := range tt {
		test := test
		s.Run(test.name, func() {
			broadcaster := s.network.Validators[0]
			for _, val := range s.network.Validators {
				vote := test.votes[val.Address.String()]
				if vote == observerTypes.VoteType_NotYetVoted {
					continue
				}

				out, err := clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetAccountCmd(), []string{val.Address.String(), "--output", "json"})
				var account authtypes.AccountI
				s.NoError(val.ClientCtx.Codec.UnmarshalInterfaceJSON(out.Bytes(), &account))
				signedTx := BuildSignedInbound

				Vote(s.T(), val, s.cfg.BondDenom, account, test.name)
				out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetBroadcastCommand(), []string{signedTx.Name(), "--broadcast-mode", "sync"})
				s.Require().NoError(err)
			}

			s.Require().NoError(s.network.WaitForNBlocks(2))
			ballotIdentifier := GetBallotIdentifier(test.name)
			out, err := clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, observerCli.CmdBallotByIdentifier(), []string{ballotIdentifier, "--output", "json"})
			s.Require().NoError(err)
			ballot := observerTypes.QueryBallotByIdentifierResponse{}
			s.NoError(broadcaster.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &ballot))
			s.Assert().Equal(len(test.votes), len(ballot.Voters))
			for _, vote := range ballot.Voters {
				s.Assert().Equal(test.votes[vote.VoterAddress], vote.VoteType)
			}
			s.Assert().Equal(test.finalresult, ballot.BallotStatus)
		})

	}

}
