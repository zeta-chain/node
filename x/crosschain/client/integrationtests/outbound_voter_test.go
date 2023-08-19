//go:build TESTNET
// +build TESTNET

package integrationtests

import (
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	crosschaincli "github.com/zeta-chain/zetacore/x/crosschain/client/cli"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observercli "github.com/zeta-chain/zetacore/x/observer/client/cli"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (s *IntegrationTestSuite) TestCCTXOutBoundVoter() {
	type Vote struct {
		voterAddress string
		voteType     observertypes.VoteType
		isFakeVote   bool
	}
	tt := []struct {
		name                  string
		votes                 []Vote
		zetaMinted            string // TODO : calculate this value
		correctBallotResult   observertypes.BallotStatus
		cctxStatus            crosschaintypes.CctxStatus
		falseBallotIdentifier string
	}{
		{
			name: "All observers voted success or not voted",
			votes: []Vote{
				{voterAddress: "zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t", voteType: observertypes.VoteType_NotYetVoted, isFakeVote: false},
			},
			correctBallotResult: observertypes.BallotStatus_BallotFinalized_SuccessObservation,
			cctxStatus:          crosschaintypes.CctxStatus_OutboundMined,
			zetaMinted:          "7991636132140714751",
		},
		{
			name: "1 fake vote but ballot still success",
			votes: []Vote{
				{voterAddress: "zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
				{voterAddress: "zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
			},
			correctBallotResult: observertypes.BallotStatus_BallotFinalized_SuccessObservation,
			cctxStatus:          crosschaintypes.CctxStatus_OutboundMined,
			zetaMinted:          "7990439496224753106",
		},
		{
			name: "Half success and half false",
			votes: []Vote{
				{voterAddress: "zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
				{voterAddress: "zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
				{voterAddress: "zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
				{voterAddress: "zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
				{voterAddress: "zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
			},
			correctBallotResult: observertypes.BallotStatus_BallotInProgress,
			cctxStatus:          crosschaintypes.CctxStatus_PendingOutbound,
			zetaMinted:          "7993442360774956232",
		},
		{
			name: "Fake ballot has more votes outbound gets finalized",
			votes: []Vote{
				{voterAddress: "zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
				{voterAddress: "zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
				{voterAddress: "zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
				{voterAddress: "zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
				{voterAddress: "zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
				{voterAddress: "zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
				{voterAddress: "zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
				{voterAddress: "zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: true},
			},
			correctBallotResult: observertypes.BallotStatus_BallotInProgress,
			cctxStatus:          crosschaintypes.CctxStatus_OutboundMined,
			zetaMinted:          "7987124742653889020",
		},
		{
			name: "5 success 5 Failed votes ",
			votes: []Vote{
				{voterAddress: "zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca", voteType: observertypes.VoteType_SuccessObservation, isFakeVote: false},
				{voterAddress: "zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt", voteType: observertypes.VoteType_FailureObservation, isFakeVote: false},
				{voterAddress: "zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4", voteType: observertypes.VoteType_FailureObservation, isFakeVote: false},
				{voterAddress: "zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy", voteType: observertypes.VoteType_FailureObservation, isFakeVote: false},
				{voterAddress: "zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav", voteType: observertypes.VoteType_FailureObservation, isFakeVote: false},
				{voterAddress: "zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t", voteType: observertypes.VoteType_FailureObservation, isFakeVote: false},
			},
			correctBallotResult: observertypes.BallotStatus_BallotInProgress,
			cctxStatus:          crosschaintypes.CctxStatus_PendingOutbound,
			zetaMinted:          "7991636132140714751",
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
			out, err := clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, crosschaincli.CmdShowSend(), []string{cctxIdentifier, "--output", "json"})
			cctx := crosschaintypes.QueryGetCctxResponse{}
			s.NoError(broadcaster.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &cctx))
			s.Assert().Equal(crosschaintypes.CctxStatus_PendingOutbound, cctx.CrossChainTx.CctxStatus.Status)

			fakeVotes := []string{}
			for _, val := range s.network.Validators {
				valVote := Vote{}
				for _, vote := range test.votes {
					if vote.voterAddress == val.Address.String() {
						valVote = vote
					}
				}
				if valVote.voteType == observertypes.VoteType_NotYetVoted {
					continue
				}
				out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetAccountCmd(), []string{val.Address.String(), "--output", "json"})
				var account authtypes.AccountI
				s.NoError(val.ClientCtx.Codec.UnmarshalInterfaceJSON(out.Bytes(), &account))
				outTxhash := test.name
				if valVote.isFakeVote {
					outTxhash = outTxhash + "falseVote"
					fakeVotes = append(fakeVotes, val.Address.String())
				}
				votestring := ""
				switch valVote.voteType {
				case observertypes.VoteType_SuccessObservation:
					votestring = "0"
				case observertypes.VoteType_FailureObservation:
					votestring = "1"
				}

				signedTx := BuildSignedOutboundVote(s.T(), val, s.cfg.BondDenom, account, cctxIdentifier, outTxhash, test.zetaMinted, votestring)
				out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetBroadcastCommand(), []string{signedTx.Name(), "--broadcast-mode", "sync"})
				s.Require().NoError(err)
			}
			s.Require().NoError(s.network.WaitForNBlocks(2))
			out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, crosschaincli.CmdShowSend(), []string{cctxIdentifier, "--output", "json"})
			cctx = crosschaintypes.QueryGetCctxResponse{}
			s.NoError(broadcaster.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &cctx))
			s.Assert().Equal(test.cctxStatus, cctx.CrossChainTx.CctxStatus.Status)
			outboundBallotIdentifier := GetBallotIdentifierOutBound(cctxIdentifier, test.name, test.zetaMinted)
			out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, observercli.CmdBallotByIdentifier(), []string{outboundBallotIdentifier, "--output", "json"})
			s.Require().NoError(err)
			ballot := observertypes.QueryBallotByIdentifierResponse{}
			s.NoError(broadcaster.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &ballot))

			s.Require().Equal(test.correctBallotResult, ballot.BallotStatus)
			for _, vote := range test.votes {
				for _, ballotvote := range ballot.Voters {
					if vote.voterAddress == ballotvote.VoterAddress {
						if !vote.isFakeVote {
							s.Assert().Equal(vote.voteType, ballotvote.VoteType)
						} else {
							s.Assert().Equal(observertypes.VoteType_NotYetVoted, ballotvote.VoteType)
						}
						break
					}
				}
			}
			if len(fakeVotes) > 0 {
				outboundFakeBallotIdentifier := GetBallotIdentifierOutBound(cctxIdentifier, test.name+"falseVote", test.zetaMinted)
				out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, observercli.CmdBallotByIdentifier(), []string{outboundFakeBallotIdentifier, "--output", "json"})
				s.Require().NoError(err)
				fakeBallot := observertypes.QueryBallotByIdentifierResponse{}
				s.NoError(broadcaster.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &fakeBallot))
				for _, vote := range test.votes {
					if vote.isFakeVote {
						for _, ballotvote := range fakeBallot.Voters {
							if vote.voterAddress == ballotvote.VoterAddress {
								s.Assert().Equal(vote.voteType, ballotvote.VoteType)
								break
							}
						}
					}
				}

			}

		})

	}
}
