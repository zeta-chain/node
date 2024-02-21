package integrationtests

import (
	"encoding/json"
	"fmt"
	"strings"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	crosschaincli "github.com/zeta-chain/zetacore/x/crosschain/client/cli"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observercli "github.com/zeta-chain/zetacore/x/observer/client/cli"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

type messageLog struct {
	Events []event `json:"events"`
}

type event struct {
	Type       string      `json:"type"`
	Attributes []attribute `json:"attributes"`
}

type attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// fetchAttribute fetches the attribute from the tx response
func fetchAttribute(rawLog string, key string) (string, error) {
	var logs []messageLog
	err := json.Unmarshal([]byte(rawLog), &logs)
	if err != nil {
		return "", err
	}

	var attributes []string
	for _, log := range logs {
		for _, event := range log.Events {
			for _, attr := range event.Attributes {
				attributes = append(attributes, attr.Key)
				if strings.EqualFold(attr.Key, key) {
					address := attr.Value

					// trim the quotes
					address = address[1 : len(address)-1]

					return address, nil
				}

			}
		}
	}

	return "", fmt.Errorf("attribute %s not found, attributes:  %+v", key, attributes)
}

type txRes struct {
	RawLog string `json:"raw_log"`
}

func ExtractRawLog(str string) (string, error) {
	var data txRes

	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return "", err
	}

	return data.RawLog, nil
}

