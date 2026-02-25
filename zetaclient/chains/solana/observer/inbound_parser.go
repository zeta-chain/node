package observer

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/coin"
	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

// InboundEventParser handles parsing of Solana instructions into inbound events
type InboundEventParser struct {
	gatewayID     solana.PublicKey
	senderChainID int64
	logger        zerolog.Logger
	tx            *solana.Transaction
	txResult      *rpc.GetTransactionResult

	seenDeposit    bool
	seenDepositSPL bool
	seenCall       bool
	events         []*clienttypes.InboundEvent
	eventIndex     uint32
}

// NewInboundEventParser creates a new InboundEventParser
// resolvedTx is an optional pre-resolved transaction (e.g., with address lookup tables resolved).
// If provided, it will be used instead of extracting a fresh transaction from txResult.
func NewInboundEventParser(
	txResult *rpc.GetTransactionResult,
	gatewayID solana.PublicKey,
	senderChainID int64,
	logger zerolog.Logger,
	resolvedTx *solana.Transaction,
) (*InboundEventParser, error) {
	var tx *solana.Transaction
	var err error

	if resolvedTx != nil {
		// Use the pre-resolved transaction
		tx = resolvedTx
	} else {
		// Extract transaction from txResult
		tx, err = txResult.Transaction.GetTransaction()
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling transaction")
		}
	}

	return &InboundEventParser{
		gatewayID:     gatewayID,
		senderChainID: senderChainID,
		logger:        logger,
		tx:            tx,
		txResult:      txResult,
		events:        make([]*clienttypes.InboundEvent, 0),
		eventIndex:    0,
	}, nil
}

