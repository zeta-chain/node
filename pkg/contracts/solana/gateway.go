// Package solana provides structures and constants that are used when interacting with the gateway program on Solana chain.
package solana

import (
	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
	idlgateway "github.com/zeta-chain/protocol-contracts-solana/go-idl/generated"
)

const (
	// PDASeed is the seed for the Solana gateway program derived address
	PDASeed = "meta"

	// accountsNumDeposit is the number of accounts required for Solana gateway deposit instruction
	// [signer, pda, system_program]
	accountsNumDeposit = 3

	// accountsNumberDepositSPL is the number of accounts required for Solana gateway deposit spl instruction
	// [signer, pda, whitelist_entry, mint_account, token_program, from, to, system_program]
	accountsNumberDepositSPL = 8

	// accountsNumberCall is the number of accounts required for Solana gateway call instruction
	// [signer]
	accountsNumberCall = 1
)

var (
	// DiscriminatorInitialize returns the discriminator for Solana gateway 'initialize' instruction
	DiscriminatorInitialize = idlgateway.IDLGateway.GetDiscriminator("initialize")

	DiscriminatorUpdateTss = idlgateway.IDLGateway.GetDiscriminator("update_tss")

	// DiscriminatorDeposit returns the discriminator for Solana gateway 'deposit' instruction
	DiscriminatorDeposit = idlgateway.IDLGateway.GetDiscriminator("deposit")

	// DiscriminatorDeposit returns the discriminator for Solana gateway 'deposit_and_call' instruction
	DiscriminatorDepositAndCall = idlgateway.IDLGateway.GetDiscriminator("deposit_and_call")

	// DiscriminatorDepositSPL returns the discriminator for Solana gateway 'deposit_spl_token' instruction
	DiscriminatorDepositSPL = idlgateway.IDLGateway.GetDiscriminator("deposit_spl_token")

	// DiscriminatorDepositSPLAndCall returns the discriminator for Solana gateway 'deposit_spl_token_and_call' instruction
	DiscriminatorDepositSPLAndCall = idlgateway.IDLGateway.GetDiscriminator("deposit_spl_token_and_call")

	// DiscriminatorCall returns the discriminator for Solana gateway 'call' instruction
	DiscriminatorCall = idlgateway.IDLGateway.GetDiscriminator("call")

	// DiscriminatorWithdraw returns the discriminator for Solana gateway 'withdraw' instruction
	DiscriminatorWithdraw = idlgateway.IDLGateway.GetDiscriminator("withdraw")

	// DiscriminatorExecute returns the discriminator for Solana gateway 'execute' instruction
	DiscriminatorExecute = idlgateway.IDLGateway.GetDiscriminator("execute")

	// DiscriminatorExecuteRevert returns the discriminator for Solana gateway 'execute_revert' instruction
	DiscriminatorExecuteRevert = idlgateway.IDLGateway.GetDiscriminator("execute_revert")

	// DiscriminatorIncrementNonce returns the discriminator for Solana gateway 'increment_nonce' instruction
	DiscriminatorIncrementNonce = idlgateway.IDLGateway.GetDiscriminator("increment_nonce")

	// DiscriminatorExecuteSPL returns the discriminator for Solana gateway 'execute_spl_token' instruction
	DiscriminatorExecuteSPL = idlgateway.IDLGateway.GetDiscriminator("execute_spl_token")

	// DiscriminatorExecuteSPLRevert returns the discriminator for Solana gateway 'execute_spl_token_revert' instruction
	DiscriminatorExecuteSPLRevert = idlgateway.IDLGateway.GetDiscriminator("execute_spl_token_revert")

	// DiscriminatorWithdrawSPL returns the discriminator for Solana gateway 'withdraw_spl_token' instruction
	DiscriminatorWithdrawSPL = idlgateway.IDLGateway.GetDiscriminator("withdraw_spl_token")

	// DiscriminatorWhitelist returns the discriminator for Solana gateway 'whitelist_spl_mint' instruction
	DiscriminatorWhitelistSplMint = idlgateway.IDLGateway.GetDiscriminator("whitelist_spl_mint")
)

// ParseGatewayWithPDA parses the gateway id and program derived address from the given string
func ParseGatewayWithPDA(gatewayAddress string) (solana.PublicKey, solana.PublicKey, error) {
	var gatewayID, pda solana.PublicKey

	// decode gateway address
	gatewayID, err := solana.PublicKeyFromBase58(gatewayAddress)
	if err != nil {
		return gatewayID, pda, errors.Wrap(err, "unable to decode address")
	}

	// compute gateway PDA
	seed := []byte(PDASeed)
	pda, _, err = solana.FindProgramAddress([][]byte{seed}, gatewayID)

	return gatewayID, pda, err
}

// ComputePdaAddress computes the PDA address for the custom program PDA with provided seed
func ComputePdaAddress(connected solana.PublicKey, seed []byte) (solana.PublicKey, error) {
	pdaComputed, _, err := solana.FindProgramAddress([][]byte{seed}, connected)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return pdaComputed, nil
}

// ComputeConnectedPdaAddress computes the PDA address for the custom program PDA with seed "connected"
func ComputeConnectedPdaAddress(connected solana.PublicKey) (solana.PublicKey, error) {
	return ComputePdaAddress(connected, []byte("connected"))
}
