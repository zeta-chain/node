package signer

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/coin"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

const (
	// solanaTransactionTimeout is the timeout for waiting for an outbound to be confirmed.
	// Transaction referencing a blockhash older than 150 blocks (60 ~90 secs) will expire and be rejected by Solana.
	solanaTransactionTimeout = 2 * time.Minute

	// broadcastBackoff is the initial backoff duration for retrying broadcast
	broadcastBackoff = 1 * time.Second

	// broadcastRetries is the maximum number of retries for broadcasting a transaction
	// 6 retries will span over 1 + 2 + 4 + 8 + 16 + 32 + 64 = 127 seconds, good enough for the 2 minute timeout
	broadcastRetries = 7

	// pdaNonceWaitTimeout is the timeout for waiting for the PDA nonce to arrive
	// given 1~2 seconds finality at the 'confirmed' level, 1 minute can cover 30~60 (the lookahead) parallel CCTXs processing
	pdaNonceWaitTimeout = 1 * time.Minute

	// broadcastOutboundCommitment is the commitment level for broadcasting solana outbound.
	// Commitment "finalized" eliminate all risk but the tradeoff is pretty severe and effectively
	// reduces the expiration of tx by about 13 seconds. The "confirmed" level has very low risk of
	// belonging to a dropped fork.
	// see: https://solana.com/developers/guides/advanced/confirmation#use-an-appropriate-preflight-commitment-level
	broadcastOutboundCommitment = rpc.CommitmentConfirmed

	// SolanaMaxComputeBudget is the max compute budget for a transaction.
	SolanaMaxComputeBudget = 1_400_000
)

type Outbound struct {
	Tx          *solana.Transaction
	FallbackMsg *contracts.MsgIncrementNonce
}

type outboundGetter func() (*Outbound, error)

// Signer deals with signing Solana transactions and implements the ChainSigner interface
type Signer struct {
	*base.Signer

	// client is the Solana RPC client that interacts with the Solana chain
	client interfaces.SolanaRPCClient

	// relayerKey is the private key of the relayer account for Solana chain
	// relayerKey is optional, the signer will not relay transactions if it is not set
	relayerKey *solana.PrivateKey

	// gatewayID is the program ID of gateway program on Solana chain
	gatewayID solana.PublicKey

	// pda is the program derived address of the gateway program
	pda solana.PublicKey
}

// New Signer constructor.
func New(
	baseSigner *base.Signer,
	solClient interfaces.SolanaRPCClient,
	gatewayAddress string,
	relayerKey *keys.RelayerKey,
) (*Signer, error) {
	// parse gateway ID and PDA
	gatewayID, pda, err := contracts.ParseGatewayWithPDA(gatewayAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse gateway address %s", gatewayAddress)
	}

	var rk *solana.PrivateKey

	if relayerKey != nil {
		pk, err := solana.PrivateKeyFromBase58(relayerKey.PrivateKey)
		if err != nil {
			return nil, errors.Wrap(err, "unable to construct solana private key")
		}

		rk = &pk
		baseSigner.Logger().Std.Info().
			Stringer("relayer_key", rk.PublicKey()).
			Msg("loaded relayer key")
	} else {
		baseSigner.Logger().Std.Info().Msg("solana relayer key was not provided")
	}

	return &Signer{
		Signer:     baseSigner,
		client:     solClient,
		gatewayID:  gatewayID,
		relayerKey: rk,
		pda:        pda,
	}, nil
}

// HasRelayerKey returns true if the signer has a relayer key
func (signer *Signer) HasRelayerKey() bool {
	return signer.relayerKey != nil
}

