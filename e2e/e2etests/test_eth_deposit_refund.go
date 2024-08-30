package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestEtherDepositAndCallRefund(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the deposit amount
	value := parseBigInt(r, args[0])

	evmClient := r.EVMClient

	nonce, err := evmClient.PendingNonceAt(r.Ctx, r.EVMAddress())
	require.NoError(r, err)

	gasLimit := uint64(23000) // in units
	gasPrice, err := evmClient.SuggestGasPrice(r.Ctx)
	require.NoError(r, err)

	data := append(r.BTCZRC20Addr.Bytes(), []byte("hello sailors")...) // this data
	tx := ethtypes.NewTransaction(nonce, r.TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := evmClient.NetworkID(r.Ctx)
	require.NoError(r, err)

	deployerPrivkey, err := r.Account.PrivateKey()
	require.NoError(r, err)

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	require.NoError(r, err)

	err = evmClient.SendTransaction(r.Ctx, signedTx)
	require.NoError(r, err)

	r.Logger.Info("EVM tx sent: %s; to %s, nonce %d", signedTx.Hash().String(), signedTx.To().Hex(), signedTx.Nonce())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)

	r.Logger.Info("EVM tx receipt: %d", receipt.Status)
	r.Logger.Info("  tx hash: %s", receipt.TxHash.String())
	r.Logger.Info("  to: %s", signedTx.To().String())
	r.Logger.Info("  value: %d", signedTx.Value())
	r.Logger.Info("  block num: %d", receipt.BlockNumber)

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.Info("cctx status message: %s", cctx.CctxStatus.StatusMessage)
	revertTxHash := cctx.GetCurrentOutboundParam().Hash
	r.Logger.Info("EVM revert tx receipt: status %d", receipt.Status)

	tx, _, err = r.EVMClient.TransactionByHash(r.Ctx, ethcommon.HexToHash(revertTxHash))
	require.NoError(r, err)

	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(revertTxHash))
	require.NoError(r, err)

	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Reverted)
	utils.RequireTxSuccessful(r, receipt)

	require.Equal(r, r.EVMAddress(), *tx.To(), "expected tx to %s; got %s", r.EVMAddress().Hex(), tx.To().Hex())

	// the received value must be lower than the original value because of the paid fees for the revert tx
	// we check that the value is still greater than 0
	invariant := tx.Value().Cmp(value) != -1 || tx.Value().Cmp(big.NewInt(0)) != 1
	require.False(
		r,
		invariant,
		"expected tx value %s; should be non-null and lower than %s",
		tx.Value().String(),
		value.String(),
	)

	r.Logger.Info("REVERT tx receipt: %d", receipt.Status)
	r.Logger.Info("  tx hash: %s", receipt.TxHash.String())
	r.Logger.Info("  to: %s", tx.To().String())
	r.Logger.Info("  value: %s", tx.Value().String())
	r.Logger.Info("  block num: %d", receipt.BlockNumber)
}
