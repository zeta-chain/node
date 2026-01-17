package types_test

import (
	"testing"

	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/ptr"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestMsgDeployFungibleCoinZRC4_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgDeployFungibleCoinZRC20
		err  error
	}{
		{
			name: "invalid address",
			msg: types.NewMsgDeployFungibleCoinZRC20(
				"invalid_address",
				"test erc20",
				1,
				6,
				"test",
				"test",
				coin.CoinType_ERC20,
				10,
				ptr.Ptr(sdkmath.NewUint(1000)),
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid gas limit",
			msg: types.NewMsgDeployFungibleCoinZRC20(
				sample.AccAddress(),
				"test erc20",
				1,
				6,
				"test",
				"test",
				coin.CoinType_ERC20,
				-1,
				ptr.Ptr(sdkmath.NewUint(1000)),
			),
			err: sdkerrors.ErrInvalidGasLimit,
		},
		{
			name: "invalid decimals",
			msg: types.NewMsgDeployFungibleCoinZRC20(
				sample.AccAddress(),
				"test erc20",
				1,
				78,
				"test",
				"test",
				coin.CoinType_ERC20,
				10,
				ptr.Ptr(sdkmath.NewUint(1000)),
			),
			err: cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "decimals must be less than 78"),
		},
		{
			name: "nil liquidity cap",
			msg: &types.MsgDeployFungibleCoinZRC20{
				Creator:        sample.AccAddress(),
				ERC20:          "test erc20",
				ForeignChainId: 1,
				Decimals:       6,
				Name:           "test",
				Symbol:         "test",
				CoinType:       coin.CoinType_ERC20,
				GasLimit:       10,
			},
		},
		{
			name: "nil liquidity cap inner",
			msg: &types.MsgDeployFungibleCoinZRC20{
				Creator:        sample.AccAddress(),
				ERC20:          "test erc20",
				ForeignChainId: 1,
				Decimals:       6,
				Name:           "test",
				Symbol:         "test",
				CoinType:       coin.CoinType_ERC20,
				GasLimit:       10,
				LiquidityCap:   &sdkmath.Uint{},
			},
			err: cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "liquidity cap is nil"),
		},
		{
			name: "valid message",
			msg: types.NewMsgDeployFungibleCoinZRC20(
				sample.AccAddress(),
				"test erc20",
				1,
				6,
				"test",
				"test",
				coin.CoinType_ERC20,
				10,
				ptr.Ptr(sdkmath.NewUint(1000)),
			),
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

func TestMsgDeployFungibleCoinZRC4_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgDeployFungibleCoinZRC20
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgDeployFungibleCoinZRC20{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgDeployFungibleCoinZRC20{
				Creator: "invalid",
			},
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

func TestMsgDeployFungibleCoinZRC4_Type(t *testing.T) {
	msg := types.MsgDeployFungibleCoinZRC20{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgDeployFungibleCoinZRC20, msg.Type())
}

func TestMsgDeployFungibleCoinZRC4_Route(t *testing.T) {
	msg := types.MsgDeployFungibleCoinZRC20{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgDeployFungibleCoinZRC4_GetSignBytes(t *testing.T) {
	msg := types.MsgDeployFungibleCoinZRC20{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
