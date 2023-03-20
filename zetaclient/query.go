package zetaclient

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/query"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc"
)

func (b *ZetaCoreBridge) GetInboundPermissions() (types.PermissionFlags, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.PermissionFlags(context.Background(), &types.QueryGetPermissionFlagsRequest{})
	if err != nil {
		b.logger.Error().Err(err).Msg("Query permissions failed")
		return types.PermissionFlags{}, err
	}
	return resp.PermissionFlags, nil

}

//func (b *ZetaCoreBridge) GetAccountDetails(address string) (string, error) {
//	client := authtypes.NewQueryClient(b.grpcConn)
//	resp, err := client.Account(context.Background(), &authtypes.QueryAccountRequest{
//		Address: address,
//	})
//	if err != nil {
//		b.logger.Error().Err(err).Msg("Query account failed")
//		return "", err
//	}
//
//	err := resp.UnpackInterfaces
//	return resp.Account.GetTypeUrl(), nil
//
//}

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
