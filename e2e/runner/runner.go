package runner

import (
	"context"
	"sync"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	zetaeth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zeta.eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/connectorzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/wzeta.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/zeta-chain/zetacore/e2e/contracts/contextapp"
	"github.com/zeta-chain/zetacore/e2e/contracts/erc20"
	"github.com/zeta-chain/zetacore/e2e/contracts/zevmswap"
	"github.com/zeta-chain/zetacore/e2e/txserver"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// E2ERunner stores all the clients and addresses needed for E2E test
// Exposes a method to run E2E test
// It also provides some helper functions
type E2ERunner struct {
	// accounts
	DeployerAddress       ethcommon.Address
	DeployerPrivateKey    string
	TSSAddress            ethcommon.Address
	BTCTSSAddress         btcutil.Address
	BTCDeployerAddress    *btcutil.AddressWitnessPubKeyHash
	FungibleAdminMnemonic string

	// rpc clients
	ZevmClient   *ethclient.Client
	GoerliClient *ethclient.Client
	BtcRPCClient *rpcclient.Client

	// grpc clients
	CctxClient     crosschaintypes.QueryClient
	FungibleClient fungibletypes.QueryClient
	AuthClient     authtypes.QueryClient
	BankClient     banktypes.QueryClient
	ObserverClient observertypes.QueryClient

	// zeta client
	ZetaTxServer txserver.ZetaTxServer

	// evm auth
	GoerliAuth *bind.TransactOpts
	ZevmAuth   *bind.TransactOpts

	// contracts
	ZetaEthAddr          ethcommon.Address
	ZetaEth              *zetaeth.ZetaEth
	ConnectorEthAddr     ethcommon.Address
	ConnectorEth         *zetaconnectoreth.ZetaConnectorEth
	ERC20CustodyAddr     ethcommon.Address
	ERC20Custody         *erc20custody.ERC20Custody
	USDTERC20Addr        ethcommon.Address
	USDTERC20            *erc20.USDT
	USDTZRC20Addr        ethcommon.Address
	USDTZRC20            *zrc20.ZRC20
	ETHZRC20Addr         ethcommon.Address
	ETHZRC20             *zrc20.ZRC20
	BTCZRC20Addr         ethcommon.Address
	BTCZRC20             *zrc20.ZRC20
	UniswapV2FactoryAddr ethcommon.Address
	UniswapV2Factory     *uniswapv2factory.UniswapV2Factory
	UniswapV2RouterAddr  ethcommon.Address
	UniswapV2Router      *uniswapv2router.UniswapV2Router02
	ConnectorZEVMAddr    ethcommon.Address
	ConnectorZEVM        *connectorzevm.ZetaConnectorZEVM
	WZetaAddr            ethcommon.Address
	WZeta                *wzeta.WETH9

	TestDAppAddr       ethcommon.Address
	ZEVMSwapAppAddr    ethcommon.Address
	ZEVMSwapApp        *zevmswap.ZEVMSwapApp
	ContextAppAddr     ethcommon.Address
	ContextApp         *contextapp.ContextApp
	SystemContractAddr ethcommon.Address
	SystemContract     *systemcontract.SystemContract

	// config
	CctxTimeout    time.Duration
	ReceiptTimeout time.Duration

	// other
	Name          string
	Ctx           context.Context
	CtxCancel     context.CancelFunc
	Logger        *Logger
	WG            sync.WaitGroup
	BitcoinParams *chaincfg.Params
	mutex         sync.Mutex
}

func NewE2ERunner(
	ctx context.Context,
	name string,
	ctxCancel context.CancelFunc,
	deployerAddress ethcommon.Address,
	deployerPrivateKey string,
	fungibleAdminMnemonic string,
	goerliClient *ethclient.Client,
	zevmClient *ethclient.Client,
	cctxClient crosschaintypes.QueryClient,
	zetaTxServer txserver.ZetaTxServer,
	fungibleClient fungibletypes.QueryClient,
	authClient authtypes.QueryClient,
	bankClient banktypes.QueryClient,
	observerClient observertypes.QueryClient,
	goerliAuth *bind.TransactOpts,
	zevmAuth *bind.TransactOpts,
	btcRPCClient *rpcclient.Client,
	logger *Logger,
) *E2ERunner {
	return &E2ERunner{
		Name:      name,
		Ctx:       ctx,
		CtxCancel: ctxCancel,

		DeployerAddress:       deployerAddress,
		DeployerPrivateKey:    deployerPrivateKey,
		FungibleAdminMnemonic: fungibleAdminMnemonic,

		ZevmClient:     zevmClient,
		GoerliClient:   goerliClient,
		ZetaTxServer:   zetaTxServer,
		CctxClient:     cctxClient,
		FungibleClient: fungibleClient,
		AuthClient:     authClient,
		BankClient:     bankClient,
		ObserverClient: observerClient,

		GoerliAuth:   goerliAuth,
		ZevmAuth:     zevmAuth,
		BtcRPCClient: btcRPCClient,

		Logger: logger,

		WG: sync.WaitGroup{},
	}
}

// CopyAddressesFrom copies addresses from another E2ETestRunner that initialized the contracts
func (runner *E2ERunner) CopyAddressesFrom(other *E2ERunner) (err error) {
	// copy TSS address
	runner.TSSAddress = other.TSSAddress
	runner.BTCTSSAddress = other.BTCTSSAddress

	// copy addresses
	runner.ZetaEthAddr = other.ZetaEthAddr
	runner.ConnectorEthAddr = other.ConnectorEthAddr
	runner.ERC20CustodyAddr = other.ERC20CustodyAddr
	runner.USDTERC20Addr = other.USDTERC20Addr
	runner.USDTZRC20Addr = other.USDTZRC20Addr
	runner.ETHZRC20Addr = other.ETHZRC20Addr
	runner.BTCZRC20Addr = other.BTCZRC20Addr
	runner.UniswapV2FactoryAddr = other.UniswapV2FactoryAddr
	runner.UniswapV2RouterAddr = other.UniswapV2RouterAddr
	runner.ConnectorZEVMAddr = other.ConnectorZEVMAddr
	runner.WZetaAddr = other.WZetaAddr
	runner.TestDAppAddr = other.TestDAppAddr
	runner.ZEVMSwapAppAddr = other.ZEVMSwapAppAddr
	runner.ContextAppAddr = other.ContextAppAddr
	runner.SystemContractAddr = other.SystemContractAddr

	// create instances of contracts
	runner.ZetaEth, err = zetaeth.NewZetaEth(runner.ZetaEthAddr, runner.GoerliClient)
	if err != nil {
		return err
	}
	runner.ConnectorEth, err = zetaconnectoreth.NewZetaConnectorEth(runner.ConnectorEthAddr, runner.GoerliClient)
	if err != nil {
		return err
	}
	runner.ERC20Custody, err = erc20custody.NewERC20Custody(runner.ERC20CustodyAddr, runner.GoerliClient)
	if err != nil {
		return err
	}
	runner.USDTERC20, err = erc20.NewUSDT(runner.USDTERC20Addr, runner.GoerliClient)
	if err != nil {
		return err
	}
	runner.USDTZRC20, err = zrc20.NewZRC20(runner.USDTZRC20Addr, runner.ZevmClient)
	if err != nil {
		return err
	}
	runner.ETHZRC20, err = zrc20.NewZRC20(runner.ETHZRC20Addr, runner.ZevmClient)
	if err != nil {
		return err
	}
	runner.BTCZRC20, err = zrc20.NewZRC20(runner.BTCZRC20Addr, runner.ZevmClient)
	if err != nil {
		return err
	}
	runner.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(runner.UniswapV2FactoryAddr, runner.ZevmClient)
	if err != nil {
		return err
	}
	runner.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(runner.UniswapV2RouterAddr, runner.ZevmClient)
	if err != nil {
		return err
	}
	runner.ConnectorZEVM, err = connectorzevm.NewZetaConnectorZEVM(runner.ConnectorZEVMAddr, runner.ZevmClient)
	if err != nil {
		return err
	}
	runner.WZeta, err = wzeta.NewWETH9(runner.WZetaAddr, runner.ZevmClient)
	if err != nil {
		return err
	}

	runner.ZEVMSwapApp, err = zevmswap.NewZEVMSwapApp(runner.ZEVMSwapAppAddr, runner.ZevmClient)
	if err != nil {
		return err
	}
	runner.ContextApp, err = contextapp.NewContextApp(runner.ContextAppAddr, runner.ZevmClient)
	if err != nil {
		return err
	}
	runner.SystemContract, err = systemcontract.NewSystemContract(runner.SystemContractAddr, runner.ZevmClient)
	if err != nil {
		return err
	}
	return nil
}

// Lock locks the mutex
func (runner *E2ERunner) Lock() {
	runner.mutex.Lock()
}

// Unlock unlocks the mutex
func (runner *E2ERunner) Unlock() {
	runner.mutex.Unlock()
}

// PrintContractAddresses prints the addresses of the contracts
// the printed contracts are grouped in a zevm and evm section
// there is a padding used to print the addresses at the same position
func (runner *E2ERunner) PrintContractAddresses() {
	// zevm contracts
	runner.Logger.Print(" --- 📜zEVM contracts ---")
	runner.Logger.Print("SystemContract: %s", runner.SystemContractAddr.Hex())
	runner.Logger.Print("ETHZRC20:       %s", runner.ETHZRC20Addr.Hex())
	runner.Logger.Print("USDTZRC20:      %s", runner.USDTZRC20Addr.Hex())
	runner.Logger.Print("BTCZRC20:       %s", runner.BTCZRC20Addr.Hex())
	runner.Logger.Print("UniswapFactory: %s", runner.UniswapV2FactoryAddr.Hex())
	runner.Logger.Print("UniswapRouter:  %s", runner.UniswapV2RouterAddr.Hex())
	runner.Logger.Print("ConnectorZEVM:  %s", runner.ConnectorZEVMAddr.Hex())
	runner.Logger.Print("WZeta:          %s", runner.WZetaAddr.Hex())

	runner.Logger.Print("ZEVMSwapApp:    %s", runner.ZEVMSwapAppAddr.Hex())
	runner.Logger.Print("ContextApp:     %s", runner.ContextAppAddr.Hex())
	runner.Logger.Print("TestDapp:       %s", runner.TestDAppAddr.Hex())

	// evm contracts
	runner.Logger.Print(" --- 📜EVM contracts ---")
	runner.Logger.Print("ZetaEth:        %s", runner.ZetaEthAddr.Hex())
	runner.Logger.Print("ConnectorEth:   %s", runner.ConnectorEthAddr.Hex())
	runner.Logger.Print("ERC20Custody:   %s", runner.ERC20CustodyAddr.Hex())
	runner.Logger.Print("USDTERC20:      %s", runner.USDTERC20Addr.Hex())
}
