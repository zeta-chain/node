package signer

import (
	"math/big"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestSigner_SignAdminTx(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner := newTestSuite(t)

	// Setup txData struct
	cctx := getCCTX(t)

	txData, skip, err := NewOutboundData(ctx, cctx, 123, zerolog.Logger{})

	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignAdminTx CmdWhitelistAsset", func(t *testing.T) {
		cmd := constant.CmdWhitelistAsset
		params := ConnectorAddress.Hex()
		// Call SignAdminTx
		tx, err := evmSigner.SignAdminTx(ctx, txData, cmd, params)
		require.NoError(t, err)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Revert tx calls erc20 custody contract with 0 gas token
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, big.NewInt(0))
	})

	t.Run("SignAdminTx CmdMigrateTSSFunds", func(t *testing.T) {
		cmd := constant.CmdMigrateTSSFunds
		// Call SignAdminTx
		tx, err := evmSigner.SignAdminTx(ctx, txData, cmd, "")
		require.NoError(t, err)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
	})
}

func TestSigner_SignWhitelistERC20Cmd(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner := newTestSuite(t)

	// Setup txData struct
	cctx := getCCTX(t)

	txData, skip, err := NewOutboundData(ctx, cctx, 123, zerolog.Logger{})

	require.NoError(t, err)
	require.False(t, skip)

	t.Run("signWhitelistERC20Cmd - should successfully sign", func(t *testing.T) {
		// Call signWhitelistERC20Cmd
		tx, err := evmSigner.signWhitelistERC20Cmd(ctx, txData, sample.EthAddress().Hex())
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, zeroValue)
	})

	t.Run("signWhitelistERC20Cmd - should fail on invalid erc20 address", func(t *testing.T) {
		tx, err := evmSigner.signWhitelistERC20Cmd(ctx, txData, "")
		require.Nil(t, tx)
		require.ErrorContains(t, err, "invalid erc20 address")
	})

	t.Run("signWhitelistERC20Cmd - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		evmSigner.tss.Pause()

		// Call signWhitelistERC20Cmd
		tx, err := evmSigner.signWhitelistERC20Cmd(ctx, txData, sample.EthAddress().Hex())
		require.ErrorContains(t, err, "sign whitelist error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignMigrateTSSFundsCmd(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner := newTestSuite(t)

	// Setup txData struct
	cctx := getCCTX(t)

	txData, skip, err := NewOutboundData(ctx, cctx, 123, zerolog.Logger{})

	require.False(t, skip)
	require.NoError(t, err)

	t.Run("signMigrateTSSFundsCmd - should successfully sign", func(t *testing.T) {
		// Call signMigrateTSSFundsCmd
		tx, err := evmSigner.signMigrateTSSFundsCmd(ctx, txData)
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
	})

	t.Run("signMigrateTSSFundsCmd - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		evmSigner.tss.Pause()

		// Call signMigrateTSSFundsCmd
		tx, err := evmSigner.signMigrateTSSFundsCmd(ctx, txData)
		require.ErrorContains(t, err, "signMigrateTSSFundsCmd error")
		require.Nil(t, tx)
	})
}
