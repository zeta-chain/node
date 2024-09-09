package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestEtherWithdraw tests the withdrawal of ether
func TestEtherWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	withdrawalAmount, ok := new(big.Int).SetString(args[0], 10)
	require.True(r, ok, "Invalid withdrawal amount specified for TestEtherWithdraw.")

	// approve 1 unit of the gas token to cover the gas fee transfer
	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.ETHZRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)

	r.Logger.EVMTransaction(*tx, "approve")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.EVMReceipt(*receipt, "approve")

	// withdraw
	tx = r.WithdrawEther(withdrawalAmount)

	// verify the withdrawal value
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")

	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// Previous binary doesn't take EIP-1559 into account, so this will fail.
	// Thus, we need to skip this check for upgrade tests
	if !r.IsRunningUpgrade() {
		withdrawalReceipt := mustFetchEthReceipt(r, cctx)
		require.Equal(r, uint8(ethtypes.DynamicFeeTxType), withdrawalReceipt.Type, "receipt type mismatch")
	}

	r.Logger.Info("TestEtherWithdraw completed")
}

func mustFetchEthReceipt(r *runner.E2ERunner, cctx *crosschaintypes.CrossChainTx) *ethtypes.Receipt {
	hash := cctx.GetCurrentOutboundParam().Hash
	require.NotEmpty(r, hash, "outbound hash is empty")

	receipt, err := r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(hash))
	require.NoError(r, err)

	return receipt
}
