package balances

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/cmd/zetatool/chains"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// ChainBalance represents the balance information for a single chain
type ChainBalance struct {
	Chain   string
	Address string
	Balance string
	Symbol  string
	Error   string
}

// PrintTSSBalances fetches and prints TSS address balances across all chains
func PrintTSSBalances(
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

	// Get Bitcoin chain ID based on network
	btcChainID := chains.GetBTCChainID(network)

	// Query TSS addresses using FinalizedZetaHeight
	req := &observertypes.QueryGetTssAddressByFinalizedHeightRequest{
		FinalizedZetaHeight: tss.FinalizedZetaHeight,
		BitcoinChainId:      btcChainID,
	}
	tssAddrRes, err := observerClient.GetTssAddressByFinalizedHeight(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get TSS addresses: %w", err)
	}

	// Use addresses from the query response
	evmAddr := common.HexToAddress(tssAddrRes.Eth)
	btcAddr := tssAddrRes.Btc
	suiAddr := tssAddrRes.Sui

	// Fetch balances in parallel
	var wg sync.WaitGroup
	results := make(chan ChainBalance, 15)

	// EVM chains
	evmChains := []struct {
		name   string
		rpc    string
		symbol string
	}{
		{"Ethereum", cfg.EthereumRPC, "ETH"},
		{"BSC", cfg.BscRPC, "BNB"},
		{"Polygon", cfg.PolygonRPC, "POL"},
		{"Base", cfg.BaseRPC, "ETH"},
		{"Arbitrum", cfg.ArbitrumRPC, "ETH"},
		{"Optimism", cfg.OptimismRPC, "ETH"},
		{"Avalanche", cfg.AvalancheRPC, "AVAX"},
	}

	for _, chain := range evmChains {
		if chain.rpc == "" {
			results <- ChainBalance{
				Chain:   chain.name,
				Address: evmAddr.Hex(),
				Error:   "RPC not configured",
			}
			continue
		}
		wg.Add(1)
		go func(name, rpc, symbol string) {
			defer wg.Done()
			balance, err := chains.GetEVMBalance(ctx, rpc, evmAddr)
			if err != nil {
				results <- ChainBalance{
					Chain:   name,
					Address: evmAddr.Hex(),
					Error:   err.Error(),
				}
				return
			}
			results <- ChainBalance{
				Chain:   name,
				Address: evmAddr.Hex(),
				Balance: chains.FormatEVMBalance(balance),
				Symbol:  symbol,
			}
		}(chain.name, chain.rpc, chain.symbol)
	}

	// Bitcoin - uses mempool.space public API
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Skip Bitcoin for localnet (mempool.space doesn't support regtest)
		if network == chains.NetworkLocalnet {
			results <- ChainBalance{
				Chain:   "Bitcoin",
				Address: btcAddr,
				Error:   "Localnet not supported (uses regtest)",
			}
			return
		}
		balance, err := chains.GetBTCBalance(ctx, btcAddr, network)
		if err != nil {
			results <- ChainBalance{
				Chain:   "Bitcoin",
				Address: btcAddr,
				Error:   err.Error(),
			}
			return
		}
		results <- ChainBalance{
			Chain:   "Bitcoin",
			Address: btcAddr,
			Balance: fmt.Sprintf("%.8f", balance),
			Symbol:  "BTC",
		}
	}()

	// Sui
	wg.Add(1)
	go func() {
		defer wg.Done()
		if cfg.SuiRPC == "" {
			results <- ChainBalance{
				Chain:   "Sui",
				Address: suiAddr,
				Error:   "RPC not configured",
			}
			return
		}
		balance, err := chains.GetSuiBalance(ctx, cfg.SuiRPC, suiAddr)
		if err != nil {
			results <- ChainBalance{
				Chain:   "Sui",
				Address: suiAddr,
				Error:   err.Error(),
			}
			return
		}
		results <- ChainBalance{
			Chain:   "Sui",
			Address: suiAddr,
			Balance: chains.FormatSuiBalance(balance),
			Symbol:  "SUI",
		}
	}()

	// Solana - query gateway PDA balance
	wg.Add(1)
	go func() {
		defer wg.Done()
		if cfg.SolanaRPC == "" {
			results <- ChainBalance{
				Chain:   "Solana",
				Address: "N/A",
				Error:   "RPC not configured",
			}
			return
		}

		// Get Solana chain params to find gateway address
		solanaChainID := chains.GetSolanaChainID(network)
		chainParamsReq := &observertypes.QueryGetChainParamsForChainRequest{
			ChainId: solanaChainID,
		}
		chainParamsRes, err := observerClient.GetChainParamsForChain(ctx, chainParamsReq)
		if err != nil {
			results <- ChainBalance{
				Chain:   "Solana",
				Address: "N/A",
				Error:   fmt.Sprintf("failed to get chain params: %v", err),
			}
			return
		}

		gatewayAddress := chainParamsRes.ChainParams.GatewayAddress
		if gatewayAddress == "" {
			results <- ChainBalance{
				Chain:   "Solana",
				Address: "N/A",
				Error:   "Gateway address not configured",
			}
			return
		}

		// Query gateway PDA balance
		balance, err := chains.GetSolanaGatewayBalance(ctx, cfg.SolanaRPC, gatewayAddress)
		if err != nil {
			results <- ChainBalance{
				Chain:   "Solana",
				Address: gatewayAddress,
				Error:   err.Error(),
			}
			return
		}

		results <- ChainBalance{
			Chain:   "Solana",
			Address: gatewayAddress,
			Balance: chains.FormatSolanaBalance(balance),
			Symbol:  "SOL",
		}
	}()

	// TON - query gateway contract balance
	wg.Add(1)
	go func() {
		defer wg.Done()
		if cfg.TonRPC == "" {
			results <- ChainBalance{
				Chain:   "TON",
				Address: "N/A",
				Error:   "RPC not configured",
			}
			return
		}

		// Get TON chain params to find gateway address
		tonChainID := chains.GetTONChainID(network)
		chainParamsReq := &observertypes.QueryGetChainParamsForChainRequest{
			ChainId: tonChainID,
		}
		chainParamsRes, err := observerClient.GetChainParamsForChain(ctx, chainParamsReq)
		if err != nil {
			results <- ChainBalance{
				Chain:   "TON",
				Address: "N/A",
				Error:   fmt.Sprintf("failed to get chain params: %v", err),
			}
			return
		}

		gatewayAddress := chainParamsRes.ChainParams.GatewayAddress
		if gatewayAddress == "" {
			results <- ChainBalance{
				Chain:   "TON",
				Address: "N/A",
				Error:   "Gateway address not configured",
			}
			return
		}

		// Query gateway balance
		balance, err := chains.GetTONGatewayBalance(ctx, cfg.TonRPC, gatewayAddress)
		if err != nil {
			results <- ChainBalance{
				Chain:   "TON",
				Address: gatewayAddress,
				Error:   err.Error(),
			}
			return
		}

		results <- ChainBalance{
			Chain:   "TON",
			Address: gatewayAddress,
			Balance: chains.FormatTONBalance(balance),
			Symbol:  "TON",
		}
	}()

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	balances := make([]ChainBalance, 0, 15)
	for result := range results {
		balances = append(balances, result)
	}

	// Print results in a formatted table
	printBalanceTable(balances)

	return nil
}

