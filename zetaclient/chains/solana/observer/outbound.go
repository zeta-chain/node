package observer

import (
	"context"
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/compliance"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// WatchOutbound watches solana chain for outgoing txs status
// TODO(revamp): move ticker function to ticker file
func (ob *Observer) WatchOutbound(ctx context.Context) error {
	// get app context
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	// create outbound ticker based on chain params
	chainID := ob.Chain().ChainId
	ticker, err := clienttypes.NewDynamicTicker(
		fmt.Sprintf("Solana_WatchOutbound_%d", chainID),
		ob.GetChainParams().OutboundTicker,
	)
	if err != nil {
		ob.Logger().Outbound.Error().Err(err).Msg("error creating ticker")
		return err
	}

	ob.Logger().Outbound.Info().Msgf("WatchOutbound started for chain %d", chainID)
	sampledLogger := ob.Logger().Outbound.Sample(&zerolog.BasicSampler{N: 10})
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !app.IsOutboundObservationEnabled() {
				sampledLogger.Info().Msgf("WatchOutbound: outbound observation is disabled for chain %d", chainID)
				continue
			}

			// process outbound trackers
			err := ob.ProcessOutboundTrackers(ctx)
			if err != nil {
				ob.Logger().
					Outbound.Error().
					Err(err).
					Msgf("WatchOutbound: error ProcessOutboundTrackers for chain %d", chainID)
			}

			ticker.UpdateInterval(ob.GetChainParams().OutboundTicker, ob.Logger().Outbound)
		case <-ob.StopChannel():
			ob.Logger().Outbound.Info().Msgf("WatchOutbound: watcher stopped for chain %d", chainID)
			return nil
		}
	}
}

// ProcessOutboundTrackers processes Solana outbound trackers
func (ob *Observer) ProcessOutboundTrackers(ctx context.Context) error {
	chainID := ob.Chain().ChainId
	trackers, err := ob.ZetacoreClient().GetAllOutboundTrackerByChain(ctx, chainID, interfaces.Ascending)
	if err != nil {
		return errors.Wrap(err, "GetAllOutboundTrackerByChain error")
	}

	// prepare logger fields
	logger := ob.Logger().Outbound.With().
		Str("method", "ProcessOutboundTrackers").
		Int64("chain", chainID).
		Logger()

	// process outbound trackers
	for _, tracker := range trackers {
		// go to next tracker if this one already has a finalized tx
		nonce := tracker.Nonce
		if ob.IsTxFinalized(tracker.Nonce) {
			continue
		}

		// get original cctx parameters
		cctx, err := ob.ZetacoreClient().GetCctxByNonce(ctx, chainID, tracker.Nonce)
		if err != nil {
			// take a rest if zetacore RPC breaks
			return errors.Wrapf(err, "GetCctxByNonce error for chain %d nonce %d", chainID, tracker.Nonce)
		}
		coinType := cctx.InboundParams.CoinType

		// check each txHash and save its txResult if it's finalized and legit
		txCount := 0
		var txResult *rpc.GetTransactionResult
		for _, txHash := range tracker.HashList {
			if result, ok := ob.CheckFinalizedTx(ctx, txHash.TxHash, nonce, coinType); ok {
				txCount++
				txResult = result
				logger.Info().Msgf("confirmed outbound %s for chain %d nonce %d", txHash.TxHash, chainID, nonce)
				if txCount > 1 {
					logger.Error().
						Msgf("checkFinalizedTx passed, txCount %d chain %d nonce %d txResult %v", txCount, chainID, nonce, txResult)
				}
			}
		}
		// should be only one finalized txHash for each nonce
		if txCount == 1 {
			ob.SetTxResult(nonce, txResult)
		} else if txCount > 1 {
			// should not happen. We can't tell which txHash is true. It might happen (e.g. bug, glitchy/hacked endpoint)
			ob.Logger().Outbound.Error().Msgf("finalized multiple (%d) outbound for chain %d nonce %d", txCount, chainID, nonce)
		}
	}

	return nil
}

