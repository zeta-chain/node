package simulation_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"strings"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"
	"github.com/zeta-chain/node/simulation/sim"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/app"

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
	sim.GetSimulatorFlags()
}

type StoreKeysPrefixes struct {
	A        storetypes.StoreKey
	B        storetypes.StoreKey
	Prefixes [][]byte
}

const (
	SimAppChainID  = "simulation_777-1"
	SimBlockMaxGas = 815000000
	//github.com/zeta-chain/node/issues/3004
	// TODO : Support pebbleDB for simulation tests
	SimDBBackend = "goleveldb"
	SimDBName    = "simulation"
)

// interBlockCacheOpt returns a BaseApp option function that sets the persistent
// inter-block write-through cache.
func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}

// TestAppStateDeterminism runs a full application simulation , and produces multiple blocks as per the config
// It checks the determinism of the application by comparing the apphash at the end of each run to other runs
// The following test certifies that , for the same set of operations ( irrespective of what the operations are ) ,
// we would reach the same final state if the initial state is the same
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
	config.DBBackend = SimDBBackend
	config.BlockMaxGas = SimBlockMaxGas

	numSeeds := 3
	numTimesToRunPerSeed := 5

	// We will be overriding the random seed and just run a single simulation on the provided seed value
	if config.Seed != cosmossimcli.DefaultSeedValue {
		numSeeds = 1
	}

	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	appOptions := make(cosmossimutils.AppOptionsMap, 0)
	appOptions[server.FlagInvCheckPeriod] = sim.FlagPeriodValue

	t.Log("Running tests for numSeeds: ", numSeeds, " numTimesToRunPerSeed: ", numTimesToRunPerSeed)

	for i := 0; i < numSeeds; i++ {
		if config.Seed == cosmossimcli.DefaultSeedValue {
			config.Seed = rand.Int63()
		}
		// For the same seed, the simApp hash produced at the end of each run should be the same
		for j := 0; j < numTimesToRunPerSeed; j++ {
			db, dir, logger, _, err := cosmossimutils.SetupSimulation(
				config,
				SimDBBackend,
				SimDBName,
				sim.FlagVerboseValue,
				sim.FlagEnabledValue,
			)
			require.NoError(t, err)
			appOptions[flags.FlagHome] = dir

			simApp, err := sim.NewSimApp(
				logger,
				db,
				appOptions,
				interBlockCacheOpt(),
				baseapp.SetChainID(SimAppChainID),
			)

			t.Logf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			blockedAddresses := simApp.ModuleAccountAddrs()

			// Random seed is used to produce a random initial state for the simulation
			_, _, err = simulation.SimulateFromSeed(
				t,
				os.Stdout,
				simApp.BaseApp,
				sim.AppStateFn(
					simApp.AppCodec(),
					simApp.SimulationManager(),
					simApp.BasicManager().DefaultGenesis(simApp.AppCodec()),
				),
				cosmossim.RandomAccounts,
				cosmossimutils.SimulationOperations(simApp, simApp.AppCodec(), config),
				blockedAddresses,
				config,
				simApp.AppCodec(),
			)
			require.NoError(t, err)

			sim.PrintStats(db)

			appHash := simApp.LastCommitID().Hash
			appHashList[j] = appHash

			// Clean up resources
			t.Cleanup(func() {
				require.NoError(t, db.Close())
				require.NoError(t, os.RemoveAll(dir))
			})

			if j != 0 {
				require.Equal(
					t,
					string(appHashList[0]),
					string(appHashList[j]),
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n",
					config.Seed,
					i+1,
					numSeeds,
					j+1,
					numTimesToRunPerSeed,
				)
			}
		}
	}
}

