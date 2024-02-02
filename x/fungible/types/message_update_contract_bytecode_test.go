package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestMsgUpdateContractBytecode_ValidateBasic(t *testing.T) {
	tt := []struct {
		name      string
		msg       types.MsgUpdateContractBytecode
		wantError bool
	}{
		{
			name: "valid",
			msg: types.MsgUpdateContractBytecode{
				Creator:         sample.AccAddress(),
				ContractAddress: sample.EthAddress().Hex(),
				NewCodeHash:     sample.Hash().Hex(),
			},
			wantError: false,
		},
		{
			name: "invalid creator",
			msg: types.MsgUpdateContractBytecode{
				Creator:         "invalid",
				ContractAddress: sample.EthAddress().Hex(),
				NewCodeHash:     sample.Hash().Hex(),
			},
			wantError: true,
		},
		{
			name: "invalid contract address",
			msg: types.MsgUpdateContractBytecode{
				Creator:         sample.AccAddress(),
				ContractAddress: "invalid",
				NewCodeHash:     sample.Hash().Hex(),
			},
			wantError: true,
		},
		{
			name: "invalid new code hash",
			msg: types.MsgUpdateContractBytecode{
				Creator:         sample.AccAddress(),
				ContractAddress: sample.EthAddress().Hex(),
				NewCodeHash:     "invalid",
			},
			wantError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
