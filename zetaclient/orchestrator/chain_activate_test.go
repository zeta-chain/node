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
	appContext := context.New(cfg)

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
		[]chains.Chain{},
		headerSupportedChains,
		zerolog.Logger{},
	)
	return appContext
}

func Test_ActivateChains(t *testing.T) {
	// define test chain and chain params
	evmChain := chains.Ethereum
	evmChainParams := sample.ChainParams(evmChain.ChainId)

	// test cases
	tests := []struct {
		name           string
		evmCfg         config.EVMConfig
		btcCfg         config.BTCConfig
		evmChain       chains.Chain
		btcChain       chains.Chain
		evmChainParams *observertypes.ChainParams
		dbPath         string
		fail           bool
	}{
		{
			name: "should activate newly supported chains that are not in existing observer map",
			evmCfg: config.EVMConfig{
				Chain:    evmChain,
				Endpoint: "http://localhost:8545",
			},
			btcCfg:         config.BTCConfig{}, // btc chain is not needed for this test
			evmChain:       evmChain,
			evmChainParams: evmChainParams,
			dbPath:         testutils.SQLiteMemory,
			fail:           false,
		},
		{
			name: "should not activate chain if dbPath is invalid",
			evmCfg: config.EVMConfig{
				Chain:    evmChain,
				Endpoint: "http://localhost:8545",
			},
			btcCfg:         config.BTCConfig{}, // btc chain is not needed for this test
			evmChain:       evmChain,
			evmChainParams: evmChainParams,
			dbPath:         "", // invalid db path
			fail:           true,
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

			// create new signer and observer maps
			newSignerMap := make(map[int64]interfaces.ChainSigner)
			newObserverMap := make(map[int64]interfaces.ChainObserver)
			oc.CreateSignerObserverEVM(newSignerMap, newObserverMap)

			// activate chains
			oc.ActivateChains(newSignerMap, newObserverMap)

			// assert signer/observer map
			ob, err1 := oc.GetUpdatedChainObserver(tt.evmChain.ChainId)
			signer, err2 := oc.GetUpdatedSigner(tt.evmChain.ChainId)

			if tt.fail {
				require.Error(t, err1)
				require.Error(t, err2)
				require.Nil(t, ob)
				require.Nil(t, signer)
			} else {
				require.NoError(t, err1)
				require.NoError(t, err2)
				require.NotNil(t, ob)
				require.NotNil(t, signer)
			}
		})
	}
}

func Test_DeactivateChains(t *testing.T) {
	// define test chain and chain params
	evmChain := chains.Ethereum
	evmChainParams := sample.ChainParams(evmChain.ChainId)

	// test cases
	tests := []struct {
		name           string
		evmCfg         config.EVMConfig
		btcCfg         config.BTCConfig
		evmChain       chains.Chain
		btcChain       chains.Chain
		evmChainParams *observertypes.ChainParams
		dbPath         string
	}{
		{
			name: "should deactivate chains that are not in new observer map",
			evmCfg: config.EVMConfig{
				Chain:    evmChain,
				Endpoint: "http://localhost:8545",
			},
			btcCfg:         config.BTCConfig{}, // btc chain is not needed for this test
			evmChain:       evmChain,
			evmChainParams: evmChainParams,
			dbPath:         testutils.SQLiteMemory,
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

			// create new signer and observer maps
			newSignerMap := make(map[int64]interfaces.ChainSigner)
			newObserverMap := make(map[int64]interfaces.ChainObserver)
			oc.CreateSignerObserverEVM(newSignerMap, newObserverMap)

			// activate chains
			oc.ActivateChains(newSignerMap, newObserverMap)

			// assert signer/observer map
			ob, err := oc.GetUpdatedChainObserver(tt.evmChain.ChainId)
			require.NoError(t, err)
			require.NotNil(t, ob)

			// create new config and set EVM chain params as empty
			newCfg := appCtx.Config()
			newCfg.EVMChainConfigs = make(map[int64]config.EVMConfig)
			appCtx.SetConfig(newCfg)

			// create maps again based on newly updated config
			newSignerMap = make(map[int64]interfaces.ChainSigner)
			newObserverMap = make(map[int64]interfaces.ChainObserver)
			oc.CreateSignerObserverEVM(newSignerMap, newObserverMap)

			// deactivate chains
			oc.DeactivateChains(newObserverMap)

			// assert signer/observer map
			ob, err1 := oc.GetUpdatedChainObserver(tt.evmChain.ChainId)
			signer, err2 := oc.GetUpdatedSigner(tt.evmChain.ChainId)
			require.Error(t, err1)
			require.Error(t, err2)
			require.Nil(t, ob)
			require.Nil(t, signer)
		})
	}
}

func Test_CreateSignerObserverEVM(t *testing.T) {
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
			name: "should create observers for EVM chain",
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
			oc.CreateSignerObserverEVM(signerMap, observerMap)

			// assert signer/observer map
			require.Len(t, signerMap, tt.numObserverCreated)
			require.Len(t, observerMap, tt.numObserverCreated)

			// assert signer/observer chain ID
			if tt.numObserverCreated > 0 {
				require.NotNil(t, signerMap[evmChain.ChainId])
				require.NotNil(t, observerMap[evmChain.ChainId])
			}
		})
	}
}

func Test_CreateSignerObserverBTC(t *testing.T) {
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
			oc.CreateSignerObserverBTC(signerMap, observerMap)

			// assert signer/observer map
			require.Len(t, signerMap, tt.numObserverCreated)
			require.Len(t, observerMap, tt.numObserverCreated)
		})
	}
}
