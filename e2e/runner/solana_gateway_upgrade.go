package runner

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os/exec"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
)

// SolanaVerifyGatewayContractsUpgrade upgrades the Solana contracts and verifies the upgrade
func (r *E2ERunner) SolanaVerifyGatewayContractsUpgrade(deployerPrivateKey string) {
	r.Logger.Print("🏃 Upgrading Solana gateway contracts")

	pdaComputed := r.ComputePdaAddress()
	pdaInfo, err := r.SolanaClient.GetAccountInfoWithOpts(r.Ctx, pdaComputed, &rpc.GetAccountInfoOpts{
		Commitment: rpc.CommitmentConfirmed,
	})
	require.NoError(r, err)

	// deserialize the PDA info
	pdaDataBefore := solanacontracts.PdaInfo{}
	err = borsh.Deserialize(&pdaDataBefore, pdaInfo.Bytes())
	require.NoError(r, err)

	err = triggerSolanaUpgrade()
	require.NoError(r, err, "failed to trigger Solana upgrade")
	r.Logger.Print("⚙️ Solana upgrade completed")

	pdaInfo, err = r.SolanaClient.GetAccountInfoWithOpts(r.Ctx, pdaComputed, &rpc.GetAccountInfoOpts{
		Commitment: rpc.CommitmentConfirmed,
	})
	require.NoError(r, err)

	// deserialize the PDA info
	pdaDataAfter := solanacontracts.PdaInfo{}
	err = borsh.Deserialize(&pdaDataAfter, pdaInfo.Bytes())
	require.NoError(r, err)

	// Verify that data does not change
	require.Equal(r, pdaDataBefore.Nonce, pdaDataAfter.Nonce)
	require.Equal(
		r,
		ethcommon.BytesToAddress(pdaDataBefore.TssAddress[:]),
		ethcommon.BytesToAddress(pdaDataAfter.TssAddress[:]),
	)
	require.Equal(r, pdaDataBefore.Authority, pdaDataAfter.Authority)
	require.Equal(r, pdaDataBefore.ChainID, pdaDataAfter.ChainID)
	require.Equal(r, pdaDataBefore.DepositPaused, pdaDataAfter.DepositPaused)

	r.VerifyUpgradedInstruction(deployerPrivateKey)
}

func (r *E2ERunner) VerifyUpgradedInstruction(deployerPrivateKey string) {
	privkey, err := solana.PrivateKeyFromBase58(deployerPrivateKey)
	require.NoError(r, err)
	// Calculate the instruction discriminator for "upgraded"
	// Anchor uses the first 8 bytes of the sha256 hash of "global:upgraded"
	// Manually generating the discriminator as there is just one extra function in the new program
	discriminator := getAnchorDiscriminator("upgraded")
	// Build instruction
	data := append(discriminator, []byte{}...)

	var instConnected solana.GenericInstruction
	accountSliceConnected := make([]*solana.AccountMeta, 0, 1)
	accountSliceConnected = append(accountSliceConnected, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
	instConnected.ProgID = r.GatewayProgram
	instConnected.AccountValues = accountSliceConnected
	instConnected.DataBytes = data

	// create and sign the transaction
	signedTx := r.CreateSignedTransaction([]solana.Instruction{&instConnected}, privkey, []solana.PrivateKey{})

	// broadcast the transaction and wait for finalization
	_, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("upgrade  logs: %v", out.Meta.LogMessages)

	decoded, err := base64.StdEncoding.DecodeString(out.Meta.ReturnData.Data.String())
	require.NoError(r, err)
	require.True(r, decoded[0] == 1)
}

func getAnchorDiscriminator(methodName string) []byte {
	// In Anchor, the namespace is "global" + method name
	namespace := fmt.Sprintf("global:%s", methodName)

	// Calculate SHA256
	hash := sha256.Sum256([]byte(namespace))

	// Return first 8 bytes
	return hash[:8]
}

// triggerSolanaUpgrade triggers the Solana upgrade by creating a file `execute-update` on the Solana container
// The shell script on the Solana container will remove the file after completing the upgrade
// Refer: contrib/localnet/solana/start-solana.sh
func triggerSolanaUpgrade() error {
	// Create the execute-update file on Solana container
	createCmd := exec.Command("ssh", "root@solana", "touch", "/data/execute-update")
	if err := createCmd.Run(); err != nil {
		return fmt.Errorf("failed to create execute-update file: %w", err)
	}

	// Start checking for file removal with timeout
	timeout := time.After(2 * time.Minute)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for Solana upgrade to complete")

		case <-ticker.C:
			// Check if file still exists
			checkCmd := exec.Command("ssh", "root@solana", "test", "-f", "/data/execute-update")
			if err := checkCmd.Run(); err != nil {
				// If the command fails, it means the file doesn't exist (upgrade completed)
				return nil
			}
		}
	}
}
