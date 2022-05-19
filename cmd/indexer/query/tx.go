package query

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/types/query"
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

// query events of subtype at block blockNum
// if blockNum <0, then query from
// each tx_response will be processed by the function processTxResponses
func (q *ZetaQuerier) VisitAllTxEvents(subtype string, blockNum int64, processTxResponses func(txRes *sdk.TxResponse) error) (uint64, error) {
	const PAGE_LIMIT = 2
	client := txtypes.NewServiceClient(q.grpcConn)
	var offset, processed uint64

	events := []string{fmt.Sprintf("message.Subtype='%s'", subtype)}

	// first call
	// NOTE: OrderBy 0 appears to be ASC block height
	offset = 0
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
			for _, v := range res.TxResponses {
				err = processTxResponses(v)
				if err != nil {
					return processed, err
				}
				processed += 1
			}
		}
	}

	return processed, nil
}
