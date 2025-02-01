package ballot

import (
	"encoding/hex"
	"fmt"

	cosmosmath "cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/zeta-chain/node/cmd/zetatool/context"

	"github.com/zeta-chain/node/cmd/zetatool/cctx"

	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	solanarpc "github.com/zeta-chain/node/zetaclient/chains/solana/rpc"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

func solanaInboundBallotIdentifier(ctx *context.Context) (cctx.CCTXDetails, error) {
	var (
		inboundHash    = ctx.GetInboundHash()
		cctxDetails    = cctx.NewCCTXDetails()
		inboundChain   = ctx.GetInboundChain()
		zetacoreClient = ctx.GetZetaCoreClient()
		zetaChainID    = ctx.GetConfig().ZetaChainID
		cfg            = ctx.GetConfig()
		logger         = ctx.GetLogger()
		goCtx          = ctx.GetContext()
	)
	solClient := solrpc.New(cfg.SolanaRPC)
	if solClient == nil {
		return cctxDetails, fmt.Errorf("error creating rpc client")
	}

	signature := solana.MustSignatureFromBase58(inboundHash)

	txResult, err := solanarpc.GetTransaction(goCtx, solClient, signature)
	if err != nil {
		return cctxDetails, fmt.Errorf("error getting transaction: %w", err)
	}

	chainParams, err := zetacoreClient.GetChainParamsForChainID(goCtx, inboundChain.ChainId)
	if err != nil {
		return cctxDetails, fmt.Errorf("failed to get chain params: %w", err)
	}

	gatewayID, _, err := solanacontracts.ParseGatewayWithPDA(chainParams.GatewayAddress)
	if err != nil {
		return cctxDetails, fmt.Errorf("cannot parse gateway address: %s, err: %w", chainParams.GatewayAddress, err)
	}

	events, err := observer.FilterInboundEvents(txResult,
		gatewayID,
		inboundChain.ChainId,
		logger,
	)

	if err != nil {
		return cctxDetails, fmt.Errorf("failed to filter solana inbound events: %w", err)
	}

	msg := &crosschaintypes.MsgVoteInbound{}

	// build inbound vote message from events and post to zetacore
	for _, event := range events {
		msg, err = voteMsgFromSolEvent(event, zetaChainID)
		if err != nil {
			return cctxDetails, fmt.Errorf("failed to create vote message: %w", err)
		}
	}

	cctxDetails.CCCTXIdentifier = msg.Digest()
	cctxDetails.Status = cctx.PendingInboundVoting

	return cctxDetails, nil
}

// voteMsgFromSolEvent builds a MsgVoteInbound from an inbound event
func voteMsgFromSolEvent(event *clienttypes.InboundEvent,
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
