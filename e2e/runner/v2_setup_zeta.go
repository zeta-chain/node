package runner

import (
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/zetacore/e2e/contracts/erc1967proxy"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

// SetZEVMContractsV2 set contracts for the ZEVM
func (r *E2ERunner) SetZEVMContractsV2() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage)
	}

	r.Logger.Print("⚙️ setting up ZEVM v2 network")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("ZEVM v2 network took %s\n", time.Since(startTime))
	}()

	r.Logger.Info("Deploying Gateway ZEVM")
	gatewayZEVMAddr, txGateway, _, err := gatewayzevm.DeployGatewayZEVM(r.EVMAuth, r.EVMClient)
	require.NoError(r, err)

	ensureTxReceipt(txGateway, "Gateway deployment failed")

	gatewayZEVMABI, err := gatewayzevm.GatewayZEVMMetaData.GetAbi()
	require.NoError(r, err)

	// Encode the initializer data
	initializerData, err := gatewayZEVMABI.Pack("initialize", r.WZetaAddr)
	require.NoError(r, err)

	// Deploy the proxy contract
	proxyAddress, txProxy, _, err := erc1967proxy.DeployERC1967Proxy(
		r.EVMAuth,
		r.EVMClient,
		gatewayZEVMAddr,
		initializerData,
	)
	require.NoError(r, err)

	r.GatewayZEVMAddr = proxyAddress
	r.GatewayZEVM, err = gatewayzevm.NewGatewayZEVM(proxyAddress, r.EVMClient)
	require.NoError(r, err)
	r.Logger.Info("Gateway ZEVM contract address: %s, tx hash: %s", gatewayZEVMAddr.Hex(), txGateway.Hash().Hex())

	// Set the gateway address in the protocol
	err = r.ZetaTxServer.UpdateGatewayAddress(utils.AdminPolicyName, r.GatewayZEVMAddr.Hex())
	require.NoError(r, err)

	ensureTxReceipt(txProxy, "Gateway proxy deployment failed")
}
