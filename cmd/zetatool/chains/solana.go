package chains

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	cosmosmath "cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"

	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
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

const (
	// lamportsPerSOL is the number of lamports in 1 SOL
	lamportsPerSOL = 1_000_000_000
)

// solanaRPCRequest represents a JSON-RPC request to Solana
type solanaRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// solanaRPCResponse represents a JSON-RPC response from Solana
type solanaRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Context struct {
			Slot uint64 `json:"slot"`
		} `json:"context"`
		Value uint64 `json:"value"`
	} `json:"result"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// GetSolanaGatewayBalance fetches the SOL balance of the gateway PDA
// The gateway PDA holds all deposited SOL funds
func GetSolanaGatewayBalance(ctx context.Context, rpcURL string, gatewayAddress string) (uint64, error) {
	// Parse gateway address and derive PDA
	_, pda, err := contracts.ParseGatewayWithPDA(gatewayAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to parse gateway address: %w", err)
	}

	// Query balance using JSON-RPC
	balance, err := getSolanaBalance(ctx, rpcURL, pda)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

// getSolanaBalance queries the balance of a Solana account using JSON-RPC
func getSolanaBalance(ctx context.Context, rpcURL string, pubkey solana.PublicKey) (uint64, error) {
	reqBody := solanaRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getBalance",
		Params:  []interface{}{pubkey.String()},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", rpcURL, bytes.NewReader(reqBytes))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("RPC returned status %d", resp.StatusCode)
	}

	var rpcResp solanaRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if rpcResp.Error != nil {
		return 0, fmt.Errorf("RPC error: %s", rpcResp.Error.Message)
	}

	return rpcResp.Result.Value, nil
}

// FormatSolanaBalance converts lamports to SOL with 9 decimal places
func FormatSolanaBalance(lamports uint64) string {
	sol := float64(lamports) / lamportsPerSOL
	return fmt.Sprintf("%.9f", sol)
}
