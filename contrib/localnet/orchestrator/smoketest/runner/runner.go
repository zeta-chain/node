package runner

import (
	"sync"

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
	WG sync.WaitGroup
}

// SmokeTest is a function representing a smoke test
// It takes a SmokeTestRunner as an argument
type SmokeTest func(*SmokeTestRunner)

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

		WG: sync.WaitGroup{},
	}
}

// RunSmokeTests runs a list of smoke tests
func (sm *SmokeTestRunner) RunSmokeTests(smokeTests []SmokeTest) {
	for _, smokeTest := range smokeTests {
		sm.RunSmokeTest(smokeTest)
	}
}

// RunSmokeTest runs a smoke test
func (sm *SmokeTestRunner) RunSmokeTest(smokeTest SmokeTest) {
	// run smoke test
	smokeTest(sm)

	// check supplies
	sm.CheckZRC20ReserveAndSupply()
}
