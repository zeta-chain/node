package metaclient

import (
	"context"
	"github.com/Meta-Protocol/metacore/common"
	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/rs/zerolog/log"
)

func (b *MetachainBridge) PostNonce(chain common.Chain, nonce uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgNonceVoter(signerAddress, chain.String(), nonce)
	metaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("PostSend broadcast fail")
		return "", err
	}
	return metaTxHash, nil
}
func (b *MetachainBridge) PostSend(sender string, senderChain string, receiver string, receiverChain string, mBurnt string, mMint string, message string, inTxHash string, inBlockHeight uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgSendVoter(signerAddress, sender, senderChain, receiver, receiverChain, mBurnt, mMint, message, inTxHash, inBlockHeight)

	metaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("PostSend broadcast fail")
		return "", err
	}
	return metaTxHash, nil
}

func (b *MetachainBridge) PostReceiveConfirmation(sendHash string, outTxHash string, outBlockHeight uint64, mMint string) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgReceiveConfirmation(signerAddress, sendHash, outTxHash, outBlockHeight, mMint)
	metaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("PostReceiveConfirmation broadcast fail")
		return "", err
	}
	return metaTxHash, nil
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
