package signer

import (
	"fmt"
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

	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})

	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignAdminTx CmdWhitelistERC20", func(t *testing.T) {
		cmd := constant.CmdWhitelistERC20
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

	t.Run("SignAdminTx CmdMigrateERC20CustodyFunds", func(t *testing.T) {
		cmd := constant.CmdMigrateERC20CustodyFunds
		params := fmt.Sprintf(
			"%s,%s,%s",
			sample.EthAddress().Hex(),
			sample.EthAddress().Hex(),
			big.NewInt(100).String(),
		)
		// Call SignAdminTx
		tx, err := evmSigner.SignAdminTx(ctx, txData, cmd, params)
		require.NoError(t, err)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Revert tx calls erc20 custody contract with 0 gas token
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, big.NewInt(0))
	})

	t.Run("SignAdminTx CmdUpdateERC20CustodyPauseStatus", func(t *testing.T) {
		cmd := constant.CmdUpdateERC20CustodyPauseStatus
		params := constant.OptionPause
		// Call SignAdminTx
		tx, err := evmSigner.SignAdminTx(ctx, txData, cmd, params)
		require.NoError(t, err)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Revert tx calls erc20 custody contract with 0 gas token
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, big.NewInt(0))
	})

	t.Run("SignAdminTx CmdMigrateTssFunds", func(t *testing.T) {
		cmd := constant.CmdMigrateTssFunds
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

	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})

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

func TestSigner_SignMigrateERC20CustodyFundsCmd(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner := newTestSuite(t)

	// Setup txData struct
	cctx := getCCTX(t)

	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})

	require.NoError(t, err)
	require.False(t, skip)

	t.Run("signMigrateERC20CustodyFundsCmd - should successfully sign", func(t *testing.T) {
		// Call signWhitelistERC20Cmd

		params := fmt.Sprintf(
			"%s,%s,%s",
			sample.EthAddress().Hex(),
			sample.EthAddress().Hex(),
			big.NewInt(100).String(),
		)

		tx, err := evmSigner.signMigrateERC20CustodyFundsCmd(ctx, txData, params)
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, zeroValue)
	})

	t.Run("signMigrateERC20CustodyFundsCmd - should fail on invalid params", func(t *testing.T) {

		params := fmt.Sprintf("%s,%s", sample.EthAddress().Hex(), sample.EthAddress().Hex())

		_, err := evmSigner.signMigrateERC20CustodyFundsCmd(ctx, txData, params)
		require.ErrorContains(t, err, "invalid params")
	})

	t.Run("signMigrateERC20CustodyFundsCmd - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		evmSigner.tss.Pause()

		params := fmt.Sprintf(
			"%s,%s,%s",
			sample.EthAddress().Hex(),
			sample.EthAddress().Hex(),
			big.NewInt(100).String(),
		)

		// Call signWhitelistERC20Cmd
		tx, err := evmSigner.signMigrateERC20CustodyFundsCmd(ctx, txData, params)
		require.ErrorContains(t, err, "tss is paused")
		require.Nil(t, tx)
	})
}

func TestSigner_SignUpdateERC20CustodyPauseStatusCmd(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner := newTestSuite(t)

	// Setup txData struct
	cctx := getCCTX(t)

	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})

	require.False(t, skip)
	require.NoError(t, err)

	t.Run("signUpdateERC20CustodyPauseStatusCmd - should successfully sign for pause", func(t *testing.T) {
		// Call signWhitelistERC20Cmd

		params := constant.OptionPause

		tx, err := evmSigner.signUpdateERC20CustodyPauseStatusCmd(ctx, txData, params)
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, zeroValue)
	})

	t.Run("signUpdateERC20CustodyPauseStatusCmd - should successfully sign for unpause", func(t *testing.T) {
		// Call signWhitelistERC20Cmd

		params := constant.OptionUnpause

		tx, err := evmSigner.signUpdateERC20CustodyPauseStatusCmd(ctx, txData, params)
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, zeroValue)
	})

	t.Run("signUpdateERC20CustodyPauseStatusCmd - should fail on invalid params", func(t *testing.T) {
		params := "invalid"

		_, err := evmSigner.signUpdateERC20CustodyPauseStatusCmd(ctx, txData, params)
		require.ErrorContains(t, err, "invalid params")
	})

	t.Run("signUpdateERC20CustodyPauseStatusCmd - should fail on empty params", func(t *testing.T) {
		params := ""

		_, err := evmSigner.signUpdateERC20CustodyPauseStatusCmd(ctx, txData, params)
		require.ErrorContains(t, err, "invalid params")
	})

	t.Run("signUpdateERC20CustodyPauseStatusCmd - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		evmSigner.tss.Pause()

		params := constant.OptionPause

		// Call signWhitelistERC20Cmd
		tx, err := evmSigner.signUpdateERC20CustodyPauseStatusCmd(ctx, txData, params)
		require.ErrorContains(t, err, "tss is paused")
		require.Nil(t, tx)
	})
}

func TestSigner_SignMigrateTssFundsCmd(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner := newTestSuite(t)

	// Setup txData struct
	cctx := getCCTX(t)

	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})

	require.False(t, skip)
	require.NoError(t, err)

	t.Run("signMigrateTssFundsCmd - should successfully sign", func(t *testing.T) {
		// Call signMigrateTssFundsCmd
		tx, err := evmSigner.signMigrateTssFundsCmd(ctx, txData)
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
	})

	t.Run("signMigrateTssFundsCmd - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		evmSigner.tss.Pause()

		// Call signMigrateTssFundsCmd
		tx, err := evmSigner.signMigrateTssFundsCmd(ctx, txData)
		require.ErrorContains(t, err, "signMigrateTssFundsCmd error")
		require.Nil(t, tx)
	})
}
