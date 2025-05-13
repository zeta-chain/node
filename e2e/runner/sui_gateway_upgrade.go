package runner

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"regexp"
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

// reGatewayPackageID is the regex to extract the gateway package ID from the output
//
// │ Published Objects:                                                                               │
// │  ┌──                                                                                             │
// │  │ PackageID: 0x67500cf6e39f3d8937c4f15a298da72358410c84357ee33e0030c13872f5339e                 │
// │  │ Version: 2                                                                                    │
// │  │ Digest: 9sLQJdSZcQp8jHVJenBGKjipxspLKZ1P4QBnTxffA4Gb                                          │
// │  │ Modules: evm, gateway                                                                         │
// │  └──                                                                                             │
// ╰──────────────────────────────────────────────────────────────────────────────────────────────────╯
var reGatewayPackageID = regexp.MustCompile(`│\s*PackageID: *(0x[0-9a-fA-F]+)\s*│`)

// SuiVerifyGatewayPackageUpgrade upgrades the Sui gateway package and verifies the upgrade
func (r *E2ERunner) SuiVerifyGatewayPackageUpgrade() {
	// retrieve original gateway object data
	gatewayDataBefore, err := r.suiGetObjectData(r.Ctx, r.SuiGateway.ObjectID())
	require.NoError(r, err)

	// upgrade the Sui gateway package
	newGatewayPackageID, err := r.upgradeSuiGatewayPackage()
	require.NoError(r, err)

	// call the new method 'upgraded' in the new gateway package
	r.moveCallUpgraded(r.Ctx, newGatewayPackageID)

	// retrieve new gateway object data
	gatewayDataAfter, err := r.suiGetObjectData(r.Ctx, r.SuiGateway.ObjectID())
	require.NoError(r, err)

	// gateway data should remain unchanged
	require.Equal(r, gatewayDataBefore, gatewayDataAfter)
}

// upgradeSuiGatewayPackage upgrades the Sui gateway package by deploying new compiled gateway package
func (r *E2ERunner) upgradeSuiGatewayPackage() (string, error) {
	// construct the CLI command for package upgrade
	// #nosec G204, inputs are controlled in E2E test
	cmdUpgrade := exec.Command("sui", []string{
		"client",
		"upgrade",
		"--skip-dependency-verification",
		"--upgrade-capability",
		r.SuiGatewayUpgradeCap,
	}...)
	cmdUpgrade.Dir = suiGatewayUpgradedPath

	// run command and show output
	startTime := time.Now()
	output, err := cmdUpgrade.Output()
	require.NoError(r, err)

	r.Logger.Info("sui gateway package upgrade took %f seconds: \n%s", time.Since(startTime).Seconds(), string(output))

	// scan through the output line by line to find the new gateway package ID
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		matches := reGatewayPackageID.FindStringSubmatch(line)
		if len(matches) >= 2 {
			// return the first capture group which contains the PackageID
			return matches[1], nil
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
	err = os.WriteFile(moveTomlPath, []byte(updatedContent), 0644)
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
