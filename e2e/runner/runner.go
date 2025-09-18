package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/protocol-contracts/pkg/coreregistry.sol"
	erc20custodyv2 "github.com/zeta-chain/protocol-contracts/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/wzeta.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/zetaconnector.eth.sol"
	zetaconnnectornative "github.com/zeta-chain/protocol-contracts/pkg/zetaconnectornative.sol"
	connectorzevm "github.com/zeta-chain/protocol-contracts/pkg/zetaconnectorzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/zetaeth.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/zrc20.sol"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/contracts/erc20"
	"github.com/zeta-chain/node/e2e/contracts/testdappv2"
	"github.com/zeta-chain/node/e2e/contracts/zevmswap"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/pkg/contracts/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/zeta-chain/node/pkg/contracts/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/node/x/lightclient/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	btcclient "github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
)

type E2ERunnerOption func(*E2ERunner)

// Important ENV
const (
	EnvKeyLocalnetMode = "LOCALNET_MODE"

	LocalnetModeUpgrade      = "upgrade"
	LocalNetModeTSSMigration = "tss-migration"

	// NodeSyncTolerance is the time tolerance for the ZetaChain nodes behind a RPC to be synced
	NodeSyncTolerance = constant.ZetaBlockTime * 5
)

func WithZetaTxServer(txServer *txserver.ZetaTxServer) E2ERunnerOption {
	return func(r *E2ERunner) {
		r.ZetaTxServer = txServer
	}
}

func WithTestFilter(testFilter *regexp.Regexp) E2ERunnerOption {
	return func(r *E2ERunner) {
		r.TestFilter = testFilter
	}
}

