package chains

import (
	"encoding/hex"
	"fmt"

	cosmosmath "cosmossdk.io/math"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

// voteMsgFromSolEvent builds a MsgVoteInbound from an inbound event
func VoteMsgFromSolEvent(event *clienttypes.InboundEvent,
	zetaChainID int64) (*crosschaintypes.MsgVoteInbound, error) {
	// decode event memo bytes to get the receiver
	err := event.DecodeMemo()
	if err != nil {
		return nil, fmt.Errorf("failed to decode memo: %w", err)
	}

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
		0, // not a smart contract call
		crosschaintypes.ProtocolContractVersion_V1,
		false, // not relevant for v1
		crosschaintypes.InboundStatus_SUCCESS,
	), nil
}
