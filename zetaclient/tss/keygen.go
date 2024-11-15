package tss

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keygen"
	"gitlab.com/thorchain/tss/go-tss/tss"
	"golang.org/x/crypto/sha3"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/ticker"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

const (
	receiveSuccess = chains.ReceiveStatus_success
	receiveFailed  = chains.ReceiveStatus_failed
)

type keygenCeremony struct {
	tss           *tss.TssServer
	zetacore      *zetacore.Client
	lastSeenBlock int64
	logger        zerolog.Logger
}

// KeygenCeremony runs TSS keygen ceremony as a blocking thread.
// Most likely the keygen is already generated, so this function will be a noop.
func KeygenCeremony(ctx context.Context, tssServer *tss.TssServer, zc *zetacore.Client, logger zerolog.Logger) error {
	const interval = time.Second

	ceremony := keygenCeremony{
		tss:      tssServer,
		zetacore: zc,
		logger:   logger.With().Str(logs.FieldModule, "tss_keygen").Logger(),
	}

	task := func(ctx context.Context, t *ticker.Ticker) error {
		shouldRetry, err := ceremony.iteration(ctx)
		switch {
		case shouldRetry:
			if err != nil {
				logger.Error().Err(err).Msg("Keygen error. Retrying...")
			}

			// continue the ticker
			return nil
		case err != nil:
			return errors.Wrap(err, "keygen ceremony failed")
		default:
			// keygen ceremony is complete (or noop)
			t.Stop()
			return nil
		}
	}

	return ticker.Run(ctx, interval, task, ticker.WithLogger(logger, "tss_keygen"))
}

// iteration runs ceremony iteration every time interval.
// - Get the keygen task from zetacore
// - If the keygen is already generated, return (false, nil) => ceremony is complete
// - If the keygen is pending, ensure we're on the right block
// - Iteration also ensured that the logic is invoked ONLY once per block (regardless of the interval)
func (k *keygenCeremony) iteration(ctx context.Context) (shouldRetry bool, err error) {
	keygenTask, err := k.zetacore.GetKeyGen(ctx)
	switch {
	case err != nil:
		return true, errors.Wrap(err, "unable to get keygen via RPC")
	case keygenTask.Status == observertypes.KeygenStatus_KeyGenSuccess:
		// all good, tss key is already generated
		return false, nil
	case keygenTask.Status == observertypes.KeygenStatus_KeyGenFailed:
		// come back later to try again (zetacore will make status=pending)
		return true, nil
	case keygenTask.Status == observertypes.KeygenStatus_PendingKeygen:
		// okay, let's try to generate the TSS key
	default:
		return false, fmt.Errorf("unexpected keygen status %q", keygenTask.Status.String())
	}

	keygenHeight := keygenTask.BlockNumber

	zetaHeight, err := k.zetacore.GetBlockHeight(ctx)
	switch {
	case err != nil:
		return true, errors.Wrap(err, "unable to get zeta height")
	case k.blockThrottled(zetaHeight):
		return true, nil
	case zetaHeight < keygenHeight:
		k.logger.Info().
			Int64("keygen.height", keygenHeight).
			Int64("zeta_height", zetaHeight).
			Msgf("Waiting for keygen block to arrive or new keygen block to be set")
		return true, nil
	case zetaHeight > keygenHeight:
		k.logger.Info().
			Int64("keygen.height", keygenHeight).
			Int64("zeta_height", zetaHeight).
			Msgf("Waiting for keygen finalization")
		return true, nil
	}

	// Now we know that the keygen status is PENDING, and we are the KEYGEN block.
	// Let's perform TSS Keygen and then post successful/failed vote to zetacore
	newPubKey, err := k.performKeygen(ctx, keygenTask)
	if err != nil {
		k.logger.Error().Err(err).Msg("Keygen failed. Broadcasting failed TSS vote")

		// Vote for failure
		failedVoteHash, err := k.zetacore.PostVoteTSS(ctx, "", keygenTask.BlockNumber, receiveFailed)
		if err != nil {
			return false, errors.Wrap(err, "failed to broadcast failed TSS vote")
		}

		k.logger.Info().
			Str("keygen.failed_vote_tx_hash", failedVoteHash).
			Msg("Broadcasted failed TSS keygen vote")

		return true, nil
	}

	successVoteHash, err := k.zetacore.PostVoteTSS(ctx, newPubKey, keygenTask.BlockNumber, receiveSuccess)
	if err != nil {
		return false, errors.Wrap(err, "failed to broadcast successful TSS vote")
	}

	k.logger.Info().
		Str("keygen.success_vote_tx_hash", successVoteHash).
		Msg("Broadcasted successful TSS keygen vote")

	k.logger.Info().Msg("Performing TSS key-sign test")

	if err = TestKeySign(k.tss, newPubKey, k.logger); err != nil {
		k.logger.Error().Err(err).Msg("Failed to test TSS keygen")
		// signing can fail even if tss keygen is successful
	}

	return false, nil
}

