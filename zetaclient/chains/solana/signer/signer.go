package signer

import (
	"context"
	"runtime/debug"
	"time"

	"cosmossdk.io/errors"
	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/rpc"
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

	// SolanaMaxComputeBudget is the max compute budget for a transaction.
	SolanaMaxComputeBudget = 1_400_000
)

type Outbound struct {
	Tx         *solana.Transaction
	FallbackTx *solana.Transaction
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
		baseSigner.Logger().Std.Info().Stringer("relayer_key", rk.PublicKey()).Msg("Loaded relayer key")
	} else {
		baseSigner.Logger().Std.Info().Msg("Solana relayer key is not provided")
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
			signer.Logger().
				Std.Error().
				Str(logs.FieldMethod, "TryProcessOutbound").
				Str(logs.FieldCctx, cctx.Index).
				Interface("panic", err).
				Str("stack_trace", string(debug.Stack())).
				Msg("caught panic error")
		}
	}()

	// prepare logger
	params := cctx.GetCurrentOutboundParam()
	logger := signer.Logger().Std.With().
		Str("method", "TryProcessOutbound").
		Int64("chain", signer.Chain().ChainId).
		Uint64("nonce", params.TssNonce).
		Str("cctx", cctx.Index).
		Logger()

	// support gas token only for Solana outbound
	chainID := signer.Chain().ChainId
	nonce := params.TssNonce
	coinType := cctx.InboundParams.CoinType
	isRevert := (cctx.CctxStatus.Status == types.CctxStatus_PendingRevert && cctx.RevertOptions.CallOnRevert)

	var outboundGetter outboundGetter

	switch coinType {
	case coin.CoinType_Cmd:
		whitelistTxGetter, err := signer.prepareWhitelistTx(ctx, cctx, height)
		if err != nil {
			logger.Error().Err(err).Msgf("TryProcessOutbound: Fail to sign whitelist outbound")
			return
		}

		outboundGetter = whitelistTxGetter

	case coin.CoinType_Gas:
		isRevert := (cctx.CctxStatus.Status == types.CctxStatus_PendingRevert && cctx.RevertOptions.CallOnRevert)
		if cctx.IsWithdrawAndCall() || isRevert {
			executeTxGetter, err := signer.prepareExecuteTx(ctx, cctx, height, logger)
			if err != nil {
				logger.Error().Err(err).Msgf("TryProcessOutbound: Fail to sign execute outbound")
				return
			}

			outboundGetter = executeTxGetter
		} else {
			withdrawTxGetter, err := signer.prepareWithdrawTx(ctx, cctx, height, logger)
			if err != nil {
				logger.Error().Err(err).Msgf("TryProcessOutbound: Fail to sign withdraw outbound")
				return
			}

			outboundGetter = withdrawTxGetter
		}

	case coin.CoinType_ERC20:
		if cctx.IsWithdrawAndCall() || isRevert {
			executeSPLTxGetter, err := signer.prepareExecuteSPLTx(ctx, cctx, height, logger)
			if err != nil {
				logger.Error().Err(err).Msgf("TryProcessOutbound: Fail to sign execute spl outbound")
				return
			}

			outboundGetter = executeSPLTxGetter
		} else {
			withdrawSPLTxGetter, err := signer.prepareWithdrawSPLTx(ctx, cctx, height, logger)
			if err != nil {
				logger.Error().Err(err).Msgf("TryProcessOutbound: Fail to sign withdraw spl outbound")
				return
			}

			outboundGetter = withdrawSPLTxGetter
		}
	case coin.CoinType_NoAssetCall:
		executeTxGetter, err := signer.prepareExecuteTx(ctx, cctx, height, logger)
		if err != nil {
			logger.Error().Err(err).Msgf("TryProcessOutbound: Fail to sign execute outbound")
			return
		}

		outboundGetter = executeTxGetter
	default:
		logger.Error().
			Msgf("TryProcessOutbound: can only send SOL to the Solana network")
		return
	}

	// skip relaying the transaction if this signer hasn't set the relayer key
	if !signer.HasRelayerKey() {
		logger.Warn().Msgf("TryProcessOutbound: no relayer key configured")
		return
	}

	// set relayer balance metrics
	signer.SetRelayerBalanceMetrics(ctx)

	// Get transactions from getters
	// This is when the recent block hash timer starts
	outbound, err := outboundGetter()
	if err != nil {
		logger.Error().Err(err).Msgf("TryProcessOutbound: Failed to get transaction")
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
	recent, err := signer.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
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
		logs.FieldMethod: "broadcastOutbound",
		logs.FieldNonce:  nonce,
		logs.FieldTx:     tx.Signatures[0].String(),
	}

	// try broacasting tx with increasing backoff (1s, 2s, 4s, 8s, 16s, 32s, 64s)
	// to tolerate tx nonce mismatch with PDA nonce or unknown RPC error
	backOff := broadcastBackoff
	for i := 0; i < broadcastRetries; i++ {
		time.Sleep(backOff)

		// PDA nonce may already be increased by other relayer, no need to retry
		pdaInfo, err := signer.client.GetAccountInfo(ctx, signer.pda)
		if err != nil {
			logger.Error().Err(err).Fields(lf).Msgf("unable to get PDA account info")
		} else {
			pda, err := contracts.DeserializePdaInfo(pdaInfo)
			if err != nil {
				logger.Error().Err(err).Fields(lf).Msgf("unable to deserialize PDA info")
			} else if pda.Nonce > nonce {
				logger.Info().Err(err).Fields(lf).Msgf("PDA nonce %d is greater than outbound nonce, stop retrying", pda.Nonce)
				break
			}
		}

		// broadcast the signed tx to the Solana network with preflight check
		txSig, err := signer.client.SendTransactionWithOpts(
			ctx,
			tx,
			// Commitment "finalized" is too conservative for preflight check and
			// it results in repeated broadcast attempts that only 1 will succeed.
			// Commitment "processed" will simulate tx against more recent state
			// thus fails faster once a tx is already broadcasted and processed by the cluster.
			// This reduces the number of "failed" txs due to repeated broadcast attempts.
			rpc.TransactionOpts{PreflightCommitment: rpc.CommitmentProcessed},
		)
		if err != nil {
			if outbound.FallbackTx != nil && shouldUseFallbackTx(err, signer.GetGatewayAddress()) {
				tx = outbound.FallbackTx
			}
			logger.Warn().Err(err).Fields(lf).Msgf("SendTransactionWithOpts failed")
			backOff *= 2
			continue
		}
		logger.Info().Fields(lf).Msgf("broadcasted Solana outbound successfully")

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

	// Create and sign fallback transaction
	fallbackInst, err := signer.createIncrementNonceInstruction(*msgIn)
	if err != nil {
		return nil, errors.Wrap(err, "error creating increment nonce instruction")
	}

	fallbackTx, err := signer.signTx(ctx, fallbackInst, 0)
	if err != nil {
		return nil, errors.Wrap(err, "error signing fallback instruction")
	}

	return &Outbound{
		Tx:         tx,
		FallbackTx: fallbackTx,
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
		signer.Logger().Std.Error().Err(err).Msgf("cannot parse gateway address: %s", address)
		return
	}

	// noop
	if signer.gatewayID.Equals(gatewayID) {
		return
	}

	signer.Logger().Std.Info().
		Str("signer.old_gateway_address", signer.gatewayID.String()).
		Str("signer.new_gateway_address", gatewayID.String()).
		Msg("Updated gateway address")

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
		signer.Logger().Std.Error().Err(err).Msg("GetBalance error")
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
