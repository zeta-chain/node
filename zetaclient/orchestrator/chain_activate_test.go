package orchestrator_test

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	context "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/orchestrator"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

// createTestAppContext creates a test app context with provided chain params and flags
func createTestAppContext(
	evmCfg config.EVMConfig,
	btcCfg config.BTCConfig,
	evmChain chains.Chain,
	btcChain chains.Chain,
	evmChainParams *observertypes.ChainParams,
	btcChainParams *observertypes.ChainParams,
) *context.AppContext {
	// create config
	cfg := config.NewConfig()
	cfg.EVMChainConfigs[evmChain.ChainId] = evmCfg
	cfg.BitcoinConfig = btcCfg

	// chains enabled
	chainsEnabled := []chains.Chain{evmChain}

	// create chain param map
	chainParamMap := make(map[int64]*observertypes.ChainParams)
	if evmChainParams != nil {
		chainParamMap[evmChain.ChainId] = evmChainParams
	}
	if btcChainParams != nil {
		chainParamMap[btcChain.ChainId] = btcChainParams
		chainsEnabled = append(chainsEnabled, btcChain)
	}

	// create app context
	appContext := context.NewAppContext(cfg)

	// create sample crosschain flags and header supported chains
	ccFlags := sample.CrosschainFlags()
	headerSupportedChains := sample.HeaderSupportedChains()

	// feed app context fields
	appContext.Update(
		observertypes.Keygen{},
		"testpubkey",
		chainsEnabled,
		chainParamMap,
		&chaincfg.MainNetParams,
		*ccFlags,
		headerSupportedChains,
		zerolog.Logger{},
	)
	return appContext
}

func Test_CreateObserversEVM(t *testing.T) {
	// define test chains and chain params
	evmChain := chains.Ethereum
	btcChain := chains.BitcoinMainnet
	evmChainParams := sample.ChainParams(evmChain.ChainId)

	// test cases
	tests := []struct {
		name               string
		evmCfg             config.EVMConfig
		btcCfg             config.BTCConfig
		evmChain           chains.Chain
		btcChain           chains.Chain
		evmChainParams     *observertypes.ChainParams
		dbPath             string
		numObserverCreated int
	}{
		{
			name: "should create observers for EVM chain and BTC chain",
			evmCfg: config.EVMConfig{
				Chain:    evmChain,
				Endpoint: "http://localhost:8545",
			},
			btcCfg:             config.BTCConfig{},
			evmChain:           evmChain,
			btcChain:           btcChain,
			evmChainParams:     evmChainParams,
			dbPath:             testutils.SQLiteMemory,
			numObserverCreated: 1,
		},
		{
			name: "should not create observer for EVM chain if chain params not found",
			evmCfg: config.EVMConfig{
				Chain:    evmChain,
				Endpoint: "http://localhost:8545",
			},
			btcCfg:             config.BTCConfig{},
			evmChain:           evmChain,
			btcChain:           btcChain,
			evmChainParams:     nil,
			dbPath:             testutils.SQLiteMemory,
			numObserverCreated: 0,
		},
		{
			name: "should not create observer for EVM chain if endpoint is invalid",
			evmCfg: config.EVMConfig{
				Chain:    evmChain,
				Endpoint: "invalid_endpoint",
			},
			btcCfg:             config.BTCConfig{},
			evmChain:           evmChain,
			btcChain:           btcChain,
			evmChainParams:     evmChainParams,
			dbPath:             testutils.SQLiteMemory,
			numObserverCreated: 0,
		},
		{
			name: "should not create observer for EVM chain if db path is invalid",
			evmCfg: config.EVMConfig{
				Chain:    evmChain,
				Endpoint: "http://localhost:8545",
			},
			btcCfg:             config.BTCConfig{},
			evmChain:           evmChain,
			btcChain:           btcChain,
			evmChainParams:     evmChainParams,
			dbPath:             "",
			numObserverCreated: 0,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create app context
			appCtx := createTestAppContext(tt.evmCfg, tt.btcCfg, tt.evmChain, tt.btcChain, tt.evmChainParams, nil)

			// create orchestrator
			ztacoreClient := mocks.NewMockZetacoreClient()
			oc := orchestrator.NewOrchestrator(appCtx, ztacoreClient, nil, base.Logger{}, tt.dbPath, nil)

			// create observers
			signerMap := make(map[int64]interfaces.ChainSigner)
			observerMap := make(map[int64]interfaces.ChainObserver)
			oc.CreateObserversEVM(signerMap, observerMap)

			// assert signer/observer map
			require.Len(t, signerMap, tt.numObserverCreated)
			require.Len(t, observerMap, tt.numObserverCreated)

			// assert signer/observer chain ID
			if tt.numObserverCreated > 0 {
				require.NotNil(t, signerMap[evmChain.ChainId])
			}
		})
	}
}

func Test_CreateObserversBTC(t *testing.T) {
	// define test chains and chain params
	evmChain := chains.Ethereum
	btcChain := chains.BitcoinMainnet
	btcChainParams := sample.ChainParams(btcChain.ChainId)

	// test cases
	tests := []struct {
		name               string
		evmCfg             config.EVMConfig
		btcCfg             config.BTCConfig
		evmChain           chains.Chain
		btcChain           chains.Chain
		btcChainParams     *observertypes.ChainParams
		dbPath             string
		numObserverCreated int
	}{
		{
			name:               "should not create observer for BTC chain if btc config is missing",
			evmCfg:             config.EVMConfig{},
			btcCfg:             config.BTCConfig{}, // empty config in file
			evmChain:           evmChain,
			btcChain:           btcChain,
			btcChainParams:     btcChainParams,
			dbPath:             testutils.SQLiteMemory,
			numObserverCreated: 0,
		},
		{
			name:   "should not create observer for BTC chain if chain is not enabled",
			evmCfg: config.EVMConfig{},
			btcCfg: config.BTCConfig{
				RPCUsername: "user",
			},
			evmChain:           evmChain,
			btcChain:           btcChain,
			btcChainParams:     nil, // disabled btc chain
			dbPath:             testutils.SQLiteMemory,
			numObserverCreated: 0,
		},
		{
			name:   "should not create observer for BTC chain if failed to Ping endpoint",
			evmCfg: config.EVMConfig{},
			btcCfg: config.BTCConfig{
				RPCUsername: "user",
			},
			evmChain:           evmChain,
			btcChain:           btcChain,
			btcChainParams:     btcChainParams,
			dbPath:             testutils.SQLiteMemory,
			numObserverCreated: 0,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create app context
			appCtx := createTestAppContext(tt.evmCfg, tt.btcCfg, tt.evmChain, tt.btcChain, nil, tt.btcChainParams)

			// create orchestrator
			ztacoreClient := mocks.NewMockZetacoreClient()
			oc := orchestrator.NewOrchestrator(appCtx, ztacoreClient, nil, base.Logger{}, tt.dbPath, nil)

			// create observers
			signerMap := make(map[int64]interfaces.ChainSigner)
			observerMap := make(map[int64]interfaces.ChainObserver)
			oc.CreateObserversBTC(signerMap, observerMap)

			// assert signer/observer map
			require.Len(t, signerMap, tt.numObserverCreated)
			require.Len(t, observerMap, tt.numObserverCreated)
		})
	}
}
