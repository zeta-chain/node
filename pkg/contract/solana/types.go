package solana // PdaInfo represents the PDA for the gateway program
import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type PdaInfo struct {
	// Discriminator is the unique identifier for the PDA
	Discriminator [8]byte

	// Nonce is the current nonce for the PDA
	Nonce uint64

	// TssAddress is the TSS address for the PDA
	TssAddress [20]byte

	// Authority is the authority for the PDA
	Authority [32]byte

	// ChainId is the Solana chain id
	ChainID uint64
}

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

// RecoverSigner recover the ECDSA signer from given message hash and signature
func RecoverSigner(msgHash []byte, msgSig []byte) (signer common.Address, err error) {
	// recover the public key
	pubKey, err := crypto.SigToPub(msgHash, msgSig)
	if err != nil {
		return
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}
