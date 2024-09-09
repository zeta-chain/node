package orchestrator

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	zctx "github.com/zeta-chain/node/zetaclient/context"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/testutil/sample"
	crosschainkeeper "github.com/zeta-chain/node/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func Test_GetUpdatedSigner(t *testing.T) {
	// initial parameters for orchestrator creation
	var (
		evmChain = chains.Ethereum
		btcChain = chains.BitcoinMainnet
		solChain = chains.SolanaMainnet
	)

	var (
		evmChainParams = mocks.MockChainParams(evmChain.ChainId, 100)
		btcChainParams = mocks.MockChainParams(btcChain.ChainId, 100)
		solChainParams = mocks.MockChainParams(solChain.ChainId, 100)
	)

	solChainParams.GatewayAddress = solanacontracts.SolanaGatewayProgramID

	// new chain params in AppContext
	evmChainParamsNew := mocks.MockChainParams(evmChainParams.ChainId, 100)
	evmChainParamsNew.ConnectorContractAddress = testutils.OtherAddress1
	evmChainParamsNew.Erc20CustodyContractAddress = testutils.OtherAddress2

	// new solana chain params in AppContext
	solChainParamsNew := mocks.MockChainParams(solChain.ChainId, 100)
	solChainParamsNew.GatewayAddress = sample.SolanaAddress(t)

	t.Run("signer should not be found", func(t *testing.T) {
		orchestrator := mockOrchestrator(t, nil, evmChain, btcChain, evmChainParams, btcChainParams)
		appContext := createAppContext(t, evmChain, btcChain, evmChainParamsNew, btcChainParams)
		// BSC signer should not be found
		_, err := orchestrator.resolveSigner(appContext, chains.BscMainnet.ChainId)
		require.ErrorContains(t, err, "signer not found")
	})

	t.Run("should be able to update connector and erc20 custody address", func(t *testing.T) {
		orchestrator := mockOrchestrator(t, nil, evmChain, btcChain, evmChainParams, btcChainParams)
		appContext := createAppContext(t, evmChain, btcChain, evmChainParamsNew, btcChainParams)

		// update signer with new connector and erc20 custody address
		signer, err := orchestrator.resolveSigner(appContext, evmChain.ChainId)
		require.NoError(t, err)

		require.Equal(t, testutils.OtherAddress1, signer.GetZetaConnectorAddress().Hex())
		require.Equal(t, testutils.OtherAddress2, signer.GetERC20CustodyAddress().Hex())
	})

	t.Run("should be able to update solana gateway address", func(t *testing.T) {
		orchestrator := mockOrchestrator(t, nil,
			evmChain, btcChain, solChain,
			evmChainParams, btcChainParams, solChainParams,
		)

		appContext := createAppContext(t,
			evmChain, btcChain, solChain,
			evmChainParams, btcChainParams, solChainParamsNew,
		)

		// update signer with new gateway address
		signer, err := orchestrator.resolveSigner(appContext, solChain.ChainId)
		require.NoError(t, err)
		require.Equal(t, solChainParamsNew.GatewayAddress, signer.GetGatewayAddress())
	})
}

