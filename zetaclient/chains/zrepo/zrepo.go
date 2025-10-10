// Package zrepo provides an abstraction layer for interactions with the zetacore client.
//
// The functions inside this module return very descriptive errors.
// There is no need to wrap these errors with "failed to..." messages.
package zrepo

import (
	"context"
	"fmt"
	"time"

	cometbft "github.com/cometbft/cometbft/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	"github.com/zeta-chain/node/pkg/chains"
	zetaerrors "github.com/zeta-chain/node/pkg/errors"
	cc "github.com/zeta-chain/node/x/crosschain/types"
	fungible "github.com/zeta-chain/node/x/fungible/types"
	observer "github.com/zeta-chain/node/x/observer/types"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/mode"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// MonitoringErrHandlerRoutineTimeout is the timeout for the handleMonitoring routine that waits for an error from the monitorVote channel
const monitoringErrHandlerRoutineTimeout = 5 * time.Minute

// MonitoringErrorWatcher is the function type for watching inbound vote monitoring errors.
type MonitoringErrorWatcher func(ctx context.Context, monitorErrCh <-chan zetaerrors.ErrTxMonitor, zetaTxHash string)

// ZetaRepo implements the Repository pattern by wrapping a zetacore client.
// Each chain module must instantiate its own ZetaRepo.
type ZetaRepo struct {
	client ZetacoreClient

	connectedChain chains.Chain

	clientMode mode.ClientMode
}

// New constructs a new ZetaRepo object.
func New(client ZetacoreClient, connectedChain chains.Chain, clientMode mode.ClientMode) *ZetaRepo {
	if client == nil {
		return nil
	}
	if clientMode.IsDryMode() {
		client = newDryZetacoreClient(client)
	}
	return &ZetaRepo{client, connectedChain, clientMode}
}

// ------------------------------------------------------------------------------------------------
// Getters
// ------------------------------------------------------------------------------------------------

func (repo *ZetaRepo) ZetaChain() chains.Chain {
	return repo.client.Chain()
}

func (repo *ZetaRepo) GetCCTX(ctx context.Context, nonce uint64) (*cc.CrossChainTx, error) {
	cctx, err := repo.client.GetCctxByNonce(ctx, repo.connectedChain.ChainId, nonce)
	if err != nil {
		outerErr := fmt.Errorf("%w %d", ErrClientGetCCTX, nonce)
		return nil, newClientError(outerErr, err)
	}
	return cctx, nil
}

func (repo *ZetaRepo) GetPendingCCTXs(ctx context.Context) ([]*cc.CrossChainTx, error) {
	cctxs, _, err := repo.client.ListPendingCCTX(ctx, repo.connectedChain)
	if err != nil {
		return nil, newClientError(ErrClientGetPendingCCTXs, err)
	}
	return cctxs, nil
}

func (repo *ZetaRepo) GetPendingNonces(ctx context.Context) (*observer.PendingNonces, error) {
	nonces, err := repo.client.GetPendingNoncesByChain(ctx, repo.connectedChain.ChainId)
	if err != nil {
		return nil, newClientError(ErrClientGetPendingNonces, err)
	}
	return &nonces, nil
}

func (repo *ZetaRepo) GetInboundTrackers(ctx context.Context) ([]cc.InboundTracker, error) {
	trackers, err := repo.client.GetInboundTrackersForChain(ctx, repo.connectedChain.ChainId)
	if err != nil {
		return nil, newClientError(ErrClientGetInboundTrackers, err)
	}
	return trackers, nil
}

func (repo *ZetaRepo) GetOutboundTrackers(ctx context.Context) ([]cc.OutboundTracker, error) {
	trackers, err := repo.client.GetOutboundTrackers(ctx, repo.connectedChain.ChainId)
	if err != nil {
		return nil, newClientError(ErrClientGetOutboundTrackers, err)
	}
	return trackers, nil
}

// TODO: We should probably move this to the TSS repository.
// See: https://github.com/zeta-chain/node/issues/4304
func (repo *ZetaRepo) GetBTCTSSAddress(ctx context.Context) (string, error) {
	address, err := repo.client.GetBTCTSSAddress(ctx, repo.connectedChain.ChainId)
	if err != nil {
		chainID := repo.connectedChain.ChainId
		outerErr := fmt.Errorf("%w for chain %d", ErrClientGetBTCTSSAddress, chainID)
		return "", newClientError(outerErr, err)
	}
	return address, nil
}

func (repo *ZetaRepo) HasVoted(ctx context.Context,
	ballotIndex string,
	voterAddress string,
) (bool, error) {
	return repo.client.HasVoted(ctx, ballotIndex, voterAddress)
}

// ------------------------------------------------------------------------------------------------
// Voting & Posting Trackers
// ------------------------------------------------------------------------------------------------

// PostOutboundTracker posts an outbound tracker.
// It returns the hash of the associated ZetaChain transaction.
func (repo *ZetaRepo) PostOutboundTracker(ctx context.Context, logger zerolog.Logger,
	nonce uint64,
	txHash string,
) (string, error) {
	// Does not post outbound trackers in dry mode.
	if repo.clientMode.IsDryMode() {
		logger.Info().Stringer(logs.FieldMode, mode.DryMode).Msg("skipping outbound tracker")
		return "", nil
	}

	zhash, err := repo.client.PostOutboundTracker(ctx, repo.connectedChain.ChainId, nonce, txHash)
	if err != nil {
		err = newClientError(ErrClientPostOutboundTracker, err)
		logger.Error().Err(err).Send()
		return "", err
	}

	if zhash == "" {
		logger.Info().Msg("outbound tracker already exists")
	} else {
		logger.Info().Str(logs.FieldZetaTx, zhash).Msg("added outbound tracker")
	}

	return zhash, nil
}

// VoteGasPrice votes on gas prices.
// It returns the hash of the vote transaction.
func (repo *ZetaRepo) VoteGasPrice(ctx context.Context, logger zerolog.Logger,
	gasPrice uint64,
	priorityFee uint64,
	block uint64,
) (string, error) {
	// Does not vote in dry mode.
	if repo.clientMode.IsDryMode() {
		logger.Info().Stringer(logs.FieldMode, mode.DryMode).Msg("skipping gas price vote")
		return "", nil
	}

	zhash, err := repo.client.PostVoteGasPrice(ctx, repo.connectedChain, gasPrice, priorityFee, block)
	if err != nil {
		err = newClientError(ErrClientVoteGasPrice, err)
		logger.Error().Err(err).Send()
		return "", err
	}
	return zhash, nil
}

// VoteInbound votes on an inbound.
// It skips invalid messages and tries to vote on the ballot even if the inbound already has a CCTX.
// It returns the index of the ballot.
func (repo *ZetaRepo) VoteInbound(ctx context.Context, logger zerolog.Logger,
	msg *cc.MsgVoteInbound,
	retryGasLimit uint64,
	monitorErrWatcher MonitoringErrorWatcher,
) (string, error) {
	logger = logger.With().
		Str(logs.FieldTx, msg.InboundHash).
		Stringer(logs.FieldCoinType, msg.CoinType).
		Stringer("confirmation_mode", msg.ConfirmationMode).
		Logger()

	{
		// A CCTX is created after an inbound ballot is finalized.
		// - If the CCTX already exists, we try voting if the finalized ballot is still present.
		// - If the CCTX exists but the ballot does not exist, we do not vote.

		cctxIndex := msg.Digest()

		cctxExists, err := repo.CCTXExists(ctx, cctxIndex)
		if err != nil {
			return "", err
		}

		if cctxExists {
			ballotExists, err := repo.ballotExists(ctx, cctxIndex)
			if err != nil {
				return "", err
			}

			if !ballotExists {
				logger.Info().Msg("not voting on inbound; CCTX exists but the ballot does not")
				return cctxIndex, nil
			}
		}
	}

	// Validates the message to avoid unnecessary retries.
	if err := msg.ValidateBasic(); err != nil {
		logger.Warn().Err(err).Msg("invalid vote-inbound message")
		return "", nil
	}

	// Does not vote in dry mode.
	if repo.clientMode.IsDryMode() {
		logger.Info().Stringer(logs.FieldMode, mode.DryMode).Msg("skipping inbound vote")
		return "", nil
	}

	// ctxWithTimeout is a context with timeout used for monitoring the vote transaction
	// Note: the canceller is not used because we want to allow the goroutines to run until they time out
	ctxWithTimeout, _ := zctx.CopyWithTimeout(ctx, context.Background(), monitoringErrHandlerRoutineTimeout)

	// Post vote to zetacore.
	const gasLimit = zetacore.PostVoteInboundGasLimit
	monitorErrCh := make(chan zetaerrors.ErrTxMonitor, 1)
	zhash, ballot, err := repo.client.PostVoteInbound(ctxWithTimeout, gasLimit, retryGasLimit, msg, monitorErrCh)
	if err != nil {
		err = newClientError(ErrClientVoteInbound, err)
		logger.Error().Err(err).Send()
		return "", err
	}

	logger = logger.With().Str(logs.FieldBallotIndex, ballot).Logger()
	if zhash == "" {
		logger.Info().Msg("already voted on the inbound")
	} else {
		logger.Info().Str(logs.FieldZetaTx, zhash).Msg("posted inbound vote")

		// watch for monitoring error for this vote
		if monitorErrWatcher != nil {
			go func() {
				monitorErrWatcher(ctxWithTimeout, monitorErrCh, zhash)
			}()
		}
	}

	return ballot, nil
}

// VoteOutbound votes on an outbound.
// It returns the hash of the vote transaction and the index of the ballot.
func (repo *ZetaRepo) VoteOutbound(ctx context.Context, logger zerolog.Logger,
	gasLimit uint64,
	retryGasLimit uint64,
	msg *cc.MsgVoteOutbound,
) (string, string, error) {
	// Does not vote in dry mode.
	if repo.clientMode.IsDryMode() {
		logger.Info().Stringer(logs.FieldMode, mode.DryMode).Msg("skipping outbound vote")
		return "", "", nil
	}

	zhash, ballot, err := repo.client.PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)
	if err != nil {
		err = newClientError(ErrClientVoteOutbound, err)
		logger.Error().Err(err).Send()
		return "", "", err
	}

	logger = logger.With().Str(logs.FieldBallotIndex, ballot).Logger()
	if zhash == "" {
		logger.Info().Msg("already voted on the outbound")
	} else {
		logger.Info().Str(logs.FieldZetaTx, zhash).Msg("posted outbound vote")
	}

	return zhash, ballot, nil
}

