package orchestrator

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/ptr"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

func TestSigners(t *testing.T) {
	var (
		ts         = metrics.NewTelemetryServer()
		tss        = mocks.NewTSSMainnet()
		log        = zerolog.New(zerolog.NewTestWriter(t))
		baseLogger = base.Logger{Std: log, Compliance: log}
	)

	t.Run("CreateSignerMap", func(t *testing.T) {
		// ARRANGE
		// Given a BTC server
		_, btcConfig := testutils.NewBtcServer(t)

		// Given a zetaclient config with ETH, MATIC, and BTC chains
		cfg := config.New(false)

		cfg.EVMChainConfigs[chains.Ethereum.ChainId] = config.EVMConfig{
			Chain:    chains.Ethereum,
			Endpoint: mocks.EVMRPCEnabled,
		}

		cfg.EVMChainConfigs[chains.Polygon.ChainId] = config.EVMConfig{
			Chain:    chains.Polygon,
			Endpoint: mocks.EVMRPCEnabled,
		}

		cfg.BitcoinConfig = btcConfig

		// Given AppContext
		app := zctx.New(cfg, log)
		ctx := zctx.WithAppContext(context.Background(), app)

		// Given chain & chainParams "fetched" from zetacore
		// (note that slice LACKS polygon chain on purpose)
		supportedChains := []chains.Chain{
			chains.Ethereum,
			chains.BitcoinMainnet,
		}

		evmParams := map[int64]*observertypes.ChainParams{
			chains.Ethereum.ChainId: ptr.Ptr(mocks.MockChainParams(chains.Ethereum.ChainId, 10)),
		}

		btcParams := &observertypes.ChainParams{
			ChainId: chains.BitcoinMainnet.ChainId,
		}

		mustUpdateAppContext(t, app, supportedChains, evmParams, btcParams)

		// ACT
		signers, err := CreateSignerMap(ctx, tss, baseLogger, ts)

		// ASSERT
		assert.NoError(t, err)
		assert.NotEmpty(t, signers)

		// Okay, now we want to check that signers for EVM and BTC were created
		assert.Equal(t, 2, len(signers))
		hasSigner(t, signers, chains.Ethereum.ChainId)
		hasSigner(t, signers, chains.BitcoinMainnet.ChainId)

		t.Run("Add polygon in the runtime", func(t *testing.T) {
			// ARRANGE
			supportedChains = []chains.Chain{
				chains.Ethereum,
				chains.Polygon,
				chains.BitcoinMainnet,
			}

			evmParams = map[int64]*observertypes.ChainParams{
				chains.Ethereum.ChainId: ptr.Ptr(mocks.MockChainParams(chains.Ethereum.ChainId, 10)),
				chains.Polygon.ChainId:  ptr.Ptr(mocks.MockChainParams(chains.Polygon.ChainId, 10)),
			}

			btcParams = &observertypes.ChainParams{
				ChainId: chains.BitcoinMainnet.ChainId,
			}

			mustUpdateAppContext(t, app, supportedChains, evmParams, btcParams)

			// ACT
			sm := signerMap(signers)
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, ts, &sm)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 1, added)
			assert.Equal(t, 0, removed)

			hasSigner(t, signers, chains.Ethereum.ChainId)
			hasSigner(t, signers, chains.Polygon.ChainId)
			hasSigner(t, signers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("Disable ethereum in the runtime", func(t *testing.T) {
			// ARRANGE
			supportedChains = []chains.Chain{
				chains.Polygon,
				chains.BitcoinMainnet,
			}

			evmParams = map[int64]*observertypes.ChainParams{
				chains.Polygon.ChainId: ptr.Ptr(mocks.MockChainParams(chains.Polygon.ChainId, 10)),
			}

			btcParams = &observertypes.ChainParams{
				ChainId: chains.BitcoinMainnet.ChainId,
			}

			mustUpdateAppContext(t, app, supportedChains, evmParams, btcParams)

			// ACT
			sm := signerMap(signers)
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, ts, &sm)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 0, added)
			assert.Equal(t, 1, removed)

			missesSigner(t, signers, chains.Ethereum.ChainId)
			hasSigner(t, signers, chains.Polygon.ChainId)
			hasSigner(t, signers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("Re-enable ethereum in the runtime", func(t *testing.T) {
			// ARRANGE
			supportedChains = []chains.Chain{
				chains.Ethereum,
				chains.Polygon,
				chains.BitcoinMainnet,
			}

			evmParams = map[int64]*observertypes.ChainParams{
				chains.Ethereum.ChainId: ptr.Ptr(mocks.MockChainParams(chains.Ethereum.ChainId, 10)),
				chains.Polygon.ChainId:  ptr.Ptr(mocks.MockChainParams(chains.Polygon.ChainId, 10)),
			}

			btcParams = &observertypes.ChainParams{
				ChainId: chains.BitcoinMainnet.ChainId,
			}

			mustUpdateAppContext(t, app, supportedChains, evmParams, btcParams)

			// ACT
			sm := signerMap(signers)
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, ts, &sm)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 1, added)
			assert.Equal(t, 0, removed)

			hasSigner(t, signers, chains.Ethereum.ChainId)
			hasSigner(t, signers, chains.Polygon.ChainId)
			hasSigner(t, signers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("Disable btc in the runtime", func(t *testing.T) {
			// ARRANGE
			// Given updated data from zetacore containing polygon chain
			supportedChains = []chains.Chain{
				chains.Ethereum,
				chains.Polygon,
			}

			evmParams = map[int64]*observertypes.ChainParams{
				chains.Ethereum.ChainId: ptr.Ptr(mocks.MockChainParams(chains.Ethereum.ChainId, 10)),
				chains.Polygon.ChainId:  ptr.Ptr(mocks.MockChainParams(chains.Polygon.ChainId, 10)),
			}

			btcParams = &observertypes.ChainParams{}

			mustUpdateAppContext(t, app, supportedChains, evmParams, btcParams)

			// ACT
			sm := signerMap(signers)
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, ts, &sm)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 0, added)
			assert.Equal(t, 1, removed)

			hasSigner(t, signers, chains.Ethereum.ChainId)
			hasSigner(t, signers, chains.Polygon.ChainId)
			missesSigner(t, signers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("Re-enable btc in the runtime", func(t *testing.T) {
			// ARRANGE
			// Given updated data from zetacore containing polygon chain
			supportedChains = []chains.Chain{
				chains.Ethereum,
				chains.Polygon,
				chains.BitcoinMainnet,
			}

			evmParams = map[int64]*observertypes.ChainParams{
				chains.Ethereum.ChainId: ptr.Ptr(mocks.MockChainParams(chains.Ethereum.ChainId, 10)),
				chains.Polygon.ChainId:  ptr.Ptr(mocks.MockChainParams(chains.Polygon.ChainId, 10)),
			}

			btcParams = &observertypes.ChainParams{
				ChainId: chains.BitcoinMainnet.ChainId,
			}

			mustUpdateAppContext(t, app, supportedChains, evmParams, btcParams)

			// ACT
			sm := signerMap(signers)
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, ts, &sm)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 1, added)
			assert.Equal(t, 0, removed)

			hasSigner(t, signers, chains.Ethereum.ChainId)
			hasSigner(t, signers, chains.Polygon.ChainId)
			hasSigner(t, signers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("No changes", func(t *testing.T) {
			// ACT
			sm := signerMap(signers)
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, ts, &sm)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 0, added)
			assert.Equal(t, 0, removed)
		})
	})
}

func mustUpdateAppContext(
	_ *testing.T,
	app *zctx.AppContext,
	chains []chains.Chain,
	evmParams map[int64]*observertypes.ChainParams,
	utxoParams *observertypes.ChainParams,
) {
	app.Update(
		ptr.Ptr(app.GetKeygen()),
		chains,
		evmParams,
		utxoParams,
		app.GetCurrentTssPubKey(),
		app.GetCrossChainFlags(),
		app.GetAdditionalChains(),
		nil,
		false,
	)
}

func hasSigner(t *testing.T, signers map[int64]interfaces.ChainSigner, chainId int64) {
	signer, ok := signers[chainId]
	assert.True(t, ok, "missing signer for chain %d", chainId)
	assert.NotEmpty(t, signer)
}

func missesSigner(t *testing.T, signers map[int64]interfaces.ChainSigner, chainId int64) {
	_, ok := signers[chainId]
	assert.False(t, ok, "unexpected signer for chain %d", chainId)
}
