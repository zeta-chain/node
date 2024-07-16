package runner

import (
	"time"

	"github.com/gagliardetto/solana-go"
)

func (r *E2ERunner) SetupSolanaAccount() {
	r.Logger.Print("⚙️ setting up Solana account")
	startTime := time.Now()
	defer func() {
		r.Logger.Print("✅ Solana account setup in %s", time.Since(startTime))
	}()

	r.SetSolanaAddress()
}

// SetSolanaAddress imports the deployer's private key
func (r *E2ERunner) SetSolanaAddress() {
	privateKey := solana.MustPrivateKeyFromBase58(r.Account.RawBase58PrivateKey.String())
	r.SolanaDeployerAddress = privateKey.PublicKey()

	r.Logger.Info("SolanaDeployerAddress: %s", r.SolanaDeployerAddress)
}
