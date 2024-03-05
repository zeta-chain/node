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

func TestSigner_SetTransactionData(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	txData := OutBoundTransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := txData.SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)
}
