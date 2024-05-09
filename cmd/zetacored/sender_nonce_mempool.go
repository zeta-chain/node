// This is fork of SenderNonceMempool from cosmos sdk 0.47 (check: https://github.com/cosmos/cosmos-sdk/blob/v0.47.10/types/mempool/sender_nonce.go)
// only change is part where signatures are checked

// this is just for illustration, if we go with similar approach we would use priority nonce mempool, which has same issue but is more complex
// so testing with this one for now

package main

import (
	"context"
	crand "crypto/rand" // #nosec // crypto/rand is used for seed generation
	"encoding/binary"
	"fmt"
	"math/rand" // #nosec // math/rand is used for random selection and seeded from crypto/rand

	"github.com/huandu/skiplist"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	_ mempool.Mempool  = (*SenderNonceMempool)(nil)
	_ mempool.Iterator = (*senderNonceMempoolIterator)(nil)
)

var DefaultMaxTx = 0

// SenderNonceMempool is a mempool that prioritizes transactions within a sender
// by nonce, the lowest first, but selects a random sender on each iteration.
// The mempool is iterated by:
//
// 1) Maintaining a separate list of nonce ordered txs per sender
// 2) For each select iteration, randomly choose a sender and pick the next nonce ordered tx from their list
// 3) Repeat 1,2 until the mempool is exhausted
//
// Note that PrepareProposal could choose to stop iteration before reaching the
// end if maxBytes is reached.
type SenderNonceMempool struct {
	senders    map[string]*skiplist.SkipList
	rnd        *rand.Rand
	maxTx      int
	existingTx map[txKey]bool
}

type SenderNonceOptions func(mp *SenderNonceMempool)

type txKey struct {
	address string
	nonce   uint64
}

// NewSenderNonceMempool creates a new mempool that prioritizes transactions by
// nonce, the lowest first, picking a random sender on each iteration.
func NewSenderNonceMempool(opts ...SenderNonceOptions) *SenderNonceMempool {
	senderMap := make(map[string]*skiplist.SkipList)
	existingTx := make(map[txKey]bool)
	snp := &SenderNonceMempool{
		senders:    senderMap,
		maxTx:      DefaultMaxTx,
		existingTx: existingTx,
	}

	var seed int64
	err := binary.Read(crand.Reader, binary.BigEndian, &seed)
	if err != nil {
		panic(err)
	}

	snp.setSeed(seed)

	for _, opt := range opts {
		opt(snp)
	}

	return snp
}

// SenderNonceSeedOpt Option To add a Seed for random type when calling the
// constructor NewSenderNonceMempool.
//
// Example:
//
//	random_seed := int64(1000)
//	NewSenderNonceMempool(SenderNonceSeedTxOpt(random_seed))
func SenderNonceSeedOpt(seed int64) SenderNonceOptions {
	return func(snp *SenderNonceMempool) {
		snp.setSeed(seed)
	}
}

// SenderNonceMaxTxOpt Option To set limit of max tx when calling the constructor
// NewSenderNonceMempool.
//
// Example:
//
//	NewSenderNonceMempool(SenderNonceMaxTxOpt(100))
func SenderNonceMaxTxOpt(maxTx int) SenderNonceOptions {
	return func(snp *SenderNonceMempool) {
		snp.maxTx = maxTx
	}
}

func (snm *SenderNonceMempool) setSeed(seed int64) {
	s1 := rand.NewSource(seed)
	snm.rnd = rand.New(s1) //#nosec // math/rand is seeded from crypto/rand by default
}

// NextSenderTx returns the next transaction for a given sender by nonce order,
// i.e. the next valid transaction for the sender. If no such transaction exists,
// nil will be returned.
func (mp *SenderNonceMempool) NextSenderTx(sender string) sdk.Tx {
	senderIndex, ok := mp.senders[sender]
	if !ok {
		return nil
	}

	cursor := senderIndex.Front()
	return cursor.Value.(sdk.Tx)
}

// Insert adds a tx to the mempool. It returns an error if the tx does not have
// at least one signer. Note, priority is ignored.
func (snm *SenderNonceMempool) Insert(_ context.Context, tx sdk.Tx) error {
	if snm.maxTx > 0 && snm.CountTx() >= snm.maxTx {
		return mempool.ErrMempoolTxMaxCapacity
	}
	if snm.maxTx < 0 {
		return nil
	}

	sendersWithNonce, err := GetSendersWithNonce(tx)
	if err != nil {
		return err
	}

	sender := sendersWithNonce[0].Sender
	nonce := sendersWithNonce[0].Nonce

	senderTxs, found := snm.senders[sender]
	if !found {
		senderTxs = skiplist.New(skiplist.Uint64)
		snm.senders[sender] = senderTxs
	}

	senderTxs.Set(nonce, tx)

	key := txKey{nonce: nonce, address: sender}
	snm.existingTx[key] = true

	return nil
}

