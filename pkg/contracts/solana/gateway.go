// Package solana privides structures and constants that are used when interacting with the gateway program on Solana chain.
package solana

import (
	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
)

const (
	// SolanaGatewayProgramID is the program ID of the Solana gateway program
	SolanaGatewayProgramID = "94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d"

	// PDASeed is the seed for the Solana gateway program derived address
	PDASeed = "meta"

	// AccountsNumberOfDeposit is the number of accounts required for Solana gateway deposit instruction
	// [signer, pda, system_program]
	AccountsNumDeposit = 3
)

// DiscriminatorInitialize returns the discriminator for Solana gateway 'initialize' instruction
func DiscriminatorInitialize() [8]byte {
	return [8]byte{175, 175, 109, 31, 13, 152, 155, 237}
}

// DiscriminatorDeposit returns the discriminator for Solana gateway 'deposit' instruction
func DiscriminatorDeposit() [8]byte {
	return [8]byte{242, 35, 198, 137, 82, 225, 242, 182}
}

// DiscriminatorDepositSPL returns the discriminator for Solana gateway 'deposit_spl_token' instruction
func DiscriminatorDepositSPL() [8]byte {
	return [8]byte{86, 172, 212, 121, 63, 233, 96, 144}
}

// DiscriminatorWithdraw returns the discriminator for Solana gateway 'withdraw' instruction
func DiscriminatorWithdraw() [8]byte {
	return [8]byte{183, 18, 70, 156, 148, 109, 161, 34}
}

// DiscriminatorWithdrawSPL returns the discriminator for Solana gateway 'withdraw_spl_token' instruction
func DiscriminatorWithdrawSPL() [8]byte {
	return [8]byte{156, 234, 11, 89, 235, 246, 32}
}

// ParseGatewayAddressAndPda parses the gateway id and program derived address from the given string
func ParseGatewayIDAndPda(address string) (solana.PublicKey, solana.PublicKey, error) {
	var gatewayID, pda solana.PublicKey

	// decode gateway address
	gatewayID, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return gatewayID, pda, errors.Wrap(err, "unable to decode address")
	}

	// compute gateway PDA
	seed := []byte(PDASeed)
	pda, _, err = solana.FindProgramAddress([][]byte{seed}, gatewayID)

	return gatewayID, pda, err
}
