package zetaclient

import (
	"context"
	"fmt"
	"sort"

	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/query"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc"
)

type Order string

const (
	NoOrder    Order = ""
	Ascending  Order = "ASC"
	Descending Order = "DESC"
)

func (b *ZetaCoreBridge) GetCrosschainFlags() (zetaObserverTypes.CrosschainFlags, error) {
	client := zetaObserverTypes.NewQueryClient(b.grpcConn)
	resp, err := client.CrosschainFlags(context.Background(), &zetaObserverTypes.QueryGetCrosschainFlagsRequest{})
	if err != nil {
		return zetaObserverTypes.CrosschainFlags{}, err
	}
	return resp.CrosschainFlags, nil
}

func (b *ZetaCoreBridge) GetCoreParamsForChainID(externalChainID int64) (*zetaObserverTypes.CoreParams, error) {
	client := zetaObserverTypes.NewQueryClient(b.grpcConn)
	resp, err := client.GetCoreParamsForChain(context.Background(), &zetaObserverTypes.QueryGetCoreParamsForChainRequest{ChainId: externalChainID})
	if err != nil {
		return &zetaObserverTypes.CoreParams{}, err
	}
	return resp.CoreParams, nil
}

func (b *ZetaCoreBridge) GetCoreParams() ([]*zetaObserverTypes.CoreParams, error) {
	client := zetaObserverTypes.NewQueryClient(b.grpcConn)
	var err error
	resp := &zetaObserverTypes.QueryGetCoreParamsResponse{}
	for i := 0; i <= DefaultRetryCount; i++ {
		resp, err = client.GetCoreParams(context.Background(), &zetaObserverTypes.QueryGetCoreParamsRequest{})
		if err == nil {
			return resp.CoreParams.CoreParams, nil
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil, fmt.Errorf("failed to get core params | err %s", err.Error())
}

func (b *ZetaCoreBridge) GetObserverParams() (zetaObserverTypes.Params, error) {
	client := zetaObserverTypes.NewQueryClient(b.grpcConn)
	resp, err := client.Params(context.Background(), &zetaObserverTypes.QueryParamsRequest{})
	if err != nil {
		return zetaObserverTypes.Params{}, err
	}
	return resp.Params, nil
}

func (b *ZetaCoreBridge) GetUpgradePlan() (*upgradetypes.Plan, error) {
	client := upgradetypes.NewQueryClient(b.grpcConn)

	resp, err := client.CurrentPlan(context.Background(), &upgradetypes.QueryCurrentPlanRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Plan, nil
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
		return nil, err
	}
	return resp.CrossChainTx, nil
}

func (b *ZetaCoreBridge) GetCctxByHash(sendHash string) (*types.CrossChainTx, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.Cctx(context.Background(), &types.QueryGetCctxRequest{Index: sendHash})
	if err != nil {
		return nil, err
	}
	return resp.CrossChainTx, nil
}

func (b *ZetaCoreBridge) GetCctxByNonce(chainID int64, nonce uint64) (*types.CrossChainTx, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.CctxByNonce(context.Background(), &types.QueryGetCctxByNonceRequest{
		ChainID: chainID,
		Nonce:   nonce,
	})
	if err != nil {
		return nil, err
	}
	return resp.CrossChainTx, nil
}

func (b *ZetaCoreBridge) GetObserverList(chain common.Chain) ([]string, error) {
	var err error

	client := zetaObserverTypes.NewQueryClient(b.grpcConn)
	for i := 0; i <= DefaultRetryCount; i++ {
		resp, err := client.ObserversByChain(context.Background(), &zetaObserverTypes.QueryObserversByChainRequest{ObservationChain: chain.ChainName.String()})
		if err == nil {
			return resp.Observers, nil
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil, err
}

func (b *ZetaCoreBridge) GetAllPendingCctx(chainID int64) ([]*types.CrossChainTx, error) {
	client := types.NewQueryClient(b.grpcConn)
	maxSizeOption := grpc.MaxCallRecvMsgSize(32 * 1024 * 1024)
	resp, err := client.CctxAllPending(context.Background(), &types.QueryAllCctxPendingRequest{ChainId: chainID}, maxSizeOption)
	if err != nil {
		return nil, err
	}
	return resp.CrossChainTx, nil
}

func (b *ZetaCoreBridge) GetLastBlockHeight() ([]*types.LastBlockHeight, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastBlockHeightAll(context.Background(), &types.QueryAllLastBlockHeightRequest{})
	if err != nil {
		b.logger.Error().Err(err).Msg("query GetBlockHeight error")
		return nil, err
	}
	return resp.LastBlockHeight, nil
}

func (b *ZetaCoreBridge) GetLatestZetaBlock() (*tmtypes.Block, error) {
	client := tmservice.NewServiceClient(b.grpcConn)
	res, err := client.GetLatestBlock(context.Background(), &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, err
	}
	return res.Block, nil
}

func (b *ZetaCoreBridge) GetNodeInfo() (*tmservice.GetNodeInfoResponse, error) {
	var err error

	client := tmservice.NewServiceClient(b.grpcConn)
	for i := 0; i <= DefaultRetryCount; i++ {
		res, err := client.GetNodeInfo(context.Background(), &tmservice.GetNodeInfoRequest{})
		if err == nil {
			return res, nil
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil, err
}

func (b *ZetaCoreBridge) GetLastBlockHeightByChain(chain common.Chain) (*types.LastBlockHeight, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastBlockHeight(context.Background(), &types.QueryGetLastBlockHeightRequest{Index: chain.ChainName.String()})
	if err != nil {
		return nil, err
	}
	return resp.LastBlockHeight, nil
}

func (b *ZetaCoreBridge) GetZetaBlockHeight() (int64, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastZetaHeight(context.Background(), &types.QueryLastZetaHeightRequest{})
	if err != nil {
		return 0, err
	}
	return resp.Height, nil
}

func (b *ZetaCoreBridge) GetNonceByChain(chain common.Chain) (*types.ChainNonces, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.ChainNonces(context.Background(), &types.QueryGetChainNoncesRequest{Index: chain.ChainName.String()})
	if err != nil {
		return nil, err
	}
	return resp.ChainNonces, nil
}

func (b *ZetaCoreBridge) GetAllNodeAccounts() ([]*zetaObserverTypes.NodeAccount, error) {
	client := zetaObserverTypes.NewQueryClient(b.grpcConn)
	resp, err := client.NodeAccountAll(context.Background(), &zetaObserverTypes.QueryAllNodeAccountRequest{})
	if err != nil {
		return nil, err
	}
	b.logger.Debug().Msgf("GetAllNodeAccounts: %d", len(resp.NodeAccount))
	return resp.NodeAccount, nil
}

func (b *ZetaCoreBridge) GetKeyGen() (*zetaObserverTypes.Keygen, error) {
	var err error

	client := zetaObserverTypes.NewQueryClient(b.grpcConn)
	for i := 0; i <= ExtendedRetryCount; i++ {
		resp, err := client.Keygen(context.Background(), &zetaObserverTypes.QueryGetKeygenRequest{})
		if err == nil {
			return resp.Keygen, nil
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil, fmt.Errorf("failed to get keygen | err %s", err.Error())

}

func (b *ZetaCoreBridge) GetBallot(ballotIdentifier string) (*zetaObserverTypes.QueryBallotByIdentifierResponse, error) {
	client := zetaObserverTypes.NewQueryClient(b.grpcConn)
	resp, err := client.BallotByIdentifier(context.Background(), &zetaObserverTypes.QueryBallotByIdentifierRequest{BallotIdentifier: ballotIdentifier})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *ZetaCoreBridge) GetInboundTrackersForChain(chainID int64) ([]types.InTxTracker, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.InTxTrackerAllByChain(context.Background(), &types.QueryAllInTxTrackerByChainRequest{ChainId: chainID})
	if err != nil {
		return nil, err
	}
	return resp.InTxTracker, nil
}

func (b *ZetaCoreBridge) GetCurrentTss() (*types.TSS, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.TSS(context.Background(), &types.QueryGetTSSRequest{})
	if err != nil {
		return nil, err
	}
	return resp.TSS, nil
}

func (b *ZetaCoreBridge) GetEthTssAddress() (string, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.GetTssAddress(context.Background(), &types.QueryGetTssAddressRequest{})
	if err != nil {
		return "", err
	}
	return resp.Eth, nil
}

func (b *ZetaCoreBridge) GetBtcTssAddress() (string, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.GetTssAddress(context.Background(), &types.QueryGetTssAddressRequest{})
	if err != nil {
		return "", err
	}
	return resp.Btc, nil
}

func (b *ZetaCoreBridge) GetTssHistory() ([]types.TSS, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.TssHistory(context.Background(), &types.QueryTssHistoryRequest{})
	if err != nil {
		return nil, err
	}
	return resp.TssList, nil
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

func (b *ZetaCoreBridge) GetAllOutTxTrackerByChain(chain common.Chain, order Order) ([]types.OutTxTracker, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.OutTxTrackerAllByChain(context.Background(), &types.QueryAllOutTxTrackerByChainRequest{
		Chain: chain.ChainId,
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      2000,
			CountTotal: false,
			Reverse:    false,
		},
	})
	if err != nil {
		return nil, err
	}
	if order == Ascending {
		sort.SliceStable(resp.OutTxTracker, func(i, j int) bool {
			return resp.OutTxTracker[i].Nonce < resp.OutTxTracker[j].Nonce
		})
	}
	if order == Descending {
		sort.SliceStable(resp.OutTxTracker, func(i, j int) bool {
			return resp.OutTxTracker[i].Nonce > resp.OutTxTracker[j].Nonce
		})
	}
	return resp.OutTxTracker, nil
}

func (b *ZetaCoreBridge) GetClientParams(chainID int64) (zetaObserverTypes.QueryGetCoreParamsForChainResponse, error) {
	client := zetaObserverTypes.NewQueryClient(b.grpcConn)
	resp, err := client.GetCoreParamsForChain(context.Background(), &zetaObserverTypes.QueryGetCoreParamsForChainRequest{ChainId: chainID})
	if err != nil {
		return zetaObserverTypes.QueryGetCoreParamsForChainResponse{}, err
	}
	return *resp, nil
}

func (b *ZetaCoreBridge) GetPendingNoncesByChain(chainID int64) (types.PendingNonces, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.PendingNoncesByChain(context.Background(), &types.QueryPendingNoncesByChainRequest{ChainId: chainID})
	if err != nil {
		return types.PendingNonces{}, err
	}
	return resp.PendingNonces, nil
}

func (b *ZetaCoreBridge) GetSupportedChains() ([]*common.Chain, error) {
	client := zetaObserverTypes.NewQueryClient(b.grpcConn)
	resp, err := client.SupportedChains(context.Background(), &zetaObserverTypes.QuerySupportedChains{})
	if err != nil {
		return nil, err
	}
	return resp.GetChains(), nil
}

func (b *ZetaCoreBridge) GetPendingNonces() (*types.QueryAllPendingNoncesResponse, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.PendingNoncesAll(context.Background(), &types.QueryAllPendingNoncesRequest{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *ZetaCoreBridge) Prove(blockHash string, txHash string, txIndex int64, proof *common.Proof, chainID int64) (bool, error) {
	client := zetaObserverTypes.NewQueryClient(b.grpcConn)
	resp, err := client.Prove(context.Background(), &zetaObserverTypes.QueryProveRequest{
		BlockHash: blockHash,
		TxIndex:   txIndex,
		Proof:     proof,
		ChainId:   chainID,
		TxHash:    txHash,
	})
	if err != nil {
		return false, err
	}
	return resp.Valid, nil
}
