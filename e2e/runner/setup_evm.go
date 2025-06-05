package runner

import (
	"math/big"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/pkg/chains"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
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

	r.GatewayEVMAddr = gatewayProxyAddress
	r.GatewayEVM, err = gatewayevm.NewGatewayEVM(gatewayProxyAddress, r.EVMClient)
	require.NoError(r, err)
	r.Logger.Info("Gateway EVM contract address: %s, tx hash: %s", gatewayEVMAddr.Hex(), txGateway.Hash().Hex())

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

	// deploy test dapp v2
	testDAppV2Addr, txTestDAppV2, _, err := testdappv2.DeployTestDAppV2(r.EVMAuth, r.EVMClient, false, r.GatewayEVMAddr)
	require.NoError(r, err)

	r.TestDAppV2EVMAddr = testDAppV2Addr
	r.TestDAppV2EVM, err = testdappv2.NewTestDAppV2(testDAppV2Addr, r.EVMClient)
	require.NoError(r, err)

	// check contract deployment receipt
	ensureTxReceipt(txERC20, "ERC20 deployment failed")
	ensureTxReceipt(txDonation, "EVM donation tx failed")
	ensureTxReceipt(gatewayProxyTx, "Gateway proxy deployment failed")
	ensureTxReceipt(erc20ProxyTx, "ERC20Custody proxy deployment failed")
	ensureTxReceipt(txSetCustody, "Set custody in Gateway failed")
	ensureTxReceipt(txTestDAppV2, "TestDAppV2 deployment failed")

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
}

func (r *E2ERunner) SetupConnectorV2() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage)
	}

	r.Logger.Info("⚙️ setting up EVM network")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("EVM setup took %s\n", time.Since(startTime))
	}()

	zetaConnnectorNativeAddress, txZetaConnnectorNativeHash, _, err := zetaconnnectornative.DeployZetaConnectorNative(
		r.EVMAuth,
		r.EVMClient,
	)
	require.NoError(r, err)

	ensureTxReceipt(txZetaConnnectorNativeHash, "ZetaConnectorNative deployment failed")

	zetaConnnectorNativeABI, err := zetaconnnectornative.ZetaConnectorNativeMetaData.GetAbi()
	require.NoError(r, err)

	// Encode the initializer data
	initializerData, err := zetaConnnectorNativeABI.Pack("initialize", r.GatewayEVMAddr, r.ZetaEthAddr, r.TSSAddress, r.Account.EVMAddress())
	require.NoError(r, err)

	// Deploy zetaConnnectorNative proxy contract
	zetaConnnectorNativeProxyAddress, zetaConnnectorNativeProxyTx, _, err := erc1967proxy.DeployERC1967Proxy(
		r.EVMAuth,
		r.EVMClient,
		zetaConnnectorNativeAddress,
		initializerData,
	)
	require.NoError(r, err)

	// check contract deployment receipt
	ensureTxReceipt(zetaConnnectorNativeProxyTx, "ZetaConnectorNative proxy deployment failed")
	r.ConnectorNativeAddr = zetaConnnectorNativeAddress
	r.ConnectorNative, err = zetaconnnectornative.NewZetaConnectorNative(zetaConnnectorNativeProxyAddress, r.EVMClient)
	require.NoError(r, err)

	r.Logger.Print(
		"ZetaConnectorNative contract address: %s, tx hash: %s",
		zetaConnnectorNativeAddress.Hex(),
		txZetaConnnectorNativeHash.Hash().Hex(),
	)

	r.Logger.Print("ZetaConnectorNative contract address: %s", r.ConnectorNativeAddr.Hex())
}

func (r *E2ERunner) MigrateConnector() {
	r.Logger.Print("⚙️ Migrating ZetaConnectorNative to v2")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("ZetaConnectorNative migration took %s\n", time.Since(startTime))
	}()
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage)
	}

	pauseV2Tx, err := r.ConnectorNative.Pause(r.EVMAuth)
	require.NoError(r, err)
	ensureTxReceipt(pauseV2Tx, "ZetaConnectorNative pause failed")

	pauseV1Tx, err := r.ConnectorEth.Pause(r.EVMAuth)
	require.NoError(r, err)
	ensureTxReceipt(pauseV1Tx, "ZetaConnectorEth pause failed")

	balance, err := r.ConnectorEth.GetLockedAmount(&bind.CallOpts{})
	require.NoError(r, err, "GetLockedAmount failed")

	r.Logger.Print("ZetaConnectorEth locked amount: %s", balance.String())

	msgMigrateConnectorFunds := crosschaintypes.NewMsgMigrateConnectorFunds(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		chains.GoerliLocalnet.ChainId,
		r.ConnectorNativeAddr.Hex(),
		sdkmath.NewUintFromBigInt(balance),
	)

	res, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msgMigrateConnectorFunds)
	require.NoError(r, err)

	event, ok := txserver.EventOfType[*crosschaintypes.EventConnectorFundsMigration](res.Events)
	require.True(r, ok, "no EventERC20CustodyFundsMigration in %s", res.TxHash)

	cctxRes, err := r.CctxClient.Cctx(r.Ctx, &crosschaintypes.QueryGetCctxRequest{Index: event.CctxIndex})
	require.NoError(r, err)

	cctx := cctxRes.CrossChainTx
	r.Logger.CCTX(*cctx, "migration")

	// wait for the cctx to be mined
	r.WaitForMinedCCTXFromIndex(event.CctxIndex)

	r.Logger.Print("CCTX %s migrated to new connector %s", cctx.Index, r.ConnectorNativeAddr.Hex())

	// check if the new connector has the funds
	newConnectorBalance, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.ConnectorNativeAddr)
	require.NoError(r, err, "BalanceOf failed for new connector")

	// Verify that the migration was successful
	require.Equal(r, balance, newConnectorBalance,
		"Migration failed: old connector balance (%s) != new connector balance (%s)",
		balance.String(), newConnectorBalance.String())

	r.Logger.Print("✅ Migration verification successful: %s ZETA tokens migrated", newConnectorBalance.String())
}
