package zetaclient

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/zetaobserver/types"
	"math/big"
	"time"
)

func (b *ZetaCoreBridge) PostZetaConversionRate(chain common.Chain, rate string, blockNum uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgZetaConversionRateVoter(signerAddress, chain.String(), rate, blockNum)
	zetaTxHash, err := b.Broadcast(msg)
	if err != nil {
		b.logger.Error().Err(err).Msg("PostZetaConversionRate broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) PostGasBalance(chain common.Chain, gasBalance string, blockNum uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgGasBalanceVoter(signerAddress, chain.String(), gasBalance, blockNum)
	zetaTxHash, err := b.Broadcast(msg)
	if err != nil {
		b.logger.Error().Err(err).Msg("PostGasPrice broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) PostGasPrice(chain common.Chain, gasPrice uint64, supply string, blockNum uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgGasPriceVoter(signerAddress, chain.String(), gasPrice, supply, blockNum)
	zetaTxHash, err := b.Broadcast(msg)
	if err != nil {
		b.logger.Error().Err(err).Msg("PostGasPrice broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) AddTxHashToOutTxTracker(chain string, nonce uint64, txHash string) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgAddToOutTxTracker(signerAddress, chain, nonce, txHash)
	zetaTxHash, err := b.Broadcast(msg)
	if err != nil {
		b.logger.Error().Err(err).Msg("PostGasPrice broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) PostNonce(chain common.Chain, nonce uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgNonceVoter(signerAddress, chain.String(), nonce)
	zetaTxHash, err := b.Broadcast(msg)
	if err != nil {
		b.logger.Error().Err(err).Msg("PostNonce broadcast fail")
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
			b.logger.Error().Err(err).Msg("PostSend broadcast fail; re-trying...")
		} else {
			return zetaTxHash, nil
		}
		time.Sleep(1 * time.Second)
	}
	return zetaTxHash, fmt.Errorf("postSend: re-try fails")
}

// FIXME: pass nonce
func (b *ZetaCoreBridge) PostReceiveConfirmation(sendHash string, outTxHash string, outBlockHeight uint64, mMint *big.Int, status common.ReceiveStatus, chain string, nonce int) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgReceiveConfirmation(signerAddress, sendHash, outTxHash, outBlockHeight, sdk.NewUintFromBigInt(mMint), status, chain, uint64(nonce))
	//b.logger.Info().Msgf("PostReceiveConfirmation msg digest: %s", msg.Digest())
	var zetaTxHash string
	for i := 0; i < 2; i++ {
		zetaTxHash, err := b.Broadcast(msg)
		if err != nil {
			b.logger.Error().Err(err).Msg("PostReceiveConfirmation broadcast fail; re-trying...")
		} else {
			return zetaTxHash, nil
		}
		time.Sleep(1 * time.Second)
	}
	return zetaTxHash, fmt.Errorf("postReceiveConfirmation: re-try fails")
}

func (b *ZetaCoreBridge) GetAllCctx() ([]*types.CrossChainTx, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.CctxAll(context.Background(), &types.QueryAllCctxRequest{})
	if err != nil {
		b.logger.Error().Err(err).Msg("query CctxAll error")
		return nil, err
	}
	return resp.CrossChainTx, nil
}

func (b *ZetaCoreBridge) GetCctxByHash(sendHash string) (*types.CrossChainTx, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.Cctx(context.Background(), &types.QueryGetCctxRequest{Index: sendHash})
	if err != nil {
		b.logger.Error().Err(err).Msg("GetCctxByHash error")
		return nil, err
	}
	return resp.CrossChainTx, nil
}

func (b *ZetaCoreBridge) GetObserverList(chain common.Chain, observationType string) ([]string, error) {
	client := zetaObserverTypes.NewQueryClient(b.grpcConn)
	b.logger.Info().Msgf("GetObserverList resp: %s, %s", chain.String, observationType)
	resp, err := client.ObserversByChainAndType(context.Background(), &zetaObserverTypes.QueryObserversByChainAndTypeRequest{
		ObservationChain: chain.String(),
		ObservationType:  observationType,
	})
	if err != nil {
		b.logger.Error().Err(err).Msg("query GetObserverList error")
		return nil, err
	}
	return resp.Observers, nil
}

func (b *ZetaCoreBridge) GetAllPendingCctx() ([]*types.CrossChainTx, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.CctxAllPending(context.Background(), &types.QueryAllCctxPendingRequest{})
	if err != nil {
		b.logger.Error().Err(err).Msg("query CctxAllPending error")
		return nil, err
	}
	return resp.CrossChainTx, nil
}

func (b *ZetaCoreBridge) GetLastBlockHeight() ([]*types.LastBlockHeight, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastBlockHeightAll(context.Background(), &types.QueryAllLastBlockHeightRequest{})
	if err != nil {
		b.logger.Warn().Err(err).Msg("query GetLastBlockHeight error")
		return nil, err
	}
	return resp.LastBlockHeight, nil
}

func (b *ZetaCoreBridge) GetLatestZetaBlock() (*tmtypes.Block, error) {
	client := tmservice.NewServiceClient(b.grpcConn)
	res, err := client.GetLatestBlock(context.Background(), &tmservice.GetLatestBlockRequest{})
	if err != nil {
		fmt.Printf("GetLatestBlock grpc err: %s\n", err)
		return nil, err
	}
	return res.Block, nil
}

func (b *ZetaCoreBridge) GetLastBlockHeightByChain(chain common.Chain) (*types.LastBlockHeight, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastBlockHeight(context.Background(), &types.QueryGetLastBlockHeightRequest{Index: chain.String()})
	if err != nil {
		b.logger.Error().Err(err).Msg("query GetLastBlockHeight error")
		return nil, err
	}
	return resp.LastBlockHeight, nil
}

func (b *ZetaCoreBridge) GetZetaBlockHeight() (uint64, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastMetaHeight(context.Background(), &types.QueryLastMetaHeightRequest{})
	if err != nil {
		b.logger.Warn().Err(err).Msg("query GetLastBlockHeight error")
		return 0, err
	}
	return resp.Height, nil
}

func (b *ZetaCoreBridge) GetNonceByChain(chain common.Chain) (*types.ChainNonces, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.ChainNonces(context.Background(), &types.QueryGetChainNoncesRequest{Index: chain.String()})
	if err != nil {
		b.logger.Error().Err(err).Msg("query QueryGetChainNoncesRequest error")
		return nil, err
	}
	return resp.ChainNonces, nil
}

func (b *ZetaCoreBridge) SetNodeKey(pubkeyset common.PubKeySet, conskey string) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgSetNodeKeys(signerAddress, pubkeyset, conskey)
	zetaTxHash, err := b.Broadcast(msg)
	if err != nil {
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) GetAllNodeAccounts() ([]*types.NodeAccount, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.NodeAccountAll(context.Background(), &types.QueryAllNodeAccountRequest{})
	if err != nil {
		b.logger.Error().Err(err).Msg("query GetAllNodeAccounts error")
		return nil, err
	}
	b.logger.Info().Msgf("GetAllNodeAccounts: %d", len(resp.NodeAccount))

	return resp.NodeAccount, nil
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
		b.logger.Err(err).Msg("SetNodeKey broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) GetOutTxTracker(chain common.Chain, nonce uint64) (*types.OutTxTracker, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.OutTxTracker(context.Background(), &types.QueryGetOutTxTrackerRequest{
		Index: fmt.Sprintf("%s-%d", chain.String(), nonce),
	})
	if err != nil {
		return nil, err
	}
	return &resp.OutTxTracker, nil
}

func (b *ZetaCoreBridge) GetAllOutTxTrackerByChain(chain common.Chain) ([]types.OutTxTracker, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.OutTxTrackerAllByChain(context.Background(), &types.QueryAllOutTxTrackerByChainRequest{
		Chain: chain.String(),
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      300,
			CountTotal: false,
			Reverse:    false,
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.OutTxTracker, nil
}
