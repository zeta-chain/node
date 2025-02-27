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

	header := map[string]string{}
	err = sui.RequestSuiFromFaucet(faucetURL, deployerAddress, header)
	require.NoError(r, err, "sui faucet request to %s", faucetURL)

	client := r.Clients.Sui

	publishTx, err := client.Publish(r.Ctx, models.PublishRequest{
		Sender:          deployerAddress,
		CompiledModules: []string{suicontract.GatewayBytecodeBase64()},
		Dependencies: []string{
			"0x0000000000000000000000000000000000000000000000000000000000000001",
			"0x0000000000000000000000000000000000000000000000000000000000000002",
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
	var packageID, gatewayID, whitelistCapID, withdrawCapID string
	for _, change := range resp.ObjectChanges {
		if change.Type == "published" {
			packageID = change.PackageId
		}
	}
	require.NotEmpty(r, packageID, "packageID not found")

	// find gateway objectID
	gatewayType := fmt.Sprintf("%s::gateway::Gateway", packageID)
	for _, change := range resp.ObjectChanges {
		if change.Type == changeTypeCreated && change.ObjectType == gatewayType {
			gatewayID = change.ObjectId
		}
	}
	require.NotEmpty(r, gatewayID, "gatewayID not found")

	// find WhitelistCap objectID
	whitelistType := fmt.Sprintf("%s::gateway::WhitelistCap", packageID)
	for _, change := range resp.ObjectChanges {
		if change.Type == changeTypeCreated && change.ObjectType == whitelistType {
			whitelistCapID = change.ObjectId
		}
	}
	require.NotEmpty(r, whitelistCapID, "whitelistID not found")

	// find WithdrawCap objectID
	withdrawCapType := fmt.Sprintf("%s::gateway::WithdrawCap", packageID)
	for _, change := range resp.ObjectChanges {
		if change.Type == changeTypeCreated && change.ObjectType == withdrawCapType {
			withdrawCapID = change.ObjectId
		}
	}

	// set sui gateway
	r.SuiGateway = zetasui.NewGateway(packageID, gatewayID)

	// deploy fake USDC
	fakeUSDCCoinType, treasuryCap := r.deployFakeUSDC()
	r.whitelistSuiFakeUSDC(deployerSigner, fakeUSDCCoinType, whitelistCapID)

	r.SuiTokenCoinType = fakeUSDCCoinType
	r.SuiTokenTreasuryCap = treasuryCap

	// send withdraw cap to TSS
	r.sendWithdrawCapToTSS(deployerSigner, withdrawCapID)

	// set the chain params
	err = r.setSuiChainParams()
	require.NoError(r, err)
}

// deployFakeUSDC deploys the FakeUSDC contract on Sui
// it returns the coinType to be used as asset value for zrc20 and treasuryCap object ID that allows to mint tokens
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

	signature, err := deployerSigner.SignTxBlock(publishReq)
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
		100000,
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

func (r *E2ERunner) sendWithdrawCapToTSS(signer *zetasui.SignerSecp256k1, withdrawCapID string) {
	tx, err := r.Clients.Sui.TransferObject(r.Ctx, models.TransferObjectRequest{
		Signer:    signer.Address(),
		ObjectId:  withdrawCapID,
		Recipient: r.SuiTSSAddress,
		GasBudget: "5000000000",
	})
	require.NoError(r, err)

	r.suiExecuteTx(signer, tx)
}
