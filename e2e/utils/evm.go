package utils

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	DefaultReceiptTimeout = 30 * time.Second
)

// MustWaitForTxReceipt waits until a broadcasted tx to be mined and return its receipt
func MustWaitForTxReceipt(
	ctx context.Context,
	client *ethclient.Client,
	tx *ethtypes.Transaction,
	logger infoLogger,
	timeout time.Duration,
) *ethtypes.Receipt {
	if timeout == 0 {
		timeout = DefaultReceiptTimeout
	}

	t := TestingFromContext(ctx)

	start := time.Now()
	for i := 0; ; i++ {
		require.False(t, time.Since(start) > timeout, "waiting tx receipt timeout with timeout %s", timeout.String())

		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			if !errors.Is(err, ethereum.NotFound) && i%10 == 0 {
				logger.Info("fetching tx %s receipt error: %s ", tx.Hash().Hex(), err.Error())
			}
			time.Sleep(1 * time.Second)
			continue
		}
		if receipt != nil {
			return receipt
		}
	}
}
