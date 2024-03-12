package runner

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/zeta-chain/zetacore/e2e/contracts/contextapp"
	"github.com/zeta-chain/zetacore/e2e/contracts/erc20"
	"github.com/zeta-chain/zetacore/e2e/contracts/zevmswap"
	"github.com/zeta-chain/zetacore/e2e/txserver"

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
	ZEVMClient   *ethclient.Client
	EVMClient    *ethclient.Client
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
	EVMAuth  *bind.TransactOpts
	ZEVMAuth *bind.TransactOpts

	// contracts
	ZetaEthAddr          ethcommon.Address
	ZetaEth              *zetaeth.ZetaEth
	ConnectorEthAddr     ethcommon.Address
	ConnectorEth         *zetaconnectoreth.ZetaConnectorEth
	ERC20CustodyAddr     ethcommon.Address
	ERC20Custody         *erc20custody.ERC20Custody
	ERC20Addr            ethcommon.Address
	ERC20                *erc20.ERC20
	ERC20ZRC20Addr       ethcommon.Address
	ERC20ZRC20           *zrc20.ZRC20
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
	evmClient *ethclient.Client,
	zevmClient *ethclient.Client,
	cctxClient crosschaintypes.QueryClient,
	zetaTxServer txserver.ZetaTxServer,
	fungibleClient fungibletypes.QueryClient,
	authClient authtypes.QueryClient,
	bankClient banktypes.QueryClient,
	observerClient observertypes.QueryClient,
	evmAuth *bind.TransactOpts,
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

		ZEVMClient:     zevmClient,
		EVMClient:      evmClient,
		ZetaTxServer:   zetaTxServer,
		CctxClient:     cctxClient,
		FungibleClient: fungibleClient,
		AuthClient:     authClient,
		BankClient:     bankClient,
		ObserverClient: observerClient,

		EVMAuth:      evmAuth,
		ZEVMAuth:     zevmAuth,
		BtcRPCClient: btcRPCClient,

		Logger: logger,

		WG: sync.WaitGroup{},
	}
}

// E2ETestFunc is a function representing a E2E test
// It takes a E2ERunner as an argument
type E2ETestFunc func(*E2ERunner, []string)

// E2ETest represents a E2E test with a name, args, description and test func
type E2ETest struct {
	Name           string
	Description    string
	Args           []string
	ArgsDefinition []ArgDefinition
	E2ETest        E2ETestFunc
}

// NewE2ETest creates a new instance of E2ETest with specified parameters.
func NewE2ETest(name, description string, argsDefinition []ArgDefinition, e2eTestFunc E2ETestFunc) E2ETest {
	return E2ETest{
		Name:           name,
		Description:    description,
		ArgsDefinition: argsDefinition,
		E2ETest:        e2eTestFunc,
		Args:           []string{},
	}
}

// ArgDefinition defines a structure for holding an argument's description along with it's default value.
type ArgDefinition struct {
	Description  string
	DefaultValue string
}

// DefaultArgs extracts and returns array of default arguments from the ArgsDefinition.
func (e E2ETest) DefaultArgs() []string {
	defaultArgs := make([]string, len(e.ArgsDefinition))
	for i, spec := range e.ArgsDefinition {
		defaultArgs[i] = spec.DefaultValue
	}
	return defaultArgs
}

// ArgsDescription returns a string representing the arguments description in a readable format.
func (e E2ETest) ArgsDescription() string {
	argsDescription := ""
	for _, def := range e.ArgsDefinition {
		argDesc := fmt.Sprintf("%s (%s)", def.Description, def.DefaultValue)
		if argsDescription != "" {
			argsDescription += ", "
		}
		argsDescription += argDesc
	}
	return argsDescription
}

// E2ETestRunConfig defines the basic configuration for initiating an E2E test, including its name and optional runtime arguments.
type E2ETestRunConfig struct {
	Name string
	Args []string
}

// GetE2ETestsToRunByName prepares a list of E2ETests to run based on given test names without arguments
func (runner *E2ERunner) GetE2ETestsToRunByName(availableTests []E2ETest, testNames ...string) ([]E2ETest, error) {
	tests := []E2ETestRunConfig{}
	for _, testName := range testNames {
		tests = append(tests, E2ETestRunConfig{
			Name: testName,
			Args: []string{},
		})
	}
	return runner.GetE2ETestsToRunByConfig(availableTests, tests)
}

// GetE2ETestsToRunByConfig prepares a list of E2ETests to run based on provided test names and their corresponding arguments
func (runner *E2ERunner) GetE2ETestsToRunByConfig(availableTests []E2ETest, testConfigs []E2ETestRunConfig) ([]E2ETest, error) {
	tests := []E2ETest{}
	for _, testSpec := range testConfigs {
		e2eTest, found := findE2ETestByName(availableTests, testSpec.Name)
		if !found {
			return nil, fmt.Errorf("e2e test %s not found", testSpec.Name)
		}
		e2eTestToRun := NewE2ETest(
			e2eTest.Name,
			e2eTest.Description,
			e2eTest.ArgsDefinition,
			e2eTest.E2ETest,
		)
		// update e2e test args
		e2eTestToRun.Args = testSpec.Args
		tests = append(tests, e2eTestToRun)
	}

	return tests, nil
}

