package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keygen"
	"gitlab.com/thorchain/tss/go-tss/tss"
	"golang.org/x/crypto/sha3"

	"github.com/zeta-chain/node/pkg/chains"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/metrics"
	mc "github.com/zeta-chain/node/zetaclient/tss"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// GenerateTSS generates a new TSS if keygen is set.
// If a TSS was generated successfully in the past,and the keygen was successful, the function will return without doing anything.
// If a keygen has been set the functions will wait for the correct block to arrive and generate a new TSS.
// In case of a successful keygen a TSS success vote is broadcasted to zetacore and the newly generate TSS is tested. The generated keyshares are stored in the correct directory
// In case of a failed keygen a TSS failed vote is broadcasted to zetacore.
func GenerateTSS(
	ctx context.Context,
	logger zerolog.Logger,
	zetaCoreClient *zetacore.Client,
	keygenTssServer *tss.TssServer) error {
	keygenLogger := logger.With().Str("module", "keygen").Logger()
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}
	// If Keygen block is set it will try to generate new TSS at the block
	// This is a blocking thread and will wait until the ceremony is complete successfully
	// If the TSS generation is unsuccessful , it will loop indefinitely until a new TSS is generated
	// Set TSS block to 0 using genesis file to disable this feature
	// Note : The TSS generation is done through the "hotkey" or "Zeta-clientGrantee" This key needs to be present on the machine for the TSS signing to happen .
	// "ZetaClientGrantee" key is different from the "operator" key .The "Operator" key gives all zetaclient related permissions such as TSS generation ,reporting and signing, INBOUND and OUTBOUND vote signing, to the "ZetaClientGrantee" key.
	// The votes to signify a successful TSS generation (Or unsuccessful) is signed by the operator key and broadcast to zetacore by the zetcalientGrantee key on behalf of the operator .
	ticker := time.NewTicker(time.Second * 1)
	triedKeygenAtBlock := false
	lastBlock := int64(0)
	for range ticker.C {
		// Break out of loop only when TSS is generated successfully, either at the keygenBlock or if it has been generated already , Block set as zero in genesis file
		// This loop will try keygen at the keygen block and then wait for keygen to be successfully reported by all nodes before breaking out of the loop.
		// If keygen is unsuccessful, it will reset the triedKeygenAtBlock flag and try again at a new keygen block.
		keyGen := app.GetKeygen()
		if keyGen.Status == observertypes.KeygenStatus_KeyGenSuccess {
			return nil
		}
		// Arrive at this stage only if keygen is unsuccessfully reported by every node . This will reset the flag and to try again at a new keygen block
		if keyGen.Status == observertypes.KeygenStatus_KeyGenFailed {
			triedKeygenAtBlock = false
			continue
		}
		// Try generating TSS at keygen block , only when status is pending keygen and generation has not been tried at the block
		if keyGen.Status == observertypes.KeygenStatus_PendingKeygen {
			// Return error if RPC is not working
			currentBlock, err := zetaCoreClient.GetBlockHeight(ctx)
			if err != nil {
				keygenLogger.Error().Err(err).Msg("GetBlockHeight RPC  error")
				continue
			}
			// Reset the flag if the keygen block has passed and a new keygen block has been set . This condition is only reached if the older keygen is stuck at PendingKeygen for some reason
			if keyGen.BlockNumber > currentBlock {
				triedKeygenAtBlock = false
			}
			if !triedKeygenAtBlock {
				// If not at keygen block do not try to generate TSS
				if currentBlock != keyGen.BlockNumber {
					if currentBlock > lastBlock {
						lastBlock = currentBlock
						keygenLogger.Info().
							Msgf("Waiting For Keygen Block to arrive or new keygen block to be set. Keygen Block : %d Current Block : %d ChainID %s ", keyGen.BlockNumber, currentBlock, app.Config().ChainID)
					}
					continue
				}
				// Try keygen only once at a particular block, irrespective of whether it is successful or failure
				triedKeygenAtBlock = true
				newPubkey, err := keygenTSS(ctx, keyGen, *keygenTssServer, zetaCoreClient, keygenLogger)
				if err != nil {
					keygenLogger.Error().Err(err).Msg("keygenTSS error")
					tssFailedVoteHash, err := zetaCoreClient.PostVoteTSS(ctx,
						"", keyGen.BlockNumber, chains.ReceiveStatus_failed)
					if err != nil {
						keygenLogger.Error().Err(err).Msg("Failed to broadcast Failed TSS Vote to zetacore")
						return err
					}
					keygenLogger.Info().Msgf("TSS Failed Vote: %s", tssFailedVoteHash)
					continue
				}
				// If TSS is successful , broadcast the vote to zetacore and also set the Pubkey
				tssSuccessVoteHash, err := zetaCoreClient.PostVoteTSS(ctx,
					newPubkey,
					keyGen.BlockNumber,
					chains.ReceiveStatus_success,
				)
				if err != nil {
					keygenLogger.Error().Err(err).Msg("TSS successful but unable to broadcast vote to zeta-core")
					return err
				}
				keygenLogger.Info().Msgf("TSS successful Vote: %s", tssSuccessVoteHash)

				err = TestTSS(newPubkey, *keygenTssServer, keygenLogger)
				if err != nil {
					keygenLogger.Error().Err(err).Msgf("TestTSS error: %s", newPubkey)
				}
				continue
			}
		}
		keygenLogger.Debug().
			Msgf("Waiting for TSS to be generated or Current Keygen to be be finalized. Keygen Block : %d ", keyGen.BlockNumber)
	}
	return errors.New("unexpected state for TSS generation")
}

