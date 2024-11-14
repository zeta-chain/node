package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

const (
	// MaxSignaturesPerTicker is the maximum number of signatures to process on a ticker
	MaxSignaturesPerTicker = 100
)

type Deposit struct {
	Sender string
	Amount uint64
	Memo   []byte
	Slot   uint64
	Asset  string
}

// ParseInboundAsDeposit tries to parse an instruction as a 'deposit'.
// It returns nil if the instruction can't be parsed as a 'deposit'.
func ParseInboundAsDeposit(
	tx *solana.Transaction,
	instructionIndex int,
	slot uint64,
) (*Deposit, error) {
	// get instruction by index
	instruction := tx.Message.Instructions[instructionIndex]

	// try deserializing instruction as a 'deposit'
	var inst DepositInstructionParams
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

	return &Deposit{
		Sender: sender,
		Amount: inst.Amount,
		Memo:   inst.Memo,
		Slot:   slot,
		Asset:  "", // no asset for gas token SOL
	}, nil
}

// ParseInboundAsDepositSPL tries to parse an instruction as a 'deposit_spl_token'.
// It returns nil if the instruction can't be parsed as a 'deposit_spl_token'.
func ParseInboundAsDepositSPL(
	tx *solana.Transaction,
	instructionIndex int,
	slot uint64,
) (*Deposit, error) {
	// get instruction by index
	instruction := tx.Message.Instructions[instructionIndex]

	// try deserializing instruction as a 'deposit_spl_token'
	var inst DepositSPLInstructionParams
	err := borsh.Deserialize(&inst, instruction.Data)
	if err != nil {
		return nil, nil
	}

	// check if the instruction is a deposit spl or not, if not, skip parsing
	if inst.Discriminator != DiscriminatorDepositSPL {
		return nil, nil
	}

	// get the sender and spl addresses
	sender, spl, err := getSignerAndSPLFromDepositSPLAccounts(tx, &instruction)
	if err != nil {
		return nil, err
	}

	return &Deposit{
		Sender: sender,
		Amount: inst.Amount,
		Memo:   inst.Memo,
		Slot:   slot,
		Asset:  spl,
	}, nil
}

// getSignerDeposit returns the signer address of the deposit instruction
func getSignerDeposit(tx *solana.Transaction, inst *solana.CompiledInstruction) (string, error) {
	instructionAccounts, err := inst.ResolveInstructionAccounts(&tx.Message)
	if err != nil {
		return "", err
	}

	// there should be 3 accounts for a deposit instruction
	if len(instructionAccounts) != accountsNumDeposit {
		return "", fmt.Errorf("want %d accounts, got %d", accountsNumDeposit, len(instructionAccounts))
	}

	// the accounts are [signer, pda, system_program]
	// check if first account is signer
	if !instructionAccounts[0].IsSigner {
		return "", fmt.Errorf("unexpected signer %s", instructionAccounts[0].PublicKey.String())
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

	// there should be 7 accounts for a deposit spl instruction
	if len(instructionAccounts) != accountsNumberDepositSPL {
		return "", "", fmt.Errorf(
			"want %d accounts, got %d",
			accountsNumberDepositSPL,
			len(instructionAccounts),
		)
	}
	// the accounts are [signer, pda, whitelist_entry, mint_account, token_program, from, to]
	// check if first account is signer
	if !instructionAccounts[0].IsSigner {
		return "", "", fmt.Errorf("unexpected signer %s", instructionAccounts[0].PublicKey.String())
	}

	signer := instructionAccounts[0].PublicKey.String()
	spl := instructionAccounts[3].PublicKey.String()

	return signer, spl, nil
}
