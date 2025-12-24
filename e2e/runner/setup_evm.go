package runner

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	erc20custodyv2 "github.com/zeta-chain/protocol-contracts/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"
	zetaconnnectornative "github.com/zeta-chain/protocol-contracts/pkg/zetaconnectornative.sol"

	"github.com/zeta-chain/node/e2e/contracts/erc1967proxy"
	"github.com/zeta-chain/node/e2e/contracts/erc20"
	"github.com/zeta-chain/node/e2e/contracts/testdappv2"
	"github.com/zeta-chain/node/e2e/utils"
)

// SetupEVM setup contracts on EVM with v2 contracts
func (r *E2ERunner) SetupEVM() {
	r.Logger.Info("⚙️ setting up EVM network")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("EVM setup took %s\n", time.Since(startTime))
	}()

	r.deployERC20Contract()
	r.donateTx()
	r.deployGatewayEVM()
	r.deployERC20Custody()
	r.DeployTestDAppV2EVM()
	r.DeployZetaConnectorNative()
	r.finalizeEVMSetup()
}

func (r *E2ERunner) ensureTxReceiptEVM(tx *ethtypes.Transaction, failMessage string) {
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt, failMessage)
}

// deployERC20Contract deploys the ERC20 test token contract
func (r *E2ERunner) deployERC20Contract() {
	r.Logger.Info("Deploying ERC20 contract")

	erc20Addr, txERC20, erc20contract, err := erc20.DeployERC20(r.EVMAuth, r.EVMClient, "TESTERC20", "TESTERC20", 6)
	require.NoError(r, err)

	r.ERC20 = erc20contract
	r.ERC20Addr = erc20Addr
	r.Logger.Info("ERC20 contract address: %s, tx hash: %s", erc20Addr.Hex(), txERC20.Hash().Hex())

	r.ensureTxReceiptEVM(txERC20, "ERC20 deployment failed")
}

// donateTx donates ether to TSS address to avoid account errors
func (r *E2ERunner) donateTx() {
	r.Logger.InfoLoud("Deploy Gateway and ERC20Custody ERC20\n")

	// Donate to the TSS address to avoid account errors because deploying gas token ZRC20 will automatically mint
	// gas token on ZetaChain to initialize the pool
	txDonation, err := r.DonateEtherToTSS(big.NewInt(101000000000000000))
	require.NoError(r, err)

	r.ensureTxReceiptEVM(txDonation, "EVM donation tx failed")
}

// deployGatewayEVM deploys the Gateway EVM contract with proxy
func (r *E2ERunner) deployGatewayEVM() {
	r.Logger.Info("Deploying Gateway EVM")

	gatewayEVMAddr, txGateway, _, err := gatewayevm.DeployGatewayEVM(r.EVMAuth, r.EVMClient)
	require.NoError(r, err)
	r.ensureTxReceiptEVM(txGateway, "Gateway deployment failed")

	gatewayEVMABI, err := gatewayevm.GatewayEVMMetaData.GetAbi()
	require.NoError(r, err)

	initializerData, err := gatewayEVMABI.Pack("initialize", r.TSSAddress, r.ZetaEthAddr, r.Account.EVMAddress())
	require.NoError(r, err)

	gatewayProxyAddress, gatewayProxyTx, _, err := erc1967proxy.DeployERC1967Proxy(
		r.EVMAuth,
		r.EVMClient,
		gatewayEVMAddr,
		initializerData,
	)
	require.NoError(r, err)
	r.ensureTxReceiptEVM(gatewayProxyTx, "Gateway proxy deployment failed")

	r.GatewayEVMAddr = gatewayProxyAddress
	r.GatewayEVM, err = gatewayevm.NewGatewayEVM(gatewayProxyAddress, r.EVMClient)
	require.NoError(r, err)
	r.Logger.Info("Gateway EVM contract address: %s, tx hash: %s", gatewayEVMAddr.Hex(), txGateway.Hash().Hex())

	updateAdditionalFeeTx, err := r.GatewayEVM.UpdateAdditionalActionFee(r.EVMAuth, big.NewInt(2e5))
	require.NoError(r, err)
	r.ensureTxReceiptEVM(updateAdditionalFeeTx, "Updating additional fee failed")
}

// deployERC20Custody deploys the ERC20Custody contract with proxy
func (r *E2ERunner) deployERC20Custody() {
	r.Logger.Info("Deploying ERC20Custody contract")

	erc20CustodyAddr, txCustody, _, err := erc20custodyv2.DeployERC20Custody(r.EVMAuth, r.EVMClient)
	require.NoError(r, err)
	r.ensureTxReceiptEVM(txCustody, "ERC20Custody deployment failed")

	erc20CustodyABI, err := erc20custodyv2.ERC20CustodyMetaData.GetAbi()
	require.NoError(r, err)

	initializerData, err := erc20CustodyABI.Pack("initialize", r.GatewayEVMAddr, r.TSSAddress, r.Account.EVMAddress())
	require.NoError(r, err)

	// Deploy erc20custody proxy contract
	erc20CustodyProxyAddress, erc20ProxyTx, _, err := erc1967proxy.DeployERC1967Proxy(
		r.EVMAuth,
		r.EVMClient,
		erc20CustodyAddr,
		initializerData,
	)
	require.NoError(r, err)
	r.ensureTxReceiptEVM(erc20ProxyTx, "ERC20Custody proxy deployment failed")

	// Initialize the custody contract instance
	r.ERC20CustodyAddr = erc20CustodyProxyAddress
	r.ERC20Custody, err = erc20custodyv2.NewERC20Custody(erc20CustodyProxyAddress, r.EVMClient)
	require.NoError(r, err)
	r.Logger.Info(
		"ERC20Custody contract address: %s, tx hash: %s",
		erc20CustodyAddr.Hex(),
		txCustody.Hash().Hex(),
	)

	// Set custody contract in gateway
	txSetCustody, err := r.GatewayEVM.SetCustody(r.EVMAuth, erc20CustodyProxyAddress)
	require.NoError(r, err)
	r.ensureTxReceiptEVM(txSetCustody, "Set custody in Gateway failed")
}