// TryProcessOutbound - signer interface implementation
// This function will attempt to build and sign a Solana transaction using the TSS signer.
// It will then broadcast the signed transaction to the Solana chain.
func (signer *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *types.CrossChainTx,
	zetacoreClient interfaces.ZetacoreClient,
	height uint64,
) {
	outboundID := base.OutboundIDFromCCTX(cctx)
	signer.MarkOutbound(outboundID, true)

	// end outbound process on panic
	defer func() {
		signer.MarkOutbound(outboundID, false)
		if err := recover(); err != nil {
			signer.Logger().Std.Error().
				Str(logs.FieldCctx, cctx.Index).
				Any("panic", err).
				Str("stack_trace", string(debug.Stack())).
				Msg("caught panic error")
		}
	}()

	// prepare logger
	params := cctx.GetCurrentOutboundParam()
	logger := signer.Logger().Std.With().
		Uint64("nonce", params.TssNonce).
		Str("cctx", cctx.Index).
		Logger()

	var (
		chainID        = signer.Chain().ChainId
		nonce          = params.TssNonce
		coinType       = cctx.InboundParams.CoinType
		isRevert       = (cctx.CctxStatus.Status == types.CctxStatus_PendingRevert && cctx.RevertOptions.CallOnRevert)
		cancelTx       = !signer.PassesCompliance(cctx)
		outboundGetter outboundGetter
	)

	switch coinType {
	case coin.CoinType_Cmd:
		whitelistTxGetter, err := signer.prepareWhitelistTx(ctx, cctx, height)
		if err != nil {
			logger.Error().Err(err).Msg("failed to sign whitelist outbound")
			return
		}

		outboundGetter = whitelistTxGetter

	case coin.CoinType_Gas:
		isRevert := (cctx.CctxStatus.Status == types.CctxStatus_PendingRevert && cctx.RevertOptions.CallOnRevert)
		if cctx.IsWithdrawAndCall() || isRevert {
			executeTxGetter, err := signer.prepareExecuteTx(ctx, cctx, height, cancelTx, logger)
			if err != nil {
				logger.Error().Err(err).Msg("failed to sign execute outbound")
				return
			}

			outboundGetter = executeTxGetter
		} else {
			withdrawTxGetter, err := signer.prepareWithdrawTx(ctx, cctx, height, cancelTx, logger)
			if err != nil {
				logger.Error().Err(err).Msg("failed to sign withdraw outbound")
				return
			}

			outboundGetter = withdrawTxGetter
		}

	case coin.CoinType_ERC20:
		if cctx.IsWithdrawAndCall() || isRevert {
			executeSPLTxGetter, err := signer.prepareExecuteSPLTx(ctx, cctx, height, cancelTx, logger)
			if err != nil {
				logger.Error().Err(err).Msg("failed to sign execute spl outbound")
				return
			}

			outboundGetter = executeSPLTxGetter
		} else {
			withdrawSPLTxGetter, err := signer.prepareWithdrawSPLTx(ctx, cctx, height, cancelTx, logger)
			if err != nil {
				logger.Error().Err(err).Msg("failed to sign withdraw spl outbound")
				return
			}

			outboundGetter = withdrawSPLTxGetter
		}
	case coin.CoinType_NoAssetCall:
		executeTxGetter, err := signer.prepareExecuteTx(ctx, cctx, height, cancelTx, logger)
		if err != nil {
			logger.Error().Err(err).Msg("failed to sign execute outbound")
			return
		}

		outboundGetter = executeTxGetter
	default:
		logger.Error().Msg("can only send SOL to the Solana network")
		return
	}

	// skip relaying the transaction if this signer hasn't set the relayer key
	if !signer.HasRelayerKey() {
		logger.Warn().Msg("no relayer key configured")
		return
	}

	// set relayer balance metrics
	signer.SetRelayerBalanceMetrics(ctx)

	// wait for the exact PDA nonce to arrive with timeout
	ctxWait, cancel := context.WithTimeout(ctx, pdaNonceWaitTimeout)
	defer cancel()

	if err := signer.waitExactGatewayNonce(ctxWait, params.TssNonce); err != nil {
		logger.Error().Err(err).Msg("failed to wait for gateway nonce")
		return
	}

	// Get transactions from getters
	// This is when the recent block hash timer starts
	outbound, err := outboundGetter()
	if err != nil {
		logger.Error().Err(err).Msg("failed to get transaction")
		return
	}

	// broadcast the signed tx to the Solana network
	signer.broadcastOutbound(ctx, outbound, chainID, nonce, logger, zetacoreClient)
}

// signTx creates and signs solana tx containing provided instruction with relayer key.
func (signer *Signer) signTx(
	ctx context.Context,
	inst *solana.GenericInstruction,
	limit uint64,
) (*solana.Transaction, error) {
	// get a recent blockhash
	recent, err := signer.client.GetLatestBlockhash(ctx, broadcastOutboundCommitment)
	if err != nil {
		return nil, errors.Wrap(err, "getLatestBlockhash error")
	}

	// if limit is provided, prepend compute unit limit instruction
	var instructions []solana.Instruction
	if limit > 0 {
		limit = min(limit, SolanaMaxComputeBudget)
		// #nosec G115 always in range
		limitInst := computebudget.NewSetComputeUnitLimitInstruction(uint32(limit)).Build()
		instructions = append(instructions, limitInst)
	}

	instructions = append(instructions, inst)

	// create a transaction that wraps the instruction
	tx, err := solana.NewTransaction(
		// TODO: outbound now uses 5K lamports as the fixed fee, we could explore priority fee and compute budget
		// https://github.com/zeta-chain/node/issues/2599
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(signer.relayerKey.PublicKey()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create new tx")
	}

	// relayer signs the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(signer.relayerKey.PublicKey()) {
			return signer.relayerKey
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "signer unable to sign transaction")
	}

	return tx, nil
}

