package e2etests

import (
	"math/big"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestBitcoinWithdrawRBF tests the RBF (Replace-By-Fee) feature in Zetaclient.
// It needs block mining to be stopped and runs as the last test in the suite.
//
// IMPORTANT: the test requires to simulate a stuck tx in the Bitcoin regnet.
// The challenge to simulate a stuck tx is to create overwhelming traffic in the local mempool.
//
// To work around this:
//  1. change the 'minTxConfirmations' to 1 to not include outbound tx right away (production should use 0)
//     here: https://github.com/zeta-chain/node/blob/5c2a8ffbc702130fd9538b1cd7640d0e04d3e4f6/zetaclient/chains/bitcoin/observer/outbound.go#L27
//  2. stop block mining to let the pending tx sit in the mempool for longer time
func TestBitcoinWithdrawRBF(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse arguments
	defaultReceiver := r.GetBtcAddress().EncodeAddress()
	to, amount := utils.ParseBitcoinWithdrawArgs(r, args, defaultReceiver, r.GetBitcoinChainID())

	// initiate a withdraw CCTX, and wait for CCTX creation
	tx := r.WithdrawBTC(to, amount, gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)}, true)
	utils.WaitForZetaBlocks(r.Ctx, r, r.ZEVMClient, 1, 20*time.Second)
	cctx := utils.GetCCTXByInboundHash(r.Ctx, r.CctxClient, tx.Hash().Hex())

	// wait for the 1st outbound tracker hash to come in
	nonce := cctx.GetCurrentOutboundParam().TssNonce
	hashes := utils.WaitOutboundTracker(r.Ctx, r.CctxClient, r.GetBitcoinChainID(), nonce, 1, r.Logger, 3*time.Minute)
	txHash, err := chainhash.NewHashFromStr(hashes[0])
	r.Logger.Info("got 1st tracker hash: %s", txHash)

	// get original tx
	require.NoError(r, err)
	txResult, err := r.BtcRPCClient.GetTransaction(r.Ctx, txHash)
	require.NoError(r, err)
	require.Zero(r, txResult.Confirmations)

	// wait for RBF tx to kick in
	hashes = utils.WaitOutboundTracker(r.Ctx, r.CctxClient, r.GetBitcoinChainID(), nonce, 2, r.Logger, 3*time.Minute)
	txHashRBF, err := chainhash.NewHashFromStr(hashes[1])
	require.NoError(r, err)
	r.Logger.Info("got 2nd tracker hash: %s", txHashRBF)

	// resume block mining
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// waiting for CCTX to be mined
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// ensure the original tx is dropped
	utils.MustHaveDroppedBitcoinTx(r.Ctx, r.BtcRPCClient, txHash)

	// ensure the RBF tx is mined
	rawResult := utils.MustHaveMinedBitcoinTx(r.Ctx, r.BtcRPCClient, txHashRBF)

	// ensure RBF fee rate > old rate
	params := cctx.GetCurrentOutboundParam()
	oldRate, err := strconv.ParseInt(params.GasPrice, 10, 64)
	require.NoError(r, err)

	_, newRate, err := r.BtcRPCClient.GetTransactionFeeAndRate(r.Ctx, rawResult)
	require.NoError(r, err)
	require.Greater(r, newRate, oldRate, "RBF fee rate should be higher than the original tx")
}
