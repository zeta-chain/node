package orchestrator

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	zctx "github.com/zeta-chain/node/zetaclient/context"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
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
		tonChain = chains.TONMainnet
	)

	var (
		evmChainParams = mocks.MockChainParams(evmChain.ChainId, 100)
		btcChainParams = mocks.MockChainParams(btcChain.ChainId, 100)
		solChainParams = mocks.MockChainParams(solChain.ChainId, 100)
		tonChainParams = mocks.MockChainParams(tonChain.ChainId, 100)
	)

	solChainParams.GatewayAddress = testutils.GatewayAddresses[solChain.ChainId]

	// new chain params in AppContext
	evmChainParamsNew := mocks.MockChainParams(evmChainParams.ChainId, 100)
	evmChainParamsNew.ConnectorContractAddress = testutils.OtherAddress1
	evmChainParamsNew.Erc20CustodyContractAddress = testutils.OtherAddress2

	// new solana chain params in AppContext
	solChainParamsNew := mocks.MockChainParams(solChain.ChainId, 100)
	solChainParamsNew.GatewayAddress = sample.SolanaAddress(t)

	tonChainParamsNew := mocks.MockChainParams(tonChain.ChainId, 100)
	tonChainParamsNew.GatewayAddress = sample.GenerateTONAccountID().ToRaw()

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

	t.Run("should be able to update ton gateway address", func(t *testing.T) {
		orchestrator := mockOrchestrator(t, nil,
			evmChain, btcChain, solChain, tonChain,
			evmChainParams, btcChainParams, solChainParams, tonChainParams,
		)

		appContext := createAppContext(t,
			evmChain, btcChain, solChain, tonChain,
			evmChainParams, btcChainParams, solChainParamsNew, tonChainParamsNew,
		)

		// update signer with new gateway address
		signer, err := orchestrator.resolveSigner(appContext, tonChain.ChainId)
		require.NoError(t, err)
		require.Equal(t, tonChainParamsNew.GatewayAddress, signer.GetGatewayAddress())
	})
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
		case chains.IsTONChain(cp.ChainId, nil):
			observers[cp.ChainId] = mocks.NewTONObserver(cp)
			signers[cp.ChainId] = mocks.NewTONSigner()
		default:
			t.Fatalf("mock orchestrator: unsupported chain %d", cp.ChainId)
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
	for _, c := range supportedChains {
		switch {
		case chains.IsEVMChain(c.ChainId, nil):
			cfg.EVMChainConfigs[c.ChainId] = config.EVMConfig{Endpoint: "localhost"}
		case chains.IsBitcoinChain(c.ChainId, nil):
			cfg.BTCChainConfigs[c.ChainId] = config.BTCConfig{RPCHost: "localhost"}
		case chains.IsSolanaChain(c.ChainId, nil):
			cfg.SolanaConfig = config.SolanaConfig{Endpoint: "localhost"}
		case chains.IsTONChain(c.ChainId, nil):
			cfg.TONConfig = config.TONConfig{LiteClientConfigURL: "localhost"}
		default:
			t.Fatalf("create app context: unsupported chain %d", c.ChainId)
		}
	}

	// chain params
	params := map[int64]*observertypes.ChainParams{}
	for i := range obsParams {
		cp := obsParams[i]
		params[cp.ChainId] = cp
	}

	// new AppContext
	appContext := zctx.New(cfg, nil, zerolog.New(zerolog.NewTestWriter(t)))

	ccFlags := sample.CrosschainFlags()
	opFlags := sample.OperationalFlags()

	// feed chain params
	err := appContext.Update(
		supportedChains,
		nil,
		params,
		*ccFlags,
		opFlags,
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
