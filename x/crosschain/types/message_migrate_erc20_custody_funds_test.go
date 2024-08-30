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

func TestNewMsgMigrateERC20CustodyFunds_ValidateBasic(t *testing.T) {
	keeper.SetConfig(false)
	tests := []struct {
		name  string
		msg   *types.MsgMigrateERC20CustodyFunds
		error bool
	}{
		{
			name: "invalid creator",
			msg: types.NewMsgMigrateERC20CustodyFunds(
				"invalid address",
				chains.DefaultChainsList()[0].ChainId,
				sample.EthAddress().String(),
				sample.EthAddress().String(),
				sdkmath.NewUintFromString("100000"),
			),
			error: true,
		},
		{
			name: "invalid amount",
			msg: types.NewMsgMigrateERC20CustodyFunds(
				sample.AccAddress(),
				chains.DefaultChainsList()[0].ChainId,
				sample.EthAddress().String(),
				sample.EthAddress().String(),
				sdkmath.NewUintFromString("0"),
			),
			error: true,
		},
		{
			name: "valid msg",
			msg: types.NewMsgMigrateERC20CustodyFunds(
				sample.AccAddress(),
				chains.DefaultChainsList()[0].ChainId,
				sample.EthAddress().String(),
				sample.EthAddress().String(),
				sdkmath.NewUintFromString("100000"),
			),
		},
		{
			name: "invalid erc20 address",
			msg: types.NewMsgMigrateERC20CustodyFunds(
				sample.AccAddress(),
				chains.DefaultChainsList()[0].ChainId,
				sample.EthAddress().String(),
				"invalid address",
				sdkmath.NewUintFromString("100000"),
			),
			error: true,
		},
		{
			name: "invalid new custody address",
			msg: types.NewMsgMigrateERC20CustodyFunds(
				sample.AccAddress(),
				chains.DefaultChainsList()[0].ChainId,
				"invalid address",
				sample.EthAddress().String(),
				sdkmath.NewUintFromString("100000"),
			),
			error: true,
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

func TestNewMsgMigrateERC20CustodyFunds_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgMigrateERC20CustodyFunds
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgMigrateERC20CustodyFunds{
				Creator:           signer,
				ChainId:           chains.DefaultChainsList()[0].ChainId,
				NewCustodyAddress: sample.EthAddress().String(),
				Erc20Address:      sample.EthAddress().String(),
				Amount:            sdkmath.NewUintFromString("100000"),
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgMigrateERC20CustodyFunds{
				Creator:           "invalid_address",
				ChainId:           chains.DefaultChainsList()[0].ChainId,
				NewCustodyAddress: sample.EthAddress().String(),
				Erc20Address:      sample.EthAddress().String(),
				Amount:            sdkmath.NewUintFromString("100000"),
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

func TestNewMsgMigrateERC20CustodyFunds_Type(t *testing.T) {
	msg := types.MsgMigrateERC20CustodyFunds{
		Creator:           sample.AccAddress(),
		ChainId:           chains.DefaultChainsList()[0].ChainId,
		NewCustodyAddress: sample.EthAddress().String(),
		Erc20Address:      sample.EthAddress().String(),
		Amount:            sdkmath.NewUintFromString("100000"),
	}
	require.Equal(t, types.TypeMsgMigrateERC20CustodyFunds, msg.Type())
}

func TestNewMsgMigrateERC20CustodyFunds_Route(t *testing.T) {
	msg := types.MsgMigrateERC20CustodyFunds{
		Creator:           sample.AccAddress(),
		ChainId:           chains.DefaultChainsList()[0].ChainId,
		NewCustodyAddress: sample.EthAddress().String(),
		Erc20Address:      sample.EthAddress().String(),
		Amount:            sdkmath.NewUintFromString("100000"),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestNewMsgMigrateERC20CustodyFunds_GetSignBytes(t *testing.T) {
	msg := types.MsgMigrateERC20CustodyFunds{
		Creator:           sample.AccAddress(),
		ChainId:           chains.DefaultChainsList()[0].ChainId,
		NewCustodyAddress: sample.EthAddress().String(),
		Erc20Address:      sample.EthAddress().String(),
		Amount:            sdkmath.NewUintFromString("100000"),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
