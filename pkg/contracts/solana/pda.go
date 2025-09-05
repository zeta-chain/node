package solana

import ethcommon "github.com/ethereum/go-ethereum/common"

// PdaInfo represents the PDA for the gateway program
type PdaInfo struct {
	// Discriminator is the unique identifier for the PDA
	Discriminator [8]byte

	// Nonce is the current nonce for the PDA
	Nonce uint64

	// TssAddress is the TSS address for the PDA
	TssAddress ethcommon.Address

	// Authority is the authority for the PDA
	Authority [32]byte

	// ChainId is the Solana chain id
	ChainID uint64

	// DepositPaused is the flag to indicate if the deposit is paused
	DepositPaused bool
}