// VoteOutboundIfConfirmed checks outbound status and returns (continueKeysign, error)
func (ob *Observer) VoteOutboundIfConfirmed(ctx context.Context, cctx *crosschaintypes.CrossChainTx) (bool, error) {
	// get outbound params
	params := cctx.GetCurrentOutboundParam()
	nonce := params.TssNonce
	coinType := cctx.InboundParams.CoinType

	// early return if outbound is not finalized yet
	txResult := ob.GetTxResult(nonce)
	if txResult == nil {
		return true, nil
	}

	// extract tx signature from tx result
	tx, err := txResult.Transaction.GetTransaction()
	if err != nil {
		// should not happen
		return false, errors.Wrapf(err, "GetTransaction error for nonce %d", nonce)
	}
	txSig := tx.Signatures[0]

	// parse gateway instruction from tx result
	inst, err := ParseGatewayInstruction(txResult, ob.gatewayID, coinType)
	if err != nil {
		// should never happen as it was already successfully parsed in CheckFinalizedTx
		return false, errors.Wrapf(err, "ParseGatewayInstruction error for sig %s", txSig)
	}

	// the amount and status of the outbound
	outboundAmount := new(big.Int).SetUint64(inst.TokenAmount())
	// status was already verified as successful in CheckFinalizedTx
	outboundStatus := chains.ReceiveStatus_success

	// compliance check, special handling the cancelled cctx
	if compliance.IsCctxRestricted(cctx) {
		// use cctx's amount to bypass the amount check in zetacore
		outboundAmount = cctx.GetCurrentOutboundParam().Amount.BigInt()
	}

	// post vote to zetacore
	ob.PostVoteOutbound(ctx, cctx.Index, txSig.String(), txResult, outboundAmount, outboundStatus, nonce, coinType)
	return false, nil
}

// PostVoteOutbound posts vote to zetacore for the finalized outbound
func (ob *Observer) PostVoteOutbound(
	ctx context.Context,
	cctxIndex string,
	outboundHash string,
	txResult *rpc.GetTransactionResult,
	valueReceived *big.Int,
	status chains.ReceiveStatus,
	nonce uint64,
	coinType coin.CoinType,
) {
	// create outbound vote message
	msg := ob.CreateMsgVoteOutbound(cctxIndex, outboundHash, txResult, valueReceived, status, nonce, coinType)

	// prepare logger fields
	logFields := map[string]any{
		"chain": ob.Chain().ChainId,
		"nonce": nonce,
		"tx":    outboundHash,
	}

	// so we set retryGasLimit to 0 because the solana gateway withdrawal will always succeed
	// and the vote msg won't trigger ZEVM interaction
	const (
		gasLimit      = zetacore.PostVoteOutboundGasLimit
		retryGasLimit = 0
	)

	// post vote to zetacore
	zetaTxHash, ballot, err := ob.ZetacoreClient().PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)
	if err != nil {
		ob.Logger().Outbound.Error().Err(err).Fields(logFields).Msg("PostVoteOutbound: error posting outbound vote")
		return
	}

	// print vote tx hash and ballot
	if zetaTxHash != "" {
		logFields["vote"] = zetaTxHash
		logFields["ballot"] = ballot
		ob.Logger().Outbound.Info().Fields(logFields).Msg("PostVoteOutbound: posted outbound vote successfully")
	}
}

// CreateMsgVoteOutbound creates a vote outbound message for Solana chain
func (ob *Observer) CreateMsgVoteOutbound(
	cctxIndex string,
	outboundHash string,
	txResult *rpc.GetTransactionResult,
	valueReceived *big.Int,
	status chains.ReceiveStatus,
	nonce uint64,
	coinType coin.CoinType,
) *crosschaintypes.MsgVoteOutbound {
	const (
		// Solana implements a different gas fee model than Ethereum, below values are not used.
		// Solana tx fee is based on both static fee and dynamic fee (priority fee), setting
		// zero values to by pass incorrectly funded gas stability pool.
		outboundGasUsed  = 0
		outboundGasPrice = 0
		outboundGasLimit = 0
	)

	creator := ob.ZetacoreClient().GetKeys().GetOperatorAddress()

	return crosschaintypes.NewMsgVoteOutbound(
		creator.String(),
		cctxIndex,
		outboundHash,
		txResult.Slot, // instead of using block, Solana explorer uses slot for indexing
		outboundGasUsed,
		math.NewInt(outboundGasPrice),
		outboundGasLimit,
		math.NewUintFromBigInt(valueReceived),
		status,
		ob.Chain().ChainId,
		nonce,
		coinType,
	)
}

