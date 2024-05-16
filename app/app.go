package app

import (
	"io"
	"net/http"
	"os"
	"time"

	cosmoserrors "cosmossdk.io/errors"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	appparams "cosmossdk.io/simapp/params"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
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
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	evmante "github.com/evmos/ethermint/app/ante"
	ethermint "github.com/evmos/ethermint/types"
	"github.com/evmos/ethermint/x/evm"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/evm/vm/geth"
	"github.com/evmos/ethermint/x/feemarket"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/zeta-chain/zetacore/app/ante"
	"github.com/zeta-chain/zetacore/docs/openapi"
	srvflags "github.com/zeta-chain/zetacore/server/flags"

	authoritymodule "github.com/zeta-chain/zetacore/x/authority"
	authoritykeeper "github.com/zeta-chain/zetacore/x/authority/keeper"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"

	lightclientmodule "github.com/zeta-chain/zetacore/x/lightclient"
	lightclientkeeper "github.com/zeta-chain/zetacore/x/lightclient/keeper"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"

	crosschainmodule "github.com/zeta-chain/zetacore/x/crosschain"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"

	emissionsmodule "github.com/zeta-chain/zetacore/x/emissions"
	emissionskeeper "github.com/zeta-chain/zetacore/x/emissions/keeper"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"

	fungiblemodule "github.com/zeta-chain/zetacore/x/fungible"
	fungiblekeeper "github.com/zeta-chain/zetacore/x/fungible/keeper"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"

	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	observermodule "github.com/zeta-chain/zetacore/x/observer"
	observerkeeper "github.com/zeta-chain/zetacore/x/observer/keeper"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"

	"github.com/zeta-chain/zetacore/x/ibccrosschain"
	ibccrosschainkeeper "github.com/zeta-chain/zetacore/x/ibccrosschain/keeper"
	ibccrosschaintypes "github.com/zeta-chain/zetacore/x/ibccrosschain/types"
)

const Name = "zetacore"

func init() {
	// manually update the power reduction by replacing micro (u) -> atto (a) evmos
	sdk.DefaultPowerReduction = ethermint.PowerReduction
	// modify fee market parameter defaults through global
	//feemarkettypes.DefaultMinGasPrice = v5.MainnetMinGasPrices
	//feemarkettypes.DefaultMinGasMultiplier = v5.MainnetMinGasMultiplier
}

var (
	AccountAddressPrefix = "zeta"
	NodeDir              = ".zetacored"

	// AddrLen is the allowed length (in bytes) for an address.
	//
	// NOTE: In the SDK, the default value is 255.
	AddrLen = 20
)

var (
	// DefaultNodeHome default home directories for wasmd
	DefaultNodeHome = os.ExpandEnv("$HOME/") + NodeDir

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = AccountAddressPrefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = AccountAddressPrefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler
	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
		upgradeclient.LegacyProposalHandler,
		upgradeclient.LegacyCancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
	)
	return govProposalHandlers
}

