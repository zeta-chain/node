package orchestrator

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

// MockOrchestrator creates a mock orchestrator for testing
func MockOrchestrator(
	t *testing.T,
	appContext *context.AppContext,
	zetacoreClient interfaces.ZetacoreClient,
	evmChain, btcChain chains.Chain,
	evmChainParams, btcChainParams *observertypes.ChainParams,
) *Orchestrator {
	// create mock signers and clients
	evmSigner := mocks.NewEVMSigner(
		evmChain,
		ethcommon.HexToAddress(evmChainParams.ConnectorContractAddress),
		ethcommon.HexToAddress(evmChainParams.Erc20CustodyContractAddress),
	)
	btcSigner := mocks.NewBTCSigner()
	evmObserver := mocks.NewEVMObserver(evmChainParams)
	btcObserver := mocks.NewBTCObserver(btcChainParams)

	// create orchestrator
	orchestrator := &Orchestrator{
		appContext:     appContext,
		zetacoreClient: zetacoreClient,
		signerMap: map[int64]interfaces.ChainSigner{
			evmChain.ChainId: evmSigner,
			btcChain.ChainId: btcSigner,
		},
		observerMap: map[int64]interfaces.ChainObserver{
			evmChain.ChainId: evmObserver,
			btcChain.ChainId: btcObserver,
		},
	}
	return orchestrator
}

// CreateTestAppContext creates a test app context for orchestrator testing
func CreateTestAppContext(
	evmChain, btcChain chains.Chain,
	evmChainParams, btcChainParams *observertypes.ChainParams,
) *context.AppContext {
	// new config
	cfg := config.NewConfig()
	cfg.EVMChainConfigs[evmChain.ChainId] = config.EVMConfig{
		Chain: evmChain,
	}
	cfg.BitcoinConfig = config.BTCConfig{
		RPCHost: "localhost",
	}
	// new app context
	appContext := context.New(cfg)
	chainParamsMap := make(map[int64]*observertypes.ChainParams)
	chainParamsMap[evmChain.ChainId] = evmChainParams
	chainParamsMap[btcChain.ChainId] = btcChainParams
	ccFlags := sample.CrosschainFlags()
	verificationFlags := sample.HeaderSupportedChains()

	// feed chain params
	appContext.Update(
		observertypes.Keygen{},
		"",
		[]chains.Chain{evmChain, btcChain},
		chainParamsMap,
		&chaincfg.RegressionNetParams,
		*ccFlags,
		verificationFlags,
		zerolog.Logger{},
	)
	return appContext
}

func Test_GetUpdatedSigner(t *testing.T) {
	// initial parameters for orchestrator creation
	evmChain := chains.Ethereum
	btcChain := chains.BitcoinMainnet
	evmChainParams := &observertypes.ChainParams{
		ChainId:                     evmChain.ChainId,
		ConnectorContractAddress:    testutils.ConnectorAddresses[evmChain.ChainId].Hex(),
		Erc20CustodyContractAddress: testutils.CustodyAddresses[evmChain.ChainId].Hex(),
	}
	btcChainParams := &observertypes.ChainParams{}

	// new chain params in app context
	evmChainParamsNew := &observertypes.ChainParams{
		ChainId:                     evmChain.ChainId,
		ConnectorContractAddress:    testutils.OtherAddress1,
		Erc20CustodyContractAddress: testutils.OtherAddress2,
	}

	t.Run("signer should not be found", func(t *testing.T) {
		appCtx := CreateTestAppContext(evmChain, btcChain, evmChainParamsNew, btcChainParams)
		orchestrator := MockOrchestrator(t, appCtx, nil, evmChain, btcChain, evmChainParams, btcChainParams)
		// BSC signer should not be found
		_, err := orchestrator.GetUpdatedSigner(chains.BscMainnet.ChainId)
		require.ErrorContains(t, err, "signer not found")
	})
	t.Run("should be able to update connector and erc20 custody address", func(t *testing.T) {
		appCtx := CreateTestAppContext(evmChain, btcChain, evmChainParamsNew, btcChainParams)
		orchestrator := MockOrchestrator(t, appCtx, nil, evmChain, btcChain, evmChainParams, btcChainParams)

		// update signer with new connector and erc20 custody address
		signer, err := orchestrator.GetUpdatedSigner(evmChain.ChainId)
		require.NoError(t, err)
		require.Equal(t, testutils.OtherAddress1, signer.GetZetaConnectorAddress().Hex())
		require.Equal(t, testutils.OtherAddress2, signer.GetERC20CustodyAddress().Hex())
	})
}

