package base_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/mode"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// defaultConfirmationCount is the default confirmation count for unit tests
	defaultConfirmationCount = 2
)

type testSuite struct {
	*base.Observer
	db       *db.DB
	tss      *mocks.TSS
	zetacore *mocks.ZetacoreClient
}

type testSuiteOpts struct {
	ConfirmationParams *observertypes.ConfirmationParams
}

type opt func(t *testSuiteOpts)

// withConfirmationParams is an option to set custom confirmation params
func withConfirmationParams(confParams observertypes.ConfirmationParams) opt {
	return func(t *testSuiteOpts) {
		t.ConfirmationParams = &confParams
	}
}

// newTestSuite creates a new observer for testing
func newTestSuite(t *testing.T, chain chains.Chain, opts ...opt) *testSuite {
	// create test suite with options
	var testOpts testSuiteOpts
	for _, opt := range opts {
		opt(&testOpts)
	}

	// constructor parameters
	chainParams := *sample.ChainParams(chain.ChainId)
	chainParams.IsSupported = true
	chainParams.ConfirmationParams = &observertypes.ConfirmationParams{
		SafeInboundCount:  defaultConfirmationCount,
		SafeOutboundCount: defaultConfirmationCount,
	}
	if testOpts.ConfirmationParams != nil {
		chainParams.ConfirmationParams = testOpts.ConfirmationParams
	}
	zetacoreClient := mocks.NewZetacoreClient(t)
	tss := mocks.NewTSS(t)

	database := createDatabase(t)

	// create observer
	logger := base.DefaultLogger()
	ob, err := base.NewObserver(
		chain,
		chainParams,
		zetacoreClient,
		tss,
		base.DefaultBlockCacheSize,
		nil,
		database,
		logger,
		mode.StandardMode,
	)
	require.NoError(t, err)

	return &testSuite{
		Observer: ob,
		db:       database,
		tss:      tss,
		zetacore: zetacoreClient,
	}
}

func TestNewObserver(t *testing.T) {
	// constructor parameters
	chain := chains.Ethereum
	chainParams := *sample.ChainParams(chain.ChainId)
	appContext := zctx.New(config.New(false), nil, zerolog.Nop())
	zetacoreClient := mocks.NewZetacoreClient(t)
	tss := mocks.NewTSS(t)
	blockCacheSize := base.DefaultBlockCacheSize

	database := createDatabase(t)

	// test cases
	tests := []struct {
		name           string
		chain          chains.Chain
		chainParams    observertypes.ChainParams
		appContext     *zctx.AppContext
		zetacoreClient zrepo.ZetacoreClient
		tssSigner      interfaces.TSSSigner
		blockCacheSize int
		fail           bool
		message        string
	}{
		{
			name:           "should be able to create new observer",
			chain:          chain,
			chainParams:    chainParams,
			appContext:     appContext,
			zetacoreClient: zetacoreClient,
			tssSigner:      tss,
			blockCacheSize: blockCacheSize,
			fail:           false,
		},
		{
			name:           "should return error on invalid block cache size",
			chain:          chain,
			chainParams:    chainParams,
			appContext:     appContext,
			zetacoreClient: zetacoreClient,
			tssSigner:      tss,
			blockCacheSize: 0,
			fail:           true,
			message:        "error creating block cache",
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob, err := base.NewObserver(
				tt.chain,
				tt.chainParams,
				tt.zetacoreClient,
				tt.tssSigner,
				tt.blockCacheSize,
				nil,
				database,
				base.DefaultLogger(),
				mode.StandardMode,
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
		ob := newTestSuite(t, chains.Ethereum)

		// stop observer
		ob.Stop()
	})
}

func TestObserverGetterAndSetter(t *testing.T) {
	chain := chains.Ethereum

	t.Run("should be able to update last block", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// update last block
		newLastBlock := uint64(100)
		ob.Observer.WithLastBlock(newLastBlock)
		require.Equal(t, newLastBlock, ob.LastBlock())
	})

	t.Run("should be able to update last block scanned", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// update last block scanned
		newLastBlockScanned := uint64(100)
		ob.Observer.WithLastBlockScanned(newLastBlockScanned)
		require.Equal(t, newLastBlockScanned, ob.LastBlockScanned())
	})

	t.Run("should be able to update last tx scanned", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// update last tx scanned
		newLastTxScanned := sample.EthAddress().String()
		ob.Observer.WithLastTxScanned(newLastTxScanned)
		require.Equal(t, newLastTxScanned, ob.LastTxScanned())
	})

	t.Run("should be able to get logger", func(t *testing.T) {
		ob := newTestSuite(t, chain)
		logger := ob.Logger()

		// should be able to print log
		logger.Chain.Info().Msg("print chain log")
		logger.Inbound.Info().Msg("print inbound log")
		logger.Outbound.Info().Msg("print outbound log")
		logger.Compliance.Info().Msg("print compliance log")
		logger.Sampled.Info().Msg("print sampled log")
	})
}