// TestFullAppSimulation runs a full simApp simulation with the provided configuration.
// At the end of the run it tries to export the genesis state to make sure the export works.
func TestFullAppSimulation(t *testing.T) {

	config := sim.NewConfigFromFlags()

	config.ChainID = SimAppChainID
	config.BlockMaxGas = SimBlockMaxGas
	config.DBBackend = SimDBBackend

	db, dir, logger, skip, err := cosmossimutils.SetupSimulation(
		config,
		SimDBBackend,
		SimDBName,
		sim.FlagVerboseValue,
		sim.FlagEnabledValue,
	)
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			require.NoError(t, err, "Error closing new database")
		}
		if err := os.RemoveAll(dir); err != nil {
			require.NoError(t, err, "Error removing directory")
		}
	})
	appOptions := make(cosmossimutils.AppOptionsMap, 0)
	appOptions[server.FlagInvCheckPeriod] = sim.FlagPeriodValue
	appOptions[flags.FlagHome] = dir

	simApp, err := sim.NewSimApp(logger, db, appOptions, interBlockCacheOpt(), baseapp.SetChainID(SimAppChainID))
	require.NoError(t, err)

	blockedAddresses := simApp.ModuleAccountAddrs()
	_, _, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		simApp.BaseApp,
		sim.AppStateFn(
			simApp.AppCodec(),
			simApp.SimulationManager(),
			simApp.BasicManager().DefaultGenesis(simApp.AppCodec()),
		),
		cosmossim.RandomAccounts,
		cosmossimutils.SimulationOperations(simApp, simApp.AppCodec(), config),
		blockedAddresses,
		config,
		simApp.AppCodec(),
	)
	require.NoError(t, simErr)

	// check export works as expected
	exported, err := simApp.ExportAppStateAndValidators(false, nil, nil)
	require.NoError(t, err)
	if config.ExportStatePath != "" {
		err := os.WriteFile(config.ExportStatePath, exported.AppState, 0o600)
		require.NoError(t, err)
	}

	sim.PrintStats(db)
}

