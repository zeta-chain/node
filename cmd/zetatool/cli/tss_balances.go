package cli

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/cmd/zetatool/chains"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	pkgchains "github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/rpc"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// chainBalance represents the balance information for a single chain
type chainBalance struct {
	Chain   string
	Address string
	Balance string
	Symbol  string
	Error   string
	VM      pkgchains.Vm
}

// networkSymbols maps chain network to its native token symbol
var networkSymbols = map[pkgchains.Network]string{
	pkgchains.Network_eth:        "ETH",
	pkgchains.Network_bsc:        "BNB",
	pkgchains.Network_polygon:    "POL",
	pkgchains.Network_base:       "ETH",
	pkgchains.Network_arbitrum:   "ETH",
	pkgchains.Network_optimism:   "ETH",
	pkgchains.Network_avalanche:  "AVAX",
	pkgchains.Network_worldchain: "ETH",
	pkgchains.Network_btc:        "BTC",
	pkgchains.Network_solana:     "SOL",
	pkgchains.Network_ton:        "TON",
	pkgchains.Network_sui:        "SUI",
}

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

	// Use custom rpc config if provided
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

	// Create ZetaCore client and fetch TSS
	zetacoreClient, err := rpc.NewCometBFTClients(cfg.ZetaChainRPC)
	if err != nil {
		return fmt.Errorf("failed to create zetacore client: %w", err)
	}

	ctx := context.Background()

	tssHistoryRes, err := zetacoreClient.Observer.TssHistory(ctx, &observertypes.QueryTssHistoryRequest{})
	if err != nil {
		return fmt.Errorf("failed to fetch TSS history: %w", err)
	}

	if len(tssHistoryRes.TssList) == 0 {
		return fmt.Errorf("no TSS entries found")
	}

	for i, tss := range tssHistoryRes.TssList {
		if i > 0 {
			fmt.Println() // Add spacing between TSS entries
		}
		fmt.Printf("=== TSS %d of %d ===\n", i+1, len(tssHistoryRes.TssList))

		if err := printTSSBalances(ctx, cfg, tss, network, zetacoreClient.Observer); err != nil {
			fmt.Printf("Error fetching balances for TSS (height %d): %v\n", tss.FinalizedZetaHeight, err)
			continue
		}
	}

	return nil
}

// getSymbolForChain returns the native token symbol for a chain
func getSymbolForChain(chain pkgchains.Chain) string {
	if symbol, ok := networkSymbols[chain.Network]; ok {
		return symbol
	}
	return ""
}

// getRPCForChain returns the RPC URL for a given chain from config
func getRPCForChain(cfg *config.Config, chain pkgchains.Chain) string {
	switch chain.Network {
	case pkgchains.Network_eth:
		return cfg.EthereumRPC
	case pkgchains.Network_bsc:
		return cfg.BscRPC
	case pkgchains.Network_polygon:
		return cfg.PolygonRPC
	case pkgchains.Network_base:
		return cfg.BaseRPC
	case pkgchains.Network_arbitrum:
		return cfg.ArbitrumRPC
	case pkgchains.Network_optimism:
		return cfg.OptimismRPC
	case pkgchains.Network_avalanche:
		return cfg.AvalancheRPC
	case pkgchains.Network_worldchain:
		return cfg.WorldRPC
	case pkgchains.Network_solana:
		return cfg.SolanaRPC
	case pkgchains.Network_ton:
		return cfg.TonRPC
	case pkgchains.Network_sui:
		return cfg.SuiRPC
	default:
		return ""
	}
}

