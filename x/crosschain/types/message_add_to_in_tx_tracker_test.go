package types_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgAddToInTxTracker_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgAddToInTxTracker
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgAddToInTxTracker{
				Creator:  "invalid_address",
				ChainId:  common.GoerliChain().ChainId,
				TxHash:   "hash",
				CoinType: common.CoinType_Gas,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid chain id",
			msg: types.MsgAddToInTxTracker{
				Creator:  sample.AccAddress(),
				ChainId:  42,
				TxHash:   "hash",
				CoinType: common.CoinType_Gas,
			},
			err: errorsmod.Wrapf(types.ErrInvalidChainID, "chain id (%d)", 42),
		},
		{
			name: "valid",
			msg: types.MsgAddToInTxTracker{
				Creator:  sample.AccAddress(),
				ChainId:  common.GoerliChain().ChainId,
				TxHash:   "hash",
				CoinType: common.CoinType_Gas,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
