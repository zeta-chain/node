package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMessageUpdateTssAddress_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgUpdateTssAddress
		error bool
	}{
		{
			name: "invalid creator",
			msg: types.MsgUpdateTssAddress{
				Creator:   "invalid_address",
				TssPubkey: "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p",
			},
			error: true,
		},
		{
			name: "invalid pubkey",
			msg: types.MsgUpdateTssAddress{
				Creator:   "zeta15ruj2tc76pnj9xtw64utktee7cc7w6vzaes73z",
				TssPubkey: "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm",
			},
			error: true,
		},
		{
			name: "valid msg",
			msg: types.MsgUpdateTssAddress{
				Creator:   "zeta15ruj2tc76pnj9xtw64utktee7cc7w6vzaes73z",
				TssPubkey: "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p",
			},
			error: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observerTypes.SetConfig(false)
			err := tt.msg.ValidateBasic()
			if tt.error {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}
