package app

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ModuleBasics defines the module BasicManager that is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
func newBasicManagerFromManager(app *App) module.BasicManager {
	var moduleBasics []module.AppModuleBasic
	for _, m := range app.mm.Modules {
		m, ok := m.(module.AppModuleBasic)
		if !ok {
			fmt.Printf("module %s is not an instance of module.AppModuleBasic\n", m)
			continue
		}
		if m.Name() == govtypes.ModuleName {
			m = gov.NewAppModuleBasic(getGovProposalHandlers())
		}

		if m.Name() == genutiltypes.ModuleName {
			m = genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator)
		}

		moduleBasics = append(moduleBasics, m)
	}
	basicManager := module.NewBasicManager(moduleBasics...)
	//basicManager.RegisterLegacyAminoCodec(app.cdc)
	basicManager.RegisterInterfaces(app.interfaceRegistry)

	return basicManager
}

func simulationModules(
	app *App,
	appCodec codec.Codec,
	_ bool,
) []module.AppModuleSimulation {
	return []module.AppModuleSimulation{
		auth.NewAppModule(appCodec, app.AccountKeeper, RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		slashing.NewAppModule(
			appCodec,
			app.SlashingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.GetSubspace(slashingtypes.ModuleName),
		),
		params.NewAppModule(app.ParamsKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
	}
}

// RandomGenesisAccounts defines the default RandomGenesisAccountsFn used on the SDK.
// It creates a slice of BaseAccount, ContinuousVestingAccount and DelayedVestingAccount.
func RandomGenesisAccounts(simState *module.SimulationState) authtypes.GenesisAccounts {
	genesisAccs := make(authtypes.GenesisAccounts, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		bacc := authtypes.NewBaseAccountWithAddress(acc.Address)

		// Only consider making a vesting account once the initial bonded validator
		// set is exhausted due to needing to track DelegatedVesting.
		if !(int64(i) > simState.NumBonded && simState.Rand.Intn(100) < 50) {
			genesisAccs[i] = bacc
			continue
		}

		initialVesting := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, simState.Rand.Int63n(simState.InitialStake.Int64())))
		var endTime int64

		startTime := simState.GenTimestamp.Unix()

		// Allow for some vesting accounts to vest very quickly while others very slowly.
		if simState.Rand.Intn(100) < 50 {
			endTime = int64(simulation.RandIntBetween(simState.Rand, int(startTime)+1, int(startTime+(60*60*24*30))))
		} else {
			endTime = int64(simulation.RandIntBetween(simState.Rand, int(startTime)+1, int(startTime+(60*60*12))))
		}

		bva := vestingtypes.NewBaseVestingAccount(bacc, initialVesting, endTime)

		if simState.Rand.Intn(100) < 50 {
			genesisAccs[i] = vestingtypes.NewContinuousVestingAccountRaw(bva, startTime)
		} else {
			genesisAccs[i] = vestingtypes.NewDelayedVestingAccountRaw(bva)
		}
	}

	return genesisAccs
}
