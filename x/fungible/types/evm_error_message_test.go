package types_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestEvmErrorMessage(t *testing.T) {
	t.Run("TestEvmErrorMessage", func(t *testing.T) {
		contractAddress := "0xE97Ac2CA30D30de65a6FE0Ab20EDC39a623c18df"
		msg := types.NewEvmErrorMessage("method", common.HexToAddress(contractAddress), "args", "message")
		msg.AddError("error_cause")
		msg.AddRevertReason("revert_reason")

		s, err := msg.ToJSON()
		require.NoError(t, err)

		require.Equal(t, `{"message":"message","method":"method","contract":"0xE97Ac2CA30D30de65a6FE0Ab20EDC39a623c18df","args":"args","error":"error_cause","revert_reason":"revert_reason"}`, s)
	})

}
