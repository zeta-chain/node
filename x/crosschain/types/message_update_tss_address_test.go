package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/testutil/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMessageUpdateTssAddress_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   crosschaintypes.MsgUpdateTssAddress
		error bool
	}{
		{
			name: "invalid creator",
			msg: crosschaintypes.MsgUpdateTssAddress{
				Creator:   "invalid_address",
				TssPubkey: "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p",
			},
			error: true,
		},
		{
			name: "invalid pubkey",
			msg: crosschaintypes.MsgUpdateTssAddress{
				Creator:   "zeta15ruj2tc76pnj9xtw64utktee7cc7w6vzaes73z",
				TssPubkey: "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm",
			},
			error: true,
		},
		{
			name: "valid msg",
			msg: crosschaintypes.MsgUpdateTssAddress{
				Creator:   "zeta15ruj2tc76pnj9xtw64utktee7cc7w6vzaes73z",
				TssPubkey: "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p",
			},
			error: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keeper.SetConfig(false)
			err := tt.msg.ValidateBasic()
			if tt.error {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
