package evm

import (
	"math/big"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/require"
)

func TestEVMBlockCache(t *testing.T) {
	// create client
	blockCache, err := lru.New(1000)
	require.NoError(t, err)
	blockCacheV3, err := lru.New(1000)
	require.NoError(t, err)
	ob := ChainClient{
		blockCache:   blockCache,
		blockCacheV3: blockCacheV3,
	}

	// delete non-existing block should not panic
	blockNumber := int64(10388180)
	// #nosec G701 possible nummber
	ob.RemoveCachedBlock(uint64(blockNumber))

	// add a block
	header := &ethtypes.Header{
		Number: big.NewInt(blockNumber),
	}
	block := ethtypes.NewBlock(header, nil, nil, nil, nil)
	ob.blockCache.Add(blockNumber, block)

	// delete the block should not panic
	ob.RemoveCachedBlock(uint64(blockNumber))
}
