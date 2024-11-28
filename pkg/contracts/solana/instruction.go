package solana

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
	"github.com/pkg/errors"
)

const (
	// MsgWithdrawSPLTokenSuccess is the success message for withdraw_spl_token instruction
	// #nosec G101 not a hardcoded credential
	MsgWithdrawSPLTokenSuccess = "withdraw spl token successfully"

	// MsgWithdrawSPLTokenNonExistentAta is the log message printed when recipient ATA does not exist
	MsgWithdrawSPLTokenNonExistentAta = "recipient ATA account does not exist"
)

// InitializeParams contains the parameters for a gateway initialize instruction
type InitializeParams struct {
	// Discriminator is the unique identifier for the initialize instruction
	Discriminator [8]byte

	// TssAddress is the TSS address
	TssAddress [20]byte

	// ChainID is the chain ID for the gateway program
	ChainID uint64
}

// InitializeRentPayerParams contains the parameters for a gateway initialize_rent_payer instruction
type InitializeRentPayerParams struct {
	// Discriminator is the unique identifier for the initialize_rent_payer instruction
	Discriminator [8]byte
}

// DepositInstructionParams contains the parameters for a gateway deposit instruction
type DepositInstructionParams struct {
	// Discriminator is the unique identifier for the deposit instruction
	Discriminator [8]byte

	// Amount is the lamports amount for the deposit
	Amount uint64

	// Memo is the memo for the deposit
	Memo []byte
}

// DepositSPLInstructionParams contains the parameters for a gateway deposit spl instruction
type DepositSPLInstructionParams struct {
	// Discriminator is the unique identifier for the deposit instruction
	Discriminator [8]byte

	// Amount is the lamports amount for the deposit
	Amount uint64

	// Memo is the memo for the deposit
	Memo []byte
}

// OutboundInstruction is the interface for all gateway outbound instructions
type OutboundInstruction interface {
	// Signer returns the signer of the instruction
	Signer() (common.Address, error)

	// GatewayNonce returns the nonce of the instruction
	GatewayNonce() uint64

	// TokenAmount returns the amount of the instruction
	TokenAmount() uint64

	// Failed returns true if the instruction logs indicate failure
	Failed(logMessages []string) bool
}

var _ OutboundInstruction = (*WithdrawInstructionParams)(nil)

// WithdrawInstructionParams contains the parameters for a gateway withdraw instruction
type WithdrawInstructionParams struct {
	// Discriminator is the unique identifier for the withdraw instruction
	Discriminator [8]byte

	// Amount is the lamports amount for the withdraw
	Amount uint64

	// Signature is the ECDSA signature (by TSS) for the withdraw
	Signature [64]byte

	// RecoveryID is the recovery ID used to recover the public key from ECDSA signature
	RecoveryID uint8

	// MessageHash is the hash of the message signed by TSS
	MessageHash [32]byte

	// Nonce is the nonce for the withdraw
	Nonce uint64
}

// Signer returns the signer of the signature contained
func (inst *WithdrawInstructionParams) Signer() (signer common.Address, err error) {
	var signature [65]byte
	copy(signature[:], inst.Signature[:64])
	signature[64] = inst.RecoveryID

	return RecoverSigner(inst.MessageHash[:], signature[:])
}

// GatewayNonce returns the nonce of the instruction
func (inst *WithdrawInstructionParams) GatewayNonce() uint64 {
	return inst.Nonce
}

// TokenAmount returns the amount of the instruction
func (inst *WithdrawInstructionParams) TokenAmount() uint64 {
	return inst.Amount
}

// Failed always returns false for a 'withdraw' without checking the logs
func (inst *WithdrawInstructionParams) Failed(_ []string) bool {
	return false
}

// ParseInstructionWithdraw tries to parse the instruction as a 'withdraw'.
// It returns nil if the instruction can't be parsed as a 'withdraw'.
func ParseInstructionWithdraw(instruction solana.CompiledInstruction) (*WithdrawInstructionParams, error) {
	// try deserializing instruction as a 'withdraw'
	inst := &WithdrawInstructionParams{}
	err := borsh.Deserialize(inst, instruction.Data)
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing instruction")
	}

	// check the discriminator to ensure it's a 'withdraw' instruction
	if inst.Discriminator != DiscriminatorWithdraw {
		return nil, fmt.Errorf("not a withdraw instruction: %v", inst.Discriminator)
	}

	return inst, nil
}

type WithdrawSPLInstructionParams struct {
	// Discriminator is the unique identifier for the withdraw instruction
	Discriminator [8]byte

	// Decimals is decimals for spl token
	Decimals uint8

	// Amount is the lamports amount for the withdraw
	Amount uint64

	// Signature is the ECDSA signature (by TSS) for the withdraw
	Signature [64]byte

	// RecoveryID is the recovery ID used to recover the public key from ECDSA signature
	RecoveryID uint8

	// MessageHash is the hash of the message signed by TSS
	MessageHash [32]byte

	// Nonce is the nonce for the withdraw
	Nonce uint64
}

