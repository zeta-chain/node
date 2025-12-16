package chains

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

// GetSuiBalance fetches the SUI balance for a given address
// Returns the balance in MIST (1 SUI = 10^9 MIST)
func GetSuiBalance(ctx context.Context, rpcURL string, address string) (uint64, error) {
	client := sui.NewSuiClient(rpcURL)

	var (
		totalBalance uint64
		cursor       = any(nil)
	)

	// Query all SUI coins owned by the address
	// There can be multiple coins, so we paginate through them
	for {
		resp, err := client.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
			Owner:    address,
			CoinType: "0x2::sui::SUI",
			Cursor:   cursor,
		})
		if err != nil {
			return 0, err
		}

		// Sum up balances from all coins
		for _, coin := range resp.Data {
			balance, err := strconv.ParseUint(coin.Balance, 10, 64)
			if err != nil {
				return 0, err
			}
			totalBalance += balance
		}

		// Check if there are more pages
		if !resp.HasNextPage {
			break
		}
		cursor = resp.NextCursor
	}

	return totalBalance, nil
}

// FormatSuiBalance converts MIST to SUI with 9 decimal places
func FormatSuiBalance(mist uint64) string {
	suiVal := float64(mist) / nanoSUIPerSUI
	return fmt.Sprintf("%.9f", suiVal)
}
