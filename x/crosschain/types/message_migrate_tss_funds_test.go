package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestNewMsgMigrateTssFunds_ValidateBasic(t *testing.T) {
	keeper.SetConfig(false)
	tests := []struct {
		name  string
		msg   *types.MsgMigrateTssFunds
		error bool
	}{
		{
			name: "invalid creator",
			msg: types.NewMsgMigrateTssFunds(
				"invalid address",
				chains.DefaultChainsList()[0].ChainId,
				sdkmath.NewUintFromString("100000"),
			),
			error: true,
		},
		{
			name: "invalid amount",
			msg: types.NewMsgMigrateTssFunds(
				sample.AccAddress(),
				chains.DefaultChainsList()[0].ChainId,
				sdkmath.NewUintFromString("0"),
			),
			error: true,
		},
		{
			name: "valid msg",
			msg: types.NewMsgMigrateTssFunds(
				sample.AccAddress(),
				chains.DefaultChainsList()[0].ChainId,
				sdkmath.NewUintFromString("100000"),
			),
			error: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestNewMsgMigrateTssFunds_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgMigrateTssFunds
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgMigrateTssFunds{
				Creator: signer,
				ChainId: chains.DefaultChainsList()[0].ChainId,
				Amount:  sdkmath.NewUintFromString("100000"),
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgMigrateTssFunds{
				Creator: "invalid_address",
				ChainId: chains.DefaultChainsList()[0].ChainId,
				Amount:  sdkmath.NewUintFromString("100000"),
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

func TestNewMsgMigrateTssFunds_Type(t *testing.T) {
	msg := types.MsgMigrateTssFunds{
		Creator: sample.AccAddress(),
		ChainId: chains.DefaultChainsList()[0].ChainId,
		Amount:  sdkmath.NewUintFromString("100000"),
	}
	require.Equal(t, types.TypeMsgMigrateTssFunds, msg.Type())
}

func TestNewMsgMigrateTssFunds_Route(t *testing.T) {
	msg := types.MsgMigrateTssFunds{
		Creator: sample.AccAddress(),
		ChainId: chains.DefaultChainsList()[0].ChainId,
		Amount:  sdkmath.NewUintFromString("100000"),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestNewMsgMigrateTssFunds_GetSignBytes(t *testing.T) {
	msg := types.MsgMigrateTssFunds{
		Creator: sample.AccAddress(),
		ChainId: chains.DefaultChainsList()[0].ChainId,
		Amount:  sdkmath.NewUintFromString("100000"),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
