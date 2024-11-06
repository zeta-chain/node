package e2etests

import (
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
)

func TestSolanaDepositSPL(r *runner.E2ERunner, _ []string) {
	// require.Len(r, args, 1)

	// load deployer private key
	privKey, err := solana.PrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	require.NoError(r, err)

	wall, err := solana.WalletFromPrivateKeyBase58(r.SPLPrivateKey.String())
	require.NoError(r, err)
	r.DepositSPL(&privKey, *wall)
}
