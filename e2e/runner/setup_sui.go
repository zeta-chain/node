package runner

import (
	"fmt"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	suicontract "github.com/zeta-chain/node/e2e/contracts/sui"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

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

	// deploy fake USDC
	_, _ = r.deployFakeUSDC()

	// set Sui contract values
	r.GatewayPackageID = packageID
	r.GatewayObjectID = objectID

	// set the chain params
	err = r.ensureSuiChainParams()
	require.NoError(r, err)
}

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
		if change.Type == "created" && strings.Contains(change.ObjectType, "TreasuryCap") {
			treasuryCap = change.ObjectId
		}
	}
	require.NotEmpty(r, treasuryCap, "find objectID")

	return packageID, treasuryCap
}

func (r *E2ERunner) ensureSuiChainParams() error {
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
