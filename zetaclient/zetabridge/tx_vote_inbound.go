package zetabridge

import (
	"cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/common"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

const (
	// PostVoteInboundGasLimit is the gas limit for voting on observed inbound tx
	PostVoteInboundGasLimit = 400_000

	// PostVoteInboundExecutionGasLimit is the gas limit for voting on observed inbound tx and executing it
	PostVoteInboundExecutionGasLimit = 4_000_000

	// PostVoteInboundMessagePassingExecutionGasLimit is the gas limit for voting on, and executing ,observed inbound tx related to message passing (coin_type == zeta)
	PostVoteInboundMessagePassingExecutionGasLimit = 1_000_000
)

// GetInBoundVoteMessage returns a new MsgVoteOnObservedInboundTx
func GetInBoundVoteMessage(
	sender string,
	senderChain int64,
	txOrigin string,
	receiver string,
	receiverChain int64,
	amount math.Uint,
	message string,
	inTxHash string,
	inBlockHeight uint64,
	gasLimit uint64,
	coinType common.CoinType,
	asset string,
	signerAddress string,
	eventIndex uint,
) *types.MsgVoteOnObservedInboundTx {
	msg := types.NewMsgVoteOnObservedInboundTx(
		signerAddress,
		sender,
		senderChain,
		txOrigin,
		receiver,
		receiverChain,
		amount,
		message,
		inTxHash,
		inBlockHeight,
		gasLimit,
		coinType,
		asset,
		eventIndex,
	)
	return msg
}