// Signer returns the signer of the signature contained
func (inst *WithdrawSPLInstructionParams) Signer() (signer common.Address, err error) {
	var signature [65]byte
	copy(signature[:], inst.Signature[:64])
	signature[64] = inst.RecoveryID

	return RecoverSigner(inst.MessageHash[:], signature[:])
}

// GatewayNonce returns the nonce of the instruction
func (inst *WithdrawSPLInstructionParams) GatewayNonce() uint64 {
	return inst.Nonce
}

// TokenAmount returns the amount of the instruction
func (inst *WithdrawSPLInstructionParams) TokenAmount() uint64 {
	return inst.Amount
}

// Failed returns true if the logs of the 'withdraw_spl_token' instruction indicate failure.
//
// Note: SPL token transfer cannot be done if the recipient ATA does not exist.
func (inst *WithdrawSPLInstructionParams) Failed(logMessages []string) bool {
	// Assumption: only one of the two messages will be present in the logs.
	// If both messages are present, it could imply a program bug or a malicious attack.
	// In such case, the function treats the transaction as successful to minimize the attack surface,
	// bacause a fabricated failure could be used to trick zetacore into refunding the withdrawer (if implemented in the future).
	return !containsLogMessage(logMessages, MsgWithdrawSPLTokenSuccess) &&
		containsLogMessage(logMessages, MsgWithdrawSPLTokenNonExistentAta)
}

// ParseInstructionWithdrawSPL tries to parse the instruction as a 'withdraw_spl_token'.
// It returns nil if the instruction can't be parsed as a 'withdraw_spl_token'.
func ParseInstructionWithdrawSPL(instruction solana.CompiledInstruction) (*WithdrawSPLInstructionParams, error) {
	// try deserializing instruction as a 'withdraw_spl_token'
	inst := &WithdrawSPLInstructionParams{}
	err := borsh.Deserialize(inst, instruction.Data)
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing instruction")
	}

	// check the discriminator to ensure it's a 'withdraw_spl_token' instruction
	if inst.Discriminator != DiscriminatorWithdrawSPL {
		return nil, fmt.Errorf("not a withdraw_spl_token instruction: %v", inst.Discriminator)
	}

	return inst, nil
}

// RecoverSigner recover the ECDSA signer from given message hash and signature
func RecoverSigner(msgHash []byte, msgSig []byte) (signer common.Address, err error) {
	// recover the public key
	pubKey, err := crypto.SigToPub(msgHash, msgSig)
	if err != nil {
		return
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}

var _ OutboundInstruction = (*WhitelistInstructionParams)(nil)

// WhitelistInstructionParams contains the parameters for a gateway whitelist_spl_mint instruction
type WhitelistInstructionParams struct {
	// Discriminator is the unique identifier for the whitelist instruction
	Discriminator [8]byte

	// Signature is the ECDSA signature (by TSS) for the whitelist
	Signature [64]byte

	// RecoveryID is the recovery ID used to recover the public key from ECDSA signature
	RecoveryID uint8

	// MessageHash is the hash of the message signed by TSS
	MessageHash [32]byte

	// Nonce is the nonce for the whitelist
	Nonce uint64
}

// Signer returns the signer of the signature contained
func (inst *WhitelistInstructionParams) Signer() (signer common.Address, err error) {
	var signature [65]byte
	copy(signature[:], inst.Signature[:64])
	signature[64] = inst.RecoveryID

	return RecoverSigner(inst.MessageHash[:], signature[:])
}

// GatewayNonce returns the nonce of the instruction
func (inst *WhitelistInstructionParams) GatewayNonce() uint64 {
	return inst.Nonce
}

// TokenAmount returns the amount of the instruction
func (inst *WhitelistInstructionParams) TokenAmount() uint64 {
	return 0
}

// Failed always returns false for a 'whitelist_spl_mint' without checking the logs
func (inst *WhitelistInstructionParams) Failed(_ []string) bool {
	return true
}

// ParseInstructionWhitelist tries to parse the instruction as a 'whitelist_spl_mint'.
// It returns nil if the instruction can't be parsed as a 'whitelist_spl_mint'.
func ParseInstructionWhitelist(instruction solana.CompiledInstruction) (*WhitelistInstructionParams, error) {
	// try deserializing instruction as a 'whitelist_spl_mint'
	inst := &WhitelistInstructionParams{}
	err := borsh.Deserialize(inst, instruction.Data)
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing instruction")
	}

	// check the discriminator to ensure it's a 'whitelist_spl_mint' instruction
	if inst.Discriminator != DiscriminatorWhitelistSplMint {
		return nil, fmt.Errorf("not a whitelist_spl_mint instruction: %v", inst.Discriminator)
	}

	return inst, nil
}

// containsLogMessage returns true if any of the log messages contains the 'msgSearch'
func containsLogMessage(logMessages []string, msgSearch string) bool {
	return slices.IndexFunc(logMessages, func(msg string) bool {
		return strings.Contains(msg, msgSearch)
	}) != -1
}
