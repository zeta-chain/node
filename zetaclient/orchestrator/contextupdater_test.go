package orchestrator

import (
	"testing"
	"time"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/ptr"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/testutils/testlog"
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
			logger                 = testlog.New(t).Logger
			chainList, chainParams = parseChainsWithParams(t, eth, ethParams)
			ctx, app               = newAppContext(t, logger, chainList, chainParams)
			zetacore               = mocks.NewZetacoreClient(t)
		)

		// Given zetacore client that has eth and btc chains
		newChains := []chains.Chain{eth, btc}
		newParams := []*observertypes.ChainParams{&ethParams, &btcParams}
		ccFlags := observertypes.CrosschainFlags{
			IsInboundEnabled:  true,
			IsOutboundEnabled: true,
		}
		opFlags := observertypes.OperationalFlags{
			RestartHeight:         123,
			SignerBlockTimeOffset: ptr.Ptr(time.Second),
			MinimumVersion:        "",
		}

		on(zetacore, "GetBlockHeight", 1).Return(int64(123), nil)
		on(zetacore, "GetUpgradePlan", 1).Return(nil, nil)
		on(zetacore, "GetSupportedChains", 1).Return(newChains, nil)
		on(zetacore, "GetAdditionalChains", 1).Return(nil, nil)
		on(zetacore, "GetChainParams", 1).Return(newParams, nil)
		on(zetacore, "GetCrosschainFlags", 1).Return(ccFlags, nil)
		on(zetacore, "GetOperationalFlags", 1).Return(opFlags, nil)
		on(zetacore, "GetBaseGasPrice", 1).Return(int64(1000), nil)
		on(zetacore, "GetNumberOfUnconfirmedTxs", 1).Return(int(1), nil)

		// ACT
		err := UpdateAppContext(ctx, app, zetacore, logger)

		// ASSERT
		require.NoError(t, err)

		// New chains should be added
		_, err = app.GetChain(btc.ChainId)
		require.NoError(t, err)

		// Check OP flags
		require.Equal(t, opFlags.RestartHeight, app.GetOperationalFlags().RestartHeight)

		// Check base gas price
		require.EqualValues(t, 1000, app.GetCurrentBaseFee())

		// Check unconfirmed tx count
		require.EqualValues(t, 1, app.GetUnconfirmedTxCount())
	})

	t.Run("Upgrade plan detected", func(t *testing.T) {
		// ARRANGE
		var (
			logger                 = testlog.New(t).Logger
			chainList, chainParams = parseChainsWithParams(t, eth, ethParams)
			ctx, app               = newAppContext(t, logger, chainList, chainParams)
			zetacore               = mocks.NewZetacoreClient(t)
		)

		on(zetacore, "GetBlockHeight", 1).Return(int64(123), nil)
		on(zetacore, "GetUpgradePlan", 1).Return(&upgradetypes.Plan{
			Name:   "hello",
			Height: 124,
		}, nil)

		// ACT
		err := UpdateAppContext(ctx, app, zetacore, logger)

		// ASSERT
		require.ErrorIs(t, err, ErrUpgradeRequired)
	})
}
