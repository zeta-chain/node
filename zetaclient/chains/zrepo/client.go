package zrepo

import (
	"context"

	cometbft "github.com/cometbft/cometbft/types"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/go-tss/blame"

	"github.com/zeta-chain/node/pkg/chains"
	zetaerrors "github.com/zeta-chain/node/pkg/errors"
	crosschain "github.com/zeta-chain/node/x/crosschain/types"
	fungible "github.com/zeta-chain/node/x/fungible/types"
	observer "github.com/zeta-chain/node/x/observer/types"
	keys "github.com/zeta-chain/node/zetaclient/keys/interfaces"
)

type ChainID = int64
type Nonce = uint64

type ZetacoreClient interface {
	zetacoreReaderClient
	zetacoreWriterClient
}

// zetacoreReaderClient contains the functions that do not mutate ZetaChain state.
type zetacoreReaderClient interface {
	Chain() chains.Chain

	NewBlockSubscriber(context.Context) (chan cometbft.EventDataNewBlock, error)

	HasVoted(_ context.Context, ballotIndex string, voterAddress string) (bool, error)

	ListPendingCCTX(context.Context, chains.Chain) ([]*crosschain.CrossChainTx, uint64, error)

	GetCctxByNonce(context.Context, ChainID, Nonce) (*crosschain.CrossChainTx, error)

	GetCctxByHash(context.Context, string) (*crosschain.CrossChainTx, error)

	GetOutboundTracker(context.Context, ChainID, Nonce) (*crosschain.OutboundTracker, error)

	GetOutboundTrackers(context.Context, ChainID) ([]crosschain.OutboundTracker, error)

	GetInboundTrackersForChain(context.Context, ChainID) ([]crosschain.InboundTracker, error)

	GetKeys() keys.ObserverKeys

	GetForeignCoinsFromAsset(context.Context, ChainID, eth.Address) (fungible.ForeignCoins, error)

	GetPendingNoncesByChain(context.Context, ChainID) (observer.PendingNonces, error)

	GetBallotByID(context.Context, string) (*observer.QueryBallotByIdentifierResponse, error)

	GetBTCTSSAddress(context.Context, ChainID) (string, error)
}

// zetacoreWriterClient contains the functions that mutate ZetaChain state.
type zetacoreWriterClient interface {
	PostVoteGasPrice(_ context.Context,
		_ chains.Chain,
		gasPrice uint64,
		priorityFee uint64,
		block uint64,
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
		monitorErrCh chan<- zetaerrors.ErrTxMonitor,
	) (string, string, error)

	PostOutboundTracker(_ context.Context,
		_ ChainID,
		_ Nonce,
		txHash string,
	) (string, error)
}

// ------------------------------------------------------------------------------------------------
// dry-mode
// ------------------------------------------------------------------------------------------------

const unreachableMsg = "called an unreachable dryZetacoreClient function"

// dryZetacoreClient is a dry-wrapper for the zetacore client.
// It overrides mutating functions from the underlying client that panic when called.
type dryZetacoreClient struct {
	zetacoreReaderClient

	// writerClient is deliberately not embedded so the compiler can ensure that all mutating
	// methods are explicitly overridden.
	writerClient zetacoreWriterClient
}

func newDryZetacoreClient(client ZetacoreClient) *dryZetacoreClient {
	return &dryZetacoreClient{zetacoreReaderClient: client, writerClient: client}
}

func (*dryZetacoreClient) PostVoteGasPrice(context.Context, chains.Chain, uint64, uint64, uint64,
) (string, error) {
	panic(unreachableMsg)
}

func (*dryZetacoreClient) PostVoteTSS(context.Context, string, int64, chains.ReceiveStatus,
) (string, error) {
	panic(unreachableMsg)
}

func (*dryZetacoreClient) PostVoteBlameData(context.Context, *blame.Blame, ChainID, string,
) (string, error) {
	panic(unreachableMsg)
}

func (*dryZetacoreClient) PostVoteOutbound(context.Context, uint64, uint64,
	*crosschain.MsgVoteOutbound,
) (string, string, error) {
	panic(unreachableMsg)
}

func (*dryZetacoreClient) PostVoteInbound(context.Context, uint64, uint64,
	*crosschain.MsgVoteInbound, chan<- zetaerrors.ErrTxMonitor,
) (string, string, error) {
	panic(unreachableMsg)
}

func (*dryZetacoreClient) PostOutboundTracker(context.Context, ChainID, Nonce, string,
) (string, error) {
	panic(unreachableMsg)
}
