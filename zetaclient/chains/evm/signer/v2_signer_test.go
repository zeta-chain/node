package signer

import (
	"math/big"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestSigner_SignOutboundFromCCTXV2_ArbitraryCallCancels verifies that a V2 CCTX
// requesting an arbitrary call is short-circuited to SignCancel (a TSS-to-TSS
// zero-value self-transfer) rather than packed into GatewayEVM.execute.
func TestSigner_SignOutboundFromCCTXV2_ArbitraryCallCancels(t *testing.T) {
	ctx := makeCtx(t)
	evmSigner := newTestSuite(t)

	cctx := getCCTX(t)
	txData, skip, err := NewOutboundData(ctx, cctx, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	// Force the outbound's call options to request an arbitrary call.
	txData.callOptions.IsArbitraryCall = true

	// SignOutboundFromCCTXV2 should route arbitrary calls to SignCancel; mock
	// the signature for the cancel digest accordingly.
	digest := getCancelDigest(t, evmSigner.Signer, txData)
	mockSignature(t, evmSigner.Signer, txData.nonce, digest)

	tx, err := evmSigner.SignOutboundFromCCTXV2(cctx, txData)
	require.NoError(t, err)
	require.NotNil(t, tx)

	// Cancel produces a TSS self-transfer: to == TSS, value == 0.
	verifyTxBodyBasics(t, tx, evmSigner.tss.PubKey().AddressEVM(), txData.nonce, big.NewInt(0))
	require.Empty(t, tx.Data())
}
