package signer

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
	"math/big"
	"testing"
)

func TestSigner_SignConnectorOnReceive(t *testing.T) {
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

	t.Run("SignConnectorOnReceive - should successfully sign", func(t *testing.T) {
		// Call SignConnectorOnReceive
		tx, err := evmSigner.SignConnectorOnReceive(ctx, txData)
		require.NoError(t, err)

		// Verify Signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())
	})
	t.Run("SignConnectorOnReceive - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call SignConnectorOnReceive
		tx, err := evmSigner.SignConnectorOnReceive(ctx, txData)
		require.ErrorContains(t, err, "sign onReceive error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignConnectorOnRevert(t *testing.T) {
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

	t.Run("SignConnectorOnRevert - should successfully sign", func(t *testing.T) {
		// Call SignConnectorOnRevert
		tx, err := evmSigner.SignConnectorOnRevert(ctx, txData)
		require.NoError(t, err)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Revert tx calls connector contract with 0 gas token
		verifyTxBodyBasics(t, tx, evmSigner.zetaConnectorAddress, txData.nonce, big.NewInt(0))
	})
	t.Run("SignConnectorOnRevert - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call SignConnectorOnRevert
		tx, err := evmSigner.SignConnectorOnRevert(ctx, txData)
		require.ErrorContains(t, err, "sign onRevert error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignCancel(t *testing.T) {
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

	t.Run("SignCancel - should successfully sign", func(t *testing.T) {
		// Call SignConnectorOnRevert
		tx, err := evmSigner.SignCancel(ctx, txData)
		require.NoError(t, err)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Cancel tx sends 0 gas token to TSS self address
		verifyTxBodyBasics(t, tx, evmSigner.TSS().EVMAddress(), txData.nonce, big.NewInt(0))
	})
	t.Run("SignCancel - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call SignCancel
		tx, err := evmSigner.SignCancel(ctx, txData)
		require.ErrorContains(t, err, "SignCancel error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignGasWithdraw(t *testing.T) {
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

	t.Run("SignGasWithdraw - should successfully sign", func(t *testing.T) {
		// Call SignGasWithdraw
		tx, err := evmSigner.SignGasWithdraw(ctx, txData)
		require.NoError(t, err)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
	})
	t.Run("SignGasWithdraw - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call SignGasWithdraw
		tx, err := evmSigner.SignGasWithdraw(ctx, txData)
		require.ErrorContains(t, err, "SignGasWithdraw error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignERC20Withdraw(t *testing.T) {
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

	t.Run("SignERC20Withdraw - should successfully sign", func(t *testing.T) {
		// Call SignERC20Withdraw
		tx, err := evmSigner.SignERC20Withdraw(ctx, txData)
		require.NoError(t, err)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Withdraw tx calls erc20 custody contract with 0 gas token
		verifyTxBodyBasics(t, tx, evmSigner.er20CustodyAddress, txData.nonce, big.NewInt(0))
	})

	t.Run("SignERC20Withdraw - should fail if keysign fails", func(t *testing.T) {
		// pause tss to make keysign fail
		tss.Pause()

		// Call SignERC20Withdraw
		tx, err := evmSigner.SignERC20Withdraw(ctx, txData)
		require.ErrorContains(t, err, "sign withdraw error")
		require.Nil(t, tx)
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
