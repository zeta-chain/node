package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
)

const TypeMsgAddInboundTracker = "AddInboundTracker"

var _ sdk.Msg = &MsgAddInboundTracker{}

func NewMsgAddInboundTracker(creator string, chain int64, coinType coin.CoinType, txHash string) *MsgAddInboundTracker {
	return &MsgAddInboundTracker{
		Creator:  creator,
		ChainId:  chain,
		TxHash:   txHash,
		CoinType: coinType,
	}
}

func (msg *MsgAddInboundTracker) Route() string {
	return RouterKey
}

func (msg *MsgAddInboundTracker) Type() string {
	return TypeMsgAddInboundTracker
}

func (msg *MsgAddInboundTracker) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddInboundTracker) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddInboundTracker) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	chain := chains.GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return errorsmod.Wrapf(ErrInvalidChainID, "chain id (%d)", msg.ChainId)
	}
	if msg.Proof != nil && !chain.SupportMerkleProof() {
		return errorsmod.Wrapf(ErrProofVerificationFail, "chain id %d does not support proof-based trackers", msg.ChainId)
	}
	_, ok := coin.CoinType_value[msg.CoinType.String()]
	if !ok {
		return errorsmod.Wrapf(ErrProofVerificationFail, "coin-type not supported")
	}
	return nil
}
