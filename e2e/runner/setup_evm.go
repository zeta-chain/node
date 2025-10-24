package runner

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	erc20custodyv2 "github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"
	zetaconnnectornative "github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnectornative.sol"

	"github.com/zeta-chain/node/e2e/contracts/erc1967proxy"
	"github.com/zeta-chain/node/e2e/contracts/erc20"
	"github.com/zeta-chain/node/e2e/contracts/testdappv2"
	"github.com/zeta-chain/node/e2e/utils"
)

// SetupEVM setup contracts on EVM with v2 contracts
func (r *E2ERunner) SetupEVM() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage)
	}

	r.Logger.Info("⚙️ setting up EVM network")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("EVM setup took %s\n", time.Since(startTime))
	}()

	r.Logger.Info("Deploying ERC20 contract")
	erc20Addr, txERC20, erc20contract, err := erc20.DeployERC20(r.EVMAuth, r.EVMClient, "TESTERC20", "TESTERC20", 6)
	require.NoError(r, err)

	r.ERC20 = erc20contract
	r.ERC20Addr = erc20Addr
	r.Logger.Info("ERC20 contract address: %s, tx hash: %s", erc20Addr.Hex(), txERC20.Hash().Hex())

	r.Logger.InfoLoud("Deploy Gateway and ERC20Custody ERC20\n")

	// donate to the TSS address to avoid account errors because deploying gas token ZRC20 will automatically mint
	// gas token on ZetaChain to initialize the pool d
	txDonation, err := r.DonateEtherToTSS(big.NewInt(101000000000000000))
	require.NoError(r, err)

	r.Logger.Info("Deploying Gateway EVM")
	gatewayEVMAddr, txGateway, _, err := gatewayevm.DeployGatewayEVM(r.EVMAuth, r.EVMClient)
	require.NoError(r, err)

	ensureTxReceipt(txGateway, "Gateway deployment failed")

	gatewayEVMABI, err := gatewayevm.GatewayEVMMetaData.GetAbi()
	require.NoError(r, err)

	// Encode the initializer data
	initializerData, err := gatewayEVMABI.Pack("initialize", r.TSSAddress, r.ZetaEthAddr, r.Account.EVMAddress())
	require.NoError(r, err)

	// Deploy gateway proxy contract
	gatewayProxyAddress, gatewayProxyTx, _, err := erc1967proxy.DeployERC1967Proxy(
		r.EVMAuth,
		r.EVMClient,
		gatewayEVMAddr,
		initializerData,
	)
	require.NoError(r, err)
	ensureTxReceipt(gatewayProxyTx, "Gateway proxy deployment failed")

	r.GatewayEVMAddr = gatewayProxyAddress
	r.GatewayEVM, err = gatewayevm.NewGatewayEVM(gatewayProxyAddress, r.EVMClient)
	require.NoError(r, err)
	r.Logger.Info("Gateway EVM contract address: %s, tx hash: %s", gatewayEVMAddr.Hex(), txGateway.Hash().Hex())

	updateAdditionalFeeTx, err := r.GatewayEVM.UpdateAdditionalActionFee(r.EVMAuth, big.NewInt(2e5))
	require.NoError(r, err)
	ensureTxReceipt(updateAdditionalFeeTx, "Updating additional fee failed")

	// Deploy erc20custody proxy contract
	r.Logger.Info("Deploying ERC20Custody contract")
	erc20CustodyAddr, txCustody, _, err := erc20custodyv2.DeployERC20Custody(r.EVMAuth, r.EVMClient)
	require.NoError(r, err)

	ensureTxReceipt(txCustody, "ERC20Custody deployment failed")

	erc20CustodyABI, err := erc20custodyv2.ERC20CustodyMetaData.GetAbi()
	require.NoError(r, err)

	// Encode the initializer data
	initializerData, err = erc20CustodyABI.Pack("initialize", r.GatewayEVMAddr, r.TSSAddress, r.Account.EVMAddress())
	require.NoError(r, err)

	// Deploy erc20custody proxy contract
	erc20CustodyProxyAddress, erc20ProxyTx, _, err := erc1967proxy.DeployERC1967Proxy(
		r.EVMAuth,
		r.EVMClient,
		erc20CustodyAddr,
		initializerData,
	)
	require.NoError(r, err)

	r.ERC20CustodyAddr = erc20CustodyProxyAddress
	r.ERC20Custody, err = erc20custodyv2.NewERC20Custody(erc20CustodyProxyAddress, r.EVMClient)
	require.NoError(r, err)
	r.Logger.Info(
		"ERC20Custody contract address: %s, tx hash: %s",
		erc20CustodyAddr.Hex(),
		txCustody.Hash().Hex(),
	)

	ensureTxReceipt(txCustody, "ERC20Custody deployment failed")

	// set custody contract in gateway
	txSetCustody, err := r.GatewayEVM.SetCustody(r.EVMAuth, erc20CustodyProxyAddress)
	require.NoError(r, err)

	r.DeployTestDAppV2(ensureTxReceipt)
	// Deploy zetaConnectorNative contract
	r.DeployZetaConnectorNative(ensureTxReceipt)

	// check contract deployment receipt
	ensureTxReceipt(txERC20, "ERC20 deployment failed")
	ensureTxReceipt(txDonation, "EVM donation tx failed")
	ensureTxReceipt(gatewayProxyTx, "Gateway proxy deployment failed")
	ensureTxReceipt(erc20ProxyTx, "ERC20Custody proxy deployment failed")
	ensureTxReceipt(txSetCustody, "Set custody in Gateway failed")

	// check isZetaChain is false
	isZetaChain, err := r.TestDAppV2EVM.IsZetaChain(&bind.CallOpts{})
	require.NoError(r, err)
	require.False(r, isZetaChain)

	// whitelist the ERC20
	txWhitelist, err := r.ERC20Custody.Whitelist(r.EVMAuth, r.ERC20Addr)
	require.NoError(r, err)

	// set legacy supported (calling deposit directly in ERC20Custody)
	txSetLegacySupported, err := r.ERC20Custody.SetSupportsLegacy(r.EVMAuth, true)
	require.NoError(r, err)

	ensureTxReceipt(txWhitelist, "ERC20 whitelist failed")
	ensureTxReceipt(txSetLegacySupported, "Set legacy support failed")

	r.Logger.Info("Granting PAUSER_ROLE to TSS address")
	pauserRoleHash := crypto.Keccak256Hash([]byte("PAUSER_ROLE"))
	txGrantPauserRole, err := r.ERC20Custody.GrantRole(r.EVMAuth, pauserRoleHash, r.TSSAddress)
	require.NoError(r, err)

	ensureTxReceipt(txGrantPauserRole, "Failed to grant PAUSER_ROLE to TSS address")
}

