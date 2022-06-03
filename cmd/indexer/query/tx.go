package query

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"google.golang.org/grpc"

	tmservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/query"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

type ZetaQuerier struct {
	grpcConn *grpc.ClientConn
}

func NewZetaQuerier(chainIP string) (*ZetaQuerier, error) {
	grpcConn, err := grpc.Dial(
		fmt.Sprintf("%s:9090", chainIP),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Error().Err(err).Msg("ZetaCore grpc dial fail")
		return nil, err
	}
	return &ZetaQuerier{grpcConn: grpcConn}, nil
}

func (q *ZetaQuerier) LatestBlock() (*tmtypes.Block, error) {
	client := tmservice.NewServiceClient(q.grpcConn)
	res, err := client.GetLatestBlock(context.Background(), &tmservice.GetLatestBlockRequest{})
	if err != nil {
		fmt.Printf("GetLatestBlock grpc err: %s\n", err)
		return nil, err
	}
	return res.Block, nil
}

func (q *ZetaQuerier) BlockByHeight(bn int64) (*tmtypes.Block, error) {
	client := tmservice.NewServiceClient(q.grpcConn)
	res, err := client.GetBlockByHeight(context.Background(), &tmservice.GetBlockByHeightRequest{
		Height: bn,
	})
	if err != nil {
		fmt.Printf("GetLatestBlock grpc err: %s\n", err)
		return nil, err
	}
	return res.Block, nil
}

// queries native txs that belong to a block height
func (q *ZetaQuerier) TxResponsesByBlock(bn int64) ([]*sdk.TxResponse, error) {
	client := txtypes.NewServiceClient(q.grpcConn)
	events := []string{fmt.Sprintf("tx.height=%d", bn)}

	res, err := client.GetTxsEvent(context.Background(), &txtypes.GetTxsEventRequest{
		Events:     events,
		Pagination: nil,
		OrderBy:    0,
	})

	if err != nil {
		fmt.Printf("GetTxsEvent grpc err: %s\n", err)
		return nil, err
	}
	return res.TxResponses, nil
}

// queries native txs that belong to a block height
func (q *ZetaQuerier) TxByHash(hash string) (*txtypes.Tx, error) {
	client := txtypes.NewServiceClient(q.grpcConn)

	res, err := client.GetTx(context.Background(), &txtypes.GetTxRequest{
		Hash: hash,
	})

	if err != nil {
		fmt.Printf("GetTxsEvent grpc err: %s\n", err)
		return nil, err
	}
	return res.Tx, nil
}

// query events of subtype at block blockNum
// if blockNum <0, then query from block 0
// each tx_response will be processed by the function processTxResponses
func (q *ZetaQuerier) VisitAllTxEvents(subtype string, blockNum int64, processTxResponses func(txRes *sdk.TxResponse) error) (uint64, error) {
	const PAGE_LIMIT = 50
	client := txtypes.NewServiceClient(q.grpcConn)
	var offset, processed uint64

	events := []string{fmt.Sprintf("message.%s='%s'", types.SubTypeKey, subtype)}
	if blockNum >= 0 {
		events = append(events, fmt.Sprintf("tx.height=%d", blockNum))
	}

	// first call
	// NOTE: OrderBy 0 appears to be ASC block height
	processed = 0
	res, err := client.GetTxsEvent(context.Background(), &txtypes.GetTxsEventRequest{
		Events: events,
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      PAGE_LIMIT,
			CountTotal: false,
			Reverse:    false,
		},
		OrderBy: 0,
	})
	if err != nil {
		log.Error().Err(err).Msg("GetTxsEvent grpc fail")
		return processed, err
	}
	for _, v := range res.TxResponses {
		err = processTxResponses(v)
		if err != nil {
			return processed, err
		}
		processed += 1
	}

	// subsequent calls if necessary for paging
	if res.Pagination.Total > PAGE_LIMIT {
		for offset = PAGE_LIMIT; offset < res.Pagination.Total; offset += PAGE_LIMIT {
			res, err = client.GetTxsEvent(context.Background(), &txtypes.GetTxsEventRequest{
				Events: events,
				Pagination: &query.PageRequest{
					Key:        nil,
					Offset:     offset,
					Limit:      PAGE_LIMIT,
					CountTotal: false,
					Reverse:    false,
				},
				OrderBy: 0,
			})
			if err != nil {
				log.Error().Err(err).Msgf("GetTxsEvent error of %v", events)
			}
			for _, v := range res.TxResponses {
				err = processTxResponses(v)
				if err != nil {
					fmt.Printf("processTxResponses error: %s", err)
				}
				processed += 1
			}
		}
	}

	return processed, nil
}