// E2ERunner stores all the clients and addresses needed for E2E test
// Exposes a method to run E2E test
// It also provides some helper functions
type E2ERunner struct {
	// accounts
	Account               config.Account
	TSSAddress            ethcommon.Address
	BTCTSSAddress         btcutil.Address
	SuiTSSAddress         string
	SolanaDeployerAddress solana.PublicKey
	FeeCollectorAddress   types.AccAddress

	// all clients.
	// a reference to this type is required to enable creating a new E2ERunner.
	Clients Clients

	// rpc clients
	ZEVMClient   *ethclient.Client
	EVMClient    *ethclient.Client
	BtcRPCClient *btcclient.Client
	SolanaClient *rpc.Client

	// zetacored grpc clients
	AuthorityClient    authoritytypes.QueryClient
	CctxClient         crosschaintypes.QueryClient
	FungibleClient     fungibletypes.QueryClient
	AuthClient         authtypes.QueryClient
	BankClient         banktypes.QueryClient
	StakingClient      stakingtypes.QueryClient
	ObserverClient     observertypes.QueryClient
	LightclientClient  lightclienttypes.QueryClient
	DistributionClient distributiontypes.QueryClient
	EmissionsClient    emissionstypes.QueryClient

	// optional zeta (cosmos) client
	// typically only in test runners that need it
	// (like admin tests)
	ZetaTxServer *txserver.ZetaTxServer

	// evm auth
	EVMAuth  *bind.TransactOpts
	ZEVMAuth *bind.TransactOpts

	// programs on Solana
	GatewayProgram      solana.PublicKey
	SPLAddr             solana.PublicKey
	ConnectedProgram    solana.PublicKey
	ConnectedSPLProgram solana.PublicKey

	// TON related
	TONGateway ton.AccountID

	// contract Sui
	SuiGateway *sui.Gateway

	// SuiGatewayUpgradeCap is the upgrade cap used for upgrading the Sui gateway package
	SuiGatewayUpgradeCap string

	// SuiTokenCoinType is the coin type identifying the fungible token for SUI
	SuiTokenCoinType string

	// SuiTokenTreasuryCap is the treasury cap for the SUI token that allows minting, only using in local tests
	SuiTokenTreasuryCap string

	// SuiExample contains the example package information for Sui authenticated call
	SuiExample config.SuiExample

	// contracts evm
	ZetaEthAddr       ethcommon.Address
	ZetaEth           *zetaeth.ZetaEth
	ERC20CustodyAddr  ethcommon.Address
	ERC20Custody      *erc20custodyv2.ERC20Custody
	ERC20Addr         ethcommon.Address
	ERC20             *erc20.ERC20
	EvmTestDAppAddr   ethcommon.Address
	GatewayEVMAddr    ethcommon.Address
	GatewayEVM        *gatewayevm.GatewayEVM
	TestDAppV2EVMAddr ethcommon.Address
	TestDAppV2EVM     *testdappv2.TestDAppV2
	// ConnectorNative is the V2 connector for EVM chains
	ConnectorNativeAddr ethcommon.Address
	ConnectorNative     *zetaconnnectornative.ZetaConnectorNative
	// ConnectorEthAddr is the V1 connector for EVM chains
	ConnectorEthAddr ethcommon.Address
	ConnectorEth     *zetaconnectoreth.ZetaConnectorEth

	// contracts zevm
	// zrc20 contracts
	ERC20ZRC20Addr    ethcommon.Address
	ERC20ZRC20        *zrc20.ZRC20
	SPLZRC20Addr      ethcommon.Address
	SPLZRC20          *zrc20.ZRC20
	ETHZRC20Addr      ethcommon.Address
	ETHZRC20          *zrc20.ZRC20
	BTCZRC20Addr      ethcommon.Address
	BTCZRC20          *zrc20.ZRC20
	SOLZRC20Addr      ethcommon.Address
	SOLZRC20          *zrc20.ZRC20
	TONZRC20Addr      ethcommon.Address
	TONZRC20          *zrc20.ZRC20
	SUIZRC20Addr      ethcommon.Address
	SUIZRC20          *zrc20.ZRC20
	SuiTokenZRC20Addr ethcommon.Address
	SuiTokenZRC20     *zrc20.ZRC20

	// other contracts
	UniswapV2FactoryAddr ethcommon.Address
	UniswapV2Factory     *uniswapv2factory.UniswapV2Factory
	UniswapV2RouterAddr  ethcommon.Address
	UniswapV2Router      *uniswapv2router.UniswapV2Router02
	ConnectorZEVMAddr    ethcommon.Address
	ConnectorZEVM        *connectorzevm.ZetaConnectorZEVM
	WZetaAddr            ethcommon.Address
	WZeta                *wzeta.WETH9
	ZEVMSwapAppAddr      ethcommon.Address
	ZEVMSwapApp          *zevmswap.ZEVMSwapApp
	SystemContractAddr   ethcommon.Address
	SystemContract       *systemcontract.SystemContract
	ZevmTestDAppAddr     ethcommon.Address
	GatewayZEVMAddr      ethcommon.Address
	GatewayZEVM          *gatewayzevm.GatewayZEVM
	TestDAppV2ZEVMAddr   ethcommon.Address
	TestDAppV2ZEVM       *testdappv2.TestDAppV2
	CoreRegistryAddr     ethcommon.Address
	CoreRegistry         *coreregistry.CoreRegistry

	// config
	CctxTimeout    time.Duration
	ReceiptTimeout time.Duration

	// other
	Name             string
	Ctx              context.Context
	CtxCancel        context.CancelCauseFunc
	Logger           *Logger
	BitcoinParams    *chaincfg.Params
	TestFilter       *regexp.Regexp
	mutex            sync.Mutex
	zetacoredVersion string
}

