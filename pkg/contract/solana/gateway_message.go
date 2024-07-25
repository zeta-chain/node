package solana

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gagliardetto/solana-go"
)

// MsgWithdraw is the message for the Solana gateway withdraw instruction
type MsgWithdraw struct {
	// ChainID is the chain ID of Solana chain
	ChainID uint64

	// Nonce is the nonce for the withdraw
	Nonce uint64

	// Amount is the lamports amount for the withdraw
	Amount uint64

	// To is the recipient address for the withdraw
	To solana.PublicKey

	// Signature is the signature of the message
	Signature [65]byte
}

// NewMsgWithdraw returns a new MsgWithdraw
func NewMsgWithdraw(chainID, nonce, amount uint64, to solana.PublicKey) *MsgWithdraw {
	return &MsgWithdraw{
		ChainID: chainID,
		Nonce:   nonce,
		Amount:  amount,
		To:      to,
	}
}

// Hash packs the withdraw message and computes the hash
func (msg *MsgWithdraw) Hash() [32]byte {
	var message []byte
	buff := make([]byte, 8)

	binary.BigEndian.PutUint64(buff, msg.ChainID)
	message = append(message, buff...)

	binary.BigEndian.PutUint64(buff, msg.Nonce)
	message = append(message, buff...)

	binary.BigEndian.PutUint64(buff, msg.Amount)
	message = append(message, buff...)

	message = append(message, msg.To.Bytes()...)

	return crypto.Keccak256Hash(message)
}

// WithSignature attaches the signature to the message
func (msg *MsgWithdraw) WithSignature(signature [65]byte) *MsgWithdraw {
	msg.Signature = signature
	return msg
}

// Sig64 returns the 64-byte [R+S] core part of the signature
func (msg *MsgWithdraw) SigRS() [64]byte {
	var sig [64]byte
	copy(sig[:], msg.Signature[:64])
	return sig
}

// SigV returns the V part (recovery ID) of the signature
func (msg *MsgWithdraw) SigV() uint8 {
	return msg.Signature[64]
}
