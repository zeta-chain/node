package app

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmos "github.com/cometbft/cometbft/libs/os"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmosencoding "github.com/cosmos/evm/encoding"
	cosmosevmtypes "github.com/cosmos/evm/types"
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/cosmos/evm/x/feemarket"
	feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	"github.com/cosmos/evm/x/vm"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	_ "github.com/ethereum/go-ethereum/eth/tracers/native" // register native tracers
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	"github.com/zeta-chain/node/app/ante"
	"github.com/zeta-chain/node/docs/openapi"
	srvflags "github.com/zeta-chain/node/server/flags"
	authoritymodule "github.com/zeta-chain/node/x/authority"
	authoritykeeper "github.com/zeta-chain/node/x/authority/keeper"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	crosschainmodule "github.com/zeta-chain/node/x/crosschain"
	crosschainkeeper "github.com/zeta-chain/node/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionsmodule "github.com/zeta-chain/node/x/emissions"
	emissionskeeper "github.com/zeta-chain/node/x/emissions/keeper"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	fungiblemodule "github.com/zeta-chain/node/x/fungible"
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	lightclientmodule "github.com/zeta-chain/node/x/lightclient"
	lightclientkeeper "github.com/zeta-chain/node/x/lightclient/keeper"
	lightclienttypes "github.com/zeta-chain/node/x/lightclient/types"
	observermodule "github.com/zeta-chain/node/x/observer"
	observerkeeper "github.com/zeta-chain/node/x/observer/keeper"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// TODO: enable back IBC
// IBC has been turned off for v19, all necessary code has been commented out
// to enable IBC, uncomment the following imports and all logic using these packages in the code
// https://github.com/zeta-chain/node/issues/2573
// "github.com/cosmos/cosmos-sdk/x/capability"
// capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
// capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
// "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
// transferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
// transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
// ibccore "github.com/cosmos/ibc-go/v8/modules/core"
// ibctypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
// ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
// ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
// "github.com/zeta-chain/node/x/ibccrosschain"
// ibccrosschainkeeper "github.com/zeta-chain/node/x/ibccrosschain/keeper"
// ibccrosschaintypes "github.com/zeta-chain/node/x/ibccrosschain/types"

const Name = "zetacore"

func init() {
	// manually update the power reduction by replacing micro (u) -> atto (a) evmos
	sdk.DefaultPowerReduction = cosmosevmtypes.AttoPowerReduction
	// modify fee market parameter defaults through global
	//feemarkettypes.DefaultMinGasPrice = v5.MainnetMinGasPrices
	//feemarkettypes.DefaultMinGasMultiplier = v5.MainnetMinGasMultiplier
}

var (
	NodeDir                    = ".zetacored"
	DefaultNodeHome            = os.ExpandEnv("$HOME/") + NodeDir
	TransactionGasLimit uint64 = 10_000_000
)

func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler
	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
		//ibcclientclient.UpdateClientProposalHandler,
		//ibcclientclient.UpgradeProposalHandler,
	)
	return govProposalHandlers
}

type GenesisState map[string]json.RawMessage

// module account permissions
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	distrtypes.ModuleName:          nil,
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:            {authtypes.Burner},
	//ibctransfertypes.ModuleName:                     {authtypes.Minter, authtypes.Burner},
	crosschaintypes.ModuleName: {authtypes.Minter, authtypes.Burner},
	//ibccrosschaintypes.ModuleName:                   nil,
	evmtypes.ModuleName:                             {authtypes.Minter, authtypes.Burner},
	fungibletypes.ModuleName:                        {authtypes.Minter, authtypes.Burner},
	emissionstypes.ModuleName:                       nil,
	emissionstypes.UndistributedObserverRewardsPool: nil,
	emissionstypes.UndistributedTSSRewardsPool:      nil,
	feemarkettypes.ModuleName:                       nil,
}

// module accounts that are NOT allowed to receive tokens
var blockedReceivingModAcc = map[string]bool{
	distrtypes.ModuleName:          true,
	authtypes.FeeCollectorName:     true,
	stakingtypes.BondedPoolName:    true,
	stakingtypes.NotBondedPoolName: true,
	govtypes.ModuleName:            true,
	evmtypes.ModuleName:            true,
	feemarkettypes.ModuleName:      true,
}

