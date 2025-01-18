package inbound

import (
	"context"
	"encoding/hex"
	"fmt"

	cosmosmath "cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/pkg/rpc"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	solanarpc "github.com/zeta-chain/node/zetaclient/chains/solana/rpc"

	//solanarpc "github.com/zeta-chain/node/zetaclient/chains/solana/rpc"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

func solanaInboundBallotIdentifier(ctx context.Context,
	cfg config.Config,
	zetacoreClient rpc.Clients,
	inboundHash string,
	inboundChain chains.Chain,
	zetaChainID int64,
	logger zerolog.Logger) (string, error) {
	solClient := solrpc.New(cfg.SolanaRPC)
	if solClient == nil {
		return "", fmt.Errorf("error creating rpc client")
	}

	signature := solana.MustSignatureFromBase58(inboundHash)

	txResult, err := solanarpc.GetTransaction(ctx, solClient, signature)
	if err != nil {
		return "", fmt.Errorf("error getting transaction: %w", err)
	}

	chainParams, err := zetacoreClient.GetChainParamsForChainID(context.Background(), inboundChain.ChainId)
	if err != nil {
		return "", fmt.Errorf("failed to get chain params %w", err)
	}

	gatewayID, _, err := solanacontracts.ParseGatewayWithPDA(chainParams.GatewayAddress)
	if err != nil {
		return "", fmt.Errorf("cannot parse gateway address %s, err %w", chainParams.GatewayAddress, err)
	}

	observerLogger := &base.ObserverLogger{
		Inbound: logger,
	}

	events, err := observer.FilterSolanaInboundEvents(txResult,
		observerLogger,
		gatewayID,
		inboundChain.ChainId,
	)

	msg := &crosschaintypes.MsgVoteInbound{}

	// build inbound vote message from events and post to zetacore
	for _, event := range events {
		msg, err = voteMsgFromSolEvent(event, zetaChainID)
		if err != nil {
			return "", fmt.Errorf("failed to create vote message: %w", err)
		}
	}

	return fmt.Sprintf("ballot identifier: %s", msg.Digest()), nil
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
	), nil
}
