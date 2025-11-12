package signer

import (
	"math/big"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSigner_SignConnectorOnReceive(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner := newTestSuite(t)

	// Setup txData struct

	cctx := getCCTX(t)
	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignConnectorOnReceive - should successfully sign", func(t *testing.T) {
		// Call SignConnectorOnReceive
		tx, err := evmSigner.SignConnectorOnReceive(ctx, txData)
		require.NoError(t, err)

		// Verify Signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())
	})
	t.Run("SignConnectorOnReceive - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		evmSigner.tss.Pause()

		// Call SignConnectorOnReceive
		tx, err := evmSigner.SignConnectorOnReceive(ctx, txData)
		require.ErrorContains(t, err, "sign onReceive error")
		require.Nil(t, tx)
		evmSigner.tss.Unpause()
	})

	t.Run("SignOutbound - should successfully sign LegacyTx", func(t *testing.T) {
		// Call SignOutbound
		tx, err := evmSigner.SignConnectorOnReceive(ctx, txData)
		require.NoError(t, err)

		// Verify Signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.client.Signer())

		// check that by default tx type is legacy tx
		assert.Equal(t, ethtypes.LegacyTxType, int(tx.Type()))
	})

	t.Run("SignOutbound - should successfully sign DynamicFeeTx", func(t *testing.T) {
		t.Skip("Skipped due to https://github.com/zeta-chain/node/issues/3221")
		// ARRANGE
		const (
			gwei        = 1_000_000_000
			priorityFee = 1 * gwei
			gasPrice    = 3 * gwei
		)

		// Given a CCTX with gas price and priority fee
		cctx := getCCTX(t)
		cctx.OutboundParams[0].GasPrice = big.NewInt(gasPrice).String()
		cctx.OutboundParams[0].GasPriorityFee = big.NewInt(priorityFee).String()

		// Given outbound data
		txData, skip, err := NewOutboundData(ctx, cctx, makeLogger(t))
		require.False(t, skip)
		require.NoError(t, err)

		// Given a working TSS
		evmSigner.tss.Unpause()

		// ACT
		tx, err := evmSigner.SignConnectorOnReceive(ctx, txData)
		require.NoError(t, err)

		// ASSERT
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// check that by default tx type is a dynamic fee tx
		assert.Equal(t, ethtypes.DynamicFeeTxType, int(tx.Type()))

		// check that the gasPrice & priorityFee are set correctly
		assert.Equal(t, int64(gasPrice), tx.GasFeeCap().Int64())
		assert.Equal(t, int64(priorityFee), tx.GasTipCap().Int64())
	})
}

func TestSigner_SignConnectorOnRevert(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner := newTestSuite(t)

	// Setup txData struct
	cctx := getCCTX(t)
	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignConnectorOnRevert - should successfully sign", func(t *testing.T) {
		// Call SignConnectorOnRevert
		tx, err := evmSigner.SignConnectorOnRevert(ctx, txData)
		require.NoError(t, err)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Revert tx calls connector contract with 0 gas token
		verifyTxBodyBasics(t, tx, evmSigner.zetaConnectorAddress, txData.nonce, big.NewInt(0))
	})
	t.Run("SignConnectorOnRevert - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		evmSigner.tss.Pause()

		// Call SignConnectorOnRevert
		tx, err := evmSigner.SignConnectorOnRevert(ctx, txData)
		require.ErrorContains(t, err, "sign onRevert error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignCancel(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner := newTestSuite(t)

	// Setup txData struct
	cctx := getCCTX(t)
	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignCancel - should successfully sign", func(t *testing.T) {
		// Call SignConnectorOnRevert
		tx, err := evmSigner.SignCancel(ctx, txData)
		require.NoError(t, err)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Cancel tx sends 0 gas token to TSS self address
		verifyTxBodyBasics(t, tx, evmSigner.tss.PubKey().AddressEVM(), txData.nonce, big.NewInt(0))
	})
	t.Run("SignCancel - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		evmSigner.tss.Pause()

		// Call SignCancel
		tx, err := evmSigner.SignCancel(ctx, txData)
		require.ErrorContains(t, err, "SignCancel error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignGasWithdraw(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner := newTestSuite(t)

	// Setup txData struct
	cctx := getCCTX(t)
	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignGasWithdraw - should successfully sign", func(t *testing.T) {
		// Call SignGasWithdraw
		tx, err := evmSigner.SignGasWithdraw(ctx, txData)
		require.NoError(t, err)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
	})
	t.Run("SignGasWithdraw - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		evmSigner.tss.Pause()

		// Call SignGasWithdraw
		tx, err := evmSigner.SignGasWithdraw(ctx, txData)
		require.ErrorContains(t, err, "SignGasWithdraw error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignERC20Withdraw(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner := newTestSuite(t)

	// Setup txData struct
	cctx := getCCTX(t)
	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignERC20WithdrawTx - should successfully sign", func(t *testing.T) {
		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20Withdraw(ctx, txData)
		require.NoError(t, err)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Withdraw tx calls erc20 custody contract with 0 gas token
		verifyTxBodyBasics(t, tx, evmSigner.er20CustodyAddress, txData.nonce, big.NewInt(0))
	})

	t.Run("SignERC20WithdrawTx - should fail if keysign fails", func(t *testing.T) {
		// pause tss to make keysign fail
		evmSigner.tss.Pause()

		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20Withdraw(ctx, txData)
		require.ErrorContains(t, err, "sign withdraw error")
		require.Nil(t, tx)
	})
}
