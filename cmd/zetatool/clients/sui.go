package clients

import (
	"context"
	"fmt"
	"strconv"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
)

const (
	// nanoSUIPerSUI is the number of MIST in 1 SUI (10^9)
	nanoSUIPerSUI = 1_000_000_000
)

// SuiClientAdapter wraps the Sui SDK client to implement SuiClient interface
type SuiClientAdapter struct {
	client sui.ISuiAPI
}

// NewSuiClientAdapter creates a new SuiClientAdapter
func NewSuiClientAdapter(rpcURL string) (*SuiClientAdapter, error) {
	client := sui.NewSuiClient(rpcURL)
	return &SuiClientAdapter{client: client}, nil
}

// GetBalance fetches the SUI balance for a given address
// Returns the balance in MIST (1 SUI = 10^9 MIST)
func (s *SuiClientAdapter) GetBalance(ctx context.Context, address string) (uint64, error) {
	var (
		totalBalance uint64
		cursor       = any(nil)
	)

	for {
		resp, err := s.client.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
			Owner:    address,
			CoinType: "0x2::sui::SUI",
			Cursor:   cursor,
		})
		if err != nil {
			return 0, err
		}

		for _, coin := range resp.Data {
			balance, err := strconv.ParseUint(coin.Balance, 10, 64)
			if err != nil {
				return 0, err
			}
			totalBalance += balance
		}

		if !resp.HasNextPage {
			break
		}
		cursor = resp.NextCursor
	}

	return totalBalance, nil
}

// GetSuiBalance fetches the SUI balance for a given address (standalone function)
// Returns the balance in MIST (1 SUI = 10^9 MIST)
func GetSuiBalance(ctx context.Context, rpcURL string, address string) (uint64, error) {
	adapter, err := NewSuiClientAdapter(rpcURL)
	if err != nil {
		return 0, err
	}
	return adapter.GetBalance(ctx, address)
}

// FormatSuiBalance converts MIST to SUI with 9 decimal places
func FormatSuiBalance(mist uint64) string {
	suiVal := float64(mist) / nanoSUIPerSUI
	return fmt.Sprintf("%.9f", suiVal)
}
