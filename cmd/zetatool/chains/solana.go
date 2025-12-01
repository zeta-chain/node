package chains

import (
	"encoding/hex"

	cosmosmath "cosmossdk.io/math"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

// VoteMsgFromSolEvent builds a MsgVoteInbound from an inbound event
func VoteMsgFromSolEvent(event *clienttypes.InboundEvent,
	zetaChainID int64) (*crosschaintypes.MsgVoteInbound, error) {
	// create inbound vote message
	return crosschaintypes.NewMsgVoteInbound(
		"",
		event.Sender,
		event.SenderChainID,
		event.Sender,
		event.Receiver,
		zetaChainID,
		cosmosmath.NewUint(event.Amount),
		hex.EncodeToString(event.Memo),
		event.TxHash,
		event.BlockNumber,
		0,
		event.CoinType,
		event.Asset,
		uint64(event.Index),
		crosschaintypes.ProtocolContractVersion_V2,
		false,
		crosschaintypes.InboundStatus_SUCCESS,
		crosschaintypes.ConfirmationMode_SAFE,
		crosschaintypes.WithCrossChainCall(event.IsCrossChainCall),
	), nil
}
