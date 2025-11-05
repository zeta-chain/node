package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// SuiVerifyGatewayPackageUpgrade upgrades the Sui gateway package and verifies the upgrade
func (r *E2ERunner) SuiVerifyGatewayPackageUpgrade() {
	r.Logger.Print("🏃 Upgrading Sui gateway package")

	// retrieve original gateway object data
	gatewayDataBefore, err := r.suiGetObjectData(r.Ctx, r.SuiGateway.ObjectID())
	require.NoError(r, err)

	// upgrade the Sui gateway package
	r.suiUpgradeGatewayPackage()

	r.Logger.Print("⚙️ Sui gateway upgrade completed: %s", r.SuiGateway.PackageID())

	// call the new method 'upgraded' in the new gateway package
	r.moveCallUpgraded(r.Ctx, r.SuiGateway.PackageID())

	// retrieve new gateway object data
	gatewayDataAfter, err := r.suiGetObjectData(r.Ctx, r.SuiGateway.ObjectID())
	require.NoError(r, err)

	// gateway data should remain unchanged
	require.Equal(r, gatewayDataBefore, gatewayDataAfter)

	// deposit from new gateway package should be observed
	r.Logger.Print("🏃 Verifying Sui deposit from new package")
	r.suiVerifyDepositFromPackage(r.SuiGateway.PackageID(), big.NewInt(10000000000))

	// deposit from previous gateway package should be observed
	r.Logger.Print("🏃 Verifying Sui deposit from previous package")
	r.suiVerifyDepositFromPackage(r.SuiGateway.Previous().PackageID(), big.NewInt(2000000))

	// deprecate previous gateway package
	previousPackageID := r.SuiGateway.Previous().PackageID()
	r.Logger.Print("🏃 Deprecating previous package %s", previousPackageID)
	r.suiDeprecatePreviousPackage()

	// deposit from deprecated package should not be observed
	r.Logger.Print("🏃 Verifying Sui deposit from deprecated package")
	r.suiVerifyDepositFromDeprecatePackage(previousPackageID, big.NewInt(2000000))
}

// suiUpgradeGatewayPackage upgrades the Sui gateway package by deploying new compiled gateway package
func (r *E2ERunner) suiUpgradeGatewayPackage() {
	// build the upgraded gateway package
	r.suiBuildGatewayUpgraded()

	r.Logger.Print("gateway package ID v1: %s", r.SuiGateway.PackageID())

	// construct the CLI command for package upgrade
	// #nosec G204, inputs are controlled in E2E test
	cmdUpgrade := exec.Command("sui", []string{
		"client",
		"upgrade",
		"--json", // output in JSON format for easier parsing
		"--upgrade-capability",
		r.SuiGatewayUpgradeCap,
	}...)
	cmdUpgrade.Dir = r.WorkDirPrefixed(suiGatewayUpgradedPathV2)

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
	packageID := ""
	for _, change := range response.ObjectChanges {
		if change.Type == "published" {
			packageID = change.PackageId
		}
	}
	require.NotEmpty(r, packageID, "new gateway package ID not found")

	// replace v1 package ID with v2 package ID in the 'published-at' field
	publishedAtOld := fmt.Sprintf(`published-at = "%s"`, r.SuiGateway.PackageID())
	publishedAtNew := fmt.Sprintf(`published-at = "%s"`, packageID)
	r.suiPatchMoveConfig(r.WorkDirPrefixed(suiGatewayUpgradedPathV2), publishedAtOld, publishedAtNew)

	// find withdraw cap ID
	withdrawCapID, found := r.suiGetOwnedObjectID(r.SuiTSSAddress, r.SuiGateway.WithdrawCapType())
	require.True(r, found, "withdraw cap object not found")

	// update runner gateway package ID
	previousPackageID := r.SuiGateway.PackageID()
	originalPackageID := previousPackageID
	r.SuiGateway, err = sui.NewGatewayFromPairID(
		sui.MakePairID(packageID, r.SuiGateway.ObjectID(), withdrawCapID, previousPackageID, originalPackageID),
	)
	require.NoError(r, err)

	r.Logger.Print("gateway package ID v2: %s", r.SuiGateway.PackageID())

	// update the chain params
	err = r.setSuiChainParams(false)
	require.NoError(r, err)

	// wait 2 Zeta blocks to ensure zetaclient picks up the new chain params
	utils.WaitForZetaBlocks(r.Ctx, r, r.ZEVMClient, 2, 10*time.Second)
}

// suiDeprecatePreviousPackage deprecates the previous Sui gateway package
func (r *E2ERunner) suiDeprecatePreviousPackage() {
	// find withdraw cap ID
	withdrawCapID, found := r.suiGetOwnedObjectID(r.SuiTSSAddress, r.SuiGateway.WithdrawCapType())
	require.True(r, found, "withdraw cap object not found")

	// deprecate the previous package by setting it to empty
	var (
		err               error
		packageID         = r.SuiGateway.PackageID()
		gatewayID         = r.SuiGateway.ObjectID()
		originalPackageID = r.SuiGateway.Original().PackageID()
	)
	r.SuiGateway, err = sui.NewGatewayFromPairID(
		sui.MakePairID(packageID, gatewayID, withdrawCapID, "", originalPackageID),
	)
	require.NoError(r, err)

	// update the chain params
	err = r.setSuiChainParams(false)
	require.NoError(r, err)

	// wait 2 Zeta blocks to ensure zetaclient picks up the new chain params
	utils.WaitForZetaBlocks(r.Ctx, r, r.ZEVMClient, 2, 10*time.Second)
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

// suiVerifyDepositFromPackage verifies the deposit from given Sui gateway packageID and amount
func (r *E2ERunner) suiVerifyDepositFromPackage(packageID string, amount *big.Int) {
	oldBalance, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	// make a deposit from gateway package
	resp := r.SuiDepositSUI(packageID, r.EVMAddress(), math.NewUintFromBigInt(amount))
	r.Logger.Info("Sui deposit tx: %s from package: %s", resp.Digest, packageID)

	// wait for the CCTX to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, resp.Digest, r.CctxClient, r.Logger, r.CctxTimeout)
	require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// wait for the SUI ZRC20 balance to be updated
	change := utils.NewExactChange(amount)
	utils.WaitAndVerifyZRC20BalanceChange(r, r.SUIZRC20, r.EVMAddress(), oldBalance, change, r.Logger)
	require.Equal(r, amount.Uint64(), cctx.InboundParams.Amount.Uint64())

	// only one single CCTX should be created
	cctxs := utils.GetCCTXByInboundHash(r.Ctx, r.CctxClient, resp.Digest)
	require.Len(r, cctxs, 1)
}

// suiVerifyDepositFromDeprecatePackage verifies the deposit from the deprecated Sui gateway package
func (r *E2ERunner) suiVerifyDepositFromDeprecatePackage(packageID string, amount *big.Int) {
	// make a deposit from gateway package
	resp := r.SuiDepositSUI(packageID, r.EVMAddress(), math.NewUintFromBigInt(amount))
	r.Logger.Info("Sui deposit tx: %s from deprecated package: %s", resp.Digest, packageID)

	// wait for 2 zeta blocks
	utils.WaitForZetaBlocks(r.Ctx, r, r.ZEVMClient, 2, 20*time.Second)

	// query cctx by inbound hash, no CCTX should be created
	in := &crosschaintypes.QueryInboundHashToCctxDataRequest{InboundHash: resp.Digest}
	_, err := r.CctxClient.InboundHashToCctxData(r.Ctx, in)
	require.ErrorIs(r, err, status.Error(codes.NotFound, "not found"))
}
