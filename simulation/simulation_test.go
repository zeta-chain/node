package simulation_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"sync"
	"testing"

	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	evidencetypes "cosmossdk.io/x/evidence/types"
	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	cosmossimutils "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/kv"
	cosmossim "github.com/cosmos/cosmos-sdk/types/simulation"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	cosmossimcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"
	"github.com/zeta-chain/node/app"
	zetasimulation "github.com/zeta-chain/node/simulation"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// AppChainID hardcoded chainID for simulation

func init() {
	zetasimulation.GetSimulatorFlags()
}

type StoreKeysPrefixes struct {
	A            storetypes.StoreKey
	B            storetypes.StoreKey
	SkipPrefixes [][]byte
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
// It does the following
// 1. It runs the simulation multiple times with the same seed value
// 2. It checks the apphash at the end of each run
// 3. It compares the apphash at the end of each run to check for determinism
// 4. Repeat steps 1-3 for multiple seeds

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
					t,
					simApp.AppCodec(),
					simApp.SimulationManager(),
					simApp.BasicManager().DefaultGenesis(simApp.AppCodec()),
					nil,
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
// This test does the following
// 1. It runs a full simulation with the provided configuration
// 2. It exports the state and validators
// 3. Verifies that the run and export were successful
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
			t,
			simApp.AppCodec(),
			simApp.SimulationManager(),
			simApp.BasicManager().DefaultGenesis(simApp.AppCodec()),
			nil,
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

// TestAppImportExport tests the application simulation after importing the state exported from a previous.At a high level,it does the following
//  1. It runs a full simulation and exports the state
//  2. It creates a new app, and db
//  3. It imports the exported state into the new app
//  4. It compares the key value pairs for the two apps.The comparison function takes a list of keys to skip as an input as well
//     a. First app which ran the simulation
//     b. Second app which imported the state

// This can verify the export and import process do not modify the state in anyway irrespective of the operations performed
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
			t,
			simApp.AppCodec(),
			simApp.SimulationManager(),
			simApp.BasicManager().DefaultGenesis(simApp.AppCodec()),
			nil,
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
	ctxSimApp := simApp.NewContext(true).WithBlockHeight(simApp.LastBlockHeight()).WithChainID(SimAppChainID)

	ctxNewSimApp := newSimApp.NewContext(true).WithBlockHeight(simApp.LastBlockHeight()).WithChainID(SimAppChainID)

	// Use genesis state from the first app to initialize the second app
	newSimApp.ModuleManager().InitGenesis(ctxNewSimApp, newSimApp.AppCodec(), genesisState)
	newSimApp.StoreConsensusParams(ctxNewSimApp, exported.ConsensusParams)

	t.Log("comparing stores")

	// The ordering of the keys is not important, we compare the same prefix for both simulations
	storeKeysPrefixes := []StoreKeysPrefixes{
		// Interaction with EVM module,
		// such as deploying contracts or interacting with them such as setting gas price,
		// causes the state for the auth module to change on export.The order of keys within the store is modified.
		// We will need to explore this further to find a definitive answer
		// TODO:https://github.com/zeta-chain/node/issues/3263

		//{simApp.GetKey(authtypes.StoreKey), newSimApp.GetKey(authtypes.StoreKey), [][]byte{}},
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
		{simApp.GetKey(crosschaintypes.StoreKey), newSimApp.GetKey(crosschaintypes.StoreKey), [][]byte{
			// We update the timestamp for cctx when importing the genesis state which results in a different value
			crosschaintypes.KeyPrefix(crosschaintypes.CCTXKey),
		}},

		{simApp.GetKey(observertypes.StoreKey), newSimApp.GetKey(observertypes.StoreKey), [][]byte{
			// The order of ballots when importing is not preserved which causes the value to be different.
			observertypes.KeyPrefix(observertypes.BallotListKey),
		}},
		{simApp.GetKey(fungibletypes.StoreKey), newSimApp.GetKey(fungibletypes.StoreKey), [][]byte{}},
	}

	for _, skp := range storeKeysPrefixes {
		storeA := ctxSimApp.KVStore(skp.A)
		storeB := ctxNewSimApp.KVStore(skp.B)

		failedKVAs, failedKVBs := DiffKVStores(storeA, storeB, skp.SkipPrefixes)
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

// TestAppSimulationAfterImport tests the application simulation after importing the state exported from a previous simulation run.
// It does the following steps
// 1. It runs a full simulation and exports the state
// 2. It creates a new app, and db
// 3. It imports the exported state into the new app
// 4. It runs a simulation on the new app and verifies that there is no error in the second simulation
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
			t,
			simApp.AppCodec(),
			simApp.SimulationManager(),
			simApp.BasicManager().DefaultGenesis(simApp.AppCodec()),
			nil,
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

	// Setup a new app with new database and directory
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

	// Initialize the new app with the exported genesis state of the first run
	t.Log("Importing genesis into the new app")
	newSimApp.InitChain(&abci.RequestInitChain{
		ChainId:       SimAppChainID,
		AppStateBytes: exported.AppState,
	})

	// Run simulation on the new app
	stopEarly, simParams, simErr = simulation.SimulateFromSeed(
		t,
		os.Stdout,
		newSimApp.BaseApp,
		zetasimulation.AppStateFn(
			t,
			nil,
			nil,
			nil,
			exported.AppState,
		),
		cosmossim.RandomAccounts,
		cosmossimutils.SimulationOperations(newSimApp, newSimApp.AppCodec(), config),
		blockedAddresses,
		config,
		simApp.AppCodec(),
	)
	require.NoError(t, simErr)
}

// that differ from one another. It also skips value comparison for a set of provided prefixes.
func DiffKVStores(a, b storetypes.KVStore, prefixesToSkip [][]byte) (diffA, diffB []kv.Pair) {
	iterA := a.Iterator(nil, nil)
	defer iterA.Close()

	iterB := b.Iterator(nil, nil)
	defer iterB.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	kvAs := make([]kv.Pair, 0)
	go func() {
		defer wg.Done()
		kvAs = getKVPairs(iterA, prefixesToSkip)
	}()

	wg.Add(1)
	kvBs := make([]kv.Pair, 0)
	go func() {
		defer wg.Done()
		kvBs = getKVPairs(iterB, prefixesToSkip)
	}()

	wg.Wait()

	if len(kvAs) != len(kvBs) {
		fmt.Printf(
			"KV stores are different: %d key/value pairs in store A and %d key/value pairs in store B\n",
			len(kvAs),
			len(kvBs),
		)
		// Find out which keypair is missing in store A
		for _, kvB := range kvBs {
			found := false
			for _, kvA := range kvAs {
				if bytes.Equal(kvA.Key, kvB.Key) {
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("Key %s is missing in store A\n", kvB.Key)
			}
		}
	}

	return getDiffFromKVPair(kvAs, kvBs)
}

// getDiffFromKVPair compares two KVstores and returns all the key/value pairs
func getDiffFromKVPair(kvAs, kvBs []kv.Pair) (diffA, diffB []kv.Pair) {
	// we assume that kvBs is equal or larger than kvAs
	// if not, we swap the two
	if len(kvAs) > len(kvBs) {
		kvAs, kvBs = kvBs, kvAs
		// we need to swap the diffA and diffB as well
		defer func() {
			diffA, diffB = diffB, diffA
		}()
	}

	// in case kvAs is empty we can return early
	// since there is nothing to compare
	// if kvAs == kvBs, then diffA and diffB will be empty
	if len(kvAs) == 0 {
		return []kv.Pair{}, kvBs
	}

	index := make(map[string][]byte, len(kvBs))
	for _, kv := range kvBs {
		index[string(kv.Key)] = kv.Value
	}

	for _, kvA := range kvAs {
		if kvBValue, ok := index[string(kvA.Key)]; !ok {
			diffA = append(diffA, kvA)
			diffB = append(
				diffB,
				kv.Pair{Key: kvA.Key},
			) // the key is missing from kvB so we append a pair with an empty value
		} else if !bytes.Equal(kvA.Value, kvBValue) {
			diffA = append(diffA, kvA)
			diffB = append(diffB, kv.Pair{Key: kvA.Key, Value: kvBValue})
		} else {
			// values are equal, so we remove the key from the index
			delete(index, string(kvA.Key))
		}
	}

	// add the remaining keys from kvBs
	for key, value := range index {
		diffA = append(
			diffA,
			kv.Pair{Key: []byte(key)},
		) // the key is missing from kvA so we append a pair with an empty value
		diffB = append(diffB, kv.Pair{Key: []byte(key), Value: value})
	}

	return diffA, diffB
}

func getKVPairs(iter dbm.Iterator, prefixesToSkip [][]byte) (kvs []kv.Pair) {
	for iter.Valid() {
		key, value := iter.Key(), iter.Value()

		// do not add the KV pair if the key is prefixed to be skipped.
		skip := false
		for _, prefix := range prefixesToSkip {
			if bytes.HasPrefix(key, prefix) {
				skip = true
				break
			}
		}

		if !skip {
			kvs = append(kvs, kv.Pair{Key: key, Value: value})
		}

		iter.Next()
	}

	return kvs
}