func (s *IntegrationTestSuite) TestCCTXInboundVoter() {
	broadcaster := s.network.Validators[0]

	var systemContractAddr string
	// Initialize system contract
	{
		out, err := clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetAccountCmd(), []string{broadcaster.Address.String(), "--output", "json"})
		s.Require().NoError(err)
		var account authtypes.AccountI
		s.NoError(broadcaster.ClientCtx.Codec.UnmarshalInterfaceJSON(out.Bytes(), &account))
		signedTx := BuildSignedDeploySystemContract(s.T(), broadcaster, s.cfg.BondDenom, account)
		res, err := clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetBroadcastCommand(), []string{signedTx.Name(), "--broadcast-mode", "block"})
		s.Require().NoError(err)

		rawLog, err := ExtractRawLog(res.String())
		s.Require().NoError(err)

		systemContractAddr, err = fetchAttribute(rawLog, "system_contract")
		s.Require().NoError(err)

		// update system contract
		out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetAccountCmd(), []string{broadcaster.Address.String(), "--output", "json"})
		s.Require().NoError(err)
		s.NoError(broadcaster.ClientCtx.Codec.UnmarshalInterfaceJSON(out.Bytes(), &account))
		signedTx = BuildSignedUpdateSystemContract(s.T(), broadcaster, s.cfg.BondDenom, account, systemContractAddr)
		res, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetBroadcastCommand(), []string{signedTx.Name(), "--broadcast-mode", "block"})
		s.Require().NoError(err)
	}

	// Deploy ETH ZRC20
	{
		out, err := clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetAccountCmd(), []string{broadcaster.Address.String(), "--output", "json"})
		s.Require().NoError(err)
		var account authtypes.AccountI
		s.NoError(broadcaster.ClientCtx.Codec.UnmarshalInterfaceJSON(out.Bytes(), &account))
		signedTx := BuildSignedDeployETHZRC20(s.T(), broadcaster, s.cfg.BondDenom, account)
		_, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetBroadcastCommand(), []string{signedTx.Name(), "--broadcast-mode", "block"})
		s.Require().NoError(err)
	}

	tt := []struct {
		name                  string
		votes                 map[string]observerTypes.VoteType
		ballotResult          observerTypes.BallotStatus
		cctxStatus            crosschaintypes.CctxStatus
		falseBallotIdentifier string
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
			ballotResult: observerTypes.BallotStatus_BallotFinalized_SuccessObservation,
			cctxStatus:   crosschaintypes.CctxStatus_PendingOutbound,
		},
		{
			name: "5 votes only ballot does not get finalized",
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
			ballotResult: observerTypes.BallotStatus_BallotInProgress,
			cctxStatus:   crosschaintypes.CctxStatus_PendingRevert,
		},
		{
			name: "1 false vote but correct ballot is still finalized",
			votes: map[string]observerTypes.VoteType{
				"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax": observerTypes.VoteType_SuccessObservation,
				"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2": observerTypes.VoteType_SuccessObservation,
				"zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4": observerTypes.VoteType_SuccessObservation,
				"zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c": observerTypes.VoteType_SuccessObservation,
				"zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca": observerTypes.VoteType_SuccessObservation,
				"zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt": observerTypes.VoteType_SuccessObservation,
				"zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4": observerTypes.VoteType_SuccessObservation,
				"zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy": observerTypes.VoteType_FailureObservation,
				"zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav": observerTypes.VoteType_SuccessObservation,
				"zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t": observerTypes.VoteType_NotYetVoted,
			},
			ballotResult: observerTypes.BallotStatus_BallotFinalized_SuccessObservation,
			cctxStatus:   crosschaintypes.CctxStatus_PendingOutbound,
		},
		{
			name: "2 ballots with 5 votes each no ballot gets finalized",
			votes: map[string]observerTypes.VoteType{
				"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax": observerTypes.VoteType_SuccessObservation,
				"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2": observerTypes.VoteType_SuccessObservation,
				"zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4": observerTypes.VoteType_SuccessObservation,
				"zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c": observerTypes.VoteType_SuccessObservation,
				"zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca": observerTypes.VoteType_SuccessObservation,
				"zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt": observerTypes.VoteType_FailureObservation,
				"zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4": observerTypes.VoteType_FailureObservation,
				"zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy": observerTypes.VoteType_FailureObservation,
				"zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav": observerTypes.VoteType_FailureObservation,
				"zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t": observerTypes.VoteType_FailureObservation,
			},
			ballotResult: observerTypes.BallotStatus_BallotInProgress,
			cctxStatus:   crosschaintypes.CctxStatus_PendingRevert,
		},
		{
			name: "majority wrong votes incorrect ballot finalized / correct ballot still in progress",
			votes: map[string]observerTypes.VoteType{
				"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax": observerTypes.VoteType_SuccessObservation,
				"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2": observerTypes.VoteType_SuccessObservation,
				"zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4": observerTypes.VoteType_SuccessObservation,
				"zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c": observerTypes.VoteType_FailureObservation,
				"zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca": observerTypes.VoteType_FailureObservation,
				"zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt": observerTypes.VoteType_FailureObservation,
				"zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4": observerTypes.VoteType_FailureObservation,
				"zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy": observerTypes.VoteType_FailureObservation,
				"zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav": observerTypes.VoteType_FailureObservation,
				"zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t": observerTypes.VoteType_FailureObservation,
			},
			ballotResult:          observerTypes.BallotStatus_BallotInProgress,
			cctxStatus:            crosschaintypes.CctxStatus_PendingOutbound,
			falseBallotIdentifier: "majority wrong votes incorrect ballot finalized / correct ballot still in progress" + "falseVote",
		},
		{
			name: "7 votes only just crossed threshold",
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
			ballotResult: observerTypes.BallotStatus_BallotFinalized_SuccessObservation,
			cctxStatus:   crosschaintypes.CctxStatus_PendingOutbound,
		},
	}
	for i, test := range tt {
		test := test
		s.Run(test.name, func() {
			// Vote the gas price
			for _, val := range s.network.Validators {
				out, err := clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetAccountCmd(), []string{val.Address.String(), "--output", "json"})
				s.Require().NoError(err)

				var account authtypes.AccountI
				s.NoError(val.ClientCtx.Codec.UnmarshalInterfaceJSON(out.Bytes(), &account))
				signedTx := BuildSignedGasPriceVote(s.T(), val, s.cfg.BondDenom, account)
				_, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetBroadcastCommand(), []string{signedTx.Name(), "--broadcast-mode", "sync"})
				s.Require().NoError(err)
			}

			s.Require().NoError(s.network.WaitForNBlocks(2))
			out, err := clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, observercli.CmdListPendingNonces(), []string{"--output", "json"})
			s.Require().NoError(err)
			out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, observercli.CmdGetSupportedChains(), []string{"--output", "json"})
			s.Require().NoError(err)
			out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, crosschaincli.CmdListGasPrice(), []string{"--output", "json"})
			s.Require().NoError(err)

			// Vote the inbound tx
			for _, val := range s.network.Validators {
				vote := test.votes[val.Address.String()]
				if vote == observerTypes.VoteType_NotYetVoted {
					continue
				}
				out, err := clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetAccountCmd(), []string{val.Address.String(), "--output", "json"})
				var account authtypes.AccountI
				s.NoError(val.ClientCtx.Codec.UnmarshalInterfaceJSON(out.Bytes(), &account))

				message := test.name
				if vote == observerTypes.VoteType_FailureObservation {
					message = message + "falseVote"
				}
				signedTx := BuildSignedInboundVote(s.T(), val, s.cfg.BondDenom, account, message, i)
				out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, authcli.GetBroadcastCommand(), []string{signedTx.Name(), "--broadcast-mode", "block"})
				s.Require().NoError(err)
				fmt.Println(out.String())
			}
			s.Require().NoError(s.network.WaitForNBlocks(2))

			// Get the ballot
			ballotIdentifier := GetBallotIdentifier(test.name, i)
			out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, observercli.CmdBallotByIdentifier(), []string{ballotIdentifier, "--output", "json"})
			s.Require().NoError(err)
			ballot := observerTypes.QueryBallotByIdentifierResponse{}
			s.NoError(broadcaster.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &ballot))

			// Check the vote in the ballot
			s.Require().Equal(len(test.votes), len(ballot.Voters))
			for _, vote := range ballot.Voters {
				if test.votes[vote.VoterAddress] == observerTypes.VoteType_FailureObservation {
					s.Assert().Equal(observerTypes.VoteType_NotYetVoted.String(), vote.VoteType.String())
					continue
				}
				s.Assert().Equal(test.votes[vote.VoterAddress].String(), vote.VoteType.String(), "incorrect vote for voter: %s", vote.VoterAddress)
			}
			s.Require().Equal(test.ballotResult.String(), ballot.BallotStatus.String())

			// Get the cctx and check its status
			cctxIdentifier := ballotIdentifier
			if test.falseBallotIdentifier != "" {
				cctxIdentifier = GetBallotIdentifier(test.falseBallotIdentifier, i)
			}
			out, err = clitestutil.ExecTestCLICmd(broadcaster.ClientCtx, crosschaincli.CmdShowSend(), []string{cctxIdentifier, "--output", "json"})
			cctx := crosschaintypes.QueryGetCctxResponse{}
			if test.cctxStatus == crosschaintypes.CctxStatus_PendingRevert {
				s.Require().Contains(out.String(), "not found")
			} else {
				s.NoError(broadcaster.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &cctx))
				s.Require().Equal(test.cctxStatus.String(), cctx.CrossChainTx.CctxStatus.Status.String(), cctx.CrossChainTx.CctxStatus.StatusMessage)
			}
		})
	}

}
