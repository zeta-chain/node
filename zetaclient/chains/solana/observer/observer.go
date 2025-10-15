package observer

import (
	"context"

	sol "github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/solana/repo"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// SolanaClient is the interface for the Solana RPC client.
type SolanaClient interface {
	GetVersion(context.Context) (*solrpc.GetVersionResult, error)

	GetHealth(context.Context) (string, error)

	GetSlot(context.Context, solrpc.CommitmentType) (uint64, error)

	GetBlockTime(_ context.Context, block uint64) (*sol.UnixTimeSeconds, error)

	GetAccountInfo(context.Context, sol.PublicKey) (*solrpc.GetAccountInfoResult, error)

	GetAccountInfoWithOpts(
		context.Context,
		sol.PublicKey,
		*solrpc.GetAccountInfoOpts,
	) (*solrpc.GetAccountInfoResult, error)

	GetBalance(_ context.Context,
		account sol.PublicKey,
		_ solrpc.CommitmentType,
	) (*solrpc.GetBalanceResult, error)

	GetLatestBlockhash(context.Context,
		solrpc.CommitmentType,
	) (*solrpc.GetLatestBlockhashResult, error)

	GetRecentPrioritizationFees(_ context.Context,
		accounts sol.PublicKeySlice,
	) ([]solrpc.PriorizationFeeResult, error)

	GetTransaction(context.Context,
		sol.Signature,
		*solrpc.GetTransactionOpts,
	) (*solrpc.GetTransactionResult, error)

	GetConfirmedTransactionWithOpts(context.Context,
		sol.Signature,
		*solrpc.GetTransactionOpts,
	) (*solrpc.TransactionWithMeta, error)

	GetSignaturesForAddressWithOpts(context.Context,
		sol.PublicKey,
		*solrpc.GetSignaturesForAddressOpts,
	) ([]*solrpc.TransactionSignature, error)

	SendTransactionWithOpts(context.Context,
		*sol.Transaction,
		solrpc.TransactionOpts,
	) (sol.Signature, error)
}

// Observer is the observer for the Solana chain
type Observer struct {
	// base.Observer implements the base chain observer
	*base.Observer

	// solanaClient is the Solana RPC client that interacts with the Solana chain
	solanaClient SolanaClient

	solanaRepo *repo.SolanaRepo

	// gatewayID is the program ID of gateway program on Solana chain
	gatewayID sol.PublicKey

	// pda is the program derived address of the gateway program
	pda sol.PublicKey

	// finalizedTxResults indexes tx results with the outbound hash
	finalizedTxResults map[string]*solrpc.GetTransactionResult
}

// New Observer constructor
func New(baseObserver *base.Observer,
	solanaClient SolanaClient,
	gatewayAddress string,
) (*Observer, error) {
	// parse gateway ID and PDA
	gatewayID, pda, err := contracts.ParseGatewayWithPDA(gatewayAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse gateway address %s", gatewayAddress)
	}

	solanaRepo := repo.New(solanaClient)

	// create solana observer
	ob := &Observer{
		Observer:           baseObserver,
		solanaClient:       solanaClient,
		solanaRepo:         solanaRepo,
		gatewayID:          gatewayID,
		pda:                pda,
		finalizedTxResults: make(map[string]*solrpc.GetTransactionResult),
	}

	ob.Observer.LoadLastTxScanned()

	return ob, nil
}

// LoadLastTxScanned loads the last scanned tx from the database.
func (ob *Observer) LoadLastTxScanned() error {
	ob.Observer.LoadLastTxScanned()
	ob.Logger().Chain.Info().
		Str(logs.FieldTx, ob.LastTxScanned()).
		Msg("chain starts scanning from tx")

	return nil
}

// SetTxResult sets the tx result for the given nonce
func (ob *Observer) SetTxResult(nonce uint64, result *solrpc.GetTransactionResult) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	ob.finalizedTxResults[ob.OutboundID(nonce)] = result
}

// GetTxResult returns the tx result for the given nonce
func (ob *Observer) GetTxResult(nonce uint64) *solrpc.GetTransactionResult {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.finalizedTxResults[ob.OutboundID(nonce)]
}

// IsTxFinalized returns true if there is a finalized tx for nonce
func (ob *Observer) IsTxFinalized(nonce uint64) bool {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.finalizedTxResults[ob.OutboundID(nonce)] != nil
}

// CheckRPCStatus checks the RPC status of the Solana chain
func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	blockTime, err := ob.solanaRepo.HealthCheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to check rpc status")
	}

	metrics.ReportBlockLatency(ob.Chain().Name, blockTime)

	return nil
}
