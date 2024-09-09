package e2etests

import (
	"math/big"
	"strconv"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func withdrawBTCZRC20(r *runner.E2ERunner, to btcutil.Address, amount *big.Int) *btcjson.TxRawResult {
	tx, err := r.BTCZRC20.Approve(
		r.ZEVMAuth,
		r.BTCZRC20Addr,
		big.NewInt(amount.Int64()*2),
	) // approve more to cover withdraw fee
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// mine blocks if testing on regnet
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// withdraw 'amount' of BTC from ZRC20 to BTC address
	tx, err = r.BTCZRC20.Withdraw(r.ZEVMAuth, []byte(to.EncodeAddress()), amount)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// mine 10 blocks to confirm the withdrawal tx
	_, err = r.GenerateToAddressIfLocalBitcoin(10, to)
	require.NoError(r, err)

	// get cctx and check status
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// get bitcoin tx according to the outTxHash in cctx
	outTxHash := cctx.GetCurrentOutboundParam().Hash
	hash, err := chainhash.NewHashFromStr(outTxHash)
	require.NoError(r, err)

	rawTx, err := r.BtcRPCClient.GetRawTransactionVerbose(hash)
	require.NoError(r, err)

	r.Logger.Info("raw tx:")
	r.Logger.Info("  TxIn: %d", len(rawTx.Vin))
	for idx, txIn := range rawTx.Vin {
		r.Logger.Info("  TxIn %d:", idx)
		r.Logger.Info("    TxID:Vout:  %s:%d", txIn.Txid, txIn.Vout)
		r.Logger.Info("    ScriptSig: %s", txIn.ScriptSig.Hex)
	}
	r.Logger.Info("  TxOut: %d", len(rawTx.Vout))
	for _, txOut := range rawTx.Vout {
		r.Logger.Info("  TxOut %d:", txOut.N)
		r.Logger.Info("    Value: %.8f", txOut.Value)
		r.Logger.Info("    ScriptPubKey: %s", txOut.ScriptPubKey.Hex)
	}

	return rawTx
}

// verifyTransferAmountFromCCTX verifies the transfer amount from the CCTX on EVM
func verifyTransferAmountFromCCTX(r *runner.E2ERunner, cctx *crosschaintypes.CrossChainTx, amount int64) {
	r.Logger.Info("outTx hash %s", cctx.GetCurrentOutboundParam().Hash)

	receipt, err := r.EVMClient.TransactionReceipt(
		r.Ctx,
		ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash),
	)
	require.NoError(r, err)

	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	for _, log := range receipt.Logs {
		event, err := r.ERC20.ParseTransfer(*log)
		if err != nil {
			continue
		}
		r.Logger.Info("  logs: from %s, to %s, value %d", event.From.Hex(), event.To.Hex(), event.Value)
		require.Equal(r, amount, event.Value.Int64(), "value is not correct")
	}
}

// verifySolanaWithdrawalAmountFromCCTX verifies the withdrawn amount on Solana for given CCTX
func verifySolanaWithdrawalAmountFromCCTX(r *runner.E2ERunner, cctx *crosschaintypes.CrossChainTx, amount uint64) {
	txHash := cctx.GetCurrentOutboundParam().Hash
	r.Logger.Info("outbound hash %s", txHash)

	// convert txHash to signature
	sig, err := solana.SignatureFromBase58(txHash)
	require.NoError(r, err)

	// query transaction by signature
	txResult, err := r.SolanaClient.GetTransaction(r.Ctx, sig, &rpc.GetTransactionOpts{})
	require.NoError(r, err)

	// unmarshal transaction
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(r, err)

	// 1st instruction is the withdraw
	instruction := tx.Message.Instructions[0]
	instWithdrae, err := solanacontracts.ParseInstructionWithdraw(instruction)
	require.NoError(r, err)

	// verify the amount
	require.Equal(r, amount, instWithdrae.TokenAmount(), "withdraw amount is not correct")
}

// Parse helpers ==========================================>

func parseFloat(t require.TestingT, s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	require.NoError(t, err, "unable to parse float %q", s)
	return f
}

func parseInt(t require.TestingT, s string) int {
	v, err := strconv.Atoi(s)
	require.NoError(t, err, "unable to parse int from %q", s)

	return v
}

func parseBigInt(t require.TestingT, s string) *big.Int {
	v, ok := big.NewInt(0).SetString(s, 10)
	require.True(t, ok, "unable to parse big.Int from %q", s)

	return v
}

// bigIntFromFloat64 takes float64 (e.g. 0.001) that represents btc amount
// and converts it to big.Int for downstream usage.
func btcAmountFromFloat64(t require.TestingT, amount float64) *big.Int {
	satoshi, err := btcutil.NewAmount(amount)
	require.NoError(t, err)

	return big.NewInt(int64(satoshi))
}

func parseBitcoinWithdrawArgs(r *runner.E2ERunner, args []string, defaultReceiver string) (btcutil.Address, *big.Int) {
	require.NotEmpty(r, args, "args list is empty")

	receiverRaw := defaultReceiver
	if args[0] != "" {
		receiverRaw = args[0]
	}

	receiver, err := chains.DecodeBtcAddress(receiverRaw, r.GetBitcoinChainID())
	require.NoError(r, err, "unable to decode btc address")

	withdrawalAmount := parseFloat(r, args[1])
	amount := btcAmountFromFloat64(r, withdrawalAmount)

	return receiver, amount
}
