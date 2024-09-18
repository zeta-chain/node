package liteapi

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// Client extends tongo's liteapi.Client with some high-level tools
type Client struct {
	*liteapi.Client
}

// GetFirstTransaction scrolls through the transactions of the given account to find the first one.
// Note that it will fail in case of old transactions. Ideally, use archival node.
// Also returns the number of scrolled transactions for this account i.e. total transactions
func (c *Client) GetFirstTransaction(ctx context.Context, acc ton.AccountID) (*ton.Transaction, int, error) {
	const pageSize = 100

	state, err := c.GetAccountState(ctx, acc)
	if err != nil {
		return nil, 0, errors.Wrap(err, "unable to get account state")
	}

	if state.Account.Status() != tlb.AccountActive {
		return nil, 0, errors.New("account is not active")
	}

	var tx *ton.Transaction

	// logical time and hash of the last transaction
	lt, hash, scrolled := state.LastTransLt, state.LastTransHash, 0

	for {
		hashBits := ton.Bits256(hash)

		txs, err := c.GetTransactions(ctx, pageSize, acc, lt, hashBits)
		if err != nil {
			return nil, scrolled, errors.Wrapf(err, "unable to get transactions [lt %d, hash %s]", lt, hashBits.Hex())
		}

		if len(txs) == 0 {
			break
		}

		scrolled += len(txs)

		tx = &txs[len(txs)-1]

		// Not we take the latest item in the list (oldest tx in the page)
		// and set it as the new last tx
		lt, hash = tx.PrevTransLt, tx.PrevTransHash
	}

	if tx == nil {
		return nil, scrolled, fmt.Errorf("no transactions found [lt %d, hash %s]", lt, ton.Bits256(hash).Hex())
	}

	return tx, scrolled, nil
}

func TransactionHashToString(lt uint64, hash ton.Bits256) string {
	return fmt.Sprintf("%d:%s", lt, hash.Hex())
}

func TransactionHashFromString(hash string) (uint64, ton.Bits256, error) {
	parts := strings.Split(hash, ":")
	if len(parts) != 2 {
		return 0, ton.Bits256{}, fmt.Errorf("invalid hash string format")
	}

	lt, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, ton.Bits256{}, fmt.Errorf("invalid logical time: %w", err)
	}

	var hashBits ton.Bits256

	if err = hashBits.FromHex(parts[1]); err != nil {
		return 0, ton.Bits256{}, fmt.Errorf("invalid hash: %w", err)
	}

	return lt, hashBits, nil
}
