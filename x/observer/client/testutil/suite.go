package testutil

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/testutil/network"
	"github.com/zeta-chain/zetacore/x/observer/client/cli"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	minOBsDel, ok := sdk.NewIntFromString("10000000000000000")
	s.Require().True(ok)
	s.cfg.StakingTokens = minOBsDel.Mul(sdk.NewInt(10))
	s.cfg.BondedTokens = minOBsDel
	genesisState := s.cfg.GenesisState
	var observerGenesis types.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[types.ModuleName], &observerGenesis))

	observerGenesis.Params.ObserverParams = types.NewParams().ObserverParams

	observerGenesisBz, err := s.cfg.Codec.MarshalJSON(&observerGenesis)
	s.Require().NoError(err)
	genesisState[types.ModuleName] = observerGenesisBz
	s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestAddObserversCmd() {
	val := s.network.Validators[0]
	cmd := cli.CmdAddObserver()
	args := []string{
		types.ObserverChain_Eth.String(),
		types.ObservationType_InBoundTx.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	_, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, args)
	s.Require().NoError(err)

}
