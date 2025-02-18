package solana

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
	"github.com/pkg/errors"
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

// DepositInstructionParams contains the parameters for a gateway deposit instruction
type DepositInstructionParams struct {
	// Discriminator is the unique identifier for the deposit instruction
	Discriminator [8]byte

	// Amount is the lamports amount for the deposit
	Amount uint64

	// Receiver is the receiver for the deposit
	Receiver [20]byte
}

// DepositAndCallInstructionParams contains the parameters for a gateway deposit_and_call instruction
type DepositAndCallInstructionParams struct {
	// Discriminator is the unique identifier for the deposit_and_call instruction
	Discriminator [8]byte

	// Amount is the lamports amount for the deposit_and_call
	Amount uint64

	// Receiver is the receiver for the deposit_and_call
	Receiver [20]byte

	// Memo is the memo for the deposit_and_call
	Memo []byte
}

// DepositSPLInstructionParams contains the parameters for a gateway deposit_spl instruction
type DepositSPLInstructionParams struct {
	// Discriminator is the unique identifier for the deposit_spl instruction
	Discriminator [8]byte

	// Amount is the lamports amount for the deposit_spl
	Amount uint64

	// Receiver is the receiver for the deposit_spl
	Receiver [20]byte
}

// DepositSPLAndCallInstructionParams contains the parameters for a gateway deposit_spl_and_call instruction
type DepositSPLAndCallInstructionParams struct {
	// Discriminator is the unique identifier for the deposit_spl_and_call instruction
	Discriminator [8]byte

	// Amount is the lamports amount for the deposit_spl_and_call
	Amount uint64

	// Receiver is the receiver for the deposit_spl_and_call
	Receiver [20]byte

	// Memo is the memo for the deposit_spl_and_call
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

	// Failed returns true if outbound failed
	Failed() bool
}

var _ OutboundInstruction = (*IncrementNonceInstructionParams)(nil)

// IncrementNonceInstructionParams contains the parameters for a gateway increment_nonce instruction
type IncrementNonceInstructionParams struct {
	// Discriminator is the unique identifier for the increment_nonce instruction
	Discriminator [8]byte

	// Amount is the lamports amount for the increment_nonce
	Amount uint64

	// Signature is the ECDSA signature (by TSS) for the increment_nonce
	Signature [64]byte

	// RecoveryID is the recovery ID used to recover the public key from ECDSA signature
	RecoveryID uint8

	// MessageHash is the hash of the message signed by TSS
	MessageHash [32]byte

	// Nonce is the nonce for the increment_nonce
	Nonce uint64
}

// Failed returns true if outbound failed
func (inst *IncrementNonceInstructionParams) Failed() bool {
	return true
}

// Signer returns the signer of the signature contained
func (inst *IncrementNonceInstructionParams) Signer() (signer common.Address, err error) {
	var signature [65]byte
	copy(signature[:], inst.Signature[:64])
	signature[64] = inst.RecoveryID

	return RecoverSigner(inst.MessageHash[:], signature[:])
}

// GatewayNonce returns the nonce of the instruction
func (inst *IncrementNonceInstructionParams) GatewayNonce() uint64 {
	return inst.Nonce
}

// TokenAmount returns the amount of the instruction
func (inst *IncrementNonceInstructionParams) TokenAmount() uint64 {
	return inst.Amount
}

// ParseInstructionIncrementNonce tries to parse the instruction as a 'increment_nonce'.
// It returns nil if the instruction can't be parsed as a 'increment_nonce'.
func ParseInstructionIncrementNonce(
	instruction solana.CompiledInstruction,
) (*IncrementNonceInstructionParams, error) {
	// try deserializing instruction as a 'increment_nonce'
	inst := &IncrementNonceInstructionParams{}
	err := borsh.Deserialize(inst, instruction.Data)
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing instruction")
	}

	// check the discriminator to ensure it's a 'increment_nonce' instruction
	if inst.Discriminator != DiscriminatorIncrementNonce {
		return nil, fmt.Errorf("not an increment_nonce instruction: %v", inst.Discriminator)
	}

	return inst, nil
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

