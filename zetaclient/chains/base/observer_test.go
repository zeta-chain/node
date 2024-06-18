package base_test

import (
	"os"
	"testing"

	lru "github.com/hashicorp/golang-lru"
	"github.com/rs/zerolog/log"
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
func createObserver(t *testing.T, dbPath string) *base.Observer {
	// constructor parameters
	chain := chains.Ethereum
	chainParams := *sample.ChainParams(chain.ChainId)
	zetacoreContext := context.NewZetacoreContext(config.NewConfig())
	zetacoreClient := mocks.NewMockZetacoreClient()
	tss := mocks.NewTSSMainnet()

	// create observer
	logger := base.DefaultLogger()
	ob, err := base.NewObserver(
		chain,
		chainParams,
		zetacoreContext,
		zetacoreClient,
		tss,
		base.DefaultBlockCacheSize,
		base.DefaultHeadersCacheSize,
		dbPath,
		nil,
		logger,
	)
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
	headersCacheSize := base.DefaultHeadersCacheSize
	dbPath := createTempDir(t)

	// test cases
	tests := []struct {
		name             string
		chain            chains.Chain
		chainParams      observertypes.ChainParams
		zetacoreContext  *context.ZetacoreContext
		zetacoreClient   interfaces.ZetacoreClient
		tss              interfaces.TSSSigner
		blockCacheSize   int
		headersCacheSize int
		dbPath           string
		fail             bool
		message          string
	}{
		{
			name:             "should be able to create new observer",
			chain:            chain,
			chainParams:      chainParams,
			zetacoreContext:  zetacoreContext,
			zetacoreClient:   zetacoreClient,
			tss:              tss,
			blockCacheSize:   blockCacheSize,
			headersCacheSize: headersCacheSize,
			dbPath:           dbPath,
			fail:             false,
		},
		{
			name:             "should return error on invalid block cache size",
			chain:            chain,
			chainParams:      chainParams,
			zetacoreContext:  zetacoreContext,
			zetacoreClient:   zetacoreClient,
			tss:              tss,
			blockCacheSize:   0,
			headersCacheSize: headersCacheSize,
			dbPath:           dbPath,
			fail:             true,
			message:          "error creating block cache",
		},
		{
			name:             "should return error on invalid header cache size",
			chain:            chain,
			chainParams:      chainParams,
			zetacoreContext:  zetacoreContext,
			zetacoreClient:   zetacoreClient,
			tss:              tss,
			blockCacheSize:   blockCacheSize,
			headersCacheSize: 0,
			dbPath:           dbPath,
			fail:             true,
			message:          "error creating header cache",
		},
		{
			name:             "should return error on invalid db path",
			chain:            chain,
			chainParams:      chainParams,
			zetacoreContext:  zetacoreContext,
			zetacoreClient:   zetacoreClient,
			tss:              tss,
			blockCacheSize:   blockCacheSize,
			headersCacheSize: headersCacheSize,
			dbPath:           "/invalid/123db",
			fail:             true,
			message:          "error opening observer db",
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
				tt.headersCacheSize,
				tt.dbPath,
				nil,
				base.DefaultLogger(),
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

func TestObserverGetterAndSetter(t *testing.T) {
	dbPath := createTempDir(t)

	t.Run("should be able to update chain", func(t *testing.T) {
		ob := createObserver(t, dbPath)

		// update chain
		newChain := chains.BscMainnet
		ob = ob.WithChain(chains.BscMainnet)
		require.Equal(t, newChain, ob.Chain())
	})
	t.Run("should be able to update chain params", func(t *testing.T) {
		ob := createObserver(t, dbPath)

		// update chain params
		newChainParams := *sample.ChainParams(chains.BscMainnet.ChainId)
		ob = ob.WithChainParams(newChainParams)
		require.True(t, observertypes.ChainParamsEqual(newChainParams, ob.ChainParams()))
	})
	t.Run("should be able to update zetacore context", func(t *testing.T) {
		ob := createObserver(t, dbPath)

		// update zetacore context
		newZetacoreContext := context.NewZetacoreContext(config.NewConfig())
		ob = ob.WithZetacoreContext(newZetacoreContext)
		require.Equal(t, newZetacoreContext, ob.ZetacoreContext())
	})
	t.Run("should be able to update zetacore client", func(t *testing.T) {
		ob := createObserver(t, dbPath)

		// update zetacore client
		newZetacoreClient := mocks.NewMockZetacoreClient()
		ob = ob.WithZetacoreClient(newZetacoreClient)
		require.Equal(t, newZetacoreClient, ob.ZetacoreClient())
	})
	t.Run("should be able to update tss", func(t *testing.T) {
		ob := createObserver(t, dbPath)

		// update tss
		newTSS := mocks.NewTSSAthens3()
		ob = ob.WithTSS(newTSS)
		require.Equal(t, newTSS, ob.TSS())
	})
	t.Run("should be able to update last block", func(t *testing.T) {
		ob := createObserver(t, dbPath)

		// update last block
		newLastBlock := uint64(100)
		ob = ob.WithLastBlock(newLastBlock)
		require.Equal(t, newLastBlock, ob.LastBlock())
	})
	t.Run("should be able to update last block scanned", func(t *testing.T) {
		ob := createObserver(t, dbPath)

		// update last block scanned
		newLastBlockScanned := uint64(100)
		ob = ob.WithLastBlockScanned(newLastBlockScanned)
		require.Equal(t, newLastBlockScanned, ob.LastBlockScanned())
	})
	t.Run("should be able to replace block cache", func(t *testing.T) {
		ob := createObserver(t, dbPath)

		// update block cache
		newBlockCache, err := lru.New(200)
		require.NoError(t, err)

		ob = ob.WithBlockCache(newBlockCache)
		require.Equal(t, newBlockCache, ob.BlockCache())
	})
	t.Run("should be able to replace headers cache", func(t *testing.T) {
		ob := createObserver(t, dbPath)

		// update headers cache
		newHeadersCache, err := lru.New(200)
		require.NoError(t, err)

		ob = ob.WithHeaderCache(newHeadersCache)
		require.Equal(t, newHeadersCache, ob.HeaderCache())
	})
	t.Run("should be able to get database", func(t *testing.T) {
		ob := createObserver(t, dbPath)

		db := ob.DB()
		require.NotNil(t, db)
	})
	t.Run("should be able to get logger", func(t *testing.T) {
		ob := createObserver(t, dbPath)
		logger := ob.Logger()

		// should be able to print log
		logger.Chain.Info().Msg("print chain log")
		logger.Inbound.Info().Msg("print inbound log")
		logger.Outbound.Info().Msg("print outbound log")
		logger.GasPrice.Info().Msg("print gasprice log")
		logger.Headers.Info().Msg("print headers log")
		logger.Compliance.Info().Msg("print compliance log")
	})
}

func TestOpenDB(t *testing.T) {
	dbPath := createTempDir(t)
	ob := createObserver(t, dbPath)

	t.Run("should be able to open db", func(t *testing.T) {
		err := ob.OpenDB(dbPath)
		require.NoError(t, err)
	})
	t.Run("should return error on invalid db path", func(t *testing.T) {
		err := ob.OpenDB("/invalid/123db")
		require.ErrorContains(t, err, "error creating db path")
	})
}

func TestLoadLastBlockScanned(t *testing.T) {
	chain := chains.Ethereum
	envvar := base.EnvVarLatestBlockByChain(chain)

	t.Run("should be able to load last block scanned", func(t *testing.T) {
		// create db and write 100 as last block scanned
		dbPath := createTempDir(t)
		ob := createObserver(t, dbPath)
		ob.WriteLastBlockScannedToDB(100)

		// read last block scanned
		fromLatest, err := ob.LoadLastBlockScanned(log.Logger)
		require.NoError(t, err)
		require.EqualValues(t, 100, ob.LastBlockScanned())
		require.False(t, fromLatest)
	})
	t.Run("should use latest block if last block scanned not found", func(t *testing.T) {
		// create empty db
		dbPath := createTempDir(t)
		ob := createObserver(t, dbPath)

		// read last block scanned
		fromLatest, err := ob.LoadLastBlockScanned(log.Logger)
		require.NoError(t, err)
		require.True(t, fromLatest)
	})
	t.Run("should overwrite last block scanned if env var is set", func(t *testing.T) {
		// create db and write 100 as last block scanned
		dbPath := createTempDir(t)
		ob := createObserver(t, dbPath)
		ob.WriteLastBlockScannedToDB(100)

		// set env var
		os.Setenv(envvar, "101")

		// read last block scanned
		fromLatest, err := ob.LoadLastBlockScanned(log.Logger)
		require.NoError(t, err)
		require.EqualValues(t, 101, ob.LastBlockScanned())
		require.False(t, fromLatest)

		// set env var to 'latest'
		os.Setenv(envvar, base.EnvVarLatestBlock)

		// read last block scanned
		fromLatest, err = ob.LoadLastBlockScanned(log.Logger)
		require.NoError(t, err)
		require.True(t, fromLatest)
	})
	t.Run("should return error on invalid env var", func(t *testing.T) {
		// create db and write 100 as last block scanned
		dbPath := createTempDir(t)
		ob := createObserver(t, dbPath)

		// set invalid env var
		os.Setenv(envvar, "invalid")

		// read last block scanned
		fromLatest, err := ob.LoadLastBlockScanned(log.Logger)
		require.Error(t, err)
		require.False(t, fromLatest)
	})
}

func TestReadWriteLastBlockScannedToDB(t *testing.T) {
	t.Run("should be able to write and read last block scanned to db", func(t *testing.T) {
		// create db and write 100 as last block scanned
		dbPath := createTempDir(t)
		ob := createObserver(t, dbPath)
		err := ob.WriteLastBlockScannedToDB(100)
		require.NoError(t, err)

		lastBlockScanned, err := ob.ReadLastBlockScannedFromDB()
		require.NoError(t, err)
		require.EqualValues(t, 100, lastBlockScanned)
	})
	t.Run("should return error when last block scanned not found in db", func(t *testing.T) {
		// create empty db
		dbPath := createTempDir(t)
		ob := createObserver(t, dbPath)

		lastScannedBlock, err := ob.ReadLastBlockScannedFromDB()
		require.Error(t, err)
		require.Zero(t, lastScannedBlock)
	})
}
