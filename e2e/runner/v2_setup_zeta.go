package runner

import (
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/contracts/erc1967proxy"
	"github.com/zeta-chain/node/pkg/contracts/testdappv2"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// SetZEVMContractsV2 set contracts for the ZEVM
func (r *E2ERunner) SetZEVMContractsV2() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage+" tx hash: "+tx.Hash().Hex())
	}

	r.Logger.Print("⚙️ setting up ZEVM v2 network")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("ZEVM v2 network took %s\n", time.Since(startTime))
	}()

	r.Logger.Info("Deploying Gateway ZEVM")
	gatewayZEVMAddr, txGateway, _, err := gatewayzevm.DeployGatewayZEVM(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	ensureTxReceipt(txGateway, "Gateway deployment failed")

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
	proxyAddress, txProxy, _, err := erc1967proxy.DeployERC1967Proxy(
		r.ZEVMAuth,
		r.ZEVMClient,
		gatewayZEVMAddr,
		initializerData,
	)
	require.NoError(r, err)

	r.GatewayZEVMAddr = proxyAddress
	r.GatewayZEVM, err = gatewayzevm.NewGatewayZEVM(proxyAddress, r.ZEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Gateway ZEVM contract address: %s, tx hash: %s", gatewayZEVMAddr.Hex(), txGateway.Hash().Hex())

	// Set the gateway address in the protocol
	err = r.ZetaTxServer.UpdateGatewayAddress(utils.AdminPolicyName, r.GatewayZEVMAddr.Hex())
	require.NoError(r, err)

	// deploy test dapp v2
	testDAppV2Addr, txTestDAppV2, _, err := testdappv2.DeployTestDAppV2(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	r.TestDAppV2ZEVMAddr = testDAppV2Addr
	r.TestDAppV2ZEVM, err = testdappv2.NewTestDAppV2(testDAppV2Addr, r.ZEVMClient)
	require.NoError(r, err)

	ensureTxReceipt(txProxy, "Gateway proxy deployment failed")
	ensureTxReceipt(txTestDAppV2, "TestDAppV2 deployment failed")
}

// UpdateChainParamsV2Contracts update the erc20 custody contract and gateway address in the chain params
// this operation is used when transitioning to new smart contract architecture where a new ERC20 custody contract is deployed
func (r *E2ERunner) UpdateChainParamsV2Contracts() {
	res, err := r.ObserverClient.GetChainParams(r.Ctx, &observertypes.QueryGetChainParamsRequest{})
	require.NoError(r, err)

	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	// find old chain params
	var (
		chainParams *observertypes.ChainParams
		found       bool
	)
	for _, cp := range res.ChainParams.ChainParams {
		if cp.ChainId == evmChainID.Int64() {
			chainParams = cp
			found = true
			break
		}
	}
	require.True(r, found, "Chain params not found for chain id %d", evmChainID)

	// update with the new ERC20 custody contract address
	chainParams.Erc20CustodyContractAddress = r.ERC20CustodyV2Addr.Hex()

	// update with the new gateway address
	chainParams.GatewayAddress = r.GatewayEVMAddr.Hex()

	// update the chain params
	_, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, observertypes.NewMsgUpdateChainParams(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		chainParams,
	))
	require.NoError(r, err)
}