// CheckFinalizedTx checks if a txHash is finalized for given nonce and coinType
// returns (tx result, true) if finalized or (nil, false) otherwise
func (ob *Observer) CheckFinalizedTx(
	ctx context.Context,
	txHash string,
	nonce uint64,
	coinType coin.CoinType,
) (*rpc.GetTransactionResult, bool) {
	// prepare logger fields
	chainID := ob.Chain().ChainId
	logger := ob.Logger().Outbound.With().
		Str(logs.FieldMethod, "CheckFinalizedTx").
		Int64(logs.FieldChain, chainID).
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, txHash).Logger()

	// convert txHash to signature
	sig, err := solana.SignatureFromBase58(txHash)
	if err != nil {
		logger.Error().Err(err).Msg("SignatureFromBase58 error")
		return nil, false
	}

	// query transaction using "finalized" commitment to avoid re-org
	txResult, err := ob.solClient.GetTransaction(ctx, sig, &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		logger.Error().Err(err).Msg("GetTransaction error")
		return nil, false
	}

	// the tx must be successful in order to effectively increment the nonce
	if txResult.Meta.Err != nil {
		logger.Error().Any("Err", txResult.Meta.Err).Msg("tx is not successful")
		return nil, false
	}

	// parse gateway instruction from tx result
	inst, err := ParseGatewayInstruction(txResult, ob.gatewayID, coinType)
	if err != nil {
		logger.Error().Err(err).Msg("ParseGatewayInstruction error")
		return nil, false
	}
	txNonce := inst.GatewayNonce()

	// recover ECDSA signer from instruction
	signerECDSA, err := inst.Signer()
	if err != nil {
		logger.Error().Err(err).Msg("cannot get instruction signer")
		return nil, false
	}

	// check tx authorization
	if signerECDSA != ob.TSS().EVMAddress() {
		logger.Error().Msgf("tx signer %s is not matching current TSS address %s", signerECDSA, ob.TSS().EVMAddress())
		return nil, false
	}

	// check tx nonce
	if txNonce != nonce {
		logger.Error().Msgf("tx nonce %d is not matching tracker nonce", txNonce)
		return nil, false
	}

	return txResult, true
}

// ParseGatewayInstruction parses the outbound instruction from tx result
func ParseGatewayInstruction(
	txResult *rpc.GetTransactionResult,
	gatewayID solana.PublicKey,
	coinType coin.CoinType,
) (contracts.OutboundInstruction, error) {
	// unmarshal transaction
	tx, err := txResult.Transaction.GetTransaction()
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling transaction")
	}

	// there should be only one single instruction ('withdraw' or 'withdraw_spl_token')
	if len(tx.Message.Instructions) != 1 {
		return nil, fmt.Errorf("want 1 instruction, got %d", len(tx.Message.Instructions))
	}
	instruction := tx.Message.Instructions[0]

	// get the program ID
	programID, err := tx.Message.Program(instruction.ProgramIDIndex)
	if err != nil {
		return nil, errors.Wrap(err, "error getting program ID")
	}

	// the instruction should be an invocation of the gateway program
	if !programID.Equals(gatewayID) {
		return nil, fmt.Errorf("programID %s is not matching gatewayID %s", programID, gatewayID)
	}

	// parse the instruction as a 'withdraw' or 'withdraw_spl_token'
	switch coinType {
	case coin.CoinType_Gas:
		return contracts.ParseInstructionWithdraw(instruction)
	default:
		return nil, fmt.Errorf("unsupported outbound coin type %s", coinType)
	}
}
