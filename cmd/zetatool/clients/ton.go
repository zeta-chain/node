package clients

import (
	"context"
	"fmt"

	"github.com/tonkeeper/tongo/ton"

	tonrpc "github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
)

const (
	// nanoTONPerTON is the number of nanoTON in 1 TON (10^9)
	nanoTONPerTON = 1_000_000_000
)

// TONClientAdapter wraps the TON RPC client to implement TONClient interface
type TONClientAdapter struct {
	client *tonrpc.Client
}

// NewTONClientAdapter creates a new TONClientAdapter
func NewTONClientAdapter(rpcURL string) (*TONClientAdapter, error) {
	client := tonrpc.New(rpcURL, 0)
	return &TONClientAdapter{client: client}, nil
}

// GetAccountBalance fetches the TON balance for a given address
// Returns the balance in nanoTON (1 TON = 10^9 nanoTON)
func (t *TONClientAdapter) GetAccountBalance(ctx context.Context, address string) (uint64, error) {
	accountID, err := ton.ParseAccountID(address)
	if err != nil {
		return 0, fmt.Errorf("failed to parse address: %w", err)
	}

	account, err := t.client.GetAccountState(ctx, accountID)
	if err != nil {
		return 0, fmt.Errorf("failed to get account state: %w", err)
	}

	return account.Balance, nil
}

// GetTONGatewayBalance fetches the TON balance of the gateway contract (standalone function)
// Returns the balance in nanoTON (1 TON = 10^9 nanoTON)
func GetTONGatewayBalance(ctx context.Context, rpcURL string, gatewayAddress string) (uint64, error) {
	adapter, err := NewTONClientAdapter(rpcURL)
	if err != nil {
		return 0, err
	}
	return adapter.GetAccountBalance(ctx, gatewayAddress)
}

// FormatTONBalance converts nanoTON to TON with 9 decimal places
func FormatTONBalance(nanoTON uint64) string {
	tonVal := float64(nanoTON) / nanoTONPerTON
	return fmt.Sprintf("%.9f", tonVal)
}
