package app_test

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

	dbm "github.com/cometbft/cometbft-db"

	"github.com/cosmos/cosmos-sdk/store"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	simulation2 "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"

	// "github.com/cosmos/gaia/v11/app/helpers"
	// "github.com/cosmos/gaia/v11/app/params"
	"github.com/zeta-chain/node/app/sim"
)

// AppChainID hardcoded chainID for simulation

func init() {
	sim.GetSimulatorFlags()
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
	if !sim.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := sim.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false
	config.ChainID = SimAppChainID

	numSeeds := 3
	numTimesToRunPerSeed := 5

	// We will be overriding the random seed and just run a single simulation on the provided seed value
	if config.Seed != simcli.DefaultSeedValue {
		numSeeds = 1
	}

	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)
	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[server.FlagInvCheckPeriod] = sim.FlagPeriodValue

	for i := 0; i < numSeeds; i++ {
		if config.Seed == simcli.DefaultSeedValue {
			config.Seed = rand.Int63()
		}

		fmt.Println("config.Seed: ", config.Seed)

		for j := 0; j < numTimesToRunPerSeed; j++ {
			var logger log.Logger
			if sim.FlagVerboseValue {
				logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
			} else {
				logger = log.NewNopLogger()
			}

			db := dbm.NewMemDB()
			dir, err := os.MkdirTemp("", "zeta-simulation")
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
				sim.AppStateFn(app.AppCodec(), app.SimulationManager(), app.ModuleBasics.DefaultGenesis(app.AppCodec())),
				simulation2.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
				simtestutil.SimulationOperations(app, app.AppCodec(), config),
				blockedAddresses,
				config,
				app.AppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				sim.PrintStats(db)
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
	config := sim.NewConfigFromFlags()
	config.ChainID = SimAppChainID
	config.BlockMaxGas = SimBlockMaxGas
	config.DBBackend = "memdb"
	//config.ExportStatePath = "/Users/tanmay/.zetacored/simulation_state_export.json"

	db, dir, logger, skip, err := simtestutil.SetupSimulation(config, "mem-db", "Simulation", sim.FlagVerboseValue, sim.FlagEnabledValue)
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		require.NoError(t, db.Close())
		require.NoError(t, os.RemoveAll(dir))
	}()
	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[server.FlagInvCheckPeriod] = sim.FlagPeriodValue

	app, err := NewSimApp(logger, db, appOptions, interBlockCacheOpt(), baseapp.SetChainID(SimAppChainID))
	require.NoError(t, err)

	blockedAddresses := app.ModuleAccountAddrs()
	_, _, simerr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		app.BaseApp,
		sim.AppStateFn(app.AppCodec(), app.SimulationManager(), app.ModuleBasics.DefaultGenesis(app.AppCodec())),
		simulation2.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
		simtestutil.SimulationOperations(app, app.AppCodec(), config),
		blockedAddresses,
		config,
		app.AppCodec(),
	)
	require.NoError(t, simerr)

	// check export works as expected
	_, err = app.ExportAppStateAndValidators(false, nil, nil)
	require.NoError(t, err)

	if config.Commit {
		simtestutil.PrintStats(db)
	}
}