// performKeygen performs TSS keygen flow via go-tss server. Returns the new TSS public key or error.
// If fails, then it will post blame data to zetacore and return an error.
func (k *keygenCeremony) performKeygen(ctx context.Context, keygenTask observertypes.Keygen) (string, error) {
	k.logger.Warn().
		Int64("keygen.block", keygenTask.BlockNumber).
		Strs("keygen.tss_signers", keygenTask.GranteePubkeys).
		Msg("Performing a keygen!")

	req := keygen.NewRequest(keygenTask.GranteePubkeys, keygenTask.BlockNumber, Version, Algo)

	res, err := k.tss.Keygen(req)
	switch {
	case err != nil:
		// returns error on network failure or other non-recoverable errors
		// if the keygen is unsuccessful, the error will be nil
		return "", errors.Wrap(err, "unable to perform keygen")
	case res.Status == tsscommon.Success && res.PubKey != "":
		// desired outcome
		k.logger.Info().
			Interface("keygen.response", res).
			Interface("keygen.tss_public_key", res.PubKey).
			Msg("Keygen successfully generated!")
		return res.PubKey, nil
	}

	// Something went wrong, let's post blame results and then FAIL
	k.logger.Error().
		Str("keygen.blame_round", res.Blame.Round).
		Str("keygen.fail_reason", res.Blame.FailReason).
		Interface("keygen.blame_nodes", res.Blame.BlameNodes).
		Msg("Keygen failed! Sending blame data to zetacore")

	// increment blame counter
	for _, node := range res.Blame.BlameNodes {
		metrics.TssNodeBlamePerPubKey.WithLabelValues(node.Pubkey).Inc()
	}

	blameDigest, err := digestReq(req)
	if err != nil {
		return "", errors.Wrap(err, "unable to create digest")
	}

	blameIndex := fmt.Sprintf("keygen-%s-%d", blameDigest, keygenTask.BlockNumber)
	chainID := k.zetacore.Chain().ChainId

	zetaHash, err := k.zetacore.PostVoteBlameData(ctx, &res.Blame, chainID, blameIndex)
	if err != nil {
		return "", errors.Wrap(err, "unable to post blame data to zetacore")
	}

	k.logger.Info().Str("keygen.blame_tx_hash", zetaHash).Msg("Posted blame data to zetacore")

	return "", errors.Errorf("keygen failed: %s", res.Blame.FailReason)
}

// returns true if the block is throttled i.e. we should wait for the next block.
func (k *keygenCeremony) blockThrottled(currentBlock int64) bool {
	switch {
	case currentBlock == 0:
		return false
	case k.lastSeenBlock == currentBlock:
		return true
	default:
		k.lastSeenBlock = currentBlock
		return false
	}
}

func digestReq(req keygen.Request) (string, error) {
	bytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(bytes)
	digest := hex.EncodeToString(hasher.Sum(nil))

	return digest, nil
}
