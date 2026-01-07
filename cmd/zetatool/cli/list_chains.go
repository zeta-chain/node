package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	zetatoolcommon "github.com/zeta-chain/node/cmd/zetatool/common"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/rpc"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func NewListChainsCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "list-chains [chain]",
		Short: "List available chains",
		Long: `List available chains with their names and chain IDs.

Without arguments, lists all chains known to zetatool.
With a chain argument, resolves the network type and queries zetacore for supported chains.

The chain argument can be:
  - A chain ID (e.g., 7001)
  - A chain name (e.g., zeta_mainnet)

Examples:
  zetatool list-chains                     # List all known chains
  zetatool list-chains 7000                # List supported chains for mainnet
  zetatool list-chains zeta_testnet        # List supported chains for testnet`,
		Args: cobra.MaximumNArgs(1),
		RunE: listChains,
	}
}

func listChains(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		allChains := chains.CombineDefaultChainsList([]chains.Chain{})
		printChainList(allChains)
		return nil
	}

	chainArg := args[0]
	chain, err := zetatoolcommon.ResolveChain(chainArg)
	if err != nil {
		return fmt.Errorf("failed to resolve chain %q: %w", chainArg, err)
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
	supportedChainsRes, err := zetacoreClient.Observer.SupportedChains(ctx, &observertypes.QuerySupportedChains{})
	if err != nil {
		return fmt.Errorf("failed to get supported chains: %w", err)
	}

	printChainList(supportedChainsRes.Chains)
	return nil
}

func printChainList(chainList []chains.Chain) {
	fmt.Println()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Chain ID", "Name", "Network", "VM", "External"})

	for _, c := range chainList {
		external := "no"
		if c.IsExternal {
			external = "yes"
		}
		t.AppendRow(table.Row{c.ChainId, c.Name, c.Network.String(), c.Vm.String(), external})
	}

	t.Render()
}