var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(getGovProposalHandlers()),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibctm.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		consensus.AppModuleBasic{},
		evm.AppModuleBasic{},
		feemarket.AppModuleBasic{},
		authoritymodule.AppModuleBasic{},
		lightclientmodule.AppModuleBasic{},
		crosschainmodule.AppModuleBasic{},
		ibccrosschain.AppModuleBasic{},
		observermodule.AppModuleBasic{},
		fungiblemodule.AppModuleBasic{},
		emissionsmodule.AppModuleBasic{},
		groupmodule.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:                      nil,
		distrtypes.ModuleName:                           nil,
		stakingtypes.BondedPoolName:                     {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName:                  {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:                             {authtypes.Burner},
		ibctransfertypes.ModuleName:                     {authtypes.Minter, authtypes.Burner},
		crosschaintypes.ModuleName:                      {authtypes.Minter, authtypes.Burner},
		ibccrosschaintypes.ModuleName:                   nil,
		evmtypes.ModuleName:                             {authtypes.Minter, authtypes.Burner},
		fungibletypes.ModuleName:                        {authtypes.Minter, authtypes.Burner},
		emissionstypes.ModuleName:                       nil,
		emissionstypes.UndistributedObserverRewardsPool: nil,
		emissionstypes.UndistributedTssRewardsPool:      nil,
	}

	// module accounts that are NOT allowed to receive tokens
	blockedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName:          true,
		authtypes.FeeCollectorName:     true,
		stakingtypes.BondedPoolName:    true,
		stakingtypes.NotBondedPoolName: true,
		govtypes.ModuleName:            true,
	}
)

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
	configurator module.Configurator

	// sdk keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             govkeeper.Keeper
	CrisisKeeper          crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper
	TransferKeeper        ibctransferkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	GroupKeeper           groupkeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper

	// scoped keepers
	ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper      capabilitykeeper.ScopedKeeper
	ScopedIBCCrosschainKeeper capabilitykeeper.ScopedKeeper

	// evm keepers
	EvmKeeper       *evmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper

	// zetachain keepers
	AuthorityKeeper     authoritykeeper.Keeper
	LightclientKeeper   lightclientkeeper.Keeper
	CrosschainKeeper    crosschainkeeper.Keeper
	IBCCrosschainKeeper ibccrosschainkeeper.Keeper
	ObserverKeeper      *observerkeeper.Keeper
	FungibleKeeper      fungiblekeeper.Keeper
	EmissionsKeeper     emissionskeeper.Keeper

	transferModule transfer.AppModule
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
	encodingConfig appparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {

	appCodec := encodingConfig.Codec
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(Name, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(encodingConfig.TxConfig.TxEncoder())

	keys := sdk.NewKVStoreKeys(
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
		ibcexported.StoreKey,
		ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey,
		authzkeeper.StoreKey,
		evmtypes.StoreKey,
		feemarkettypes.StoreKey,
		authoritytypes.StoreKey,
		lightclienttypes.StoreKey,
		crosschaintypes.StoreKey,
		ibccrosschaintypes.StoreKey,
		observertypes.StoreKey,
		fungibletypes.StoreKey,
		emissionstypes.StoreKey,
		consensusparamtypes.StoreKey,
		crisistypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey, feemarkettypes.TransientKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &App{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}
	if homePath == "" {
		homePath = DefaultNodeHome
	}

	// get authority address
	authAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, keys[consensusparamtypes.StoreKey], authAddr)
	bApp.SetParamStore(&app.ConsensusParamsKeeper)

	customProposalHandler := NewCustomProposalHandler(bApp.Mempool(), bApp)
	app.SetPrepareProposal(customProposalHandler.PrepareProposalHandler())

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)

	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	// add keepers
	// use custom Ethermint account for contracts
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey],
		ethermint.ProtoAccount,
		maccPerms,
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authAddr,
	)

	logger.Info("bank keeper blocklist addresses", "addresses", app.BlockedAddrs())

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.BlockedAddrs(),
		authAddr,
	)

	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		authAddr,
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		authtypes.FeeCollectorName,
		authAddr,
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		app.LegacyAmino(),
		keys[slashingtypes.StoreKey],
		app.StakingKeeper,
		authAddr,
	)

	app.CrisisKeeper = *crisiskeeper.NewKeeper(
		appCodec,
		keys[crisistypes.StoreKey],
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		authAddr,
	)

	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
		authAddr,
	)

	// IBC keepers

	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		app.GetSubspace(ibcexported.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	ibcRouter := porttypes.NewRouter()

	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
	)
	app.transferModule = transfer.NewAppModule(app.TransferKeeper)

	// create IBC module from bottom to top of stack
	var transferStack porttypes.IBCModule
	transferStack = transfer.NewIBCModule(app.TransferKeeper)

	// Add transfer stack to IBC Router
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack)

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
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks(), app.ObserverKeeper.Hooks()),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
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

	// Create Ethermint keepers
	tracer := cast.ToString(appOpts.Get(srvflags.EVMTracer))
	feeSs := app.GetSubspace(feemarkettypes.ModuleName)
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		keys[feemarkettypes.StoreKey],
		tkeys[feemarkettypes.TransientKey],
		feeSs,
		app.ConsensusParamsKeeper,
	)
	evmSs := app.GetSubspace(evmtypes.ModuleName)
	app.EvmKeeper = evmkeeper.NewKeeper(
		appCodec,
		keys[evmtypes.StoreKey],
		tkeys[evmtypes.TransientKey],
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		&app.FeeMarketKeeper,
		nil,
		geth.NewEVM,
		tracer,
		evmSs,
		app.ConsensusParamsKeeper,
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

	scopedIBCCrosschainKeeper := app.CapabilityKeeper.ScopeToModule(ibccrosschaintypes.ModuleName)
	app.ScopedIBCCrosschainKeeper = scopedIBCCrosschainKeeper

	app.IBCCrosschainKeeper = *ibccrosschainkeeper.NewKeeper(
		appCodec,
		keys[ibccrosschaintypes.StoreKey],
		keys[ibccrosschaintypes.MemStoreKey],
		&app.CrosschainKeeper,
		app.TransferKeeper,
	)

	ibcRouter.AddRoute(ibccrosschaintypes.ModuleName, ibccrosschain.NewIBCModule(app.IBCCrosschainKeeper))

	app.CrosschainKeeper.SetIBCCrosschainKeeper(app.IBCCrosschainKeeper)

	app.GroupKeeper = groupkeeper.NewKeeper(
		keys[group.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
		group.Config{
			MaxExecutionPeriod: 2 * time.Hour, // Two hours.
			MaxMetadataLen:     255,
		},
	)

	// register the proposal types
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper))
	govConfig := govtypes.DefaultConfig()
	govKeeper := govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.MsgServiceRouter(),
		govConfig,
		authAddr,
	)
	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register governance hooks
		),
	)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	app.EvmKeeper = app.EvmKeeper.SetHooks(evmkeeper.NewMultiEvmHooks(
		app.CrosschainKeeper.Hooks(),
		app.FungibleKeeper.EVMHooks(),
	))

	// seal the IBC router
	app.IBCKeeper.SetRouter(ibcRouter)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		app.transferModule,
		ibc.NewAppModule(app.IBCKeeper),
		transfer.NewAppModule(app.TransferKeeper),
		groupmodule.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, interfaceRegistry),
		feemarket.NewAppModule(app.FeeMarketKeeper, feeSs),
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper, evmSs),
		authoritymodule.NewAppModule(appCodec, app.AuthorityKeeper),
		lightclientmodule.NewAppModule(appCodec, app.LightclientKeeper),
		crosschainmodule.NewAppModule(appCodec, app.CrosschainKeeper),
		ibccrosschain.NewAppModule(appCodec, app.IBCCrosschainKeeper),
		observermodule.NewAppModule(appCodec, *app.ObserverKeeper),
		fungiblemodule.NewAppModule(appCodec, app.FungibleKeeper),
		emissionsmodule.NewAppModule(appCodec, app.EmissionsKeeper, app.GetSubspace(emissionstypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		capabilitytypes.ModuleName,
		upgradetypes.ModuleName,
		evmtypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		group.ModuleName,
		vestingtypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		feemarkettypes.ModuleName,
		crosschaintypes.ModuleName,
		ibccrosschaintypes.ModuleName,
		observertypes.ModuleName,
		fungibletypes.ModuleName,
		emissionstypes.ModuleName,
		authz.ModuleName,
		authoritytypes.ModuleName,
		lightclienttypes.ModuleName,
		consensusparamtypes.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		capabilitytypes.ModuleName,
		banktypes.ModuleName,
		authtypes.ModuleName,
		upgradetypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		vestingtypes.ModuleName,
		govtypes.ModuleName,
		paramstypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		genutiltypes.ModuleName,
		group.ModuleName,
		crisistypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		crosschaintypes.ModuleName,
		ibccrosschaintypes.ModuleName,
		observertypes.ModuleName,
		fungibletypes.ModuleName,
		emissionstypes.ModuleName,
		authz.ModuleName,
		authoritytypes.ModuleName,
		lightclienttypes.ModuleName,
		consensusparamtypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	// NOTE: Cross-chain module must be initialized after observer module, as pending nonces in crosschain needs the tss pubkey from observer module
	app.mm.SetOrderInitGenesis(InitGenesisModuleList()...)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	maxGasWanted := cast.ToUint64(appOpts.Get(srvflags.EVMMaxTxGasWanted))
	options := ante.HandlerOptions{
		AccountKeeper:   app.AccountKeeper,
		BankKeeper:      app.BankKeeper,
		EvmKeeper:       app.EvmKeeper,
		FeeMarketKeeper: app.FeeMarketKeeper,
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:  evmante.DefaultSigVerificationGasConsumer,
		MaxTxGasWanted:  maxGasWanted,
		DisabledAuthzMsgs: []string{
			sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{}), // disable the Msg types that cannot be included on an authz.MsgExec msgs field
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

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	return app
}

// Name returns the name of the App
func (app *App) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
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

// LegacyAmino returns SimApp's amino codec.
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
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register app's OpenAPI routes.
	if apiConfig.Swagger {
		openapi.RegisterOpenAPIService(apiSvr.Router)
	}
}

func (app *App) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
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
	tmservice.RegisterTendermintService(
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
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(evmtypes.ModuleName)
	paramsKeeper.Subspace(feemarkettypes.ModuleName)
	paramsKeeper.Subspace(group.ModuleName)
	paramsKeeper.Subspace(observertypes.ModuleName)
	paramsKeeper.Subspace(emissionstypes.ModuleName)
	return paramsKeeper
}

// VerifyAddressFormat verifies the address is compatible with ethereum
func VerifyAddressFormat(bz []byte) error {
	if len(bz) == 0 {
		return cosmoserrors.Wrap(sdkerrors.ErrUnknownAddress, "invalid address; cannot be empty")
	}
	if len(bz) != AddrLen {
		return cosmoserrors.Wrapf(
			sdkerrors.ErrUnknownAddress,
			"invalid address length; got: %d, expect: %d", len(bz), AddrLen,
		)
	}

	return nil
}

// SimulationManager implements the SimulationApp interface
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

func (app *App) BlockedAddrs() map[string]bool {
	blockList := make(map[string]bool)
	for k, v := range blockedReceivingModAcc {
		addr := authtypes.NewModuleAddress(k)
		blockList[addr.String()] = v
	}
	return blockList
}
