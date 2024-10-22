package simulation_test

import (
	"os"
	"testing"

	cosmossim "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/stretchr/testify/require"
	simutils "github.com/zeta-chain/node/tests/simulation/sim"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	cosmossimutils "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

func BenchmarkFullAppSimulation(b *testing.B) {
	b.ReportAllocs()

	config := simutils.NewConfigFromFlags()

	config.ChainID = SimAppChainID
	config.BlockMaxGas = SimBlockMaxGas
	config.DBBackend = SimDBBackend

	db, dir, logger, skip, err := cosmossimutils.SetupSimulation(
		config,
		SimDBBackend,
		SimDBName,
		simutils.FlagVerboseValue,
		simutils.FlagEnabledValue,
	)
	if skip {
		b.Skip("skipping application simulation")
	}
	require.NoError(b, err, "simulation setup failed")

	defer func() {
		require.NoError(b, db.Close())
		require.NoError(b, os.RemoveAll(dir))
	}()

	appOptions := make(cosmossimutils.AppOptionsMap, 0)
	appOptions[server.FlagInvCheckPeriod] = simutils.FlagPeriodValue
	appOptions[flags.FlagHome] = dir
	simApp, err := simutils.NewSimApp(logger, db, appOptions, interBlockCacheOpt(), baseapp.SetChainID(SimAppChainID))
	require.NoError(b, err)

	// Run randomized simulation:
	blockedAddresses := simApp.ModuleAccountAddrs()
	_, simParams, simErr := simulation.SimulateFromSeed(
		b,
		os.Stdout,
		simApp.BaseApp,
		simutils.AppStateFn(
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
	require.NoError(b, simErr)

	// export state and simParams before the simulation error is checked
	err = simutils.CheckExportSimulation(simApp, config, simParams)
	require.NoError(b, err)

	if simErr != nil {
		b.Fatal(simErr)
	}

	simutils.PrintStats(db)
}
