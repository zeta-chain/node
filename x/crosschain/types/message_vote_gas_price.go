package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/zeta-chain/node/pkg/authz"
)

var _ sdk.Msg = &MsgVoteGasPrice{}

// NewMsgVoteGasPrice creates a new MsgVoteGasPrice instance.
// Notes:
//   - Price is in `wei` for EVM and `satoshi per byte` for UTXO.
//   - PriorityFeePrice is only for EVM that supports EIP-1559 ("London" hard fork).
func NewMsgVoteGasPrice(creator string, chain int64, price, priorityFee, blockNumber uint64) *MsgVoteGasPrice {
	return &MsgVoteGasPrice{
		Creator:     creator,
		ChainId:     chain,
		Price:       price,
		PriorityFee: priorityFee,
		BlockNumber: blockNumber,
	}
}

func (msg *MsgVoteGasPrice) Route() string {
	return RouterKey
}

func (msg *MsgVoteGasPrice) Type() string {
	return authz.GasPriceVoter.String()
}

func (msg *MsgVoteGasPrice) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVoteGasPrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVoteGasPrice) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.ChainId < 0 {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidChainID, "chain id (%d)", msg.ChainId)
	}
	return nil
}
