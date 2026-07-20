// Package migration holds the shared native-fund migration formulas used both by
// the on-chain TSS migration CCTX builder and the off-chain emergency drain tool.
// Keeping the constants and math in one place guarantees every signer computes
// byte-identical amounts.
package migration

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"

	"github.com/zeta-chain/node/pkg/gas"
	btccommon "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

const (
	// GasMultiplierEVM is multiplied to the median gas price to get the gas price for the migration.
	// This is done to avoid the migration tx getting stuck in the mempool.
	GasMultiplierEVM = "2.5"

	// BufferAmountEVM is the buffer amount (wei) added to the fee for the migration transaction.
	BufferAmountEVM = "2100000000"
)

// ErrInsufficientFunds is returned when the balance cannot cover the migration fee.
var ErrInsufficientFunds = errorsmod.Register("migration", 1, "insufficient funds for migration")

// ComputeEVMMigration computes the amount to send, the gas price and gas limit for
// draining an EVM native balance to a destination address.
//
//	gasLimit = 21_000                          (gas.EVMSend)
//	gasPrice = medianGasPrice × 2.5            (GasMultiplierEVM)
//	fee      = gasLimit × gasPrice + 2.1 gwei  (BufferAmountEVM)
//	amount   = balance − fee
func ComputeEVMMigration(
	balance, medianGasPrice sdkmath.Uint,
) (amount sdkmath.Uint, gasPrice sdkmath.Uint, gasLimit uint64, err error) {
	gasLimit = gas.EVMSend
	gasPrice, err = gas.MultiplyGasPrice(medianGasPrice, GasMultiplierEVM)
	if err != nil {
		return sdkmath.ZeroUint(), sdkmath.ZeroUint(), 0, err
	}

	fee := sdkmath.NewUint(gasLimit).
		Mul(gasPrice).
		Add(sdkmath.NewUintFromString(BufferAmountEVM))
	if fee.GT(balance) {
		return sdkmath.ZeroUint(), sdkmath.ZeroUint(), 0, errorsmod.Wrap(
			ErrInsufficientFunds,
			fmt.Sprintf("balance: %s, fee: %s", balance.String(), fee.String()),
		)
	}

	return balance.Sub(fee), gasPrice, gasLimit, nil
}

// ComputeBTCMigration computes the single-output amount and fee (in satoshis) for a
// Bitcoin sweep spending numInputs pinned UTXOs totalling totalInputSats to payee.
//
// The size estimate reuses common.EstimateOutboundSize, which conservatively assumes
// nonce-mark + change outputs; the drain has a single output, so the fee is slightly
// over-estimated, which is safe for a sweep.
func ComputeBTCMigration(
	totalInputSats, feeRate int64,
	numInputs int,
	payee btcutil.Address,
) (outputSats, feeSats int64, err error) {
	// #nosec G115 numInputs is a small bounded count
	txSize, err := btccommon.EstimateOutboundSize(int64(numInputs), []btcutil.Address{payee})
	if err != nil {
		return 0, 0, err
	}

	feeSats = txSize * feeRate
	outputSats = totalInputSats - feeSats
	if outputSats < 0 {
		return 0, 0, errorsmod.Wrap(
			ErrInsufficientFunds,
			fmt.Sprintf("total input: %d, fee: %d", totalInputSats, feeSats),
		)
	}

	return outputSats, feeSats, nil
}
