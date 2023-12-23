package runner

import (
	"fmt"
	"sync"
	"time"

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
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/contextapp"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/erc20"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/zevmswap"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/txserver"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// SmokeTestRunner stores all the clients and addresses needed for smoke test
// Exposes a method to run smoke test
// It also provides some helper functions
type SmokeTestRunner struct {
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
	TestDAppAddr         ethcommon.Address
	ZEVMSwapAppAddr      ethcommon.Address
	ZEVMSwapApp          *zevmswap.ZEVMSwapApp
	ContextAppAddr       ethcommon.Address
	ContextApp           *contextapp.ContextApp
	SystemContractAddr   ethcommon.Address
	SystemContract       *systemcontract.SystemContract

	// other
	Logger *Logger
	WG     sync.WaitGroup
	mutex  sync.Mutex
}

func NewSmokeTestRunner(
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
) *SmokeTestRunner {
	return &SmokeTestRunner{
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

// SmokeTestFunc is a function representing a smoke test
// It takes a SmokeTestRunner as an argument
type SmokeTestFunc func(*SmokeTestRunner)

// SmokeTest represents a smoke test with a name
type SmokeTest struct {
	Name        string
	Description string
	SmokeTest   SmokeTestFunc
}

// RunSmokeTestsFromNames runs a list of smoke tests by name in a list of smoke tests
func (sm *SmokeTestRunner) RunSmokeTestsFromNames(smokeTests []SmokeTest, smokeTestNames ...string) error {
	for _, smokeTestName := range smokeTestNames {
		smokeTest, ok := findSmokeTest(smokeTestName, smokeTests)
		if !ok {
			return fmt.Errorf("smoke test %s not found", smokeTestName)
		}
		if err := sm.RunSmokeTest(smokeTest); err != nil {
			return err
		}
	}

	return nil
}

// RunSmokeTests runs a list of smoke tests
func (sm *SmokeTestRunner) RunSmokeTests(smokeTests []SmokeTest) (err error) {
	for _, smokeTest := range smokeTests {
		if err := sm.RunSmokeTest(smokeTest); err != nil {
			return err
		}
	}
	return nil
}

// RunSmokeTest runs a smoke test
func (sm *SmokeTestRunner) RunSmokeTest(smokeTestWithName SmokeTest) (err error) {
	// return an error on panic
	// https://github.com/zeta-chain/node/issues/1500
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s failed: %v", smokeTestWithName.Name, r)
		}
	}()

	startTime := time.Now()
	sm.Logger.Print("⏳running - %s", smokeTestWithName.Description)

	// run smoke test
	smokeTestWithName.SmokeTest(sm)

	// check supplies
	//sm.CheckZRC20ReserveAndSupply()

	sm.Logger.Print("✅ completed in %s - %s", time.Since(startTime), smokeTestWithName.Description)

	return err
}

// findSmokeTest finds a smoke test by name
func findSmokeTest(name string, smokeTests []SmokeTest) (SmokeTest, bool) {
	for _, test := range smokeTests {
		if test.Name == name {
			return test, true
		}
	}
	return SmokeTest{}, false
}

// CopyAddressesFrom copies addresses from another SmokeTestRunner that initialized the contracts
func (sm *SmokeTestRunner) CopyAddressesFrom(other *SmokeTestRunner) (err error) {
	// copy TSS address
	sm.TSSAddress = other.TSSAddress
	sm.BTCTSSAddress = other.BTCTSSAddress

	// copy addresses
	sm.ZetaEthAddr = other.ZetaEthAddr
	sm.ConnectorEthAddr = other.ConnectorEthAddr
	sm.ERC20CustodyAddr = other.ERC20CustodyAddr
	sm.USDTERC20Addr = other.USDTERC20Addr
	sm.USDTZRC20Addr = other.USDTZRC20Addr
	sm.ETHZRC20Addr = other.ETHZRC20Addr
	sm.BTCZRC20Addr = other.BTCZRC20Addr
	sm.UniswapV2FactoryAddr = other.UniswapV2FactoryAddr
	sm.UniswapV2RouterAddr = other.UniswapV2RouterAddr
	sm.TestDAppAddr = other.TestDAppAddr
	sm.ZEVMSwapAppAddr = other.ZEVMSwapAppAddr
	sm.ContextAppAddr = other.ContextAppAddr
	sm.SystemContractAddr = other.SystemContractAddr

	// create instances of contracts
	sm.ZetaEth, err = zetaeth.NewZetaEth(sm.ZetaEthAddr, sm.GoerliClient)
	if err != nil {
		return err
	}
	sm.ConnectorEth, err = zetaconnectoreth.NewZetaConnectorEth(sm.ConnectorEthAddr, sm.GoerliClient)
	if err != nil {
		return err
	}
	sm.ERC20Custody, err = erc20custody.NewERC20Custody(sm.ERC20CustodyAddr, sm.GoerliClient)
	if err != nil {
		return err
	}
	sm.USDTERC20, err = erc20.NewUSDT(sm.USDTERC20Addr, sm.GoerliClient)
	if err != nil {
		return err
	}
	sm.USDTZRC20, err = zrc20.NewZRC20(sm.USDTZRC20Addr, sm.ZevmClient)
	if err != nil {
		return err
	}
	sm.ETHZRC20, err = zrc20.NewZRC20(sm.ETHZRC20Addr, sm.ZevmClient)
	if err != nil {
		return err
	}
	sm.BTCZRC20, err = zrc20.NewZRC20(sm.BTCZRC20Addr, sm.ZevmClient)
	if err != nil {
		return err
	}
	sm.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(sm.UniswapV2FactoryAddr, sm.ZevmClient)
	if err != nil {
		return err
	}
	sm.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(sm.UniswapV2RouterAddr, sm.ZevmClient)
	if err != nil {
		return err
	}
	sm.ZEVMSwapApp, err = zevmswap.NewZEVMSwapApp(sm.ZEVMSwapAppAddr, sm.ZevmClient)
	if err != nil {
		return err
	}
	sm.ContextApp, err = contextapp.NewContextApp(sm.ContextAppAddr, sm.ZevmClient)
	if err != nil {
		return err
	}
	sm.SystemContract, err = systemcontract.NewSystemContract(sm.SystemContractAddr, sm.ZevmClient)
	if err != nil {
		return err
	}
	return nil
}

// Lock locks the mutex
func (sm *SmokeTestRunner) Lock() {
	sm.mutex.Lock()
}

// Unlock unlocks the mutex
func (sm *SmokeTestRunner) Unlock() {
	sm.mutex.Unlock()
}
