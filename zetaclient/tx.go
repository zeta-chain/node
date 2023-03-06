package zetaclient

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"math/big"
	"time"

	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

const (
	PostGasPriceGasLimit            = 1_500_000
	AddTxHashToOutTxTrackerGasLimit = 200_000
	PostNonceGasLimit               = 200_000
	PostSendEVMGasLimit             = 1_000_000 // likely emit a lot of logs, so costly
	PostSendNonEVMGasLimit          = 1_000_000
	PostReceiveConfirmationGasLimit = 200_000
	DefaultGasLimit                 = 200_000
)

func (b *ZetaCoreBridge) WrapMessageWithAuthz(msg sdk.Msg) sdk.Msg {
	msgUrl := sdk.MsgTypeURL(msg)
	authzSigner := GetSigner(msgUrl)
	authzMessage := authz.NewMsgExec(authzSigner.GranteeAddress, []sdk.Msg{msg})
	return &authzMessage
}

func (b *ZetaCoreBridge) PostGasPrice(chain common.Chain, gasPrice uint64, supply string, blockNum uint64) (string, error) {

	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgGasPriceVoter(signerAddress, chain.ChainId, gasPrice, supply, blockNum)
	authzMsg := b.WrapMessageWithAuthz(msg)

	zetaTxHash, err := b.Broadcast(PostGasPriceGasLimit, authzMsg)
	if err != nil {
		b.logger.Error().Err(err).Msg("PostGasPrice broadcast fail")
		return "", err
	}
	b.logger.Debug().Str("zetaTxHash", zetaTxHash).Msg("PostGasPrice broadcast success")

	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) AddTxHashToOutTxTracker(chainID int64, nonce uint64, txHash string) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgAddToOutTxTracker(signerAddress, chainID, nonce, txHash)
	zetaTxHash, err := b.Broadcast(AddTxHashToOutTxTrackerGasLimit, msg)
	if err != nil {
		b.logger.Error().Err(err).Msg("AddTxHashToOutTxTracker broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) PostNonce(chain common.Chain, nonce uint64) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgNonceVoter(signerAddress, chain.ChainId, nonce)
	zetaTxHash, err := b.Broadcast(PostNonceGasLimit, msg)
	if err != nil {
		b.logger.Error().Err(err).Msg("PostNonce broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) PostSend(sender string, senderChain int64, txOrigin string, receiver string, receiverChain int64, amount math.Uint, message string, inTxHash string, inBlockHeight uint64, gasLimit uint64, coinType common.CoinType, zetaGasLimit uint64, asset string) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgSendVoter(signerAddress, sender, senderChain, txOrigin, receiver, receiverChain, amount, message, inTxHash, inBlockHeight, gasLimit, coinType, asset)
	var zetaTxHash string
	for i := 0; i < 2; i++ {
		zetaTxHash, err := b.Broadcast(zetaGasLimit, msg)
		if err != nil {
			b.logger.Error().Err(err).Msg("PostSend broadcast fail; re-trying...")
		} else {
			return zetaTxHash, nil
		}
		time.Sleep(1 * time.Second)
	}
	return zetaTxHash, fmt.Errorf("postSend: re-try fails")
}

func (b *ZetaCoreBridge) PostReceiveConfirmation(sendHash string, outTxHash string, outBlockHeight uint64, amount *big.Int, status common.ReceiveStatus, chain common.Chain, nonce int, coinType common.CoinType) (string, error) {
	lastReport, found := b.lastOutTxReportTime[outTxHash]
	if found && time.Since(lastReport) < 10*time.Minute {
		return "", fmt.Errorf("PostReceiveConfirmation: outTxHash %s already reported in last 10min; last report %s", outTxHash, lastReport)
	}

	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgReceiveConfirmation(signerAddress, sendHash, outTxHash, outBlockHeight, math.NewUintFromBigInt(amount), status, chain.ChainId, uint64(nonce), coinType)
	//b.logger.Info().Msgf("PostReceiveConfirmation msg digest: %s", msg.Digest())
	var zetaTxHash string
	// FIXME: remove this gas limit stuff; in the special ante handler with no gas limit, add
	// NewMsgReceiveConfirmation to it.
	var gasLimit uint64 = PostReceiveConfirmationGasLimit
	if status == common.ReceiveStatus_Failed {
		gasLimit = PostSendEVMGasLimit
	}
	for i := 0; i < 2; i++ {
		zetaTxHash, err := b.Broadcast(gasLimit, msg)
		if err != nil {
			b.logger.Error().Err(err).Msg("PostReceiveConfirmation broadcast fail; re-trying...")
		} else {
			b.lastOutTxReportTime[outTxHash] = time.Now() // update last report time when bcast succeeds
			return zetaTxHash, nil
		}
		time.Sleep(1 * time.Second)
	}
	return zetaTxHash, fmt.Errorf("postReceiveConfirmation: re-try fails")
}

func (b *ZetaCoreBridge) SetNodeKey(pubkeyset common.PubKeySet, conskey string) (string, error) {
	address, err := b.keys.GetSignerInfo(common.ObserverGranteeKey).GetAddress()
	if err != nil {
		return "", err
	}
	signerAddress := address.String()
	msg := types.NewMsgSetNodeKeys(signerAddress, pubkeyset, conskey)
	zetaTxHash, err := b.Broadcast(DefaultGasLimit, msg)
	if err != nil {
		return "", err
	}
	b.logger.Debug().Msgf("SetNodeKey txhash: %s", zetaTxHash)

	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) SetTSS(chain common.Chain, address string, pubkey string) (string, error) {
	signerAddress := b.keys.GetOperatorAddress().String()
	msg := types.NewMsgCreateTSSVoter(signerAddress, chain.ChainName.String(), address, pubkey)
	zetaTxHash, err := b.Broadcast(DefaultGasLimit, msg)
	if err != nil {
		b.logger.Err(err).Msg("SetNodeKey broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}
