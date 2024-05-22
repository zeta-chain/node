package querytests

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/suite"

	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/testutil/network"
)

func TestCLIQuerySuite(t *testing.T) {
	cfg := network.DefaultConfig(NewTestNetworkFixture)
	suite.Run(t, NewCLITestSuite(cfg))
}

func NewTestNetworkFixture() network.TestFixture {
	encoding := app.MakeEncodingConfig()
	appCtr := func(val network.ValidatorI) servertypes.Application {
		return app.New(
			val.GetCtx().Logger, tmdb.NewMemDB(), nil, true, map[int64]bool{}, val.GetCtx().Config.RootDir, 0,
			encoding,
			simtestutil.EmptyAppOptions{},
			baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			baseapp.SetChainID("athens_8888-2"),
		)
	}

	return network.TestFixture{
		AppConstructor: appCtr,
		GenesisState:   app.ModuleBasics.DefaultGenesis(encoding.Codec),
		EncodingConfig: testutil.TestEncodingConfig{
			InterfaceRegistry: encoding.InterfaceRegistry,
			Codec:             encoding.Codec,
			TxConfig:          encoding.TxConfig,
			Amino:             encoding.Amino,
		},
	}
}
