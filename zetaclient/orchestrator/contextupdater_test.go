package orchestrator

import (
	"context"
	"testing"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func Test_UpdateAppContext(t *testing.T) {
	var (
		eth       = chains.Ethereum
		ethParams = mocks.MockChainParams(eth.ChainId, 100)

		btc       = chains.BitcoinMainnet
		btcParams = mocks.MockChainParams(btc.ChainId, 100)
	)

	t.Run("Updates app context", func(t *testing.T) {
		var (
			ctx      = context.Background()
			app      = createAppContext(t, eth, ethParams)
			zetacore = mocks.NewZetacoreClient(t)
			logger   = zerolog.New(zerolog.NewTestWriter(t))
		)

		// Given zetacore client that has eth and btc chains
		newChains := []chains.Chain{eth, btc}
		newParams := []*observertypes.ChainParams{&ethParams, &btcParams}
		ccFlags := observertypes.CrosschainFlags{
			IsInboundEnabled:  true,
			IsOutboundEnabled: true,
		}

		zetacore.On("GetBlockHeight", mock.Anything).Return(int64(123), nil)
		zetacore.On("GetUpgradePlan", mock.Anything).Return(nil, nil)
		zetacore.On("GetSupportedChains", mock.Anything).Return(newChains, nil)
		zetacore.On("GetAdditionalChains", mock.Anything).Return(nil, nil)
		zetacore.On("GetChainParams", mock.Anything).Return(newParams, nil)
		zetacore.On("GetKeyGen", mock.Anything).Return(observertypes.Keygen{}, nil)
		zetacore.On("GetCrosschainFlags", mock.Anything).Return(ccFlags, nil)
		zetacore.On("GetTSS", mock.Anything).Return(observertypes.TSS{TssPubkey: "0x123"}, nil)

		// ACT
		err := UpdateAppContext(ctx, app, zetacore, logger)

		// ASSERT
		require.NoError(t, err)

		// New chains should be added
		_, err = app.GetChain(btc.ChainId)
		require.NoError(t, err)
	})

	t.Run("Upgrade plan detected", func(t *testing.T) {
		// ARRANGE
		var (
			ctx      = context.Background()
			app      = createAppContext(t, eth, ethParams)
			zetacore = mocks.NewZetacoreClient(t)
			logger   = zerolog.New(zerolog.NewTestWriter(t))
		)

		zetacore.On("GetBlockHeight", mock.Anything).Return(int64(123), nil)
		zetacore.On("GetUpgradePlan", mock.Anything).Return(&upgradetypes.Plan{
			Name:   "hello",
			Height: 124,
		}, nil)

		// ACT
		err := UpdateAppContext(ctx, app, zetacore, logger)

		// ASSERT
		require.ErrorIs(t, err, ErrUpgradeRequired)
	})
}
