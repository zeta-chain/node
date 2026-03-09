package cli

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/cmd/zetatool/clients"
	zetatoolcommon "github.com/zeta-chain/node/cmd/zetatool/common"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	pkgchains "github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/rpc"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	btccommon "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

const (
	// conservativeFeeRate is a conservative fee rate in sat/vB for migration calculations.
	// We use 50 sat/vB which is 5x the default testnet rate to ensure outbound goes through.
	conservativeFeeRate = 50

	// reservedRBFFees is the amount reserved for potential RBF fee bumping (0.01 BTC)
	reservedRBFFees = 0.01

	// nonceMarkBuffer is a buffer for nonce mark output in BTC (0.00003 BTC = 3000 satoshis)
	nonceMarkBuffer = 0.00003

	// satoshisPerBTC is the number of satoshis in 1 BTC
	satoshisPerBTC = 100_000_000
)

// chainBalance represents the balance information for a single chain
type chainBalance struct {
	ChainName          string
	ChainID            int64
	Address            string
	Balance            string
	MigrationAmount    string
	MigrationAmountRaw string // Raw amount in wei (EVM) or satoshis (BTC) for direct use in migration command
	Symbol             string
	PendingNonceLow    int64
	PendingNonceHigh   int64
	NonceAvailable     bool // true if pending nonces were successfully retrieved for this chain
	Error              string
	VM                 pkgchains.Vm
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

const (
	// FlagMigrationAmounts is the flag to show migration amount columns
	FlagMigrationAmounts = "migration-amounts"
	// FlagShowNonces is the flag to show pending nonce low and high columns
	FlagShowNonces = "show-nonces"
)

// NewTSSBalancesCMD creates a new command to check TSS address balances across all chains
func NewTSSBalancesCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tss-balances <chain>",
		Short: "Check TSS address balances across all chains",
		Long: `Check the balance of TSS (Threshold Signature Scheme) addresses across all supported chains.

The chain argument can be:
  - A chain ID (e.g., 7000, 1, 56)
  - A chain name (e.g., zeta_mainnet, eth_mainnet)

The network type (mainnet/testnet/etc) is inferred from the chain.

Examples:
  zetatool tss-balances 7000
  zetatool tss-balances zeta_mainnet
  zetatool tss-balances zeta_testnet --config custom_config.json
  zetatool tss-balances zeta_testnet --raw-amounts`,
		Args: cobra.ExactArgs(1),
		RunE: getTSSBalances,
	}

	cmd.Flags().Bool(FlagMigrationAmounts, false, "Show migration amount and raw migration amount columns")
	cmd.Flags().Bool(FlagShowNonces, false, "Show pending nonce low and high columns")

	return cmd
}