// broadcastOutbound sends the signed transaction to the Solana network
func (signer *Signer) broadcastOutbound(
	ctx context.Context,
	outbound *Outbound,
	chainID int64,
	nonce uint64,
	logger zerolog.Logger,
	zetacoreClient interfaces.ZetacoreClient,
) {
	tx := outbound.Tx
	// prepare logger fields
	lf := map[string]any{
		logs.FieldNonce: nonce,
		logs.FieldTx:    tx.Signatures[0].String(),
	}

	// try broacasting tx with increasing backoff (1s, 2s, 4s, 8s, 16s, 32s, 64s)
	// to tolerate tx nonce mismatch with PDA nonce or unknown RPC error
	backOff := broadcastBackoff
	for range broadcastRetries {
		time.Sleep(backOff)

		// query the gateway PDA nonce
		pdaNonce, err := signer.getGatewayNonce(ctx)
		if err != nil {
			logger.Error().Err(err).Fields(lf).Msg("unable to get PDA nonce")
			backOff *= 2
			continue
		}
		lf["pda_nonce"] = pdaNonce

		// PDA nonce may already be increased by other relayers, no need to retry
		if pdaNonce > nonce {
			logger.Info().Fields(lf).Msg("PDA nonce is greater than outbound nonce, stop retrying")
			break
		}

		// broadcast the signed tx to the Solana network with preflight check
		// the PDA nonce MUST be equal to 'nonce' if arrived here, guaranteed by upstream code
		txSig, err := signer.client.SendTransactionWithOpts(
			ctx,
			tx,
			rpc.TransactionOpts{PreflightCommitment: broadcastOutboundCommitment},
		)
		if err != nil {
			shouldUseFallbackTx, failureReason := parseRPCErrorForFallback(err, signer.GetGatewayAddress())
			if outbound.FallbackMsg != nil && shouldUseFallbackTx {
				// create and sign fallback transaction
				outbound.FallbackMsg.SetFailureReason(failureReason)
				fallbackInst, err := signer.createIncrementNonceInstruction(*outbound.FallbackMsg)
				if err != nil {
					logger.Error().Err(err).Fields(lf).Msg("error creating increment nonce instruction")
					break
				}

				fallbackTx, err := signer.signTx(ctx, fallbackInst, 0)
				if err != nil {
					logger.Error().Err(err).Fields(lf).Msg("error signing increment nonce instruction")
					break
				}
				tx = fallbackTx
			}
			logger.Warn().Err(err).Fields(lf).Msg("error calling SendTransactionWithOpts")
			backOff *= 2
			continue
		}
		logger.Info().Fields(lf).Msg("broadcasted Solana outbound successfully")

		// successful broadcast; report to the outbound tracker
		signer.reportToOutboundTracker(ctx, zetacoreClient, chainID, nonce, txSig, logger)
		break
	}
}

// createOutboundWithFallback is a helper function that creates an outbound with a main and a fallback transaction
// and signs them with relayer key
func (signer *Signer) createOutboundWithFallback(
	ctx context.Context,
	mainInst *solana.GenericInstruction,
	msgIn *contracts.MsgIncrementNonce,
	computeLimit uint64,
) (*Outbound, error) {
	// Create and sign main transaction
	tx, err := signer.signTx(ctx, mainInst, computeLimit)
	if err != nil {
		return nil, errors.Wrap(err, "error signing main instruction")
	}

	return &Outbound{
		Tx:          tx,
		FallbackMsg: msgIn,
	}, nil
}

// SetGatewayAddress sets the gateway address
func (signer *Signer) SetGatewayAddress(address string) {
	// noop
	if address == "" || signer.gatewayID.String() == address {
		return
	}

	// parse gateway ID and PDA
	gatewayID, pda, err := contracts.ParseGatewayWithPDA(address)
	if err != nil {
		signer.Logger().Std.Error().
			Err(err).
			Str("address", address).
			Msg("error parsing the gateway address")
		return
	}

	// noop
	if signer.gatewayID.Equals(gatewayID) {
		return
	}

	signer.Logger().Std.Info().
		Stringer("signer_old_gateway_address", signer.gatewayID).
		Stringer("signer_new_gateway_address", gatewayID).
		Msg("updated the gateway address")

	signer.Lock()
	signer.gatewayID = gatewayID
	signer.pda = pda
	signer.Unlock()
}

// GetGatewayAddress returns the gateway address
func (signer *Signer) GetGatewayAddress() string {
	return signer.gatewayID.String()
}

