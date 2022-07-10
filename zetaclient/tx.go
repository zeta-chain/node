package zetaclient

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"time"
)

func (b *ZetaCoreBridge) PostZetaConversionRate(chain common.Chain, rate string, blockNum uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgZetaConversionRateVoter(signerAddress, chain.String(), rate, blockNum)
	zetaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("PostZetaConversionRate broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) PostGasBalance(chain common.Chain, gasBalance string, blockNum uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgGasBalanceVoter(signerAddress, chain.String(), gasBalance, blockNum)
	zetaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("PostGasPrice broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) PostGasPrice(chain common.Chain, gasPrice uint64, supply string, blockNum uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgGasPriceVoter(signerAddress, chain.String(), gasPrice, supply, blockNum)
	zetaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("PostGasPrice broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) PostNonce(chain common.Chain, nonce uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgNonceVoter(signerAddress, chain.String(), nonce)
	zetaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("PostNonce broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}
func (b *ZetaCoreBridge) PostSend(sender string, senderChain string, receiver string, receiverChain string, mBurnt string, mMint string, message string, inTxHash string, inBlockHeight uint64, gasLimit uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgSendVoter(signerAddress, sender, senderChain, receiver, receiverChain, mBurnt, mMint, message, inTxHash, inBlockHeight, gasLimit)
	var zetaTxHash string
	for i := 0; i < 2; i++ {
		zetaTxHash, err := b.Broadcast(msg)
		if err != nil {
			log.Err(err).Msg("PostSend broadcast fail; re-trying...")
		} else {
			return zetaTxHash, nil
		}
		time.Sleep(1 * time.Second)
	}
	return zetaTxHash, fmt.Errorf("PostSend: re-try fails!")
}

// FIXME: pass nonce
func (b *ZetaCoreBridge) PostReceiveConfirmation(sendHash string, outTxHash string, outBlockHeight uint64, mMint string, status common.ReceiveStatus, chain string) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgReceiveConfirmation(signerAddress, sendHash, outTxHash, outBlockHeight, mMint, status, chain)
	log.Info().Msgf("PostReceiveConfirmation msg digest: %s", msg.Digest())
	var zetaTxHash string
	for i := 0; i < 2; i++ {
		zetaTxHash, err := b.Broadcast(msg)
		if err != nil {
			log.Err(err).Msg("PostReceiveConfirmation broadcast fail; re-trying...")
		} else {
			return zetaTxHash, nil
		}
		time.Sleep(1 * time.Second)
	}
	return zetaTxHash, fmt.Errorf("PostReceiveConfirmation: re-try fails!")
}

func (b *ZetaCoreBridge) GetAllSend() ([]*types.Send, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.SendAll(context.Background(), &types.QueryAllSendRequest{})
	if err != nil {
		log.Error().Err(err).Msg("query SendAll error")
		return nil, err
	}
	return resp.Send, nil
}

func (b *ZetaCoreBridge) GetSendByHash(sendHash string) (*types.Send, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.Send(context.Background(), &types.QueryGetSendRequest{Index: sendHash})
	if err != nil {
		log.Error().Err(err).Msg("GetSendByHash error")
		return nil, err
	}
	return resp.Send, nil
}

func (b *ZetaCoreBridge) GetAllPendingSend() ([]*types.Send, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.SendAllPending(context.Background(), &types.QueryAllSendPendingRequest{})
	if err != nil {
		log.Error().Err(err).Msg("query SendAllPending error")
		return nil, err
	}
	return resp.Send, nil
}

func (b *ZetaCoreBridge) GetAllReceive() ([]*types.Receive, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.ReceiveAll(context.Background(), &types.QueryAllReceiveRequest{})
	if err != nil {
		log.Error().Err(err).Msg("query GetAllReceive error")
		return nil, err
	}
	return resp.Receive, nil
}

func (b *ZetaCoreBridge) GetLastBlockHeight() ([]*types.LastBlockHeight, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastBlockHeightAll(context.Background(), &types.QueryAllLastBlockHeightRequest{})
	if err != nil {
		log.Warn().Err(err).Msg("query GetLastBlockHeight error")
		return nil, err
	}
	return resp.LastBlockHeight, nil
}

func (b *ZetaCoreBridge) GetZetaBlockHeight() (uint64, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastMetaHeight(context.Background(), &types.QueryLastMetaHeightRequest{})
	if err != nil {
		log.Warn().Err(err).Msg("query GetLastBlockHeight error")
		return 0, err
	}
	return resp.Height, nil
}

func (b *ZetaCoreBridge) GetLastBlockHeightByChain(chain common.Chain) (*types.LastBlockHeight, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastBlockHeight(context.Background(), &types.QueryGetLastBlockHeightRequest{Index: chain.String()})
	if err != nil {
		log.Error().Err(err).Msg("query GetLastBlockHeight error")
		return nil, err
	}
	return resp.LastBlockHeight, nil
}

func (b *ZetaCoreBridge) GetNonceByChain(chain common.Chain) (*types.ChainNonces, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.ChainNonces(context.Background(), &types.QueryGetChainNoncesRequest{Index: chain.String()})
	if err != nil {
		log.Error().Err(err).Msg("query GetNonceByChain error")
		return nil, err
	}
	return resp.ChainNonces, nil
}

func (b *ZetaCoreBridge) SetNodeKey(pubkeyset common.PubKeySet, conskey string) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgSetNodeKeys(signerAddress, pubkeyset, conskey)
	zetaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("SetNodeKey broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) GetKeyGen() (*types.Keygen, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.Keygen(context.Background(), &types.QueryGetKeygenRequest{})
	if err != nil {
		//log.Error().Err(err).Msg("query GetKeyGen error")
		return nil, err
	}
	return resp.Keygen, nil
}

func (b *ZetaCoreBridge) SetTSS(chain common.Chain, address string, pubkey string) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgCreateTSSVoter(signerAddress, chain.String(), address, pubkey)
	zetaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("SetNodeKey broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}
