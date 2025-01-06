package types_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestEvmErrorMessage(t *testing.T) {
	t.Run("TestEvmErrorMessage", func(t *testing.T) {
		address := sample.EthAddress()
		msg := types.EvmErrorMessage("errorMsg", "method", address, "args")
		msg = types.EvmErrorMessageAddErrorString(msg, "error_cause")
		msg = types.EvmErrorMessageAddRevertReason(msg, "revert_reason")

		require.Equal(t, fmt.Sprintf(
			"message:%s,method:%s,contract:%s,args:%s,error:%s,revertReason:%s",
			"errorMsg",
			"method",
			address.String(),
			"args",
			"error_cause",
			"revert_reason",
		), msg)
	})

}
