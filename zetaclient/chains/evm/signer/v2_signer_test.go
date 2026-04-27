package signer

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/coin"
)

// TestSigner_SignOutboundFromCCTXV2_ArbitraryCallCancels verifies that a V2
// gas-withdraw-and-call CCTX with IsArbitraryCall=true is short-circuited to
// SignCancel (a TSS-to-TSS zero-value self-transfer) rather than packed into
// GatewayEVM.execute.
func TestSigner_SignOutboundFromCCTXV2_ArbitraryCallCancels(t *testing.T) {
	ctx := makeCtx(t)
	evmSigner := newTestSuite(t)

	// Build a OutboundTypeGasWithdrawAndCall CCTX (CoinType=Gas with
	// IsCrossChainCall=true) — the GatewayEVM.execute dispatch path that must
	// be cancelled when IsArbitraryCall=true.
	cctx := getCCTX(t)
	cctx.InboundParams.IsCrossChainCall = true

	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	// Force the outbound's call options to request an arbitrary call.
	txData.callOptions.IsArbitraryCall = true

	digest := getCancelDigest(t, evmSigner.Signer, txData)
	mockSignature(t, evmSigner.Signer, txData.nonce, digest)

	tx, err := evmSigner.SignOutboundFromCCTXV2(cctx, txData)
	require.NoError(t, err)
	require.NotNil(t, tx)

	// Cancel produces a TSS self-transfer: to == TSS, value == 0, no calldata.
	verifyTxBodyBasics(t, tx, evmSigner.tss.PubKey().AddressEVM(), txData.nonce, big.NewInt(0))
	require.Empty(t, tx.Data())
}

// TestSigner_SignOutboundFromCCTXV2_NoAssetArbitraryCallCancels verifies that
// a V2 OutboundTypeCall (no-asset call from GatewayZEVM.call) with
// IsArbitraryCall=true is short-circuited to SignCancel — this is the path
// the live mainnet drain used.
func TestSigner_SignOutboundFromCCTXV2_NoAssetArbitraryCallCancels(t *testing.T) {
	ctx := makeCtx(t)
	evmSigner := newTestSuite(t)

	// Build a OutboundTypeCall CCTX (CoinType=NoAssetCall, amount=0,
	// PendingOutbound).
	cctx := getCCTX(t)
	cctx.InboundParams.CoinType = coin.CoinType_NoAssetCall
	cctx.GetCurrentOutboundParam().CoinType = coin.CoinType_NoAssetCall
	cctx.GetCurrentOutboundParam().Amount = math.ZeroUint()

	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	txData.callOptions.IsArbitraryCall = true

	digest := getCancelDigest(t, evmSigner.Signer, txData)
	mockSignature(t, evmSigner.Signer, txData.nonce, digest)

	tx, err := evmSigner.SignOutboundFromCCTXV2(cctx, txData)
	require.NoError(t, err)
	require.NotNil(t, tx)

	// Cancel produces a TSS self-transfer: to == TSS, value == 0, no calldata.
	verifyTxBodyBasics(t, tx, evmSigner.tss.PubKey().AddressEVM(), txData.nonce, big.NewInt(0))
	require.Empty(t, tx.Data())
}

// TestSigner_SignOutboundFromCCTXV2_PlainWithdrawNotCancelled is a regression
// guard for the case where GatewayZEVM.withdraw() emits Withdrawn with
// CallOptions.isArbitraryCall=true even though there is no payload. Plain V2
// withdraws must NOT be cancelled — they don't reach signGatewayExecute.
func TestSigner_SignOutboundFromCCTXV2_PlainWithdrawNotCancelled(t *testing.T) {
	ctx := makeCtx(t)
	evmSigner := newTestSuite(t)

	// CCTX is OutboundTypeGasWithdraw (CoinType=Gas, IsCrossChainCall=false).
	cctx := getCCTX(t)
	require.False(t, cctx.InboundParams.IsCrossChainCall)

	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	// Simulate the protocol-contract behavior: GatewayZEVM.withdraw() emits
	// isArbitraryCall=true on plain withdraws.
	txData.callOptions.IsArbitraryCall = true

	digest := getGasWithdrawDigest(t, evmSigner.Signer, txData)
	mockSignature(t, evmSigner.Signer, txData.nonce, digest)

	tx, err := evmSigner.SignOutboundFromCCTXV2(cctx, txData)
	require.NoError(t, err)
	require.NotNil(t, tx)

	// Gas withdraw sends to the receiver, not the TSS — i.e. NOT a self-transfer.
	verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
	require.NotEqual(t, evmSigner.tss.PubKey().AddressEVM(), *tx.To())
}
