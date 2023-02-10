package types

//import (
//	sdk "github.com/cosmos/cosmos-sdk/types"
//)
//
//var _ sdk.Msg = &MsgSetSupportedChains{}
//
//func (msg *MsgSetSupportedChains) Route() string {
//	return RouterKey
//}
//
//func (msg *MsgSetSupportedChains) Type() string {
//	return "SetSupportedChains"
//}
//
//func (msg *MsgSetSupportedChains) GetSigners() []sdk.AccAddress {
//	creator, err := sdk.AccAddressFromBech32(msg.Creator)
//	if err != nil {
//		panic(err)
//	}
//	return []sdk.AccAddress{creator}
//}
//
//func (msg *MsgSetSupportedChains) GetSignBytes() []byte {
//	bz := ModuleCdc.MustMarshalJSON(msg)
//	return sdk.MustSortJSON(bz)
//}
//
//func (msg *MsgSetSupportedChains) ValidateBasic() error {
//	// TODO :ADD validation
//	return nil
//}
