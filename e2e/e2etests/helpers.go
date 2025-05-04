package e2etests

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// randomPayload generates a random payload to be used in gateway calls for testing purposes
func randomPayload(r *runner.E2ERunner) string {
	bytes := make([]byte, 50)
	_, err := rand.Read(bytes)
	require.NoError(r, err)

	return hex.EncodeToString(bytes)
}

func withdrawBTCZRC20(r *runner.E2ERunner, to btcutil.Address, amount *big.Int) *btcjson.TxRawResult {
	// approve and withdraw on ZRC20 contract
	receipt := approveAndWithdrawBTCZRC20(r, to, amount)

	// mine blocks if testing on regnet
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// get cctx and check status
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// get bitcoin tx according to the outTxHash in cctx
	outTxHash := cctx.GetCurrentOutboundParam().Hash
	hash, err := chainhash.NewHashFromStr(outTxHash)
	require.NoError(r, err)

	rawTx, err := r.BtcRPCClient.GetRawTransactionVerbose(r.Ctx, hash)
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

// approveAndWithdrawBTCZRC20 is a helper function to call 'approve' and 'withdraw' on BTCZRC20 contract
func approveAndWithdrawBTCZRC20(r *runner.E2ERunner, to btcutil.Address, amount *big.Int) *ethtypes.Receipt {
	// ensure enough balance to cover the withdrawal
	_, gasFee, err := r.BTCZRC20.WithdrawGasFee(&bind.CallOpts{})
	require.NoError(r, err)
	minimumAmount := new(big.Int).Add(amount, gasFee)
	currentBalance, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	require.NoError(r, err)
	require.Greater(
		r,
		currentBalance.Int64(),
		minimumAmount.Int64(),
		"current balance must be greater than amount + gasFee",
	)

	// approve more to cover withdraw fee
	tx, err := r.BTCZRC20.Approve(
		r.ZEVMAuth,
		r.BTCZRC20Addr,
		big.NewInt(amount.Int64()*2),
	)
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// withdraw 'amount' of BTC from ZRC20 to BTC address
	tx, err = r.BTCZRC20.Withdraw(r.ZEVMAuth, []byte(to.EncodeAddress()), amount)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	return receipt
}

// bigAdd is shorthand for new(big.Int).Add(x, y)
func bigAdd(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Add(x, y)
}

// bigSub is shorthand for new(big.Int).Sub(x, y)
func bigSub(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Sub(x, y)
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := d.Seconds() - float64(minutes*60)
	return fmt.Sprintf("%dm%.1fs", minutes, seconds)
}
