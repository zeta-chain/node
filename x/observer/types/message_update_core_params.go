package types

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
)

const TypeMsgUpdateClientParams = "update_client_params"

var _ sdk.Msg = &MsgUpdateCoreParams{}

func NewMsgUpdateCoreParams(creator string, coreParams *CoreParams) *MsgUpdateCoreParams {
	return &MsgUpdateCoreParams{
		Creator:    creator,
		CoreParams: coreParams,
	}
}

func (msg *MsgUpdateCoreParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdateCoreParams) Type() string {
	return TypeMsgUpdateClientParams
}

func (msg *MsgUpdateCoreParams) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateCoreParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCoreParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return ValidateCoreParams(msg.CoreParams)
}

// ValidateCoreParams performs some basic checks on core params
func ValidateCoreParams(params *CoreParams) error {
	if params == nil {
		return fmt.Errorf("core params cannot be nil")
	}
	chain := common.GetChainFromChainID(params.ChainId)
	if chain == nil {
		return fmt.Errorf("ChainId %d not supported", params.ChainId)
	}
	// zeta chain skips the rest of the checks for now
	if chain.IsZetaChain() {
		return nil
	}

	if params.ConfirmationCount == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "ConfirmationCount must be greater than 0")
	}
	if params.GasPriceTicker <= 0 || params.GasPriceTicker > 300 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "GasPriceTicker %d out of range", params.GasPriceTicker)
	}
	if params.InTxTicker <= 0 || params.InTxTicker > 300 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "InTxTicker %d out of range", params.InTxTicker)
	}
	if params.OutTxTicker <= 0 || params.OutTxTicker > 300 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "OutTxTicker %d out of range", params.OutTxTicker)
	}
	if params.OutboundTxScheduleInterval == 0 || params.OutboundTxScheduleInterval > 100 { // 600 secs
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "OutboundTxScheduleInterval %d out of range", params.OutboundTxScheduleInterval)
	}
	if params.OutboundTxScheduleLookahead == 0 || params.OutboundTxScheduleLookahead > 500 { // 500 cctxs
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "OutboundTxScheduleLookahead %d out of range", params.OutboundTxScheduleLookahead)
	}

	// chain type specific checks
	if common.IsBitcoinChain(params.ChainId) {
		if params.WatchUtxoTicker == 0 || params.WatchUtxoTicker > 300 {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "WatchUtxoTicker %d out of range", params.WatchUtxoTicker)
		}
	}
	if common.IsEVMChain(params.ChainId) {
		if !validCoreContractAddress(params.ZetaTokenContractAddress) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ZetaTokenContractAddress %s", params.ZetaTokenContractAddress)
		}
		if !validCoreContractAddress(params.ConnectorContractAddress) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ConnectorContractAddress %s", params.ConnectorContractAddress)
		}
		if !validCoreContractAddress(params.Erc20CustodyContractAddress) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid Erc20CustodyContractAddress %s", params.Erc20CustodyContractAddress)
		}
	}
	return nil
}

func validCoreContractAddress(address string) bool {
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	return ethcommon.IsHexAddress(address)
}