// ------------------------------------------------------------------------------------------------
// Misc.
// ------------------------------------------------------------------------------------------------

// WatchNewBlocks subscribes to new block events.
func (repo *ZetaRepo) WatchNewBlocks(ctx context.Context) (chan cometbft.EventDataNewBlock, error) {
	ch, err := repo.client.NewBlockSubscriber(ctx)
	if err != nil {
		return nil, newClientError(ErrClientNewBlockSubscriber, err)
	}
	return ch, nil
}

// TODO: This function seems out of place.
// See: https://github.com/zeta-chain/node/issues/4304
func (repo *ZetaRepo) GetKeysAddress() (string, error) {
	address, err := repo.client.GetKeys().GetAddress()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrGetKeysAddress, err)
	}
	return address.String(), nil
}

// TODO: This function seems out of place.
// See: https://github.com/zeta-chain/node/issues/4304
func (repo *ZetaRepo) GetOperatorAddress() string {
	return repo.client.GetKeys().GetOperatorAddress().String()
}

func (repo *ZetaRepo) GetForeignCoinsFromAsset(ctx context.Context,
	asset string,
) (*fungible.ForeignCoins, error) {
	chainID := repo.connectedChain.ChainId
	address := ethcommon.HexToAddress(asset)
	coins, err := repo.client.GetForeignCoinsFromAsset(ctx, chainID, address)
	if err != nil {
		outerErr := fmt.Errorf("%w %s", ErrClientGetForeignCoinsForAsset, asset)
		return nil, newClientError(outerErr, err)
	}
	return &coins, nil
}

