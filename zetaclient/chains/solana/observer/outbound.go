package observer

import (
	"context"
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	contract "github.com/zeta-chain/zetacore/pkg/contract/solana"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

// GetTxID returns a unique id for Solana outbound
func (ob *Observer) GetTxID(nonce uint64) string {
	tssAddr := ob.TSS().EVMAddress().String()
	return fmt.Sprintf("%d-%s-%d", ob.Chain().ChainId, tssAddr, nonce)
}

// WatchOutbound watches evm chain for outgoing txs status
// TODO(revamp): move ticker function to ticker file
// TODO(revamp): move inner logic to a separate function
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
			if !app.IsOutboundObservationEnabled(ob.GetChainParams()) {
				sampledLogger.Info().Msgf("WatchOutbound: outbound observation is disabled for chain %d", chainID)
				continue
			}
			trackers, err := ob.ZetacoreClient().GetAllOutboundTrackerByChain(ctx, chainID, interfaces.Ascending)
			if err != nil {
				continue
			}

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
					ob.Logger().Outbound.Error().
						Err(err).
						Msgf("WatchOutbound: can't find cctx for chain %d nonce %d", chainID, tracker.Nonce)
					break
				}
				coinType := cctx.InboundParams.CoinType

				txCount := 0
				var txResult *rpc.GetTransactionResult
				for _, txHash := range tracker.HashList {
					if result, ok := ob.CheckFinalizedTx(ctx, txHash.TxHash, nonce, coinType); ok {
						txCount++
						txResult = result
						ob.Logger().Outbound.Info().
							Msgf("WatchOutbound: confirmed outbound %s for chain %d nonce %d", txHash.TxHash, chainID, nonce)
						if txCount > 1 {
							ob.Logger().Outbound.Error().Msgf(
								"WatchOutbound: checkFinalizedTx passed, txCount %d chain %d nonce %d txResult %v", txCount, chainID, nonce, txResult)
						}
					}
				}
				// should be only one finalized txHash for each nonce
				if txCount == 1 {
					ob.SetTxResult(nonce, txResult)
				} else if txCount > 1 {
					// should not happen. We can't tell which txHash is true. It might happen (e.g. glitchy/hacked endpoint)
					ob.Logger().Outbound.Error().Msgf("WatchOutbound: finalized multiple (%d) outbound for chain %d nonce %d", txCount, chainID, nonce)
				}
			}
			ticker.UpdateInterval(ob.GetChainParams().OutboundTicker, ob.Logger().Outbound)
		case <-ob.StopChannel():
			ob.Logger().Outbound.Info().Msgf("WatchOutbound: watcher stopped for chain %d", chainID)
			return nil
		}
	}
}

