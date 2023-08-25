package zetaclient

import (
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"gitlab.com/thorchain/tss/go-tss/blame"

	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

const (
	PostGasPriceGasLimit            = 1_500_000
	AddTxHashToOutTxTrackerGasLimit = 200_000
	PostNonceGasLimit               = 200_000
	PostSendEVMGasLimit             = 1_000_000 // likely emit a lot of logs, so costly
	PostSendNonEVMGasLimit          = 1_000_000
	PostReceiveConfirmationGasLimit = 200_000
	PostBlameDataGasLimit           = 200_000
	DefaultGasLimit                 = 200_000
	DefaultRetryCount               = 5
	DefaultRetryInterval            = 5
)

func (b *ZetaCoreBridge) WrapMessageWithAuthz(msg sdk.Msg) (sdk.Msg, AuthZSigner) {
	msgURL := sdk.MsgTypeURL(msg)
	authzSigner := GetSigner(msgURL)
	authzMessage := authz.NewMsgExec(authzSigner.GranteeAddress, []sdk.Msg{msg})
	return &authzMessage, authzSigner
}

func (b *ZetaCoreBridge) PostGasPrice(chain common.Chain, gasPrice uint64, supply string, blockNum uint64) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgGasPriceVoter(signerAddress, chain.ChainId, gasPrice, supply, blockNum)
	authzMsg, authzSigner := b.WrapMessageWithAuthz(msg)

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

func (b *ZetaCoreBridge) AddTxHashToOutTxTracker(chainID int64, nonce uint64, txHash string) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgAddToOutTxTracker(signerAddress, chainID, nonce, txHash)
	authzMsg, authzSigner := b.WrapMessageWithAuthz(msg)
	zetaTxHash, err := b.Broadcast(AddTxHashToOutTxTrackerGasLimit, authzMsg, authzSigner)
	if err != nil {
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) PostSend(sender string, senderChain int64, txOrigin string, receiver string, receiverChain int64, amount math.Uint, message string, inTxHash string, inBlockHeight uint64, gasLimit uint64, coinType common.CoinType, zetaGasLimit uint64, asset string) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgVoteOnObservedInboundTx(signerAddress, sender, senderChain, txOrigin, receiver, receiverChain, amount, message, inTxHash, inBlockHeight, gasLimit, coinType, asset)
	authzMsg, authzSigner := b.WrapMessageWithAuthz(msg)

	for i := 0; i < DefaultRetryCount; i++ {
		zetaTxHash, err := b.Broadcast(zetaGasLimit, authzMsg, authzSigner)
		if err == nil {
			return zetaTxHash, nil
		}
		b.logger.Debug().Err(err).Msgf("PostSend broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}

	return "", fmt.Errorf("post send failed after %d retries", DefaultRetryInterval)
}

func (b *ZetaCoreBridge) PostReceiveConfirmation(
	sendHash string,
	outTxHash string,
	outBlockHeight uint64,
	outTxGasUsed uint64,
	amount *big.Int,
	status common.ReceiveStatus,
	chain common.Chain,
	nonce uint64,
	coinType common.CoinType,
) (string, error) {
	lastReport, found := b.lastOutTxReportTime[outTxHash]
	if found && time.Since(lastReport) < 10*time.Minute {
		return "", fmt.Errorf("PostReceiveConfirmation: outTxHash %s already reported in last 10min; last report %s", outTxHash, lastReport)
	}

	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgVoteOnObservedOutboundTx(
		signerAddress,
		sendHash,
		outTxHash,
		outBlockHeight,
		outTxGasUsed,
		math.NewUintFromBigInt(amount),
		status,
		chain.ChainId,
		nonce,
		coinType,
	)
	authzMsg, authzSigner := b.WrapMessageWithAuthz(msg)
	// FIXME: remove this gas limit stuff; in the special ante handler with no gas limit, add
	// NewMsgReceiveConfirmation to it.
	var gasLimit uint64 = PostReceiveConfirmationGasLimit
	if status == common.ReceiveStatus_Failed {
		gasLimit = PostSendEVMGasLimit
	}
	for i := 0; i < DefaultRetryCount; i++ {
		zetaTxHash, err := b.Broadcast(gasLimit, authzMsg, authzSigner)
		if err == nil {
			b.lastOutTxReportTime[outTxHash] = time.Now() // update last report time when bcast succeeds
			return zetaTxHash, nil
		}
		b.logger.Debug().Err(err).Msgf("PostReceive broadcast fail | Retry count : %d", i+1)
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return "", fmt.Errorf("post receive failed after %d retries", DefaultRetryCount)
}

func (b *ZetaCoreBridge) SetTSS(tssPubkey string, keyGenZetaHeight int64, status common.ReceiveStatus) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgCreateTSSVoter(signerAddress, tssPubkey, keyGenZetaHeight, status)
	authzMsg, authzSigner := b.WrapMessageWithAuthz(msg)
	err := error(nil)
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
	b.logger.Info().Msg("UpdateConfig started")
	ticker := time.NewTicker(time.Duration(cfg.ConfigUpdateTicker) * time.Second)
	for {
		select {
		case <-ticker.C:
			b.logger.Debug().Msg("Running Updater")
			err := b.UpdateConfigFromCore(cfg)
			if err != nil {
				b.logger.Err(err).Msg("UpdateConfig error")
				return
			}
		case <-b.stop:
			b.logger.Info().Msg("UpdateConfig stopped")
			return
		}
	}
}

func (b *ZetaCoreBridge) PostBlameData(blame *blame.Blame, chainID int64, index string) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	zetaBlame := &observerTypes.Blame{
		Index:         index,
		FailureReason: blame.FailReason,
		Nodes:         observerTypes.ConvertNodes(blame.BlameNodes),
	}
	msg := observerTypes.NewMsgAddBlameVoteMsg(signerAddress, chainID, zetaBlame)
	authzMsg, authzSigner := b.WrapMessageWithAuthz(msg)
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