func TestTSSAddressString(t *testing.T) {
	btcSomething := chains.BitcoinMainnet
	btcSomething.ChainId = 123123123

	tests := []struct {
		name         string
		chain        chains.Chain
		addrExpected string
	}{
		{
			name:         "should return TSS BTC address for Bitcoin chain",
			chain:        chains.BitcoinMainnet,
			addrExpected: "btc",
		},
		{
			name:         "should return TSS EVM address for EVM chain",
			chain:        chains.Ethereum,
			addrExpected: "eth",
		},
		{
			name:         "should return TSS EVM address for other non-BTC chain",
			chain:        chains.SolanaDevnet,
			addrExpected: "eth",
		},
		{
			name:         "should return empty address for unknown BTC chain",
			chain:        btcSomething,
			addrExpected: "",
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create observer
			ob := newTestSuite(t, tt.chain)

			// get TSS address
			addr := ob.TSSAddressString()
			switch tt.addrExpected {
			case "":
				require.Equal(t, "", addr)
			case "btc":
				require.True(t, strings.HasPrefix(addr, "bc"))
			case "eth":
				require.True(t, strings.HasPrefix(addr, "0x"))
			default:
				t.Fail()
			}
		})
	}
}

func TestOutboundID(t *testing.T) {
	tests := []struct {
		name  string
		chain chains.Chain
		nonce uint64
	}{
		{
			name:  "should get correct outbound id for Ethereum chain",
			chain: chains.Ethereum,
			nonce: 100,
		},
		{
			name:  "should get correct outbound id for Bitcoin chain",
			chain: chains.BitcoinMainnet,
			nonce: 200,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create observer
			ob := newTestSuite(t, tt.chain)

			// get outbound id
			outboundID := ob.OutboundID(tt.nonce)

			// expected outbound id
			exepctedID := fmt.Sprintf("%d-%s-%d", tt.chain.ChainId, ob.TSSAddressString(), tt.nonce)
			require.Equal(t, exepctedID, outboundID)
		})
	}
}

func TestLoadLastBlockScanned(t *testing.T) {
	chain := chains.Ethereum
	envvar := base.EnvVarLatestBlockByChain(chain)

	t.Run("should be able to load last block scanned", func(t *testing.T) {
		// create observer and open db
		ob := newTestSuite(t, chain)

		// create db and write 100 as last block scanned
		err := ob.WriteLastBlockScannedToDB(100)
		require.NoError(t, err)

		// read last block scanned
		err = ob.LoadLastBlockScanned()
		require.NoError(t, err)
		require.EqualValues(t, 100, ob.LastBlockScanned())
	})

	t.Run("latest block scanned should be 0 if not found in db", func(t *testing.T) {
		// create observer and open db
		ob := newTestSuite(t, chain)

		// read last block scanned
		err := ob.LoadLastBlockScanned()
		require.NoError(t, err)
		require.EqualValues(t, 0, ob.LastBlockScanned())
	})

	t.Run("should overwrite last block scanned if env var is set", func(t *testing.T) {
		// create observer and open db
		ob := newTestSuite(t, chain)

		// create db and write 100 as last block scanned
		ob.WriteLastBlockScannedToDB(100)

		// set env var
		os.Setenv(envvar, "101")

		// read last block scanned
		err := ob.LoadLastBlockScanned()
		require.NoError(t, err)
		require.EqualValues(t, 101, ob.LastBlockScanned())
	})

	t.Run("last block scanned should remain 0 if env var is set to latest", func(t *testing.T) {
		// create observer and open db
		ob := newTestSuite(t, chain)

		// create db and write 100 as last block scanned
		ob.WriteLastBlockScannedToDB(100)

		// set env var to 'latest'
		os.Setenv(envvar, base.EnvVarLatestBlock)

		// last block scanned should remain 0
		err := ob.LoadLastBlockScanned()
		require.NoError(t, err)
		require.EqualValues(t, 0, ob.LastBlockScanned())
	})

	t.Run("should return error on invalid env var", func(t *testing.T) {
		// create observer and open db
		ob := newTestSuite(t, chain)

		// set invalid env var
		os.Setenv(envvar, "invalid")

		// read last block scanned
		err := ob.LoadLastBlockScanned()
		require.Error(t, err)
	})
}

