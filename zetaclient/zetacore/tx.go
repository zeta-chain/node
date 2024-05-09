package zetacore

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/go-tss/blame"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	clientauthz "github.com/zeta-chain/zetacore/zetaclient/authz"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
	appcontext "github.com/zeta-chain/zetacore/zetaclient/context"
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
	coinType coin.CoinType,
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

// GasPriceMultiplier returns the gas price multiplier for the given chain
func GasPriceMultiplier(chainID int64) (float64, error) {
	if chains.IsEVMChain(chainID) {
		return clientcommon.EVMOuttxGasPriceMultiplier, nil
	} else if chains.IsBitcoinChain(chainID) {
		return clientcommon.BTCOuttxGasPriceMultiplier, nil
	}
	return 0, fmt.Errorf("cannot get gas price multiplier for unknown chain %d", chainID)
}

func (c *Client) WrapMessageWithAuthz(msg sdk.Msg) (sdk.Msg, clientauthz.Signer, error) {
	msgURL := sdk.MsgTypeURL(msg)

	// verify message validity
	if err := msg.ValidateBasic(); err != nil {
		return nil, clientauthz.Signer{}, fmt.Errorf("%s invalid msg | %s", msgURL, err.Error())
	}

	authzSigner := clientauthz.GetSigner(msgURL)
	authzMessage := authz.NewMsgExec(authzSigner.GranteeAddress, []sdk.Msg{msg})
	return &authzMessage, authzSigner, nil
}

