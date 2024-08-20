package rpc

import (
	"context"
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
)

const (
	// rpcLatencyThreshold is the threshold for RPC latency to be considered unhealthy
	// 100s is a reasonable threshold for most EVM chains
	rpcLatencyThreshold = 100
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
func CheckRPCStatus(ctx context.Context, client interfaces.EVMRPCClient, logger zerolog.Logger) error {
	// query latest block number
	bn, err := client.BlockNumber(ctx)
	if err != nil {
		return errors.Wrap(err, "BlockNumber error: RPC down?")
	}

	// query suggested gas price
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return errors.Wrap(err, "SuggestGasPrice error: RPC down?")
	}

	// query latest block header
	header, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(bn))
	if err != nil {
		return errors.Wrap(err, "HeaderByNumber error: RPC down?")
	}

	// latest block should not be too old
	// #nosec G115 always in range
	blockTime := time.Unix(int64(header.Time), 0).UTC()
	elapsedSeconds := time.Since(blockTime).Seconds()
	if elapsedSeconds > rpcLatencyThreshold {
		return errors.Errorf(
			"Latest block %d is %.0fs old, RPC stale or chain stuck (check explorer)?",
			bn,
			elapsedSeconds,
		)
	}

	logger.Info().
		Msgf("RPC Status [OK]: latest block %d, timestamp %s (%.0fs ago), gas price %s", header.Number, blockTime.String(), elapsedSeconds, gasPrice.String())
	return nil
}
