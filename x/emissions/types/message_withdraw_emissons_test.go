package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
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
