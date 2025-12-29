package app

import (
	"cosmossdk.io/x/evidence"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/evm/x/feemarket"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	"github.com/cosmos/evm/x/vm"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	authoritymodule "github.com/zeta-chain/node/x/authority"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	crosschainmodule "github.com/zeta-chain/node/x/crosschain"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionsmodule "github.com/zeta-chain/node/x/emissions"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	fungiblemodule "github.com/zeta-chain/node/x/fungible"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	lightclientmodule "github.com/zeta-chain/node/x/lightclient"
	lightclienttypes "github.com/zeta-chain/node/x/lightclient/types"
	observermodule "github.com/zeta-chain/node/x/observer"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// ModuleBasics defines the module BasicManager that is in charge of setting up basic,
// non-dependent module elements, such as codec registration
// and genesis verification.
// https://github.com/zeta-chain/node/issues/3021
// TODO: Use app.mm to create the basic manager instead
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
	bank.AppModuleBasic{},
	//capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	distr.AppModuleBasic{},
	gov.NewAppModuleBasic(getGovProposalHandlers()),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	//ibc.AppModuleBasic{},
	//ibctm.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	//transfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	consensus.AppModuleBasic{},
	vm.AppModuleBasic{},
	feemarket.AppModuleBasic{},
	authoritymodule.AppModuleBasic{},
	lightclientmodule.AppModuleBasic{},
	crosschainmodule.AppModuleBasic{},
	//ibccrosschain.AppModuleBasic{},
	observermodule.AppModuleBasic{},
	fungiblemodule.AppModuleBasic{},
	emissionsmodule.AppModuleBasic{},
	groupmodule.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
)

// simulationModules returns a list of modules to include in the simulation
func simulationModules(
	app *App,
	appCodec codec.Codec,
	_ bool,
) []module.AppModuleSimulation {
	return []module.AppModuleSimulation{
		auth.NewAppModule(
			appCodec,
			app.AccountKeeper,
			authsims.RandomGenesisAccounts,
			app.GetSubspace(authtypes.ModuleName),
		),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		// Todo : Enable gov module simulation
		// https://github.com/zeta-chain/node/issues/3007
		//gov.NewAppModule(
		//	appCodec,
		//	&app.GovKeeper,
		//	app.AccountKeeper,
		//	app.BankKeeper,
		//	app.GetSubspace(govtypes.ModuleName),
		//),
		staking.NewAppModule(
			appCodec,
			app.StakingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(stakingtypes.ModuleName),
		),
		distr.NewAppModule(
			appCodec,
			app.DistrKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.GetSubspace(distrtypes.ModuleName),
		),
		slashing.NewAppModule(
			appCodec,
			app.SlashingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.GetSubspace(slashingtypes.ModuleName),
			app.interfaceRegistry,
		),
		params.NewAppModule(app.ParamsKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		vm.NewAppModule(app.EvmKeeper, app.AccountKeeper, app.AccountKeeper.AddressCodec()),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		groupmodule.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		crosschainmodule.NewAppModule(appCodec, app.CrosschainKeeper),
		observermodule.NewAppModule(appCodec, *app.ObserverKeeper),
		fungiblemodule.NewAppModule(appCodec, app.FungibleKeeper),
		emissionsmodule.NewAppModule(appCodec, app.EmissionsKeeper, app.GetSubspace(emissionstypes.ModuleName)),
	}
}

// Order user by cosmos/gaia
// https://github.com/cosmos/gaia/blob/main/app/modules.go

// OrderInitGenesis returns the module list for genesis initialization
// NOTE: Capability module must occur first so that it can initialize any capabilities
// TODO: enable back IBC
// all commented lines in this function are modules related to IBC
// https://github.com/zeta-chain/node/issues/2573
func OrderInitGenesis() []string {
	return []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		evmtypes.ModuleName,
		// Comments from cosmos
		// The feemarket module should ideally be initialized before the genutil module in theory:
		// The feemarket antehandler performs checks in DeliverTx, which is called by gentx.
		// When the fee > 0, gentx needs to pay the fee. However, this is not expected.
		// To resolve this issue, we should initialize the feemarket module after genutil, ensuring that the
		// min fee is empty when gentx is called.
		// A similar issue existed for the 'globalfee' module, which was previously used instead of 'feemarket'.
		// For more details, please refer to the following link: https://github.com/cosmos/gaia/issues/2489

		// https://github.com/zeta-chain/node/issues/3791
		feemarkettypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		group.ModuleName,
		observertypes.ModuleName,
		crosschaintypes.ModuleName,
		fungibletypes.ModuleName,
		emissionstypes.ModuleName,
		authoritytypes.ModuleName,
		lightclienttypes.ModuleName,
		consensusparamtypes.ModuleName,
		// crisis needs to be last so that the genesis state is consistent
		// when it checks invariants
		crisistypes.ModuleName,
	}
}

// During begin block slashing happens after distr.BeginBlocker so that
// there is nothing left over in the validator fee pool, so as to keep the
// CanWithdrawInvariant invariant.
// NOTE: staking module is required if HistoricalEntries param > 0
func orderBeginBlockers() []string {
	return []string{
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		feemarkettypes.ModuleName,
		evmtypes.ModuleName,
		group.ModuleName,
		crosschaintypes.ModuleName,
		observertypes.ModuleName,
		fungibletypes.ModuleName,
		emissionstypes.ModuleName,
		authoritytypes.ModuleName,
		lightclienttypes.ModuleName,
		consensusparamtypes.ModuleName,
	}
}

func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		group.ModuleName,
		crosschaintypes.ModuleName,
		observertypes.ModuleName,
		fungibletypes.ModuleName,
		emissionstypes.ModuleName,
		authoritytypes.ModuleName,
		lightclienttypes.ModuleName,
		consensusparamtypes.ModuleName,
	}
}
