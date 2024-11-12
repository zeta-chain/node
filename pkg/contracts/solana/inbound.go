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

// GetSignerDeposit returns the signer address of the deposit instruction
// Note: solana-go is not able to parse the AccountMeta 'is_signer' ATM. This is a workaround.
func getSignerDeposit(tx *solana.Transaction, inst *solana.CompiledInstruction) (string, error) {
	// there should be 3 accounts for a deposit instruction
	if len(inst.Accounts) != accountsNumDeposit {
		return "", fmt.Errorf("want %d accounts, got %d", accountsNumDeposit, len(inst.Accounts))
	}

	// sender is the signer account
	return tx.Message.AccountKeys[0].String(), nil
}

func getSignerAndSPLFromDepositSPLAccounts(
	tx *solana.Transaction,
	inst *solana.CompiledInstruction,
) (string, string, error) {
	// there should be 7 accounts for a deposit spl instruction
	if len(inst.Accounts) != accountsNumberDepositSPL {
		return "", "", fmt.Errorf(
			"want %d accounts, got %d",
			accountsNumberDepositSPL,
			len(inst.Accounts),
		)
	}

	// the accounts are [signer, pda, whitelist_entry, mint_account, token_program, from, to]
	signer := tx.Message.AccountKeys[0]
	spl := tx.Message.AccountKeys[inst.Accounts[3]]

	return signer.String(), spl.String(), nil
}
