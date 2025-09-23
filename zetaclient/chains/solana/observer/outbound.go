package observer

import (
	"context"
	stderrors "errors"
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/compliance"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

var (
	gasOutboundParsers = []func(solana.CompiledInstruction) (contracts.OutboundInstruction, error){
		func(inst solana.CompiledInstruction) (contracts.OutboundInstruction, error) {
			return contracts.ParseInstructionWithdraw(inst)
		},
		func(inst solana.CompiledInstruction) (contracts.OutboundInstruction, error) {
			return contracts.ParseInstructionExecute(inst)
		},
		func(inst solana.CompiledInstruction) (contracts.OutboundInstruction, error) {
			return contracts.ParseInstructionExecuteRevert(inst)
		},
	}

	splOutboundParsers = []func(solana.CompiledInstruction) (contracts.OutboundInstruction, error){
		func(inst solana.CompiledInstruction) (contracts.OutboundInstruction, error) {
			return contracts.ParseInstructionWithdrawSPL(inst)
		},
		func(inst solana.CompiledInstruction) (contracts.OutboundInstruction, error) {
			return contracts.ParseInstructionExecuteSPL(inst)
		},
		func(inst solana.CompiledInstruction) (contracts.OutboundInstruction, error) {
			return contracts.ParseInstructionExecuteSPLRevert(inst)
		},
	}
)

// ProcessOutboundTrackers processes Solana outbound trackers
func (ob *Observer) ProcessOutboundTrackers(ctx context.Context) error {
	chainID := ob.Chain().ChainId
	trackers, err := ob.ZetacoreClient().GetOutboundTrackers(ctx, chainID)
	if err != nil {
		return errors.Wrap(err, "GetOutboundTrackers error")
	}

	logger := ob.Logger().Outbound

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

				logger.Info().
					Str(logs.FieldTx, txHash.TxHash).
					Uint64(logs.FieldNonce, nonce).
					Msg("confirmed outbound")
			}
		}

		// should be only one finalized txHash for each nonce
		if txCount == 1 {
			ob.SetTxResult(nonce, txResult)
		} else if txCount > 1 {
			// Should not happen. We can't tell which txHash is true.
			// It might happen (e.g. bug, glitchy/hacked endpoint)
			logger.Error().
				Uint64(logs.FieldNonce, nonce).
				Int("count", txCount).
				Msg("finalized multiple outbounds")
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
	inst, err := ParseGatewayInstruction(txResult, ob.gatewayID, coinType, nonce)
	if err != nil {
		// should never happen as it was already successfully parsed in CheckFinalizedTx
		return false, errors.Wrapf(err, "ParseGatewayInstruction error for sig %s", txSig)
	}

	// the amount and status of the outbound
	outboundAmount := new(big.Int).SetUint64(inst.TokenAmount())

	// status was already verified as successful in CheckFinalizedTx
	outboundStatus := chains.ReceiveStatus_success
	if inst.InstructionDiscriminator() == contracts.DiscriminatorIncrementNonce {
		outboundStatus = chains.ReceiveStatus_failed
	}

	// cancelled transaction means the outbound is failed
	// - set amount to CCTX's amount to bypass amount check in zetacore
	// - set status to failed to revert the CCTX in zetacore
	if compliance.IsCCTXRestricted(cctx) {
		outboundAmount = cctx.GetCurrentOutboundParam().Amount.BigInt()
		outboundStatus = chains.ReceiveStatus_failed
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
	logger := ob.Logger().Outbound.With().
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, outboundHash).
		Logger()

	// create outbound vote message
	msg := ob.CreateMsgVoteOutbound(cctxIndex, outboundHash, txResult, valueReceived, status, nonce, coinType)

	// so we set retryGasLimit to 0 because the solana gateway withdrawal will always succeed
	// and the vote msg won't trigger ZEVM interaction
	const gasLimit = zetacore.PostVoteOutboundGasLimit

	retryGasLimit := zetacore.PostVoteOutboundRetryGasLimit
	if msg.Status == chains.ReceiveStatus_failed {
		retryGasLimit = zetacore.PostVoteOutboundRevertGasLimit
	}

	// post vote to zetacore
	zetaTxHash, ballot, err := ob.ZetacoreClient().PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)
	if err != nil {
		logger.Error().Err(err).Msg("error posting outbound vote")
		return
	}

	// print vote tx hash and ballot
	if zetaTxHash != "" {
		logger.Info().
			Str(logs.FieldZetaTx, zetaTxHash).
			Str(logs.FieldBallotIndex, ballot).
			Msg("posted outbound vote")
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
		crosschaintypes.ConfirmationMode_SAFE,
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
		Int64(logs.FieldChain, chainID).
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, txHash).Logger()

	// convert txHash to signature
	sig, err := solana.SignatureFromBase58(txHash)
	if err != nil {
		logger.Error().Err(err).Msg("error calling SignatureFromBase58")
		return nil, false
	}

	// query transaction using "finalized" commitment to avoid re-org
	txResult, err := ob.solanaClient.GetTransaction(ctx, sig, &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		logger.Error().Err(err).Msg("error calling GetTransaction")
		return nil, false
	}

	// the tx must be successful in order to effectively increment the nonce
	if txResult.Meta.Err != nil {
		logger.Error().Any("tx_meta_err", txResult.Meta.Err).Msg("tx was not successful")
		return nil, false
	}

	// parse gateway instruction from tx result
	inst, err := ParseGatewayInstruction(txResult, ob.gatewayID, coinType, nonce)
	if err != nil {
		logger.Error().Err(err).Msg("error calling ParseGatewayInstruction")
		return nil, false
	}

	// recover ECDSA signer from instruction
	signerECDSA, err := inst.Signer()
	if err != nil {
		logger.Error().Err(err).Msg("cannot get instruction signer")
		return nil, false
	}

	// check tx authorization
	if signerECDSA != ob.TSS().PubKey().AddressEVM() {
		logger.Error().
			Stringer("signer", signerECDSA).
			Stringer("address", ob.TSS().PubKey().AddressEVM()).
			Msg("tx signer is not matching current TSS address")
		return nil, false
	}

	return txResult, true
}

