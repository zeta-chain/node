package solana

// DiscriminatorDeposit returns the discriminator for Solana gateway deposit instruction
func DiscriminatorDeposit() []byte {
	return []byte{242, 35, 198, 137, 82, 225, 242, 182}
}

const (
	// PDASeed is the seed for the Solana gateway program derived address
	PDASeed = "meta"

	// AccountsNumberOfDeposit is the number of accounts required for Solana gateway deposit instruction
	// [signer, pda, system_program, gateway_program]
	AccountsNumDeposit = 4

	// MaxSignaturesPerTicker is the maximum number of signatures to process on a ticker
	MaxSignaturesPerTicker = 100
)