// printBalanceTable prints the balance results in a formatted table
func printBalanceTable(balances []ChainBalance) {
	// Define chain order for consistent output
	chainOrder := []string{
		"Ethereum", "BSC", "Polygon", "Base",
		"Arbitrum", "Optimism", "Avalanche",
		"Bitcoin", "Sui", "Solana", "TON",
	}

	// Create a map for quick lookup
	balanceMap := make(map[string]ChainBalance)
	for _, b := range balances {
		balanceMap[b.Chain] = b
	}

	// Print header
	fmt.Println("Chain Balances:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-12s %-44s %s\n", "Chain", "Address", "Balance")
	fmt.Println(strings.Repeat("-", 80))

	// Print each chain in order
	for _, chain := range chainOrder {
		b, ok := balanceMap[chain]
		if !ok {
			continue
		}

		addr := b.Address
		if len(addr) > 42 {
			addr = addr[:20] + "..." + addr[len(addr)-18:]
		}

		if b.Error != "" {
			fmt.Printf("%-12s %-44s %s\n", b.Chain, addr, b.Error)
		} else {
			fmt.Printf("%-12s %-44s %s %s\n", b.Chain, addr, b.Balance, b.Symbol)
		}
	}

	fmt.Println(strings.Repeat("-", 80))
}
