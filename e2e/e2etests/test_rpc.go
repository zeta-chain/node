package e2etests

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"
)

// TestRpc performs sanity checks on core JSON-RPC methods (eth_getTransactionByHash, eth_getTransactionReceipt,
// eth_getBlockByNumber, eth_getBlockByHash, debug_traceTransaction, debug_traceBlockByNumber).
//
// If no transaction hashes are provided via args, it seeds test transactions and discovers existing
// transactions from this user via CCTX queries (both ZEVM and EVM chains). This enables testing
// RPC methods against transactions submitted before upgrades to detect any regressions.
func TestRpc(r *runner.E2ERunner, args []string) {
	r.Logger.Info("starting JSON-RPC tests")

	rpcWrapper := NewZEVMRPC(r.ZEVMClient)

	txHashes := []string{}
	if len(args) == 1 && args[0] != "" {
		txHashes = strings.Split(args[0], ",")
	}

	if len(txHashes) == 0 {
		r.Logger.Info("seeding example transactions")
		amount := utils.ParseBigInt(r, "10000000000000000")

		// perform the deposit
		tx := r.ETHDeposit(r.EVMAddress(), amount, gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)}, true)

		// wait for the cctx to be mined
		cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
		r.Logger.CCTX(*cctx, "deposit")
		require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

		amount = utils.ParseBigInt(r, "9000000000000000")

		r.ApproveETHZRC20(r.GatewayZEVMAddr)

		// perform the withdraw
		tx = r.ETHWithdraw(r.EVMAddress(), amount, gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})

		// wait for the cctx to be mined
		cctx2 := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
		r.Logger.CCTX(*cctx2, "withdraw")
		require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

		cctxsRes, err := r.CctxClient.CctxAll(r.Ctx, &crosschaintypes.QueryAllCctxRequest{})
		require.NoError(r, err)

		zetaChainID, err := r.ZEVMClient.ChainID(r.Ctx)
		require.NoError(r, err)

		evmChainID, err := r.EVMClient.ChainID(r.Ctx)
		require.NoError(r, err)

		// this will also fetch txs submitted by this user before upgrade
		for _, cctx := range cctxsRes.CrossChainTx {
			if cctx.InboundParams.Sender == r.ZEVMAuth.From.Hex() && cctx.InboundParams.GetSenderChainId() == zetaChainID.Int64() {
				txHashes = append(txHashes, cctx.InboundParams.ObservedHash)
			}

			if cctx.InboundParams.Sender == r.EVMAuth.From.Hex() && cctx.InboundParams.GetSenderChainId() == evmChainID.Int64() {
				txHashes = append(txHashes, cctx.GetCurrentOutboundParam().Hash)
			}
		}
	}

	// testing just most important rpc methods, if there are some errors in rpc layer probably most will fail
	// just testing sanity checks, like presence of fields, so this can be executed for provided hashes on livenet
	for _, txHash := range txHashes {
		receipt, err := rpcWrapper.EthGetTransactionReceipt(r.Ctx, common.HexToHash(txHash))
		require.NoError(r, err)
		require.NotNil(r, receipt)
		require.NotEmpty(r, receipt.BlockHash)
		require.NotEmpty(r, receipt.BlockNumber)
		require.NotEmpty(r, receipt.TransactionHash)

		txInfo, err := rpcWrapper.EthGetTransactionByHash(r.Ctx, common.HexToHash(txHash))
		require.NoError(r, err)
		require.NotNil(r, txInfo)
		require.NotEmpty(r, txInfo.Hash)
		require.Equal(r, txHash, txInfo.Hash)
		require.NotEmpty(r, txInfo.From)
		require.NotEmpty(r, txInfo.Nonce)
		require.NotEmpty(r, txInfo.Gas)

		blockNumber, err := hexutil.DecodeBig(receipt.BlockNumber)
		require.NoError(r, err)

		blockByNumber, err := rpcWrapper.EthGetBlockByNumber(r.Ctx, blockNumber, false)
		require.NoError(r, err)
		require.NotNil(r, blockByNumber)
		require.Equal(r, receipt.BlockNumber, blockByNumber.Number)
		require.NotEmpty(r, blockByNumber.Hash)
		require.NotEmpty(r, blockByNumber.ParentHash)
		require.NotEmpty(r, blockByNumber.Timestamp)

		blockByHash, err := rpcWrapper.EthGetBlockByHash(r.Ctx, common.HexToHash(receipt.BlockHash), false)
		require.NoError(r, err)
		require.NotNil(r, blockByHash)
		require.Equal(r, receipt.BlockHash, blockByHash.Hash)
		require.NotEmpty(r, blockByHash.Number)
		require.NotEmpty(r, blockByHash.ParentHash)
		require.NotEmpty(r, blockByHash.Timestamp)

		traceTx, err := rpcWrapper.DebugTraceTransaction(r.Ctx, common.HexToHash(txHash))
		require.NoError(r, err)
		require.NotNil(r, traceTx)
		require.NotEmpty(r, traceTx.Type)
		require.NotEmpty(r, traceTx.From)
		require.NotEmpty(r, traceTx.To)

		traceBlock, err := rpcWrapper.DebugTraceBlockByNumber(r.Ctx, blockNumber)
		require.NoError(r, err)
		require.NotNil(r, traceBlock)
		require.GreaterOrEqual(r, len(traceBlock), 1)
		require.NotNil(r, traceBlock[0].Result)
		require.NotEmpty(r, traceBlock[0].Result.From)
		require.NotEmpty(r, traceBlock[0].Result.To)
	}

	r.Logger.Info("JSON-RPC tests completed successfully!")
}
