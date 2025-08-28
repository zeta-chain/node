package runner

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/contracts/sui"
)

// SuiVerifyGatewayPackageUpgrade upgrades the Sui gateway package and verifies the upgrade
func (r *E2ERunner) SuiVerifyGatewayPackageUpgrade() {
	r.Logger.Print("🏃 Upgrading Sui gateway package")

	// retrieve original gateway object data
	gatewayDataBefore, err := r.suiGetObjectData(r.Ctx, r.SuiGateway.ObjectID())
	require.NoError(r, err)

	// upgrade the Sui gateway package
	newGatewayPackageID, err := r.suiUpgradeGatewayPackage()
	require.NoError(r, err)

	r.Logger.Print("⚙️ Sui gateway package upgrade completed")

	// call the new method 'upgraded' in the new gateway package
	r.moveCallUpgraded(r.Ctx, newGatewayPackageID)

	// retrieve new gateway object data
	gatewayDataAfter, err := r.suiGetObjectData(r.Ctx, r.SuiGateway.ObjectID())
	require.NoError(r, err)

	// gateway data should remain unchanged
	require.Equal(r, gatewayDataBefore, gatewayDataAfter)
}

// suiUpgradeGatewayPackage upgrades the Sui gateway package by deploying new compiled gateway package
func (r *E2ERunner) suiUpgradeGatewayPackage() (packageID string, err error) {
	// build the upgraded gateway package
	r.suiBuildGatewayUpgraded()

	// construct the CLI command for package upgrade
	// #nosec G204, inputs are controlled in E2E test
	cmdUpgrade := exec.Command("sui", []string{
		"client",
		"upgrade",
		"--json", // output in JSON format for easier parsing
		"--upgrade-capability",
		r.SuiGatewayUpgradeCap,
	}...)
	cmdUpgrade.Dir = r.WorkDirPrefixed(suiGatewayUpgradedPath)

	// run command and show output
	startTime := time.Now()
	output, err := cmdUpgrade.Output()
	require.NoError(r, err, "Sui upgrade gateway package failed: \n%s", string(output))

	r.Logger.Info("Sui gateway package upgrade took %f seconds: \n%s", time.Since(startTime).Seconds(), string(output))

	// convert output to transaction block response struct
	response := &models.SuiTransactionBlockResponse{}
	err = json.Unmarshal(output, response)
	require.NoError(r, err)

	// find packageID
	for _, change := range response.ObjectChanges {
		if change.Type == "published" {
			return change.PackageId, nil
		}
	}

	return "", errors.New("new gateway package ID not found")
}

// moveCallUpgraded performs a move call to 'upgraded' method on the new Sui gateway package
func (r *E2ERunner) moveCallUpgraded(ctx context.Context, gatewayPackageID string) {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "unable to get deployer signer")

	tx, err := r.Clients.Sui.MoveCall(ctx, models.MoveCallRequest{
		Signer:          signer.Address(),
		PackageObjectId: gatewayPackageID,
		Module:          sui.GatewayModule,
		Function:        "upgraded",
		TypeArguments:   []any{},
		Arguments:       []any{},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)

	r.suiExecuteTx(signer, tx)
}

// suiPatchMoveConfig replaces the given 'text' in the 'Move.toml' file with the given 'value'
func (r *E2ERunner) suiPatchMoveConfig(path, text, value string) {
	moveTomlFile := filepath.Join(path, "Move.toml")

	// read the entire Move.toml file
	// #nosec G304 -- this is a config file for example package
	content, err := os.ReadFile(moveTomlFile)
	require.NoError(r, err, "unable to read "+moveTomlFile)
	contentStr := string(content)

	// replace the text with the specified value
	updatedContent := strings.Replace(contentStr, text, value, 1)

	// write the updated content back to the file
	err = os.WriteFile(moveTomlFile, []byte(updatedContent), 0600)
	require.NoError(r, err, "unable to write to "+moveTomlFile)
}

// suiGetObjectData retrieves the object data for the given object ID
func (r *E2ERunner) suiGetObjectData(ctx context.Context, objectID string) (models.SuiParsedData, error) {
	object, err := r.Clients.Sui.SuiGetObject(ctx, models.SuiGetObjectRequest{
		ObjectId: objectID,
		Options:  models.SuiObjectDataOptions{ShowContent: true},
	})
	require.NoError(r, err)
	require.NotNil(r, object.Data)
	require.NotNil(r, object.Data.Content)

	return *object.Data.Content, nil
}
