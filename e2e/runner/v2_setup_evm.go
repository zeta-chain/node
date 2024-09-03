package runner

import (
	"math/big"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	erc20custodyv2 "github.com/zeta-chain/protocol-contracts/v2/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/contracts/erc1967proxy"
	"github.com/zeta-chain/node/pkg/contracts/testdappv2"
)

// SetupEVMV2 setup contracts on EVM with v2 contracts
func (r *E2ERunner) SetupEVMV2() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage)
	}

	r.Logger.Info("⚙️ setting up EVM v2 network")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("EVM v2 setup took %s\n", time.Since(startTime))
	}()

	r.Logger.InfoLoud("Deploy Gateway and ERC20Custody ERC20\n")

	// donate to the TSS address to avoid account errors because deploying gas token ZRC20 will automatically mint
	// gas token on ZetaChain to initialize the pool
	txDonation, err := r.SendEther(r.TSSAddress, big.NewInt(101000000000000000), []byte(constant.DonationMessage))
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

	// Deploy the proxy contract
	proxyAddress, txProxy, _, err := erc1967proxy.DeployERC1967Proxy(
		r.EVMAuth,
		r.EVMClient,
		gatewayEVMAddr,
		initializerData,
	)
	require.NoError(r, err)

	r.GatewayEVMAddr = proxyAddress
	r.GatewayEVM, err = gatewayevm.NewGatewayEVM(proxyAddress, r.EVMClient)
	require.NoError(r, err)
	r.Logger.Info("Gateway EVM contract address: %s, tx hash: %s", gatewayEVMAddr.Hex(), txGateway.Hash().Hex())

	r.Logger.Info("Deploying ERC20Custody contract")
	erc20CustodyNewAddr, txCustody, erc20CustodyNew, err := erc20custodyv2.DeployERC20Custody(
		r.EVMAuth,
		r.EVMClient,
		r.GatewayEVMAddr,
		r.TSSAddress,
		r.Account.EVMAddress(),
	)
	require.NoError(r, err)

	r.ERC20CustodyV2Addr = erc20CustodyNewAddr
	r.ERC20CustodyV2 = erc20CustodyNew
	r.Logger.Info(
		"ERC20CustodyV2 contract address: %s, tx hash: %s",
		erc20CustodyNewAddr.Hex(),
		txCustody.Hash().Hex(),
	)

	ensureTxReceipt(txCustody, "ERC20CustodyV2 deployment failed")

	// set custody contract in gateway
	txSetCustody, err := r.GatewayEVM.SetCustody(r.EVMAuth, erc20CustodyNewAddr)
	require.NoError(r, err)

	// deploy test dapp v2
	testDAppV2Addr, txTestDAppV2, _, err := testdappv2.DeployTestDAppV2(r.EVMAuth, r.EVMClient)
	require.NoError(r, err)

	r.TestDAppV2EVMAddr = testDAppV2Addr
	r.TestDAppV2EVM, err = testdappv2.NewTestDAppV2(testDAppV2Addr, r.EVMClient)
	require.NoError(r, err)

	// check contract deployment receipt
	ensureTxReceipt(txDonation, "EVM donation tx failed")
	ensureTxReceipt(txProxy, "Gateway proxy deployment failed")
	ensureTxReceipt(txSetCustody, "Set custody in Gateway failed")
	ensureTxReceipt(txTestDAppV2, "TestDAppV2 deployment failed")

	// whitelist the ERC20
	txWhitelist, err := r.ERC20CustodyV2.Whitelist(r.EVMAuth, r.ERC20Addr)
	require.NoError(r, err)

	// set legacy supported (calling deposit directly in ERC20Custody)
	txSetLegacySupported, err := r.ERC20CustodyV2.SetSupportsLegacy(r.EVMAuth, true)
	require.NoError(r, err)

	ensureTxReceipt(txWhitelist, "ERC20 whitelist failed")
	ensureTxReceipt(txSetLegacySupported, "Set legacy support failed")
}
