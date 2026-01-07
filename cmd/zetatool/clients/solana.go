package clients

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"

	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	solrepo "github.com/zeta-chain/node/zetaclient/chains/solana/repo"
)

// SolanaClientAdapter wraps solrepo.SolanaRepo to implement SolanaClient interface
type SolanaClientAdapter struct {
	repo *solrepo.SolanaRepo
}

// NewSolanaClientAdapter creates a new SolanaClientAdapter
func NewSolanaClientAdapter(rpcURL string) (*SolanaClientAdapter, error) {
	client := solrpc.New(rpcURL)
	if client == nil {
		return nil, fmt.Errorf("failed to create solana RPC client")
	}

	repo := solrepo.New(client)

	return &SolanaClientAdapter{
		repo: repo,
	}, nil
}

// GetTransaction retrieves a transaction by its signature
func (s *SolanaClientAdapter) GetTransaction(ctx context.Context, signature solana.Signature) (*solrpc.GetTransactionResult, error) {
	return s.repo.GetTransaction(ctx, signature)
}

// GetRawRepo returns the underlying SolanaRepo for advanced operations
func (s *SolanaClientAdapter) GetRawRepo() *solrepo.SolanaRepo {
	return s.repo
}

// GetSolanaGatewayBalance fetches the SOL balance of the gateway PDA
func GetSolanaGatewayBalance(ctx context.Context, rpcURL string, gatewayAddress string) (uint64, error) {
	_, pda, err := contracts.ParseGatewayWithPDA(gatewayAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to parse gateway address: %w", err)
	}

	client := solrpc.New(rpcURL)

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