// RunE2ETests runs a list of e2e tests
func (runner *E2ERunner) RunE2ETests(e2eTests []E2ETest) (err error) {
	for _, e2eTest := range e2eTests {
		if err := runner.RunE2ETest(e2eTest, true); err != nil {
			return err
		}
	}
	return nil
}

// RunE2ETestsFromNamesIntoReport runs a list of e2e tests by name in a list of e2e tests and returns a report
// The function doesn't return an error, it returns a report with the error
func (runner *E2ERunner) RunE2ETestsIntoReport(e2eTests []E2ETest) (TestReports, error) {
	// go through all tests
	reports := make(TestReports, 0, len(e2eTests))
	for _, test := range e2eTests {
		// get info before test
		balancesBefore, err := runner.GetAccountBalances(true)
		if err != nil {
			return nil, err
		}
		timeBefore := time.Now()

		// run test
		testErr := runner.RunE2ETest(test, false)
		if testErr != nil {
			runner.Logger.Print("test %s failed: %s", test.Name, testErr.Error())
		}

		// wait 5 sec to make sure we get updated balances
		time.Sleep(5 * time.Second)

		// get info after test
		balancesAfter, err := runner.GetAccountBalances(true)
		if err != nil {
			return nil, err
		}
		timeAfter := time.Now()

		// create report
		report := TestReport{
			Name:     test.Name,
			Success:  testErr == nil,
			Time:     timeAfter.Sub(timeBefore),
			GasSpent: GetAccountBalancesDiff(balancesBefore, balancesAfter),
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// RunE2ETest runs a e2e test
func (runner *E2ERunner) RunE2ETest(e2eTest E2ETest, checkAccounting bool) (err error) {
	// return an error on panic
	// https://github.com/zeta-chain/node/issues/1500
	defer func() {
		if r := recover(); r != nil {
			// print stack trace
			stack := make([]byte, 4096)
			n := runtime.Stack(stack, false)
			err = fmt.Errorf("%s failed: %v, stack trace %s", e2eTest.Name, r, stack[:n])
		}
	}()

	startTime := time.Now()
	runner.Logger.Print("⏳running - %s", e2eTest.Description)

	// run e2e test, if args are not provided, use default args
	args := e2eTest.Args
	if len(args) == 0 {
		args = e2eTest.DefaultArgs()
	}
	e2eTest.E2ETest(runner, args)

	//check supplies
	if checkAccounting {
		if err := runner.CheckZRC20ReserveAndSupply(); err != nil {
			return err
		}
	}

	runner.Logger.Print("✅ completed in %s - %s", time.Since(startTime), e2eTest.Description)

	return err
}

// findE2ETest finds a e2e test by name
func findE2ETestByName(e2eTests []E2ETest, e2eTestName string) (E2ETest, bool) {
	for _, test := range e2eTests {
		if test.Name == e2eTestName {
			return test, true
		}
	}
	return E2ETest{}, false
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
	runner.ERC20Addr = other.ERC20Addr
	runner.ERC20ZRC20Addr = other.ERC20ZRC20Addr
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
	runner.ZetaEth, err = zetaeth.NewZetaEth(runner.ZetaEthAddr, runner.EVMClient)
	if err != nil {
		return err
	}
	runner.ConnectorEth, err = zetaconnectoreth.NewZetaConnectorEth(runner.ConnectorEthAddr, runner.EVMClient)
	if err != nil {
		return err
	}
	runner.ERC20Custody, err = erc20custody.NewERC20Custody(runner.ERC20CustodyAddr, runner.EVMClient)
	if err != nil {
		return err
	}
	runner.ERC20, err = erc20.NewERC20(runner.ERC20Addr, runner.EVMClient)
	if err != nil {
		return err
	}
	runner.ERC20ZRC20, err = zrc20.NewZRC20(runner.ERC20ZRC20Addr, runner.ZEVMClient)
	if err != nil {
		return err
	}
	runner.ETHZRC20, err = zrc20.NewZRC20(runner.ETHZRC20Addr, runner.ZEVMClient)
	if err != nil {
		return err
	}
	runner.BTCZRC20, err = zrc20.NewZRC20(runner.BTCZRC20Addr, runner.ZEVMClient)
	if err != nil {
		return err
	}
	runner.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(runner.UniswapV2FactoryAddr, runner.ZEVMClient)
	if err != nil {
		return err
	}
	runner.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(runner.UniswapV2RouterAddr, runner.ZEVMClient)
	if err != nil {
		return err
	}
	runner.ConnectorZEVM, err = connectorzevm.NewZetaConnectorZEVM(runner.ConnectorZEVMAddr, runner.ZEVMClient)
	if err != nil {
		return err
	}
	runner.WZeta, err = wzeta.NewWETH9(runner.WZetaAddr, runner.ZEVMClient)
	if err != nil {
		return err
	}

	runner.ZEVMSwapApp, err = zevmswap.NewZEVMSwapApp(runner.ZEVMSwapAppAddr, runner.ZEVMClient)
	if err != nil {
		return err
	}
	runner.ContextApp, err = contextapp.NewContextApp(runner.ContextAppAddr, runner.ZEVMClient)
	if err != nil {
		return err
	}
	runner.SystemContract, err = systemcontract.NewSystemContract(runner.SystemContractAddr, runner.ZEVMClient)
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
	runner.Logger.Print("ERC20ZRC20:     %s", runner.ERC20ZRC20Addr.Hex())
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
	runner.Logger.Print("ERC20:      %s", runner.ERC20Addr.Hex())
}