func NewE2ERunner(
	ctx context.Context,
	name string,
	ctxCancel context.CancelCauseFunc,
	account config.Account,
	clients Clients,
	logger *Logger,
	opts ...E2ERunnerOption,
) *E2ERunner {
	r := &E2ERunner{
		Name:      name,
		CtxCancel: ctxCancel,

		Account: account,

		FeeCollectorAddress: authtypes.NewModuleAddress(authtypes.FeeCollectorName),

		Clients: clients,

		ZEVMClient:         clients.Zevm,
		EVMClient:          clients.Evm,
		AuthorityClient:    clients.Zetacore.Authority,
		CctxClient:         clients.Zetacore.Crosschain,
		FungibleClient:     clients.Zetacore.Fungible,
		AuthClient:         clients.Zetacore.Auth,
		BankClient:         clients.Zetacore.Bank,
		StakingClient:      clients.Zetacore.Staking,
		ObserverClient:     clients.Zetacore.Observer,
		LightclientClient:  clients.Zetacore.Lightclient,
		DistributionClient: clients.Zetacore.Distribution,
		EmissionsClient:    clients.Zetacore.Emissions,

		EVMAuth:      clients.EvmAuth,
		ZEVMAuth:     clients.ZevmAuth,
		BtcRPCClient: clients.BtcRPC,
		SolanaClient: clients.Solana,

		Logger: logger,
	}

	r.Ctx = utils.WithTesting(ctx, r)

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// CopyAddressesFrom copies addresses from another E2ETestRunner that initialized the contracts
func (r *E2ERunner) CopyAddressesFrom(other *E2ERunner) (err error) {
	// copy TSS address
	r.TSSAddress = other.TSSAddress
	r.BTCTSSAddress = other.BTCTSSAddress
	r.SuiTSSAddress = other.SuiTSSAddress

	// copy addresses
	r.ZetaEthAddr = other.ZetaEthAddr
	r.ConnectorEthAddr = other.ConnectorEthAddr
	r.ERC20CustodyAddr = other.ERC20CustodyAddr
	r.ERC20Addr = other.ERC20Addr
	r.ERC20ZRC20Addr = other.ERC20ZRC20Addr
	r.ETHZRC20Addr = other.ETHZRC20Addr
	r.BTCZRC20Addr = other.BTCZRC20Addr
	r.SOLZRC20Addr = other.SOLZRC20Addr
	r.TONZRC20Addr = other.TONZRC20Addr
	r.SUIZRC20Addr = other.SUIZRC20Addr
	r.SuiTokenZRC20Addr = other.SuiTokenZRC20Addr
	r.UniswapV2FactoryAddr = other.UniswapV2FactoryAddr
	r.UniswapV2RouterAddr = other.UniswapV2RouterAddr
	r.ConnectorZEVMAddr = other.ConnectorZEVMAddr
	r.WZetaAddr = other.WZetaAddr
	r.EvmTestDAppAddr = other.EvmTestDAppAddr
	r.ZEVMSwapAppAddr = other.ZEVMSwapAppAddr
	r.SystemContractAddr = other.SystemContractAddr
	r.ZevmTestDAppAddr = other.ZevmTestDAppAddr

	r.GatewayProgram = other.GatewayProgram

	r.TONGateway = other.TONGateway

	r.SuiGateway = other.SuiGateway
	r.SuiGatewayUpgradeCap = other.SuiGatewayUpgradeCap
	r.SuiTokenCoinType = other.SuiTokenCoinType
	r.SuiTokenTreasuryCap = other.SuiTokenTreasuryCap
	r.SuiExample = other.SuiExample

	// create instances of contracts
	r.ZetaEth, err = zetaeth.NewZetaEth(r.ZetaEthAddr, r.EVMClient)
	if err != nil {
		return err
	}
	r.ConnectorEth, err = zetaconnectoreth.NewZetaConnectorEth(r.ConnectorEthAddr, r.EVMClient)
	if err != nil {
		return err
	}
	r.ERC20Custody, err = erc20custodyv2.NewERC20Custody(r.ERC20CustodyAddr, r.EVMClient)
	if err != nil {
		return err
	}
	r.ERC20, err = erc20.NewERC20(r.ERC20Addr, r.EVMClient)
	if err != nil {
		return err
	}
	r.ERC20ZRC20, err = zrc20.NewZRC20(r.ERC20ZRC20Addr, r.ZEVMClient)
	if err != nil {
		return err
	}
	r.ETHZRC20, err = zrc20.NewZRC20(r.ETHZRC20Addr, r.ZEVMClient)
	if err != nil {
		return err
	}
	r.BTCZRC20, err = zrc20.NewZRC20(r.BTCZRC20Addr, r.ZEVMClient)
	if err != nil {
		return err
	}
	r.SOLZRC20, err = zrc20.NewZRC20(r.SOLZRC20Addr, r.ZEVMClient)
	if err != nil {
		return err
	}
	r.TONZRC20, err = zrc20.NewZRC20(r.TONZRC20Addr, r.ZEVMClient)
	if err != nil {
		return err
	}
	r.SUIZRC20, err = zrc20.NewZRC20(r.SUIZRC20Addr, r.ZEVMClient)
	if err != nil {
		return err
	}
	r.SuiTokenZRC20, err = zrc20.NewZRC20(r.SuiTokenZRC20Addr, r.ZEVMClient)
	if err != nil {
		return err
	}

	r.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(r.UniswapV2FactoryAddr, r.ZEVMClient)
	if err != nil {
		return err
	}
	r.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(r.UniswapV2RouterAddr, r.ZEVMClient)
	if err != nil {
		return err
	}
	r.ConnectorZEVM, err = connectorzevm.NewZetaConnectorZEVM(r.ConnectorZEVMAddr, r.ZEVMClient)
	if err != nil {
		return err
	}
	r.WZeta, err = wzeta.NewWETH9(r.WZetaAddr, r.ZEVMClient)
	if err != nil {
		return err
	}

	r.ZEVMSwapApp, err = zevmswap.NewZEVMSwapApp(r.ZEVMSwapAppAddr, r.ZEVMClient)
	if err != nil {
		return err
	}
	r.SystemContract, err = systemcontract.NewSystemContract(r.SystemContractAddr, r.ZEVMClient)
	if err != nil {
		return err
	}

	// v2 contracts
	r.GatewayEVMAddr = other.GatewayEVMAddr
	r.GatewayEVM, err = gatewayevm.NewGatewayEVM(r.GatewayEVMAddr, r.EVMClient)
	if err != nil {
		return err
	}

	r.TestDAppV2EVMAddr = other.TestDAppV2EVMAddr
	r.TestDAppV2EVM, err = testdappv2.NewTestDAppV2(r.TestDAppV2EVMAddr, r.EVMClient)
	if err != nil {
		return err
	}

	r.GatewayZEVMAddr = other.GatewayZEVMAddr
	r.GatewayZEVM, err = gatewayzevm.NewGatewayZEVM(r.GatewayZEVMAddr, r.ZEVMClient)
	if err != nil {
		return err
	}
	r.TestDAppV2ZEVMAddr = other.TestDAppV2ZEVMAddr
	r.TestDAppV2ZEVM, err = testdappv2.NewTestDAppV2(r.TestDAppV2ZEVMAddr, r.ZEVMClient)
	if err != nil {
		return err
	}
	r.ConnectorNativeAddr = other.ConnectorNativeAddr
	r.ConnectorNative, err = zetaconnnectornative.NewZetaConnectorNative(r.ConnectorNativeAddr, r.EVMClient)
	if err != nil {
		return err
	}
	return nil
}

// Lock locks the mutex
func (r *E2ERunner) Lock() {
	r.mutex.Lock()
}

// Unlock unlocks the mutex
func (r *E2ERunner) Unlock() {
	r.mutex.Unlock()
}

// PrintContractAddresses prints the addresses of the contracts
// the printed contracts are grouped in a zevm and evm section
// there is a padding used to print the addresses at the same position
func (r *E2ERunner) PrintContractAddresses() {
	r.Logger.Print("Zetacored version: %s ", r.GetZetacoredVersion())

	r.Logger.Print(" --- 📜Solana addresses ---")
	r.Logger.Print("GatewayProgram:      %s", r.GatewayProgram.String())
	r.Logger.Print("SPL:                 %s", r.SPLAddr.String())
	r.Logger.Print("ConnectedProgram:    %s", r.ConnectedProgram.String())
	r.Logger.Print("ConnectedSPLProgram: %s", r.ConnectedSPLProgram.String())

	r.Logger.Print(" --- 📜TON addresses ---")
	if !r.TONGateway.IsZero() {
		r.Logger.Print("Gateway:        %s", r.TONGateway.ToRaw())
	} else {
		r.Logger.Print("Gateway:        not set! 💤")
	}

	r.Logger.Print(" --- 📜Sui addresses ---")
	if r.SuiGateway != nil {
		r.Logger.Print("GatewayPackageID:  %s", r.SuiGateway.PackageID())
		r.Logger.Print("GatewayObjectID:   %s", r.SuiGateway.ObjectID())
		r.Logger.Print("GatewayUpgradeCap: %s", r.SuiGatewayUpgradeCap)
		r.Logger.Print("ExamplePackageID:  %s", r.SuiExample.PackageID)
	} else {
		r.Logger.Print("💤 Sui tests disabled")
	}

	// zevm contracts
	r.Logger.Print(" --- 📜zEVM contracts ---")
	r.Logger.Print("SystemContract: %s", r.SystemContractAddr.Hex())
	r.Logger.Print("ETHZRC20:       %s", r.ETHZRC20Addr.Hex())
	r.Logger.Print("ERC20ZRC20:     %s", r.ERC20ZRC20Addr.Hex())
	r.Logger.Print("BTCZRC20:       %s", r.BTCZRC20Addr.Hex())
	r.Logger.Print("SOLZRC20:       %s", r.SOLZRC20Addr.Hex())
	r.Logger.Print("SPLZRC20:       %s", r.SPLZRC20Addr.Hex())
	r.Logger.Print("TONZRC20:       %s", r.TONZRC20Addr.Hex())
	r.Logger.Print("SUIZRC20:       %s", r.SUIZRC20Addr.Hex())
	r.Logger.Print("SuiTokenZRC20:  %s", r.SuiTokenZRC20Addr.Hex())
	r.Logger.Print("UniswapFactory: %s", r.UniswapV2FactoryAddr.Hex())
	r.Logger.Print("UniswapRouter:  %s", r.UniswapV2RouterAddr.Hex())
	r.Logger.Print("ConnectorZEVM:  %s", r.ConnectorZEVMAddr.Hex())
	r.Logger.Print("WZeta:          %s", r.WZetaAddr.Hex())
	r.Logger.Print("GatewayZEVM:    %s", r.GatewayZEVMAddr.Hex())
	r.Logger.Print("TestDAppV2ZEVM: %s", r.TestDAppV2ZEVMAddr.Hex())
	r.Logger.Print("CoreRegistry:   %s", r.CoreRegistryAddr.Hex())

	// evm contracts
	r.Logger.Print(" --- 📜EVM contracts ---")
	r.Logger.Print("ZetaEth:           	%s", r.ZetaEthAddr.Hex())
	r.Logger.Print("ConnectorEthLegacy:	%s", r.ConnectorEthAddr.Hex())
	r.Logger.Print("ConnectorEth:  		%s", r.ConnectorNativeAddr.Hex())
	r.Logger.Print("ERC20Custody  		%s", r.ERC20CustodyAddr.Hex())
	r.Logger.Print("ERC20:            	%s", r.ERC20Addr.Hex())
	r.Logger.Print("GatewayEVM:       	%s", r.GatewayEVMAddr.Hex())
	r.Logger.Print("TestDAppV2EVM:     	%s", r.TestDAppV2EVMAddr.Hex())

	r.Logger.Print(" --- 📜Legacy contracts ---")

	r.Logger.Print("ZEVMSwapApp:    %s", r.ZEVMSwapAppAddr.Hex())
	r.Logger.Print("TestDappZEVM:   %s", r.ZevmTestDAppAddr.Hex())
	r.Logger.Print("TestDappEVM:    %s", r.EvmTestDAppAddr.Hex())
}

// IsRunningUpgrade returns true if the test is running an upgrade test suite.
func (r *E2ERunner) IsRunningUpgrade() bool {
	return os.Getenv(EnvKeyLocalnetMode) == LocalnetModeUpgrade
}

func (r *E2ERunner) IsRunningTssMigration() bool {
	return os.Getenv(EnvKeyLocalnetMode) == LocalnetModeUpgrade
}

func (r *E2ERunner) IsRunningUpgradeOrTSSMigration() bool {
	return os.Getenv(EnvKeyLocalnetMode) == LocalnetModeUpgrade ||
		os.Getenv(EnvKeyLocalnetMode) == LocalNetModeTSSMigration
}

// Errorf logs an error message. Mimics the behavior of testing.T.Errorf
func (r *E2ERunner) Errorf(format string, args ...any) {
	r.Logger.Error(format, args...)
}

// FailNow implemented to mimic the behavior of testing.T.FailNow
func (r *E2ERunner) FailNow() {
	err := fmt.Errorf("(*E2ERunner).FailNow for runner %q. Exiting", r.Name)

	r.Logger.Error("Failure: %s", err.Error())
	r.CtxCancel(err)

	// this panic ensures that the test routine exits fast.
	// it should be caught and handled gracefully so long
	// as the test is being executed by RunE2ETest().
	panic(err)
}

func (r *E2ERunner) requireTxSuccessful(receipt *ethtypes.Receipt, msgAndArgs ...any) {
	utils.RequireTxSuccessful(r, receipt, msgAndArgs...)
}

// EVMAddress is shorthand to get the EVM address of the account
func (r *E2ERunner) EVMAddress() ethcommon.Address {
	return r.Account.EVMAddress()
}

func (r *E2ERunner) ZetaAddress() types.AccAddress {
	return types.AccAddress(r.EVMAddress().Bytes())
}

func (r *E2ERunner) GetSolanaPrivKey() solana.PrivateKey {
	privkey, err := solana.PrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	require.NoError(r, err)
	return privkey
}

func (r *E2ERunner) GetZetacoredVersion() string {
	if r.zetacoredVersion != "" {
		return r.zetacoredVersion
	}
	nodeInfo, err := r.Clients.Zetacore.GetNodeInfo(r.Ctx)
	require.NoError(r, err, "get node info")
	r.zetacoredVersion = constant.NormalizeVersion(nodeInfo.ApplicationVersion.Version)
	return r.zetacoredVersion
}

func (r *E2ERunner) WorkDirPrefixed(path string) string {
	prefix := utils.WorkDir(r)
	return filepath.Join(prefix, path)
}