// SetRelayerBalanceMetrics sets the relayer balance metrics
func (signer *Signer) SetRelayerBalanceMetrics(ctx context.Context) {
	if !signer.HasRelayerKey() {
		return
	}

	result, err := signer.client.GetBalance(ctx, signer.relayerKey.PublicKey(), rpc.CommitmentFinalized)
	if err != nil {
		signer.Logger().Std.Error().Err(err).Msg("error calling GetBalance")
		return
	}
	solBalance := float64(result.Value) / float64(solana.LAMPORTS_PER_SOL)
	metrics.RelayerKeyBalance.WithLabelValues(signer.Chain().Name).Set(solBalance)
}

// IsPendingOutboundFromZetaChain checks if the sender chain is ZetaChain and if status is PendingOutbound
// TODO(revamp): move to another package more general for cctx functions
func IsPendingOutboundFromZetaChain(
	cctx *types.CrossChainTx,
	zetacoreClient interfaces.ZetacoreClient,
) bool {
	return cctx.InboundParams.SenderChainId == zetacoreClient.Chain().ChainId &&
		cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound
}

type SignableMessage[T any] interface {
	Hash() [32]byte
	SetSignature([65]byte) T
}

// signMsgWithFallback TSS signs solana outbound with fallback increment nonce
func signMsgWithFallback[T SignableMessage[T]](
	ctx context.Context,
	signer *Signer,
	height, nonce uint64,
	msg T,
	msgIn *contracts.MsgIncrementNonce,
) (T, *contracts.MsgIncrementNonce, error) {
	msgHash := msg.Hash()
	msgInHash := msgIn.Hash()

	signature, err := signer.TSS().
		SignBatch(ctx, [][]byte{msgHash[:], msgInHash[:]}, height, nonce, signer.Chain().ChainId)
	if err != nil {
		var zero T
		return zero, nil, errors.Wrap(err, "key-sign failed")
	}

	return msg.SetSignature(signature[0]), msgIn.SetSignature(signature[1]), nil
}

// waitExactGatewayNonce waits for exact given gateway nonce to arrive
//
// the reasons are:
//  1. any pre-signed Solana tx expires after 150 blocks, so we should avoid pre-signing any tx to maximize the lifetime of signed txs
//  2. there can be up to 'lookahead' CCTX processing goroutines running in parallel, so waiting for PDA nonce helps to order the CCTX
//     processing goroutines by nonce and avoid nonce mismatch
//  3. less nonce mismatch will reduce CCTX retries and TSS keysign requests
func (signer *Signer) waitExactGatewayNonce(ctx context.Context, nonce uint64) error {
	logger := signer.Logger().Std.With().
		Str("method", "waitExactGatewayNonce").
		Int64("chain", signer.Chain().ChainId).
		Uint64("nonce", nonce).
		Logger()

	for {
		if ctx.Err() != nil {
			return errors.Wrap(ctx.Err(), "context cancelled")
		}

		// check timeout to avoid infinite waiting
		if deadline, ok := ctx.Deadline(); ok {
			if time.Now().After(deadline) {
				return errors.New("timeout reached on waiting for gateway nonce")
			}
		}

		// query the gateway PDA nonce
		pdaNonce, err := signer.getGatewayNonce(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("unable to get gateway nonce")
			time.Sleep(time.Second) // prevent RPC spamming
			continue
		}

		switch {
		case pdaNonce > nonce:
			return errors.Wrapf(err, "PDA nonce %d is greater than outbound nonce %d", pdaNonce, nonce)
		case pdaNonce == nonce:
			return nil
		default:
			logger.Info().Uint64("pda_nonce", pdaNonce).Msg("waiting for PDA nonce to arrive")

			// calculate how far behind the PDA nonce and sleep accordingly
			//  - base sleep time of 1 second, multiplied by the nonce difference
			//  - 'lookahead' parameter should keep this from getting too out of control
			// #nosec G115 always in range
			sleepDuration := time.Second * time.Duration(nonce-pdaNonce)
			time.Sleep(sleepDuration)
		}
	}
}

// getGatewayNonce queries the gateway nonce from the PDA account information
func (signer *Signer) getGatewayNonce(ctx context.Context) (uint64, error) {
	// query the gateway PDA account information
	pdaInfo, err := signer.client.GetAccountInfoWithOpts(
		ctx,
		signer.pda,
		&rpc.GetAccountInfoOpts{Commitment: broadcastOutboundCommitment},
	)
	if err != nil {
		return 0, errors.Wrap(err, "unable to get gateway PDA account info")
	}

	// deserialize the PDA account information
	pda, err := contracts.DeserializePdaInfo(pdaInfo)
	if err != nil {
		return 0, errors.Wrap(err, "unable to deserialize PDA info")
	}

	return pda.Nonce, nil
}
