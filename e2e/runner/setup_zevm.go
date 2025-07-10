package runner

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/coreregistry.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/wzeta.sol"
	connectorzevm "github.com/zeta-chain/protocol-contracts/pkg/zetaconnectorzevm.sol"

	"github.com/zeta-chain/node/e2e/contracts/erc1967proxy"
	"github.com/zeta-chain/node/e2e/contracts/testdappv2"
	e2eutils "github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/contracts/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/zeta-chain/node/pkg/contracts/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

// SetupZEVM setup protocol contracts for the ZEVM
func (r *E2ERunner) SetupZEVM() {
	r.Logger.Print("⚙️ setting up ZEVM protocol contracts")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("ZEVM protocol contracts took %s\n", time.Since(startTime))
	}()

	// Deploy system contracts and setup core protocol components
	r.deploySystemContracts()
	r.setupUniswapContracts()
	r.setupSystemContract()
	r.deployGatewayZEVM()
	r.deployCoreRegistry()
	r.deployTestDAppV2()
}

// ensureTxReceipt is a helper function to ensure transaction success
func (r *E2ERunner) ensureTxReceipt(tx *ethtypes.Transaction, failMessage string) {
	receipt := e2eutils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt, failMessage+" tx hash: "+tx.Hash().Hex())
}

// deploySystemContracts deploys the core system contracts and ZRC20 contracts
func (r *E2ERunner) deploySystemContracts() {
	addresses, err := r.ZetaTxServer.DeploySystemContracts(
		e2eutils.OperationalPolicyName,
		e2eutils.AdminPolicyName,
	)
	require.NoError(r, err)

	// Store addresses for later use
	r.UniswapV2FactoryAddr = ethcommon.HexToAddress(addresses.UniswapV2FactoryAddr)
	r.UniswapV2RouterAddr = ethcommon.HexToAddress(addresses.UniswapV2RouterAddr)
	r.ConnectorZEVMAddr = ethcommon.HexToAddress(addresses.ZEVMConnectorAddr)
	r.WZetaAddr = ethcommon.HexToAddress(addresses.WZETAAddr)
}

// setupUniswapContracts initializes Uniswap factory and router contracts
func (r *E2ERunner) setupUniswapContracts() {
	var err error

	// Setup UniswapV2Factory
	r.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(r.UniswapV2FactoryAddr, r.ZEVMClient)
	require.NoError(r, err)

	// Setup UniswapV2Router
	r.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(r.UniswapV2RouterAddr, r.ZEVMClient)
	require.NoError(r, err)

	// Setup ZevmConnector
	r.ConnectorZEVM, err = connectorzevm.NewZetaConnectorZEVM(r.ConnectorZEVMAddr, r.ZEVMClient)
	require.NoError(r, err)

	// Setup WZeta
	r.WZeta, err = wzeta.NewWETH9(r.WZetaAddr, r.ZEVMClient)
	require.NoError(r, err)
}

// setupSystemContract queries and initializes the system contract
func (r *E2ERunner) setupSystemContract() {
	// Query system contract address from the chain
	systemContractRes, err := r.FungibleClient.SystemContract(
		r.Ctx,
		&fungibletypes.QueryGetSystemContractRequest{},
	)
	require.NoError(r, err)

	systemContractAddr := ethcommon.HexToAddress(systemContractRes.SystemContract.SystemContract)
	systemContract, err := systemcontract.NewSystemContract(
		systemContractAddr,
		r.ZEVMClient,
	)
	require.NoError(r, err)

	r.SystemContract = systemContract
	r.SystemContractAddr = systemContractAddr
}