func Test_GetUpdatedChainObserver(t *testing.T) {
	// initial parameters for orchestrator creation
	evmChain := chains.Ethereum
	btcChain := chains.BitcoinMainnet
	evmChainParams := &observertypes.ChainParams{
		ChainId:                     evmChain.ChainId,
		ConnectorContractAddress:    testutils.ConnectorAddresses[evmChain.ChainId].Hex(),
		Erc20CustodyContractAddress: testutils.CustodyAddresses[evmChain.ChainId].Hex(),
	}
	btcChainParams := &observertypes.ChainParams{
		ChainId: btcChain.ChainId,
	}

	// new chain params in app context
	evmChainParamsNew := &observertypes.ChainParams{
		ChainId:                     evmChain.ChainId,
		ConfirmationCount:           10,
		GasPriceTicker:              11,
		InboundTicker:               12,
		OutboundTicker:              13,
		WatchUtxoTicker:             14,
		ZetaTokenContractAddress:    testutils.OtherAddress1,
		ConnectorContractAddress:    testutils.OtherAddress2,
		Erc20CustodyContractAddress: testutils.OtherAddress3,
		OutboundScheduleInterval:    15,
		OutboundScheduleLookahead:   16,
		BallotThreshold:             sdk.OneDec(),
		MinObserverDelegation:       sdk.OneDec(),
		IsSupported:                 true,
	}
	btcChainParamsNew := &observertypes.ChainParams{
		ChainId:                     btcChain.ChainId,
		ConfirmationCount:           3,
		GasPriceTicker:              300,
		InboundTicker:               60,
		OutboundTicker:              60,
		WatchUtxoTicker:             30,
		ZetaTokenContractAddress:    testutils.OtherAddress1,
		ConnectorContractAddress:    testutils.OtherAddress2,
		Erc20CustodyContractAddress: testutils.OtherAddress3,
		OutboundScheduleInterval:    60,
		OutboundScheduleLookahead:   200,
		BallotThreshold:             sdk.OneDec(),
		MinObserverDelegation:       sdk.OneDec(),
		IsSupported:                 true,
	}

	t.Run("evm chain observer should not be found", func(t *testing.T) {
		appCtx := CreateTestAppContext(evmChain, btcChain, evmChainParamsNew, btcChainParams)
		orchestrator := MockOrchestrator(t, appCtx, nil, evmChain, btcChain, evmChainParams, btcChainParams)

		// BSC chain observer should not be found
		_, err := orchestrator.GetUpdatedChainObserver(chains.BscMainnet.ChainId)
		require.ErrorContains(t, err, "chain observer not found")
	})
	t.Run("chain params in evm chain observer should be updated successfully", func(t *testing.T) {
		appCtx := CreateTestAppContext(evmChain, btcChain, evmChainParamsNew, btcChainParams)
		orchestrator := MockOrchestrator(t, appCtx, nil, evmChain, btcChain, evmChainParams, btcChainParams)

		// update evm chain observer with new chain params
		chainOb, err := orchestrator.GetUpdatedChainObserver(evmChain.ChainId)
		require.NoError(t, err)
		require.NotNil(t, chainOb)
		require.True(t, observertypes.ChainParamsEqual(*evmChainParamsNew, chainOb.GetChainParams()))
	})
	t.Run("btc chain observer should not be found", func(t *testing.T) {
		appCtx := CreateTestAppContext(btcChain, btcChain, evmChainParams, btcChainParamsNew)
		orchestrator := MockOrchestrator(t, appCtx, nil, evmChain, btcChain, evmChainParams, btcChainParams)

		// BTC testnet chain observer should not be found
		_, err := orchestrator.GetUpdatedChainObserver(chains.BitcoinTestnet.ChainId)
		require.ErrorContains(t, err, "chain observer not found")
	})
	t.Run("chain params in btc chain observer should be updated successfully", func(t *testing.T) {
		appCtx := CreateTestAppContext(btcChain, btcChain, evmChainParams, btcChainParamsNew)
		orchestrator := MockOrchestrator(t, appCtx, nil, evmChain, btcChain, evmChainParams, btcChainParams)

		// update btc chain observer with new chain params
		chainOb, err := orchestrator.GetUpdatedChainObserver(btcChain.ChainId)
		require.NoError(t, err)
		require.NotNil(t, chainOb)
		require.True(t, observertypes.ChainParamsEqual(*btcChainParamsNew, chainOb.GetChainParams()))
	})
}

