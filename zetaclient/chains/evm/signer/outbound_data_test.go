package signer

import (
	"context"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestSigner_SetChainAndSender(t *testing.T) {
	// setup inputs
	cctx := getCCTX(t)
	txData := &OutboundData{}
	logger := zerolog.Logger{}

	t.Run("SetChainAndSender PendingRevert", func(t *testing.T) {
		cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
		skipTx := txData.SetChainAndSender(cctx, logger)

		require.False(t, skipTx)
		require.Equal(t, ethcommon.HexToAddress(cctx.InboundParams.Sender), txData.to)
		require.Equal(t, big.NewInt(cctx.InboundParams.SenderChainId), txData.toChainID)
	})

	t.Run("SetChainAndSender PendingOutbound", func(t *testing.T) {
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		skipTx := txData.SetChainAndSender(cctx, logger)

		require.False(t, skipTx)
		require.Equal(t, ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver), txData.to)
		require.Equal(t, big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId), txData.toChainID)
	})

	t.Run("SetChainAndSender Should skip cctx", func(t *testing.T) {
		cctx.CctxStatus.Status = types.CctxStatus_PendingInbound
		skipTx := txData.SetChainAndSender(cctx, logger)
		require.True(t, skipTx)
	})
}

func TestSigner_SetupGas(t *testing.T) {
	cctx := getCCTX(t)
	evmSigner, err := getNewEvmSigner(nil)
	require.NoError(t, err)

	txData := &OutboundData{}
	logger := zerolog.Logger{}

	t.Run("SetupGas_success", func(t *testing.T) {
		chain := chains.BscMainnet
		err := txData.SetupGas(cctx, logger, evmSigner.EvmClient(), chain)
		require.NoError(t, err)
	})

	t.Run("SetupGas_error", func(t *testing.T) {
		cctx.GetCurrentOutboundParam().GasPrice = "invalidGasPrice"
		chain := chains.BscMainnet
		err := txData.SetupGas(cctx, logger, evmSigner.EvmClient(), chain)
		require.ErrorContains(t, err, "cannot convert gas price")
	})
}

func TestSigner_NewOutboundData(t *testing.T) {
	app := zctx.New(config.New(false), zerolog.Nop())
	ctx := zctx.WithAppContext(context.Background(), app)

	// Setup evm signer
	evmSigner, err := getNewEvmSigner(nil)
	require.NoError(t, err)

	mockObserver, err := getNewEvmChainObserver(t, nil)
	require.NoError(t, err)

	t.Run("NewOutboundData success", func(t *testing.T) {
		cctx := getCCTX(t)
		_, skip, err := NewOutboundData(ctx, cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
		require.False(t, skip)
		require.NoError(t, err)
	})

	t.Run("NewOutboundData skip", func(t *testing.T) {
		cctx := getCCTX(t)
		cctx.CctxStatus.Status = types.CctxStatus_Aborted
		_, skip, err := NewOutboundData(ctx, cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
		require.NoError(t, err)
		require.True(t, skip)
	})

	t.Run("NewOutboundData unknown chain", func(t *testing.T) {
		cctx := getInvalidCCTX(t)
		require.NoError(t, err)
		_, skip, err := NewOutboundData(ctx, cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
		require.ErrorContains(t, err, "unknown chain")
		require.True(t, skip)
	})

	t.Run("NewOutboundData setup gas error", func(t *testing.T) {
		cctx := getCCTX(t)
		require.NoError(t, err)
		cctx.GetCurrentOutboundParam().GasPrice = "invalidGasPrice"
		_, skip, err := NewOutboundData(ctx, cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
		require.True(t, skip)
		require.ErrorContains(t, err, "cannot convert gas price")
	})
}
