package base_test

import (
	"os"
	"testing"

	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

// create a temporary directory for testing
func createTempDir(t *testing.T) string {
	tempPath, err := os.MkdirTemp("", "tempdir-")
	require.NoError(t, err)
	return tempPath
}

// createObserver creates a new observer for testing
func createObserver(t *testing.T) *base.Observer {
	// constructor parameters
	chain := chains.Ethereum
	chainParams := *sample.ChainParams(chain.ChainId)
	zetacoreContext := context.NewZetacoreContext(config.NewConfig())
	zetacoreClient := mocks.NewMockZetacoreClient()
	tss := mocks.NewTSSMainnet()
	blockCacheSize := base.DefaultBlockCacheSize
	dbPath := createTempDir(t)

	// create observer
	ob, err := base.NewObserver(chain, chainParams, zetacoreContext, zetacoreClient, tss, blockCacheSize, dbPath, nil)
	require.NoError(t, err)

	return ob
}

func TestNewObserver(t *testing.T) {
	// constructor parameters
	chain := chains.Ethereum
	chainParams := *sample.ChainParams(chain.ChainId)
	zetacoreContext := context.NewZetacoreContext(config.NewConfig())
	zetacoreClient := mocks.NewMockZetacoreClient()
	tss := mocks.NewTSSMainnet()
	blockCacheSize := base.DefaultBlockCacheSize
	dbPath := createTempDir(t)

	// test cases
	tests := []struct {
		name            string
		chain           chains.Chain
		chainParams     observertypes.ChainParams
		zetacoreContext *context.ZetacoreContext
		zetacoreClient  interfaces.ZetacoreClient
		tss             interfaces.TSSSigner
		blockCacheSize  int
		dbPath          string
		fail            bool
		message         string
	}{
		{
			name:            "should be able to create new observer",
			chain:           chain,
			chainParams:     chainParams,
			zetacoreContext: zetacoreContext,
			zetacoreClient:  zetacoreClient,
			tss:             tss,
			blockCacheSize:  blockCacheSize,
			dbPath:          dbPath,
			fail:            false,
		},
		{
			name:            "should return error on invalid block cache size",
			chain:           chain,
			chainParams:     chainParams,
			zetacoreContext: zetacoreContext,
			zetacoreClient:  zetacoreClient,
			tss:             tss,
			blockCacheSize:  0,
			dbPath:          dbPath,
			fail:            true,
			message:         "error creating block cache",
		},
		{
			name:            "should return error on invalid db path",
			chain:           chain,
			chainParams:     chainParams,
			zetacoreContext: zetacoreContext,
			zetacoreClient:  zetacoreClient,
			tss:             tss,
			blockCacheSize:  blockCacheSize,
			dbPath:          "/invalid/123db",
			fail:            true,
			message:         "error opening observer db",
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob, err := base.NewObserver(
				tt.chain,
				tt.chainParams,
				tt.zetacoreContext,
				tt.zetacoreClient,
				tt.tss,
				tt.blockCacheSize,
				tt.dbPath,
				nil,
			)
			if tt.fail {
				require.ErrorContains(t, err, tt.message)
				require.Nil(t, ob)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, ob)
		})
	}
}

func TestObserverSetters(t *testing.T) {
	t.Run("should be able to update chain", func(t *testing.T) {
		ob := createObserver(t)

		// update chain
		newChain := chains.BscMainnet
		ob = ob.WithChain(chains.BscMainnet)
		require.Equal(t, newChain, ob.Chain())
	})
	t.Run("should be able to update chain params", func(t *testing.T) {
		ob := createObserver(t)

		// update chain params
		newChainParams := *sample.ChainParams(chains.BscMainnet.ChainId)
		ob = ob.WithChainParams(newChainParams)
		require.True(t, observertypes.ChainParamsEqual(newChainParams, ob.ChainParams()))
	})
	t.Run("should be able to update zetacore client", func(t *testing.T) {
		ob := createObserver(t)

		// update zetacore client
		newZetacoreClient := mocks.NewMockZetacoreClient()
		ob = ob.WithZetacoreClient(newZetacoreClient)
		require.Equal(t, newZetacoreClient, ob.ZetacoreClient())
	})
	t.Run("should be able to update last block", func(t *testing.T) {
		ob := createObserver(t)

		// update last block
		newLastBlock := uint64(100)
		ob = ob.WithLastBlock(newLastBlock)
		require.Equal(t, newLastBlock, ob.LastBlock())
	})
	t.Run("should be able to update last block scanned", func(t *testing.T) {
		ob := createObserver(t)

		// update last block scanned
		newLastBlockScanned := uint64(100)
		ob = ob.WithLastBlockScanned(newLastBlockScanned)
		require.Equal(t, newLastBlockScanned, ob.LastBlockScanned())
	})
	t.Run("should be able to update block cache", func(t *testing.T) {
		ob := createObserver(t)

		// update block cache
		newBlockCache, err := lru.New(200)
		require.NoError(t, err)

		ob = ob.WithBlockCache(newBlockCache)
		require.Equal(t, newBlockCache, ob.BlockCache())
	})
}

func TestOpenDB(t *testing.T) {
	ob := createObserver(t)
	dbPath := createTempDir(t)

	t.Run("should be able to open db", func(t *testing.T) {
		err := ob.OpenDB(dbPath)
		require.NoError(t, err)
	})
	t.Run("should return error on invalid db path", func(t *testing.T) {
		err := ob.OpenDB("/invalid/123db")
		require.ErrorContains(t, err, "error creating db path")
	})
}

func TestReadWriteLastBlockScannedToDB(t *testing.T) {
	t.Run("should be able to write and read last block scanned to db", func(t *testing.T) {
		ob := createObserver(t)
		err := ob.WriteLastBlockScannedToDB(100)
		require.NoError(t, err)

		lastBlockScanned, err := ob.ReadLastBlockScannedFromDB()
		require.NoError(t, err)
		require.EqualValues(t, 100, lastBlockScanned)
	})
	t.Run("should return error when last block scanned not found in db", func(t *testing.T) {
		ob := createObserver(t)
		lastScannedBlock, err := ob.ReadLastBlockScannedFromDB()
		require.Error(t, err)
		require.Zero(t, lastScannedBlock)
	})
}
