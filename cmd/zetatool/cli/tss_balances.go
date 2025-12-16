package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/cmd/zetatool/balances"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/rpc"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// NewTSSBalancesCMD creates a new command to check TSS address balances across all chains
func NewTSSBalancesCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "tss-balances <network>",
		Short: "Check TSS address balances across all chains",
		Long: `Check the balance of TSS (Threshold Signature Scheme) addresses across all supported chains.

The network argument must be one of: mainnet, testnet, localnet, devnet

Example:
  zetatool tss-balances mainnet
  zetatool tss-balances testnet
  zetatool tss-balances localnet --config custom_config.json`,
		Args: cobra.ExactArgs(1),
		RunE: getTSSBalances,
	}
}

func getTSSBalances(cmd *cobra.Command, args []string) error {
	network := args[0]

	// Get custom config file if provided
	configFile, err := cmd.Flags().GetString(config.FlagConfig)
	if err != nil {
		return fmt.Errorf("failed to read value for flag %s: %w", config.FlagConfig, err)
	}

	// Get config based on network
	cfg, err := config.GetConfigByNetwork(network, configFile)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	// Validate that we have a ZetaChain RPC endpoint
	if cfg.ZetaChainRPC == "" {
		return fmt.Errorf("ZetaChainRPC is not configured for network %s", network)
	}

	// Create ZetaCore client and fetch TSS
	zetacoreClient, err := rpc.NewCometBFTClients(cfg.ZetaChainRPC)
	if err != nil {
		return fmt.Errorf("failed to create zetacore client: %w", err)
	}

	ctx := context.Background()

	// Fetch TSS history to get all TSS entries
	tssHistoryRes, err := zetacoreClient.Observer.TssHistory(ctx, &observertypes.QueryTssHistoryRequest{})
	if err != nil {
		return fmt.Errorf("failed to fetch TSS history: %w", err)
	}

	if len(tssHistoryRes.TssList) == 0 {
		return fmt.Errorf("no TSS entries found")
	}

	// Iterate over all TSS entries and print balances for each
	for i, tss := range tssHistoryRes.TssList {
		if i > 0 {
			fmt.Println() // Add spacing between TSS entries
		}
		fmt.Printf("=== TSS %d of %d ===\n", i+1, len(tssHistoryRes.TssList))

		if err := balances.PrintTSSBalances(ctx, cfg, tss, network, zetacoreClient.Observer); err != nil {
			fmt.Printf("Error fetching balances for TSS (height %d): %v\n", tss.FinalizedZetaHeight, err)
			// Continue with next TSS instead of failing completely
			continue
		}
	}

	return nil
}
