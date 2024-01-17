package zetaclient

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/pkg/errors"

	"github.com/zeta-chain/go-tss/blame"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

const (
	// DefaultGasLimit is the default gas limit used for broadcasting txs
	DefaultGasLimit = 200_000

	// PostGasPriceGasLimit is the gas limit for voting new gas price
	PostGasPriceGasLimit = 1_500_000

	// AddTxHashToOutTxTrackerGasLimit is the gas limit for adding tx hash to out tx tracker
	AddTxHashToOutTxTrackerGasLimit = 200_000

	// PostVoteInboundGasLimit is the gas limit for voting on observed inbound tx
	PostVoteInboundGasLimit = 400_000

	// PostVoteInboundExecutionGasLimit is the gas limit for voting on observed inbound tx and executing it
	PostVoteInboundExecutionGasLimit = 4_000_000

	// PostVoteInboundMessagePassingExecutionGasLimit is the gas limit for voting on, and executing ,observed inbound tx related to message passing (coin_type == zeta)
	PostVoteInboundMessagePassingExecutionGasLimit = 1_000_000

	// PostVoteOutboundGasLimit is the gas limit for voting on observed outbound tx
	PostVoteOutboundGasLimit = 400_000

	// PostBlameDataGasLimit is the gas limit for voting on blames
	PostBlameDataGasLimit = 200_000

	// DefaultRetryCount is the number of retries for broadcasting a tx
	DefaultRetryCount = 5

	// ExtendedRetryCount is an extended number of retries for broadcasting a tx, used in keygen operations
	ExtendedRetryCount = 15

	// DefaultRetryInterval is the interval between retries in seconds
	DefaultRetryInterval = 5

	// MonitorTxResultInterval is the interval between retries for monitoring tx result in seconds
	MonitorTxResultInterval = 5

	// MonitorTxResultRetryCount is the number of retries to fetch monitoring tx result
	MonitorTxResultRetryCount = 10
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

func (b *ZetaCoreBridge) WrapMessageWithAuthz(msg sdk.Msg) (sdk.Msg, AuthZSigner, error) {
	msgURL := sdk.MsgTypeURL(msg)

	// verify message validity
	if err := msg.ValidateBasic(); err != nil {
		return nil, AuthZSigner{}, fmt.Errorf("%s invalid msg | %s", msgURL, err.Error())
	}

	authzSigner := GetSigner(msgURL)
	authzMessage := authz.NewMsgExec(authzSigner.GranteeAddress, []sdk.Msg{msg})
	return &authzMessage, authzSigner, nil
}

func (b *ZetaCoreBridge) PostGasPrice(chain common.Chain, gasPrice uint64, supply string, blockNum uint64) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgGasPriceVoter(signerAddress, chain.ChainId, gasPrice, supply, blockNum)

	authzMsg, authzSigner, err := b.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	for i := 0; i < DefaultRetryCount; i++ {
		zetaTxHash, err := b.Broadcast(PostGasPriceGasLimit, authzMsg, authzSigner)
		if err == nil {
			return zetaTxHash, nil
		}
		b.logger.Debug().Err(err).Msgf("PostGasPrice broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}

	return "", fmt.Errorf("post gasprice failed after %d retries", DefaultRetryInterval)
}

func (b *ZetaCoreBridge) AddTxHashToOutTxTracker(
	chainID int64,
	nonce uint64,
	txHash string,
	proof *common.Proof,
	blockHash string,
	txIndex int64,
) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgAddToOutTxTracker(signerAddress, chainID, nonce, txHash, proof, blockHash, txIndex)

	authzMsg, authzSigner, err := b.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	zetaTxHash, err := b.Broadcast(AddTxHashToOutTxTrackerGasLimit, authzMsg, authzSigner)
	if err != nil {
		return "", err
	}
	return zetaTxHash, nil
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
			// monitor the result of the transaction
			go b.MonitorVoteInboundTxResult(zetaTxHash, retryGasLimit, msg)

			return zetaTxHash, ballotIndex, nil
		}
		b.logger.Debug().Err(err).Msgf("PostVoteInbound broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return "", ballotIndex, fmt.Errorf("post send failed after %d retries", DefaultRetryInterval)
}

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

	// FIXME: remove this gas limit stuff; in the special ante handler with no gas limit, add
	// NewMsgReceiveConfirmation to it.
	var gasLimit uint64 = PostVoteOutboundGasLimit
	if status == common.ReceiveStatus_Failed {
		gasLimit = PostVoteInboundExecutionGasLimit
	}
	for i := 0; i < DefaultRetryCount; i++ {
		zetaTxHash, err := b.Broadcast(gasLimit, authzMsg, authzSigner)
		if err == nil {
			//// monitor the result of the transaction
			//go b.MonitorTxResult(zetaTxHash, true)
			return zetaTxHash, ballotIndex, nil
		}
		b.logger.Debug().Err(err).Msgf("PostVoteOutbound broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return "", ballotIndex, fmt.Errorf("post receive failed after %d retries", DefaultRetryCount)
}

func (b *ZetaCoreBridge) SetTSS(tssPubkey string, keyGenZetaHeight int64, status common.ReceiveStatus) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgCreateTSSVoter(signerAddress, tssPubkey, keyGenZetaHeight, status)

	authzMsg, authzSigner, err := b.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	zetaTxHash := ""
	for i := 0; i <= DefaultRetryCount; i++ {
		zetaTxHash, err = b.Broadcast(DefaultGasLimit, authzMsg, authzSigner)
		if err == nil {
			return zetaTxHash, nil
		}
		b.logger.Debug().Err(err).Msgf("SetTSS broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}

	return "", fmt.Errorf("set tss failed | err %s", err.Error())
}

func (b *ZetaCoreBridge) ConfigUpdater(cfg *config.Config) {
	b.logger.Info().Msg("ConfigUpdater started")
	ticker := time.NewTicker(time.Duration(cfg.ConfigUpdateTicker) * time.Second)
	for {
		select {
		case <-ticker.C:
			b.logger.Debug().Msg("Running Updater")
			err := b.UpdateConfigFromCore(cfg, false)
			if err != nil {
				b.logger.Err(err).Msg("ConfigUpdater failed to update config")
			}
		case <-b.stop:
			b.logger.Info().Msg("ConfigUpdater stopped")
			return
		}
	}
}

func (b *ZetaCoreBridge) PostBlameData(blame *blame.Blame, chainID int64, index string) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	zetaBlame := observerTypes.Blame{
		Index:         index,
		FailureReason: blame.FailReason,
		Nodes:         observerTypes.ConvertNodes(blame.BlameNodes),
	}
	msg := observerTypes.NewMsgAddBlameVoteMsg(signerAddress, chainID, zetaBlame)

	authzMsg, authzSigner, err := b.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	var gasLimit uint64 = PostBlameDataGasLimit

	for i := 0; i < DefaultRetryCount; i++ {
		zetaTxHash, err := b.Broadcast(gasLimit, authzMsg, authzSigner)
		if err == nil {
			return zetaTxHash, nil
		}
		b.logger.Error().Err(err).Msgf("PostBlame broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return "", fmt.Errorf("post blame data failed after %d retries", DefaultRetryCount)
}

func (b *ZetaCoreBridge) PostAddBlockHeader(chainID int64, blockHash []byte, height int64, header common.HeaderData) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()

	msg := observerTypes.NewMsgAddBlockHeader(signerAddress, chainID, blockHash, height, header)

	authzMsg, authzSigner, err := b.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	var gasLimit uint64 = DefaultGasLimit
	for i := 0; i < DefaultRetryCount; i++ {
		zetaTxHash, err := b.Broadcast(gasLimit, authzMsg, authzSigner)
		if err == nil {
			return zetaTxHash, nil
		}
		b.logger.Error().Err(err).Msgf("PostAddBlockHeader broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return "", fmt.Errorf("post add block header failed after %d retries", DefaultRetryCount)
}

// MonitorVoteInboundTxResult monitors the result of a tx (used for inbound and outbound vote txs)
// retryGasLimit is the gas limit used to resend the tx if it fails because of insufficient gas
// if retryGasLimit is 0, the tx is not resent
func (b *ZetaCoreBridge) MonitorVoteInboundTxResult(zetaTxHash string, retryGasLimit uint64, msg *types.MsgVoteOnObservedInboundTx) {
	var lastErr error

	for i := 0; i < MonitorTxResultRetryCount; i++ {
		time.Sleep(MonitorTxResultInterval * time.Second)

		// query tx result from ZetaChain
		txResult, err := b.QueryTxResult(zetaTxHash)

		if err == nil {
			if strings.Contains(txResult.RawLog, "failed to execute message") {
				// the inbound vote tx shouldn't fail to execute
				// this shouldn't happen
				b.logger.Error().Msgf(
					"MonitorInboundTxResult: failed to execute vote, txHash: %s, log %s",
					zetaTxHash,
					txResult.RawLog,
				)
			} else if strings.Contains(txResult.RawLog, "out of gas") {
				// if the tx fails with out of gas error, resend the tx with more gas if retryGasLimit > 0
				b.logger.Debug().Msgf(
					"MonitorInboundTxResult: out of gas, txHash: %s, log %s",
					zetaTxHash,
					txResult.RawLog,
				)
				if retryGasLimit > 0 {
					_, _, err := b.PostVoteInbound(retryGasLimit, 0, msg)
					if err != nil {
						b.logger.Error().Err(err).Msgf(
							"MonitorInboundTxResult: failed to resend tx, txHash: %s, log %s",
							zetaTxHash,
							txResult.RawLog,
						)
					} else {
						b.logger.Info().Msgf(
							"MonitorInboundTxResult: successfully resent tx, txHash: %s, log %s",
							zetaTxHash,
							txResult.RawLog,
						)
					}
				}
			} else {
				b.logger.Debug().Msgf(
					"MonitorInboundTxResult: successful txHash %s, log %s",
					zetaTxHash,
					txResult.RawLog,
				)
			}
			return
		}
		lastErr = err
	}

	b.logger.Error().Err(lastErr).Msgf(
		"MonitorInboundTxResult: unable to query tx result for txHash %s, err %s",
		zetaTxHash,
		lastErr.Error(),
	)
	return
}
