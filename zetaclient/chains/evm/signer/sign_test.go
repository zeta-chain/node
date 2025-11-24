package signer

import (
	"context"
	"math/big"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	connectorevm "github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.base.sol"
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

	// Mock the digest to be signed
	digest := getConnectorOnReceiveDigest(t, evmSigner.Signer, txData)
	mockSignature(t, evmSigner.Signer, txData.nonce, digest)

	t.Run("SignConnectorOnReceive - should successfully sign", func(t *testing.T) {
		// Call SignConnectorOnReceive
		tx, err := evmSigner.SignConnectorOnReceive(txData)
		require.NoError(t, err)

		// Verify Signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())
	})

	t.Run("SignConnectorOnReceive - should fail if signature is not available", func(t *testing.T) {
		// use txData with a different nonce, digest will change
		txDataOther := *txData
		txDataOther.nonce++

		// Call SignConnectorOnReceive
		tx, err := evmSigner.SignConnectorOnReceive(&txDataOther)
		require.ErrorIs(t, err, ErrWaitForSignature)
		require.Nil(t, tx)
	})

	t.Run("SignOutbound - should successfully sign LegacyTx", func(t *testing.T) {
		// Call SignOutbound
		tx, err := evmSigner.SignConnectorOnReceive(txData)
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
		tx, err := evmSigner.SignConnectorOnReceive(txData)
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

	// Mock the digest to be signed
	digest := getConnectorOnRevertDigest(t, evmSigner.Signer, txData)
	mockSignature(t, evmSigner.Signer, txData.nonce, digest)

	t.Run("SignConnectorOnRevert - should successfully sign", func(t *testing.T) {
		// Call SignConnectorOnRevert
		tx, err := evmSigner.SignConnectorOnRevert(txData)
		require.NoError(t, err)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Revert tx calls connector contract with 0 gas token
		verifyTxBodyBasics(t, tx, evmSigner.zetaConnectorAddress, txData.nonce, big.NewInt(0))
	})
	t.Run("SignConnectorOnRevert - should fail if signature is not available", func(t *testing.T) {
		// use txData with a different nonce, digest will change
		txDataOther := *txData
		txDataOther.nonce++

		// Call SignConnectorOnRevert
		tx, err := evmSigner.SignConnectorOnRevert(&txDataOther)
		require.ErrorIs(t, err, ErrWaitForSignature)
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

	// Mock the digest to be signed
	digest := getCancelDigest(t, evmSigner.Signer, txData)
	mockSignature(t, evmSigner.Signer, txData.nonce, digest)

	t.Run("SignCancel - should successfully sign", func(t *testing.T) {
		// Call SignCancel
		tx, err := evmSigner.SignCancel(txData)
		require.NoError(t, err)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Cancel tx sends 0 gas token to TSS self address
		verifyTxBodyBasics(t, tx, evmSigner.tss.PubKey().AddressEVM(), txData.nonce, big.NewInt(0))
	})
	t.Run("SignCancel - should fail if signature is not available", func(t *testing.T) {
		// use txData with a different nonce, digest will change
		txDataOther := *txData
		txDataOther.nonce++

		// Call SignCancel
		tx, err := evmSigner.SignCancel(&txDataOther)
		require.ErrorIs(t, err, ErrWaitForSignature)
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

	// Mock the digest to be signed
	digest := getGasWithdrawDigest(t, evmSigner.Signer, txData)
	mockSignature(t, evmSigner.Signer, txData.nonce, digest)

	t.Run("SignGasWithdraw - should successfully sign", func(t *testing.T) {
		// Call SignGasWithdraw
		tx, err := evmSigner.SignGasWithdraw(txData)
		require.NoError(t, err)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
	})
	t.Run("SignGasWithdraw - should fail if signature is not available", func(t *testing.T) {
		// use txData with a different nonce, digest will change
		txDataOther := *txData
		txDataOther.nonce++

		// Call SignGasWithdraw
		tx, err := evmSigner.SignGasWithdraw(&txDataOther)
		require.ErrorIs(t, err, ErrWaitForSignature)
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

	// Mock the digest to be signed
	digest := getERC20WithdrawDigest(t, evmSigner.Signer, txData)
	mockSignature(t, evmSigner.Signer, txData.nonce, digest)

	t.Run("SignERC20WithdrawTx - should successfully sign", func(t *testing.T) {
		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20Withdraw(txData)
		require.NoError(t, err)

		// Verify tx signature
		verifyTxSender(t, tx, evmSigner.tss.PubKey().AddressEVM(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Withdraw tx calls erc20 custody contract with 0 gas token
		verifyTxBodyBasics(t, tx, evmSigner.er20CustodyAddress, txData.nonce, big.NewInt(0))
	})

	t.Run("SignERC20WithdrawTx - should fail if signature is not available", func(t *testing.T) {
		// use txData with a different nonce, digest will change
		txDataOther := *txData
		txDataOther.nonce++

		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20Withdraw(&txDataOther)
		require.ErrorIs(t, err, ErrWaitForSignature)
		require.Nil(t, tx)
	})
}

func getGasWithdrawDigest(t *testing.T, signer *Signer, txData *OutboundData) []byte {
	var (
		chainID = big.NewInt(signer.Chain().ChainId)
		to      = txData.to
		amount  = txData.amount
	)

	tx, err := newTx(chainID, nil, to, amount, txData.gas, txData.nonce)
	require.NoError(t, err)

	return signer.evmClient.Signer().Hash(tx).Bytes()
}

func getERC20WithdrawDigest(t *testing.T, signer *Signer, txData *OutboundData) []byte {
	erc20CustodyV1ABI, err := erc20custody.ERC20CustodyMetaData.GetAbi()
	require.NoError(t, err)

	data, err := erc20CustodyV1ABI.Pack("withdraw", txData.to, txData.asset, txData.amount)
	require.NoError(t, err)

	var (
		chainID = big.NewInt(txData.toChainID.Int64())
		to      = signer.er20CustodyAddress
		amount  = zeroValue
	)

	tx, err := newTx(chainID, data, to, amount, txData.gas, txData.nonce)
	require.NoError(t, err)

	return signer.evmClient.Signer().Hash(tx).Bytes()
}

func getConnectorOnReceiveDigest(t *testing.T, signer *Signer, txData *OutboundData) []byte {
	zetaConnectorABI, err := connectorevm.ZetaConnectorBaseMetaData.GetAbi()
	require.NoError(t, err)

	data, err := zetaConnectorABI.Pack("onReceive",
		txData.sender.Bytes(),
		txData.srcChainID,
		txData.to,
		txData.amount,
		txData.message,
		txData.cctxIndex)
	require.NoError(t, err)

	var (
		chainID = big.NewInt(signer.Chain().ChainId)
		to      = signer.zetaConnectorAddress
		amount  = zeroValue
	)

	tx, err := newTx(chainID, data, to, amount, txData.gas, txData.nonce)
	require.NoError(t, err)

	return signer.evmClient.Signer().Hash(tx).Bytes()
}

func getConnectorOnRevertDigest(t *testing.T, signer *Signer, txData *OutboundData) []byte {
	zetaConnectorABI, err := connectorevm.ZetaConnectorBaseMetaData.GetAbi()
	require.NoError(t, err)

	data, err := zetaConnectorABI.Pack("onRevert",
		txData.sender,
		txData.srcChainID,
		txData.to.Bytes(),
		txData.toChainID,
		txData.amount,
		txData.message,
		txData.cctxIndex)
	require.NoError(t, err)

	var (
		chainID = big.NewInt(signer.Chain().ChainId)
		to      = signer.zetaConnectorAddress
		amount  = zeroValue
	)

	tx, err := newTx(chainID, data, to, amount, txData.gas, txData.nonce)
	require.NoError(t, err)

	return signer.evmClient.Signer().Hash(tx).Bytes()
}

func getCancelDigest(t *testing.T, signer *Signer, txData *OutboundData) []byte {
	var (
		chainID = big.NewInt(signer.Chain().ChainId)
		to      = signer.TSS().PubKey().AddressEVM()
		amount  = zeroValue
	)

	tx, err := newTx(chainID, nil, to, amount, txData.gas, txData.nonce)
	require.NoError(t, err)

	return signer.evmClient.Signer().Hash(tx).Bytes()
}

func mockSignature(t *testing.T, signer *Signer, nonce uint64, digest []byte) {
	ctx := context.Background()
	chainID := signer.Chain().ChainId

	// add digest to cache
	signer.GetSignatureOrAddDigest(nonce, digest)

	// mock a batch that contains the digest
	batch := base.NewTSSKeysignBatch()
	batch.AddKeysignInfo(nonce, *base.NewTSSKeysignInfo(digest, [65]byte{}))

	// sign
	sigs, err := signer.TSS().SignBatch(ctx, batch.Digests(), 1, nonce, chainID)
	require.NoError(t, err)

	// add signatures to cache
	signer.AddBatchSignatures(*batch, sigs)
}
