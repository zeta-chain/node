package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestNewMsgMigrateTssFunds(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgMigrateTssFunds
		error bool
	}{
		{
			name: "invalid creator",
			msg: types.MsgMigrateTssFunds{
				Creator: "invalid_address",
				ChainId: common.EthChain().ChainId,
				Amount:  sdkmath.NewUintFromString("100000"),
			},
			error: true,
		},
		{
			name: "invalid chain id",
			msg: types.MsgMigrateTssFunds{
				Creator: "zeta15ruj2tc76pnj9xtw64utktee7cc7w6vzaes73z",
				ChainId: 999,
				Amount:  sdkmath.NewUintFromString("100000"),
			},
			error: true,
		},
		{
			name: "valid msg",
			msg: types.MsgMigrateTssFunds{
				Creator: "zeta15ruj2tc76pnj9xtw64utktee7cc7w6vzaes73z",
				ChainId: common.EthChain().ChainId,
				Amount:  sdkmath.NewUintFromString("100000"),
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