// IsOutboundProcessed checks outbound status and returns (isIncluded, isFinalized, error)
// It also posts vote to zetacore if the tx is finalized
// TODO(revamp): rename as it also vote the outbound
func (ob *Observer) IsOutboundProcessed(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
) (bool, bool, error) {
	// early return if outbound is not finalized yet
	params := cctx.GetCurrentOutboundParam()
	nonce := params.TssNonce
	txResult := ob.GetTxResult(nonce)
	if txResult == nil {
		return false, false, nil
	}

	// extract tx signature of the finalized tx
	tx, err := txResult.Transaction.GetTransaction()
	if err != nil {
		return true, true, errors.Wrapf(err, "GetTransaction error for nonce %d", nonce)
	}
	txSig := tx.Signatures[0]

	outboundAmount := params.Amount.BigInt()      // FIXME: parse this amount from txRes itself, not from cctx
	outboundStatus := chains.ReceiveStatus_failed // tx was failed/reverted: FIXME: see the note in function comment
	coinType := cctx.InboundParams.CoinType
	if txResult.Meta.Err == nil {
		// tx was successful
		outboundStatus = chains.ReceiveStatus_success
	}

	// post vote to zetacore
	ob.PostVoteOutbound(ctx, cctx.Index, txSig.String(), txResult, outboundAmount, outboundStatus, nonce, coinType)
	return true, true, nil
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
	chainID := ob.Chain().ChainId
	computeUnitsConsumed := uint64(0)
	cuPrice := big.NewInt(0)
	if txResult.Meta.ComputeUnitsConsumed == nil {
		ob.Logger().Outbound.Warn().Msgf("solana.GetTransaction: compute units consumed is nil")
	} else {
		computeUnitsConsumed = *txResult.Meta.ComputeUnitsConsumed
		if computeUnitsConsumed <= 0 {
			ob.Logger().Outbound.Warn().Msgf("solana.GetTransaction: compute units consumed is %d", computeUnitsConsumed)
			computeUnitsConsumed = 5001 // default to 5000, for a single signature tx; make it 5001 to distinguish
		}
		cuPrice.SetUint64(txResult.Meta.Fee / computeUnitsConsumed)
	}

	creator := ob.ZetacoreClient().GetKeys().GetOperatorAddress()

	msg := crosschaintypes.NewMsgVoteOutbound(
		creator.String(),
		cctxIndex,
		outboundHash,
		txResult.Slot, // TODO: check this; is slot equivalent to block height?
		computeUnitsConsumed,
		math.NewIntFromBigInt(cuPrice),
		200_000,                               // this is default compute unit budget;
		math.NewUintFromBigInt(valueReceived), // FIXME: parse this amount from txRes itself, not from cctx
		status,
		chainID,
		nonce, // FIXME: parse this from the txRes/tx ?
		coinType,
	)

	const gasLimit = zetacore.PostVoteOutboundGasLimit
	var retryGasLimit uint64
	if msg.Status == chains.ReceiveStatus_failed {
		retryGasLimit = zetacore.PostVoteOutboundRevertGasLimit
	}

	// post vote to zetacore
	logFields := map[string]any{
		"chain":    chainID,
		"nonce":    nonce,
		"outbound": outboundHash,
	}
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

// CheckFinalizedTx checks if a txHash is finalized for given nonce and coinType
// returns (tx result, true) if finalized or (nil, false) otherwise
func (ob *Observer) CheckFinalizedTx(
	ctx context.Context,
	txHash string,
	nonce uint64,
	coinType coin.CoinType,
) (*rpc.GetTransactionResult, bool) {
	// prepare logger
	logger := ob.Logger().Outbound.With().
		Str("method", "checkFinalizedTx").
		Int64("chain", ob.Chain().ChainId).
		Uint64("nonce", nonce).
		Str("tx", txHash).Logger()

	// convert txHash to signature
	sig, err := solana.SignatureFromBase58(txHash)
	if err != nil {
		logger.Error().Err(err).Msg("SignatureFromBase58 err")
		return nil, false
	}

	// query transaction using "finalized" commitment to avoid re-org
	txResult, err := ob.solClient.GetTransaction(ctx, sig, &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		logger.Error().Err(err).Msg("GetTransaction err")
		return nil, false
	}

	// the tx must be successful in order to effectively increment the nonce
	if txResult.Meta.Err != nil {
		logger.Error().Msg("tx is not successful")
		return nil, false
	}

	// parse gateway instruction from tx result
	inst, err := ob.ParseGatewayInstruction(txResult, coinType)
	if err != nil {
		logger.Error().Err(err).Msg("ParseGatewayInstruction err")
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
		logger.Error().Msgf("tx signer %s is not matching TSS address %s", signerECDSA, ob.TSS().EVMAddress())
		return nil, false
	}

	// check tx nonce
	if txNonce != nonce {
		logger.Error().Msgf("tx nonce %d is not matching cctx nonce %d", txNonce, nonce)
		return nil, false
	}

	return txResult, true
}

// ParseGatewayInstruction parses the instruction signer and nonce from tx result
func (ob *Observer) ParseGatewayInstruction(
	txResult *rpc.GetTransactionResult,
	coinType coin.CoinType,
) (contract.OutboundInstruction, error) {
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
	programPk, err := tx.Message.Program(instruction.ProgramIDIndex)
	if err != nil {
		return nil, errors.Wrap(err, "error getting program ID")
	}

	// the instruction should be an invocation of the gateway program
	if !programPk.Equals(ob.gatewayID) {
		return nil, errors.New("not a gateway program invocation")
	}

	// parse the instruction as a 'withdraw' or 'withdraw_spl_token'
	switch coinType {
	case coin.CoinType_Gas:
		return ob.ParseInstructionWithdraw(tx, 0)
	default:
		return nil, errors.New("unsupported outbound coin type")
	}
}

// ParseInstructionWithdraw tries to parse an instruction as a 'withdraw'.
// It returns nil if the instruction can't be parsed as a 'withdraw'.
func (ob *Observer) ParseInstructionWithdraw(
	tx *solana.Transaction,
	instructionIndex int,
) (*contract.WithdrawInstructionParams, error) {
	// locate instruction by index
	instruction := tx.Message.Instructions[instructionIndex]

	// try deserializing instruction as a 'withdraw'
	inst := &contract.WithdrawInstructionParams{}
	err := borsh.Deserialize(inst, instruction.Data)
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing instruction")
	}

	// check the discriminator to ensure it's a 'withdraw' instruction
	if inst.Discriminator != contract.DiscriminatorWithdraw() {
		return nil, errors.New("not a withdraw instruction")
	}

	return inst, nil
}
