package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgEnableCCTXFlags_ValidateBasic(t *testing.T) {
	tt := []struct {
		name string
		msg  *types.MsgEnableCCTXFlags
		err  require.ErrorAssertionFunc
	}{
		{
			name: "invalid creator address",
			msg:  types.NewMsgEnableCCTXFlags("invalid", true, true),
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid creator address")
			},
		},
		{
			name: "valid",
			msg:  types.NewMsgEnableCCTXFlags(sample.AccAddress(), true, true),
			err:  require.NoError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.err(t, tc.msg.ValidateBasic())
		})
	}
}