var (
	_ runtime.AppI            = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry
	invCheckPeriod    uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	mm           *module.Manager
	sm           *module.SimulationManager
	mb           module.BasicManager
	configurator module.Configurator

	// sdk keepers
	AccountKeeper authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	//CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper  *stakingkeeper.Keeper
	SlashingKeeper slashingkeeper.Keeper
	DistrKeeper    distrkeeper.Keeper
	GovKeeper      govkeeper.Keeper
	CrisisKeeper   crisiskeeper.Keeper
	UpgradeKeeper  *upgradekeeper.Keeper
	ParamsKeeper   paramskeeper.Keeper
	//IBCKeeper             *ibckeeper.Keeper
	//TransferKeeper        ibctransferkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	GroupKeeper           groupkeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper

	// scoped keepers
	//ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	//ScopedTransferKeeper      capabilitykeeper.ScopedKeeper
	//ScopedIBCCrosschainKeeper capabilitykeeper.ScopedKeeper

	// evm keepers
	EvmKeeper       *evmkeeper.Keeper
	Erc20Keeper     erc20keeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper

	// zetachain keepers
	AuthorityKeeper   authoritykeeper.Keeper
	LightclientKeeper lightclientkeeper.Keeper
	CrosschainKeeper  crosschainkeeper.Keeper
	//IBCCrosschainKeeper ibccrosschainkeeper.Keeper
	ObserverKeeper  *observerkeeper.Keeper
	FungibleKeeper  fungiblekeeper.Keeper
	EmissionsKeeper emissionskeeper.Keeper

	//transferModule transfer.AppModule
}

