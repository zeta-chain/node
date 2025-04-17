package runner

import (
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/zrc20.sol"

	suicontract "github.com/zeta-chain/node/e2e/contracts/sui"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	zetasui "github.com/zeta-chain/node/pkg/contracts/sui"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const changeTypeCreated = "created"

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

	// fund the TSS
	// TODO: this step might no longer necessary if a custom solution is implemented for the TSS funding
	r.RequestSuiFromFaucet(faucetURL, r.SuiTSSAddress)

	// deploy gateway package
	whitelistCapID, withdrawCapID := r.deploySUIGateway()

	// deploy SUI zrc20
	r.deploySUIZRC20()

	// deploy fake USDC and whitelist it
	fakeUSDCCoinType := r.deploySuiFakeUSDC()
	r.whitelistSuiFakeUSDC(deployerSigner, fakeUSDCCoinType, whitelistCapID)

	// deploy example contract with on_call function
	r.deploySuiExample()

	// send withdraw cap to TSS
	r.suiSendWithdrawCapToTSS(deployerSigner, withdrawCapID)

	// set the chain params
	err = r.setSuiChainParams()
	require.NoError(r, err)
}

// deploySUIGateway deploys the SUI gateway package on Sui
func (r *E2ERunner) deploySUIGateway() (whitelistCapID, withdrawCapID string) {
	const (
		filterGatewayType      = "gateway::Gateway"
		filterWithdrawCapType  = "gateway::WithdrawCap"
		filterWhitelistCapType = "gateway::WhitelistCap"
	)

	objectTypeFilters := []string{filterGatewayType, filterWhitelistCapType, filterWithdrawCapType}
	packageID, objectIDs := r.deploySuiPackage(
		[]string{suicontract.GatewayBytecodeBase64(), suicontract.EVMBytecodeBase64()},
		objectTypeFilters,
	)

	gatewayID, ok := objectIDs[filterGatewayType]
	require.True(r, ok, "gateway object not found")

	whitelistCapID, ok = objectIDs[filterWhitelistCapType]
	require.True(r, ok, "whitelistCap object not found")

	withdrawCapID, ok = objectIDs[filterWithdrawCapType]
	require.True(r, ok, "withdrawCap object not found")

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

// deploySuiFakeUSDC deploys the FakeUSDC contract on Sui
// it returns the treasuryCap object ID that allows to mint tokens
func (r *E2ERunner) deploySuiFakeUSDC() string {
	packageID, objectIDs := r.deploySuiPackage([]string{suicontract.FakeUSDCBytecodeBase64()}, []string{"TreasuryCap"})

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

// deploySuiExample deploys the example package on Sui
func (r *E2ERunner) deploySuiExample() {
	const (
		filterGlobalConfigType = "connected::GlobalConfig"
		filterPartnerType      = "connected::Partner"
		filterClockType        = "connected::Clock"
		filterPoolType         = "connected::Pool"
	)

	objectTypeFilters := []string{filterGlobalConfigType, filterPartnerType, filterClockType, filterPoolType}
	packageID, objectIDs := r.deploySuiPackage(
		[]string{suicontract.ExampleTokenBytecodeBase64(), suicontract.ExampleConnectedBytecodeBase64()},
		objectTypeFilters,
	)
	r.Logger.Info("deployed example package with packageID: %s", packageID)

	globalConfigID, ok := objectIDs[filterGlobalConfigType]
	require.True(r, ok, "globalConfig object not found")

	partnerID, ok := objectIDs[filterPartnerType]
	require.True(r, ok, "partner object not found")

	clockID, ok := objectIDs[filterClockType]
	require.True(r, ok, "clock object not found")

	poolID, ok := objectIDs[filterPoolType]
	require.True(r, ok, "pool object not found")

	r.SuiExample = NewExample(packageID, globalConfigID, partnerID, clockID, poolID)
}

// deploySuiPackage is a helper function that deploys a package on Sui
// It returns the packageID and a map of object types to their IDs
func (r *E2ERunner) deploySuiPackage(bytecodeBase64s []string, objectTypeFilters []string) (string, map[string]string) {
	client := r.Clients.Sui

	deployerSigner, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	deployerAddress := deployerSigner.Address()

	publishTx, err := client.Publish(r.Ctx, models.PublishRequest{
		Sender:          deployerAddress,
		CompiledModules: bytecodeBase64s,
		Dependencies: []string{
			"0x1", // Sui Framework
			"0x2", // Move Standard Library
		},
		GasBudget: "5000000000",
	})
	require.NoError(r, err, "create publish tx")

	signature, err := deployerSigner.SignTxBlock(publishTx)
	require.NoError(r, err, "sign transaction")

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
		GatewayAddress:              fmt.Sprintf("%s,%s", r.SuiGateway.PackageID(), r.SuiGateway.ObjectID()),
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

func (r *E2ERunner) suiSendWithdrawCapToTSS(signer *zetasui.SignerSecp256k1, withdrawCapID string) {
	tx, err := r.Clients.Sui.TransferObject(r.Ctx, models.TransferObjectRequest{
		Signer:    signer.Address(),
		ObjectId:  withdrawCapID,
		Recipient: r.SuiTSSAddress,
		GasBudget: "5000000000",
	})
	require.NoError(r, err)

	r.suiExecuteTx(signer, tx)
}
