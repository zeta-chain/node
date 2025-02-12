package runner

import (
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"

	suicontract "github.com/zeta-chain/node/e2e/contracts/sui"
)

func (r *E2ERunner) SetupSui(faucetURL string) {
	r.Logger.Print("⚙️ initializing gateway package on Sui")

	deployerSigner, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	deployerAddress := deployerSigner.Address()

	header := map[string]string{}
	err = sui.RequestSuiFromFaucet(faucetURL, deployerAddress, header)
	require.NoError(r, err, "sui faucet request to %s", faucetURL)

	client := r.Clients.Sui

	publishReq, err := client.Publish(r.Ctx, models.PublishRequest{
		Sender:          deployerAddress,
		CompiledModules: []string{suicontract.GatewayBytecodeBase64()},
		Dependencies: []string{
			"0x0000000000000000000000000000000000000000000000000000000000000001",
			"0x0000000000000000000000000000000000000000000000000000000000000002",
		},
		GasBudget: "5000000000",
	})
	require.NoError(r, err, "create publish tx")

	signature, err := deployerSigner.SignTransactionBlock(publishReq.TxBytes)
	require.NoError(r, err, "sign transaction")

	resp, err := client.SuiExecuteTransactionBlock(r.Ctx, models.SuiExecuteTransactionBlockRequest{
		TxBytes:   publishReq.TxBytes,
		Signature: []string{signature},
		Options: models.SuiTransactionBlockOptions{
			ShowEffects:        true,
			ShowBalanceChanges: true,
			ShowEvents:         true,
			ShowObjectChanges:  true,
		},
		RequestType: "WaitForLocalExecution",
	})
	require.NoError(r, err)

	var packageID, objectID string
	for _, change := range resp.ObjectChanges {
		if change.Type == "published" {
			packageID = change.PackageId
		}
	}
	require.NotEmpty(r, packageID, "find packageID")

	gatewayType := fmt.Sprintf("%s::gateway::Gateway", packageID)
	for _, change := range resp.ObjectChanges {
		if change.Type == "created" && change.ObjectType == gatewayType {
			objectID = change.ObjectId
		}
	}
	require.NotEmpty(r, objectID, "find objectID")

	// TODO: save IDs in config and configure chain
}
