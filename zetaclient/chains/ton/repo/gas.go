package repo

import (
	"context"
	"errors"

	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
)

// GetGasPrice returns the most recent gas price and the number of the last block.
func (repo *TONRepo) GetGasPrice(ctx context.Context) (uint64, uint64, error) {
	rawGasPrice, err := rpc.FetchGasConfigRPC(ctx, repo.Client)
	if err != nil {
		return 0, 0, errors.Join(ErrFetchGasPrice, err)
	}

	gasPrice, err := rpc.ParseGasPrice(rawGasPrice)
	if err != nil {
		return 0, 0, errors.Join(ErrParseGasPrice, err)
	}

	info, err := repo.Client.GetMasterchainInfo(ctx)
	if err != nil {
		return gasPrice, 0, errors.Join(ErrGetMasterchainInfo, err)
	}
	lastBlockNumber := uint64(info.Last.Seqno)

	return gasPrice, lastBlockNumber, nil
}

// // TODO
// func (repo *TONRepo) VoteGasPrice(ctx context.Context,
// 	chain chains.Chain,
// 	gasPrice uint64,
// 	blockNumber uint64,
// ) error {
// 	// There is no concept of priority fee in TON.
// 	const priorityFee = 0
//
// 	_, err := repo.ZetacoreClient.PostVoteGasPrice(ctx, chain, gasPrice, priorityFee, blockNumber)
// 	if err != nil {
// 		return errors.Join(ErrPostVoteGasPrice, err)
// 	}
//
// 	return nil
// }
