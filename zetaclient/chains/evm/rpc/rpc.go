package rpc

import (
	"context"
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

const (
	// RPCAlertLatency is the default threshold for RPC latency to be considered unhealthy and trigger an alert.
	// 100s is a reasonable threshold for most EVM chains
	RPCAlertLatency = time.Duration(100) * time.Second
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

// CheckRPCStatus checks the RPC status of the evm chain
func CheckRPCStatus(ctx context.Context, client interfaces.EVMRPCClient) (time.Time, error) {
	// query latest block number
	bn, err := client.BlockNumber(ctx)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "RPC failed on BlockNumber, RPC down?")
	}

	// query suggested gas price
	_, err = client.SuggestGasPrice(ctx)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "RPC failed on SuggestGasPrice, RPC down?")
	}

	// query latest block header
	header, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(bn))
	if err != nil {
		return time.Time{}, errors.Wrap(err, "RPC failed on HeaderByNumber, RPC down?")
	}

	// convert block time to UTC
	// #nosec G115 always in range
	blockTime := time.Unix(int64(header.Time), 0).UTC()

	return blockTime, nil
}