// printTSSBalances fetches and prints TSS address balances across all chains
func printTSSBalances(
	ctx context.Context,
	cfg *config.Config,
	tss observertypes.TSS,
	network string,
	observerClient observertypes.QueryClient,
) error {
	// Print TSS info
	fmt.Println("TSS Information:")
	fmt.Printf("  PubKey: %s\n", tss.TssPubkey)
	fmt.Printf("  Finalized Height: %d\n", tss.FinalizedZetaHeight)
	fmt.Println()

	// Query supported chains from zetacore
	supportedChainsRes, err := observerClient.SupportedChains(ctx, &observertypes.QuerySupportedChains{})
	if err != nil {
		return fmt.Errorf("failed to get supported chains: %w", err)
	}

	btcChainID := chains.GetBTCChainID(network)
	req := &observertypes.QueryGetTssAddressByFinalizedHeightRequest{
		FinalizedZetaHeight: tss.FinalizedZetaHeight,
		BitcoinChainId:      btcChainID,
	}
	tssAddrRes, err := observerClient.GetTssAddressByFinalizedHeight(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get TSS addresses: %w", err)
	}

	evmAddr := common.HexToAddress(tssAddrRes.Eth)
	btcAddr := tssAddrRes.Btc
	suiAddr := tssAddrRes.Sui

	// Filter external chains and group by VM type
	var evmChains, btcChains, solanaChains, tonChains, suiChains []pkgchains.Chain
	for _, chain := range supportedChainsRes.Chains {
		if !chain.IsExternal {
			continue
		}
		switch chain.Vm {
		case pkgchains.Vm_evm:
			evmChains = append(evmChains, chain)
		case pkgchains.Vm_no_vm:
			btcChains = append(btcChains, chain)
		case pkgchains.Vm_svm:
			solanaChains = append(solanaChains, chain)
		case pkgchains.Vm_tvm:
			tonChains = append(tonChains, chain)
		case pkgchains.Vm_mvm_sui:
			suiChains = append(suiChains, chain)
		}
	}

	// Fetch balances in parallel
	var wg sync.WaitGroup
	results := make(chan chainBalance, len(supportedChainsRes.Chains))

	// EVM chains - use TSS EVM address
	for _, chain := range evmChains {
		rpc := getRPCForChain(cfg, chain)
		if rpc == "" {
			results <- chainBalance{
				Chain:   chain.Name,
				Address: evmAddr.Hex(),
				VM:      chain.Vm,
				Error:   "RPC not configured",
			}
			continue
		}
		wg.Add(1)
		go func(c pkgchains.Chain, rpcURL string) {
			defer wg.Done()
			balance, err := chains.GetEVMBalance(ctx, rpcURL, evmAddr)
			if err != nil {
				results <- chainBalance{
					Chain:   c.Name,
					Address: evmAddr.Hex(),
					VM:      c.Vm,
					Error:   err.Error(),
				}
				return
			}
			results <- chainBalance{
				Chain:   c.Name,
				Address: evmAddr.Hex(),
				Balance: chains.FormatEVMBalance(balance),
				Symbol:  getSymbolForChain(c),
				VM:      c.Vm,
			}
		}(chain, rpc)
	}

	// Bitcoin chains - use TSS BTC address
	for _, chain := range btcChains {
		wg.Add(1)
		go func(c pkgchains.Chain) {
			defer wg.Done()
			// Skip Bitcoin for localnet (mempool.space doesn't support regtest)
			if network == config.NetworkLocalnet {
				results <- chainBalance{
					Chain:   c.Name,
					Address: btcAddr,
					VM:      c.Vm,
					Error:   "Localnet not supported (uses regtest)",
				}
				return
			}
			balance, err := chains.GetBTCBalance(ctx, btcAddr, network)
			if err != nil {
				results <- chainBalance{
					Chain:   c.Name,
					Address: btcAddr,
					VM:      c.Vm,
					Error:   err.Error(),
				}
				return
			}
			results <- chainBalance{
				Chain:   c.Name,
				Address: btcAddr,
				Balance: fmt.Sprintf("%.8f", balance),
				Symbol:  getSymbolForChain(c),
				VM:      c.Vm,
			}
		}(chain)
	}

	// Sui chains - use TSS Sui address
	for _, chain := range suiChains {
		rpc := getRPCForChain(cfg, chain)
		if rpc == "" {
			results <- chainBalance{
				Chain:   chain.Name,
				Address: suiAddr,
				VM:      chain.Vm,
				Error:   "RPC not configured",
			}
			continue
		}
		wg.Add(1)
		go func(c pkgchains.Chain, rpcURL string) {
			defer wg.Done()
			balance, err := chains.GetSuiBalance(ctx, rpcURL, suiAddr)
			if err != nil {
				results <- chainBalance{
					Chain:   c.Name,
					Address: suiAddr,
					VM:      c.Vm,
					Error:   err.Error(),
				}
				return
			}
			results <- chainBalance{
				Chain:   c.Name,
				Address: suiAddr,
				Balance: chains.FormatSuiBalance(balance),
				Symbol:  getSymbolForChain(c),
				VM:      c.Vm,
			}
		}(chain, rpc)
	}

	// Solana chains - use gateway PDA balance
	for _, chain := range solanaChains {
		rpc := getRPCForChain(cfg, chain)
		if rpc == "" {
			results <- chainBalance{
				Chain:   chain.Name,
				Address: "N/A",
				VM:      chain.Vm,
				Error:   "RPC not configured",
			}
			continue
		}
		wg.Add(1)
		go func(c pkgchains.Chain, rpcURL string) {
			defer wg.Done()
			chainParamsReq := &observertypes.QueryGetChainParamsForChainRequest{
				ChainId: c.ChainId,
			}
			chainParamsRes, err := observerClient.GetChainParamsForChain(ctx, chainParamsReq)
			if err != nil {
				results <- chainBalance{
					Chain:   c.Name,
					Address: "N/A",
					VM:      c.Vm,
					Error:   fmt.Sprintf("failed to get chain params: %v", err),
				}
				return
			}

			gatewayAddress := chainParamsRes.ChainParams.GatewayAddress
			if gatewayAddress == "" {
				results <- chainBalance{
					Chain:   c.Name,
					Address: "N/A",
					VM:      c.Vm,
					Error:   "Gateway address not configured",
				}
				return
			}

			balance, err := chains.GetSolanaGatewayBalance(ctx, rpcURL, gatewayAddress)
			if err != nil {
				results <- chainBalance{
					Chain:   c.Name,
					Address: gatewayAddress,
					VM:      c.Vm,
					Error:   err.Error(),
				}
				return
			}

			results <- chainBalance{
				Chain:   c.Name,
				Address: gatewayAddress,
				Balance: chains.FormatSolanaBalance(balance),
				Symbol:  getSymbolForChain(c),
				VM:      c.Vm,
			}
		}(chain, rpc)
	}

	// TON chains - use gateway contract balance
	for _, chain := range tonChains {
		rpc := getRPCForChain(cfg, chain)
		if rpc == "" {
			results <- chainBalance{
				Chain:   chain.Name,
				Address: "N/A",
				VM:      chain.Vm,
				Error:   "RPC not configured",
			}
			continue
		}
		wg.Add(1)
		go func(c pkgchains.Chain, rpcURL string) {
			defer wg.Done()
			chainParamsReq := &observertypes.QueryGetChainParamsForChainRequest{
				ChainId: c.ChainId,
			}
			chainParamsRes, err := observerClient.GetChainParamsForChain(ctx, chainParamsReq)
			if err != nil {
				results <- chainBalance{
					Chain:   c.Name,
					Address: "N/A",
					VM:      c.Vm,
					Error:   fmt.Sprintf("failed to get chain params: %v", err),
				}
				return
			}

			gatewayAddress := chainParamsRes.ChainParams.GatewayAddress
			if gatewayAddress == "" {
				results <- chainBalance{
					Chain:   c.Name,
					Address: "N/A",
					VM:      c.Vm,
					Error:   "Gateway address not configured",
				}
				return
			}

			balance, err := chains.GetTONGatewayBalance(ctx, rpcURL, gatewayAddress)
			if err != nil {
				results <- chainBalance{
					Chain:   c.Name,
					Address: gatewayAddress,
					VM:      c.Vm,
					Error:   err.Error(),
				}
				return
			}

			results <- chainBalance{
				Chain:   c.Name,
				Address: gatewayAddress,
				Balance: chains.FormatTONBalance(balance),
				Symbol:  getSymbolForChain(c),
				VM:      c.Vm,
			}
		}(chain, rpc)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	balances := make([]chainBalance, 0, len(supportedChainsRes.Chains))
	for result := range results {
		balances = append(balances, result)
	}

	printBalanceTable(balances)

	return nil
}

// printBalanceTable prints the balance results in a formatted table
func printBalanceTable(balances []chainBalance) {
	vmGroups := make(map[pkgchains.Vm][]chainBalance)
	for _, b := range balances {
		vmGroups[b.VM] = append(vmGroups[b.VM], b)
	}

	vmOrder := []pkgchains.Vm{
		pkgchains.Vm_evm,
		pkgchains.Vm_no_vm,   // Bitcoin
		pkgchains.Vm_svm,     // Solana
		pkgchains.Vm_tvm,     // TON
		pkgchains.Vm_mvm_sui, // Sui
	}

	vmLabels := map[pkgchains.Vm]string{
		pkgchains.Vm_evm:     "evm",
		pkgchains.Vm_no_vm:   "btc",
		pkgchains.Vm_svm:     "svm",
		pkgchains.Vm_tvm:     "tvm",
		pkgchains.Vm_mvm_sui: "sui",
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"VM", "Chain", "Address", "Balance"})

	for _, vm := range vmOrder {
		groupBalances, ok := vmGroups[vm]
		if !ok || len(groupBalances) == 0 {
			continue
		}

		for _, b := range groupBalances {
			addr := b.Address
			if len(addr) > 44 {
				addr = addr[:20] + "..." + addr[len(addr)-20:]
			}

			var balanceStr string
			if b.Error != "" {
				balanceStr = b.Error
			} else if b.Symbol != "" {
				balanceStr = fmt.Sprintf("%s %s", b.Balance, b.Symbol)
			} else {
				balanceStr = b.Balance
			}

			t.AppendRow(table.Row{vmLabels[vm], b.Chain, addr, balanceStr})
		}
	}

	fmt.Println()
	t.Render()
}
