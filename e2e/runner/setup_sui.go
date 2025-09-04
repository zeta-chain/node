package runner

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/zrc20.sol"

	"github.com/zeta-chain/node/e2e/config"
	suibin "github.com/zeta-chain/node/e2e/contracts/sui/bin"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	zetasui "github.com/zeta-chain/node/pkg/contracts/sui"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	suiclient "github.com/zeta-chain/node/zetaclient/chains/sui/client"
)

const (
	// changeTypeCreated is the type of change that indicates a new object was created
	changeTypeCreated = "created"

	// suiExamplePath is the path to the example package
	suiExamplePath = "example"

	// suiGatewayUpgradedPath is the path to the upgraded Sui gateway package
	suiGatewayUpgradedPath = "protocol-contracts-sui-upgrade"
)

var (
	// suiExampleBinToken is the path to the example token binary file
	suiExampleBinToken = fmt.Sprintf("%s/build/example/bytecode_modules/token.mv", suiExamplePath)

	// suiExampleBinConnected is the path to the example connected binary file
	suiExampleBinConnected = fmt.Sprintf("%s/build/example/bytecode_modules/connected.mv", suiExamplePath)
)

// RequestSuiFromFaucet requests SUI tokens from the faucet for the runner account
func (r *E2ERunner) RequestSuiFromFaucet(faucetURL, recipient string) {
	header := map[string]string{}
	err := sui.RequestSuiFromFaucet(faucetURL, recipient, header)
	require.NoError(r, err, "sui faucet request to %s", faucetURL)
}

// SetupSui initializes the gateway package on Sui and initialize the chain params on ZetaChain
func (r *E2ERunner) SetupSui(faucetURL string) {
	r.Logger.Print("⚙️ initializing gateway package on Sui")

	deployerSigner, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	deployerAddress := deployerSigner.Address()

	// fund deployer
	r.RequestSuiFromFaucet(faucetURL, deployerAddress)

	// import deployer private key and select it as active address
	r.suiSetupDeployerAccount()

	// fund the TSS
	// request twice from the faucet to ensure TSS has enough funds for the first withdraw
	// TODO: this step might no longer necessary if a custom solution is implemented for the TSS funding
	r.RequestSuiFromFaucet(faucetURL, r.SuiTSSAddress)
	r.RequestSuiFromFaucet(faucetURL, r.SuiTSSAddress)

	// deploy gateway package
	whitelistCapID, withdrawCapID := r.suiDeployGateway()

	// issue message context
	messageContextID := r.issueMessageContext()

	// deploy SUI zrc20
	r.deploySUIZRC20()

	// deploy fake USDC and whitelist it
	fakeUSDCCoinType := r.suiDeployFakeUSDC()
	r.whitelistSuiFakeUSDC(deployerSigner, fakeUSDCCoinType, whitelistCapID)

	// build and deploy example package with on_call function
	r.suiBuildExample()
	r.suiDeployExample(
		&r.SuiExample,
		suibin.ReadMoveBinaryBase64(r, r.WorkDirPrefixed(suiExampleBinToken)),
		suibin.ReadMoveBinaryBase64(r, r.WorkDirPrefixed(suiExampleBinConnected)),
		[]string{r.SuiGateway.PackageID()},
	)

	// send withdraw cap to TSS
	r.suiTransferObjectToTSS(deployerSigner, withdrawCapID)

	// send message context to TSS
	r.suiTransferObjectToTSS(deployerSigner, messageContextID)

	// set the chain params
	err = r.setSuiChainParams()
	require.NoError(r, err)
}

// suiSetupDeployerAccount imports a Sui deployer private key using the sui keytool import command
// and sets the deployer address as the active address.
func (r *E2ERunner) suiSetupDeployerAccount() {
	deployerSigner, err := r.Account.SuiSigner()
	require.NoError(r, err, "unable to get deployer signer")

	var (
		deployerAddress    = deployerSigner.Address()
		deployerPrivKeyHex = r.Account.RawPrivateKey.String()
	)

	// convert private key to bech32
	deployerPrivKeySecp256k1, err := zetasui.PrivateKeyBech32Secp256k1FromHex(deployerPrivKeyHex)
	require.NoError(r, err)

	// import deployer private key using sui keytool import
	// #nosec G204, inputs are controlled in E2E test
	cmdImport := exec.Command("sui", "keytool", "import", deployerPrivKeySecp256k1, "secp256k1")
	require.NoError(r, cmdImport.Run(), "unable to import sui deployer private key")

	// switch to deployer address using sui client switch
	// #nosec G204, inputs are controlled in E2E test
	cmdSwitch := exec.Command("sui", "client", "switch", "--address", deployerAddress)
	require.NoError(r, cmdSwitch.Run(), "unable to switch to deployer address")

	// ensure the deployer address is active
	// #nosec G204, inputs are controlled in E2E test
	cmdList := exec.Command("sui", "client", "active-address")
	output, err := cmdList.Output()
	require.NoError(r, err)
	require.Equal(r, deployerAddress, strings.TrimSpace(string(output)))
}