func Test_GetUpdatedChainObserver(t *testing.T) {
	// initial parameters for orchestrator creation
	var (
		evmChain = chains.Ethereum
		btcChain = chains.BitcoinMainnet
		solChain = chains.SolanaMainnet
	)

	var (
		evmChainParams = mocks.MockChainParams(evmChain.ChainId, 100)
		btcChainParams = mocks.MockChainParams(btcChain.ChainId, 100)
		solChainParams = mocks.MockChainParams(solChain.ChainId, 100)
	)

	solChainParams.GatewayAddress = solanacontracts.SolanaGatewayProgramID

	// new chain params in AppContext
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
	solChainParamsNew := &observertypes.ChainParams{
		ChainId:                     solChain.ChainId,
		ConfirmationCount:           10,
		GasPriceTicker:              5,
		InboundTicker:               6,
		OutboundTicker:              6,
		WatchUtxoTicker:             1,
		ZetaTokenContractAddress:    "",
		ConnectorContractAddress:    "",
		Erc20CustodyContractAddress: "",
		OutboundScheduleInterval:    10,
		OutboundScheduleLookahead:   10,
		BallotThreshold:             sdk.OneDec(),
		MinObserverDelegation:       sdk.OneDec(),
		IsSupported:                 true,
	}

	t.Run("evm chain observer should not be found", func(t *testing.T) {
		orchestrator := mockOrchestrator(
			t,
			nil,
			evmChain,
			btcChain,
			solChain,
			evmChainParams,
			btcChainParams,
			solChainParams,
		)
		appContext := createAppContext(t, evmChain, btcChain, evmChainParamsNew, btcChainParams)

		// BSC chain observer should not be found
		_, err := orchestrator.resolveObserver(appContext, chains.BscMainnet.ChainId)
		require.ErrorContains(t, err, "observer not found")
	})
	t.Run("chain params in evm chain observer should be updated successfully", func(t *testing.T) {
		orchestrator := mockOrchestrator(
			t,
			nil,
			evmChain,
			btcChain,
			solChain,
			evmChainParams,
			btcChainParams,
			solChainParams,
		)
		appContext := createAppContext(
			t,
			evmChain,
			btcChain,
			solChain,
			evmChainParamsNew,
			btcChainParams,
			solChainParams,
		)

		// update evm chain observer with new chain params
		chainOb, err := orchestrator.resolveObserver(appContext, evmChain.ChainId)
		require.NoError(t, err)
		require.NotNil(t, chainOb)
		require.True(t, observertypes.ChainParamsEqual(*evmChainParamsNew, chainOb.GetChainParams()))
	})

	t.Run("btc chain observer should not be found", func(t *testing.T) {
		orchestrator := mockOrchestrator(
			t,
			nil,
			evmChain,
			btcChain,
			solChain,
			evmChainParams,
			btcChainParams,
			solChainParams,
		)
		appContext := createAppContext(
			t,
			evmChain,
			btcChain,
			solChain,
			evmChainParams,
			btcChainParamsNew,
			solChainParams,
		)

		// BTC testnet chain observer should not be found
		_, err := orchestrator.resolveObserver(appContext, chains.BitcoinTestnet.ChainId)
		require.ErrorContains(t, err, "observer not found")
	})
	t.Run("chain params in btc chain observer should be updated successfully", func(t *testing.T) {
		orchestrator := mockOrchestrator(
			t,
			nil,
			evmChain, btcChain, solChain,
			evmChainParams, btcChainParams, solChainParams,
		)
		appContext := createAppContext(
			t,
			evmChain,
			btcChain,
			solChain,
			evmChainParams,
			btcChainParamsNew,
			solChainParams,
		)
		// update btc chain observer with new chain params
		chainOb, err := orchestrator.resolveObserver(appContext, btcChain.ChainId)
		require.NoError(t, err)
		require.NotNil(t, chainOb)
		require.True(t, observertypes.ChainParamsEqual(*btcChainParamsNew, chainOb.GetChainParams()))
	})
	t.Run("solana chain observer should not be found", func(t *testing.T) {
		orchestrator := mockOrchestrator(
			t,
			nil,
			evmChain, btcChain, solChain,
			evmChainParams, btcChainParams, solChainParams,
		)

		appContext := createAppContext(
			t,
			evmChain,
			btcChain,
			solChain,
			evmChainParams,
			btcChainParams,
			solChainParamsNew,
		)

		// Solana Devnet chain observer should not be found
		_, err := orchestrator.resolveObserver(appContext, chains.SolanaDevnet.ChainId)
		require.ErrorContains(t, err, "observer not found")
	})
	t.Run("chain params in solana chain observer should be updated successfully", func(t *testing.T) {
		orchestrator := mockOrchestrator(t, nil,
			evmChain, btcChain, solChain,
			evmChainParams, btcChainParams, solChainParams,
		)
		appContext := createAppContext(t,
			evmChain, btcChain, solChain,
			evmChainParams, btcChainParams, solChainParamsNew,
		)

		// update solana chain observer with new chain params
		chainOb, err := orchestrator.resolveObserver(appContext, solChain.ChainId)
		require.NoError(t, err)
		require.NotNil(t, chainOb)
		require.True(t, observertypes.ChainParamsEqual(*solChainParamsNew, chainOb.GetChainParams()))
	})
}

