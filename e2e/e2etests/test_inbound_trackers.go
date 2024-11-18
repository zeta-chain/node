package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestInboundTrackers tests inbound trackers processing in ZetaClient
// It run deposits, send inbound trackers and check cctxs are mined
// IMPORTANT: the test requires inbound observation to be disabled, ob.WatchInbound bg process should be commented out
// https://github.com/zeta-chain/node/blob/b20c3f15decf1de85e3c7192852b07f98f8dbb8c/zetaclient/chains/evm/observer/observer.go#L184
func TestInboundTrackers(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	amount := big.NewInt(1e17)

	addTrackerAndWaitForCCTX := func(coinType coin.CoinType, txHash string) {
		r.AddInboundTracker(coinType, txHash)
		cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash, r.CctxClient, r.Logger, r.CctxTimeout)
		require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
		r.Logger.CCTX(*cctx, "cctx")
	}

	// send v1 eth deposit
	r.Logger.Info("test v1 eth deposit")
	txHash := r.DepositEtherWithAmount(amount)
	addTrackerAndWaitForCCTX(coin.CoinType_Gas, txHash.Hex())

	// send v1 erc20 deposit
	r.Logger.Info("test v1 erc20 deposit")
	txHash = r.DepositERC20WithAmountAndMessage(r.EVMAddress(), amount, []byte{})
	addTrackerAndWaitForCCTX(coin.CoinType_ERC20, txHash.Hex())

	// send v2 deposit
	r.Logger.Info("test v2 deposit")
	tx := r.V2ETHDeposit(r.EVMAddress(), amount, gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})
	addTrackerAndWaitForCCTX(coin.CoinType_Gas, tx.Hash().Hex())

	// send v2 deposit and call
	r.Logger.Info("test v2 deposit and call")
	tx = r.V2ERC20DepositAndCall(
		r.TestDAppV2ZEVMAddr,
		amount,
		[]byte(randomPayload(r)),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)
	addTrackerAndWaitForCCTX(coin.CoinType_Gas, tx.Hash().Hex())

	// send v2 call
	r.Logger.Info("test v2 call")
	tx = r.V2EVMToZEMVCall(
		r.TestDAppV2ZEVMAddr,
		[]byte(randomPayload(r)),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)
	addTrackerAndWaitForCCTX(coin.CoinType_NoAssetCall, tx.Hash().Hex())
}
