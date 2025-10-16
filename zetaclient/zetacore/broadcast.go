package zetacore

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"

	"github.com/zeta-chain/node/app/ante"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/zetaclient/authz"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// paying 50% more than the current base gas price to buffer for potential block-by-block
// gas price increase due to EIP1559 feemarket on ZetaChain
var bufferMultiplier = sdkmath.LegacyMustNewDecFromStr("1.5")

var reductionRate = sdkmath.LegacyMustNewDecFromStr(ante.GasPriceReductionRate)

var accountSequenceMismatchRegex = regexp.MustCompile(`account sequence mismatch, expected ([0-9]*), got ([0-9]*)`)

// Broadcast broadcasts tx to ZetaChain. Returns txHash and error
func (c *Client) Broadcast(
	ctx context.Context,
	msgType string,
	msgDigest string,
	gasLimit uint64,
	authzWrappedMsg sdktypes.Msg,
	authzSigner authz.Signer,
) (string, error) {
	blockHeight, err := c.GetBlockHeight(ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get block height")
	}

	baseGasPrice, err := c.GetBaseGasPrice(ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get base gas price")
	}

	// shouldn't happen, but just in case
	if baseGasPrice == 0 {
		baseGasPrice = DefaultBaseGasPrice
	}

	// multiply gas price by the system tx reduction rate
	adjustedBaseGasPrice := sdkmath.LegacyNewDec(baseGasPrice).Mul(reductionRate).Mul(bufferMultiplier)

	c.mu.Lock()
	defer c.mu.Unlock()

	if blockHeight > c.blockHeight {
		c.blockHeight = blockHeight
		accountNumber, seqNumber, err := c.GetAccountNumberAndSequenceNumber(authzSigner.KeyType)
		if err != nil {
			return "", err
		}

		c.accountNumber[authzSigner.KeyType] = accountNumber

		if c.seqNumber[authzSigner.KeyType] < seqNumber {
			c.seqNumber[authzSigner.KeyType] = seqNumber
		}
	}

	factory, err := clienttx.NewFactoryCLI(c.cosmosClientContext, flag.NewFlagSet("zetaclient", 0))
	if err != nil {
		return "", err
	}

	factory = factory.
		WithAccountNumber(c.accountNumber[authzSigner.KeyType]).
		WithSequence(c.seqNumber[authzSigner.KeyType]).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	builder, err := factory.BuildUnsignedTx(authzWrappedMsg)
	if err != nil {
		return "", errors.Wrap(err, "unable to build unsigned tx")
	}

	builder.SetGasLimit(gasLimit)

	// #nosec G115 always in range
	fee := sdktypes.NewCoins(sdktypes.NewCoin(
		config.BaseDenom,
		sdkmath.NewInt(int64(gasLimit)).Mul(adjustedBaseGasPrice.Ceil().RoundInt()),
	))
	builder.SetFeeAmount(fee)

	err = c.SignTx(ctx, factory, c.cosmosClientContext.GetFromName(), builder, true)
	if err != nil {
		return "", errors.Wrap(err, "unable to sign tx")
	}

	txBytes, err := c.cosmosClientContext.TxConfig.TxEncoder()(builder.GetTx())
	if err != nil {
		return "", errors.Wrap(err, "unable to encode tx")
	}

	// broadcast to a Tendermint node
	commit, err := c.cosmosClientContext.BroadcastTxSync(txBytes)
	if err != nil {
		return "", errors.Wrap(err, "fail to broadcast tx sync")
	}

	// Code will be the tendermint ABICode,
	// it starts at 1, so if it is an error, code will not be zero.
	if commit.Code == 0 {
		// increment seqNum
		c.seqNumber[authzSigner.KeyType]++

		return commit.TxHash, nil
	}

	if commit.Code == 32 {
		matches := accountSequenceMismatchRegex.FindStringSubmatch(commit.RawLog)
		if len(matches) != 3 {
			return "", fmt.Errorf("code 32, no matches: %s", commit.RawLog)
		}

		expectedSeq, err := strconv.ParseUint(matches[1], 10, 64)
		if err != nil {
			return "", errors.Wrapf(err, "code 32, cannot parse expected seq %q", matches[1])
		}

		gotSeq, err := strconv.ParseUint(matches[2], 10, 64)
		if err != nil {
			return "", errors.Wrapf(err, "code 32, cannot parse got seq %q", matches[2])
		}

		c.logger.Warn().
			Uint64("from", gotSeq).
			Uint64("to", expectedSeq).
			Msg("reset seq number (from err msg)")

		c.seqNumber[authzSigner.KeyType] = expectedSeq
	}

	return commit.TxHash, fmt.Errorf("failed to broadcast tx (code: %d). Log: %s", commit.Code, commit.RawLog)
}

