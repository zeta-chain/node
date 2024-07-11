package observer

// DepositInstructionParams contains the parameters for a gateway deposit instruction
type DepositInstructionParams struct {
	// Discriminator is the unique identifier for the deposit instruction
	Discriminator [8]byte

	// Amount is the lamports amount for the deposit
	Amount uint64

	// Memo is the memo for the deposit
	Memo []byte
}
