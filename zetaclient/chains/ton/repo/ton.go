package repo

import (
	"context"
	"errors"

	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/node/pkg/chains"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/zetaclient/chains/ton/encoder"
	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
)

type TONRepo struct {
	// TODO: make these private before opening the pull request
	Client  TONClient
	Gateway *toncontracts.Gateway

	connectedChain chains.Chain
}

func NewTONRepo(tonClient TONClient,
	gateway *toncontracts.Gateway,
	connectedChain chains.Chain,
) *TONRepo {
	return &TONRepo{
		Client:         tonClient,
		Gateway:        gateway,
		connectedChain: connectedChain,
	}
}

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

// TODO: this function seems very wrong.
// GetLastTransaction TODO.
func (repo *TONRepo) GetLastTransaction(ctx context.Context) (string, error) {
	const limit = 20
	accountID := repo.Gateway.AccountID()
	var zeroLT uint64
	var zeroHash ton.Bits256

	txs, err := repo.Client.GetTransactions(ctx, limit, accountID, zeroLT, zeroHash)
	if err != nil {
		return "", errors.Join(ErrGetTransactions, err)
	}
	if len(txs) == 0 {
		return "", ErrNoTransactions
	}

	tx := txs[len(txs)-1]
	hash := encoder.EncodeTx(tx)

	return hash, nil
}
