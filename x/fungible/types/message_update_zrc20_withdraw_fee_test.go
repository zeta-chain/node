package types_test

import (
	"testing"

	math "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestMsgUpdateZRC20WithdrawFee_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateZRC20WithdrawFee
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateZRC20WithdrawFee{
				Creator:        "invalid_address",
				Zrc20Address:   sample.EthAddress().String(),
				NewWithdrawFee: math.NewUint(1),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid new system contract address",
			msg: types.MsgUpdateZRC20WithdrawFee{
				Creator:        sample.AccAddress(),
				Zrc20Address:   "invalid_address",
				NewWithdrawFee: math.NewUint(1),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "both withdraw fee and gas limit nil",
			msg: types.MsgUpdateZRC20WithdrawFee{
				Creator:        sample.AccAddress(),
				Zrc20Address:   sample.EthAddress().String(),
				NewGasLimit:    math.Uint{},
				NewWithdrawFee: math.Uint{},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid message",
			msg: types.MsgUpdateZRC20WithdrawFee{
				Creator:        sample.AccAddress(),
				Zrc20Address:   sample.EthAddress().String(),
				NewWithdrawFee: math.NewUint(42),
				NewGasLimit:    math.NewUint(42),
			},
		},
		{
			name: "withdraw fee can be zero",
			msg: types.MsgUpdateZRC20WithdrawFee{
				Creator:        sample.AccAddress(),
				Zrc20Address:   sample.EthAddress().String(),
				NewWithdrawFee: math.ZeroUint(),
				NewGasLimit:    math.NewUint(42),
			},
		},
		{
			name: "withdraw fee can be nil",
			msg: types.MsgUpdateZRC20WithdrawFee{
				Creator:      sample.AccAddress(),
				Zrc20Address: sample.EthAddress().String(),
				NewGasLimit:  math.NewUint(42),
			},
		},
		{
			name: "gas limit can be zero",
			msg: types.MsgUpdateZRC20WithdrawFee{
				Creator:        sample.AccAddress(),
				Zrc20Address:   sample.EthAddress().String(),
				NewGasLimit:    math.ZeroUint(),
				NewWithdrawFee: math.NewUint(42),
			},
		},
		{
			name: "gas limit can be nil",
			msg: types.MsgUpdateZRC20WithdrawFee{
				Creator:        sample.AccAddress(),
				Zrc20Address:   sample.EthAddress().String(),
				NewWithdrawFee: math.NewUint(42),
			},
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

func TestMsgUpdateZRC20WithdrawFee_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateZRC20WithdrawFee
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateZRC20WithdrawFee{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateZRC20WithdrawFee{
				Creator: "invalid",
			},
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				assert.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				assert.Panics(t, func() {
					tt.msg.GetSigners()
				})
			}
		})
	}
}

func TestMsgUpdateZRC20WithdrawFee_Type(t *testing.T) {
	msg := types.MsgUpdateZRC20WithdrawFee{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.TypeMsgUpdateZRC20WithdrawFee, msg.Type())
}

func TestMsgUpdateZRC20WithdrawFee_Route(t *testing.T) {
	msg := types.MsgUpdateZRC20WithdrawFee{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateZRC20WithdrawFee_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateZRC20WithdrawFee{
		Creator: sample.AccAddress(),
	}
	assert.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
