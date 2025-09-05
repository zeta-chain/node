package solana

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

const (
	// MaxSignaturesPerTicker is the maximum number of signatures to process on a ticker
	MaxSignaturesPerTicker = 100
)

// Inbound represents an inbound instruction from a Solana transaction to ZetaChain
type Inbound struct {
	Sender           string
	Receiver         string
	Amount           uint64
	Memo             []byte
	Slot             uint64
	Asset            string
	IsCrossChainCall bool
	RevertOptions    *RevertOptions
}

// ParseInboundAsDeposit tries to parse an instruction as a 'deposit' or 'deposit_and_call'.
// It returns nil if the instruction can't be parsed.
func ParseInboundAsDeposit(
	tx *solana.Transaction,
	instruction solana.CompiledInstruction,
	slot uint64,
) (*Inbound, error) {
	// first try to parse as deposit, then as deposit_and_call
	deposit, err := parseAsDeposit(tx, instruction, slot)
	if err != nil || deposit != nil {
		return deposit, err
	}

	return parseAsDepositAndCall(tx, instruction, slot)
}

// parseAsDeposit tries to parse instruction as deposit
func parseAsDeposit(
	tx *solana.Transaction,
	instruction solana.CompiledInstruction,
	slot uint64,
) (*Inbound, error) {
	// try deserializing instruction as a deposit
	inst := DepositInstructionParams{}
	err := borsh.Deserialize(&inst, instruction.Data)
	if err != nil {
		return nil, nil
	}

	// check if the instruction is a deposit or not, if not, skip parsing
	if inst.Discriminator != DiscriminatorDeposit {
		return nil, nil
	}

	// get the sender address (skip if unable to parse signer address)
	sender, err := getSignerDeposit(tx, &instruction)
	if err != nil {
		return nil, err
	}

	receiver, err := parseReceiver(inst.Receiver)
	if err != nil {
		return nil, err
	}

	return &Inbound{
		Sender:           sender,
		Receiver:         receiver,
		Amount:           inst.Amount,
		Memo:             []byte{},
		Slot:             slot,
		Asset:            "", // no asset for gas token SOL
		IsCrossChainCall: false,
		RevertOptions:    inst.RevertOptions,
	}, nil
}

// parseAsDepositAndCall tries to parse instruction as deposit_and_call
func parseAsDepositAndCall(
	tx *solana.Transaction,
	instruction solana.CompiledInstruction,
	slot uint64,
) (*Inbound, error) {
	// try deserializing instruction as a deposit_and_call
	instDepositAndCall := DepositAndCallInstructionParams{}
	err := borsh.Deserialize(&instDepositAndCall, instruction.Data)
	if err != nil {
		return nil, nil
	}

	// check if the instruction is a deposit_and_call or not, if not, skip parsing
	if instDepositAndCall.Discriminator != DiscriminatorDepositAndCall {
		return nil, nil
	}

	receiver, err := parseReceiver(instDepositAndCall.Receiver)
	if err != nil {
		return nil, err
	}

	// get the sender address (skip if unable to parse signer address)
	sender, err := getSignerDeposit(tx, &instruction)
	if err != nil {
		return nil, err
	}

	return &Inbound{
		Sender:           sender,
		Receiver:         receiver,
		Amount:           instDepositAndCall.Amount,
		Memo:             instDepositAndCall.Memo,
		Slot:             slot,
		Asset:            "", // no asset for gas token SOL
		IsCrossChainCall: true,
		RevertOptions:    instDepositAndCall.RevertOptions,
	}, nil
}

// ParseInboundAsDepositSPL tries to parse an instruction as a deposit_spl or deposit_spl_and_call.
// It returns nil if the instruction can't be parsed as a deposit_spl.
func ParseInboundAsDepositSPL(
	tx *solana.Transaction,
	instruction solana.CompiledInstruction,
	slot uint64,
) (*Inbound, error) {
	// first try to parse as deposit_spl, then as deposit_spl_and_call
	deposit, err := parseAsDepositSPL(tx, instruction, slot)
	if err != nil || deposit != nil {
		return deposit, err
	}

	return parseAsDepositSPLAndCall(tx, instruction, slot)
}

// parseAsDepositSPL tries to parse instruction as deposit_spl
func parseAsDepositSPL(
	tx *solana.Transaction,
	instruction solana.CompiledInstruction,
	slot uint64,
) (*Inbound, error) {
	// try deserializing instruction as a deposit_spl
	var inst DepositSPLInstructionParams
	err := borsh.Deserialize(&inst, instruction.Data)
	if err != nil {
		return nil, nil
	}

	// check if the instruction is a deposit_spl or not, if not, skip parsing
	if inst.Discriminator != DiscriminatorDepositSPL {
		return nil, nil
	}

	// get the sender and spl addresses
	sender, spl, err := getSignerAndSPLFromDepositSPLAccounts(tx, &instruction)
	if err != nil {
		return nil, err
	}

	receiver, err := parseReceiver(inst.Receiver)
	if err != nil {
		return nil, err
	}

	return &Inbound{
		Sender:           sender,
		Receiver:         receiver,
		Amount:           inst.Amount,
		Memo:             []byte{},
		Slot:             slot,
		Asset:            spl,
		IsCrossChainCall: false,
		RevertOptions:    inst.RevertOptions,
	}, nil
}

