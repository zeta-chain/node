package context

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/ptr"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/config"
	"golang.org/x/exp/maps"
)

func TestAppContext(t *testing.T) {
	var (
		testCfg = config.New(false)
		logger  = zerolog.New(zerolog.NewTestWriter(t))

		ccFlags = types.CrosschainFlags{
			IsInboundEnabled:      true,
			IsOutboundEnabled:     true,
			GasPriceIncreaseFlags: nil,
		}

		opFlags = types.OperationalFlags{
			RestartHeight:         123,
			SignerBlockTimeOffset: ptr.Ptr(time.Second),
		}
	)

	testCfg.BTCChainConfigs[111] = config.BTCConfig{RPCUsername: "satoshi"}

	ethParams := types.GetDefaultEthMainnetChainParams()
	ethParams.IsSupported = true

	btcParams := types.GetDefaultBtcMainnetChainParams()
	btcParams.IsSupported = true

	solParams := sample.ChainParamsSupported(chains.SolanaLocalnet.ChainId)

	fancyL2 := chains.Chain{
		ChainId:     123,
		Network:     0,
		NetworkType: chains.NetworkType_mainnet,
		Vm:          chains.Vm_evm,
		Consensus:   chains.Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: 1,
	}

	fancyL2Params := types.GetDefaultEthMainnetChainParams()
	fancyL2Params.ChainId = fancyL2.ChainId
	fancyL2Params.IsSupported = true

	t.Run("Update", func(t *testing.T) {
		// Given AppContext
		appContext := New(testCfg, nil, logger)

		// With expected default behavior
		_, err := appContext.GetChain(123)
		require.ErrorIs(t, err, ErrChainNotFound)

		require.Equal(t, testCfg, appContext.Config())
		require.Empty(t, appContext.GetCrossChainFlags())
		require.False(t, appContext.IsInboundObservationEnabled())
		require.False(t, appContext.IsOutboundObservationEnabled())
		require.False(t, appContext.IsMempoolCongested())

		// Given some data that is supposed to come from ZetaCore RPC
		newChains := []chains.Chain{
			chains.Ethereum,
			chains.BitcoinMainnet,
			chains.SolanaLocalnet,
		}

		chainParams := map[int64]*types.ChainParams{
			chains.Ethereum.ChainId:       ethParams,
			chains.BitcoinMainnet.ChainId: btcParams,
			chains.SolanaLocalnet.ChainId: solParams,
			fancyL2.ChainId:               fancyL2Params,
		}

		additionalChains := []chains.Chain{
			fancyL2,
		}

		// ACT
		err = appContext.Update(newChains, additionalChains, chainParams, ccFlags, opFlags, 3000)

		// ASSERT
		require.NoError(t, err)

		// Check getters
		assert.Equal(t, testCfg, appContext.Config())
		assert.Equal(t, ccFlags, appContext.GetCrossChainFlags())
		assert.True(t, appContext.IsInboundObservationEnabled())
		assert.True(t, appContext.IsOutboundObservationEnabled())
		assert.True(t, appContext.IsMempoolCongested())

		// Check ETH Chain
		ethChain, err := appContext.GetChain(1)
		assert.NoError(t, err)
		assert.True(t, ethChain.IsEVM())
		assert.False(t, ethChain.IsBitcoin())
		assert.False(t, ethChain.IsSolana())
		assert.Equal(t, ethParams, ethChain.Params())

		// Check that fancyL2 chain is added as well
		fancyL2Chain, err := appContext.GetChain(fancyL2.ChainId)
		assert.NoError(t, err)
		assert.True(t, fancyL2Chain.IsEVM())
		assert.Equal(t, fancyL2Params, fancyL2Chain.Params())

		// Check chain IDs
		expectedIDs := []int64{ethParams.ChainId, btcParams.ChainId, solParams.ChainId, fancyL2.ChainId}
		assert.ElementsMatch(t, expectedIDs, appContext.ListChainIDs())

		// Check config
		assert.Equal(t, "satoshi", appContext.Config().BTCChainConfigs[111].RPCUsername)

		// Check cc flags
		assert.True(t, appContext.GetCrossChainFlags().IsInboundEnabled)

		// Check operational flags
		assert.Equal(t, time.Second, *appContext.GetOperationalFlags().SignerBlockTimeOffset)

		// Check unconfirmed tx count
		assert.EqualValues(t, 3000, appContext.GetUnconfirmedTxCount())

		t.Run("edge-cases", func(t *testing.T) {
			for _, tt := range []struct {
				name   string
				act    func(*AppContext) error
				assert func(*testing.T, *AppContext, error)
			}{
				{
					name: "update with empty chains results in an error",
					act: func(a *AppContext) error {
						return appContext.Update(newChains, nil, nil, ccFlags, opFlags, 1)
					},
					assert: func(t *testing.T, a *AppContext, err error) {
						assert.ErrorContains(t, err, "no chain params present")
					},
				},
				{
					name: "trying to add non-supported chain results in an error",
					act: func(a *AppContext) error {
						// ASSERT
						// GIven Optimism chain params from ZetaCore, but it's not supported YET
						op := chains.OptimismMainnet
						opParams := types.GetDefaultEthMainnetChainParams()
						opParams.ChainId = op.ChainId
						opParams.IsSupported = false

						chainsWithOpt := append(newChains, op)

						chainParamsWithOpt := maps.Clone(chainParams)
						chainParamsWithOpt[opParams.ChainId] = opParams

						return a.Update(chainsWithOpt, additionalChains, chainParamsWithOpt, ccFlags, opFlags, 1)
					},
					assert: func(t *testing.T, a *AppContext, err error) {
						assert.ErrorIs(t, err, ErrChainNotSupported)
						mustBeNotFound(t, a, chains.OptimismMainnet.ChainId)
					},
				},
				{
					name: "trying to add zeta chain without chain params is allowed",
					act: func(a *AppContext) error {
						chainsWithZeta := append(newChains, chains.ZetaChainMainnet)
						return a.Update(chainsWithZeta, additionalChains, chainParams, ccFlags, opFlags, 1)
					},
					assert: func(t *testing.T, a *AppContext, err error) {
						assert.NoError(t, err)

						zc := mustBePresent(t, a, chains.ZetaChainMainnet.ChainId)
						assert.True(t, zc.IsZeta())
					},
				},
				{
					name: "trying to add zetachain with chain params is allowed but forces fake params",
					act: func(a *AppContext) error {
						zetaParams := types.GetDefaultZetaPrivnetChainParams()
						zetaParams.ChainId = chains.ZetaChainMainnet.ChainId
						zetaParams.IsSupported = true
						zetaParams.GatewayAddress = "ABC123"

						chainParamsWithZeta := maps.Clone(chainParams)
						chainParamsWithZeta[zetaParams.ChainId] = zetaParams

						chainsWithZeta := append(newChains, chains.ZetaChainMainnet)

						return a.Update(chainsWithZeta, additionalChains, chainParamsWithZeta, ccFlags, opFlags, 1)
					},
					assert: func(t *testing.T, a *AppContext, err error) {
						assert.NoError(t, err)

						zc := mustBePresent(t, a, chains.ZetaChainMainnet.ChainId)
						assert.True(t, zc.IsZeta())
						assert.Equal(t, "", zc.Params().GatewayAddress)
					},
				},
				{
					name: "trying to add new chainParams without chain results in an error",
					act: func(a *AppContext) error {
						// ASSERT
						// Given polygon chain params WITHOUT the chain itself
						maticParams := types.GetDefaultMumbaiTestnetChainParams()
						maticParams.ChainId = chains.Polygon.ChainId
						maticParams.IsSupported = true

						updatedChainParams := maps.Clone(chainParams)
						updatedChainParams[maticParams.ChainId] = maticParams
						delete(updatedChainParams, chains.ZetaChainMainnet.ChainId)

						return a.Update(newChains, additionalChains, updatedChainParams, ccFlags, opFlags, 1)
					},
					assert: func(t *testing.T, a *AppContext, err error) {
						assert.ErrorContains(t, err, "unable to locate fresh chain 137 based on chain params")
						mustBeNotFound(t, a, chains.Polygon.ChainId)
					},
				},
			} {
				t.Run(tt.name, func(t *testing.T) {
					// ACT
					errAct := tt.act(appContext)

					// ASSERT
					require.NotNil(t, tt.assert)
					tt.assert(t, appContext, errAct)
				})
			}
		})
	})
}

func mustBeNotFound(t *testing.T, a *AppContext, chainID int64) {
	t.Helper()
	_, err := a.GetChain(chainID)
	require.ErrorIs(t, err, ErrChainNotFound)
}

func mustBePresent(t *testing.T, a *AppContext, chainID int64) Chain {
	t.Helper()
	c, err := a.GetChain(chainID)
	require.NoError(t, err)

	return c
}
