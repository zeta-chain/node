package zetaclient

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"time"
)

func (b *MetachainBridge) PostGasBalance(chain common.Chain, gasBalance string, blockNum uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgGasBalanceVoter(signerAddress, chain.String(), gasBalance, blockNum)
	metaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("PostGasPrice broadcast fail")
		return "", err
	}
	return metaTxHash, nil
}

func (b *MetachainBridge) PostGasPrice(chain common.Chain, gasPrice uint64, supply string, blockNum uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgGasPriceVoter(signerAddress, chain.String(), gasPrice, supply, blockNum)
	metaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("PostGasPrice broadcast fail")
		return "", err
	}
	return metaTxHash, nil
}
func (b *MetachainBridge) PostNonce(chain common.Chain, nonce uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgNonceVoter(signerAddress, chain.String(), nonce)
	metaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("PostNonce broadcast fail")
		return "", err
	}
	return metaTxHash, nil
}
func (b *MetachainBridge) PostSend(sender string, senderChain string, receiver string, receiverChain string, mBurnt string, mMint string, message string, inTxHash string, inBlockHeight uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgSendVoter(signerAddress, sender, senderChain, receiver, receiverChain, mBurnt, mMint, message, inTxHash, inBlockHeight)
	var metaTxHash string
	for i := 0; i < 2; i++ {
		metaTxHash, err := b.Broadcast(msg)
		if err != nil {
			log.Err(err).Msg("PostSend broadcast fail; re-trying...")
		} else {
			return metaTxHash, nil
		}
		time.Sleep(1 * time.Second)
	}
	return metaTxHash, fmt.Errorf("PostSend: re-try fails!")
}

func (b *MetachainBridge) PostReceiveConfirmation(sendHash string, outTxHash string, outBlockHeight uint64, mMint string, status common.ReceiveStatus, chain string) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgReceiveConfirmation(signerAddress, sendHash, outTxHash, outBlockHeight, mMint, status, chain)
	log.Info().Msgf("PostReceiveConfirmation msg digest: %s", msg.Digest())
	var metaTxHash string
	for i := 0; i < 2; i++ {
		metaTxHash, err := b.Broadcast(msg)
		if err != nil {
			log.Err(err).Msg("PostReceiveConfirmation broadcast fail; re-trying...")
		} else {
			return metaTxHash, nil
		}
		time.Sleep(1 * time.Second)
	}
	return metaTxHash, fmt.Errorf("PostReceiveConfirmation: re-try fails!")
}

func (b *MetachainBridge) GetAllSend() ([]*types.Send, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.SendAll(context.Background(), &types.QueryAllSendRequest{})
	if err != nil {
		log.Error().Err(err).Msg("query SendAll error")
		return nil, err
	}
	return resp.Send, nil
}

func (b *MetachainBridge) GetAllPendingSend() ([]*types.Send, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.SendAllPending(context.Background(), &types.QueryAllSendPendingRequest{})
	if err != nil {
		log.Error().Err(err).Msg("query SendAllPending error")
		return nil, err
	}
	return resp.Send, nil
}

func (b *MetachainBridge) GetAllReceive() ([]*types.Receive, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.ReceiveAll(context.Background(), &types.QueryAllReceiveRequest{})
	if err != nil {
		log.Error().Err(err).Msg("query GetAllReceive error")
		return nil, err
	}
	return resp.Receive, nil
}

func (b *MetachainBridge) GetLastBlockHeight() ([]*types.LastBlockHeight, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastBlockHeightAll(context.Background(), &types.QueryAllLastBlockHeightRequest{})
	if err != nil {
		log.Warn().Err(err).Msg("query GetLastBlockHeight error")
		return nil, err
	}
	return resp.LastBlockHeight, nil
}

func (b *MetachainBridge) GetMetaBlockHeight() (uint64, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastMetaHeight(context.Background(), &types.QueryLastMetaHeightRequest{})
	if err != nil {
		log.Warn().Err(err).Msg("query GetLastBlockHeight error")
		return 0, err
	}
	return resp.Height, nil
}

func (b *MetachainBridge) GetLastBlockHeightByChain(chain common.Chain) (*types.LastBlockHeight, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastBlockHeight(context.Background(), &types.QueryGetLastBlockHeightRequest{Index: chain.String()})
	if err != nil {
		log.Error().Err(err).Msg("query GetLastBlockHeight error")
		return nil, err
	}
	return resp.LastBlockHeight, nil
}

func (b *MetachainBridge) GetNonceByChain(chain common.Chain) (*types.ChainNonces, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.ChainNonces(context.Background(), &types.QueryGetChainNoncesRequest{Index: chain.String()})
	if err != nil {
		log.Error().Err(err).Msg("query GetNonceByChain error")
		return nil, err
	}
	return resp.ChainNonces, nil
}
