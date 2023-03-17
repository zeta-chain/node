package zetaclient

import (
	"context"
	"cosmossdk.io/math"
	"fmt"
	"math/big"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/query"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc"
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

func (b *ZetaCoreBridge) PostGasPrice(chain common.Chain, gasPrice uint64, supply string, blockNum uint64) (string, error) {

	signerAddress := b.keys.GetAddress().String()
	msg := types.NewMsgGasPriceVoter(signerAddress, chain.ChainId, gasPrice, supply, blockNum)
	zetaTxHash, err := b.Broadcast(PostGasPriceGasLimit, msg)
	if err != nil {
		b.logger.Error().Err(err).Msg("PostGasPrice broadcast fail")
		return "", err
	}
	b.logger.Debug().Str("zetaTxHash", zetaTxHash).Msg("PostGasPrice broadcast success")

	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) AddTxHashToOutTxTracker(chainID int64, nonce uint64, txHash string) (string, error) {
	signerAddress := b.keys.GetAddress().String()
	msg := types.NewMsgAddToOutTxTracker(signerAddress, chainID, nonce, txHash)
	zetaTxHash, err := b.Broadcast(AddTxHashToOutTxTrackerGasLimit, msg)
	if err != nil {
		b.logger.Error().Err(err).Msg("AddTxHashToOutTxTracker broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) PostNonce(chain common.Chain, nonce uint64) (string, error) {
	signerAddress := b.keys.GetAddress().String()
	msg := types.NewMsgNonceVoter(signerAddress, chain.ChainId, nonce)
	zetaTxHash, err := b.Broadcast(PostNonceGasLimit, msg)
	if err != nil {
		b.logger.Error().Err(err).Msg("PostNonce broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) PostSend(sender string, senderChain int64, txOrigin string, receiver string, receiverChain int64, amount math.Uint, message string, inTxHash string, inBlockHeight uint64, gasLimit uint64, coinType common.CoinType, zetaGasLimit uint64, asset string) (string, error) {
	signerAddress := b.keys.GetAddress().String()
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

	address, err := b.keys.GetSignerInfo().GetAddress()
	if err != nil {
		return "", err
	}
	signerAddress := address.String()
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

func (b *ZetaCoreBridge) GetInboundPermissions() (types.PermissionFlags, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.PermissionFlags(context.Background(), &types.QueryGetPermissionFlagsRequest{})
	if err != nil {
		b.logger.Error().Err(err).Msg("Query permissions failed")
		return types.PermissionFlags{}, err
	}
	return resp.PermissionFlags, nil

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

func (b *ZetaCoreBridge) GetObserverList(chain common.Chain) ([]string, error) {
	client := zetaObserverTypes.NewQueryClient(b.grpcConn)
	resp, err := client.ObserversByChain(context.Background(), &zetaObserverTypes.QueryObserversByChainRequest{
		ObservationChain: chain.ChainName.String(),
	})
	if err != nil {
		b.logger.Error().Err(err).Msg("query GetObserverList error")
		return nil, err
	}
	return resp.Observers, nil
}

func (b *ZetaCoreBridge) GetAllPendingCctx() ([]*types.CrossChainTx, error) {
	client := types.NewQueryClient(b.grpcConn)
	maxSizeOption := grpc.MaxCallRecvMsgSize(32 * 1024 * 1024)
	resp, err := client.CctxAllPending(context.Background(), &types.QueryAllCctxPendingRequest{}, maxSizeOption)
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
		b.logger.Warn().Err(err).Msg("query GetBlockHeight error")
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
	resp, err := client.LastBlockHeight(context.Background(), &types.QueryGetLastBlockHeightRequest{Index: chain.ChainName.String()})
	if err != nil {
		b.logger.Error().Err(err).Msg("query GetBlockHeight error")
		return nil, err
	}
	return resp.LastBlockHeight, nil
}

func (b *ZetaCoreBridge) GetZetaBlockHeight() (int64, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastZetaHeight(context.Background(), &types.QueryLastZetaHeightRequest{})
	if err != nil {
		b.logger.Warn().Err(err).Msg("query GetBlockHeight error")
		return 0, err
	}
	return resp.Height, nil
}

func (b *ZetaCoreBridge) GetNonceByChain(chain common.Chain) (*types.ChainNonces, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.ChainNonces(context.Background(), &types.QueryGetChainNoncesRequest{Index: chain.ChainName.String()})
	if err != nil {
		b.logger.Error().Err(err).Msg("query QueryGetChainNoncesRequest error")
		return nil, err
	}
	return resp.ChainNonces, nil
}

func (b *ZetaCoreBridge) SetNodeKey(pubkeyset common.PubKeySet, conskey string) (string, error) {
	address, err := b.keys.GetSignerInfo().GetAddress()
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
	addr, err := b.keys.GetSignerInfo().GetAddress()
	if err != nil {
		return "", err
	}
	signerAddress := addr.String()
	msg := types.NewMsgCreateTSSVoter(signerAddress, chain.ChainName.String(), address, pubkey)
	zetaTxHash, err := b.Broadcast(DefaultGasLimit, msg)
	if err != nil {
		b.logger.Err(err).Msg("SetNodeKey broadcast fail")
		return "", err
	}
	return zetaTxHash, nil
}

func (b *ZetaCoreBridge) GetOutTxTracker(chain common.Chain, nonce uint64) (*types.OutTxTracker, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.OutTxTracker(context.Background(), &types.QueryGetOutTxTrackerRequest{
		ChainID: chain.ChainId,
		Nonce:   nonce,
	})
	if err != nil {
		return nil, err
	}
	return &resp.OutTxTracker, nil
}

func (b *ZetaCoreBridge) GetAllOutTxTrackerByChain(chain common.Chain) ([]types.OutTxTracker, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.OutTxTrackerAllByChain(context.Background(), &types.QueryAllOutTxTrackerByChainRequest{
		Chain: chain.ChainId,
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