func Test_GetPendingCctxsWithinRatelimit(t *testing.T) {
	// define test foreign chains
	ethChain := chains.Ethereum
	btcChain := chains.BitcoinMainnet
	zetaChainID := chains.ZetaChainMainnet.ChainId
	foreignChains := []chains.Chain{
		ethChain,
		btcChain,
	}

	// chain params
	ethChainParams := &observertypes.ChainParams{ChainId: ethChain.ChainId}
	btcChainParams := &observertypes.ChainParams{ChainId: btcChain.ChainId}

	// create 10 missed and 90 pending cctxs for eth chain, the coinType/amount does not matter for this test
	ethCctxsMissed := sample.CustomCctxsInBlockRange(
		t,
		1,
		10,
		zetaChainID,
		ethChain.ChainId,
		coin.CoinType_Gas,
		"",
		uint64(2e14),
		crosschaintypes.CctxStatus_PendingOutbound,
	)
	ethCctxsPending := sample.CustomCctxsInBlockRange(
		t,
		11,
		100,
		zetaChainID,
		ethChain.ChainId,
		coin.CoinType_Gas,
		"",
		uint64(2e14),
		crosschaintypes.CctxStatus_PendingOutbound,
	)
	ethCctxsAll := append(append([]*crosschaintypes.CrossChainTx{}, ethCctxsMissed...), ethCctxsPending...)

	// create 10 missed and 90 pending cctxs for btc chain, the coinType/amount does not matter for this test
	btcCctxsMissed := sample.CustomCctxsInBlockRange(
		t,
		1,
		10,
		zetaChainID,
		btcChain.ChainId,
		coin.CoinType_Gas,
		"",
		2000,
		crosschaintypes.CctxStatus_PendingOutbound,
	)
	btcCctxsPending := sample.CustomCctxsInBlockRange(
		t,
		11,
		100,
		zetaChainID,
		btcChain.ChainId,
		coin.CoinType_Gas,
		"",
		2000,
		crosschaintypes.CctxStatus_PendingOutbound,
	)
	btcCctxsAll := append(append([]*crosschaintypes.CrossChainTx{}, btcCctxsMissed...), btcCctxsPending...)

	// all missed cctxs and all pending cctxs across all foreign chains
	allCctxsMissed := crosschainkeeper.SortCctxsByHeightAndChainID(
		append(append([]*crosschaintypes.CrossChainTx{}, ethCctxsMissed...), btcCctxsMissed...))
	allCctxsPending := crosschainkeeper.SortCctxsByHeightAndChainID(
		append(append([]*crosschaintypes.CrossChainTx{}, ethCctxsPending...), btcCctxsPending...))

	// define test cases
	tests := []struct {
		name             string
		rateLimiterFlags *crosschaintypes.RateLimiterFlags
		response         *crosschaintypes.QueryRateLimiterInputResponse
		ethCctxsFallback []*crosschaintypes.CrossChainTx
		btcCctxsFallback []*crosschaintypes.CrossChainTx

		// expected result map
		fail             bool
		expectedCctxsMap map[int64][]*crosschaintypes.CrossChainTx
	}{
		{
			name:             "should return all missed and pending cctxs using fallback",
			rateLimiterFlags: &crosschaintypes.RateLimiterFlags{Enabled: false},
			response:         &crosschaintypes.QueryRateLimiterInputResponse{},
			ethCctxsFallback: ethCctxsAll,
			btcCctxsFallback: btcCctxsAll,
			expectedCctxsMap: map[int64][]*crosschaintypes.CrossChainTx{
				ethChain.ChainId: ethCctxsAll,
				btcChain.ChainId: btcCctxsAll,
			},
		},
		{
			name: "should return all missed and pending cctxs without fallback",
			rateLimiterFlags: &crosschaintypes.RateLimiterFlags{
				Enabled: true,
				Window:  100,
				Rate:    sdk.NewUint(1e18), // 1 ZETA/block
			},
			response: &crosschaintypes.QueryRateLimiterInputResponse{
				Height:       100,
				CctxsMissed:  allCctxsMissed,
				CctxsPending: allCctxsPending,
				// #nosec G701 len always positive
				TotalPending:            uint64(len(allCctxsPending) + len(allCctxsMissed)),
				PastCctxsValue:          sdk.NewInt(10).Mul(sdk.NewInt(1e18)).String(), // 10 ZETA
				PendingCctxsValue:       sdk.NewInt(90).Mul(sdk.NewInt(1e18)).String(), // 90 ZETA
				LowestPendingCctxHeight: 11,
			},
			ethCctxsFallback: nil,
			btcCctxsFallback: nil,
			expectedCctxsMap: map[int64][]*crosschaintypes.CrossChainTx{
				ethChain.ChainId: ethCctxsAll,
				btcChain.ChainId: btcCctxsAll,
			},
		},
		{
			name:             "should fail if cannot query rate limiter flags",
			rateLimiterFlags: nil,
			fail:             true,
		},
		{
			name: "should fail if cannot query rate limiter input",
			rateLimiterFlags: &crosschaintypes.RateLimiterFlags{
				Enabled: true,
				Window:  100,
				Rate:    sdk.NewUint(1e18), // 1 ZETA/block
			},
			response: nil,
			fail:     true,
		},
		{
			name: "should fail on invalid rate limiter input",
			rateLimiterFlags: &crosschaintypes.RateLimiterFlags{
				Enabled: true,
				Window:  100,
				Rate:    sdk.NewUint(1e18), // 1 ZETA/block
			},
			response: &crosschaintypes.QueryRateLimiterInputResponse{
				PastCctxsValue: "invalid",
			},
			fail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create mock zetacore client
			client := mocks.NewMockZetacoreClient()

			// load mock data
			client.WithRateLimiterFlags(tt.rateLimiterFlags)
			client.WithPendingCctx(ethChain.ChainId, tt.ethCctxsFallback)
			client.WithPendingCctx(btcChain.ChainId, tt.btcCctxsFallback)
			client.WithRateLimiterInput(tt.response)

			// create orchestrator
			orchestrator := MockOrchestrator(t, nil, client, ethChain, btcChain, ethChainParams, btcChainParams)

			// run the test
			cctxsMap, err := orchestrator.GetPendingCctxsWithinRatelimit(foreignChains)
			if tt.fail {
				require.Error(t, err)
				require.Nil(t, cctxsMap)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedCctxsMap, cctxsMap)
			}
		})
	}
}