// New returns a reference to an initialized ZetaApp.
func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	evmChainID uint64,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	encodingConfig := evmosencoding.MakeConfig(evmChainID)
	appCodec := encodingConfig.Codec
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(Name, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(encodingConfig.TxConfig.TxEncoder())

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		group.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		//ibcexported.StoreKey,
		//ibctransfertypes.StoreKey,
		//capabilitytypes.StoreKey,
		authzkeeper.StoreKey,
		evmtypes.StoreKey,
		feemarkettypes.StoreKey,
		erc20types.StoreKey,
		authoritytypes.StoreKey,
		lightclienttypes.StoreKey,
		crosschaintypes.StoreKey,
		//ibccrosschaintypes.StoreKey,
		observertypes.StoreKey,
		fungibletypes.StoreKey,
		emissionstypes.StoreKey,
		consensusparamtypes.StoreKey,
		crisistypes.StoreKey,
	)
	tKeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey, feemarkettypes.TransientKey)
	memKeys := storetypes.NewMemoryStoreKeys(
	// capabilitytypes.MemStoreKey,
	)

	app := &App{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tKeys,
		memKeys:           memKeys,
	}
	if homePath == "" {
		homePath = DefaultNodeHome
	}

	// get authority address
	authAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tKeys[paramstypes.TStoreKey])
	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey]),
		authAddr,
		runtime.EventService{},
	)
	bApp.SetParamStore(app.ConsensusParamsKeeper.ParamsStore)

	// add capability keeper and ScopeToModule for ibc module
	//app.CapabilityKeeper = capabilitykeeper.NewKeeper(
	//	appCodec,
	//	keys[capabilitytypes.StoreKey],
	//	memKeys[capabilitytypes.MemStoreKey],
	//)

	//scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	//scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	// add keepers
	// use custom Evm account for contracts
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authAddr,
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.AccountKeeper,
		app.BlockedAddrs(),
		authAddr,
		logger,
	)

	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authAddr,
		address.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		address.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		authtypes.FeeCollectorName,
		authAddr,
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		app.LegacyAmino(),
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		app.StakingKeeper,
		authAddr,
	)

	app.CrisisKeeper = *crisiskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[crisistypes.StoreKey]),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		authAddr,
		app.AccountKeeper.AddressCodec(),
	)

	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		app.BaseApp,
		authAddr,
	)

	// IBC keepers

	//app.IBCKeeper = ibckeeper.NewKeeper(
	//	appCodec,
	//	keys[ibcexported.StoreKey],
	//	app.GetSubspace(ibcexported.ModuleName),
	//	app.StakingKeeper,
	//	app.UpgradeKeeper,
	//	scopedIBCKeeper,
	//)
	//
	//ibcRouter := porttypes.NewRouter()
	//
	//app.TransferKeeper = ibctransferkeeper.NewKeeper(
	//	appCodec,
	//	keys[ibctransfertypes.StoreKey],
	//	app.GetSubspace(ibctransfertypes.ModuleName),
	//	app.IBCKeeper.ChannelKeeper,
	//	app.IBCKeeper.ChannelKeeper,
	//	&app.IBCKeeper.PortKeeper,
	//	app.AccountKeeper,
	//	app.BankKeeper,
	//	scopedTransferKeeper,
	//)
	//app.transferModule = transfer.NewAppModule(app.TransferKeeper)
	//
	//// create IBC module from bottom to top of stack
	//var transferStack porttypes.IBCModule
	//transferStack = transfer.NewIBCModule(app.TransferKeeper)
	//
	//// Add transfer stack to IBC Router
	//ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack)

	// ZetaChain keepers

	app.AuthorityKeeper = authoritykeeper.NewKeeper(
		appCodec,
		keys[authoritytypes.StoreKey],
		keys[authoritytypes.MemStoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName),
	)

	app.LightclientKeeper = lightclientkeeper.NewKeeper(
		appCodec,
		keys[lightclienttypes.StoreKey],
		keys[lightclienttypes.MemStoreKey],
		app.AuthorityKeeper,
	)

	app.ObserverKeeper = observerkeeper.NewKeeper(
		appCodec,
		keys[observertypes.StoreKey],
		keys[observertypes.MemStoreKey],
		app.StakingKeeper,
		app.SlashingKeeper,
		app.AuthorityKeeper,
		app.LightclientKeeper,
		app.BankKeeper,
		app.AccountKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks(),
			app.ObserverKeeper.Hooks(),
		),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[authzkeeper.StoreKey]),
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
	)

	app.EmissionsKeeper = *emissionskeeper.NewKeeper(
		appCodec,
		keys[emissionstypes.StoreKey],
		keys[emissionstypes.MemStoreKey],
		authtypes.FeeCollectorName,
		app.BankKeeper,
		app.StakingKeeper,
		app.ObserverKeeper,
		app.AccountKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Create Evm keepers
	tracer := cast.ToString(appOpts.Get(srvflags.EVMTracer))

	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec, authtypes.NewModuleAddress(govtypes.ModuleName),
		keys[feemarkettypes.StoreKey],
		tKeys[feemarkettypes.TransientKey],
	)

	app.EvmKeeper = evmkeeper.NewKeeper(
		appCodec,
		keys[evmtypes.StoreKey],
		tKeys[evmtypes.TransientKey],
		keys,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.FeeMarketKeeper,
		app.ConsensusParamsKeeper,
		&app.Erc20Keeper,
		tracer,
	)

	app.Erc20Keeper = erc20keeper.NewKeeper(
		keys[erc20types.StoreKey],
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.EvmKeeper,
		app.StakingKeeper,
		nil,
	)

	app.FungibleKeeper = *fungiblekeeper.NewKeeper(
		appCodec,
		keys[fungibletypes.StoreKey],
		keys[fungibletypes.MemStoreKey],
		app.AccountKeeper,
		app.EvmKeeper,
		app.BankKeeper,
		app.ObserverKeeper,
		app.AuthorityKeeper,
	)

	app.CrosschainKeeper = *crosschainkeeper.NewKeeper(
		appCodec,
		keys[crosschaintypes.StoreKey],
		keys[crosschaintypes.MemStoreKey],
		app.StakingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		app.ObserverKeeper,
		&app.FungibleKeeper,
		app.AuthorityKeeper,
		app.LightclientKeeper,
	)

	// initialize ibccrosschain keeper and set it to the crosschain keeper
	// there is a circular dependency between the two keepers, crosschain keeper must be initialized first
	//
	//scopedIBCCrosschainKeeper := app.CapabilityKeeper.ScopeToModule(ibccrosschaintypes.ModuleName)
	//app.ScopedIBCCrosschainKeeper = scopedIBCCrosschainKeeper
	//
	//app.IBCCrosschainKeeper = *ibccrosschainkeeper.NewKeeper(
	//	appCodec,
	//	keys[ibccrosschaintypes.StoreKey],
	//	keys[ibccrosschaintypes.MemStoreKey],
	//	&app.CrosschainKeeper,
	//	app.TransferKeeper,
	//)
	//
	//ibcRouter.AddRoute(ibccrosschaintypes.ModuleName, ibccrosschain.NewIBCModule(app.IBCCrosschainKeeper))
	//
	//app.CrosschainKeeper.SetIBCCrosschainKeeper(app.IBCCrosschainKeeper)

	app.GroupKeeper = groupkeeper.NewKeeper(
		keys[group.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
		group.Config{
			MaxExecutionPeriod: 48 * time.Hour,
			MaxMetadataLen:     255,
		},
	)

	// register the proposal types
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper))

	govConfig := govtypes.DefaultConfig()
	govKeeper := govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[govtypes.StoreKey]),
		app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.DistrKeeper,
		app.MsgServiceRouter(), govConfig, authAddr,
	)

	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register governance hooks
		),
	)
	// Set legacy router for backwards compatibility with gov v1beta1
	// app.GovKeeper.SetLegacyRouter(govRouter)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evidencetypes.StoreKey]),
		app.StakingKeeper, app.SlashingKeeper,
		app.AccountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)

	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	app.EvmKeeper = app.EvmKeeper.SetHooks(evmkeeper.NewMultiEvmHooks(
		app.CrosschainKeeper.Hooks(),
		app.FungibleKeeper.EVMHooks(),
	))

	// seal the IBC router
	//app.IBCKeeper.SetRouter(ibcRouter)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// add static precompiles
	app.EvmKeeper.WithStaticPrecompiles(
		NewAvailableStaticPrecompiles(
			*app.StakingKeeper,
			app.DistrKeeper,
			app.BankKeeper,
			app.Erc20Keeper,
			app.EvmKeeper,
			app.GovKeeper,
			app.SlashingKeeper,
			app.AppCodec(),
		),
	)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app, encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		//capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		gov.NewAppModule(
			appCodec,
			&app.GovKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(govtypes.ModuleName),
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
		distr.NewAppModule(
			appCodec,
			app.DistrKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.GetSubspace(distrtypes.ModuleName),
		),
		staking.NewAppModule(
			appCodec,
			app.StakingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.GetSubspace(stakingtypes.ModuleName),
		),
		upgrade.NewAppModule(app.UpgradeKeeper, app.AccountKeeper.AddressCodec()),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		//app.transferModule,
		//ibc.NewAppModule(app.IBCKeeper),
		//transfer.NewAppModule(app.TransferKeeper),
		groupmodule.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, interfaceRegistry),
		feemarket.NewAppModule(app.FeeMarketKeeper),
		vm.NewAppModule(app.EvmKeeper, app.AccountKeeper, app.AccountKeeper.AddressCodec()),
		authoritymodule.NewAppModule(appCodec, app.AuthorityKeeper),
		lightclientmodule.NewAppModule(appCodec, app.LightclientKeeper),
		crosschainmodule.NewAppModule(appCodec, app.CrosschainKeeper),
		//ibccrosschain.NewAppModule(appCodec, app.IBCCrosschainKeeper),
		observermodule.NewAppModule(appCodec, *app.ObserverKeeper),
		fungiblemodule.NewAppModule(appCodec, app.FungibleKeeper),
		emissionsmodule.NewAppModule(appCodec, app.EmissionsKeeper, app.GetSubspace(emissionstypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
	)

	// BasicModuleManager defines the module BasicManager which is in charge of setting up basic,
	// non-dependant module elements, such as codec registration and genesis verification.
	// By default, it is composed of all the modules from the module manager.
	// Additionally, app module basics can be overwritten by passing them as an argument.
	app.mb = module.NewBasicManagerFromManager(
		app.mm,
		map[string]module.AppModuleBasic{
			genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			govtypes.ModuleName: gov.NewAppModuleBasic(
				[]govclient.ProposalHandler{
					paramsclient.ProposalHandler,
				},
			),
		},
	)
	app.mb.RegisterLegacyAminoCodec(cdc)
	app.mb.RegisterInterfaces(interfaceRegistry)

	app.mm.SetOrderPreBlockers(
		upgradetypes.ModuleName,
		authtypes.ModuleName,
	)

	app.mm.SetOrderBeginBlockers(orderBeginBlockers()...)

	app.mm.SetOrderEndBlockers(orderEndBlockers()...)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	// NOTE: Cross-chain module must be initialized after observer module, as pending nonces in crosschain needs the tss pubkey from observer module
	app.mm.SetOrderInitGenesis(OrderInitGenesis()...)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	if err := app.mm.RegisterServices(app.configurator); err != nil {
		panic(err)
	}

	app.sm = module.NewSimulationManager(simulationModules(app, appCodec, skipGenesisInvariants)...)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)

	options := ante.HandlerOptions{
		AccountKeeper:   app.AccountKeeper,
		BankKeeper:      app.BankKeeper,
		EvmKeeper:       app.EvmKeeper,
		FeeMarketKeeper: app.FeeMarketKeeper,
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		MaxTxGasWanted:  TransactionGasLimit,
		DisabledAuthzMsgs: []string{
			sdk.MsgTypeURL(
				&evmtypes.MsgEthereumTx{},
			), // disable the Msg types that cannot be included on an authz.MsgExec msgs field
			sdk.MsgTypeURL(&vestingtypes.MsgCreateVestingAccount{}),
			sdk.MsgTypeURL(&vestingtypes.MsgCreatePermanentLockedAccount{}),
			sdk.MsgTypeURL(&vestingtypes.MsgCreatePeriodicVestingAccount{}),
		},
		ObserverKeeper: app.ObserverKeeper,
	}

	anteHandler, err := ante.NewAnteHandler(options)
	if err != nil {
		panic(err)
	}

	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)
	SetupHandlers(app)
	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	//app.ScopedIBCKeeper = scopedIBCKeeper
	//app.ScopedTransferKeeper = scopedTransferKeeper

	return app
}