func (c *Client) PostGasPrice(chain chains.Chain, gasPrice uint64, supply string, blockNum uint64) (string, error) {
	// apply gas price multiplier for the chain
	multiplier, err := GasPriceMultiplier(chain.ChainId)
	if err != nil {
		return "", err
	}
	// #nosec G701 always in range
	gasPrice = uint64(float64(gasPrice) * multiplier)
	signerAddress := c.keys.GetOperatorAddress().String()
	msg := types.NewMsgVoteGasPrice(signerAddress, chain.ChainId, gasPrice, supply, blockNum)

	authzMsg, authzSigner, err := c.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	for i := 0; i < DefaultRetryCount; i++ {
		zetaTxHash, err := zetacoreBroadcast(c, PostGasPriceGasLimit, authzMsg, authzSigner)
		if err == nil {
			return zetaTxHash, nil
		}
		c.logger.Debug().Err(err).Msgf("PostGasPrice broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}

	return "", fmt.Errorf("post gasprice failed after %d retries", DefaultRetryInterval)
}

func (c *Client) AddTxHashToOutTxTracker(
	chainID int64,
	nonce uint64,
	txHash string,
	proof *proofs.Proof,
	blockHash string,
	txIndex int64,
) (string, error) {
	// don't report if the tracker already contains the txHash
	tracker, err := c.GetOutTxTracker(chains.Chain{ChainId: chainID}, nonce)
	if err == nil {
		for _, hash := range tracker.HashList {
			if strings.EqualFold(hash.TxHash, txHash) {
				return "", nil
			}
		}
	}
	signerAddress := c.keys.GetOperatorAddress().String()
	msg := types.NewMsgAddToOutTxTracker(signerAddress, chainID, nonce, txHash, proof, blockHash, txIndex)

	authzMsg, authzSigner, err := c.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	zetaTxHash, err := zetacoreBroadcast(c, AddTxHashToOutTxTrackerGasLimit, authzMsg, authzSigner)
	if err != nil {
		return "", err
	}
	return zetaTxHash, nil
}

func (c *Client) SetTSS(tssPubkey string, keyGenZetaHeight int64, status chains.ReceiveStatus) (string, error) {
	signerAddress := c.keys.GetOperatorAddress().String()
	msg := observertypes.NewMsgVoteTSS(signerAddress, tssPubkey, keyGenZetaHeight, status)

	authzMsg, authzSigner, err := c.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	zetaTxHash := ""
	for i := 0; i <= DefaultRetryCount; i++ {
		zetaTxHash, err = zetacoreBroadcast(c, DefaultGasLimit, authzMsg, authzSigner)
		if err == nil {
			return zetaTxHash, nil
		}
		c.logger.Debug().Err(err).Msgf("SetTSS broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}

	return "", fmt.Errorf("set tss failed | err %s", err.Error())
}

// CoreContextUpdater is a polling goroutine that checks and updates core context at every height
func (c *Client) CoreContextUpdater(appContext *appcontext.AppContext) {
	c.logger.Info().Msg("CoreContextUpdater started")
	ticker := time.NewTicker(time.Duration(appContext.Config().ConfigUpdateTicker) * time.Second)
	sampledLogger := c.logger.Sample(&zerolog.BasicSampler{N: 10})
	for {
		select {
		case <-ticker.C:
			c.logger.Debug().Msg("Running Updater")
			err := c.UpdateZetaCoreContext(appContext.ZetaCoreContext(), false, sampledLogger)
			if err != nil {
				c.logger.Err(err).Msg("CoreContextUpdater failed to update config")
			}
		case <-c.stop:
			c.logger.Info().Msg("CoreContextUpdater stopped")
			return
		}
	}
}

func (c *Client) PostBlameData(blame *blame.Blame, chainID int64, index string) (string, error) {
	signerAddress := c.keys.GetOperatorAddress().String()
	zetaBlame := observertypes.Blame{
		Index:         index,
		FailureReason: blame.FailReason,
		Nodes:         observertypes.ConvertNodes(blame.BlameNodes),
	}
	msg := observertypes.NewMsgAddBlameVoteMsg(signerAddress, chainID, zetaBlame)

	authzMsg, authzSigner, err := c.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	var gasLimit uint64 = PostBlameDataGasLimit

	for i := 0; i < DefaultRetryCount; i++ {
		zetaTxHash, err := zetacoreBroadcast(c, gasLimit, authzMsg, authzSigner)
		if err == nil {
			return zetaTxHash, nil
		}
		c.logger.Error().Err(err).Msgf("PostBlame broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return "", fmt.Errorf("post blame data failed after %d retries", DefaultRetryCount)
}

func (c *Client) PostVoteBlockHeader(chainID int64, blockHash []byte, height int64, header proofs.HeaderData) (string, error) {
	signerAddress := c.keys.GetOperatorAddress().String()

	msg := observertypes.NewMsgVoteBlockHeader(signerAddress, chainID, blockHash, height, header)

	authzMsg, authzSigner, err := c.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	var gasLimit uint64 = DefaultGasLimit
	for i := 0; i < DefaultRetryCount; i++ {
		zetaTxHash, err := zetacoreBroadcast(c, gasLimit, authzMsg, authzSigner)
		if err == nil {
			return zetaTxHash, nil
		}
		c.logger.Error().Err(err).Msgf("PostVoteBlockHeader broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return "", fmt.Errorf("post add block header failed after %d retries", DefaultRetryCount)
}

// PostVoteInbound posts a vote on an observed inbound tx
// retryGasLimit is the gas limit used to resend the tx if it fails because of insufficient gas
// it is used when the ballot is finalized and the inbound tx needs to be processed
func (c *Client) PostVoteInbound(gasLimit, retryGasLimit uint64, msg *types.MsgVoteOnObservedInboundTx) (string, string, error) {
	authzMsg, authzSigner, err := c.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", "", err
	}

	// don't post send if has already voted before
	ballotIndex := msg.Digest()
	hasVoted, err := c.HasVoted(ballotIndex, msg.Creator)
	if err != nil {
		return "", ballotIndex, errors.Wrapf(err, "PostVoteInbound: unable to check if already voted for ballot %s voter %s", ballotIndex, msg.Creator)
	}
	if hasVoted {
		return "", ballotIndex, nil
	}

	for i := 0; i < DefaultRetryCount; i++ {
		zetaTxHash, err := zetacoreBroadcast(c, gasLimit, authzMsg, authzSigner)
		if err == nil {
			// monitor the result of the transaction and resend if necessary
			go c.MonitorVoteInboundTxResult(zetaTxHash, retryGasLimit, msg)

			return zetaTxHash, ballotIndex, nil
		}
		c.logger.Debug().Err(err).Msgf("PostVoteInbound broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return "", ballotIndex, fmt.Errorf("post send failed after %d retries", DefaultRetryInterval)
}

// MonitorVoteInboundTxResult monitors the result of a vote inbound tx
// retryGasLimit is the gas limit used to resend the tx if it fails because of insufficient gas
// if retryGasLimit is 0, the tx is not resent
func (c *Client) MonitorVoteInboundTxResult(zetaTxHash string, retryGasLimit uint64, msg *types.MsgVoteOnObservedInboundTx) {
	var lastErr error

	for i := 0; i < MonitorVoteInboundTxResultRetryCount; i++ {
		time.Sleep(MonitorVoteInboundTxResultInterval * time.Second)

		// query tx result from ZetaChain
		txResult, err := c.QueryTxResult(zetaTxHash)

		if err == nil {
			if strings.Contains(txResult.RawLog, "failed to execute message") {
				// the inbound vote tx shouldn't fail to execute
				// this shouldn't happen
				c.logger.Error().Msgf(
					"MonitorInboundTxResult: failed to execute vote, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
				)
			} else if strings.Contains(txResult.RawLog, "out of gas") {
				// if the tx fails with an out of gas error, resend the tx with more gas if retryGasLimit > 0
				c.logger.Debug().Msgf(
					"MonitorInboundTxResult: out of gas, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
				)
				if retryGasLimit > 0 {
					// new retryGasLimit set to 0 to prevent reentering this function
					_, _, err := c.PostVoteInbound(retryGasLimit, 0, msg)
					if err != nil {
						c.logger.Error().Err(err).Msgf(
							"MonitorInboundTxResult: failed to resend tx, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
						)
					} else {
						c.logger.Info().Msgf(
							"MonitorInboundTxResult: successfully resent tx, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
						)
					}
				}
			} else {
				c.logger.Debug().Msgf(
					"MonitorInboundTxResult: successful txHash %s, log %s", zetaTxHash, txResult.RawLog,
				)
			}
			return
		}
		lastErr = err
	}

	c.logger.Error().Err(lastErr).Msgf(
		"MonitorInboundTxResult: unable to query tx result for txHash %s, err %s", zetaTxHash, lastErr.Error(),
	)
}

// PostVoteOutbound posts a vote on an observed outbound tx
func (c *Client) PostVoteOutbound(
	cctxIndex string,
	outTxHash string,
	outBlockHeight uint64,
	outTxGasUsed uint64,
	outTxEffectiveGasPrice *big.Int,
	outTxEffectiveGasLimit uint64,
	amount *big.Int,
	status chains.ReceiveStatus,
	chain chains.Chain,
	nonce uint64,
	coinType coin.CoinType,
) (string, string, error) {
	signerAddress := c.keys.GetOperatorAddress().String()
	msg := types.NewMsgVoteOnObservedOutboundTx(
		signerAddress,
		cctxIndex,
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
	if msg.Status == chains.ReceiveStatus_failed {
		retryGasLimit = PostVoteOutboundRevertGasLimit
	}

	return c.PostVoteOutboundFromMsg(PostVoteOutboundGasLimit, retryGasLimit, msg)
}

// PostVoteOutboundFromMsg posts a vote on an observed outbound tx from a MsgVoteOnObservedOutboundTx
func (c *Client) PostVoteOutboundFromMsg(gasLimit, retryGasLimit uint64, msg *types.MsgVoteOnObservedOutboundTx) (string, string, error) {
	authzMsg, authzSigner, err := c.WrapMessageWithAuthz(msg)
	if err != nil {
		return "", "", err
	}

	// don't post confirmation if has already voted before
	ballotIndex := msg.Digest()
	hasVoted, err := c.HasVoted(ballotIndex, msg.Creator)
	if err != nil {
		return "", ballotIndex, errors.Wrapf(err, "PostVoteOutbound: unable to check if already voted for ballot %s voter %s", ballotIndex, msg.Creator)
	}
	if hasVoted {
		return "", ballotIndex, nil
	}
	for i := 0; i < DefaultRetryCount; i++ {
		zetaTxHash, err := zetacoreBroadcast(c, gasLimit, authzMsg, authzSigner)
		if err == nil {
			// monitor the result of the transaction and resend if necessary
			go c.MonitorVoteOutboundTxResult(zetaTxHash, retryGasLimit, msg)

			return zetaTxHash, ballotIndex, nil
		}
		c.logger.Debug().Err(err).Msgf("PostVoteOutbound broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return "", ballotIndex, fmt.Errorf("post receive failed after %d retries", DefaultRetryCount)
}

// MonitorVoteOutboundTxResult monitors the result of a vote outbound tx
// retryGasLimit is the gas limit used to resend the tx if it fails because of insufficient gas
// if retryGasLimit is 0, the tx is not resent
func (c *Client) MonitorVoteOutboundTxResult(zetaTxHash string, retryGasLimit uint64, msg *types.MsgVoteOnObservedOutboundTx) {
	var lastErr error

	for i := 0; i < MonitorVoteOutboundTxResultRetryCount; i++ {
		time.Sleep(MonitorVoteOutboundTxResultInterval * time.Second)

		// query tx result from ZetaChain
		txResult, err := c.QueryTxResult(zetaTxHash)

		if err == nil {
			if strings.Contains(txResult.RawLog, "failed to execute message") {
				// the inbound vote tx shouldn't fail to execute
				// this shouldn't happen
				c.logger.Error().Msgf(
					"MonitorVoteOutboundTxResult: failed to execute vote, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
				)
			} else if strings.Contains(txResult.RawLog, "out of gas") {
				// if the tx fails with an out of gas error, resend the tx with more gas if retryGasLimit > 0
				c.logger.Debug().Msgf(
					"MonitorVoteOutboundTxResult: out of gas, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
				)
				if retryGasLimit > 0 {
					// new retryGasLimit set to 0 to prevent reentering this function
					_, _, err := c.PostVoteOutboundFromMsg(retryGasLimit, 0, msg)

					if err != nil {
						c.logger.Error().Err(err).Msgf(
							"MonitorVoteOutboundTxResult: failed to resend tx, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
						)
					} else {
						c.logger.Info().Msgf(
							"MonitorVoteOutboundTxResult: successfully resent tx, txHash: %s, log %s", zetaTxHash, txResult.RawLog,
						)
					}
				}
			} else {
				c.logger.Debug().Msgf(
					"MonitorVoteOutboundTxResult: successful txHash %s, log %s", zetaTxHash, txResult.RawLog,
				)
			}
			return
		}
		lastErr = err
	}

	c.logger.Error().Err(lastErr).Msgf(
		"MonitorVoteOutboundTxResult: unable to query tx result for txHash %s, err %s", zetaTxHash, lastErr.Error(),
	)
}
