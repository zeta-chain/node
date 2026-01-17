package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/sdkconfig"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestMsgWhitelistAsset_ValidateBasic(t *testing.T) {
	sdkconfig.SetDefault(false)
	tests := []struct {
		name  string
		msg   *types.MsgWhitelistAsset
		error bool
	}{
		{
			name: "invalid creator",
			msg: types.NewMsgWhitelistAsset(
				"invalid_address",
				"0x0",
				1,
				"name",
				"symbol",
				6,
				10,
				sdkmath.NewUint(100),
			),
			error: true,
		},
		{
			name: "invalid asset",
			msg: types.NewMsgWhitelistAsset(
				sample.AccAddress(),
				"",
				1,
				"name",
				"symbol",
				6,
				10,
				sdkmath.NewUint(100),
			),
			error: true,
		},
		{
			name: "invalid decimals",
			msg: types.NewMsgWhitelistAsset(
				sample.AccAddress(),
				sample.EthAddress().Hex(),
				1,
				"name",
				"symbol",
				130,
				10,
				sdkmath.NewUint(100),
			),
			error: true,
		},
		{
			name: "invalid gas limit",
			msg: types.NewMsgWhitelistAsset(
				sample.AccAddress(),
				sample.EthAddress().Hex(),
				1,
				"name",
				"symbol",
				6,
				-10,
				sdkmath.NewUint(100),
			),
			error: true,
		},
		{
			name: "evm asset address with invalid checksum format",
			msg: types.NewMsgWhitelistAsset(
				sample.AccAddress(),
				"0x5a4f260a7d716c859a2736151cb38b9c58c32c64",
				1,
				"name",
				"symbol",
				6,
				10,
				sdkmath.NewUint(100),
			),
			error: true,
		},
		{
			name: "invalid liquidity cap",
			msg: &types.MsgWhitelistAsset{
				Creator:      sample.AccAddress(),
				AssetAddress: "0x5a4f260A7D716c859A2736151cB38b9c58C32c64",
				ChainId:      1,
				Name:         "name",
				Symbol:       "symbol",
				Decimals:     6,
				GasLimit:     10,
			},
			error: true,
		},
		{
			name: "valid message with evm asset address",
			msg: types.NewMsgWhitelistAsset(
				sample.AccAddress(),
				"0x5a4f260A7D716c859A2736151cB38b9c58C32c64",
				1,
				"name",
				"symbol",
				6,
				10,
				sdkmath.NewUint(100),
			),
			error: false,
		},
		{
			name: "valid message with solana asset address",
			msg: types.NewMsgWhitelistAsset(
				sample.AccAddress(),
				"Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr",
				1,
				"name",
				"symbol",
				6,
				10,
				sdkmath.NewUint(100),
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

func TestMsgWhitelistAsset_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgWhitelistAsset
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgWhitelistAsset{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgWhitelistAsset{
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

func TestMsgWhitelistAsset_Type(t *testing.T) {
	msg := types.MsgWhitelistAsset{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgWhitelistAsset, msg.Type())
}

func TestMsgWhitelistAsset_Route(t *testing.T) {
	msg := types.MsgWhitelistAsset{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgWhitelistAsset_GetSignBytes(t *testing.T) {
	msg := types.MsgWhitelistAsset{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
