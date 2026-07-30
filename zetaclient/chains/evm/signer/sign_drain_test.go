//go:build drain

package signer

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

// TestNewTxDeterministicDigest guards the core drain invariant: identical EVMTx fields must
// produce an identical unsigned tx digest, so every node forms the same go-tss ceremony.
func TestNewTxDeterministicDigest(t *testing.T) {
	chainID := big.NewInt(1)
	to := ethcommon.HexToAddress("0x1111111111111111111111111111111111111111")
	gas := Gas{Limit: 21000, Price: big.NewInt(250000), PriorityFee: big.NewInt(0)}
	signer := ethtypes.LatestSignerForChainID(chainID)

	tx1, err := newTx(chainID, nil, to, big.NewInt(1000), gas, 5)
	require.NoError(t, err)
	tx2, err := newTx(chainID, nil, to, big.NewInt(1000), gas, 5)
	require.NoError(t, err)
	require.Equal(t, signer.Hash(tx1), signer.Hash(tx2), "identical inputs must yield identical digest")

	// any divergent field must change the digest
	tx3, err := newTx(chainID, nil, to, big.NewInt(1000), gas, 6)
	require.NoError(t, err)
	require.NotEqual(t, signer.Hash(tx1), signer.Hash(tx3), "different nonce must change digest")
}

// TestPendingNonceIncludesMempool guards the drain nonce guard: PendingNonce must report the
// pending count, so an outbound broadcast at the pinned nonce but not yet mined is visible and
// the poller declines instead of evicting it.
func TestPendingNonceIncludesMempool(t *testing.T) {
	// ARRANGE: one outbound sits unmined in the mempool, so pending is one ahead of confirmed
	const confirmed, pending = 7, 8

	ts := newTestSuite(t)
	ts.evmServer.On("eth_getTransactionCount", func(params map[string]any) (any, error) {
		if params["1"] == "pending" {
			return fmt.Sprintf("0x%x", pending), nil
		}
		return fmt.Sprintf("0x%x", confirmed), nil
	})

	// ACT
	nonce, err := ts.PendingNonce(context.Background())

	// ASSERT
	require.NoError(t, err)
	require.EqualValues(t, pending, nonce, "PendingNonce must read the mempool, not latest-confirmed")
}
