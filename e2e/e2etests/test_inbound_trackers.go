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
// IMPORTANT: the test requires inbound observation to be disabled, the following line should be uncommented:
// https://github.com/zeta-chain/node/blob/9dcb42729653e033f5ba60a77dc37e5e19b092ad/zetaclient/chains/evm/observer/inbound.go#L210
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
	r.Logger.Print("🏃test v1 eth deposit")
	txHash := r.LegacyDepositEtherWithAmount(amount)
	addTrackerAndWaitForCCTX(coin.CoinType_Gas, txHash.Hex())
	r.Logger.Print("🍾v1 eth deposit observed")

	// send v1 erc20 deposit
	r.Logger.Print("🏃test v1 erc20 deposit")
	txHash = r.LegacyDepositERC20WithAmountAndMessage(r.EVMAddress(), amount, []byte{})
	addTrackerAndWaitForCCTX(coin.CoinType_ERC20, txHash.Hex())
	r.Logger.Print("🍾v1 erc20 deposit observed")

	// send v2 eth deposit
	r.Logger.Print("🏃test v2 eth deposit")
	tx := r.ETHDeposit(r.EVMAddress(), amount, gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})
	addTrackerAndWaitForCCTX(coin.CoinType_Gas, tx.Hash().Hex())
	r.Logger.Print("🍾v2 eth deposit observed")

	// send v2 eth deposit and call
	r.Logger.Print("🏃test v2 eth eposit and call")
	tx = r.ETHDepositAndCall(
		r.TestDAppV2ZEVMAddr,
		amount,
		[]byte(randomPayload(r)),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)
	addTrackerAndWaitForCCTX(coin.CoinType_Gas, tx.Hash().Hex())
	r.Logger.Print("🍾v2 eth deposit and call observed")

	// send v2 erc20 deposit
	r.Logger.Print("🏃test v2 erc20 deposit")
	r.ApproveERC20OnEVM(r.GatewayEVMAddr)
	tx = r.ERC20Deposit(r.EVMAddress(), amount, gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})
	addTrackerAndWaitForCCTX(coin.CoinType_Gas, tx.Hash().Hex())
	r.Logger.Print("🍾v2 erc20 deposit observed")

	// send v2 erc20 deposit and call
	r.Logger.Print("🏃test v2 erc20 deposit and call")
	tx = r.ERC20DepositAndCall(
		r.TestDAppV2ZEVMAddr,
		amount,
		[]byte(randomPayload(r)),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)
	addTrackerAndWaitForCCTX(coin.CoinType_Gas, tx.Hash().Hex())
	r.Logger.Print("🍾v2 erc20 deposit and call observed")

	// send v2 call
	r.Logger.Print("🏃test v2 call")
	tx = r.EVMToZEMVCall(
		r.TestDAppV2ZEVMAddr,
		[]byte(randomPayload(r)),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)
	addTrackerAndWaitForCCTX(coin.CoinType_NoAssetCall, tx.Hash().Hex())
	r.Logger.Print("🍾v2 call observed")
}
