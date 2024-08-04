package signer

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/constant"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
	"math/big"
	"testing"
)

func TestSigner_SignCommandTx(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner, err := getNewEvmSigner(nil)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver(t, nil)
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(ctx, cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignCommandTx CmdWhitelistERC20", func(t *testing.T) {
		cmd := constant.CmdWhitelistERC20
		params := ConnectorAddress.Hex()
		// Call SignCommandTx
		tx, err := evmSigner.SignCommandTx(ctx, txData, cmd, params)
		require.NoError(t, err)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Revert tx calls erc20 custody contract with 0 gas token
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, big.NewInt(0))
	})

	t.Run("SignCommandTx CmdMigrateTssFunds", func(t *testing.T) {
		cmd := constant.CmdMigrateTssFunds
		// Call SignCommandTx
		tx, err := evmSigner.SignCommandTx(ctx, txData, cmd, "")
		require.NoError(t, err)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
	})
}

func TestSigner_SignWhitelistERC20Cmd(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)

	mockObserver, err := getNewEvmChainObserver(t, tss)
	require.NoError(t, err)

	txData, skip, err := NewOutboundData(ctx, cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.NoError(t, err)
	require.False(t, skip)

	t.Run("SignWhitelistERC20Cmd - should successfully sign", func(t *testing.T) {
		// Call SignWhitelistERC20Cmd
		tx, err := evmSigner.SignWhitelistERC20Cmd(ctx, txData, sample.EthAddress().Hex())
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, zeroValue)
	})
	t.Run("SignWhitelistERC20Cmd - should fail on invalid erc20 address", func(t *testing.T) {
		tx, err := evmSigner.SignWhitelistERC20Cmd(ctx, txData, "")
		require.Nil(t, tx)
		require.ErrorContains(t, err, "invalid erc20 address")
	})
	t.Run("SignWhitelistERC20Cmd - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call SignWhitelistERC20Cmd
		tx, err := evmSigner.SignWhitelistERC20Cmd(ctx, txData, sample.EthAddress().Hex())
		require.ErrorContains(t, err, "sign whitelist error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignMigrateTssFundsCmd(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver(t, tss)
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(ctx, cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignMigrateTssFundsCmd - should successfully sign", func(t *testing.T) {
		// Call SignMigrateTssFundsCmd
		tx, err := evmSigner.SignMigrateTssFundsCmd(ctx, txData)
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
	})

	t.Run("SignMigrateTssFundsCmd - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call SignMigrateTssFundsCmd
		tx, err := evmSigner.SignMigrateTssFundsCmd(ctx, txData)
		require.ErrorContains(t, err, "SignMigrateTssFundsCmd error")
		require.Nil(t, tx)
	})
}
