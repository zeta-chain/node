package runner

import (
	"time"

	"github.com/gagliardetto/solana-go"
)

func (r *E2ERunner) SetupSolanaAccount() {
	r.Logger.Print("⚙️ setting up Solana account")
	startTime := time.Now()
	defer func() {
		r.Logger.Print("✅ Solana account setup in %s\n", time.Since(startTime))
	}()

	r.SetSolanaAddress()
}

// SetSolanaAddress imports the deployer's private key
func (r *E2ERunner) SetSolanaAddress() {
	privateKey := solana.MustPrivateKeyFromBase58(r.Account.RawPrivateKey.String())
	r.SolanaDeployerAddress = privateKey.PublicKey()

	r.Logger.Info("SolanaDeployerAddress: %s", r.BTCDeployerAddress.EncodeAddress())
}