func TestAppImportExport(t *testing.T) {
	config := sim.NewConfigFromFlags()

	config.ChainID = SimAppChainID
	config.BlockMaxGas = SimBlockMaxGas
	config.DBBackend = SimDBBackend

	db, dir, logger, skip, err := cosmossimutils.SetupSimulation(
		config,
		SimDBBackend,
		SimDBName,
		sim.FlagVerboseValue,
		sim.FlagEnabledValue,
	)
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			require.NoError(t, err, "Error closing new database")
		}
		if err := os.RemoveAll(dir); err != nil {
			require.NoError(t, err, "Error removing directory")
		}
	})

	appOptions := make(cosmossimutils.AppOptionsMap, 0)
	appOptions[server.FlagInvCheckPeriod] = sim.FlagPeriodValue
	appOptions[flags.FlagHome] = dir
	simApp, err := sim.NewSimApp(logger, db, appOptions, interBlockCacheOpt(), baseapp.SetChainID(SimAppChainID))
	require.NoError(t, err)

	// Run randomized simulation
	blockedAddresses := simApp.ModuleAccountAddrs()
	_, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		simApp.BaseApp,
		sim.AppStateFn(
			simApp.AppCodec(),
			simApp.SimulationManager(),
			simApp.BasicManager().DefaultGenesis(simApp.AppCodec()),
		),
		cosmossim.RandomAccounts,
		cosmossimutils.SimulationOperations(simApp, simApp.AppCodec(), config),
		blockedAddresses,
		config,
		simApp.AppCodec(),
	)
	require.NoError(t, simErr)

	err = sim.CheckExportSimulation(simApp, config, simParams)
	require.NoError(t, err)

	sim.PrintStats(db)

	t.Log("exporting genesis")
	// export state and simParams
	exported, err := simApp.ExportAppStateAndValidators(false, []string{}, []string{})
	require.NoError(t, err)

	t.Log("importing genesis")
	newDB, newDir, _, _, err := cosmossimutils.SetupSimulation(
		config,
		SimDBBackend+"_new",
		SimDBName+"_new",
		sim.FlagVerboseValue,
		sim.FlagEnabledValue,
	)

	require.NoError(t, err, "simulation setup failed")

	t.Cleanup(func() {
		if err := newDB.Close(); err != nil {
			require.NoError(t, err, "Error closing new database")
		}
		if err := os.RemoveAll(newDir); err != nil {
			require.NoError(t, err, "Error removing directory")
		}
	})

	newSimApp, err := sim.NewSimApp(
		logger,
		newDB,
		appOptions,
		interBlockCacheOpt(),
		baseapp.SetChainID(SimAppChainID),
	)
	require.NoError(t, err)

	var genesisState app.GenesisState
	err = json.Unmarshal(exported.AppState, &genesisState)
	require.NoError(t, err)

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%v", r)
			if !strings.Contains(err, "validator set is empty after InitGenesis") {
				panic(r)
			}
			t.Log("Skipping simulation as all validators have been unbonded")
			t.Log("err", err, "stacktrace", string(debug.Stack()))
		}
	}()

	// Create context for the old and the new sim app, which can be used to compare keys
	ctxSimApp := simApp.NewContext(true, tmproto.Header{
		Height:  simApp.LastBlockHeight(),
		ChainID: SimAppChainID,
	}).WithChainID(SimAppChainID)
	ctxNewSimApp := newSimApp.NewContext(true, tmproto.Header{
		Height:  simApp.LastBlockHeight(),
		ChainID: SimAppChainID,
	}).WithChainID(SimAppChainID)

	// Use genesis state from the first app to initialize the second app
	newSimApp.ModuleManager().InitGenesis(ctxNewSimApp, newSimApp.AppCodec(), genesisState)
	newSimApp.StoreConsensusParams(ctxNewSimApp, exported.ConsensusParams)

	t.Log("comparing stores")
	// The ordering of the keys is not important, we compare the same prefix for both simulations
	storeKeysPrefixes := []StoreKeysPrefixes{
		{simApp.GetKey(authtypes.StoreKey), newSimApp.GetKey(authtypes.StoreKey), [][]byte{}},
		{
			simApp.GetKey(stakingtypes.StoreKey), newSimApp.GetKey(stakingtypes.StoreKey),
			[][]byte{
				stakingtypes.UnbondingQueueKey, stakingtypes.RedelegationQueueKey, stakingtypes.ValidatorQueueKey,
				stakingtypes.HistoricalInfoKey, stakingtypes.UnbondingIDKey, stakingtypes.UnbondingIndexKey, stakingtypes.UnbondingTypeKey, stakingtypes.ValidatorUpdatesKey,
			},
		},
		{simApp.GetKey(slashingtypes.StoreKey), newSimApp.GetKey(slashingtypes.StoreKey), [][]byte{}},
		{simApp.GetKey(distrtypes.StoreKey), newSimApp.GetKey(distrtypes.StoreKey), [][]byte{}},
		{simApp.GetKey(banktypes.StoreKey), newSimApp.GetKey(banktypes.StoreKey), [][]byte{banktypes.BalancesPrefix}},
		{simApp.GetKey(paramtypes.StoreKey), newSimApp.GetKey(paramtypes.StoreKey), [][]byte{}},
		{simApp.GetKey(govtypes.StoreKey), newSimApp.GetKey(govtypes.StoreKey), [][]byte{}},
		{simApp.GetKey(evidencetypes.StoreKey), newSimApp.GetKey(evidencetypes.StoreKey), [][]byte{}},
		{simApp.GetKey(evmtypes.StoreKey), newSimApp.GetKey(evmtypes.StoreKey), [][]byte{}},
	}

	for _, skp := range storeKeysPrefixes {
		storeA := ctxSimApp.KVStore(skp.A)
		storeB := ctxNewSimApp.KVStore(skp.B)

		failedKVAs, failedKVBs := sdk.DiffKVStores(storeA, storeB, skp.Prefixes)
		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")

		t.Logf("compared %d different key/value pairs between %s and %s\n", len(failedKVAs), skp.A, skp.B)
		require.Equal(
			t,
			0,
			len(failedKVAs),
			cosmossimutils.GetSimulationLog(
				skp.A.Name(),
				simApp.SimulationManager().StoreDecoders,
				failedKVAs,
				failedKVBs,
			),
		)
	}
}

