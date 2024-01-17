package zetaclient

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

const (
	// PostVoteOutboundGasLimit is the gas limit for voting on observed outbound tx
	PostVoteOutboundGasLimit = 400_000

	// PostVoteOutboundRevertGasLimit is the gas limit for voting on observed outbound tx for revert (when outbound fails)
	// The value needs to be higher because reverting implies interacting with the EVM to perform swaps for the gas token
	PostVoteOutboundRevertGasLimit = 1_500_000

	// MonitorVoteOutboundTxResultInterval is the interval between retries for monitoring tx result in seconds
	MonitorVoteOutboundTxResultInterval = 5

	// MonitorVoteOutboundTxResultRetryCount is the number of retries to fetch monitoring tx result
	MonitorVoteOutboundTxResultRetryCount = 20
)

// PostVoteOutbound posts a vote on an observed outbound tx
func (b *ZetaCoreBridge) PostVoteOutbound(
	sendHash string,
	outTxHash string,
	outBlockHeight uint64,
	outTxGasUsed uint64,
	outTxEffectiveGasPrice *big.Int,
	outTxEffectiveGasLimit uint64,
	amount *big.Int,
	status common.ReceiveStatus,
	chain common.Chain,
	nonce uint64,
	coinType common.CoinType,
) (string, string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgVoteOnObservedOutboundTx(
		signerAddress,
		sendHash,
		outTxHash,
		outBlockHeight,
		outTxGasUsed,
		math.NewIntFromBigInt(outTxEffectiveGasPrice),
		outTxEffectiveGasLimit,
		math.NewUintFromBigInt(amount),
		status,
		chain.ChainId,
		nonce,
		coinType,
	)

	// when an outbound fails and a revert is required, the gas limit needs to be higher
	// this is because the revert tx needs to interact with the EVM to perform swaps for the gas token
	// the higher gas limit is only necessary when the vote is finalized and the outbound is processed
	// therefore we use a retryGasLimit with a higher value to resend the tx if it fails (when the vote is finalized)
	retryGasLimit := uint64(0)
	if msg.Status == common.ReceiveStatus_Failed {
		retryGasLimit = PostVoteOutboundRevertGasLimit
	}

	return b.PostVoteOutboundFromMsg(PostVoteOutboundGasLimit, retryGasLimit, msg)
}

// PostVoteOutboundFromMsg posts a vote on an observed outbound tx from a MsgVoteOnObservedOutboundTx
func (b *ZetaCoreBridge) PostVoteOutboundFromMsg(gasLimit, retryGasLimit uint64, msg *types.MsgVoteOnObservedOutboundTx) (string, string, error) {
	authzMsg, authzSigner, err := b.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", "", err
	}

	// don't post confirmation if has already voted before
	ballotIndex := msg.Digest()
	hasVoted, err := b.HasVoted(ballotIndex, msg.Creator)
	if err != nil {
		return "", ballotIndex, errors.Wrapf(err, "PostVoteOutbound: unable to check if already voted for ballot %s voter %s", ballotIndex, msg.Creator)
	}
	if hasVoted {
		return "", ballotIndex, nil
	}
	for i := 0; i < DefaultRetryCount; i++ {
		zetaTxHash, err := b.Broadcast(gasLimit, authzMsg, authzSigner)
		if err == nil {
			// monitor the result of the transaction and resend if necessary
			go b.MonitorVoteOutboundTxResult(zetaTxHash, retryGasLimit, msg)

			return zetaTxHash, ballotIndex, nil
		}
		b.logger.Debug().Err(err).Msgf("PostVoteOutbound broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return "", ballotIndex, fmt.Errorf("post receive failed after %d retries", DefaultRetryCount)
}

// MonitorVoteOutboundTxResult monitors the result of a vote outbound tx
// retryGasLimit is the gas limit used to resend the tx if it fails because of insufficient gas
// if retryGasLimit is 0, the tx is not resent
func (b *ZetaCoreBridge) MonitorVoteOutboundTxResult(zetaTxHash string, retryGasLimit uint64, msg *types.MsgVoteOnObservedOutboundTx) {
	var lastErr error

	for i := 0; i < MonitorVoteOutboundTxResultRetryCount; i++ {
		time.Sleep(MonitorVoteOutboundTxResultInterval * time.Second)

		// query tx result from ZetaChain
		txResult, err := b.QueryTxResult(zetaTxHash)

		if err == nil {
			if strings.Contains(txResult.RawLog, "failed to execute message") {
				// the inbound vote tx shouldn't fail to execute
				// this shouldn't happen
				b.logger.Error().Msgf(
					"MonitorVoteOutboundTxResult: failed to execute vote, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
				)
			} else if strings.Contains(txResult.RawLog, "out of gas") {
				// if the tx fails with an out of gas error, resend the tx with more gas if retryGasLimit > 0
				b.logger.Debug().Msgf(
					"MonitorVoteOutboundTxResult: out of gas, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
				)
				if retryGasLimit > 0 {
					// new retryGasLimit set to 0 to prevent reentering this function
					_, _, err := b.PostVoteOutboundFromMsg(retryGasLimit, 0, msg)

					if err != nil {
						b.logger.Error().Err(err).Msgf(
							"MonitorVoteOutboundTxResult: failed to resend tx, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
						)
					} else {
						b.logger.Info().Msgf(
							"MonitorVoteOutboundTxResult: successfully resent tx, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
						)
					}
				}
			} else {
				b.logger.Debug().Msgf(
					"MonitorVoteOutboundTxResult: successful txHash %s, log %s", zetaTxHash, txResult.RawLog,
				)
			}
			return
		}
		lastErr = err
	}

	b.logger.Error().Err(lastErr).Msgf(
		"MonitorVoteOutboundTxResult: unable to query tx result for txHash %s, err %s", zetaTxHash, lastErr.Error(),
	)
	return
}
