package zetaclient

import (
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/common"

	"github.com/pkg/errors"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

const (
	// PostVoteInboundGasLimit is the gas limit for voting on observed inbound tx
	PostVoteInboundGasLimit = 400_000

	// PostVoteInboundExecutionGasLimit is the gas limit for voting on observed inbound tx and executing it
	PostVoteInboundExecutionGasLimit = 4_000_000

	// PostVoteInboundMessagePassingExecutionGasLimit is the gas limit for voting on, and executing ,observed inbound tx related to message passing (coin_type == zeta)
	PostVoteInboundMessagePassingExecutionGasLimit = 1_000_000

	// MonitorVoteInboundTxResultInterval is the interval between retries for monitoring tx result in seconds
	MonitorVoteInboundTxResultInterval = 5

	// MonitorVoteInboundTxResultRetryCount is the number of retries to fetch monitoring tx result
	MonitorVoteInboundTxResultRetryCount = 20
)

// GetInBoundVoteMessage returns a new MsgVoteOnObservedInboundTx
func GetInBoundVoteMessage(
	sender string,
	senderChain int64,
	txOrigin string,
	receiver string,
	receiverChain int64,
	amount math.Uint,
	message string,
	inTxHash string,
	inBlockHeight uint64,
	gasLimit uint64,
	coinType common.CoinType,
	asset string,
	signerAddress string,
	eventIndex uint,
) *types.MsgVoteOnObservedInboundTx {
	msg := types.NewMsgVoteOnObservedInboundTx(
		signerAddress,
		sender,
		senderChain,
		txOrigin,
		receiver,
		receiverChain,
		amount,
		message,
		inTxHash,
		inBlockHeight,
		gasLimit,
		coinType,
		asset,
		eventIndex,
	)
	return msg
}

// PostVoteInbound posts a vote on an observed inbound tx
// retryGasLimit is the gas limit used to resend the tx if it fails because of insufficient gas
// it is used when the ballot is finalized and the inbound tx needs to be processed
func (b *ZetaCoreBridge) PostVoteInbound(gasLimit, retryGasLimit uint64, msg *types.MsgVoteOnObservedInboundTx) (string, string, error) {
	authzMsg, authzSigner, err := b.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", "", err
	}

	// don't post send if has already voted before
	ballotIndex := msg.Digest()
	hasVoted, err := b.HasVoted(ballotIndex, msg.Creator)
	if err != nil {
		return "", ballotIndex, errors.Wrapf(err, "PostVoteInbound: unable to check if already voted for ballot %s voter %s", ballotIndex, msg.Creator)
	}
	if hasVoted {
		return "", ballotIndex, nil
	}

	for i := 0; i < DefaultRetryCount; i++ {
		zetaTxHash, err := b.Broadcast(gasLimit, authzMsg, authzSigner)
		if err == nil {
			// monitor the result of the transaction and resend if necessary
			go b.MonitorVoteInboundTxResult(zetaTxHash, retryGasLimit, msg)

			return zetaTxHash, ballotIndex, nil
		}
		b.logger.Debug().Err(err).Msgf("PostVoteInbound broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return "", ballotIndex, fmt.Errorf("post send failed after %d retries", DefaultRetryInterval)
}

// MonitorVoteInboundTxResult monitors the result of a vote inbound tx
// retryGasLimit is the gas limit used to resend the tx if it fails because of insufficient gas
// if retryGasLimit is 0, the tx is not resent
func (b *ZetaCoreBridge) MonitorVoteInboundTxResult(zetaTxHash string, retryGasLimit uint64, msg *types.MsgVoteOnObservedInboundTx) {
	var lastErr error

	for i := 0; i < MonitorVoteInboundTxResultRetryCount; i++ {
		time.Sleep(MonitorVoteInboundTxResultInterval * time.Second)

		// query tx result from ZetaChain
		txResult, err := b.QueryTxResult(zetaTxHash)

		if err == nil {
			if strings.Contains(txResult.RawLog, "failed to execute message") {
				// the inbound vote tx shouldn't fail to execute
				// this shouldn't happen
				b.logger.Error().Msgf(
					"MonitorInboundTxResult: failed to execute vote, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
				)
			} else if strings.Contains(txResult.RawLog, "out of gas") {
				// if the tx fails with an out of gas error, resend the tx with more gas if retryGasLimit > 0
				b.logger.Debug().Msgf(
					"MonitorInboundTxResult: out of gas, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
				)
				if retryGasLimit > 0 {
					// new retryGasLimit set to 0 to prevent reentering this function
					_, _, err := b.PostVoteInbound(retryGasLimit, 0, msg)
					if err != nil {
						b.logger.Error().Err(err).Msgf(
							"MonitorInboundTxResult: failed to resend tx, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
						)
					} else {
						b.logger.Info().Msgf(
							"MonitorInboundTxResult: successfully resent tx, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
						)
					}
				}
			} else {
				b.logger.Debug().Msgf(
					"MonitorInboundTxResult: successful txHash %s, log %s", zetaTxHash, txResult.RawLog,
				)
			}
			return
		}
		lastErr = err
	}

	b.logger.Error().Err(lastErr).Msgf(
		"MonitorInboundTxResult: unable to query tx result for txHash %s, err %s", zetaTxHash, lastErr.Error(),
	)
	return
}