func getTSSBalances(cmd *cobra.Command, args []string) error {
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

	showMigrationAmounts, err := cmd.Flags().GetBool(FlagMigrationAmounts)
	if err != nil {
		return fmt.Errorf("failed to read value for flag %s: %w", FlagMigrationAmounts, err)
	}

	showNonces, err := cmd.Flags().GetBool(FlagShowNonces)
	if err != nil {
		return fmt.Errorf("failed to read value for flag %s: %w", FlagShowNonces, err)
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

		if err := printTSSBalances(
			ctx,
			cfg,
			tss,
			network,
			zetacoreClient.Observer,
			showMigrationAmounts,
			showNonces,
		); err != nil {
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

// calculateBTCMigrationAmount calculates the maximum amount that can be migrated from a Bitcoin TSS address
// after accounting for transaction fees, RBF reserve, and nonce mark.
func calculateBTCMigrationAmount(balance float64) (migrationAmt float64) {
	// Calculate estimated fee using conservative estimates:
	// - OutboundBytesMax (1543 vB) for maximum transaction size
	// - conservativeFeeRate (50 sat/vB) to account for network congestion
	estimatedFee := float64(conservativeFeeRate*btccommon.OutboundBytesMax) / satoshisPerBTC

	// Total overhead includes: estimated fee + RBF reserve + nonce mark buffer
	totalOverhead := estimatedFee + reservedRBFFees + nonceMarkBuffer
	migrationAmt = balance - totalOverhead

	// Ensure migration amount is not negative
	if migrationAmt < 0 {
		migrationAmt = 0
	}

	return migrationAmt
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

// fetchAllPendingNonces fetches all pending nonces from the observer client
func fetchAllPendingNonces(
	ctx context.Context,
	observerClient observertypes.QueryClient,
) ([]observertypes.PendingNonces, error) {
	res, err := observerClient.PendingNoncesAll(ctx, &observertypes.QueryAllPendingNoncesRequest{})
	if err != nil {
		return nil, err
	}
	return res.PendingNonces, nil
}

// getPendingNoncesForChainAndTss finds the pending nonces for a specific chain and TSS from the list.
func getPendingNoncesForChainAndTss(
	allNonces []observertypes.PendingNonces,
	chainID int64,
	tssPubkey string,
) (nonceLow, nonceHigh int64, ok bool) {
	for _, nonce := range allNonces {
		if nonce.ChainId == chainID && nonce.Tss == tssPubkey {
			return nonce.NonceLow, nonce.NonceHigh, true
		}
	}
	return 0, 0, false
}

// printTSSBalances fetches and prints TSS address balances across all chains
func printTSSBalances(
	ctx context.Context,
	cfg *config.Config,
	tss observertypes.TSS,
	network string,
	observerClient observertypes.QueryClient,
	showMigrationAmounts bool,
	showNonces bool,
) error {
	// Print TSS info
	fmt.Println("TSS Information:")
	fmt.Printf("  PubKey: %s\n", tss.TssPubkey)
	fmt.Printf("  Finalized Height: %d\n", tss.FinalizedZetaHeight)
	fmt.Println()

	// Print Bitcoin fee estimation info if showing migration amounts
	if showMigrationAmounts {
		fmt.Println("Bitcoin Fee Estimation Parameters:")
		fmt.Printf("  Conservative Fee Rate: %d sat/vB\n", conservativeFeeRate)
		fmt.Printf("  Max Transaction Size: %d vB\n", btccommon.OutboundBytesMax)
		fmt.Printf(
			"  Estimated Fee: %.8f BTC\n",
			float64(conservativeFeeRate*btccommon.OutboundBytesMax)/satoshisPerBTC,
		)
		fmt.Printf("  RBF Reserve: %.8f BTC\n", reservedRBFFees)
		fmt.Printf("  Nonce Mark Buffer: %.8f BTC\n", nonceMarkBuffer)
		fmt.Printf(
			"  Total Overhead: %.8f BTC\n",
			float64(conservativeFeeRate*btccommon.OutboundBytesMax)/satoshisPerBTC+reservedRBFFees+nonceMarkBuffer,
		)
		fmt.Println()
	}

	// Query supported chains from zetacore
	supportedChainsRes, err := observerClient.SupportedChains(ctx, &observertypes.QuerySupportedChains{})
	if err != nil {
		return fmt.Errorf("failed to get supported chains: %w", err)
	}

	btcChainID, err := clients.GetBTCChainID(network)
	if err != nil {
		return fmt.Errorf("failed to get BTC chain ID: %w", err)
	}
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

	var wg sync.WaitGroup
	results := make(chan chainBalance, len(supportedChainsRes.Chains))

	// EVM chains - use TSS EVM address (migration amount equals balance, fees handled by zetacore)
	for _, chain := range evmChains {
		chainRPC := getRPCForChain(cfg, chain)
		if chainRPC == "" {
			results <- chainBalance{
				ChainName: chain.Name,
				ChainID:   chain.ChainId,
				Address:   evmAddr.Hex(),
				VM:        chain.Vm,
				Error:     "RPC not configured",
			}
			continue
		}
		wg.Add(1)
		go func(c pkgchains.Chain, rpcURL string) {
			defer wg.Done()
			balance, err := clients.GetEVMBalance(ctx, rpcURL, evmAddr)
			if err != nil {
				results <- chainBalance{
					ChainName: c.Name,
					ChainID:   c.ChainId,
					Address:   evmAddr.Hex(),
					VM:        c.Vm,
					Error:     err.Error(),
				}
				return
			}
			formattedBalance := clients.FormatEVMBalance(balance)
			results <- chainBalance{
				ChainName:          c.Name,
				ChainID:            c.ChainId,
				Address:            evmAddr.Hex(),
				Balance:            formattedBalance,
				MigrationAmount:    formattedBalance, // Same as balance for EVM
				MigrationAmountRaw: balance.String(), // Raw wei amount for migration command
				Symbol:             getSymbolForChain(c),
				VM:                 c.Vm,
			}
		}(chain, chainRPC)
	}

	for _, chain := range btcChains {
		wg.Add(1)
		go func(c pkgchains.Chain) {
			defer wg.Done()
			// Skip Bitcoin for localnet (mempool.space doesn't support regtest)
			if network == config.NetworkLocalnet {
				results <- chainBalance{
					ChainName: c.Name,
					ChainID:   c.ChainId,
					Address:   btcAddr,
					VM:        c.Vm,
					Error:     "Localnet not supported (uses regtest)",
				}
				return
			}
			balance, err := clients.GetBTCBalance(ctx, btcAddr, c.ChainId)
			if err != nil {
				results <- chainBalance{
					ChainName: c.Name,
					ChainID:   c.ChainId,
					Address:   btcAddr,
					VM:        c.Vm,
					Error:     err.Error(),
				}
				return
			}
			// Calculate migration amount after fee deduction
			migrationAmt := calculateBTCMigrationAmount(balance)
			// Convert migration amount to satoshis for raw value
			migrationAmtSats := int64(migrationAmt * satoshisPerBTC)
			results <- chainBalance{
				ChainName:          c.Name,
				ChainID:            c.ChainId,
				Address:            btcAddr,
				Balance:            fmt.Sprintf("%.8f", balance),
				MigrationAmount:    fmt.Sprintf("%.8f", migrationAmt),
				MigrationAmountRaw: fmt.Sprintf("%d", migrationAmtSats), // Raw satoshi amount for migration command
				Symbol:             getSymbolForChain(c),
				VM:                 c.Vm,
			}
		}(chain)
	}

	for _, chain := range suiChains {
		chainRPC := getRPCForChain(cfg, chain)
		if chainRPC == "" {
			results <- chainBalance{
				ChainName:          chain.Name,
				ChainID:            chain.ChainId,
				Address:            suiAddr,
				MigrationAmount:    "N/A",
				MigrationAmountRaw: "N/A",
				Symbol:             getSymbolForChain(chain),
				VM:                 chain.Vm,
				Error:              "RPC not configured",
			}
			continue
		}
		wg.Add(1)
		go func(c pkgchains.Chain, rpcURL string) {
			defer wg.Done()
			balance, err := clients.GetSuiBalance(ctx, rpcURL, suiAddr)
			if err != nil {
				results <- chainBalance{
					ChainName:          c.Name,
					ChainID:            c.ChainId,
					Address:            suiAddr,
					MigrationAmount:    "N/A",
					MigrationAmountRaw: "N/A",
					Symbol:             getSymbolForChain(c),
					VM:                 c.Vm,
					Error:              err.Error(),
				}
				return
			}
			results <- chainBalance{
				ChainName:          c.Name,
				ChainID:            c.ChainId,
				Address:            suiAddr,
				Balance:            clients.FormatSuiBalance(balance),
				MigrationAmount:    "N/A",
				MigrationAmountRaw: "N/A",
				Symbol:             getSymbolForChain(c),
				VM:                 c.Vm,
			}
		}(chain, chainRPC)
	}

	for _, chain := range solanaChains {
		chainRPC := getRPCForChain(cfg, chain)
		if chainRPC == "" {
			results <- chainBalance{
				ChainName:          chain.Name,
				ChainID:            chain.ChainId,
				Address:            "N/A",
				MigrationAmount:    "N/A",
				MigrationAmountRaw: "N/A",
				Symbol:             getSymbolForChain(chain),
				VM:                 chain.Vm,
				Error:              "RPC not configured",
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
					ChainName:          c.Name,
					ChainID:            c.ChainId,
					Address:            "N/A",
					MigrationAmount:    "N/A",
					MigrationAmountRaw: "N/A",
					Symbol:             getSymbolForChain(c),
					VM:                 c.Vm,
					Error:              fmt.Sprintf("failed to get chain params: %v", err),
				}
				return
			}

			gatewayAddress := chainParamsRes.ChainParams.GatewayAddress
			if gatewayAddress == "" {
				results <- chainBalance{
					ChainName:          c.Name,
					ChainID:            c.ChainId,
					Address:            "N/A",
					MigrationAmount:    "N/A",
					MigrationAmountRaw: "N/A",
					Symbol:             getSymbolForChain(c),
					VM:                 c.Vm,
					Error:              "Gateway address not configured",
				}
				return
			}

			balance, err := clients.GetSolanaGatewayBalance(ctx, rpcURL, gatewayAddress)
			if err != nil {
				results <- chainBalance{
					ChainName:          c.Name,
					ChainID:            c.ChainId,
					Address:            gatewayAddress,
					MigrationAmount:    "N/A",
					MigrationAmountRaw: "N/A",
					Symbol:             getSymbolForChain(c),
					VM:                 c.Vm,
					Error:              err.Error(),
				}
				return
			}

			results <- chainBalance{
				ChainName:          c.Name,
				ChainID:            c.ChainId,
				Address:            gatewayAddress,
				Balance:            clients.FormatSolanaBalance(balance),
				MigrationAmount:    "N/A",
				MigrationAmountRaw: "N/A",
				Symbol:             getSymbolForChain(c),
				VM:                 c.Vm,
			}
		}(chain, chainRPC)
	}

	for _, chain := range tonChains {
		chainRPC := getRPCForChain(cfg, chain)
		if chainRPC == "" {
			results <- chainBalance{
				ChainName:          chain.Name,
				ChainID:            chain.ChainId,
				Address:            "N/A",
				MigrationAmount:    "N/A",
				MigrationAmountRaw: "N/A",
				Symbol:             getSymbolForChain(chain),
				VM:                 chain.Vm,
				Error:              "RPC not configured",
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
					ChainName:          c.Name,
					ChainID:            c.ChainId,
					Address:            "N/A",
					MigrationAmount:    "N/A",
					MigrationAmountRaw: "N/A",
					Symbol:             getSymbolForChain(c),
					VM:                 c.Vm,
					Error:              fmt.Sprintf("failed to get chain params: %v", err),
				}
				return
			}

			gatewayAddress := chainParamsRes.ChainParams.GatewayAddress
			if gatewayAddress == "" {
				results <- chainBalance{
					ChainName:          c.Name,
					ChainID:            c.ChainId,
					Address:            "N/A",
					MigrationAmount:    "N/A",
					MigrationAmountRaw: "N/A",
					Symbol:             getSymbolForChain(c),
					VM:                 c.Vm,
					Error:              "Gateway address not configured",
				}
				return
			}

			balance, err := clients.GetTONGatewayBalance(ctx, rpcURL, gatewayAddress)
			if err != nil {
				results <- chainBalance{
					ChainName:          c.Name,
					ChainID:            c.ChainId,
					Address:            gatewayAddress,
					MigrationAmount:    "N/A",
					MigrationAmountRaw: "N/A",
					Symbol:             getSymbolForChain(c),
					VM:                 c.Vm,
					Error:              err.Error(),
				}
				return
			}

			results <- chainBalance{
				ChainName:          c.Name,
				ChainID:            c.ChainId,
				Address:            gatewayAddress,
				Balance:            clients.FormatTONBalance(balance),
				MigrationAmount:    "N/A",
				MigrationAmountRaw: "N/A",
				Symbol:             getSymbolForChain(c),
				VM:                 c.Vm,
			}
		}(chain, chainRPC)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	balances := make([]chainBalance, 0, len(supportedChainsRes.Chains))
	for result := range results {
		balances = append(balances, result)
	}

	// Fetch pending nonces for all chains if flag is enabled
	if showNonces {
		allNonces, err := fetchAllPendingNonces(ctx, observerClient)
		if err != nil {
			log.Warn().Err(err).Msg("failed to fetch pending nonces")
		} else {
			for i := range balances {
				if balances[i].ChainID != 0 {
					nonceLow, nonceHigh, ok := getPendingNoncesForChainAndTss(
						allNonces,
						balances[i].ChainID,
						tss.TssPubkey,
					)
					if ok {
						balances[i].PendingNonceLow = nonceLow
						balances[i].PendingNonceHigh = nonceHigh
						balances[i].NonceAvailable = true
					}
					// NonceAvailable remains false if no matching entry was found
				}
			}
		}
	}

	printBalanceTable(balances, showMigrationAmounts, showNonces)

	return nil
}

// printBalanceTable prints the balance results in a formatted table
func printBalanceTable(balances []chainBalance, showMigrationAmounts bool, showNonces bool) {
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

	// Build header based on flags
	header := table.Row{"VM", "Chain", "Chain ID", "Address", "Balance"}
	if showMigrationAmounts {
		header = append(header, "Migration Amount", "Migration Amount (Raw)")
	}
	if showNonces {
		header = append(header, "Nonce Low", "Nonce High")
	}
	t.AppendHeader(header)

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

			var balanceStr, migrationStr string
			if b.Error != "" {
				balanceStr = b.Error
				migrationStr = b.Error
			} else if b.Symbol != "" {
				balanceStr = fmt.Sprintf("%s %s", b.Balance, b.Symbol)
				migrationStr = fmt.Sprintf("%s %s", b.MigrationAmount, b.Symbol)
			} else {
				balanceStr = b.Balance
				migrationStr = b.MigrationAmount
			}

			// Build row based on flags
			row := table.Row{vmLabels[vm], b.ChainName, b.ChainID, addr, balanceStr}
			if showMigrationAmounts {
				row = append(row, migrationStr, b.MigrationAmountRaw)
			}
			if showNonces {
				if b.NonceAvailable {
					row = append(row, b.PendingNonceLow, b.PendingNonceHigh)
				} else {
					row = append(row, "N/A", "N/A")
				}
			}
			t.AppendRow(row)
		}
	}

	fmt.Println()
	t.Render()
}
