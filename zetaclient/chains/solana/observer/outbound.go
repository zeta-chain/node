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

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	contracts "github.com/zeta-chain/zetacore/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

// WatchOutbound watches solana chain for outgoing txs status
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
				ob.Logger().
					Outbound.Error().
					Err(err).
					Msgf("WatchOutbound: GetAllOutboundTrackerByChain error for chain %d", chainID)
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
func (ob *Observer) IsOutboundProcessed(ctx context.Context, cctx *crosschaintypes.CrossChainTx) (bool, bool, error) {
	// get outbound params
	params := cctx.GetCurrentOutboundParam()
	nonce := params.TssNonce
	coinType := cctx.InboundParams.CoinType

	// early return if outbound is not finalized yet
	txResult := ob.GetTxResult(nonce)
	if txResult == nil {
		return false, false, nil
	}

	// extract tx signature from tx result
	tx, err := txResult.Transaction.GetTransaction()
	if err != nil {
		// should not happen
		return true, true, errors.Wrapf(err, "GetTransaction error for nonce %d", nonce)
	}
	txSig := tx.Signatures[0]

	// parse gateway instruction from tx result
	inst, err := ParseGatewayInstruction(txResult, ob.gatewayID, coinType)
	if err != nil {
		// should never happen as it was already successfully parsed in CheckFinalizedTx
		return true, true, errors.Wrapf(err, "ParseGatewayInstruction error for sig %s", txSig)
	}

	// the amount and status of the outbound
	outboundAmount := new(big.Int).SetUint64(inst.TokenAmount())
	// status was already verified as successful in CheckFinalizedTx
	outboundStatus := chains.ReceiveStatus_success

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

	const (
		// Solana implements a different gas fee model than Ethereum, below values are not used.
		// Solana tx fee is based on both static fee and dynamic fee (priority fee), setting
		// zero values to by pass incorrectly funded gas stability pool.
		outboundGasUsed  = 0
		outboundGasPrice = 0
		outboundGasLimit = 0

		gasLimit      = zetacore.PostVoteOutboundGasLimit
		retryGasLimit = 0
	)

	creator := ob.ZetacoreClient().GetKeys().GetOperatorAddress()
	msg := crosschaintypes.NewMsgVoteOutbound(
		creator.String(),
		cctxIndex,
		outboundHash,
		txResult.Slot, // instead of using block, Solana explorer uses slot for indexing
		outboundGasUsed,
		math.NewInt(outboundGasPrice),
		outboundGasLimit,
		math.NewUintFromBigInt(valueReceived),
		status,
		chainID,
		nonce,
		coinType,
	)

	// post vote to zetacore
	logFields := map[string]any{
		"chain": chainID,
		"nonce": nonce,
		"tx":    outboundHash,
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
	// prepare logger fields
	chainID := ob.Chain().ChainId
	logger := ob.Logger().Outbound.With().
		Str("method", "checkFinalizedTx").
		Int64("chain", chainID).
		Uint64("nonce", nonce).
		Str("tx", txHash).Logger()

	// convert txHash to signature
	sig, err := solana.SignatureFromBase58(txHash)
	if err != nil {
		logger.Error().Err(err).Msgf("SignatureFromBase58 err for chain %d nonce %d", chainID, nonce)
		return nil, false
	}

	// query transaction using "finalized" commitment to avoid re-org
	txResult, err := ob.solClient.GetTransaction(ctx, sig, &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		logger.Error().Err(err).Msgf("GetTransaction err for chain %d nonce %d", chainID, nonce)
		return nil, false
	}

	// the tx must be successful in order to effectively increment the nonce
	if txResult.Meta.Err != nil {
		logger.Error().Any("Err", txResult.Meta.Err).Msgf("tx is not successful for chain %d nonce %d", chainID, nonce)
		return nil, false
	}

	// parse gateway instruction from tx result
	inst, err := ParseGatewayInstruction(txResult, ob.gatewayID, coinType)
	if err != nil {
		logger.Error().Err(err).Msgf("ParseGatewayInstruction err for chain %d nonce %d", chainID, nonce)
		return nil, false
	}
	txNonce := inst.GatewayNonce()

	// recover ECDSA signer from instruction
	signerECDSA, err := inst.Signer()
	if err != nil {
		logger.Error().Err(err).Msgf("cannot get instruction signer for chain %d nonce %d", chainID, nonce)
		return nil, false
	}

	// check tx authorization
	if signerECDSA != ob.TSS().EVMAddress() {
		logger.Error().Msgf("tx signer %s is not matching TSS, chain %d nonce %d", signerECDSA, chainID, nonce)
		return nil, false
	}

	// check tx nonce
	if txNonce != nonce {
		logger.Error().Msgf("tx nonce %d is not matching cctx, chain %d nonce %d", txNonce, chainID, nonce)
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
	programPk, err := tx.Message.Program(instruction.ProgramIDIndex)
	if err != nil {
		return nil, errors.Wrap(err, "error getting program ID")
	}

	// the instruction should be an invocation of the gateway program
	if !programPk.Equals(gatewayID) {
		return nil, errors.New("not a gateway program invocation")
	}

	// parse the instruction as a 'withdraw' or 'withdraw_spl_token'
	switch coinType {
	case coin.CoinType_Gas:
		return contracts.ParseInstructionWithdraw(instruction)
	default:
		return nil, errors.New("unsupported outbound coin type")
	}
}
