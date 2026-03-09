package cli

import (
	"context"
	"fmt"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	zetatoolcommon "github.com/zeta-chain/node/cmd/zetatool/common"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/rpc"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// observerInfo holds observer address and resolved validator moniker
type observerInfo struct {
	ObserverAddress string
	OperatorAddress string
	Moniker         string
	Error           string
}

// NewListObserversCMD creates a command to list all observers with their validator monikers
func NewListObserversCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "list-observers <chain>",
		Short: "List observers with their validator monikers",
		Long: `List all active observers and resolve their validator monikers from the staking module.

The chain argument can be:
  - A chain ID (e.g., 7000, 7001)
  - A chain name (e.g., zeta_mainnet, zeta_testnet)

The network type (mainnet/testnet/etc) is inferred from the chain.

Examples:
  zetatool list-observers 7000
  zetatool list-observers zeta_mainnet
  zetatool list-observers zeta_testnet --config custom_config.json`,
		Args: cobra.ExactArgs(1),
		RunE: listObservers,
	}
}

func listObservers(cmd *cobra.Command, args []string) error {
	chain, err := zetatoolcommon.ResolveChain(args[0])
	if err != nil {
		return fmt.Errorf("failed to resolve chain %q: %w", args[0], err)
	}

	network := zetatoolcommon.NetworkTypeFromChain(chain)

	configFile, err := cmd.Flags().GetString(config.FlagConfig)
	if err != nil {
		return fmt.Errorf("failed to read value for flag %s: %w", config.FlagConfig, err)
	}

	cfg, err := config.GetConfigByNetwork(network, configFile)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	if cfg.ZetaChainRPC == "" {
		return fmt.Errorf("ZetaChainRPC is not configured for network %s", network)
	}

	zetacoreClient, err := rpc.NewCometBFTClients(cfg.ZetaChainRPC)
	if err != nil {
		return fmt.Errorf("failed to create zetacore client: %w", err)
	}

	ctx := context.Background()

	// Fetch the observer set
	observerSetRes, err := zetacoreClient.Observer.ObserverSet(ctx, &observertypes.QueryObserverSet{})
	if err != nil {
		return fmt.Errorf("failed to fetch observer set: %w", err)
	}

	if len(observerSetRes.Observers) == 0 {
		fmt.Println("No observers found")
		return nil
	}

	// Resolve monikers concurrently
	var wg sync.WaitGroup
	results := make(chan observerInfo, len(observerSetRes.Observers))

	for _, obs := range observerSetRes.Observers {
		wg.Add(1)
		go func(observerAddr string) {
			defer wg.Done()
			results <- resolveObserverMoniker(ctx, zetacoreClient.Staking, observerAddr)
		}(obs)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	infos := make([]observerInfo, 0, len(observerSetRes.Observers))
	for info := range results {
		infos = append(infos, info)
	}

	printObserverTable(infos)
	return nil
}

// resolveObserverMoniker converts an observer address to a validator operator address
// and queries the staking module for the validator's moniker.
func resolveObserverMoniker(
	ctx context.Context,
	stakingClient stakingtypes.QueryClient,
	observerAddr string,
) observerInfo {
	info := observerInfo{ObserverAddress: observerAddr}

	// Convert observer address (zeta1xxx) to validator operator address (zetavaloper1xxx)
	accAddr, err := sdk.AccAddressFromBech32(observerAddr)
	if err != nil {
		info.Error = fmt.Sprintf("invalid address: %v", err)
		return info
	}
	valAddr := sdk.ValAddress(accAddr.Bytes())
	info.OperatorAddress = valAddr.String()

	// Query the staking module for this validator
	valRes, err := stakingClient.Validator(ctx, &stakingtypes.QueryValidatorRequest{
		ValidatorAddr: info.OperatorAddress,
	})
	if err != nil {
		info.Error = fmt.Sprintf("validator query failed: %v", err)
		return info
	}

	info.Moniker = valRes.Validator.Description.Moniker
	return info
}

func printObserverTable(infos []observerInfo) {
	t := newTableWriter()
	t.AppendHeader(table.Row{"#", "Observer Address", "Moniker"})

	for i, info := range infos {
		moniker := info.Moniker
		if info.Error != "" {
			moniker = info.Error
		}

		t.AppendRow(table.Row{i + 1, info.ObserverAddress, moniker})
	}

	fmt.Println()
	t.Render()
}
