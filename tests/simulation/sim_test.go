package simulation_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/ethermint/app"
	evmante "github.com/zeta-chain/ethermint/app/ante"
	zetaapp "github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/app/ante"
	simutils "github.com/zeta-chain/node/tests/simulation/sim"

	dbm "github.com/cometbft/cometbft-db"

	"github.com/cosmos/cosmos-sdk/store"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	cosmossimutils "github.com/cosmos/cosmos-sdk/testutil/sims"
	cosmossim "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	cosmossimcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
)

// AppChainID hardcoded chainID for simulation

func init() {
	simutils.GetSimulatorFlags()
}

const (
	SimAppChainID  = "simulation_777-1"
	SimBlockMaxGas = 815000000
	TestAppChainID = "zetachain_777-1"
)

// NewSimApp disable feemarket on native tx, otherwise the cosmos-sdk simulation tests will fail.
func NewSimApp(logger log.Logger, db dbm.DB, appOptions servertypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp)) (*zetaapp.App, error) {

	encCdc := zetaapp.MakeEncodingConfig()
	app := zetaapp.New(
		logger,
		db,
		nil,
		false,
		map[int64]bool{},
		app.DefaultNodeHome,
		5,
		encCdc,
		appOptions,
		baseAppOptions...,
	)
	sdk.DefaultPowerReduction = sdk.OneInt()
	// disable feemarket on native tx
	options := ante.HandlerOptions{
		AccountKeeper:   app.AccountKeeper,
		BankKeeper:      app.BankKeeper,
		EvmKeeper:       app.EvmKeeper,
		FeeMarketKeeper: app.FeeMarketKeeper,
		SignModeHandler: encCdc.TxConfig.SignModeHandler(),
		SigGasConsumer:  evmante.DefaultSigVerificationGasConsumer,
		MaxTxGasWanted:  0,
		ObserverKeeper:  app.ObserverKeeper,
	}

	anteHandler, err := ante.NewAnteHandler(options)
	if err != nil {
		panic(err)
	}

	app.SetAnteHandler(anteHandler)
	if err := app.LoadLatestVersion(); err != nil {
		return nil, err
	}
	return app, nil
}

// interBlockCacheOpt returns a BaseApp option function that sets the persistent
// inter-block write-through cache.
func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}

// TODO: Make another test for the fuzzer itself, which just has noOp txs
// and doesn't depend on the application.
func TestAppStateDeterminism(t *testing.T) {
	if !simutils.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simutils.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false
	config.ChainID = SimAppChainID
	config.DBBackend = "goleveldb"

	numSeeds := 3
	numTimesToRunPerSeed := 5

	// We will be overriding the random seed and just run a single simulation on the provided seed value
	if config.Seed != cosmossimcli.DefaultSeedValue {
		numSeeds = 1
	}

	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)
	appOptions := make(cosmossimutils.AppOptionsMap, 0)
	appOptions[server.FlagInvCheckPeriod] = simutils.FlagPeriodValue

	for i := 0; i < numSeeds; i++ {
		if config.Seed == cosmossimcli.DefaultSeedValue {
			config.Seed = rand.Int63()
		}

		fmt.Println("config.Seed: ", config.Seed)

		for j := 0; j < numTimesToRunPerSeed; j++ {
			var logger log.Logger
			if simutils.FlagVerboseValue {
				logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
			} else {
				logger = log.NewNopLogger()
			}

			db, dir, logger, skip, err := cosmossimutils.SetupSimulation(config, "level-db", "Simulation", simutils.FlagVerboseValue, simutils.FlagEnabledValue)
			if skip {
				t.Skip("skipping application simulation")
			}
			require.NoError(t, err)
			appOptions[flags.FlagHome] = dir

			app, err := NewSimApp(logger, db, appOptions, interBlockCacheOpt(), baseapp.SetChainID(SimAppChainID))

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			blockedAddresses := app.ModuleAccountAddrs()

			_, _, err = simulation.SimulateFromSeed(
				t,
				os.Stdout,
				app.BaseApp,
				simutils.AppStateFn(app.AppCodec(), app.SimulationManager(), app.ModuleBasics.DefaultGenesis(app.AppCodec())),
				cosmossim.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
				cosmossimutils.SimulationOperations(app, app.AppCodec(), config),
				blockedAddresses,
				config,
				app.AppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				simutils.PrintStats(db)
			}

			appHash := app.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t, string(appHashList[0]), string(appHashList[j]),
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}

func TestFullAppSimulation(t *testing.T) {
	config := simutils.NewConfigFromFlags()
	config.ChainID = SimAppChainID
	config.BlockMaxGas = SimBlockMaxGas
	config.DBBackend = "goleveldb"
	//config.ExportStatePath = "/Users/tanmay/.zetacored/simulation_state_export.json"

	db, dir, logger, skip, err := cosmossimutils.SetupSimulation(config, "level-db", "Simulation", simutils.FlagVerboseValue, simutils.FlagEnabledValue)
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		require.NoError(t, db.Close())
		require.NoError(t, os.RemoveAll(dir))
	}()
	appOptions := make(cosmossimutils.AppOptionsMap, 0)
	appOptions[server.FlagInvCheckPeriod] = simutils.FlagPeriodValue

	app, err := NewSimApp(logger, db, appOptions, interBlockCacheOpt(), baseapp.SetChainID(SimAppChainID))
	require.NoError(t, err)

	blockedAddresses := app.ModuleAccountAddrs()
	_, _, simerr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		app.BaseApp,
		simutils.AppStateFn(app.AppCodec(), app.SimulationManager(), app.ModuleBasics.DefaultGenesis(app.AppCodec())),
		cosmossim.RandomAccounts,
		cosmossimutils.SimulationOperations(app, app.AppCodec(), config),
		blockedAddresses,
		config,
		app.AppCodec(),
	)
	require.NoError(t, simerr)

	// check export works as expected
	_, err = app.ExportAppStateAndValidators(false, nil, nil)
	require.NoError(t, err)

	if config.Commit {
		cosmossimutils.PrintStats(db)
	}
}