// Select returns an iterator ordering transactions the mempool with the lowest
// nonce of a random selected sender first.
//
// NOTE: It is not safe to use this iterator while removing transactions from
// the underlying mempool.
func (snm *SenderNonceMempool) Select(_ context.Context, _ [][]byte) mempool.Iterator {
	var senders []string

	senderCursors := make(map[string]*skiplist.Element)
	orderedSenders := skiplist.New(skiplist.String)

	// #nosec
	for s := range snm.senders {
		orderedSenders.Set(s, s)
	}

	s := orderedSenders.Front()
	for s != nil {
		sender := s.Value.(string)
		senders = append(senders, sender)
		senderCursors[sender] = snm.senders[sender].Front()
		s = s.Next()
	}

	iter := &senderNonceMempoolIterator{
		senders:       senders,
		rnd:           snm.rnd,
		senderCursors: senderCursors,
	}

	return iter.Next()
}

// CountTx returns the total count of txs in the mempool.
func (snm *SenderNonceMempool) CountTx() int {
	return len(snm.existingTx)
}

// Remove removes a tx from the mempool. It returns an error if the tx does not
// have at least one signer or the tx was not found in the pool.
func (snm *SenderNonceMempool) Remove(tx sdk.Tx) error {
	sendersWithNonce, err := GetSendersWithNonce(tx)
	if err != nil {
		return err
	}

	sender := sendersWithNonce[0].Sender
	nonce := sendersWithNonce[0].Nonce

	senderTxs, found := snm.senders[sender]

	if !found {
		return mempool.ErrTxNotFound
	}

	res := senderTxs.Remove(nonce)

	if res == nil {
		return mempool.ErrTxNotFound
	}

	if senderTxs.Len() == 0 {
		delete(snm.senders, sender)
	}

	key := txKey{nonce: nonce, address: sender}
	delete(snm.existingTx, key)

	return nil
}

type senderNonceMempoolIterator struct {
	rnd           *rand.Rand
	currentTx     *skiplist.Element
	senders       []string
	senderCursors map[string]*skiplist.Element
}

// Next returns the next iterator state which will contain a tx with the next
// smallest nonce of a randomly selected sender.
func (i *senderNonceMempoolIterator) Next() mempool.Iterator {
	for len(i.senders) > 0 {
		senderIndex := i.rnd.Intn(len(i.senders))
		sender := i.senders[senderIndex]
		senderCursor, found := i.senderCursors[sender]
		if !found {
			i.senders = removeAtIndex(i.senders, senderIndex)
			continue
		}

		if nextCursor := senderCursor.Next(); nextCursor != nil {
			i.senderCursors[sender] = nextCursor
		} else {
			i.senders = removeAtIndex(i.senders, senderIndex)
		}

		return &senderNonceMempoolIterator{
			senders:       i.senders,
			currentTx:     senderCursor,
			rnd:           i.rnd,
			senderCursors: i.senderCursors,
		}
	}

	return nil
}

func (i *senderNonceMempoolIterator) Tx() sdk.Tx {
	return i.currentTx.Value.(sdk.Tx)
}

func removeAtIndex[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func GetSendersWithNonce(tx sdk.Tx) ([]SenderWithNonce, error) {
	if txWithExtensions, ok := tx.(authante.HasExtensionOptionsTx); ok {
		opts := txWithExtensions.GetExtensionOptions()
		if len(opts) > 0 && opts[0].GetTypeUrl() == "/ethermint.evm.v1.ExtensionOptionsEthereumTx" {
			for _, msg := range tx.GetMsgs() {
				if ethMsg, ok := msg.(*evmtypes.MsgEthereumTx); ok {

					return []SenderWithNonce{
						{
							Sender: ethMsg.GetFrom().String(),
							Nonce:  ethMsg.AsTransaction().Nonce(),
						},
					}, nil
				}
			}
		}
	}

	return getSendersWithNonceDefault(tx)
}

type SenderWithNonce struct {
	Sender string
	Nonce  uint64
}

func getSendersWithNonceDefault(tx sdk.Tx) ([]SenderWithNonce, error) {
	sendersWithNonce := []SenderWithNonce{}

	sigTx, ok := tx.(signing.SigVerifiableTx)
	if !ok {
		return nil, fmt.Errorf("tx of type %T does not implement SigVerifiableTx", tx)
	}

	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return nil, err
	}

	if len(sigs) == 0 {
		return nil, fmt.Errorf("tx must have at least one signer")
	}

	for _, sig := range sigs {
		sendersWithNonce = append(sendersWithNonce, SenderWithNonce{
			Sender: sig.PubKey.Address().String(),
			Nonce:  sig.Sequence,
		})
	}

	return sendersWithNonce, nil
}
