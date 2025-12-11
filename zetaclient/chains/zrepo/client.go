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

// ZetacoreClient contains the zetacore client functions used by ZetaRepo.
type ZetacoreClient interface {
	ZetacoreReaderClient
	ZetacoreWriterClient
}

// ZetacoreReaderClient contains the zetacore client functions used by ZetaRepo that do not mutate
// ZetaChain state.
type ZetacoreReaderClient interface {
	Chain() chains.Chain

	NewBlockSubscriber(context.Context) (chan cometbft.EventDataNewBlock, error)

	GetBlockHeight(context.Context) (int64, error)

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

// ZetacoreReaderClient contains the zetacore client functions used by ZetaRepo that do mutate
// ZetaChain state.
type ZetacoreWriterClient interface {
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
