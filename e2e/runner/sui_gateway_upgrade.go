package runner

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const (
	// suiGatewayUpgradedPath is the path to the upgraded Sui gateway package
	suiGatewayUpgradedPath = "/work/protocol-contracts-sui-upgrade"
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
	// build the CLI command for package upgrade
	cmdBuild := exec.Command("sui", "move", "build")
	cmdBuild.Dir = suiGatewayUpgradedPath
	require.NoError(r, cmdBuild.Run(), "unable to build sui gateway package")

	// construct the CLI command for package upgrade
	// #nosec G204, inputs are controlled in E2E test
	cmdUpgrade := exec.Command("sui", []string{
		"client",
		"upgrade",
		"--json", // output in JSON format for easier parsing
		"--skip-dependency-verification",
		"--upgrade-capability",
		r.SuiGatewayUpgradeCap,
	}...)
	cmdUpgrade.Dir = suiGatewayUpgradedPath

	// run command and show output
	startTime := time.Now()
	output, err := cmdUpgrade.Output()
	require.NoError(r, err)

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
		Module:          r.SuiGateway.Module(),
		Function:        "upgraded",
		TypeArguments:   []any{},
		Arguments:       []any{},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)

	r.suiExecuteTx(signer, tx)
}

// suiPatchMoveConfig updates the 'published-at' field in the 'Move.toml' file
// with the original published gateway package ID
func (r *E2ERunner) suiPatchMoveConfig() {
	const moveTomlPath = suiGatewayUpgradedPath + "/Move.toml"

	// read the entire Move.toml file
	content, err := os.ReadFile(moveTomlPath)
	require.NoError(r, err, "unable to read Move.toml")
	contentStr := string(content)

	// Replace the placeholder with the actual published gateway package ID
	publishedAt := r.SuiGateway.PackageID()
	updatedContent := strings.Replace(contentStr, "ORIGINAL-PACKAGE-ID", publishedAt, 1)

	// Write the updated content back to the file
	err = os.WriteFile(moveTomlPath, []byte(updatedContent), 0600)
	require.NoError(r, err, "unable to write to Move.toml")
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
