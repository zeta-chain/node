package signer

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/constant"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
	"math/big"
	"testing"
)

func TestSigner_SignAdminTx(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner, err := getNewEvmSigner(nil)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver(t, nil)
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(ctx, cctx, mockObserver, 123, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignAdminTx CmdWhitelistERC20", func(t *testing.T) {
		cmd := constant.CmdWhitelistERC20
		params := ConnectorAddress.Hex()
		// Call SignAdminTx
		tx, err := evmSigner.SignAdminTx(ctx, txData, cmd, params)
		require.NoError(t, err)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

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
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

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

	txData, skip, err := NewOutboundData(ctx, cctx, mockObserver, 123, zerolog.Logger{})
	require.NoError(t, err)
	require.False(t, skip)

	t.Run("signWhitelistERC20Cmd - should successfully sign", func(t *testing.T) {
		// Call signWhitelistERC20Cmd
		tx, err := evmSigner.signWhitelistERC20Cmd(ctx, txData, sample.EthAddress().Hex())
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

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
		tss.Pause()

		// Call signWhitelistERC20Cmd
		tx, err := evmSigner.signWhitelistERC20Cmd(ctx, txData, sample.EthAddress().Hex())
		require.ErrorContains(t, err, "sign whitelist error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignMigrateERC20CustodyFundsCmd(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)

	mockObserver, err := getNewEvmChainObserver(t, tss)
	require.NoError(t, err)

	txData, skip, err := NewOutboundData(ctx, cctx, mockObserver, 123, zerolog.Logger{})
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
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

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
		tss.Pause()

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
	txData, skip, err := NewOutboundData(ctx, cctx, mockObserver, 123, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("signMigrateTssFundsCmd - should successfully sign", func(t *testing.T) {
		// Call signMigrateTssFundsCmd
		tx, err := evmSigner.signMigrateTssFundsCmd(ctx, txData)
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
	})

	t.Run("signMigrateTssFundsCmd - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call signMigrateTssFundsCmd
		tx, err := evmSigner.signMigrateTssFundsCmd(ctx, txData)
		require.ErrorContains(t, err, "signMigrateTssFundsCmd error")
		require.Nil(t, tx)
	})
}
