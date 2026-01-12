package clients

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/rs/zerolog"

	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	solobserver "github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	solrepo "github.com/zeta-chain/node/zetaclient/chains/solana/repo"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

// SolanaClientAdapter wraps solrepo.SolanaRepo to implement SolanaClient interface
type SolanaClientAdapter struct {
	client *solrpc.Client
	repo   *solrepo.SolanaRepo
}

// NewSolanaClientAdapter creates a new SolanaClientAdapter
func NewSolanaClientAdapter(rpcURL string) (*SolanaClientAdapter, error) {
	client := solrpc.New(rpcURL)
	if client == nil {
		return nil, fmt.Errorf("failed to create solana RPC client")
	}

	repo := solrepo.New(client)

	return &SolanaClientAdapter{
		client: client,
		repo:   repo,
	}, nil
}

// GetTransaction retrieves a transaction by its signature
func (s *SolanaClientAdapter) GetTransaction(
	ctx context.Context,
	signature solana.Signature,
) (*solrpc.GetTransactionResult, error) {
	return s.repo.GetTransaction(ctx, signature)
}

// ProcessTransactionResultWithAddressLookups resolves address lookups in a transaction
func (s *SolanaClientAdapter) ProcessTransactionResultWithAddressLookups(
	ctx context.Context,
	txResult *solrpc.GetTransactionResult,
	logger zerolog.Logger,
	signature solana.Signature,
) *solana.Transaction {
	return solobserver.ProcessTransactionResultWithAddressLookups(ctx, txResult, s.client, logger, signature)
}

// FilterInboundEvents filters inbound events from a transaction result
func (s *SolanaClientAdapter) FilterInboundEvents(
	txResult *solrpc.GetTransactionResult,
	gatewayID solana.PublicKey,
	chainID int64,
	logger zerolog.Logger,
	tx *solana.Transaction,
) ([]*clienttypes.InboundEvent, error) {
	return solobserver.FilterInboundEvents(txResult, gatewayID, chainID, logger, tx)
}

// GetSolanaGatewayBalance fetches the SOL balance of the gateway PDA
func GetSolanaGatewayBalance(ctx context.Context, rpcURL string, gatewayAddress string) (uint64, error) {
	_, pda, err := contracts.ParseGatewayWithPDA(gatewayAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to parse gateway address: %w", err)
	}

	client := solrpc.New(rpcURL)
	if client == nil {
		return 0, fmt.Errorf("failed to create solana rpc client")
	}

	result, err := client.GetBalance(ctx, pda, solrpc.CommitmentFinalized)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return result.Value, nil
}

// FormatSolanaBalance converts lamports to SOL with 9 decimal places
func FormatSolanaBalance(lamports uint64) string {
	sol := float64(lamports) / float64(solana.LAMPORTS_PER_SOL)
	return fmt.Sprintf("%.9f", sol)
}