func TestAppSimulationAfterImport(t *testing.T) {
	config := sim.NewConfigFromFlags()

	config.ChainID = SimAppChainID
	config.BlockMaxGas = SimBlockMaxGas
	config.DBBackend = SimDBBackend

	db, dir, logger, skip, err := cosmossimutils.SetupSimulation(
		config,
		SimDBBackend,
		SimDBName,
		sim.FlagVerboseValue,
		sim.FlagEnabledValue,
	)
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			require.NoError(t, err, "Error closing new database")
		}
		if err := os.RemoveAll(dir); err != nil {
			require.NoError(t, err, "Error removing directory")
		}
	})

	appOptions := make(cosmossimutils.AppOptionsMap, 0)
	appOptions[server.FlagInvCheckPeriod] = sim.FlagPeriodValue
	appOptions[flags.FlagHome] = dir
	simApp, err := sim.NewSimApp(logger, db, appOptions, interBlockCacheOpt(), baseapp.SetChainID(SimAppChainID))
	require.NoError(t, err)

	// Run randomized simulation
	blockedAddresses := simApp.ModuleAccountAddrs()
	stopEarly, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		simApp.BaseApp,
		sim.AppStateFn(
			simApp.AppCodec(),
			simApp.SimulationManager(),
			simApp.BasicManager().DefaultGenesis(simApp.AppCodec()),
		),
		cosmossim.RandomAccounts,
		cosmossimutils.SimulationOperations(simApp, simApp.AppCodec(), config),
		blockedAddresses,
		config,
		simApp.AppCodec(),
	)
	require.NoError(t, simErr)

	err = sim.CheckExportSimulation(simApp, config, simParams)
	require.NoError(t, err)

	sim.PrintStats(db)

	if stopEarly {
		t.Log("can't export or import a zero-validator genesis, exiting test")
		return
	}

	t.Log("exporting genesis")

	// export state and simParams
	exported, err := simApp.ExportAppStateAndValidators(true, []string{}, []string{})
	require.NoError(t, err)

	t.Log("importing genesis")

	newDB, newDir, _, _, err := cosmossimutils.SetupSimulation(
		config,
		SimDBBackend+"_new",
		SimDBName+"_new",
		sim.FlagVerboseValue,
		sim.FlagEnabledValue,
	)

	require.NoError(t, err, "simulation setup failed")

	t.Cleanup(func() {
		if err := newDB.Close(); err != nil {
			require.NoError(t, err, "Error closing new database")
		}
		if err := os.RemoveAll(newDir); err != nil {
			require.NoError(t, err, "Error removing directory")
		}
	})

	newSimApp, err := sim.NewSimApp(
		logger,
		newDB,
		appOptions,
		interBlockCacheOpt(),
		baseapp.SetChainID(SimAppChainID),
	)
	require.NoError(t, err)

	newSimApp.InitChain(abci.RequestInitChain{
		ChainId:       SimAppChainID,
		AppStateBytes: exported.AppState,
	})

	stopEarly, simParams, simErr = simulation.SimulateFromSeed(
		t,
		os.Stdout,
		newSimApp.BaseApp,
		sim.AppStateFn(
			simApp.AppCodec(),
			simApp.SimulationManager(),
			simApp.BasicManager().DefaultGenesis(simApp.AppCodec()),
		),
		cosmossim.RandomAccounts,
		cosmossimutils.SimulationOperations(newSimApp, newSimApp.AppCodec(), config),
		blockedAddresses,
		config,
		simApp.AppCodec(),
	)
	require.NoError(t, err)
}
