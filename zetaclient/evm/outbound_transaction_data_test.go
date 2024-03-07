package evm

import (
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	corecommon "github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestSigner_SetChainAndSender(t *testing.T) {
	// setup inputs
	cctx, err := getCCTX()
	require.NoError(t, err)

	txData := &OutBoundTransactionData{}
	logger := zerolog.Logger{}

	t.Run("SetChainAndSender PendingRevert", func(t *testing.T) {
		cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
		skipTx := txData.SetChainAndSender(cctx, logger)

		require.False(t, skipTx)
		require.Equal(t, ethcommon.HexToAddress(cctx.InboundTxParams.Sender), txData.to)
		require.Equal(t, big.NewInt(cctx.InboundTxParams.SenderChainId), txData.toChainID)
	})

	t.Run("SetChainAndSender PendingOutBound", func(t *testing.T) {
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		skipTx := txData.SetChainAndSender(cctx, logger)

		require.False(t, skipTx)
		require.Equal(t, ethcommon.HexToAddress(cctx.GetCurrentOutTxParam().Receiver), txData.to)
		require.Equal(t, big.NewInt(cctx.GetCurrentOutTxParam().ReceiverChainId), txData.toChainID)
	})

	t.Run("SetChainAndSender Should skip cctx", func(t *testing.T) {
		cctx.CctxStatus.Status = types.CctxStatus_PendingInbound
		skipTx := txData.SetChainAndSender(cctx, logger)
		require.True(t, skipTx)
	})
}

func TestSigner_SetupGas(t *testing.T) {
	cctx, err := getCCTX()
	require.NoError(t, err)

	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	txData := &OutBoundTransactionData{}
	logger := zerolog.Logger{}

	t.Run("SetupGas_success", func(t *testing.T) {
		chain := corecommon.BscMainnetChain()
		err := txData.SetupGas(cctx, logger, evmSigner.EvmClient(), &chain)
		require.NoError(t, err)
	})

	t.Run("SetupGas_error", func(t *testing.T) {
		cctx.GetCurrentOutTxParam().OutboundTxGasPrice = "invalidGasPrice"
		chain := corecommon.BscMainnetChain()
		err := txData.SetupGas(cctx, logger, evmSigner.EvmClient(), &chain)
		require.ErrorContains(t, err, "cannot convert gas price")
	})
}

func TestSigner_NewOutBoundTransactionData(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)

	t.Run("NewOutBoundTransactionData success", func(t *testing.T) {
		cctx, err := getCCTX()
		require.NoError(t, err)
		_, skip, err := NewOutBoundTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, 123)
		require.False(t, skip)
		require.NoError(t, err)
	})

	t.Run("NewOutBoundTransactionData skip", func(t *testing.T) {
		cctx, err := getCCTX()
		require.NoError(t, err)
		cctx.CctxStatus.Status = types.CctxStatus_Aborted
		_, skip, err := NewOutBoundTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, 123)
		require.NoError(t, err)
		require.True(t, skip)
	})

	t.Run("NewOutBoundTransactionData unknown chain", func(t *testing.T) {
		cctx, err := getInvalidCCTX()
		require.NoError(t, err)
		_, skip, err := NewOutBoundTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, 123)
		require.ErrorContains(t, err, "unknown chain")
		require.True(t, skip)
	})

	t.Run("NewOutBoundTransactionData setup gas error", func(t *testing.T) {
		cctx, err := getCCTX()
		require.NoError(t, err)
		cctx.GetCurrentOutTxParam().OutboundTxGasPrice = "invalidGasPrice"
		_, skip, err := NewOutBoundTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, 123)
		require.True(t, skip)
		require.ErrorContains(t, err, "cannot convert gas price")
	})
}
