package repo

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/node/pkg/chains"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/zetaclient/chains/ton/encoder"
	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
)

const PaginationLimit = 100

type TONRepo struct {
	// TODO: make these private before opening the pull request
	client TONClient

	gatewayAccountID ton.AccountID
	connectedChain   chains.Chain
}

func NewTONRepo(tonClient TONClient,
	gateway *toncontracts.Gateway,
	connectedChain chains.Chain,
) *TONRepo {
	return &TONRepo{
		client:           tonClient,
		gatewayAccountID: gateway.AccountID(),
		connectedChain:   connectedChain,
	}
}

// CheckHealth checks the client's health and returns the most recent block time.
func (repo *TONRepo) CheckHealth(ctx context.Context) (*time.Time, error) {
	blockTime, err := repo.client.HealthCheck(ctx)
	if err != nil {
		return nil, errors.Join(ErrHealthCheck, err)
	}

	return &blockTime, nil
}

// GetGasPrice returns the most recent gas price and the number of the last block.
func (repo *TONRepo) GetGasPrice(ctx context.Context) (uint64, uint64, error) {
	rawGasPrice, err := rpc.FetchGasConfigRPC(ctx, repo.client)
	if err != nil {
		return 0, 0, errors.Join(ErrFetchGasPrice, err)
	}

	gasPrice, err := rpc.ParseGasPrice(rawGasPrice)
	if err != nil {
		return 0, 0, errors.Join(ErrParseGasPrice, err)
	}

	info, err := repo.client.GetMasterchainInfo(ctx)
	if err != nil {
		return gasPrice, 0, errors.Join(ErrGetMasterchainInfo, err)
	}
	lastBlockNumber := uint64(info.Last.Seqno)

	return gasPrice, lastBlockNumber, nil
}

// GetTransactionByHash returns the transaction associated with a given encoded hash.
func (repo *TONRepo) GetTransactionByHash(ctx context.Context,
	encodedHash string,
) (*ton.Transaction, error) {
	lt, hash, err := encoder.DecodeHash(encodedHash)
	if err != nil {
		return nil, errors.Join(ErrEncoding, err)
	}

	raw, err := repo.client.GetTransaction(ctx, repo.gatewayAccountID, lt, hash)
	if err != nil {
		return nil, errors.Join(ErrGetTransaction, err)
	}

	return &raw, nil
}

// GetTransactionByIndex returns the Nth most recent transaction.
// (Or the oldest transaction available if there are fewer than N transactions in the blockchain.)
func (repo *TONRepo) GetTransactionByIndex(ctx context.Context,
	n uint32,
) (*ton.Transaction, error) {
	var zeroLT uint64
	var zeroHash ton.Bits256

	txs, err := repo.client.GetTransactions(ctx, n, repo.gatewayAccountID, zeroLT, zeroHash)
	if err != nil {
		return nil, errors.Join(ErrGetTransactions, err)
	}
	if len(txs) == 0 {
		return nil, ErrNoTransactions
	}

	tx := txs[len(txs)-1]
	return &tx, nil
}

// GetNextTransactions TODO.
// does pagination.
func (repo *TONRepo) GetNextTransactions(ctx context.Context, logger zerolog.Logger,
	encodedHash string,
) ([]ton.Transaction, error) {
	lastLT, lastHash, err := encoder.DecodeHash(encodedHash)
	if err != nil {
		return nil, errors.Join(ErrEncoding, err)
	}

	txs, err := repo.client.GetTransactionsSince(ctx, repo.gatewayAccountID, lastLT, lastHash)
	if err != nil {
		return nil, errors.Join(ErrGetTransactionsSince, err)
	}

	numberOfTransactions := len(txs)
	logger.Info().Int("transactions", numberOfTransactions).Msg("observed some transactions")

	if numberOfTransactions > PaginationLimit {
		logger.Info().
			Int("transactions", numberOfTransactions).
			Int("limit", PaginationLimit).
			Msg("number of transactions exceeds pagination limit; processing only some")
		txs = txs[:PaginationLimit]
	}

	return txs, nil
}