func TestSaveLastBlockScanned(t *testing.T) {
	t.Run("should be able to save last block scanned", func(t *testing.T) {
		// create observer and open db
		ob := newTestSuite(t, chains.Ethereum)

		// save 100 as last block scanned
		err := ob.SaveLastBlockScanned(100)
		require.NoError(t, err)

		// check last block scanned in memory
		require.EqualValues(t, 100, ob.LastBlockScanned())

		// read last block scanned from db
		lastBlockScanned, err := ob.ReadLastBlockScannedFromDB()
		require.NoError(t, err)
		require.EqualValues(t, 100, lastBlockScanned)
	})
}

func TestReadWriteDBLastBlockScanned(t *testing.T) {
	chain := chains.Ethereum
	t.Run("should be able to write and read last block scanned to db", func(t *testing.T) {
		// create observer and open db
		ob := newTestSuite(t, chain)

		// write last block scanned
		err := ob.WriteLastBlockScannedToDB(100)
		require.NoError(t, err)

		lastBlockScanned, err := ob.ReadLastBlockScannedFromDB()
		require.NoError(t, err)
		require.EqualValues(t, 100, lastBlockScanned)
	})

	t.Run("should return error when last block scanned not found in db", func(t *testing.T) {
		// create empty db
		ob := newTestSuite(t, chain)

		lastScannedBlock, err := ob.ReadLastBlockScannedFromDB()
		require.Error(t, err)
		require.Zero(t, lastScannedBlock)
	})
}
func TestLoadLastTxScanned(t *testing.T) {
	chain := chains.SolanaDevnet
	envvar := base.EnvVarLatestTxByChain(chain)
	lastTx := "5LuQMorgd11p8GWEw6pmyHCDtA26NUyeNFhLWPNk2oBoM9pkag1LzhwGSRos3j4TJLhKjswFhZkGtvSGdLDkmqsk"

	t.Run("should be able to load last tx scanned", func(t *testing.T) {
		// create observer and open db
		ob := newTestSuite(t, chain)

		// create db and write sample hash as last tx scanned
		ob.WriteLastTxScannedToDB(lastTx)

		// read last tx scanned
		ob.LoadLastTxScanned()
		require.EqualValues(t, lastTx, ob.LastTxScanned())
	})

	t.Run("latest tx scanned should be empty if not found in db", func(t *testing.T) {
		// create observer and open db
		ob := newTestSuite(t, chain)

		// read last tx scanned
		ob.LoadLastTxScanned()
		require.Empty(t, ob.LastTxScanned())
	})

	t.Run("should overwrite last tx scanned if env var is set", func(t *testing.T) {
		// create observer and open db
		ob := newTestSuite(t, chain)

		// create db and write sample hash as last tx scanned
		ob.WriteLastTxScannedToDB(lastTx)

		// set env var to other tx
		otherTx := "4Q27KQqJU1gJQavNtkvhH6cGR14fZoBdzqWdWiFd9KPeJxFpYsDRiKAwsQDpKMPtyRhppdncyURTPZyokrFiVHrx"
		os.Setenv(envvar, otherTx)

		// read last block scanned
		ob.LoadLastTxScanned()
		require.EqualValues(t, otherTx, ob.LastTxScanned())
	})
}

func TestSaveLastTxScanned(t *testing.T) {
	chain := chains.SolanaDevnet
	t.Run("should be able to save last tx scanned", func(t *testing.T) {
		// create observer and open db
		ob := newTestSuite(t, chain)

		// save random tx hash
		lastSlot := uint64(100)
		lastTx := "5LuQMorgd11p8GWEw6pmyHCDtA26NUyeNFhLWPNk2oBoM9pkag1LzhwGSRos3j4TJLhKjswFhZkGtvSGdLDkmqsk"
		err := ob.SaveLastTxScanned(lastTx, lastSlot)
		require.NoError(t, err)

		// check last tx and slot scanned in memory
		require.EqualValues(t, lastTx, ob.LastTxScanned())
		require.EqualValues(t, lastSlot, ob.LastBlockScanned())

		// read last tx scanned from db
		lastTxScanned, err := ob.ReadLastTxScannedFromDB()
		require.NoError(t, err)
		require.EqualValues(t, lastTx, lastTxScanned)
	})
}