func (r *E2ERunner) DeployTestDAppV2(ensureTxReceipt func(tx *ethtypes.Transaction, failMessage string)) {
	// deploy test dapp v2
	testDAppV2Addr, txTestDAppV2, _, err := testdappv2.DeployTestDAppV2(r.EVMAuth, r.EVMClient, false, r.GatewayEVMAddr)
	require.NoError(r, err)

	r.TestDAppV2EVMAddr = testDAppV2Addr
	r.TestDAppV2EVM, err = testdappv2.NewTestDAppV2(testDAppV2Addr, r.EVMClient)
	require.NoError(r, err)

	ensureTxReceipt(txTestDAppV2, "TestDAppV2 deployment failed")
}

// DeployZetaConnectorNative deploys the ZetaConnectorNative contract with proxy
func (r *E2ERunner) DeployZetaConnectorNative(ensureTxReceipt func(tx *ethtypes.Transaction, failMessage string)) {
	// Deploy zetaConnectorNative contract
	zetaConnectorNativeAddress, txZetaConnectorNativeHash, _, err := zetaconnnectornative.DeployZetaConnectorNative(
		r.EVMAuth,
		r.EVMClient,
	)
	require.NoError(r, err)
	ensureTxReceipt(txZetaConnectorNativeHash, "ZetaConnectorNative deployment failed")

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
	ensureTxReceipt(zetaConnnectorNativeProxyTx, "ZetaConnectorNative proxy deployment failed")

	// Initialize the connector contract instance
	r.ConnectorNativeAddr = zetaConnnectorNativeProxyAddress
	r.ConnectorNative, err = zetaconnnectornative.NewZetaConnectorNative(zetaConnnectorNativeProxyAddress, r.EVMClient)
	require.NoError(r, err)

	// Set connector in gateway
	txSetConnector, err := r.GatewayEVM.SetConnector(r.EVMAuth, zetaConnnectorNativeProxyAddress)
	require.NoError(r, err)
	ensureTxReceipt(txSetConnector, "Set connector in Gateway failed")

	r.Logger.Info(
		"ZetaConnectorNative contract address: %s, tx hash: %s",
		zetaConnectorNativeAddress.Hex(),
		txZetaConnectorNativeHash.Hash().Hex(),
	)
}

// DeployTestDAppV2EVM deploys the test DApp V2 contract for EVM
func (r *E2ERunner) DeployTestDAppV2EVM() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage)
	}

	testDAppV2Addr, txTestDAppV2, _, err := testdappv2.DeployTestDAppV2(
		r.EVMAuth,
		r.EVMClient,
		false,
		r.GatewayEVMAddr,
	)
	require.NoError(r, err)
	ensureTxReceipt(txTestDAppV2, "TestDAppV2 deployment failed")

	// Initialize the test dapp contract instance
	r.TestDAppV2EVMAddr = testDAppV2Addr
	r.TestDAppV2EVM, err = testdappv2.NewTestDAppV2(testDAppV2Addr, r.EVMClient)
	require.NoError(r, err)

	isZetaChain, err := r.TestDAppV2EVM.IsZetaChain(&bind.CallOpts{})
	require.NoError(r, err)
	require.False(r, isZetaChain)
}
