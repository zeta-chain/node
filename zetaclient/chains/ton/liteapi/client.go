package liteapi

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"

	zetaton "github.com/zeta-chain/node/zetaclient/chains/ton"
)

// Client extends liteapi.Client with some high-level tools
// Reference: https://github.com/ton-blockchain/ton/blob/master/tl/generate/scheme/tonlib_api.tl
type Client struct {
	*liteapi.Client
	blockCache *lru.Cache
}

const (
	pageSize       = 200
	blockCacheSize = 250
)

// New Client constructor.
func New(client *liteapi.Client) *Client {
	blockCache, _ := lru.New(blockCacheSize)

	return &Client{Client: client, blockCache: blockCache}
}

// NewFromSource creates a new client from a URL or a file path.
func NewFromSource(ctx context.Context, urlOrPath string) (*Client, error) {
	cfg, err := zetaton.ConfigFromSource(ctx, urlOrPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get config")
	}

	client, err := liteapi.NewClient(
		liteapi.WithConfigurationFile(*cfg),
		liteapi.WithDetectArchiveNodes(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create client")
	}

	return New(client), nil
}

// GetBlockHeader returns block header by block ID.
// Uses LRU cache for network efficiency.
// I haven't found what mode means but `0` works fine.
func (c *Client) GetBlockHeader(ctx context.Context, blockID ton.BlockIDExt, mode uint32) (tlb.BlockInfo, error) {
	if c.blockCache == nil {
		return tlb.BlockInfo{}, errors.New("block cache is not initialized")
	}

	cached, ok := c.getBlockHeaderCache(blockID)
	if ok {
		return cached, nil
	}

	header, err := c.Client.GetBlockHeader(ctx, blockID, mode)
	if err != nil {
		return tlb.BlockInfo{}, err
	}

	c.setBlockHeaderCache(blockID, header)

	return header, nil
}

func (c *Client) getBlockHeaderCache(blockID ton.BlockIDExt) (tlb.BlockInfo, bool) {
	raw, ok := c.blockCache.Get(blockID.String())
	if !ok {
		return tlb.BlockInfo{}, false
	}

	header, ok := raw.(tlb.BlockInfo)

	return header, ok
}

func (c *Client) setBlockHeaderCache(blockID ton.BlockIDExt, header tlb.BlockInfo) {
	c.blockCache.Add(blockID.String(), header)
}

// GetFirstTransaction scrolls through the transactions of the given account to find the first one.
// Note that it might fail w/o using an archival node. Also returns the number of
// scrolled transactions for this account i.e. total transactions
func (c *Client) GetFirstTransaction(ctx context.Context, acc ton.AccountID) (*ton.Transaction, int, error) {
	lt, hash, err := c.getLastTransactionHash(ctx, acc)
	if err != nil {
		return nil, 0, err
	}

	var (
		tx       *ton.Transaction
		scrolled int
	)

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

// GetTransactionsSince returns all account transactions since the given logicalTime and hash (exclusive).
// The result is ordered from oldest to newest. Used to detect new txs to observe.
func (c *Client) GetTransactionsSince(
	ctx context.Context,
	acc ton.AccountID,
	oldestLT uint64,
	oldestHash ton.Bits256,
) ([]ton.Transaction, error) {
	lt, hash, err := c.getLastTransactionHash(ctx, acc)
	if err != nil {
		return nil, err
	}

	var result []ton.Transaction

	for {
		hashBits := ton.Bits256(hash)

		// note that ton liteapi works in the reverse order.
		// Here we go from the LATEST txs to the oldest at N txs per page
		txs, err := c.GetTransactions(ctx, pageSize, acc, lt, hashBits)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to get transactions [lt %d, hash %s]", lt, hashBits.Hex())
		}

		if len(txs) == 0 {
			break
		}

		for i := range txs {
			found := txs[i].Lt == oldestLT && txs[i].Hash() == tlb.Bits256(oldestHash)
			if !found {
				continue
			}

			// early exit
			result = append(result, txs[:i]...)

			return result, nil
		}

		// otherwise, append all page results
		result = append(result, txs...)

		// prepare pagination params for the next page
		oldestIndex := len(txs) - 1

		lt, hash = txs[oldestIndex].PrevTransLt, txs[oldestIndex].PrevTransHash
	}

	// reverse the result to get the oldest tx first
	slices.Reverse(result)

	return result, nil
}

// getLastTransactionHash returns logical time and hash of the last transaction
func (c *Client) getLastTransactionHash(ctx context.Context, acc ton.AccountID) (uint64, tlb.Bits256, error) {
	state, err := c.GetAccountState(ctx, acc)
	if err != nil {
		return 0, tlb.Bits256{}, errors.Wrap(err, "unable to get account state")
	}

	if state.Account.Status() != tlb.AccountActive {
		return 0, tlb.Bits256{}, errors.New("account is not active")
	}

	return state.LastTransLt, state.LastTransHash, nil
}

func TransactionToHashString(tx *ton.Transaction) string {
	return TransactionHashToString(tx.Lt, ton.Bits256(tx.Hash()))
}

// TransactionHashToString converts logicalTime and hash to string
func TransactionHashToString(lt uint64, hash ton.Bits256) string {
	return fmt.Sprintf("%d:%s", lt, hash.Hex())
}

// TransactionHashFromString parses encoded string into logicalTime and hash
func TransactionHashFromString(encoded string) (uint64, ton.Bits256, error) {
	parts := strings.Split(encoded, ":")
	if len(parts) != 2 {
		return 0, ton.Bits256{}, fmt.Errorf("invalid encoded string format")
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