func TestReadWriteDBLastTxScanned(t *testing.T) {
	chain := chains.SolanaDevnet
	t.Run("should be able to write and read last tx scanned to db", func(t *testing.T) {
		// create observer and open db
		ob := newTestSuite(t, chain)

		// write last tx scanned
		lastTx := "5LuQMorgd11p8GWEw6pmyHCDtA26NUyeNFhLWPNk2oBoM9pkag1LzhwGSRos3j4TJLhKjswFhZkGtvSGdLDkmqsk"
		err := ob.WriteLastTxScannedToDB(lastTx)
		require.NoError(t, err)

		lastTxScanned, err := ob.ReadLastTxScannedFromDB()
		require.NoError(t, err)
		require.EqualValues(t, lastTx, lastTxScanned)
	})

	t.Run("should return error when last tx scanned not found in db", func(t *testing.T) {
		// create empty db
		ob := newTestSuite(t, chain)

		lastTxScanned, err := ob.ReadLastTxScannedFromDB()
		require.Error(t, err)
		require.Empty(t, lastTxScanned)
	})
}

func TestPostVoteInbound(t *testing.T) {
	t.Run("should be able to post vote inbound", func(t *testing.T) {
		// create observer
		ob := newTestSuite(t, chains.Ethereum)

		ob.zetacore.WithPostVoteInbound("", "sampleBallotIndex")

		// post vote inbound
		msg := sample.InboundVote(coin.CoinType_Gas, chains.Ethereum.ChainId, chains.ZetaChainMainnet.ChainId)
		ob.zetacore.MockGetCctxByHash(errors.New("not found"))
		ballot, err := ob.PostVoteInbound(context.TODO(), &msg, 100000)
		require.NoError(t, err)
		require.Equal(t, "sampleBallotIndex", ballot)
	})

	t.Run("should not post vote if message basic validation fails", func(t *testing.T) {
		// create observer
		ob := newTestSuite(t, chains.Ethereum)

		// create sample message with long Message
		msg := sample.InboundVote(coin.CoinType_Gas, chains.Ethereum.ChainId, chains.ZetaChainMainnet.ChainId)
		msg.Message = strings.Repeat("1", crosschaintypes.MaxMessageLength+1)
		ob.zetacore.MockGetCctxByHash(errors.New("not found"))

		// post vote inbound
		ballot, err := ob.PostVoteInbound(context.TODO(), &msg, 100000)
		require.NoError(t, err)
		require.Empty(t, ballot)
	})

	t.Run("should not post vote cctx already exists and ballot is not found", func(t *testing.T) {
		//Arrange
		// create observer
		ob := newTestSuite(t, chains.Ethereum)

		ob.zetacore.WithPostVoteInbound("", "sampleBallotIndex")
		msg := sample.InboundVote(coin.CoinType_Gas, chains.Ethereum.ChainId, chains.ZetaChainMainnet.ChainId)

		ob.zetacore.MockGetCctxByHash(nil)
		ob.zetacore.MockGetBallotByID(msg.Digest(), status.Error(codes.NotFound, "not found ballot"))

		var logBuffer bytes.Buffer
		consoleWriter := zerolog.ConsoleWriter{Out: &logBuffer}
		logger := zerolog.New(consoleWriter)
		ob.Observer.Logger().Inbound = logger

		// Act
		ballot, err := ob.PostVoteInbound(context.TODO(), &msg, 100000)
		// Assert
		require.NoError(t, err)
		require.Equal(t, ballot, msg.Digest())

		logOutput := logBuffer.String()
		require.Contains(t, logOutput, "inbound detected: CCTX exists but the ballot does not")
	})

	t.Run("should post vote cctx already exists but ballot is found", func(t *testing.T) {
		//Arrange
		// create observer
		ob := newTestSuite(t, chains.Ethereum)

		msg := sample.InboundVote(coin.CoinType_Gas, chains.Ethereum.ChainId, chains.ZetaChainMainnet.ChainId)
		ob.zetacore.WithPostVoteInbound(sample.ZetaIndex(t), msg.Digest())
		ob.zetacore.MockGetCctxByHash(nil)
		ob.zetacore.MockGetBallotByID(msg.Digest(), nil)

		var logBuffer bytes.Buffer
		consoleWriter := zerolog.ConsoleWriter{Out: &logBuffer}
		logger := zerolog.New(consoleWriter)
		ob.Observer.Logger().Inbound = logger

		// Act
		ballot, err := ob.PostVoteInbound(context.TODO(), &msg, 100000)
		// Assert
		require.NoError(t, err)
		require.Equal(t, ballot, msg.Digest())

		logOutput := logBuffer.String()
		require.Contains(t, logOutput, "inbound detected: vote posted")
	})
}

func createDatabase(t *testing.T) *db.DB {
	sqlDatabase, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	return sqlDatabase
}
