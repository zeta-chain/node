package zrepo

import (
	"context"

	cosmosmath "cosmossdk.io/math"
	upgrade "cosmossdk.io/x/upgrade/types"
	cometbft "github.com/cometbft/cometbft/types"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/go-tss/blame"

	"github.com/zeta-chain/node/pkg/chains"
	crosschain "github.com/zeta-chain/node/x/crosschain/types"
	fungible "github.com/zeta-chain/node/x/fungible/types"
	observer "github.com/zeta-chain/node/x/observer/types"
	keys "github.com/zeta-chain/node/zetaclient/keys/interfaces"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

var _ ZetacoreClient = &zetacore.Client{}

type ChainID = int64
type Nonce = uint64

// ZetacoreWriter contains the functions that mutate ZetaChain state.
type ZetacoreWriter interface {
	PostVoteGasPrice(_ context.Context,
		_ chains.Chain,
		gasPrice uint64,
		priorityFee uint64,
		blockNum uint64,
	) (string, error)

	PostVoteTSS(_ context.Context,
		tssPubKey string,
		keyGenZetaHeight int64,
		_ chains.ReceiveStatus,
	) (string, error)

	PostVoteBlameData(_ context.Context,
		_ *blame.Blame,
		_ ChainID,
		index string,
	) (string, error)

	PostVoteOutbound(_ context.Context,
		gasLimit uint64,
		retryGasLimit uint64,
		_ *crosschain.MsgVoteOutbound,
	) (string, string, error)

	PostVoteInbound(_ context.Context,
		gasLimit uint64,
		retryGasLimit uint64,
		_ *crosschain.MsgVoteInbound,
	) (string, string, error)

	PostOutboundTracker(_ context.Context,
		_ ChainID,
		_ Nonce,
		txHash string,
	) (string, error)
}

// ZetacoreClientRepo contains the functions used by ZetaRepo.
type ZetacoreClientRepo interface {
	ZetacoreWriter

	Chain() chains.Chain

	GetKeys() keys.ObserverKeys

	ListPendingCCTX(context.Context, chains.Chain) ([]*crosschain.CrossChainTx, uint64, error)

	GetForeignCoinsFromAsset(context.Context, ChainID, eth.Address) (fungible.ForeignCoins, error)

	GetPendingNoncesByChain(context.Context, ChainID) (observer.PendingNonces, error)

	GetCctxByNonce(context.Context, ChainID, Nonce) (*crosschain.CrossChainTx, error)

	GetCctxByHash(context.Context, string) (*crosschain.CrossChainTx, error)

	GetBallotByID(context.Context, string) (*observer.QueryBallotByIdentifierResponse, error)

	GetOutboundTracker(context.Context, ChainID, Nonce) (*crosschain.OutboundTracker, error)

	GetOutboundTrackers(context.Context, ChainID) ([]crosschain.OutboundTracker, error)

	GetInboundTrackersForChain(context.Context, ChainID) ([]crosschain.InboundTracker, error)

	NewBlockSubscriber(context.Context) (chan cometbft.EventDataNewBlock, error)

	GetBTCTSSAddress(context.Context, ChainID) (string, error)
}

// ZetacoreClient is the client interface that interacts with zetacore.
//
// TODO: this should be moved elsewhere, since it is not used by ZetaRepo.
// See: https://github.com/zeta-chain/node/issues/4300
//
//go:generate mockery --name ZetacoreClient --filename zetacore_client.go --case underscore --output ../../testutils/mocks
type ZetacoreClient interface {
	ZetacoreClientRepo

	GetSupportedChains(context.Context) ([]chains.Chain, error)

	GetAdditionalChains(context.Context) ([]chains.Chain, error)

	GetChainParams(context.Context) ([]*observer.ChainParams, error)

	GetKeyGen(context.Context) (observer.Keygen, error)

	GetTSS(context.Context) (observer.TSS, error)
	GetTSSHistory(context.Context) ([]observer.TSS, error)

	GetBlockHeight(context.Context) (int64, error)

	ListPendingCCTXWithinRateLimit(context.Context,
	) (*crosschain.QueryListPendingCctxWithinRateLimitResponse, error)

	GetRateLimiterInput(_ context.Context,
		window int64,
	) (*crosschain.QueryRateLimiterInputResponse, error)

	GetCrosschainFlags(context.Context) (observer.CrosschainFlags, error)
	GetRateLimiterFlags(context.Context) (crosschain.RateLimiterFlags, error)
	GetOperationalFlags(context.Context) (observer.OperationalFlags, error)

	GetObserverList(context.Context) ([]string, error)

	GetZetaHotKeyBalance(context.Context) (cosmosmath.Int, error)

	GetUpgradePlan(context.Context) (*upgrade.Plan, error)
}