// parseInstructionWith attempts to parse an instruction using a list of parsers
func parseInstructionWith(
	instruction solana.CompiledInstruction,
	parsers []func(solana.CompiledInstruction) (contracts.OutboundInstruction, error),
) (contracts.OutboundInstruction, error) {
	errs := make([]error, 0, len(parsers))
	for _, parser := range parsers {
		inst, err := parser(instruction)
		if err == nil {
			return inst, nil
		}
		errs = append(errs, err)
	}
	return nil, errors.Wrap(stderrors.Join(errs...), "failed to parse instruction")
}

// ParseGatewayInstruction parses the outbound instruction from tx result
func ParseGatewayInstruction(
	txResult *rpc.GetTransactionResult,
	gatewayID solana.PublicKey,
	coinType coin.CoinType,
	expectedNonce uint64,
) (contracts.OutboundInstruction, error) {
	tx, err := txResult.Transaction.GetTransaction()
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling transaction")
	}

	var (
		candidateInst contracts.OutboundInstruction
		parseErrs     []error
		matchCount    int
	)

	for _, inst := range tx.Message.Instructions {
		programID, err := tx.Message.Program(inst.ProgramIDIndex)
		if err != nil {
			parseErrs = append(parseErrs, fmt.Errorf("error getting program ID: %w", err))
			continue
		}

		// Skip non-Gateway program instructions
		if !programID.Equals(gatewayID) {
			continue
		}

		// try parsing increment_nonce
		if parsed, err := contracts.ParseInstructionIncrementNonce(inst); err == nil {
			if parsed.GatewayNonce() == expectedNonce {
				matchCount++
				candidateInst = parsed
			}
			continue
		}

		// try parsing based on coin type
		var parsed contracts.OutboundInstruction
		switch coinType {
		case coin.CoinType_Gas:
			parsed, err = parseInstructionWith(inst, gasOutboundParsers)
		case coin.CoinType_ERC20:
			parsed, err = parseInstructionWith(inst, splOutboundParsers)
		case coin.CoinType_Cmd:
			parsed, err = contracts.ParseInstructionWhitelist(inst)
		case coin.CoinType_NoAssetCall:
			parsed, err = contracts.ParseInstructionExecute(inst)
		default:
			err = fmt.Errorf("unsupported outbound coin type %s", coinType)
		}

		if err == nil {
			if parsed.GatewayNonce() == expectedNonce {
				matchCount++
				candidateInst = parsed
			}
		} else {
			parseErrs = append(parseErrs, err)
		}
	}

	if matchCount == 0 {
		if len(parseErrs) == 0 {
			return nil, fmt.Errorf("no matching outbound instruction with expected nonce %d", expectedNonce)
		}
		return nil, errors.Wrap(stderrors.Join(parseErrs...), "no matching outbound instruction with expected nonce")
	}

	// should not happen
	if matchCount > 1 {
		return nil, fmt.Errorf("multiple outbounds with same nonce detected")
	}

	return candidateInst, nil
}
