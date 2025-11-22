package signer

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
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

		// Mock the digest to be signed
		digest := getWhitelistERC20Digest(t, evmSigner.Signer, txData, params)
		mockSignature(t, evmSigner.Signer, txData.nonce, digest)

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

		// Mock the digest to be signed
		digest := getMigrateERC20CustodyFundsDigest(t, evmSigner.Signer, txData, params)
		mockSignature(t, evmSigner.Signer, txData.nonce, digest)

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

		// Mock the digest to be signed
		digest := getUpdateERC20CustodyPauseStatusDigest(t, evmSigner.Signer, txData, params)
		mockSignature(t, evmSigner.Signer, txData.nonce, digest)

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

		// Mock the digest to be signed
		digest := getMigrateTssFundsDigest(t, evmSigner.Signer, txData)
		mockSignature(t, evmSigner.Signer, txData.nonce, digest)

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
	params := sample.EthAddress().Hex()

	require.NoError(t, err)
	require.False(t, skip)

	t.Run("signWhitelistERC20Cmd - should successfully sign", func(t *testing.T) {
		// Mock the digest to be signed
		digest := getWhitelistERC20Digest(t, evmSigner.Signer, txData, params)
		mockSignature(t, evmSigner.Signer, txData.nonce, digest)

		// Call signWhitelistERC20Cmd
		tx, err := evmSigner.signWhitelistERC20Cmd(ctx, txData, params)
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
		// use txData with a different nonce, digest will change
		txDataOther := *txData
		txDataOther.nonce++

		// Call signWhitelistERC20Cmd
		tx, err := evmSigner.signWhitelistERC20Cmd(ctx, &txDataOther, params)
		require.ErrorIs(t, err, ErrWaitForSignature)
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

		// Mock the digest to be signed
		digest := getMigrateERC20CustodyFundsDigest(t, evmSigner.Signer, txData, params)
		mockSignature(t, evmSigner.Signer, txData.nonce, digest)

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
		// use txData with a different nonce, digest will change
		txDataOther := *txData
		txDataOther.nonce++

		params := fmt.Sprintf(
			"%s,%s,%s",
			sample.EthAddress().Hex(),
			sample.EthAddress().Hex(),
			big.NewInt(100).String(),
		)

		// Call signWhitelistERC20Cmd
		tx, err := evmSigner.signMigrateERC20CustodyFundsCmd(ctx, txData, params)
		require.ErrorIs(t, err, ErrWaitForSignature)
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

		// Mock the digest to be signed
		digest := getUpdateERC20CustodyPauseStatusDigest(t, evmSigner.Signer, txData, params)
		mockSignature(t, evmSigner.Signer, txData.nonce, digest)

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

		// Mock the digest to be signed
		digest := getUpdateERC20CustodyPauseStatusDigest(t, evmSigner.Signer, txData, params)
		mockSignature(t, evmSigner.Signer, txData.nonce, digest)

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
		// use txData with a different nonce, digest will change
		txDataOther := *txData
		txDataOther.nonce++

		params := constant.OptionPause

		// Call signWhitelistERC20Cmd
		tx, err := evmSigner.signUpdateERC20CustodyPauseStatusCmd(ctx, &txDataOther, params)
		require.ErrorIs(t, err, ErrWaitForSignature)
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
		// Mock the digest to be signed
		digest := getMigrateTssFundsDigest(t, evmSigner.Signer, txData)
		mockSignature(t, evmSigner.Signer, txData.nonce, digest)

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
		// use txData with a different nonce, digest will change
		txDataOther := *txData
		txDataOther.nonce++

		// Call signMigrateTssFundsCmd
		tx, err := evmSigner.signMigrateTssFundsCmd(ctx, &txDataOther)
		require.ErrorIs(t, err, ErrWaitForSignature)
		require.Nil(t, tx)
	})
}

func getWhitelistERC20Digest(t *testing.T, signer *Signer, txData *OutboundData, params string) []byte {
	erc20 := ethcommon.HexToAddress(params)
	require.NotEqual(t, erc20, (ethcommon.Address{}))

	custodyAbi, err := erc20custody.ERC20CustodyMetaData.GetAbi()
	require.NoError(t, err)

	data, err := custodyAbi.Pack("whitelist", erc20)
	require.NoError(t, err)

	var (
		chainID = big.NewInt(txData.toChainID.Int64())
		to      = txData.to
		amount  = zeroValue
	)

	tx, err := newTx(chainID, data, to, amount, txData.gas, txData.nonce)
	require.NoError(t, err)

	return signer.evmClient.Signer().Hash(tx).Bytes()
}

func getMigrateERC20CustodyFundsDigest(t *testing.T, signer *Signer, txData *OutboundData, params string) []byte {
	paramsArray := strings.Split(params, ",")
	require.Len(t, paramsArray, 3)

	newCustody := ethcommon.HexToAddress(paramsArray[0])
	erc20 := ethcommon.HexToAddress(paramsArray[1])
	amount, ok := new(big.Int).SetString(paramsArray[2], 10)
	require.True(t, ok)

	custodyAbi, err := erc20custody.ERC20CustodyMetaData.GetAbi()
	require.NoError(t, err)

	data, err := custodyAbi.Pack("withdraw", newCustody, erc20, amount)
	require.NoError(t, err)

	var (
		chainID = big.NewInt(txData.toChainID.Int64())
		to      = txData.to
	)

	tx, err := newTx(chainID, data, to, zeroValue, txData.gas, txData.nonce)
	require.NoError(t, err)

	return signer.evmClient.Signer().Hash(tx).Bytes()
}

func getUpdateERC20CustodyPauseStatusDigest(t *testing.T, signer *Signer, txData *OutboundData, params string) []byte {
	custodyAbi, err := erc20custody.ERC20CustodyMetaData.GetAbi()
	require.NoError(t, err)

	// select the action
	var functionName string
	switch params {
	case constant.OptionPause:
		functionName = "pause"
	case constant.OptionUnpause:
		functionName = "unpause"
	default:
		require.Fail(t, "invalid params: %s", params)
	}

	data, err := custodyAbi.Pack(functionName)
	require.NoError(t, err)

	var (
		chainID = big.NewInt(txData.toChainID.Int64())
		to      = txData.to
		amount  = zeroValue
	)

	tx, err := newTx(chainID, data, to, amount, txData.gas, txData.nonce)
	require.NoError(t, err)

	return signer.evmClient.Signer().Hash(tx).Bytes()
}

func getMigrateTssFundsDigest(t *testing.T, signer *Signer, txData *OutboundData) []byte {
	chainID := big.NewInt(txData.toChainID.Int64())

	tx, err := newTx(chainID, nil, txData.to, txData.amount, txData.gas, txData.nonce)
	require.NoError(t, err)

	return signer.evmClient.Signer().Hash(tx).Bytes()
}
