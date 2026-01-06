package clients

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"

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
