package querytests

import (
	"math/rand"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcfg "github.com/evmos/ethermint/cmd/config"
	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/app"
	cmdcfg "github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/testutil/network"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

type CliTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
	ballots []*observerTypes.Ballot
}

func NewCLITestSuite(cfg network.Config) *CliTestSuite {
	return &CliTestSuite{cfg: cfg}
}

func (s *CliTestSuite) Setconfig() {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	ethcfg.SetBip44CoinType(config)
	// Make sure address is compatible with ethereum
	config.SetAddressVerifier(app.VerifyAddressFormat)
	config.Seal()
}
func (s *CliTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")
	s.Setconfig()
	minOBsDel, ok := sdk.NewIntFromString("100000000000000000000")
	s.Require().True(ok)
	s.cfg.StakingTokens = minOBsDel.Mul(sdk.NewInt(int64(10)))
	s.cfg.BondedTokens = minOBsDel
	observerList := []string{"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax",
		"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2",
		"zeta1szrskhdeleyt6wmn0nfxvcvt2l6f4fn06uaga4",
		"zeta16h3y7s7030l4chcznwq3n6uz2m9wvmzu5vwt7c",
		"zeta1xl2rfsrmx8nxryty3lsjuxwdxs59cn2q65e5ca",
		"zeta1ktmprjdvc72jq0mpu8tn8sqx9xwj685qx0q6kt",
		"zeta1ygeyr8pqfjvclxay5234gulnjzv2mkz6lph9y4",
		"zeta1zegyenj7xg5nck04ykkzndm2qxdzc6v83mklsy",
		"zeta1us2qpqdcctk6q7qv2c9d9jvjxlv88jscf68kav",
		"zeta1e9fyaulgntkrnqnl0es4nyxghp3petpn2ntu3t",
	}
	network.SetupZetaGenesisState(s.T(), s.cfg.GenesisState, s.cfg.Codec, observerList, false)
	s.ballots = RandomBallotGenerator(20, observerList)
	network.AddObserverData(s.T(), 2, s.cfg.GenesisState, s.cfg.Codec, s.ballots)

	net, err := network.New(s.T(), app.NodeDir, s.cfg)
	s.Assert().NoError(err)
	s.network = net
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

}

func CreateRandomVoteList(numberOfVotes int) []observerTypes.VoteType {
	voteOptions := []observerTypes.VoteType{observerTypes.VoteType_SuccessObservation, observerTypes.VoteType_FailureObservation, observerTypes.VoteType_NotYetVoted}
	min := 0
	max := len(voteOptions) - 1
	voteList := make([]observerTypes.VoteType, numberOfVotes)
	for i := 0; i < numberOfVotes; i++ {
		voteList[i] = voteOptions[rand.Intn(max-min)+min] // #nosec G404
	}
	return voteList
}
func RandomBallotGenerator(numberOfBallots int, voterList []string) []*observerTypes.Ballot {
	ballots := make([]*observerTypes.Ballot, numberOfBallots)
	ballotStatus := []observerTypes.BallotStatus{observerTypes.BallotStatus_BallotFinalized_FailureObservation, observerTypes.BallotStatus_BallotFinalized_SuccessObservation}
	// #nosec G404 randomness is not a security issue here
	for i := 0; i < numberOfBallots; i++ {
		ballots[i] = &observerTypes.Ballot{
			Index:            "",
			BallotIdentifier: "TestBallot" + strconv.Itoa(i),
			VoterList:        voterList,
			Votes:            CreateRandomVoteList(len(voterList)),
			ObservationType:  observerTypes.ObservationType_InBoundTx,
			BallotThreshold:  sdk.MustNewDecFromStr("0.66"),
			// #nosec G404 randomness used for testing
			BallotStatus:         ballotStatus[rand.Intn(2)],
			BallotCreationHeight: 0,
		}
	}
	return ballots
}

func (s *CliTestSuite) TearDownSuite() {
	s.T().Log("tearing down genesis test suite")
	s.network.Cleanup()
}
