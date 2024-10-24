package app

import (
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
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
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
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/zeta-chain/ethermint/x/evm"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"
	"github.com/zeta-chain/ethermint/x/feemarket"
	feemarkettypes "github.com/zeta-chain/ethermint/x/feemarket/types"

	authoritymodule "github.com/zeta-chain/node/x/authority"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	crosschainmodule "github.com/zeta-chain/node/x/crosschain"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionsmodule "github.com/zeta-chain/node/x/emissions"
	emissionsModuleTypes "github.com/zeta-chain/node/x/emissions/types"
	fungiblemodule "github.com/zeta-chain/node/x/fungible"
	fungibleModuleTypes "github.com/zeta-chain/node/x/fungible/types"
	lightclientmodule "github.com/zeta-chain/node/x/lightclient"
	lightclienttypes "github.com/zeta-chain/node/x/lightclient/types"
	observermodule "github.com/zeta-chain/node/x/observer"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// ModuleBasics defines the module BasicManager that is in charge of setting up basic,
// non-dependant module elements, such as codec registration
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
	evm.AppModuleBasic{},
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

// InitGenesisModuleList returns the module list for genesis initialization
// NOTE: Capability module must occur first so that it can initialize any capabilities
// TODO: enable back IBC
// all commented lines in this function are modules related to IBC
// https://github.com/zeta-chain/node/issues/2573
func InitGenesisModuleList() []string {
	return []string{
		//capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		//ibcexported.ModuleName,
		//ibctransfertypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		paramstypes.ModuleName,
		group.ModuleName,
		genutiltypes.ModuleName,
		upgradetypes.ModuleName,
		evidencetypes.ModuleName,
		vestingtypes.ModuleName,
		observertypes.ModuleName,
		crosschaintypes.ModuleName,
		//ibccrosschaintypes.ModuleName,
		fungibleModuleTypes.ModuleName,
		emissionsModuleTypes.ModuleName,
		authz.ModuleName,
		authoritytypes.ModuleName,
		lightclienttypes.ModuleName,
		consensusparamtypes.ModuleName,
		// NOTE: crisis module must go at the end to check for invariants on each module
		crisistypes.ModuleName,
	}
}

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
		),
		params.NewAppModule(app.ParamsKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper, app.GetSubspace(evmtypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		groupmodule.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		crosschainmodule.NewAppModule(appCodec, app.CrosschainKeeper),
	}
}
