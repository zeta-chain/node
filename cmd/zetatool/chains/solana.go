package chains

import (
	"context"
	"encoding/hex"
	"fmt"

	cosmosmath "cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"

	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

// VoteMsgFromSolEvent builds a MsgVoteInbound from an inbound event
func VoteMsgFromSolEvent(event *clienttypes.InboundEvent,
	zetaChainID int64) (*crosschaintypes.MsgVoteInbound, error) {
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

// GetSolanaGatewayBalance fetches the SOL balance of the gateway PDA
func GetSolanaGatewayBalance(ctx context.Context, rpcURL string, gatewayAddress string) (uint64, error) {
	_, pda, err := contracts.ParseGatewayWithPDA(gatewayAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to parse gateway address: %w", err)
	}

	client := solrpc.New(rpcURL)

	result, err := client.GetBalance(ctx, pda, solrpc.CommitmentFinalized)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return result.Value, nil
}

// FormatSolanaBalance converts lamports to SOL with 9 decimal places
func FormatSolanaBalance(lamports uint64) string {
	sol := float64(lamports) / float64(solana.LAMPORTS_PER_SOL)
	return fmt.Sprintf("%.9f", sol)
}
