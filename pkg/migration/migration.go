// Package migration holds the shared native-fund migration formulas used both by
// the on-chain TSS migration CCTX builder and the off-chain emergency drain tool.
// Keeping the constants and math in one place guarantees every signer computes
// byte-identical amounts.
package migration

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/zeta-chain/node/pkg/gas"
	btccommon "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

const (
	// GasMultiplierEVM is multiplied to the median gas price to get the gas price for the migration.
	// This is done to avoid the migration tx getting stuck in the mempool.
	GasMultiplierEVM = "2.5"

	// BufferAmountEVM is the buffer amount (wei) added to the fee for the migration transaction.
	BufferAmountEVM = "2100000000"

	// BTCConservativeFeeRate is a conservative BTC fee rate in sat/vB used as the default for
	// migration calculations. 50 sat/vB is 5x the default testnet rate to ensure the sweep confirms.
	BTCConservativeFeeRate int64 = 50

	// BTCReservedRBFFeeSats reserves 0.01 BTC for potential RBF fee bumping of the sweep.
	BTCReservedRBFFeeSats int64 = 1_000_000

	// BTCNonceMarkBufferSats buffers 3000 satoshis for the nonce-mark output.
	BTCNonceMarkBufferSats int64 = 3_000
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
// Bitcoin sweep spending pinned UTXOs totalling totalInputSats.
//
// It backs the zetatool tss-balances migration-amount display, not the on-chain TSS migration
// path (which defers BTC fee handling to the zetaclient signer). The overhead is intentionally
// generous:
//
//	feeSats = OutboundBytesMax × feeRate    (miner fee at the max outbound tx size)
//	        + BTCReservedRBFFeeSats         (0.01 BTC reserved for RBF fee bumping)
//	        + BTCNonceMarkBufferSats        (3000 sats for the nonce-mark output)
//
// The RBF and nonce-mark reserves are a deliberate margin beyond the real miner fee that
// guarantees a single migration tx confirms; the unspent remainder survives as UTXOs. The
// multi-tx emergency drain uses ComputeBTCSweep instead, where the reserves don't scale.
func ComputeBTCMigration(totalInputSats, feeRate int64) (outputSats, feeSats int64, err error) {
	feeSats = btccommon.OutboundBytesMax*feeRate + BTCReservedRBFFeeSats + BTCNonceMarkBufferSats
	outputSats = totalInputSats - feeSats
	if outputSats < 0 {
		return 0, 0, errorsmod.Wrap(
			ErrInsufficientFunds,
			fmt.Sprintf("total input: %d, fee: %d", totalInputSats, feeSats),
		)
	}

	return outputSats, feeSats, nil
}

// ComputeBTCSweep computes the single-output amount and fee (in satoshis) for one transaction
// of the emergency drain's multi-tx BTC sweep.
//
// Unlike ComputeBTCMigration it charges the miner fee only (no RBF or nonce-mark reserve): the
// drain partitions a large UTXO set into many sweep txs, and a per-tx reserve would both
// multiply across N txs and, in a single-output sweep, be paid to miners rather than preserved.
func ComputeBTCSweep(totalInputSats, feeRate int64) (outputSats, feeSats int64, err error) {
	feeSats = btccommon.OutboundBytesMax * feeRate
	outputSats = totalInputSats - feeSats
	if outputSats < 0 {
		return 0, 0, errorsmod.Wrap(
			ErrInsufficientFunds,
			fmt.Sprintf("total input: %d, fee: %d", totalInputSats, feeSats),
		)
	}

	return outputSats, feeSats, nil
}
