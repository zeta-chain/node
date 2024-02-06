package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
)

var _ sdk.Msg = &MsgRefundAbortedCCTX{}

func NewMsgRefundAbortedCCTX(creator string, cctxIndex string) *MsgRefundAbortedCCTX {
	return &MsgRefundAbortedCCTX{
		Creator:   creator,
		CctxIndex: cctxIndex,
	}
}

func (msg *MsgRefundAbortedCCTX) Route() string {
	return RouterKey
}

func (msg *MsgRefundAbortedCCTX) Type() string {
	return common.RefundAborted.String()
}

func (msg *MsgRefundAbortedCCTX) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRefundAbortedCCTX) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRefundAbortedCCTX) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if len(msg.CctxIndex) != 66 {
		return ErrInvalidCCTXIndex
	}
	return nil
}