// Failed returns true if outbound failed
func (inst *WithdrawInstructionParams) Failed() bool {
	return false
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
		return nil, fmt.Errorf("not a withdraw instruction: %v %v", inst.Discriminator, DiscriminatorWithdraw)
	}

	return inst, nil
}

var _ OutboundInstruction = (*ExecuteInstructionParams)(nil)

// ExecuteInstructionParams contains the parameters for a gateway execute instruction
type ExecuteInstructionParams struct {
	// Discriminator is the unique identifier for the execute instruction
	Discriminator [8]byte

	// Amount is the lamports amount for the execute
	Amount uint64

	// Sender from zetachain
	Sender [20]byte

	// Data for connected program
	Data []byte

	// Signature is the ECDSA signature (by TSS) for the execute
	Signature [64]byte

	// RecoveryID is the recovery ID used to recover the public key from ECDSA signature
	RecoveryID uint8

	// MessageHash is the hash of the message signed by TSS
	MessageHash [32]byte

	// Nonce is the nonce for the execute
	Nonce uint64
}

// Failed returns true if outbound failed
func (inst *ExecuteInstructionParams) Failed() bool {
	return false
}

// Signer returns the signer of the signature contained
func (inst *ExecuteInstructionParams) Signer() (signer common.Address, err error) {
	var signature [65]byte
	copy(signature[:], inst.Signature[:64])
	signature[64] = inst.RecoveryID

	return RecoverSigner(inst.MessageHash[:], signature[:])
}

// GatewayNonce returns the nonce of the instruction
func (inst *ExecuteInstructionParams) GatewayNonce() uint64 {
	return inst.Nonce
}

// TokenAmount returns the amount of the instruction
func (inst *ExecuteInstructionParams) TokenAmount() uint64 {
	return inst.Amount
}

// ParseInstructionExecute tries to parse the instruction as a 'execute'.
// It returns nil if the instruction can't be parsed as a 'execute'.
func ParseInstructionExecute(instruction solana.CompiledInstruction) (*ExecuteInstructionParams, error) {
	// try deserializing instruction as a 'execute'
	inst := &ExecuteInstructionParams{}
	err := borsh.Deserialize(inst, instruction.Data)
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing instruction")
	}

	// check the discriminator to ensure it's a 'execute' instruction
	if inst.Discriminator != DiscriminatorExecute {
		return nil, fmt.Errorf("not an execute instruction: %v", inst.Discriminator)
	}

	return inst, nil
}

var _ OutboundInstruction = (*WithdrawSPLInstructionParams)(nil)

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

// Failed returns true if outbound failed
func (inst *WithdrawSPLInstructionParams) Failed() bool {
	return false
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

// ParseInstructionWithdraw tries to parse the instruction as a 'withdraw'.
// It returns nil if the instruction can't be parsed as a 'withdraw'.
func ParseInstructionWithdrawSPL(instruction solana.CompiledInstruction) (*WithdrawSPLInstructionParams, error) {
	// try deserializing instruction as a 'withdraw'
	inst := &WithdrawSPLInstructionParams{}
	err := borsh.Deserialize(inst, instruction.Data)
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing instruction")
	}

	// check the discriminator to ensure it's a 'withdraw' instruction
	if inst.Discriminator != DiscriminatorWithdrawSPL {
		return nil, fmt.Errorf("not a withdraw instruction: %v", inst.Discriminator)
	}

	return inst, nil
}

var _ OutboundInstruction = (*ExecuteSPLInstructionParams)(nil)

