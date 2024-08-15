package rpc

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
)

// IsTxConfirmed checks if the transaction is confirmed with given confirmations
func IsTxConfirmed(
	ctx context.Context,
	client interfaces.EVMRPCClient,
	txHash string,
	confirmations uint64,
) (bool, error) {
	// query the tx
	_, isPending, err := client.TransactionByHash(ctx, ethcommon.HexToHash(txHash))
	if err != nil {
		return false, errors.Wrapf(err, "error getting transaction for tx %s", txHash)
	}
	if isPending {
		return false, nil
	}

	// query receipt
	receipt, err := client.TransactionReceipt(ctx, ethcommon.HexToHash(txHash))
	if err != nil {
		return false, errors.Wrapf(err, "error getting transaction receipt for tx %s", txHash)
	}

	// should not happen
	if receipt == nil {
		return false, errors.Errorf("receipt is nil for tx %s", txHash)
	}

	// query last block height
	lastHeight, err := client.BlockNumber(ctx)
	if err != nil {
		return false, errors.Wrap(err, "error getting block number")
	}

	// check confirmations
	if lastHeight < receipt.BlockNumber.Uint64() {
		return false, nil
	}
	blocks := lastHeight - receipt.BlockNumber.Uint64() + 1

	return blocks >= confirmations, nil
}
