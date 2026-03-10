package cli

import (
	"context"
	"fmt"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/cmd/zetatool/clients"
	zetatoolcommon "github.com/zeta-chain/node/cmd/zetatool/common"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	pkgchains "github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/rpc"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// migrationStatus holds the tracked status for a single TSS migration CCTX
type migrationStatus struct {
	ChainID       int64
	ChainName     string
	CctxIndex     string
	CctxStatus    string
	OutboundHash  string
	ReceiptStatus string
}

// NewTrackTSSMigrationCMD creates a command to track TSS fund migration status
func NewTrackTSSMigrationCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "track-tss-migration <chain>",
		Short: "Track TSS fund migration CCTX status and outbound receipts",
		Long: `Track the status of TSS fund migration CCTXs across all chains.

For each migration CCTX, fetches the CCTX status, outbound hash, and queries the
destination chain for the transaction receipt status.

The chain argument can be:
  - A chain ID (e.g., 7000, 7001)
  - A chain name (e.g., zeta_mainnet, zeta_testnet)

The network type (mainnet/testnet/etc) is inferred from the chain.

Examples:
  zetatool track-tss-migration 7000
  zetatool track-tss-migration zeta_mainnet
  zetatool track-tss-migration zeta_testnet --config custom_config.json`,
		Args: cobra.ExactArgs(1),
		RunE: trackTSSMigration,
	}
}

func trackTSSMigration(cmd *cobra.Command, args []string) error {
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

	// Fetch all TSS fund migrators
	migratorsRes, err := zetacoreClient.Observer.TssFundsMigratorInfoAll(
		ctx,
		&observertypes.QueryTssFundsMigratorInfoAllRequest{},
	)
	if err != nil {
		return fmt.Errorf("failed to fetch TSS fund migrators: %w", err)
	}

	if len(migratorsRes.TssFundsMigrators) == 0 {
		fmt.Println("No TSS fund migration entries found")
		return nil
	}

	fmt.Printf("Found %d TSS fund migration(s)\n", len(migratorsRes.TssFundsMigrators))

	// Process each migrator concurrently
	var wg sync.WaitGroup
	results := make(chan migrationStatus, len(migratorsRes.TssFundsMigrators))

	for _, migrator := range migratorsRes.TssFundsMigrators {
		wg.Add(1)
		go func(m observertypes.TssFundMigratorInfo) {
			defer wg.Done()
			results <- processMigration(ctx, &zetacoreClient, cfg, m)
		}(migrator)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	statuses := make([]migrationStatus, 0, len(migratorsRes.TssFundsMigrators))
	for status := range results {
		statuses = append(statuses, status)
	}

	printMigrationTable(statuses)
	return nil
}

// processMigration fetches the CCTX and receipt for a single migration entry
func processMigration(
	ctx context.Context,
	zetacoreClient *rpc.Clients,
	cfg *config.Config,
	migrator observertypes.TssFundMigratorInfo,
) migrationStatus {
	status := migrationStatus{
		ChainID:   migrator.ChainId,
		CctxIndex: migrator.MigrationCctxIndex,
	}

	// Resolve chain name
	chain, found := pkgchains.GetChainFromChainID(migrator.ChainId, nil)
	if found {
		status.ChainName = chain.Name
	} else {
		status.ChainName = fmt.Sprintf("unknown(%d)", migrator.ChainId)
	}

	// Fetch the CCTX
	cctx, err := zetacoreClient.GetCctxByHash(ctx, migrator.MigrationCctxIndex)
	if err != nil {
		status.CctxStatus = fmt.Sprintf("error: %v", err)
		status.ReceiptStatus = "N/A"
		return status
	}

	if cctx.CctxStatus != nil {
		status.CctxStatus = cctx.CctxStatus.Status.String()
	} else {
		status.CctxStatus = "unknown"
	}

	// Get the first (and only) outbound
	if len(cctx.OutboundParams) == 0 {
		status.OutboundHash = "N/A"
		status.ReceiptStatus = "no outbound"
		return status
	}

	outbound := cctx.OutboundParams[0]
	status.OutboundHash = outbound.Hash

	if outbound.Hash == "" {
		status.ReceiptStatus = "no hash yet"
		return status
	}

	// Fetch receipt from the outbound chain
	status.ReceiptStatus = fetchOutboundReceipt(ctx, cfg, outbound)
	return status
}

// fetchOutboundReceipt fetches the transaction receipt from the outbound chain
func fetchOutboundReceipt(
	ctx context.Context,
	cfg *config.Config,
	outbound *crosschaintypes.OutboundParams,
) string {
	outboundChain, found := pkgchains.GetChainFromChainID(outbound.ReceiverChainId, nil)
	if !found {
		return fmt.Sprintf("unknown chain %d", outbound.ReceiverChainId)
	}

	switch outboundChain.Vm {
	case pkgchains.Vm_evm:
		return fetchEVMReceipt(ctx, cfg, outboundChain, outbound.Hash)
	case pkgchains.Vm_no_vm:
		return fetchBTCTxStatus(ctx, outbound.Hash, outboundChain.ChainId)
	default:
		return fmt.Sprintf("receipt check not supported for %s", outboundChain.Vm.String())
	}
}

// fetchEVMReceipt fetches a transaction receipt from an EVM chain and returns its status
func fetchEVMReceipt(
	ctx context.Context,
	cfg *config.Config,
	chain pkgchains.Chain,
	txHash string,
) string {
	rpcURL := clients.ResolveEVMRPC(chain, cfg)
	if rpcURL == "" {
		return "RPC not configured"
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Sprintf("dial error: %v", err)
	}
	defer client.Close()

	receipt, err := client.TransactionReceipt(ctx, ethcommon.HexToHash(txHash))
	if err != nil {
		return fmt.Sprintf("receipt error: %v", err)
	}

	return formatEVMReceiptStatus(receipt)
}

func formatEVMReceiptStatus(receipt *ethtypes.Receipt) string {
	switch receipt.Status {
	case ethtypes.ReceiptStatusSuccessful:
		return "success"
	case ethtypes.ReceiptStatusFailed:
		return "failed"
	default:
		return fmt.Sprintf("unknown(%d)", receipt.Status)
	}
}

// fetchBTCTxStatus checks if a BTC transaction exists using mempool.space API
func fetchBTCTxStatus(ctx context.Context, txHash string, chainID int64) string {
	confirmed, err := clients.IsBTCTxConfirmed(ctx, txHash, chainID)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	if confirmed {
		return "confirmed"
	}
	return "unconfirmed"
}

func printMigrationTable(statuses []migrationStatus) {
	t := newTableWriter()
	t.AppendHeader(table.Row{"Chain ID", "Chain", "Migration CCTX Index", "CCTX Status", "Outbound Hash", "Receipt Status"})

	for _, s := range statuses {
		outHash := s.OutboundHash
		if len(outHash) > 20 {
			outHash = outHash[:10] + "..." + outHash[len(outHash)-10:]
		}

		t.AppendRow(table.Row{s.ChainID, s.ChainName, s.CctxIndex, s.CctxStatus, outHash, s.ReceiptStatus})
	}

	fmt.Println()
	t.Render()
}
