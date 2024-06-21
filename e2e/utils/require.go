package utils

import (
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// RequireCCTXStatus checks if the cctx status is equal to the expected status
func RequireCCTXStatus(
	t require.TestingT,
	cctx *crosschaintypes.CrossChainTx,
	expected crosschaintypes.CctxStatus,
	msgAndArgs ...any,
) {
	msg := fmt.Sprintf("cctx status is not %q", expected.String())

	require.NotNil(t, cctx.CctxStatus)
	require.Equal(t, expected, cctx.CctxStatus.Status, msg+errSuffix(msgAndArgs...))
}

// RequireReceiptApproved checks if the receipt status is successful.
func RequireReceiptApproved(t require.TestingT, receipt *ethtypes.Receipt, msgAndArgs ...any) {
	msg := "receipt status is not successful"
	require.Equal(t, ethtypes.ReceiptStatusSuccessful, receipt.Status, msg+errSuffix(msgAndArgs...))
}

// RequireReceiptFailed checks if the receipt status is failed
func RequireReceiptFailed(t require.TestingT, receipt *ethtypes.Receipt, msgAndArgs ...any) {
	msg := "receipt status is not successful"
	require.Equal(t, ethtypes.ReceiptStatusFailed, receipt.Status, msg+errSuffix(msgAndArgs...))
}

func errSuffix(msgAndArgs ...any) string {
	if len(msgAndArgs) == 0 {
		return ""
	}

	template := "; " + msgAndArgs[0].(string)

	return fmt.Sprintf(template, msgAndArgs[1:])
}
