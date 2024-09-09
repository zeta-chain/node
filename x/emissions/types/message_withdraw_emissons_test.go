package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
)

func TestMsgWithdrawEmission_ValidateBasic(t *testing.T) {
	t.Run("invalid creator address", func(t *testing.T) {
		msg := emissionstypes.NewMsgWithdrawEmissions("invalid_address", sample.IntInRange(1, 100))
		err := msg.ValidateBasic()
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})

	t.Run("invalid negative amount", func(t *testing.T) {
		msg := emissionstypes.NewMsgWithdrawEmissions(sample.AccAddress(), sample.IntInRange(-100, -1))
		err := msg.ValidateBasic()
		require.ErrorIs(t, err, emissionstypes.ErrInvalidAmount)
	})

	t.Run("invalid zero amount", func(t *testing.T) {
		msg := emissionstypes.NewMsgWithdrawEmissions(sample.AccAddress(), sdkmath.ZeroInt())
		err := msg.ValidateBasic()
		require.ErrorIs(t, err, emissionstypes.ErrInvalidAmount)
	})

	t.Run("invalid nil amount", func(t *testing.T) {
		msg := emissionstypes.NewMsgWithdrawEmissions(sample.AccAddress(), sdkmath.Int{})
		err := msg.ValidateBasic()
		require.ErrorIs(t, err, emissionstypes.ErrInvalidAmount)
	})

	t.Run("valid withdraw message", func(t *testing.T) {
		msg := emissionstypes.NewMsgWithdrawEmissions(sample.AccAddress(), sample.IntInRange(1, 100))
		err := msg.ValidateBasic()
		require.NoError(t, err)
	})
}

func TestMsgWithdrawEmission_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *emissionstypes.MsgWithdrawEmission
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    emissionstypes.NewMsgWithdrawEmissions(signer, sample.IntInRange(1, 100)),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    emissionstypes.NewMsgWithdrawEmissions("invalid", sample.IntInRange(1, 100)),
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				require.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				require.Panics(t, func() {
					tt.msg.GetSigners()
				})
			}
		})
	}
}

func TestMsgWithdrawEmission_Type(t *testing.T) {
	msg := emissionstypes.NewMsgWithdrawEmissions(sample.AccAddress(), sample.IntInRange(1, 100))
	require.Equal(t, emissionstypes.MsgWithdrawEmissionType, msg.Type())
}

func TestMsgWithdrawEmission_Route(t *testing.T) {
	msg := emissionstypes.NewMsgWithdrawEmissions(sample.AccAddress(), sample.IntInRange(1, 100))
	require.Equal(t, emissionstypes.RouterKey, msg.Route())
}

func TestMsgWithdrawEmission_GetSignBytes(t *testing.T) {
	msg := emissionstypes.NewMsgWithdrawEmissions(sample.AccAddress(), sample.IntInRange(1, 100))
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
