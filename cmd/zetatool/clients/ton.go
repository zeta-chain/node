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

// GetTONGatewayBalance fetches the TON balance of the gateway contract
// Returns the balance in nanoTON (1 TON = 10^9 nanoTON)
func GetTONGatewayBalance(ctx context.Context, rpcURL string, gatewayAddress string) (uint64, error) {
	accountID, err := ton.ParseAccountID(gatewayAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to parse gateway address: %w", err)
	}

	// Create TON RPC client (chainID 0 is fine for balance queries)
	client := tonrpc.New(rpcURL, 0)

	account, err := client.GetAccountState(ctx, accountID)
	if err != nil {
		return 0, fmt.Errorf("failed to get account state: %w", err)
	}

	return account.Balance, nil
}

// FormatTONBalance converts nanoTON to TON with 9 decimal places
func FormatTONBalance(nanoTON uint64) string {
	tonVal := float64(nanoTON) / nanoTONPerTON
	return fmt.Sprintf("%.9f", tonVal)
}
