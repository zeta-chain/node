package simulation_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	cosmossimutils "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	cosmossim "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	cosmossimcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/app"
	zetasimulation "github.com/zeta-chain/node/simulation"
)

// AppChainID hardcoded chainID for simulation

func init() {
	zetasimulation.GetSimulatorFlags()
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
	if !zetasimulation.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := zetasimulation.NewConfigFromFlags()

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
	appOptions[server.FlagInvCheckPeriod] = zetasimulation.FlagPeriodValue

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
				zetasimulation.FlagVerboseValue,
				zetasimulation.FlagEnabledValue,
			)
			require.NoError(t, err)
			appOptions[flags.FlagHome] = dir

			simApp, err := zetasimulation.NewSimApp(
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
				zetasimulation.AppStateFn(
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

			zetasimulation.PrintStats(db)

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

	config := zetasimulation.NewConfigFromFlags()

	config.ChainID = SimAppChainID
	config.BlockMaxGas = SimBlockMaxGas
	config.DBBackend = SimDBBackend

	db, dir, logger, skip, err := cosmossimutils.SetupSimulation(
		config,
		SimDBBackend,
		SimDBName,
		zetasimulation.FlagVerboseValue,
		zetasimulation.FlagEnabledValue,
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
	appOptions[server.FlagInvCheckPeriod] = zetasimulation.FlagPeriodValue
	appOptions[flags.FlagHome] = dir

	simApp, err := zetasimulation.NewSimApp(
		logger,
		db,
		appOptions,
		interBlockCacheOpt(),
		baseapp.SetChainID(SimAppChainID),
	)
	require.NoError(t, err)

	blockedAddresses := simApp.ModuleAccountAddrs()
	_, _, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		simApp.BaseApp,
		zetasimulation.AppStateFn(
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

	zetasimulation.PrintStats(db)
}

func TestAppImportExport(t *testing.T) {
	config := zetasimulation.NewConfigFromFlags()

	config.ChainID = SimAppChainID
	config.BlockMaxGas = SimBlockMaxGas
	config.DBBackend = SimDBBackend

	db, dir, logger, skip, err := cosmossimutils.SetupSimulation(
		config,
		SimDBBackend,
		SimDBName,
		zetasimulation.FlagVerboseValue,
		zetasimulation.FlagEnabledValue,
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
	appOptions[server.FlagInvCheckPeriod] = zetasimulation.FlagPeriodValue
	appOptions[flags.FlagHome] = dir
	simApp, err := zetasimulation.NewSimApp(
		logger,
		db,
		appOptions,
		interBlockCacheOpt(),
		baseapp.SetChainID(SimAppChainID),
	)
	require.NoError(t, err)

	// Run randomized simulation
	blockedAddresses := simApp.ModuleAccountAddrs()
	_, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		simApp.BaseApp,
		zetasimulation.AppStateFn(
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

	err = zetasimulation.CheckExportSimulation(simApp, config, simParams)
	require.NoError(t, err)

	zetasimulation.PrintStats(db)

	t.Log("exporting genesis")
	// export state and simParams
	exported, err := simApp.ExportAppStateAndValidators(false, []string{}, []string{})
	require.NoError(t, err)

	newDB, newDir, _, _, err := cosmossimutils.SetupSimulation(
		config,
		SimDBBackend+"_new",
		SimDBName+"_new",
		zetasimulation.FlagVerboseValue,
		zetasimulation.FlagEnabledValue,
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

	newSimApp, err := zetasimulation.NewSimApp(
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
			require.Contains(t, err, "validator set is empty after InitGenesis", "unexpected error: %v", r)
			t.Log("Skipping simulation as all validators have been unbonded")
			t.Log("err", err, "stacktrace", string(debug.Stack()))
		}
	}()

	// Create context for the old and the new sim app, which can be used to compare keys
	ctxSimApp := simApp.NewContext(true, tmproto.Header{
		Height:  simApp.LastBlockHeight(),
		ChainID: SimAppChainID,
	})

	ctxNewSimApp := newSimApp.NewContext(true, tmproto.Header{
		Height:  simApp.LastBlockHeight(),
		ChainID: SimAppChainID,
	})

	// Use genesis state from the first app to initialize the second app
	newSimApp.ModuleManager().InitGenesis(ctxNewSimApp, newSimApp.AppCodec(), genesisState)
	newSimApp.StoreConsensusParams(ctxNewSimApp, exported.ConsensusParams)

	t.Log("comparing stores")

	// The ordering of the keys is not important, we compare the same prefix for both simulations
	storeKeysPrefixes := []StoreKeysPrefixes{
		{simApp.GetKey(authtypes.StoreKey), newSimApp.GetKey(authtypes.StoreKey), [][]byte{}},
		//{
		//	simApp.GetKey(stakingtypes.StoreKey), newSimApp.GetKey(stakingtypes.StoreKey),
		//	[][]byte{
		//		stakingtypes.UnbondingQueueKey, stakingtypes.RedelegationQueueKey, stakingtypes.ValidatorQueueKey,
		//		stakingtypes.HistoricalInfoKey, stakingtypes.UnbondingIDKey, stakingtypes.UnbondingIndexKey, stakingtypes.UnbondingTypeKey, stakingtypes.ValidatorUpdatesKey,
		//	},
		//},
		//{simApp.GetKey(slashingtypes.StoreKey), newSimApp.GetKey(slashingtypes.StoreKey), [][]byte{}},
		//{simApp.GetKey(distrtypes.StoreKey), newSimApp.GetKey(distrtypes.StoreKey), [][]byte{}},
		//{simApp.GetKey(banktypes.StoreKey), newSimApp.GetKey(banktypes.StoreKey), [][]byte{banktypes.BalancesPrefix}},
		//{simApp.GetKey(paramtypes.StoreKey), newSimApp.GetKey(paramtypes.StoreKey), [][]byte{}},
		//{simApp.GetKey(govtypes.StoreKey), newSimApp.GetKey(govtypes.StoreKey), [][]byte{}},
		//{simApp.GetKey(evidencetypes.StoreKey), newSimApp.GetKey(evidencetypes.StoreKey), [][]byte{}},
		//{simApp.GetKey(evmtypes.StoreKey), newSimApp.GetKey(evmtypes.StoreKey), [][]byte{}},
		//{simApp.GetKey(crosschaintypes.StoreKey), newSimApp.GetKey(crosschaintypes.StoreKey), [][]byte{
		//	//crosschaintypes.KeyPrefix(crosschaintypes.CCTXKey),
		//	//crosschaintypes.KeyPrefix(crosschaintypes.LastBlockHeightKey),
		//	//crosschaintypes.KeyPrefix(crosschaintypes.FinalizedInboundsKey),
		//	//crosschaintypes.KeyPrefix(crosschaintypes.GasPriceKey),
		//	//crosschaintypes.KeyPrefix(crosschaintypes.OutboundTrackerKeyPrefix),
		//	//crosschaintypes.KeyPrefix(crosschaintypes.InboundTrackerKeyPrefix),
		//	//crosschaintypes.KeyPrefix(crosschaintypes.ZetaAccountingKey),
		//	//crosschaintypes.KeyPrefix(crosschaintypes.RateLimiterFlagsKey),
		//}},

		//{simApp.GetKey(observertypes.StoreKey), newSimApp.GetKey(observertypes.StoreKey), [][]byte{
		//	//observertypes.KeyPrefix(observertypes.BlameKey),
		//	//observertypes.KeyPrefix(observertypes.VoterKey),
		//	//observertypes.KeyPrefix(observertypes.CrosschainFlagsKey),
		//	//observertypes.KeyPrefix(observertypes.LastBlockObserverCountKey),
		//	//observertypes.KeyPrefix(observertypes.NodeAccountKey),
		//	//observertypes.KeyPrefix(observertypes.KeygenKey),
		//	observertypes.KeyPrefix(observertypes.BallotListKey),
		//	//observertypes.KeyPrefix(observertypes.TSSKey),
		//	//observertypes.KeyPrefix(observertypes.ObserverSetKey),
		//	//observertypes.KeyPrefix(observertypes.AllChainParamsKey),
		//	//observertypes.KeyPrefix(observertypes.TSSHistoryKey),
		//	//observertypes.KeyPrefix(observertypes.TssFundMigratorKey),
		//	//observertypes.KeyPrefix(observertypes.PendingNoncesKeyPrefix),
		//	//observertypes.KeyPrefix(observertypes.ChainNoncesKey),
		//	//observertypes.KeyPrefix(observertypes.NonceToCctxKeyPrefix),
		//	//observertypes.KeyPrefix(observertypes.ParamsKey),
		//}},
		//{simApp.GetKey(fungibletypes.StoreKey), newSimApp.GetKey(fungibletypes.StoreKey), [][]byte{}},
	}

	for _, skp := range storeKeysPrefixes {
		storeA := ctxSimApp.KVStore(skp.A)
		storeB := ctxNewSimApp.KVStore(skp.B)

		failedKVAs, failedKVBs := DiffKVStores(storeA, storeB, skp.Prefixes, simApp.AppCodec())
		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")

		t.Logf("compared %d different key/value pairs between %s and %s\n", len(failedKVAs), skp.A, skp.B)
		require.Equal(
			t,
			0,
			len(failedKVAs),
			//cosmossimutils.GetSimulationLog(
			//	skp.A.Name(),
			//	simApp.SimulationManager().StoreDecoders,
			//	failedKVAs,
			//	failedKVBs,
			//),
		)
	}
}

// DiffKVStores compares two KVstores and returns all the key/value pairs
// that differ from one another. It also skips value comparison for a set of provided prefixes.
func DiffKVStores(a sdk.KVStore, b sdk.KVStore, prefixesToSkip [][]byte, cdc codec.Codec) (kvAs, kvBs []kv.Pair) {
	iterA := a.Iterator(nil, nil)

	defer iterA.Close()

	iterB := b.Iterator(nil, nil)

	defer iterB.Close()

	for {
		if !iterA.Valid() && !iterB.Valid() {
			return kvAs, kvBs
		}

		var kvA, kvB kv.Pair
		if iterA.Valid() {
			kvA = kv.Pair{Key: iterA.Key(), Value: iterA.Value()}

			iterA.Next()
		}

		if iterB.Valid() {
			kvB = kv.Pair{Key: iterB.Key(), Value: iterB.Value()}
		}

		compareValue := true

		unknowprefix := false
		for _, prefix := range prefixesToSkip {
			// Skip value comparison if we matched a prefix
			if bytes.HasPrefix(kvA.Key, prefix) {
				compareValue = false
				break
			}
			unknowprefix = true
		}

		if !compareValue {
			// We're skipping this key due to an exclusion prefix.  If it's present in B, iterate past it.  If it's
			// absent don't iterate.
			if bytes.Equal(kvA.Key, kvB.Key) {
				iterB.Next()
			}
			continue
		}

		// always iterate B when comparing
		iterB.Next()

		if !bytes.Equal(kvA.Value, kvB.Value) {
			fmt.Println("Value mismatch", unknowprefix)
			fmt.Println("A", string(kvA.Key))
			fmt.Println("B", string(kvB.Key))
			fmt.Println("A", string(kvA.Value))
			fmt.Println("B", string(kvB.Value))
			fmt.Println("-------------------------------------------------------------")
		}

		if !bytes.Equal(kvA.Key, kvB.Key) {
			fmt.Println("Key mismatch", unknowprefix)
		}

		if !bytes.Equal(kvA.Key, kvB.Key) || !bytes.Equal(kvA.Value, kvB.Value) {
			kvAs = append(kvAs, kvA)
			kvBs = append(kvBs, kvB)
		}
	}
}

func TestAppSimulationAfterImport(t *testing.T) {
	config := zetasimulation.NewConfigFromFlags()

	config.ChainID = SimAppChainID
	config.BlockMaxGas = SimBlockMaxGas
	config.DBBackend = SimDBBackend

	db, dir, logger, skip, err := cosmossimutils.SetupSimulation(
		config,
		SimDBBackend,
		SimDBName,
		zetasimulation.FlagVerboseValue,
		zetasimulation.FlagEnabledValue,
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
	appOptions[server.FlagInvCheckPeriod] = zetasimulation.FlagPeriodValue
	appOptions[flags.FlagHome] = dir
	simApp, err := zetasimulation.NewSimApp(
		logger,
		db,
		appOptions,
		interBlockCacheOpt(),
		baseapp.SetChainID(SimAppChainID),
	)
	require.NoError(t, err)

	// Run randomized simulation
	blockedAddresses := simApp.ModuleAccountAddrs()
	stopEarly, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		simApp.BaseApp,
		zetasimulation.AppStateFn(
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

	err = zetasimulation.CheckExportSimulation(simApp, config, simParams)
	require.NoError(t, err)

	zetasimulation.PrintStats(db)

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
		zetasimulation.FlagVerboseValue,
		zetasimulation.FlagEnabledValue,
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

	newSimApp, err := zetasimulation.NewSimApp(
		logger,
		newDB,
		appOptions,
		interBlockCacheOpt(),
		baseapp.SetChainID(SimAppChainID),
	)
	require.NoError(t, err)

	t.Log("Adding app state to new app")
	newSimApp.InitChain(abci.RequestInitChain{
		ChainId:       SimAppChainID,
		AppStateBytes: exported.AppState,
	})

	t.Log("Simulating new simulation")
	stopEarly, simParams, simErr = simulation.SimulateFromSeed(
		t,
		os.Stdout,
		newSimApp.BaseApp,
		zetasimulation.AppStateFn(
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
