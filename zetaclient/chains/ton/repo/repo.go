package repo

import (
	"context"
	"errors"
	"time"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/node/pkg/chains"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
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
	hash := rpc.TransactionToHashString(tx)

	return hash, nil
}

// GetTransactionsSince TODO.
func (repo *TONRepo) GetTransactionsSince(ctx context.Context,
	lastTx string,
) ([]ton.Transaction, error) {
	accountID := repo.Gateway.AccountID()

	lastLT, lastHash, err := rpc.TransactionHashFromString(lastTx)
	if err != nil {
		return nil, errors.Join(ErrTransactionEncoding, err)
	}

	txs, err := repo.Client.GetTransactionsSince(ctx, accountID, lastLT, lastHash)
	if err != nil {
		return nil, errors.Join(ErrGetTransactionsSince, err)
	}

	return txs, nil
}

// func (repo *TONRepo) GetInboundTrackers(ctx context.Context) ([]types.InboundTracker, error) {
// 	chainID := repo.connectedChain.ChainId
// 	trackers, err := repo.ZetacoreClient.GetInboundTrackersForChain(ctx, chainID)
// 	if err != nil {
// 		return nil, errors.Join(ErrGetInboundTrackers, err)
// 	}
// 	return trackers, nil
// }

// ------------------------------------------------------------------------------------------------

// TODO: Duplicate; remove this before merge.
type TONClient interface {
	GetConfigParam(_ context.Context, index uint32) (*boc.Cell, error)

	GetBlockHeader(_ context.Context, blockID rpc.BlockIDExt) (rpc.BlockHeader, error)

	GetMasterchainInfo(context.Context) (rpc.MasterchainInfo, error)

	HealthCheck(context.Context) (time.Time, error)

	GetTransaction(_ context.Context,
		_ ton.AccountID,
		lt uint64,
		hash ton.Bits256,
	) (ton.Transaction, error)

	GetTransactions(_ context.Context,
		count uint32,
		_ ton.AccountID,
		lt uint64,
		hash ton.Bits256,
	) ([]ton.Transaction, error)

	GetTransactionsSince(_ context.Context,
		_ ton.AccountID,
		lt uint64,
		hash ton.Bits256,
	) ([]ton.Transaction, error)
}
