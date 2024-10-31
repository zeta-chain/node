package solana

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gagliardetto/solana-go"
)

// MsgWithdraw is the message for the Solana gateway withdraw/withdraw_spl instruction
type MsgWithdraw struct {
	// chainID is the chain ID of Solana chain
	chainID uint64

	// Nonce is the nonce for the withdraw/withdraw_spl
	nonce uint64

	// amount is the lamports amount for the withdraw/withdraw_spl
	amount uint64

	// To is the recipient address for the withdraw/withdraw_spl
	to solana.PublicKey

	// signature is the signature of the message
	signature [65]byte
}

// NewMsgWithdraw returns a new withdraw message
func NewMsgWithdraw(chainID, nonce, amount uint64, to solana.PublicKey) *MsgWithdraw {
	return &MsgWithdraw{
		chainID: chainID,
		nonce:   nonce,
		amount:  amount,
		to:      to,
	}
}

// ChainID returns the chain ID of the message
func (msg *MsgWithdraw) ChainID() uint64 {
	return msg.chainID
}

// Nonce returns the nonce of the message
func (msg *MsgWithdraw) Nonce() uint64 {
	return msg.nonce
}

// Amount returns the amount of the message
func (msg *MsgWithdraw) Amount() uint64 {
	return msg.amount
}

// To returns the recipient address of the message
func (msg *MsgWithdraw) To() solana.PublicKey {
	return msg.to
}

// Hash packs the withdraw message and computes the hash
func (msg *MsgWithdraw) Hash() [32]byte {
	var message []byte
	buff := make([]byte, 8)

	message = append(message, []byte("withdraw")...)

	binary.BigEndian.PutUint64(buff, msg.chainID)
	message = append(message, buff...)

	binary.BigEndian.PutUint64(buff, msg.nonce)
	message = append(message, buff...)

	binary.BigEndian.PutUint64(buff, msg.amount)
	message = append(message, buff...)

	message = append(message, msg.to.Bytes()...)

	return crypto.Keccak256Hash(message)
}

// SetSignature attaches the signature to the message
func (msg *MsgWithdraw) SetSignature(signature [65]byte) *MsgWithdraw {
	msg.signature = signature
	return msg
}

// SigRSV returns the full 65-byte [R+S+V] signature
func (msg *MsgWithdraw) SigRSV() [65]byte {
	return msg.signature
}

// SigRS returns the 64-byte [R+S] core part of the signature
func (msg *MsgWithdraw) SigRS() [64]byte {
	var sig [64]byte
	copy(sig[:], msg.signature[:64])
	return sig
}

// SigV returns the V part (recovery ID) of the signature
func (msg *MsgWithdraw) SigV() uint8 {
	return msg.signature[64]
}

// Signer returns the signer of the message
func (msg *MsgWithdraw) Signer() (common.Address, error) {
	msgHash := msg.Hash()
	msgSig := msg.SigRSV()

	return RecoverSigner(msgHash[:], msgSig[:])
}

// MsgWhitelist is the message for the Solana gateway whitelist_spl_mint instruction
type MsgWhitelist struct {
	// whitelistCandidate is the whitelist candidate
	whitelistCandidate solana.PublicKey
	whitelistEntry     solana.PublicKey

	// chainID is the chain ID of Solana chain
	chainID uint64

	// Nonce is the nonce for the withdraw/withdraw_spl
	nonce uint64

	// signature is the signature of the message
	signature [65]byte
}

// NewMsgWhitelist returns a new whitelist_spl_mint message
func NewMsgWhitelist(
	whitelistCandidate solana.PublicKey,
	whitelistEntry solana.PublicKey,
	chainID, nonce uint64,
) *MsgWhitelist {
	return &MsgWhitelist{
		whitelistCandidate: whitelistCandidate,
		whitelistEntry:     whitelistEntry,
		chainID:            chainID,
		nonce:              nonce,
	}
}

// To returns the recipient address of the message
func (msg *MsgWhitelist) WhitelistCandidate() solana.PublicKey {
	return msg.whitelistCandidate
}

func (msg *MsgWhitelist) WhitelistEntry() solana.PublicKey {
	return msg.whitelistEntry
}

// ChainID returns the chain ID of the message
func (msg *MsgWhitelist) ChainID() uint64 {
	return msg.chainID
}

// Nonce returns the nonce of the message
func (msg *MsgWhitelist) Nonce() uint64 {
	return msg.nonce
}

// Hash packs the whitelist message and computes the hash
func (msg *MsgWhitelist) Hash() [32]byte {
	var message []byte
	buff := make([]byte, 8)

	message = append(message, []byte("whitelist_spl_mint")...)

	binary.BigEndian.PutUint64(buff, msg.chainID)
	message = append(message, buff...)

	message = append(message, msg.whitelistCandidate.Bytes()...)

	binary.BigEndian.PutUint64(buff, msg.nonce)
	message = append(message, buff...)

	return crypto.Keccak256Hash(message)
}

// SetSignature attaches the signature to the message
func (msg *MsgWhitelist) SetSignature(signature [65]byte) *MsgWhitelist {
	msg.signature = signature
	return msg
}

// SigRSV returns the full 65-byte [R+S+V] signature
func (msg *MsgWhitelist) SigRSV() [65]byte {
	return msg.signature
}

// SigRS returns the 64-byte [R+S] core part of the signature
func (msg *MsgWhitelist) SigRS() [64]byte {
	var sig [64]byte
	copy(sig[:], msg.signature[:64])
	return sig
}

// SigV returns the V part (recovery ID) of the signature
func (msg *MsgWhitelist) SigV() uint8 {
	return msg.signature[64]
}

// Signer returns the signer of the message
func (msg *MsgWhitelist) Signer() (common.Address, error) {
	msgHash := msg.Hash()
	msgSig := msg.SigRSV()

	return RecoverSigner(msgHash[:], msgSig[:])
}