// deployGatewayZEVM deploys the Gateway ZEVM contract with proxy
func (r *E2ERunner) deployGatewayZEVM() {
	r.Logger.Info("Deploying Gateway ZEVM")

	// Deploy the gateway implementation
	gatewayZEVMAddr, txGateway, _, err := gatewayzevm.DeployGatewayZEVM(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.ensureTxReceipt(txGateway, "Gateway deployment failed")

	// Get ABI for initialization
	gatewayZEVMABI, err := gatewayzevm.GatewayZEVMMetaData.GetAbi()
	require.NoError(r, err)

	// Encode the initializer data
	initializerData, err := gatewayZEVMABI.Pack("initialize", r.WZetaAddr, r.Account.EVMAddress())
	require.NoError(r, err)

	// Deploy the proxy contract
	r.Logger.Info(
		"Deploying proxy with %s and %s, address: %s",
		r.WZetaAddr.Hex(),
		r.Account.EVMAddress().Hex(),
		gatewayZEVMAddr.Hex(),
	)
	proxyAddress, txProxyGatewayZEVM, _, err := erc1967proxy.DeployERC1967Proxy(
		r.ZEVMAuth,
		r.ZEVMClient,
		gatewayZEVMAddr,
		initializerData,
	)
	require.NoError(r, err)
	r.ensureTxReceipt(txProxyGatewayZEVM, "GatewayZEVM proxy deployment failed")

	// Initialize the gateway contract instance
	r.GatewayZEVMAddr = proxyAddress
	r.GatewayZEVM, err = gatewayzevm.NewGatewayZEVM(proxyAddress, r.ZEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Gateway ZEVM contract address: %s, tx hash: %s", gatewayZEVMAddr.Hex(), txGateway.Hash().Hex())

	// Set the gateway address in the protocol
	err = r.ZetaTxServer.UpdateGatewayAddress(e2eutils.AdminPolicyName, r.GatewayZEVMAddr.Hex())
	require.NoError(r, err)
}

// deployCoreRegistry deploys the CoreRegistry contract with proxy
func (r *E2ERunner) deployCoreRegistry() {
	r.Logger.Print("Deploying CoreRegistry contract")

	// Deploy the registry implementation
	coreRegistryAddr, txCoreRegistry, _, err := coreregistry.DeployCoreRegistry(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.ensureTxReceipt(txCoreRegistry, "CoreRegistry deployment failed")

	// Get ABI for initialization
	coreRegistryABI, err := coreregistry.CoreRegistryMetaData.GetAbi()
	require.NoError(r, err)

	// Encode the initializer data
	initializerData, err := coreRegistryABI.Pack(
		"initialize",
		r.Account.EVMAddress(),
		r.Account.EVMAddress(),
		r.GatewayZEVMAddr,
	)
	require.NoError(r, err)

	// Deploy the proxy contract for the CoreRegistry
	r.Logger.Info(
		"Deploying CoreRegistry proxy with admin %s gateway %s, address: %s",
		r.Account.EVMAddress().Hex(),
		r.GatewayZEVMAddr.Hex(),
		coreRegistryAddr.Hex(),
	)
	proxyCoreRegistryAddr, txProxyCoreRegistry, _, err := erc1967proxy.DeployERC1967Proxy(
		r.ZEVMAuth,
		r.ZEVMClient,
		coreRegistryAddr,
		initializerData,
	)
	require.NoError(r, err)
	r.ensureTxReceipt(txProxyCoreRegistry, "CoreRegistry proxy deployment failed")

	// Initialize the registry contract instance
	r.CoreRegistryAddr = proxyCoreRegistryAddr
	r.CoreRegistry, err = coreregistry.NewCoreRegistry(proxyCoreRegistryAddr, r.ZEVMClient)
	require.NoError(r, err)

	// Update the gateway with the registry address
	updateRegistryTx, err := r.GatewayZEVM.SetRegistryAddress(r.ZEVMAuth, proxyCoreRegistryAddr)
	require.NoError(r, err)
	r.ensureTxReceipt(updateRegistryTx, "Gateway set registry address failed")
}

// deployTestDAppV2 deploys the test DApp V2 contract
func (r *E2ERunner) deployTestDAppV2() {
	// Deploy test dapp v2
	testDAppV2Addr, txTestDAppV2, _, err := testdappv2.DeployTestDAppV2(
		r.ZEVMAuth,
		r.ZEVMClient,
		true,
		r.GatewayEVMAddr,
	)
	require.NoError(r, err)
	r.ensureTxReceipt(txTestDAppV2, "TestDAppV2 deployment failed")

	// Initialize the test dapp contract instance
	r.TestDAppV2ZEVMAddr = testDAppV2Addr
	r.TestDAppV2ZEVM, err = testdappv2.NewTestDAppV2(testDAppV2Addr, r.ZEVMClient)
	require.NoError(r, err)

	// Verify isZetaChain is true
	isZetaChain, err := r.TestDAppV2ZEVM.IsZetaChain(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, isZetaChain)
}