// suiBuildPackage builds the Sui package under the given path using CLI command
func (r *E2ERunner) suiBuildPackage(path string) {
	tStart := time.Now()

	// build the CLI command for package upgrade
	cmdBuild := exec.Command("sui", "move", "build")
	cmdBuild.Dir = path

	// run command and show output
	output, err := cmdBuild.CombinedOutput()
	r.Logger.Info("sui move build output for: %s\n%s", path, string(output))
	require.NoError(r, err, "unable to build sui package: \n%s", string(output))

	r.Logger.Info("sui package build took %f seconds", time.Since(tStart).Seconds())
}

// suiBuildGatewayUpgraded builds the upgraded gateway package
func (r *E2ERunner) suiBuildGatewayUpgraded() {
	// in order to upgrade the gateway package, we need e patches to the Move.toml files:
	// 1. set the `published-at` so that SUI knows which deployed gateway package the upgrade applies to.
	// 2. use `gateway = 0x0` as a placeholder that will be replaced by new gateway package address.
	// 3. set the old placeholder "ORIGINAL-PACKAGE-ID" to actual package ID, will deprecate it in the future.
	publishedAt := fmt.Sprintf(`published-at = "%s"`, r.SuiGateway.PackageID())
	gatewayAddress := fmt.Sprintf(`gateway = "%s"`, r.SuiGateway.PackageID())
	r.suiPatchMoveConfig(r.WorkDirPrefixed(suiGatewayUpgradedPath), `published-at = "0x0"`, publishedAt)
	r.suiPatchMoveConfig(r.WorkDirPrefixed(suiGatewayUpgradedPath), gatewayAddress, `gateway = "0x0"`)
	r.suiPatchMoveConfig(r.WorkDirPrefixed(suiGatewayUpgradedPath), `published-at = "ORIGINAL-PACKAGE-ID"`, publishedAt)

	// build the upgraded gateway package
	r.suiBuildPackage(r.WorkDirPrefixed(suiGatewayUpgradedPath))
}

// suiBuildExample builds the example package
func (r *E2ERunner) suiBuildExample() {
	// in order to import the gateway package, we need 3 patches to the Move.toml files:
	// 1. set the actual gateway address in the gateway package, otherwise the build will fail
	// 2. set the actual gateway address to `published-at` in the gateway package, otherwise the deploy will fail
	// 3. set the actual gateway address in the example package, otherwise the build will fail.
	publishedAt := fmt.Sprintf(`published-at = "%s"`, r.SuiGateway.PackageID())
	gatewayAddress := fmt.Sprintf(`gateway = "%s"`, r.SuiGateway.PackageID())
	r.suiPatchMoveConfig(r.WorkDirPrefixed(suiGatewayUpgradedPath), `published-at = "0x0"`, publishedAt)
	r.suiPatchMoveConfig(r.WorkDirPrefixed(suiGatewayUpgradedPath), `gateway = "0x0"`, gatewayAddress)
	r.suiPatchMoveConfig(r.WorkDirPrefixed(suiExamplePath), `gateway = "0x0"`, gatewayAddress)

	r.suiBuildPackage(r.WorkDirPrefixed(suiExamplePath))
}