// parseInstruction parses a single instruction and adds any detected events
func (p *InboundEventParser) parseInstruction(instruction solana.CompiledInstruction, location string) error {
	// get the program ID
	programPk, err := p.tx.Message.Program(instruction.ProgramIDIndex)
	if err != nil {
		p.logger.
			Err(err).
			Stringer("signature", p.tx.Signatures[0]).
			Str("location", location).
			Uint16("index", instruction.ProgramIDIndex).
			Msg("no program found")
		return nil
	}

	// skip instructions that are irrelevant to the gateway program invocation
	if !programPk.Equals(p.gatewayID) {
		return nil
	}

	// try parsing the instruction as a 'deposit' if not seen yet
	if !p.seenDeposit {
		deposit, err := solanacontracts.ParseInboundAsDeposit(p.tx, instruction, p.txResult.Slot)
		if err != nil {
			return errors.Wrap(err, "error ParseInboundAsDeposit")
		}
		if deposit != nil {
			p.seenDeposit = true
			p.events = append(p.events, &clienttypes.InboundEvent{
				SenderChainID:    p.senderChainID,
				Sender:           deposit.Sender,
				Receiver:         deposit.Receiver,
				TxOrigin:         deposit.Sender,
				Amount:           deposit.Amount,
				Memo:             deposit.Memo,
				BlockNumber:      deposit.Slot,
				TxHash:           p.tx.Signatures[0].String(),
				Index:            p.eventIndex,
				CoinType:         coin.CoinType_Gas,
				Asset:            deposit.Asset,
				IsCrossChainCall: deposit.IsCrossChainCall,
				RevertOptions:    deposit.RevertOptions,
			})
			p.eventIndex++
			p.logger.Info().
				Stringer("signature", p.tx.Signatures[0]).
				Str("location", location).
				Uint32("event_index", p.eventIndex-1).
				Msg("deposit detected")
			return nil
		}
	} else {
		p.logger.Warn().
			Stringer("signature", p.tx.Signatures[0]).
			Str("location", location).
			Msg("multiple deposits detected")
	}

	// try parsing the instruction as a 'deposit_spl_token' if not seen yet
	if !p.seenDepositSPL {
		depositSPL, err := solanacontracts.ParseInboundAsDepositSPL(p.tx, instruction, p.txResult.Slot)
		if err != nil {
			return errors.Wrap(err, "error ParseInboundAsDepositSPL")
		}
		if depositSPL != nil {
			p.seenDepositSPL = true
			p.events = append(p.events, &clienttypes.InboundEvent{
				SenderChainID:    p.senderChainID,
				Sender:           depositSPL.Sender,
				Receiver:         depositSPL.Receiver,
				TxOrigin:         depositSPL.Sender,
				Amount:           depositSPL.Amount,
				Memo:             depositSPL.Memo,
				BlockNumber:      depositSPL.Slot,
				TxHash:           p.tx.Signatures[0].String(),
				Index:            p.eventIndex,
				CoinType:         coin.CoinType_ERC20,
				Asset:            depositSPL.Asset,
				IsCrossChainCall: depositSPL.IsCrossChainCall,
				RevertOptions:    depositSPL.RevertOptions,
			})
			p.eventIndex++
			p.logger.Info().
				Stringer("signature", p.tx.Signatures[0]).
				Str("location", location).
				Uint32("eventIndex", p.eventIndex-1).
				Msg("SPL deposit detected")
			return nil
		}
	} else {
		p.logger.Warn().
			Stringer("signature", p.tx.Signatures[0]).
			Str("location", location).
			Msg("multiple SPL deposits detected")
	}

	// try parsing the instruction as a 'call' if not seen yet
	if !p.seenCall {
		call, err := solanacontracts.ParseInboundAsCall(p.tx, instruction, p.txResult.Slot)
		if err != nil {
			return errors.Wrap(err, "error ParseInboundAsCall")
		}
		if call != nil {
			p.seenCall = true
			p.events = append(p.events, &clienttypes.InboundEvent{
				SenderChainID:    p.senderChainID,
				Sender:           call.Sender,
				Receiver:         call.Receiver,
				TxOrigin:         call.Sender,
				Amount:           call.Amount,
				Memo:             call.Memo,
				BlockNumber:      call.Slot,
				TxHash:           p.tx.Signatures[0].String(),
				Index:            p.eventIndex,
				CoinType:         coin.CoinType_NoAssetCall,
				Asset:            call.Asset,
				IsCrossChainCall: call.IsCrossChainCall,
				RevertOptions:    call.RevertOptions,
			})
			p.eventIndex++
			p.logger.Info().
				Stringer("signature", p.tx.Signatures[0]).
				Str("location", location).
				Uint32("eventIndex", p.eventIndex-1).
				Msg("call detected")
			return nil
		}
	} else {
		p.logger.Warn().
			Stringer("signature", p.tx.Signatures[0]).
			Str("location", location).
			Msg("multiple calls detected")
	}

	return nil
}

// Parse parses all instructions in the transaction and returns the detected events
func (p *InboundEventParser) Parse() error {
	// there should be at least one instruction and one account, otherwise skip
	if len(p.tx.Message.Instructions) == 0 {
		return nil
	}

	// parse top-level instructions
	for i, instruction := range p.tx.Message.Instructions {
		if err := p.parseInstruction(instruction, fmt.Sprintf("top-level instruction %d", i)); err != nil {
			return err
		}
	}

	// parse inner instructions
	if p.txResult.Meta != nil && p.txResult.Meta.InnerInstructions != nil {
		for _, inner := range p.txResult.Meta.InnerInstructions {
			for j, instruction := range inner.Instructions {
				compiledInstruction := solana.CompiledInstruction{
					ProgramIDIndex: instruction.ProgramIDIndex,
					Accounts:       instruction.Accounts,
					Data:           instruction.Data,
				}
				desc := fmt.Sprintf("inner instruction %d (outer %d)", j, inner.Index)
				if err := p.parseInstruction(compiledInstruction, desc); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// GetEvents returns the parsed events
func (p *InboundEventParser) GetEvents() []*clienttypes.InboundEvent {
	return p.events
}
