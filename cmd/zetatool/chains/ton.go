package chains

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tonkeeper/tongo/ton"
)

const (
	// nanoTONPerTON is the number of nanoTON in 1 TON (10^9)
	nanoTONPerTON = 1_000_000_000
)

// tonRPCRequest represents a JSON-RPC request to TON
type tonRPCRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      string         `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
}

// tonAddressInfoResponse represents a response from getAddressInformation
type tonAddressInfoResponse struct {
	OK     bool           `json:"ok"`
	Result tonAccountInfo `json:"result"`
	Error  string         `json:"error"`
}

type tonAccountInfo struct {
	Balance string `json:"balance"`
	State   string `json:"state"`
}

// GetTONGatewayBalance fetches the TON balance of the gateway contract
// Returns the balance in nanoTON (1 TON = 10^9 nanoTON)
func GetTONGatewayBalance(ctx context.Context, rpcURL string, gatewayAddress string) (uint64, error) {
	// Parse and validate the gateway address
	accountID, err := ton.ParseAccountID(gatewayAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to parse gateway address: %w", err)
	}

	// Query balance using JSON-RPC
	balance, err := getTONBalance(ctx, rpcURL, accountID)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

// getTONBalance queries the balance of a TON account using JSON-RPC
func getTONBalance(ctx context.Context, rpcURL string, accountID ton.AccountID) (uint64, error) {
	// Ensure endpoint has /jsonRPC suffix (consistent with TON RPC client)
	endpoint := strings.TrimRight(rpcURL, "/") + "/jsonRPC"

	// Create JSON-RPC request for getAddressInformation
	reqBody := tonRPCRequest{
		JSONRPC: "2.0",
		ID:      "1",
		Method:  "getAddressInformation",
		Params: map[string]any{
			"address": accountID.ToRaw(),
		},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(reqBytes))
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

	// Decode response
	var rpcResp tonAddressInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if !rpcResp.OK {
		return 0, fmt.Errorf("RPC error: %s", rpcResp.Error)
	}

	// Parse balance from string to uint64
	var balance uint64
	if _, err := fmt.Sscanf(rpcResp.Result.Balance, "%d", &balance); err != nil {
		return 0, fmt.Errorf("failed to parse balance %q: %w", rpcResp.Result.Balance, err)
	}

	return balance, nil
}

// FormatTONBalance converts nanoTON to TON with 9 decimal places
func FormatTONBalance(nanoTON uint64) string {
	tonVal := float64(nanoTON) / nanoTONPerTON
	return fmt.Sprintf("%.9f", tonVal)
}

// GetTONChainID returns the TON chain ID for the given network
func GetTONChainID(network string) int64 {
	switch network {
	case NetworkMainnet:
		return 2015140 // TON mainnet
	case NetworkTestnet:
		return 2015141 // TON testnet
	case NetworkLocalnet:
		return 2015142 // TON localnet
	default:
		return 2015141 // Default to testnet
	}
}
