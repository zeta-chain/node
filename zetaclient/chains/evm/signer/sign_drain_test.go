//go:build drain

package signer

import (
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