// BroadcastMultiple broadcasts a tx that contains multiple messages to ZetaChain, returning the txHash and error
func (c *Client) BroadcastBatch(
	ctx context.Context,
	gasLimit uint64,
	authzWrappedMsg sdktypes.Msg,
	authzSigner authz.Signer,
) (string, error) {
	blockHeight, err := c.GetBlockHeight(ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get block height")
	}

	baseGasPrice, err := c.GetBaseGasPrice(ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get base gas price")
	}

	// shouldn't happen, but just in case
	if baseGasPrice == 0 {
		baseGasPrice = DefaultBaseGasPrice
	}

	// multiply gas price by the system tx reduction rate
	adjustedBaseGasPrice := sdkmath.LegacyNewDec(baseGasPrice).Mul(reductionRate).Mul(bufferMultiplier)

	c.mu.Lock()
	defer c.mu.Unlock()

	if blockHeight > c.blockHeight {
		c.blockHeight = blockHeight
		accountNumber, seqNumber, err := c.GetAccountNumberAndSequenceNumber(authzSigner.KeyType)
		if err != nil {
			return "", err
		}

		c.accountNumber[authzSigner.KeyType] = accountNumber

		if c.seqNumber[authzSigner.KeyType] < seqNumber {
			c.seqNumber[authzSigner.KeyType] = seqNumber
		}
	}

	factory, err := clienttx.NewFactoryCLI(c.cosmosClientContext, flag.NewFlagSet("zetaclient", 0))
	if err != nil {
		return "", err
	}

	factory = factory.
		WithAccountNumber(c.accountNumber[authzSigner.KeyType]).
		WithSequence(c.seqNumber[authzSigner.KeyType]).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	builder, err := factory.BuildUnsignedTx(authzWrappedMsg)
	if err != nil {
		return "", errors.Wrap(err, "unable to build unsigned tx")
	}

	builder.SetGasLimit(gasLimit)

	// #nosec G115 always in range
	fee := sdktypes.NewCoins(sdktypes.NewCoin(
		config.BaseDenom,
		sdkmath.NewInt(int64(gasLimit)).Mul(adjustedBaseGasPrice.Ceil().RoundInt()),
	))
	builder.SetFeeAmount(fee)

	err = c.SignTx(ctx, factory, c.cosmosClientContext.GetFromName(), builder, true)
	if err != nil {
		return "", errors.Wrap(err, "unable to sign tx")
	}

	txBytes, err := c.cosmosClientContext.TxConfig.TxEncoder()(builder.GetTx())
	if err != nil {
		return "", errors.Wrap(err, "unable to encode tx")
	}

	// broadcast to a Tendermint node
	commit, err := c.cosmosClientContext.BroadcastTxSync(txBytes)
	if err != nil {
		return "", errors.Wrap(err, "fail to broadcast tx sync")
	}

	// Code will be the tendermint ABICode,
	// it starts at 1, so if it is an error, code will not be zero.
	if commit.Code == 0 {
		// increment seqNum
		c.seqNumber[authzSigner.KeyType]++

		return commit.TxHash, nil
	}

	if commit.Code == 32 {
		matches := accountSequenceMismatchRegex.FindStringSubmatch(commit.RawLog)
		if len(matches) != 3 {
			return "", fmt.Errorf("code 32, no matches: %s", commit.RawLog)
		}

		expectedSeq, err := strconv.ParseUint(matches[1], 10, 64)
		if err != nil {
			return "", errors.Wrapf(err, "code 32, cannot parse expected seq %q", matches[1])
		}

		gotSeq, err := strconv.ParseUint(matches[2], 10, 64)
		if err != nil {
			return "", errors.Wrapf(err, "code 32, cannot parse got seq %q", matches[2])
		}

		c.logger.Warn().
			Uint64("from", gotSeq).
			Uint64("to", expectedSeq).
			Msg("reset seq number (from err msg)")

		c.seqNumber[authzSigner.KeyType] = expectedSeq
	}

	return commit.TxHash, fmt.Errorf("failed to broadcast tx (code: %d). Log: %s", commit.Code, commit.RawLog)
}

// SignTx signs a tx with the given name
func (c *Client) SignTx(
	ctx context.Context,
	txf clienttx.Factory,
	name string,
	txBuilder client.TxBuilder,
	overwriteSig bool,
) error {
	return clienttx.Sign(ctx, txf, name, txBuilder, overwriteSig)
}

// QueryTxResult query the result of a tx
func (c *Client) QueryTxResult(hash string) (*sdktypes.TxResponse, error) {
	return authtx.QueryTx(c.cosmosClientContext, hash)
}

// HandleBroadcastError returns whether to retry in a few seconds, and whether to report via AddOutboundTracker
// returns (bool retry, bool report)
func HandleBroadcastError(err error, nonce uint64, toChain int64, outboundHash string) (bool, bool) {
	if err == nil {
		return false, false
	}

	msg := err.Error()
	evt := log.Warn().
		Err(err).
		Int64(logs.FieldChain, toChain).
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, outboundHash)

	switch {
	// From the literal meaning of the error message, the tx with this 'nonce' has already been processed,
	// and the latest TSS account nonce has already been incremented.
	// Theoretically, this the tx hash should not be posted to the tracker in this case, but we've already
	// encountered missed outbound tracker caused by unknown reasons (may or may not be false positive).
	//
	// To prevent missed potential outbound tracker, now we pass this hash to tracker reporter in this case.
	// The overhead is:
	// 	- It is uncertain whether this tx hash was the very FIRST accepted tx with THIS 'nonce', it might be the second...
	//  - Once decided to report this tx hash, we need to spawn extra goroutines and making extra RPC queries for monitoring.
	case strings.Contains(msg, "nonce too low"):
		const m = "nonce too low! this might be an unnecessary key-sign. increase retry interval and awaits outbound confirmation"
		evt.Msg(m)
		return false, true

	case strings.Contains(msg, "replacement transaction underpriced"):
		evt.Msg("broadcast replacement")
		return false, false

	case strings.Contains(msg, "already known"):
		// report to tracker, because there's possibilities a successful broadcast gets this error code
		evt.Msg("broadcast duplicates")
		return false, true

	default:
		evt.Msg("broadcast error; retrying...")
		return true, false
	}
}
