package contract

// PdaInfo represents the PDA for the gateway program
type PdaInfo struct {
	// Discriminator is the unique identifier for the PDA
	Discriminator [8]byte

	// Nonce is the current nonce for the PDA
	Nonce uint64

	// TssAddress is the TSS address for the PDA
	TssAddress [20]byte

	// Authority is the authority for the PDA
	Authority [32]byte

	// ChainId is the chain ID for the gateway program
	// TODO: this field exists in latest version of gateway program, but not in the current e2e test program
	// ChainId uint64
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
