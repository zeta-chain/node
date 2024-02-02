package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
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
				ChainId: common.DefaultChainsList()[0].ChainId,
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
				ChainId: common.DefaultChainsList()[0].ChainId,
				Amount:  sdkmath.NewUintFromString("100000"),
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
