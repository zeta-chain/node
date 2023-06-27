//go:build TESTNET
// +build TESTNET

package testutil

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethcfg "github.com/evmos/ethermint/cmd/config"
	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/app"
	cmdcfg "github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/network"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) Setconfig() {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	ethcfg.SetBip44CoinType(config)
	// Make sure address is compatible with ethereum
	config.SetAddressVerifier(app.VerifyAddressFormat)
	config.Seal()
}
func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")
	s.Setconfig()
	minOBsDel, ok := sdk.NewIntFromString("100000000000000000000")
	s.Require().True(ok)
	s.cfg.StakingTokens = minOBsDel.Mul(sdk.NewInt(int64(10)))
	s.cfg.BondedTokens = minOBsDel
	genesisState := s.cfg.GenesisState

	// Cross-chain genesis state
	var crossChainGenesis types.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[types.ModuleName], &crossChainGenesis))
	crossChainGenesis.Params.Enabled = true
	crossChainGenesisBz, err := s.cfg.Codec.MarshalJSON(&crossChainGenesis)
	s.Require().NoError(err)

	// Staking genesis state
	var stakingGenesisState stakingtypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenesisState))
	stakingGenesisState.Params.BondDenom = cmdcfg.BaseDenom
	stakingGenesisStateBz, err := s.cfg.Codec.MarshalJSON(&stakingGenesisState)
	s.Require().NoError(err)

	// Observer genesis state
	var observerGenesis observerTypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[observerTypes.ModuleName], &observerGenesis))
	var observerMapper []*observerTypes.ObserverMapper
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
	for _, chain := range common.DefaultChainsList() {
		observerMapper = append(observerMapper, &observerTypes.ObserverMapper{
			ObserverChain: chain,
			ObserverList:  observerList,
		})
	}
	observerGenesis.Observers = observerMapper
	observerGenesisBz, err := s.cfg.Codec.MarshalJSON(&observerGenesis)
	s.Require().NoError(err)

	genesisState[types.ModuleName] = crossChainGenesisBz
	genesisState[stakingtypes.ModuleName] = stakingGenesisStateBz
	genesisState[observerTypes.ModuleName] = observerGenesisBz
	s.cfg.GenesisState = genesisState

	s.network, err = network.New(s.T(), app.NodeDir, s.cfg)
	s.Assert().NoError(err)
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}
