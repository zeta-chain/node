package base_test

import (
	"os"
	"testing"

	lru "github.com/hashicorp/golang-lru"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

// createObserver creates a new observer for testing
func createObserver(t *testing.T) *base.Observer {
	// constructor parameters
	chain := chains.Ethereum
	chainParams := *sample.ChainParams(chain.ChainId)
	zetacoreClient := mocks.NewZetacoreClient(t)
	tss := mocks.NewTSSMainnet()

	// create observer
	logger := base.DefaultLogger()
	ob, err := base.NewObserver(
		chain,
		chainParams,
		zetacoreClient,
		tss,
		base.DefaultBlockCacheSize,
		base.DefaultHeaderCacheSize,
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
	appContext := context.New(config.New(false), zerolog.Nop())
	zetacoreClient := mocks.NewZetacoreClient(t)
	tss := mocks.NewTSSMainnet()
	blockCacheSize := base.DefaultBlockCacheSize
	headersCacheSize := base.DefaultHeaderCacheSize

	// test cases
	tests := []struct {
		name            string
		chain           chains.Chain
		chainParams     observertypes.ChainParams
		appContext      *context.AppContext
		zetacoreClient  interfaces.ZetacoreClient
		tss             interfaces.TSSSigner
		blockCacheSize  int
		headerCacheSize int
		fail            bool
		message         string
	}{
		{
			name:            "should be able to create new observer",
			chain:           chain,
			chainParams:     chainParams,
			appContext:      appContext,
			zetacoreClient:  zetacoreClient,
			tss:             tss,
			blockCacheSize:  blockCacheSize,
			headerCacheSize: headersCacheSize,
			fail:            false,
		},
		{
			name:            "should return error on invalid block cache size",
			chain:           chain,
			chainParams:     chainParams,
			appContext:      appContext,
			zetacoreClient:  zetacoreClient,
			tss:             tss,
			blockCacheSize:  0,
			headerCacheSize: headersCacheSize,
			fail:            true,
			message:         "error creating block cache",
		},
		{
			name:            "should return error on invalid header cache size",
			chain:           chain,
			chainParams:     chainParams,
			appContext:      appContext,
			zetacoreClient:  zetacoreClient,
			tss:             tss,
			blockCacheSize:  blockCacheSize,
			headerCacheSize: 0,
			fail:            true,
			message:         "error creating header cache",
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob, err := base.NewObserver(
				tt.chain,
				tt.chainParams,
				tt.zetacoreClient,
				tt.tss,
				tt.blockCacheSize,
				tt.headerCacheSize,
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

func TestStop(t *testing.T) {
	t.Run("should be able to stop observer", func(t *testing.T) {
		// create observer and initialize db
		ob := createObserver(t)
		ob.OpenDB(sample.CreateTempDir(t), "")

		// stop observer
		ob.Stop()
	})
}

func TestObserverGetterAndSetter(t *testing.T) {
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
		newZetacoreClient := mocks.NewZetacoreClient(t)
		ob = ob.WithZetacoreClient(newZetacoreClient)
		require.Equal(t, newZetacoreClient, ob.ZetacoreClient())
	})
	t.Run("should be able to update tss", func(t *testing.T) {
		ob := createObserver(t)

		// update tss
		newTSS := mocks.NewTSSAthens3()
		ob = ob.WithTSS(newTSS)
		require.Equal(t, newTSS, ob.TSS())
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
	t.Run("should be able to replace block cache", func(t *testing.T) {
		ob := createObserver(t)

		// update block cache
		newBlockCache, err := lru.New(200)
		require.NoError(t, err)

		ob = ob.WithBlockCache(newBlockCache)
		require.Equal(t, newBlockCache, ob.BlockCache())
	})
	t.Run("should be able to replace header cache", func(t *testing.T) {
		ob := createObserver(t)

		// update headers cache
		newHeadersCache, err := lru.New(200)
		require.NoError(t, err)

		ob = ob.WithHeaderCache(newHeadersCache)
		require.Equal(t, newHeadersCache, ob.HeaderCache())
	})
	t.Run("should be able to get database", func(t *testing.T) {
		// create observer and open db
		dbPath := sample.CreateTempDir(t)
		ob := createObserver(t)
		ob.OpenDB(dbPath, "")

		db := ob.DB()
		require.NotNil(t, db)
	})
	t.Run("should be able to update telemetry server", func(t *testing.T) {
		ob := createObserver(t)

		// update telemetry server
		newServer := metrics.NewTelemetryServer()
		ob = ob.WithTelemetryServer(newServer)
		require.Equal(t, newServer, ob.TelemetryServer())
	})
	t.Run("should be able to get logger", func(t *testing.T) {
		ob := createObserver(t)
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

func TestOpenCloseDB(t *testing.T) {
	dbPath := sample.CreateTempDir(t)
	ob := createObserver(t)

	t.Run("should be able to open/close db", func(t *testing.T) {
		// open db
		err := ob.OpenDB(dbPath, "")
		require.NoError(t, err)

		// close db
		err = ob.CloseDB()
		require.NoError(t, err)
	})
	t.Run("should use memory db if specified", func(t *testing.T) {
		// open db with memory
		err := ob.OpenDB(testutils.SQLiteMemory, "")
		require.NoError(t, err)

		// close db
		err = ob.CloseDB()
		require.NoError(t, err)
	})
	t.Run("should return error on invalid db path", func(t *testing.T) {
		err := ob.OpenDB("/invalid/123db", "")
		require.ErrorContains(t, err, "error creating db path")
	})
}

func TestLoadLastBlockScanned(t *testing.T) {
	chain := chains.Ethereum
	envvar := base.EnvVarLatestBlockByChain(chain)

	t.Run("should be able to load last block scanned", func(t *testing.T) {
		// create observer and open db
		dbPath := sample.CreateTempDir(t)
		ob := createObserver(t)
		err := ob.OpenDB(dbPath, "")
		require.NoError(t, err)

		// create db and write 100 as last block scanned
		ob.WriteLastBlockScannedToDB(100)

		// read last block scanned
		err = ob.LoadLastBlockScanned(log.Logger)
		require.NoError(t, err)
		require.EqualValues(t, 100, ob.LastBlockScanned())
	})
	t.Run("latest block scanned should be 0 if not found in db", func(t *testing.T) {
		// create observer and open db
		dbPath := sample.CreateTempDir(t)
		ob := createObserver(t)
		err := ob.OpenDB(dbPath, "")
		require.NoError(t, err)

		// read last block scanned
		err = ob.LoadLastBlockScanned(log.Logger)
		require.NoError(t, err)
		require.EqualValues(t, 0, ob.LastBlockScanned())
	})
	t.Run("should overwrite last block scanned if env var is set", func(t *testing.T) {
		// create observer and open db
		dbPath := sample.CreateTempDir(t)
		ob := createObserver(t)
		err := ob.OpenDB(dbPath, "")
		require.NoError(t, err)

		// create db and write 100 as last block scanned
		ob.WriteLastBlockScannedToDB(100)

		// set env var
		os.Setenv(envvar, "101")

		// read last block scanned
		err = ob.LoadLastBlockScanned(log.Logger)
		require.NoError(t, err)
		require.EqualValues(t, 101, ob.LastBlockScanned())
	})
	t.Run("last block scanned should remain 0 if env var is set to latest", func(t *testing.T) {
		// create observer and open db
		dbPath := sample.CreateTempDir(t)
		ob := createObserver(t)
		err := ob.OpenDB(dbPath, "")
		require.NoError(t, err)

		// create db and write 100 as last block scanned
		ob.WriteLastBlockScannedToDB(100)

		// set env var to 'latest'
		os.Setenv(envvar, base.EnvVarLatestBlock)

		// last block scanned should remain 0
		err = ob.LoadLastBlockScanned(log.Logger)
		require.NoError(t, err)
		require.EqualValues(t, 0, ob.LastBlockScanned())
	})
	t.Run("should return error on invalid env var", func(t *testing.T) {
		// create observer and open db
		dbPath := sample.CreateTempDir(t)
		ob := createObserver(t)
		err := ob.OpenDB(dbPath, "")
		require.NoError(t, err)

		// set invalid env var
		os.Setenv(envvar, "invalid")

		// read last block scanned
		err = ob.LoadLastBlockScanned(log.Logger)
		require.Error(t, err)
	})
}

func TestSaveLastBlockScanned(t *testing.T) {
	t.Run("should be able to save last block scanned", func(t *testing.T) {
		// create observer and open db
		dbPath := sample.CreateTempDir(t)
		ob := createObserver(t)
		err := ob.OpenDB(dbPath, "")
		require.NoError(t, err)

		// save 100 as last block scanned
		err = ob.SaveLastBlockScanned(100)
		require.NoError(t, err)

		// check last block scanned in memory
		require.EqualValues(t, 100, ob.LastBlockScanned())

		// read last block scanned from db
		lastBlockScanned, err := ob.ReadLastBlockScannedFromDB()
		require.NoError(t, err)
		require.EqualValues(t, 100, lastBlockScanned)
	})
}

func TestReadWriteLastBlockScannedToDB(t *testing.T) {
	t.Run("should be able to write and read last block scanned to db", func(t *testing.T) {
		// create observer and open db
		dbPath := sample.CreateTempDir(t)
		ob := createObserver(t)
		err := ob.OpenDB(dbPath, "")
		require.NoError(t, err)

		// write last block scanned
		err = ob.WriteLastBlockScannedToDB(100)
		require.NoError(t, err)

		lastBlockScanned, err := ob.ReadLastBlockScannedFromDB()
		require.NoError(t, err)
		require.EqualValues(t, 100, lastBlockScanned)
	})
	t.Run("should return error when last block scanned not found in db", func(t *testing.T) {
		// create empty db
		dbPath := sample.CreateTempDir(t)
		ob := createObserver(t)
		err := ob.OpenDB(dbPath, "")
		require.NoError(t, err)

		lastScannedBlock, err := ob.ReadLastBlockScannedFromDB()
		require.Error(t, err)
		require.Zero(t, lastScannedBlock)
	})
}