// ------------------------------------------------------------------------------------------------
// Auxiliary functions
// ------------------------------------------------------------------------------------------------

// checkCode returns nil if the error is a GRPC error with the given code.
//
// We pass code as a parameter because:
// - GetBallotByID returns NotFound when it does not find a ballot.
// - GetCctxByHash returns InvalidArgument when it does not find a CCTX (instead of NotFound).
func checkCode(err error, code grpccodes.Code) error {
	status, ok := grpcstatus.FromError(err) // get the GRPC status from the error
	if !ok {
		return fmt.Errorf("%w: %w", ErrNotRPCError, err) // fail if it is not a GRPC error
	}

	if status.Code() != code {
		return err // fail if it is a different code
	}

	return nil
}

// exists is the generic function used by cctxExists and ballotExists to check for CCTXs and
// ballots.
func exists[T any](ctx context.Context,
	hashOrID string, // hash of a CCTX or the ID of a ballot
	f func(context.Context, string) (*T, error), // client function
	code grpccodes.Code, // code that gets returned by the client if the entity does not exist
	outerErr error, // error that will be returned if the client fails
) (bool, error) {
	res, err := f(ctx, hashOrID)
	if err != nil {
		err = checkCode(err, code)
		if err != nil {
			return false, newClientError(outerErr, err)
		}
		return false, nil
	}
	return res != nil, nil
}

func (repo *ZetaRepo) CCTXExists(ctx context.Context, hash string) (bool, error) {
	f := repo.client.GetCctxByHash
	return exists(ctx, hash, f, grpccodes.InvalidArgument, ErrClientGetCCTXByHash)
}

func (repo *ZetaRepo) ballotExists(ctx context.Context, id string) (bool, error) {
	f := repo.client.GetBallotByID
	return exists(ctx, id, f, grpccodes.NotFound, ErrClientGetBallotByID)
}
