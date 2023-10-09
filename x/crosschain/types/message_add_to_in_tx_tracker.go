package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
)

const TypeMsgAddToInTxTracker = "AddToInTxTracker"

var _ sdk.Msg = &MsgAddToInTxTracker{}

func NewMsgAddToInTxTracker(creator string, chain int64, coinType common.CoinType, txHash string) *MsgAddToInTxTracker {
	return &MsgAddToInTxTracker{
		Creator:  creator,
		ChainId:  chain,
		TxHash:   txHash,
		CoinType: coinType,
	}
}

func (msg *MsgAddToInTxTracker) Route() string {
	return RouterKey
}

func (msg *MsgAddToInTxTracker) Type() string {
	return TypeMsgAddToInTxTracker
}

func (msg *MsgAddToInTxTracker) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddToInTxTracker) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddToInTxTracker) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	chain := common.GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return errorsmod.Wrapf(ErrInvalidChainID, "chain id (%d)", msg.ChainId)
	}
	if msg.Proof != nil && !chain.IsProvable() {
		return errorsmod.Wrapf(ErrCannotVerifyProof, "chain id %d does not support proof-based trackers", msg.ChainId)
	}
	_, err = common.GetCoinType(msg.CoinType.String())
	if err != nil {
		return err
	}

	return nil
}
