package runner

import (
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	suicontract "github.com/zeta-chain/node/e2e/contracts/sui"
	"github.com/zeta-chain/node/e2e/utils"
	suiutils "github.com/zeta-chain/node/e2e/utils/sui"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const changeTypeCreated = "created"

// RequestSuiFaucetToken requests SUI tokens from the faucet for the runner account
func (r *E2ERunner) RequestSuiFaucetToken(faucetURL string) {
	deployerSigner, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	deployerAddress := deployerSigner.Address()

	header := map[string]string{}
	err = sui.RequestSuiFromFaucet(faucetURL, deployerAddress, header)
	require.NoError(r, err, "sui faucet request to %s", faucetURL)
}

// SetupSui initializes the gateway package on Sui and initialize the chain params on ZetaChain
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

	// find packageID
	var packageID, gatewayID, whitelistID string
	for _, change := range resp.ObjectChanges {
		if change.Type == "published" {
			packageID = change.PackageId
		}
	}
	require.NotEmpty(r, packageID, "find packageID")

	// find gateway objectID
	gatewayType := fmt.Sprintf("%s::gateway::Gateway", packageID)
	for _, change := range resp.ObjectChanges {
		if change.Type == changeTypeCreated && change.ObjectType == gatewayType {
			gatewayID = change.ObjectId
		}
	}
	require.NotEmpty(r, gatewayID, "find gatewayID")

	// find whitelist objectID
	whitelistType := fmt.Sprintf("%s::gateway::WhitelistCap", packageID)
	for _, change := range resp.ObjectChanges {
		if change.Type == changeTypeCreated && change.ObjectType == whitelistType {
			whitelistID = change.ObjectId
		}
	}

	// set Sui gateway values
	r.GatewayPackageID = packageID
	r.GatewayObjectID = gatewayID

	// deploy fake USDC
	fakeUSDCCoinType, treasuryCap := r.deployFakeUSDC()
	r.whitelistFakeUSDC(deployerSigner, fakeUSDCCoinType, whitelistID)

	r.SuiTokenCoinType = fakeUSDCCoinType
	r.SuiTokenTreasuryCap = treasuryCap

	// set the chain params
	err = r.setSuiChainParams()
	require.NoError(r, err)
}

// deployFakeUSDC deploys the FakeUSDC contract on Sui
func (r *E2ERunner) deployFakeUSDC() (string, string) {
	client := r.Clients.Sui
	deployerSigner, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	deployerAddress := deployerSigner.Address()

	publishReq, err := client.Publish(r.Ctx, models.PublishRequest{
		Sender:          deployerAddress,
		CompiledModules: []string{suicontract.FakeUSDCBytecodeBase64()},
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

	var packageID, treasuryCap string
	for _, change := range resp.ObjectChanges {
		if change.Type == "published" {
			packageID = change.PackageId
		}
	}
	require.NotEmpty(r, packageID, "find packageID")

	for _, change := range resp.ObjectChanges {
		if change.Type == changeTypeCreated && strings.Contains(change.ObjectType, "TreasuryCap") {
			treasuryCap = change.ObjectId
		}
	}
	require.NotEmpty(r, treasuryCap, "find objectID")

	coinType := packageID + "::fake_usdc::FAKE_USDC"

	// strip 0x from packageID
	coinType = coinType[2:]

	return coinType, treasuryCap
}

// whitelistFakeUSDC deploys the FakeUSDC zrc20 on ZetaChain and whitelist it
func (r *E2ERunner) whitelistFakeUSDC(signer *suiutils.SignerSecp256k1, fakeUSDCCoinType, whitelistCap string) {
	// we use DeployFungibleCoinZRC20 and whitelist manually because whitelist cctx are currently not supported for Sui
	// TODO: change this logic and use MsgWhitelistERC20 once it's supported
	// https://github.com/zeta-chain/node/issues/3569

	// deploy zrc20
	liqCap := math.NewUint(10e18)
	_, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, fungibletypes.NewMsgDeployFungibleCoinZRC20(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		fakeUSDCCoinType,
		chains.SuiLocalnet.ChainId,
		6,
		"Sui's FakeUSDC",
		"USDC.SUI",
		coin.CoinType_ERC20,
		100000,
		&liqCap,
	))
	require.NoError(r, err)

	// whitelist zrc20
	tx, err := r.Clients.Sui.MoveCall(r.Ctx, models.MoveCallRequest{
		Signer:          signer.Address(),
		PackageObjectId: r.GatewayPackageID,
		Module:          "gateway",
		Function:        "whitelist",
		TypeArguments:   []any{"0x" + fakeUSDCCoinType},
		Arguments:       []any{r.GatewayObjectID, whitelistCap},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)

	r.executeSuiTx(signer, tx)
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
		GatewayAddress:              r.GatewayPackageID,
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