// DeployTestDAppV2EVM deploys the test DApp V2 contract for EVM
func (r *E2ERunner) DeployTestDAppV2EVM() {
	testDAppV2Addr, txTestDAppV2, _, err := testdappv2.DeployTestDAppV2(
		r.EVMAuth,
		r.EVMClient,
		false,
		r.GatewayEVMAddr,
		r.ZetaEthAddr,
	)
	require.NoError(r, err)
	r.ensureTxReceiptEVM(txTestDAppV2, "TestDAppV2 deployment failed")

	// Initialize the test dapp contract instance
	r.TestDAppV2EVMAddr = testDAppV2Addr
	r.TestDAppV2EVM, err = testdappv2.NewTestDAppV2(testDAppV2Addr, r.EVMClient)
	require.NoError(r, err)

	isZetaChain, err := r.TestDAppV2EVM.IsZetaChain(&bind.CallOpts{})
	require.NoError(r, err)
	require.False(r, isZetaChain)
}

// DeployZetaConnectorNative deploys the ZetaConnectorNative contract with proxy
func (r *E2ERunner) DeployZetaConnectorNative() {
	// Deploy zetaConnectorNative contract
	zetaConnectorNativeAddress, txZetaConnectorNativeHash, _, err := zetaconnnectornative.DeployZetaConnectorNative(
		r.EVMAuth,
		r.EVMClient,
	)
	require.NoError(r, err)
	r.ensureTxReceiptEVM(txZetaConnectorNativeHash, "ZetaConnectorNative deployment failed")

	// Get ABI for initialization
	zetaConnnectorNativeABI, err := zetaconnnectornative.ZetaConnectorNativeMetaData.GetAbi()
	require.NoError(r, err)

	// Encode the initializer data
	initializerData, err := zetaConnnectorNativeABI.Pack(
		"initialize",
		r.GatewayEVMAddr,
		r.ZetaEthAddr,
		r.TSSAddress,
		r.Account.EVMAddress(),
	)
	require.NoError(r, err)

	// Deploy zetaConnnectorNative proxy contract
	zetaConnnectorNativeProxyAddress, zetaConnnectorNativeProxyTx, _, err := erc1967proxy.DeployERC1967Proxy(
		r.EVMAuth,
		r.EVMClient,
		zetaConnectorNativeAddress,
		initializerData,
	)
	require.NoError(r, err)
	r.ensureTxReceiptEVM(zetaConnnectorNativeProxyTx, "ZetaConnectorNative proxy deployment failed")

	// Initialize the connector contract instance
	r.ConnectorNativeAddr = zetaConnnectorNativeProxyAddress
	r.ConnectorNative, err = zetaconnnectornative.NewZetaConnectorNative(zetaConnnectorNativeProxyAddress, r.EVMClient)
	require.NoError(r, err)

	// Set connector in gateway
	txSetConnector, err := r.GatewayEVM.SetConnector(r.EVMAuth, zetaConnnectorNativeProxyAddress)
	require.NoError(r, err)
	r.ensureTxReceiptEVM(txSetConnector, "Set connector in Gateway failed")

	r.Logger.Info(
		"ZetaConnectorNative contract address: %s, tx hash: %s",
		zetaConnectorNativeAddress.Hex(),
		txZetaConnectorNativeHash.Hash().Hex(),
	)
}

// finalizeEVMSetup performs final configuration steps
func (r *E2ERunner) finalizeEVMSetup() {
	// Whitelist the ERC20 token
	txWhitelist, err := r.ERC20Custody.Whitelist(r.EVMAuth, r.ERC20Addr)
	require.NoError(r, err)
	r.ensureTxReceiptEVM(txWhitelist, "ERC20 whitelist failed")

	// Set legacy supported (calling deposit directly in ERC20Custody)
	txSetLegacySupported, err := r.ERC20Custody.SetSupportsLegacy(r.EVMAuth, true)
	require.NoError(r, err)
	r.ensureTxReceiptEVM(txSetLegacySupported, "Set legacy support failed")

	// Grant PAUSER_ROLE to TSS address
	r.Logger.Info("Granting PAUSER_ROLE to TSS address")
	pauserRoleHash := crypto.Keccak256Hash([]byte("PAUSER_ROLE"))
	txGrantPauserRole, err := r.ERC20Custody.GrantRole(r.EVMAuth, pauserRoleHash, r.TSSAddress)
	require.NoError(r, err)
	r.ensureTxReceiptEVM(txGrantPauserRole, "Failed to grant PAUSER_ROLE to TSS address")
}