// parseAsDepositSPLAndCall tries to parse instruction as deposit_spl_and_call
func parseAsDepositSPLAndCall(
	tx *solana.Transaction,
	instruction solana.CompiledInstruction,
	slot uint64,
) (*Inbound, error) {
	// try deserializing instruction as a deposit_spl_and_call
	instDepositAndCall := DepositSPLAndCallInstructionParams{}
	err := borsh.Deserialize(&instDepositAndCall, instruction.Data)
	if err != nil {
		return nil, nil
	}

	// check if the instruction is a deposit_spl_and_call or not, if not, skip parsing
	if instDepositAndCall.Discriminator != DiscriminatorDepositSPLAndCall {
		return nil, nil
	}

	receiver, err := parseReceiver(instDepositAndCall.Receiver)
	if err != nil {
		return nil, err
	}

	// get the sender and spl addresses
	sender, spl, err := getSignerAndSPLFromDepositSPLAccounts(tx, &instruction)
	if err != nil {
		return nil, err
	}

	return &Inbound{
		Sender:           sender,
		Receiver:         receiver,
		Amount:           instDepositAndCall.Amount,
		Memo:             instDepositAndCall.Memo,
		Slot:             slot,
		Asset:            spl,
		IsCrossChainCall: true,
		RevertOptions:    instDepositAndCall.RevertOptions,
	}, nil
}

// ParseInboundAsCall tries to parse an instruction as a call.
// It returns nil if the instruction can't be parsed as a call.
func ParseInboundAsCall(
	tx *solana.Transaction,
	instruction solana.CompiledInstruction,
	slot uint64,
) (*Inbound, error) {
	// try deserializing instruction as a call
	inst := CallInstructionParams{}
	err := borsh.Deserialize(&inst, instruction.Data)
	if err != nil {
		return nil, nil
	}

	// check if the instruction is a call or not, if not, skip parsing
	if inst.Discriminator != DiscriminatorCall {
		return nil, nil
	}

	// get the sender address (skip if unable to parse signer address)
	instructionAccounts, err := instruction.ResolveInstructionAccounts(&tx.Message)
	if err != nil {
		return nil, err
	}
	if len(instructionAccounts) < accountsNumberCall {
		return nil, fmt.Errorf("want required 1 signer account, got %d", len(instructionAccounts))
	}

	// parse receiver
	receiver, err := parseReceiver(inst.Receiver)
	if err != nil {
		return nil, err
	}

	return &Inbound{
		Sender:           instructionAccounts[0].PublicKey.String(),
		Receiver:         receiver,
		Amount:           0,
		Memo:             inst.Memo,
		Slot:             slot,
		Asset:            "", // no asset for call
		IsCrossChainCall: true,
		RevertOptions:    inst.RevertOptions,
	}, nil
}

// getSignerDeposit returns the signer address of the deposit instruction
func getSignerDeposit(tx *solana.Transaction, inst *solana.CompiledInstruction) (string, error) {
	instructionAccounts, err := inst.ResolveInstructionAccounts(&tx.Message)
	if err != nil {
		return "", err
	}

	// there should be at least all mandatory accounts for a deposit instruction
	if len(instructionAccounts) < accountsNumDeposit {
		return "", fmt.Errorf("want required %d accounts, got %d", accountsNumDeposit, len(instructionAccounts))
	}

	return instructionAccounts[0].PublicKey.String(), nil
}

// getSignerAndSPLFromDepositSPLAccounts returns the signer and spl address of the deposit_spl instruction
func getSignerAndSPLFromDepositSPLAccounts(
	tx *solana.Transaction,
	inst *solana.CompiledInstruction,
) (string, string, error) {
	instructionAccounts, err := inst.ResolveInstructionAccounts(&tx.Message)
	if err != nil {
		return "", "", err
	}

	// there should be at least all mandatory accounts for a deposit spl instruction
	if len(instructionAccounts) < accountsNumberDepositSPL {
		return "", "", fmt.Errorf(
			"want required %d accounts, got %d",
			accountsNumberDepositSPL,
			len(instructionAccounts),
		)
	}

	signer := instructionAccounts[0].PublicKey.String()
	spl := instructionAccounts[3].PublicKey.String()

	return signer, spl, nil
}

// parseReceiver parses the receiver bytes into a Ethereum address string
func parseReceiver(receiver ethcommon.Address) (string, error) {
	if receiver == (ethcommon.Address{}) {
		return "", fmt.Errorf("invalid receiver address: %v", receiver)
	}

	return receiver.Hex(), nil
}
