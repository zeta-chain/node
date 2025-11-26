package repo

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type SolanaClient interface {
	GetSlot(context.Context, rpc.CommitmentType) (uint64, error)

	GetLatestBlockhash(context.Context, rpc.CommitmentType) (*rpc.GetLatestBlockhashResult, error)

	GetBlockTime(_ context.Context, block uint64) (*solana.UnixTimeSeconds, error)

	// TODO switch for GetAccoungInfoWithOpts
	GetAccountInfo(context.Context, solana.PublicKey) (*rpc.GetAccountInfoResult, error)

	GetAccountInfoWithOpts(context.Context,
		solana.PublicKey,
		*rpc.GetAccountInfoOpts,
	) (*rpc.GetAccountInfoResult, error)

	GetBalance(context.Context,
		solana.PublicKey,
		rpc.CommitmentType,
	) (*rpc.GetBalanceResult, error)

	GetTransaction(context.Context,
		solana.Signature,
		*rpc.GetTransactionOpts,
	) (*rpc.GetTransactionResult, error)

	GetSignaturesForAddressWithOpts(context.Context,
		solana.PublicKey,
		*rpc.GetSignaturesForAddressOpts,
	) ([]*rpc.TransactionSignature, error)

	GetRecentPrioritizationFees(_ context.Context,
		accounts solana.PublicKeySlice,
	) ([]rpc.PriorizationFeeResult, error)

	// This is a mutating function that does not get called when zetaclient is in dry-mode.
	SendTransactionWithOpts(context.Context,
		*solana.Transaction,
		rpc.TransactionOpts,
	) (solana.Signature, error)
}