// Name returns the name of the App
func (app *App) Name() string { return app.BaseApp.Name() }

// PreBlocker updates every pre begin block
func (app *App) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	return app.mm.PreBlock(ctx)
}

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.mm.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	// InitChainErrorMessage is the error message displayed when trying to sync testnet or mainnet from block 1 using the latest binary.
	const InitChainErrorMessage = `
Unable to sync testnet or mainnet from block 1 using the latest version.
Please use a snapshot to sync your node.
Refer to the documentation for more information:
https://www.zetachain.com/docs/nodes/start-here/syncing/`

	// The defer is used to catch panics during InitChain
	// and display a more meaningful message for people trying to sync a node from block 1 using the latest binary.
	// We exit the process after displaying the message as we do not need to start a node with empty state.

	defer func() {
		if r := recover(); r != nil {
			ctx.Logger().Error("panic occurred during InitGenesis", "error", r)
			ctx.Logger().Debug("stack trace", "stack", string(debug.Stack()))
			ctx.Logger().
				Info(InitChainErrorMessage)
			os.Exit(1)
		}
	}()
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap()); err != nil {
		panic(err)
	}

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// LegacyAmino returns app's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns Zeta app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Gaia's InterfaceRegistry
func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// AutoCliOpts returns the autocli options for the app.
func (app *App) AutoCliOpts() autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range app.mm.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservices.ExtractAutoCLIOptions(app.mm.Modules),
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	}
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register legacy tx routes.
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	app.mb.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register app's OpenAPI routes.
	if apiConfig.Swagger {
		openapi.RegisterOpenAPIService(apiSvr.Router)
	}
}

func (app *App) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(_ client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(
	appCodec codec.BinaryCodec,
	legacyAmino *codec.LegacyAmino,
	key, tkey storetypes.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	//paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	//paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(evmtypes.ModuleName)
	paramsKeeper.Subspace(feemarkettypes.ModuleName)
	paramsKeeper.Subspace(group.ModuleName)
	paramsKeeper.Subspace(observertypes.ModuleName)
	paramsKeeper.Subspace(emissionstypes.ModuleName)
	return paramsKeeper
}

// SimulationManager implements the SimulationApp interface
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

func (app *App) BasicManager() module.BasicManager {
	return app.mb
}

func (app *App) ModuleManager() *module.Manager {
	return app.mm
}

func (app *App) BlockedAddrs() map[string]bool {
	blockList := make(map[string]bool)

	for k, v := range blockedReceivingModAcc {
		addr := authtypes.NewModuleAddress(k)
		blockList[addr.String()] = v
	}

	return blockList
}