type ExecuteSPLInstructionParams struct {
	// Discriminator is the unique identifier for the execute spl instruction
	Discriminator [8]byte

	// Decimals is decimals for spl token
	Decimals uint8

	// Amount is the lamports amount for the withdraw
	Amount uint64

	// Sender from zetachain
	Sender [20]byte

	// Data for connected program
	Data []byte

	// Signature is the ECDSA signature (by TSS) for the withdraw
	Signature [64]byte

	// RecoveryID is the recovery ID used to recover the public key from ECDSA signature
	RecoveryID uint8

	// MessageHash is the hash of the message signed by TSS
	MessageHash [32]byte

	// Nonce is the nonce for the withdraw
	Nonce uint64
}

// Failed returns true if outbound failed
func (inst *ExecuteSPLInstructionParams) Failed() bool {
	return false
}

// Signer returns the signer of the signature contained
func (inst *ExecuteSPLInstructionParams) Signer() (signer common.Address, err error) {
	var signature [65]byte
	copy(signature[:], inst.Signature[:64])
	signature[64] = inst.RecoveryID

	return RecoverSigner(inst.MessageHash[:], signature[:])
}

// GatewayNonce returns the nonce of the instruction
func (inst *ExecuteSPLInstructionParams) GatewayNonce() uint64 {
	return inst.Nonce
}

// TokenAmount returns the amount of the instruction
func (inst *ExecuteSPLInstructionParams) TokenAmount() uint64 {
	return inst.Amount
}

// ParseInstructionExecuteSPL tries to parse the instruction as a 'execute_spl_token'.
// It returns nil if the instruction can't be parsed as a 'execute_spl_token'.
func ParseInstructionExecuteSPL(instruction solana.CompiledInstruction) (*ExecuteSPLInstructionParams, error) {
	// try deserializing instruction as a 'execute_spl_token'
	inst := &ExecuteSPLInstructionParams{}
	err := borsh.Deserialize(inst, instruction.Data)
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing instruction")
	}

	// check the discriminator to ensure it's a 'execute_spl_token' instruction
	if inst.Discriminator != DiscriminatorExecuteSPL {
		return nil, fmt.Errorf("not an execute_spl_token instruction: %v", inst.Discriminator)
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

// Failed returns true if outbound failed
func (inst *WhitelistInstructionParams) Failed() bool {
	return false
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

type AccountMeta struct {
	PublicKey  [32]byte
	IsWritable bool
}

// ExecuteMsg describes data and accounts passed to connected programs
type ExecuteMsg struct {
	Accounts []AccountMeta
	Data     []byte
}

// GetExecuteMsgAbi used for abi encoding/decoding execute msg
func GetExecuteMsgAbi() (abi.Arguments, error) {
	MsgAbiType, err := abi.NewType("tuple", "struct Msg", []abi.ArgumentMarshaling{
		{
			Name: "accounts",
			Type: "tuple[]",
			Components: []abi.ArgumentMarshaling{
				{Name: "publicKey", Type: "bytes32"},
				{Name: "isWritable", Type: "bool"},
			},
		},
		{Name: "data", Type: "bytes"},
	})

	if err != nil {
		return abi.Arguments{}, err
	}

	MsgAbiArgs := abi.Arguments{
		{Type: MsgAbiType, Name: "msg"},
	}

	return MsgAbiArgs, nil
}

// DecodeExecuteMsg decodes execute msg using abi decoding
func DecodeExecuteMsg(msgbz []byte) (ExecuteMsg, error) {
	args, err := GetExecuteMsgAbi()
	if err != nil {
		return ExecuteMsg{}, err
	}

	unpacked, err := args.Unpack(msgbz)
	if err != nil {
		return ExecuteMsg{}, err
	}

	jsonMsg, err := json.Marshal(unpacked[0])
	if err != nil {
		return ExecuteMsg{}, err
	}

	var msg ExecuteMsg
	err = json.Unmarshal(jsonMsg, &msg)
	if err != nil {
		return ExecuteMsg{}, err
	}

	return msg, nil
}
