package chains

import (
	"encoding/hex"
	"fmt"
	"math/big"

	cosmosmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/cmd/zetatool/context"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	zetaclientObserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
)

func BitcoinBallotIdentifier(
	ctx *context.Context,
	btcClient *client.Client,
	params *chaincfg.Params,
	tss string,
	feeRateMultiplier float64,
	txHash string,
	senderChainID int64,
	zetacoreChainID int64,
	confirmationCount uint64,
) (cctxIdentifier string, isConfirmed bool, err error) {
	var (
		goCtx = ctx.GetContext()
	)

	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return
	}
	tx, err := btcClient.GetRawTransactionVerbose(goCtx, hash)
	if err != nil {
		return
	}

	if tx.Confirmations >= confirmationCount {
		isConfirmed = true
	}

	blockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return
	}

	blockVb, err := btcClient.GetBlockVerbose(goCtx, blockHash)
	if err != nil {
		return
	}

	event, err := zetaclientObserver.GetBtcEventWithWitness(
		goCtx,
		btcClient,
		*tx,
		tss,
		uint64(blockVb.Height), // #nosec G115 always positive
		feeRateMultiplier,
		zerolog.New(zerolog.Nop()),
		params,
		common.CalcDepositorFee,
	)
	if err != nil {
		return
	}
	if event == nil {
		err = fmt.Errorf("no event built for btc sent to TSS")
		return
	}

	cctxIdentifier, err = identifierFromBtcEvent(event, senderChainID, zetacoreChainID)
	return
}

func identifierFromBtcEvent(event *zetaclientObserver.BTCInboundEvent,
	senderChainID int64,
	zetacoreChainID int64) (cctxIdentifier string, err error) {
	// decode event memo bytes
	err = event.DecodeMemoBytes(senderChainID)
	if err != nil {
		return
	}

	// convert the amount to integer (satoshis)
	amountSats, err := common.GetSatoshis(event.Value)
	if err != nil {
		return
	}
	amountInt := big.NewInt(amountSats)

	var msg *crosschaintypes.MsgVoteInbound
	switch event.MemoStd {
	case nil:
		{
			msg = voteFromLegacyMemo(event, amountInt, senderChainID, zetacoreChainID)
		}
	default:
		{
			msg = voteFromStdMemo(event, amountInt, senderChainID, zetacoreChainID)
		}
	}
	if msg == nil {
		return
	}

	cctxIdentifier = msg.Digest()
	return
}

// NewInboundVoteFromLegacyMemo creates a MsgVoteInbound message for inbound that uses legacy memo
func voteFromLegacyMemo(
	event *zetaclientObserver.BTCInboundEvent,
	amountSats *big.Int,
	senderChainID int64,
	zetacoreChainID int64,
) *crosschaintypes.MsgVoteInbound {
	message := hex.EncodeToString(event.MemoBytes)

	return crosschaintypes.NewMsgVoteInbound(
		"",
		event.FromAddress,
		senderChainID,
		event.FromAddress,
		event.ToAddress,
		zetacoreChainID,
		cosmosmath.NewUintFromBigInt(amountSats),
		message,
		event.TxHash,
		event.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		0,
		crosschaintypes.ProtocolContractVersion_V2,
		false, // not relevant for v1
		crosschaintypes.InboundStatus_SUCCESS,
		crosschaintypes.ConfirmationMode_SAFE,
		crosschaintypes.WithCrossChainCall(len(event.MemoBytes) > 0),
	)
}

func voteFromStdMemo(
	event *zetaclientObserver.BTCInboundEvent,
	amountSats *big.Int,
	senderChainID int64,
	zetacoreChainID int64,
) *crosschaintypes.MsgVoteInbound {
	// zetacore will create a revert outbound that points to the custom revert address.
	revertOptions := crosschaintypes.RevertOptions{
		RevertAddress: event.MemoStd.RevertOptions.RevertAddress,
	}

	// check if the memo is a cross-chain call, or simple token deposit
	isCrosschainCall := event.MemoStd.OpCode == memo.OpCodeCall || event.MemoStd.OpCode == memo.OpCodeDepositAndCall

	return crosschaintypes.NewMsgVoteInbound(
		"",
		event.FromAddress,
		senderChainID,
		event.FromAddress,
		event.ToAddress,
		zetacoreChainID,
		cosmosmath.NewUintFromBigInt(amountSats),
		hex.EncodeToString(event.MemoStd.Payload),
		event.TxHash,
		event.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		0,
		crosschaintypes.ProtocolContractVersion_V2,
		false, // not relevant for v1
		event.Status,
		crosschaintypes.ConfirmationMode_SAFE,
		crosschaintypes.WithRevertOptions(revertOptions),
		crosschaintypes.WithCrossChainCall(isCrosschainCall),
	)
}