func Test_GetPendingCctxsWithinRateLimit(t *testing.T) {
	ctx := context.Background()

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
				// #nosec G115 len always positive
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
			client := mocks.NewZetacoreClient(t)

			// load mock data
			client.WithRateLimiterFlags(tt.rateLimiterFlags)
			client.WithRateLimiterInput(tt.response)
			client.WithPendingCctx(ethChain.ChainId, tt.ethCctxsFallback)
			client.WithPendingCctx(btcChain.ChainId, tt.btcCctxsFallback)

			// create orchestrator
			orchestrator := mockOrchestrator(t, client, ethChain, btcChain, ethChainParams, btcChainParams)

			chainIDs := lo.Map(foreignChains, func(c chains.Chain, _ int) int64 { return c.ChainId })

			// run the test
			cctxsMap, err := orchestrator.GetPendingCctxsWithinRateLimit(ctx, chainIDs)
			if tt.fail {
				assert.Error(t, err)
				assert.Empty(t, cctxsMap)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCctxsMap, cctxsMap)
			}
		})
	}
}

func mockOrchestrator(t *testing.T, zetaClient interfaces.ZetacoreClient, chainsOrParams ...any) *Orchestrator {
	supportedChains, obsParams := parseChainsWithParams(t, chainsOrParams...)

	var (
		signers   = make(map[int64]interfaces.ChainSigner)
		observers = make(map[int64]interfaces.ChainObserver)
	)

	mustFindChain := func(chainID int64) chains.Chain {
		for _, c := range supportedChains {
			if c.ChainId == chainID {
				return c
			}
		}

		t.Fatalf("mock orchestrator: must find chain: chain %d not found", chainID)

		return chains.Chain{}
	}

	for i := range obsParams {
		cp := obsParams[i]

		switch {
		case chains.IsEVMChain(cp.ChainId, nil):
			observers[cp.ChainId] = mocks.NewEVMObserver(cp)
			signers[cp.ChainId] = mocks.NewEVMSigner(
				mustFindChain(cp.ChainId),
				ethcommon.HexToAddress(cp.ConnectorContractAddress),
				ethcommon.HexToAddress(cp.Erc20CustodyContractAddress),
			)
		case chains.IsBitcoinChain(cp.ChainId, nil):
			observers[cp.ChainId] = mocks.NewBTCObserver(cp)
			signers[cp.ChainId] = mocks.NewBTCSigner()
		case chains.IsSolanaChain(cp.ChainId, nil):
			observers[cp.ChainId] = mocks.NewSolanaObserver(cp)
			signers[cp.ChainId] = mocks.NewSolanaSigner()
		default:
			t.Fatalf("mock orcestrator: unsupported chain %d", cp.ChainId)
		}
	}

	return &Orchestrator{
		zetacoreClient: zetaClient,
		signerMap:      signers,
		observerMap:    observers,
	}
}

func createAppContext(t *testing.T, chainsOrParams ...any) *zctx.AppContext {
	supportedChains, obsParams := parseChainsWithParams(t, chainsOrParams...)

	cfg := config.New(false)

	// Mock config
	cfg.BitcoinConfig = config.BTCConfig{
		RPCHost: "localhost",
	}

	for _, c := range supportedChains {
		if chains.IsEVMChain(c.ChainId, nil) {
			cfg.EVMChainConfigs[c.ChainId] = config.EVMConfig{Chain: c}
		}
	}

	params := map[int64]*observertypes.ChainParams{}
	for i := range obsParams {
		cp := obsParams[i]
		params[cp.ChainId] = cp
	}

	// new AppContext
	appContext := zctx.New(cfg, nil, zerolog.New(zerolog.NewTestWriter(t)))

	ccFlags := sample.CrosschainFlags()

	// feed chain params
	err := appContext.Update(
		observertypes.Keygen{},
		supportedChains,
		nil,
		params,
		"tssPubKey",
		*ccFlags,
	)
	require.NoError(t, err, "failed to update app context")

	return appContext
}

// handy helper for testing
func parseChainsWithParams(t *testing.T, chainsOrParams ...any) ([]chains.Chain, []*observertypes.ChainParams) {
	var (
		supportedChains = make([]chains.Chain, 0, len(chainsOrParams))
		obsParams       = make([]*observertypes.ChainParams, 0, len(chainsOrParams))
	)

	for _, something := range chainsOrParams {
		switch tt := something.(type) {
		case *chains.Chain:
			supportedChains = append(supportedChains, *tt)
		case chains.Chain:
			supportedChains = append(supportedChains, tt)
		case *observertypes.ChainParams:
			obsParams = append(obsParams, tt)
		case observertypes.ChainParams:
			obsParams = append(obsParams, &tt)
		default:
			t.Fatalf("parse chains and params: unsupported type %T (%+v)", tt, tt)
		}
	}

	return supportedChains, obsParams
}
