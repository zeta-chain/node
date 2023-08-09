package querytests

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcfg "github.com/evmos/ethermint/cmd/config"
	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/app"
	cmdcfg "github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/testutil/network"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

type CliTestSuite struct {
	suite.Suite

	cfg             network.Config
	network         *network.Network
	crosschainState *types.GenesisState
	observerState   *observerTypes.GenesisState
}

func NewCLITestSuite(cfg network.Config) *CliTestSuite {
	return &CliTestSuite{cfg: cfg}
}

func (s *CliTestSuite) Setconfig() {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	ethcfg.SetBip44CoinType(config)
	// Make sure the address is compatible with ethereum
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
	}
	network.SetupZetaGenesisState(s.T(), s.cfg.GenesisState, s.cfg.Codec, observerList)
	s.crosschainState = network.AddCrosschainData(s.T(), 2, s.cfg.GenesisState, s.cfg.Codec)
	s.observerState = network.AddObserverData(s.T(), s.cfg.GenesisState, s.cfg.Codec, nil)
	net, err := network.New(s.T(), app.NodeDir, s.cfg)
	s.Assert().NoError(err)
	s.network = net
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

}

func (s *CliTestSuite) TearDownSuite() {
	s.T().Log("tearing down genesis test suite")
	s.network.Cleanup()
}