// suiDeployGateway deploys the SUI gateway package on Sui
func (r *E2ERunner) suiDeployGateway() (whitelistCapID, withdrawCapID string) {
	const (
		filterGatewayType      = "gateway::Gateway"
		filterWithdrawCapType  = "gateway::WithdrawCap"
		filterWhitelistCapType = "gateway::WhitelistCap"
		filterUpgradeCapType   = "0x2::package::UpgradeCap"
	)

	objectTypeFilters := []string{
		filterGatewayType,
		filterWhitelistCapType,
		filterWithdrawCapType,
		filterUpgradeCapType,
	}
	packageID, objectIDs := r.suiDeployPackage(
		[]string{suibin.GatewayBytecodeBase64(), suibin.EVMBytecodeBase64()},
		[]string{},
		objectTypeFilters,
	)

	gatewayID, ok := objectIDs[filterGatewayType]
	require.True(r, ok, "gateway object not found")

	whitelistCapID, ok = objectIDs[filterWhitelistCapType]
	require.True(r, ok, "whitelistCap object not found")

	withdrawCapID, ok = objectIDs[filterWithdrawCapType]
	require.True(r, ok, "withdrawCap object not found")

	r.SuiGatewayUpgradeCap, ok = objectIDs[filterUpgradeCapType]
	require.True(r, ok, "upgradeCap object not found")

	// set sui gateway
	r.SuiGateway = zetasui.NewGateway(packageID, gatewayID)

	return whitelistCapID, withdrawCapID
}

// deploySUIZRC20 deploys the SUI zrc20 on ZetaChain
func (r *E2ERunner) deploySUIZRC20() {
	// send message to deploy SUI zrc20
	liqCap := math.NewUint(10e18)
	adminAddr := r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName)
	_, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, fungibletypes.NewMsgDeployFungibleCoinZRC20(
		adminAddr,
		"",
		chains.SuiLocalnet.ChainId,
		9,
		"SUI",
		"SUI",
		coin.CoinType_Gas,
		10000,
		&liqCap,
	))
	require.NoError(r, err)

	// set the address in the store
	r.SetupSUIZRC20()
}

// suiDeployFakeUSDC deploys the FakeUSDC contract on Sui
// it returns the treasuryCap object ID that allows to mint tokens
func (r *E2ERunner) suiDeployFakeUSDC() string {
	packageID, objectIDs := r.suiDeployPackage(
		[]string{suibin.FakeUSDCBytecodeBase64()},
		[]string{},
		[]string{"TreasuryCap"},
	)

	treasuryCap, ok := objectIDs["TreasuryCap"]
	require.True(r, ok, "treasuryCap not found")

	coinType := packageID + "::fake_usdc::FAKE_USDC"

	// strip 0x from packageID
	coinType = coinType[2:]

	// set asset value for zrc20 and treasuryCap object ID
	r.SuiTokenCoinType = coinType
	r.SuiTokenTreasuryCap = treasuryCap

	return coinType
}

// suiDeployExample deploys the example package on Sui
func (r *E2ERunner) suiDeployExample(
	example *config.SuiExample,
	fungibleTokenBytecodeBase64, connectedBytecodeBase64 string,
	extraDependencies []string,
) {
	const (
		filterGlobalConfigType = "connected::GlobalConfig"
		filterPartnerType      = "connected::Partner"
		filterClockType        = "connected::Clock"
	)

	objectTypeFilters := []string{filterGlobalConfigType, filterPartnerType, filterClockType}
	packageID, objectIDs := r.suiDeployPackage(
		[]string{fungibleTokenBytecodeBase64, connectedBytecodeBase64},
		extraDependencies,
		objectTypeFilters,
	)
	r.Logger.Info("deployed example package with packageID: %s", packageID)

	globalConfigID, ok := objectIDs[filterGlobalConfigType]
	require.True(r, ok, "globalConfig object not found")

	partnerID, ok := objectIDs[filterPartnerType]
	require.True(r, ok, "partner object not found")

	clockID, ok := objectIDs[filterClockType]
	require.True(r, ok, "clock object not found")

	// save the example package info
	*example = config.SuiExample{
		PackageID:      config.DoubleQuotedString(packageID),
		TokenType:      config.DoubleQuotedString(packageID + "::token::TOKEN"),
		GlobalConfigID: config.DoubleQuotedString(globalConfigID),
		PartnerID:      config.DoubleQuotedString(partnerID),
		ClockID:        config.DoubleQuotedString(clockID),
	}
}