// keygenTSS generates a new TSS using the keygen request and the TSS server.
// If the keygen is successful, the function returns the new TSS pubkey.
// If the keygen is unsuccessful, the function posts blame and returns an error.
func keygenTSS(
	ctx context.Context,
	keyGen observertypes.Keygen,
	tssServer tss.TssServer,
	zetacoreClient interfaces.ZetacoreClient,
	keygenLogger zerolog.Logger,
) (string, error) {
	keygenLogger.Info().Msgf("Keygen at blocknum %d , TSS signers %s ", keyGen.BlockNumber, keyGen.GranteePubkeys)
	var req keygen.Request
	req = keygen.NewRequest(keyGen.GranteePubkeys, keyGen.BlockNumber, "0.14.0")
	res, err := tssServer.Keygen(req)
	if res.Status != tsscommon.Success || res.PubKey == "" {
		keygenLogger.Error().Msgf("keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes)
		// Need to broadcast keygen blame result here
		digest, err := digestReq(req)
		if err != nil {
			return "", err
		}
		index := fmt.Sprintf("keygen-%s-%d", digest, keyGen.BlockNumber)
		zetaHash, err := zetacoreClient.PostVoteBlameData(
			ctx,
			&res.Blame,
			zetacoreClient.Chain().ChainId,
			index,
		)
		if err != nil {
			keygenLogger.Error().Err(err).Msg("error sending blame data to core")
			return "", err
		}

		// Increment Blame counter
		for _, node := range res.Blame.BlameNodes {
			metrics.TssNodeBlamePerPubKey.WithLabelValues(node.Pubkey).Inc()
		}

		keygenLogger.Info().Msgf("keygen posted blame data tx hash: %s", zetaHash)
		return "", fmt.Errorf("keygen fail: reason %s blame nodes %s", res.Blame.FailReason, res.Blame.BlameNodes)
	}
	if err != nil {
		keygenLogger.Error().Msgf("keygen fail: reason %s ", err.Error())
		return "", err
	}
	// Keygen succeed
	keygenLogger.Info().Msgf("Keygen success! keygen response: %v", res)
	return res.PubKey, nil
}

// TestTSS tests the TSS keygen by signing a sample message with the TSS key.
func TestTSS(pubkey string, tssServer tss.TssServer, logger zerolog.Logger) error {
	keygenLogger := logger.With().Str("module", "test-keygen").Logger()
	keygenLogger.Info().Msgf("KeyGen success ! Doing a Key-sign test")
	// KeySign can fail even if TSS keygen is successful, just logging the error here to break out of outer loop and report TSS
	err := mc.TestKeysign(pubkey, tssServer)
	if err != nil {
		return err
	}
	return nil
}

func digestReq(request keygen.Request) (string, error) {
	bytes, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(bytes)
	digest := hex.EncodeToString(hasher.Sum(nil))

	return digest, nil
}
