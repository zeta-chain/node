package solana

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gagliardetto/solana-go"
)

// MsgWithdraw is the message for the Solana gateway withdraw instruction
type MsgWithdraw struct {
	// chainID is the chain ID of Solana chain
	chainID uint64

	// Nonce is the nonce for the withdraw
	nonce uint64

	// amount is the lamports amount for the withdraw
	amount uint64

	// To is the recipient address for the withdraw
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

// MsgWithdrawSPL is the message for the Solana gateway withdraw_spl instruction
type MsgWithdrawSPL struct {
	// chainID is the chain ID of Solana chain
	chainID uint64

	// Nonce is the nonce for the withdraw_spl
	nonce uint64

	// amount is the lamports amount for the withdraw_spl
	amount uint64

	// mintAccount is the address for the spl token
	mintAccount solana.PublicKey

	// decimals of spl token
	decimals uint8

	// to is the recipient address for the withdraw_spl
	to solana.PublicKey

	// recipientAta is the recipient associated token account for the withdraw_spl
	recipientAta solana.PublicKey

	// signature is the signature of the message
	signature [65]byte
}

// NewMsgWithdrawSPL returns a new withdraw spl message
func NewMsgWithdrawSPL(
	chainID, nonce, amount uint64,
	decimals uint8,
	mintAccount, to, toAta solana.PublicKey,
) *MsgWithdrawSPL {
	return &MsgWithdrawSPL{
		chainID:      chainID,
		nonce:        nonce,
		amount:       amount,
		to:           to,
		recipientAta: toAta,
		mintAccount:  mintAccount,
		decimals:     decimals,
	}
}

// ChainID returns the chain ID of the message
func (msg *MsgWithdrawSPL) ChainID() uint64 {
	return msg.chainID
}

// Nonce returns the nonce of the message
func (msg *MsgWithdrawSPL) Nonce() uint64 {
	return msg.nonce
}

// Amount returns the amount of the message
func (msg *MsgWithdrawSPL) Amount() uint64 {
	return msg.amount
}

// To returns the recipient address of the message
func (msg *MsgWithdrawSPL) To() solana.PublicKey {
	return msg.to
}

func (msg *MsgWithdrawSPL) RecipientAta() solana.PublicKey {
	return msg.recipientAta
}

func (msg *MsgWithdrawSPL) MintAccount() solana.PublicKey {
	return msg.mintAccount
}

func (msg *MsgWithdrawSPL) Decimals() uint8 {
	return msg.decimals
}

// Hash packs the withdraw spl message and computes the hash
func (msg *MsgWithdrawSPL) Hash() [32]byte {
	var message []byte
	buff := make([]byte, 8)

	message = append(message, []byte("withdraw_spl_token")...)

	binary.BigEndian.PutUint64(buff, msg.chainID)
	message = append(message, buff...)

	binary.BigEndian.PutUint64(buff, msg.nonce)
	message = append(message, buff...)

	binary.BigEndian.PutUint64(buff, msg.amount)
	message = append(message, buff...)

	message = append(message, msg.mintAccount.Bytes()...)

	message = append(message, msg.recipientAta.Bytes()...)

	return crypto.Keccak256Hash(message)
}

// SetSignature attaches the signature to the message
func (msg *MsgWithdrawSPL) SetSignature(signature [65]byte) *MsgWithdrawSPL {
	msg.signature = signature
	return msg
}

// SigRSV returns the full 65-byte [R+S+V] signature
func (msg *MsgWithdrawSPL) SigRSV() [65]byte {
	return msg.signature
}

// SigRS returns the 64-byte [R+S] core part of the signature
func (msg *MsgWithdrawSPL) SigRS() [64]byte {
	var sig [64]byte
	copy(sig[:], msg.signature[:64])
	return sig
}

// SigV returns the V part (recovery ID) of the signature
func (msg *MsgWithdrawSPL) SigV() uint8 {
	return msg.signature[64]
}

// Signer returns the signer of the message
func (msg *MsgWithdrawSPL) Signer() (common.Address, error) {
	msgHash := msg.Hash()
	msgSig := msg.SigRSV()

	return RecoverSigner(msgHash[:], msgSig[:])
}

// MsgWhitelist is the message for the Solana gateway whitelist_spl_mint instruction
type MsgWhitelist struct {
	// whitelistCandidate is the SPL token to be whitelisted in gateway program
	whitelistCandidate solana.PublicKey

	// whitelistEntry is the entry in gateway program representing whitelisted SPL token
	whitelistEntry solana.PublicKey

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
