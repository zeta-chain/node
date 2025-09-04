package utils

import (
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// RequireCCTXStatus checks if the cctx status is equal to the expected status
func RequireCCTXStatus(
	t require.TestingT,
	cctx *crosschaintypes.CrossChainTx,
	expected crosschaintypes.CctxStatus,
	msgAndArgs ...any,
) {
	msg := fmt.Sprintf(
		"cctx status is not %q cctx index %s, status: %s, status message %s, error: %s",
		expected.String(),
		cctx.Index,
		cctx.CctxStatus.Status.String(),
		cctx.CctxStatus.StatusMessage,
		cctx.CctxStatus.ErrorMessage,
	)

	require.NotNil(t, cctx.CctxStatus)
	require.Equal(t, expected, cctx.CctxStatus.Status, msg+errSuffix(msgAndArgs...))
}

// RequireTxSuccessful checks if the receipt status is successful.
// Currently, it accepts eth receipt, but we can make this more generic by using type assertion.
func RequireTxSuccessful(t require.TestingT, receipt *ethtypes.Receipt, msgAndArgs ...any) {
	msg := "receipt status is not successful: %s"
	require.Equal(
		t,
		ethtypes.ReceiptStatusSuccessful,
		receipt.Status,
		msg+errSuffix(msgAndArgs...),
		receipt.TxHash.String(),
	)
}

// RequiredTxFailed checks if the receipt status is failed.
// Currently, it accepts eth receipt, but we can make this more generic by using type assertion.
func RequiredTxFailed(t require.TestingT, receipt *ethtypes.Receipt, msgAndArgs ...any) {
	msg := "receipt status is not failed: %s"
	require.Equal(
		t,
		ethtypes.ReceiptStatusFailed,
		receipt.Status,
		msg+errSuffix(msgAndArgs...),
		receipt.TxHash.String(),
	)
}

func errSuffix(msgAndArgs ...any) string {
	if len(msgAndArgs) == 0 {
		return ""
	}

	template := "; " + msgAndArgs[0].(string)

	if len(msgAndArgs) == 1 {
		return template
	}

	return fmt.Sprintf(template, msgAndArgs[1:])
}