// suiDeployPackage is a helper function that deploys a package on Sui
// It returns the packageID and a map of object types to their IDs
func (r *E2ERunner) suiDeployPackage(
	bytecodeBase64s []string,
	extraDependencies []string,
	objectTypeFilters []string,
) (string, map[string]string) {
	client := r.Clients.Sui

	deployerSigner, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	deployerAddress := deployerSigner.Address()

	// aside from the standard framework dependencies, add extra dependencies if provided
	dependencies := append([]string{
		"0x1", // Sui Framework
		"0x2", // Move Standard Library
	}, extraDependencies...) // other dependencies

	// build the publish transaction and sign it with deployer key
	publishTx, err := client.Publish(r.Ctx, models.PublishRequest{
		Sender:          deployerAddress,
		CompiledModules: bytecodeBase64s,
		Dependencies:    dependencies,
		GasBudget:       "5000000000",
	})
	require.NoError(r, err, "create publish tx")

	signature, err := deployerSigner.SignTxBlock(publishTx)
	require.NoError(r, err, "sign transaction")

	// execute the publish transaction and wait for it to be executed
	resp, err := client.SuiExecuteTransactionBlock(r.Ctx, models.SuiExecuteTransactionBlockRequest{
		TxBytes:   publishTx.TxBytes,
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
	require.True(r, resp.Effects.Status.Status == suiclient.TxStatusSuccess, resp.Effects.Status.Error)

	// find packageID
	var packageID string
	for _, change := range resp.ObjectChanges {
		if change.Type == "published" {
			packageID = change.PackageId
		}
	}
	require.NotEmpty(r, packageID, "packageID not found")

	// find objects by type filters
	objectIDs := make(map[string]string)
	for _, filter := range objectTypeFilters {
		for _, change := range resp.ObjectChanges {
			if change.Type == changeTypeCreated && strings.Contains(change.ObjectType, filter) {
				objectIDs[filter] = change.ObjectId
			}
		}
	}

	return packageID, objectIDs
}

// issueMessageContext issues a message context object in the gateway package if not exists
func (r *E2ERunner) issueMessageContext() string {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	var (
		adminCapType   = fmt.Sprintf("%s::gateway::AdminCap", r.SuiGateway.PackageID())
		msgContextType = fmt.Sprintf("%s::gateway::MessageContext", r.SuiGateway.PackageID())
	)

	// message context object exists or not
	msgContextID, found := r.suiGetOwnedObjectID(signer.Address(), msgContextType)
	if found {
		r.Logger.Info("message context object already exists, skipping issue")
		return msgContextID
	}

	// get admin cap object ID
	adminCapID, found := r.suiGetOwnedObjectID(signer.Address(), adminCapType)
	require.True(r, found, "admin cap object not found")

	// if no message context object found, issue a new one
	tx, err := r.Clients.Sui.MoveCall(r.Ctx, models.MoveCallRequest{
		Signer:          signer.Address(),
		PackageObjectId: r.SuiGateway.PackageID(),
		Module:          zetasui.GatewayModule,
		Function:        zetasui.FuncIssueMessageContext,
		TypeArguments:   []any{},
		Arguments:       []any{r.SuiGateway.ObjectID(), adminCapID},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)
	resp := r.suiExecuteTx(signer, tx)

	// find MessageContext object
	var messageContextID string
	for _, change := range resp.ObjectChanges {
		if change.Type == changeTypeCreated && strings.Contains(change.ObjectType, msgContextType) {
			messageContextID = change.ObjectId
		}
	}
	require.NotEmpty(r, messageContextID, "MessageContext object not found")

	r.Logger.Info("message context object issued: %s", messageContextID)

	return messageContextID
}

// whitelistSuiFakeUSDC deploys the FakeUSDC zrc20 on ZetaChain and whitelist it
func (r *E2ERunner) whitelistSuiFakeUSDC(signer *zetasui.SignerSecp256k1, fakeUSDCCoinType, whitelistCap string) {
	// we use DeployFungibleCoinZRC20 and whitelist manually because whitelist cctx are currently not supported for Sui
	// TODO: change this logic and use MsgWhitelistERC20 once it's supported
	// https://github.com/zeta-chain/node/issues/3569

	// deploy zrc20
	liqCap := math.NewUint(10e18)
	res, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, fungibletypes.NewMsgDeployFungibleCoinZRC20(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		fakeUSDCCoinType,
		chains.SuiLocalnet.ChainId,
		6,
		"Sui's FakeUSDC",
		"USDC.SUI",
		coin.CoinType_ERC20,
		10000,
		&liqCap,
	))
	require.NoError(r, err)

	// extract the zrc20 address from event and set the erc20 address in the runner
	deployedEvent, ok := txserver.EventOfType[*fungibletypes.EventZRC20Deployed](res.Events)
	require.True(r, ok, "unable to find deployed zrc20 event")

	r.SuiTokenZRC20Addr = ethcommon.HexToAddress(deployedEvent.Contract)
	require.NotEqualValues(r, ethcommon.Address{}, r.SuiTokenZRC20Addr)
	r.SuiTokenZRC20, err = zrc20.NewZRC20(r.SuiTokenZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)

	// whitelist zrc20
	tx, err := r.Clients.Sui.MoveCall(r.Ctx, models.MoveCallRequest{
		Signer:          signer.Address(),
		PackageObjectId: r.SuiGateway.PackageID(),
		Module:          "gateway",
		Function:        "whitelist",
		TypeArguments:   []any{"0x" + fakeUSDCCoinType},
		Arguments:       []any{r.SuiGateway.ObjectID(), whitelistCap},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)

	r.suiExecuteTx(signer, tx)
}

// set the chain params for Sui
func (r *E2ERunner) setSuiChainParams() error {
	if r.ZetaTxServer == nil {
		return errors.New("ZetaTxServer is not initialized")
	}

	creator := r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName)

	chainID := chains.SuiLocalnet.ChainId

	chainParams := &observertypes.ChainParams{
		ChainId:                     chainID,
		ZetaTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    constant.EVMZeroAddress,
		Erc20CustodyContractAddress: constant.EVMZeroAddress,
		GasPriceTicker:              5,
		WatchUtxoTicker:             0,
		InboundTicker:               2,
		OutboundTicker:              2,
		OutboundScheduleInterval:    2,
		OutboundScheduleLookahead:   5,
		BallotThreshold:             observertypes.DefaultBallotThreshold,
		MinObserverDelegation:       observertypes.DefaultMinObserverDelegation,
		IsSupported:                 true,
		GatewayAddress:              r.SuiGateway.ToAddress(),
		ConfirmationParams: &observertypes.ConfirmationParams{
			SafeInboundCount:  1,
			SafeOutboundCount: 1,
			FastInboundCount:  1,
			FastOutboundCount: 1,
		},
		ConfirmationCount: 1, // still need to be provided for now
	}
	if err := r.ZetaTxServer.UpdateChainParams(chainParams); err != nil {
		return errors.Wrap(err, "unable to broadcast solana chain params tx")
	}

	resetMsg := observertypes.NewMsgResetChainNonces(creator, chainID, 0, 0)
	if _, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, resetMsg); err != nil {
		return errors.Wrap(err, "unable to broadcast solana chain nonce reset tx")
	}

	query := &observertypes.QueryGetChainParamsForChainRequest{ChainId: chainID}

	const duration = 2 * time.Second

	for i := 0; i < 10; i++ {
		_, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, query)
		if err == nil {
			r.Logger.Print("⚙️ Sui chain params are set")
			return nil
		}

		time.Sleep(duration)
	}

	return errors.New("unable to set Sui chain params")
}

// suiTransferObjectToTSS transfers an object to the TSS
func (r *E2ERunner) suiTransferObjectToTSS(signer *zetasui.SignerSecp256k1, objectID string) {
	tx, err := r.Clients.Sui.TransferObject(r.Ctx, models.TransferObjectRequest{
		Signer:    signer.Address(),
		ObjectId:  objectID,
		Recipient: r.SuiTSSAddress,
		GasBudget: "5000000000",
	})
	require.NoError(r, err)

	r.suiExecuteTx(signer, tx)
}

// suiGetOwnedObjectID gets the first owned object ID by owner address and struct type
func (r *E2ERunner) suiGetOwnedObjectID(ownerAddress, structType string) (string, bool) {
	res, err := r.Clients.Sui.SuiXGetOwnedObjects(r.Ctx, models.SuiXGetOwnedObjectsRequest{
		Address: ownerAddress,
		Query: models.SuiObjectResponseQuery{
			Filter: map[string]any{
				"StructType": structType,
			},
		},
		Limit: 1,
	})
	require.NoError(r, err)

	if len(res.Data) == 0 {
		return "", false
	}

	return res.Data[0].Data.ObjectId, true
}
