package signer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/compliance"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/outboundprocessor"
)

const (
	// solanaTransactionTimeout is the timeout for waiting for an outbound to be confirmed.
	// Transaction referencing a blockhash older than 150 blocks (~1 minute) will expire and be rejected by Solana.
	solanaTransactionTimeout = 1 * time.Minute

	// broadcastBackoff is the initial backoff duration for retrying broadcast
	broadcastBackoff = 1 * time.Second

	// broadcastRetries is the maximum number of retries for broadcasting a transaction
	// 6 retries will span over 1 + 2 + 4 + 8 + 16 + 32 = 63 seconds, good enough for the 1 minute timeout
	broadcastRetries = 6
)

var _ interfaces.ChainSigner = (*Signer)(nil)

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

// NewSigner creates a new Solana signer
func NewSigner(
	chain chains.Chain,
	chainParams observertypes.ChainParams,
	solClient interfaces.SolanaRPCClient,
	tss interfaces.TSSSigner,
	relayerKey *keys.RelayerKey,
	logger base.Logger,
) (*Signer, error) {
	// create base signer
	baseSigner := base.NewSigner(chain, tss, logger)

	// parse gateway ID and PDA
	gatewayID, pda, err := contracts.ParseGatewayWithPDA(chainParams.GatewayAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse gateway address %s", chainParams.GatewayAddress)
	}

	// create Solana signer
	signer := &Signer{
		Signer:    baseSigner,
		client:    solClient,
		gatewayID: gatewayID,
		pda:       pda,
	}

	// construct Solana private key if present
	if relayerKey != nil {
		privKey, err := solana.PrivateKeyFromBase58(relayerKey.PrivateKey)
		if err != nil {
			return nil, errors.Wrap(err, "unable to construct solana private key")
		}
		signer.relayerKey = &privKey
		logger.Std.Info().Msgf("Solana relayer address: %s", signer.relayerKey.PublicKey())
	} else {
		logger.Std.Info().Msg("Solana relayer key is not provided")
	}

	return signer, nil
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
	outboundProc *outboundprocessor.Processor,
	outboundID string,
	_ interfaces.ChainObserver,
	zetacoreClient interfaces.ZetacoreClient,
	height uint64,
) {
	// end outbound process on panic
	defer func() {
		outboundProc.EndTryProcess(outboundID)
		if err := recover(); err != nil {
			signer.Logger().Std.Error().Msgf("TryProcessOutbound: %s, caught panic error: %v", cctx.Index, err)
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

	var tx *solana.Transaction

	switch coinType {
	case coin.CoinType_Cmd:
		whitelistTx, err := signer.prepareWhitelistTx(ctx, cctx, height)
		if err != nil {
			logger.Error().Err(err).Msgf("TryProcessOutbound: Fail to sign whitelist outbound")
			return
		}

		tx = whitelistTx

	case coin.CoinType_Gas:
		withdrawTx, err := signer.prepareWithdrawTx(ctx, cctx, height, logger)
		if err != nil {
			logger.Error().Err(err).Msgf("TryProcessOutbound: Fail to sign withdraw outbound")
			return
		}

		tx = withdrawTx

	case coin.CoinType_ERC20:
		withdrawSPLTx, err := signer.prepareWithdrawSPLTx(ctx, cctx, height, logger)
		if err != nil {
			logger.Error().Err(err).Msgf("TryProcessOutbound: Fail to sign withdraw spl outbound")
			return
		}

		tx = withdrawSPLTx
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

	// broadcast the signed tx to the Solana network
	signer.broadcastOutbound(ctx, tx, chainID, nonce, logger, zetacoreClient)
}

// broadcastOutbound sends the signed transaction to the Solana network
func (signer *Signer) broadcastOutbound(
	ctx context.Context,
	tx *solana.Transaction,
	chainID int64,
	nonce uint64,
	logger zerolog.Logger,
	zetacoreClient interfaces.ZetacoreClient,
) {
	// prepare logger fields
	lf := map[string]any{
		logs.FieldMethod: "broadcastOutbound",
		logs.FieldNonce:  nonce,
		logs.FieldTx:     tx.Signatures[0].String(),
	}

	// try broacasting tx with increasing backoff (1s, 2s, 4s, 8s, 16s, 32s)
	// to tolerate tx nonce mismatch with PDA nonce or unknown RPC error
	backOff := broadcastBackoff
	for i := 0; i < broadcastRetries; i++ {
		time.Sleep(backOff)

		// PDA nonce may already be increased by other relayer, no need to retry
		pdaInfo, err := signer.client.GetAccountInfo(ctx, signer.pda)
		if err != nil {
			logger.Error().Err(err).Fields(lf).Msgf("unable to get PDA account info")
		} else {
			pda := contracts.PdaInfo{}
			err = borsh.Deserialize(&pda, pdaInfo.Bytes())
			if err != nil {
				logger.Error().Err(err).Msgf("unable to deserialize PDA info")
			} else if pda.Nonce > nonce {
				logger.Info().Err(err).Msgf("PDA nonce %d is greater than outbound nonce, stop retrying", pda.Nonce)
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
			logger.Warn().Err(err).Fields(lf).Msgf("SendTransactionWithOpts failed")
			backOff *= 2
			continue
		}

		// successful broadcast; report to the outbound tracker
		signer.reportToOutboundTracker(ctx, zetacoreClient, chainID, nonce, txSig, logger)
		break
	}
}

func (signer *Signer) prepareWithdrawTx(
	ctx context.Context,
	cctx *types.CrossChainTx,
	height uint64,
	logger zerolog.Logger,
) (*solana.Transaction, error) {
	params := cctx.GetCurrentOutboundParam()
	// compliance check
	cancelTx := compliance.IsCctxRestricted(cctx)
	if cancelTx {
		compliance.PrintComplianceLog(
			logger,
			signer.Logger().Compliance,
			true,
			signer.Chain().ChainId,
			cctx.Index,
			cctx.InboundParams.Sender,
			params.Receiver,
			"SOL",
		)
	}

	// sign gateway withdraw message by TSS
	msg, err := signer.createAndSignMsgWithdraw(ctx, params, height, cancelTx)
	if err != nil {
		return nil, err
	}

	// sign the withdraw transaction by relayer key
	tx, err := signer.signWithdrawTx(ctx, *msg)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (signer *Signer) prepareWithdrawSPLTx(
	ctx context.Context,
	cctx *types.CrossChainTx,
	height uint64,
	logger zerolog.Logger,
) (*solana.Transaction, error) {
	params := cctx.GetCurrentOutboundParam()
	// compliance check
	cancelTx := compliance.IsCctxRestricted(cctx)
	if cancelTx {
		compliance.PrintComplianceLog(
			logger,
			signer.Logger().Compliance,
			true,
			signer.Chain().ChainId,
			cctx.Index,
			cctx.InboundParams.Sender,
			params.Receiver,
			"SPL",
		)
	}

	// get mint details to get decimals
	mint, err := signer.decodeMintAccountDetails(ctx, cctx.InboundParams.Asset)
	if err != nil {
		return nil, err
	}

	// sign gateway withdraw spl message by TSS
	msg, err := signer.createAndSignMsgWithdrawSPL(
		ctx,
		params,
		height,
		cctx.InboundParams.Asset,
		mint.Decimals,
		cancelTx,
	)
	if err != nil {
		return nil, err
	}

	// sign the withdraw transaction by relayer key
	tx, err := signer.signWithdrawSPLTx(ctx, *msg)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (signer *Signer) prepareWhitelistTx(
	ctx context.Context,
	cctx *types.CrossChainTx,
	height uint64,
) (*solana.Transaction, error) {
	params := cctx.GetCurrentOutboundParam()
	relayedMsg := strings.Split(cctx.RelayedMessage, ":")
	if len(relayedMsg) != 2 {
		return nil, fmt.Errorf("TryProcessOutbound: invalid relayed msg")
	}

	pk, err := solana.PublicKeyFromBase58(relayedMsg[1])
	if err != nil {
		return nil, err
	}

	seed := [][]byte{[]byte("whitelist"), pk.Bytes()}
	whitelistEntryPDA, _, err := solana.FindProgramAddress(seed, signer.gatewayID)
	if err != nil {
		return nil, err
	}

	// sign gateway whitelist message by TSS
	msg, err := signer.createAndSignMsgWhitelist(ctx, params, height, pk, whitelistEntryPDA)
	if err != nil {
		return nil, err
	}

	// sign the whitelist transaction by relayer key
	tx, err := signer.signWhitelistTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (signer *Signer) decodeMintAccountDetails(ctx context.Context, asset string) (token.Mint, error) {
	info, err := signer.client.GetAccountInfo(ctx, solana.MustPublicKeyFromBase58(asset))
	if err != nil {
		return token.Mint{}, err
	}

	var mint token.Mint
	// Account{}.Data.GetBinary() returns the *decoded* binary data
	// regardless the original encoding (it can handle them all).
	err = bin.NewBinDecoder(info.Value.Data.GetBinary()).Decode(&mint)
	if err != nil {
		return token.Mint{}, err
	}

	return mint, nil
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

// TODO: get rid of below four functions for Solana and Bitcoin
// https://github.com/zeta-chain/node/issues/2532
func (signer *Signer) SetZetaConnectorAddress(_ ethcommon.Address) {
}

func (signer *Signer) SetERC20CustodyAddress(_ ethcommon.Address) {
}

func (signer *Signer) GetZetaConnectorAddress() ethcommon.Address {
	return ethcommon.Address{}
}

func (signer *Signer) GetERC20CustodyAddress() ethcommon.Address {
	return ethcommon.Address{}
}
