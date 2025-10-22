package runner

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/sui/client"
)

const (
	// onCallRevertMessage is the message that triggers a revert in the 'on_call' function
	onCallRevertMessage = "revert"
)

// SuiGetSUIBalance returns the SUI balance of an address
func (r *E2ERunner) SuiGetSUIBalance(addr string) uint64 {
	resp, err := r.Clients.Sui.SuiXGetBalance(r.Ctx, models.SuiXGetBalanceRequest{
		Owner:    addr,
		CoinType: string(sui.SUI),
	})
	require.NoError(r, err)

	balance, err := strconv.ParseUint(resp.TotalBalance, 10, 64)
	require.NoError(r, err)

	return balance
}

// SuiGetFungibleTokenBalance returns the fungible token balance of an address
func (r *E2ERunner) SuiGetFungibleTokenBalance(addr string) uint64 {
	resp, err := r.Clients.Sui.SuiXGetBalance(r.Ctx, models.SuiXGetBalanceRequest{
		Owner:    addr,
		CoinType: "0x" + r.SuiTokenCoinType,
	})
	require.NoError(r, err)

	balance, err := strconv.ParseUint(resp.TotalBalance, 10, 64)
	require.NoError(r, err)

	return balance
}

// SuiWithdraw calls Withdraw on ZEVM Gateway with given ZRC20
func (r *E2ERunner) SuiWithdraw(
	receiver string,
	amount *big.Int,
	zrc20 ethcommon.Address,
	revertOptions gatewayzevm.RevertOptions,
) *ethtypes.Transaction {
	receiverBytes, err := hex.DecodeString(receiver[2:])
	require.NoError(r, err, "receiver: "+receiver[2:])

	tx, err := r.GatewayZEVM.Withdraw(
		r.ZEVMAuth,
		receiverBytes,
		amount,
		zrc20,
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// SuiWithdrawAndCall calls WithdrawAndCall on ZEVM Gateway with given ZRC20
func (r *E2ERunner) SuiWithdrawAndCall(
	receiver string,
	amount *big.Int,
	zrc20 ethcommon.Address,
	payload sui.CallPayload,
	gasLimit *big.Int,
	revertOptions gatewayzevm.RevertOptions,
) *ethtypes.Transaction {
	receiverBytes, err := hex.DecodeString(receiver[2:])
	require.NoError(r, err, "receiver: "+receiver[2:])

	// ACT
	payloadBytes, err := payload.PackABI()
	require.NoError(r, err)

	tx, err := r.GatewayZEVM.WithdrawAndCall0(
		r.ZEVMAuth,
		receiverBytes,
		amount,
		zrc20,
		payloadBytes,
		gatewayzevm.CallOptions{
			GasLimit:        gasLimit,
			IsArbitraryCall: false, // always authenticated call
		},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// SuiDepositSUI calls Deposit on Sui
func (r *E2ERunner) SuiDepositSUI(
	receiver ethcommon.Address,
	amount math.Uint,
) models.SuiTransactionBlockResponse {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	// retrieve SUI coin object to deposit
	coinObjectID := r.suiSplitSUI(signer, amount)

	// create the tx
	return r.suiExecuteDeposit(signer, string(sui.SUI), coinObjectID, receiver)
}

// SuiDepositAndCallSUI calls DepositAndCall on Sui
func (r *E2ERunner) SuiDepositAndCallSUI(
	receiver ethcommon.Address,
	amount math.Uint,
	payload []byte,
) models.SuiTransactionBlockResponse {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	// retrieve SUI coin object to deposit
	coinObjectID := r.suiSplitSUI(signer, amount)

	// create the tx
	return r.suiExecuteDepositAndCall(signer, string(sui.SUI), coinObjectID, receiver, payload)
}

// SuiDepositFungibleToken calls Deposit with fungible token on Sui
func (r *E2ERunner) SuiDepositFungibleToken(
	receiver ethcommon.Address,
	amount math.Uint,
) models.SuiTransactionBlockResponse {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	// retrieve SUI coin object to deposit
	coinObjectID := r.suiSplitUSDC(signer, amount)

	// create the tx
	return r.suiExecuteDeposit(signer, "0x"+r.SuiTokenCoinType, coinObjectID, receiver)
}

// SuiFungibleTokenDepositAndCall calls DepositAndCall with fungible token on Sui
func (r *E2ERunner) SuiFungibleTokenDepositAndCall(
	receiver ethcommon.Address,
	amount math.Uint,
	payload []byte,
) models.SuiTransactionBlockResponse {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	// retrieve SUI coin object to deposit
	coinObjectID := r.suiSplitUSDC(signer, amount)

	// create the tx
	return r.suiExecuteDepositAndCall(signer, "0x"+r.SuiTokenCoinType, coinObjectID, receiver, payload)
}

// SuiMintUSDC mints FakeUSDC on Sui to a receiver
// this function requires the signer to be the owner of the trasuryCap
func (r *E2ERunner) SuiMintUSDC(
	amount,
	receiver string,
) models.SuiTransactionBlockResponse {
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	// extract the package ID from the coin type
	splitted := strings.Split(r.SuiTokenCoinType, "::")
	require.Len(r, splitted, 3, "coinType should be in format <packageID>::<module>::<name> - %s", r.SuiTokenCoinType)
	packageID := "0x" + splitted[0]

	// create the tx
	tx, err := r.Clients.Sui.MoveCall(r.Ctx, models.MoveCallRequest{
		Signer:          signer.Address(),
		PackageObjectId: packageID,
		Module:          "fake_usdc",
		Function:        "mint",
		TypeArguments:   []any{},
		Arguments:       []any{r.SuiTokenTreasuryCap, amount, receiver},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)

	return r.suiExecuteTx(signer, tx)
}

// SuiCreateExampleWACPayload creates a payload for on_call function in Sui the example package
// The payload message contains below three fields in order:
// field 0: [42-byte ZEVM sender], checksum address, with 0x prefix
// field 1: [32-byte Sui target package], without 0x prefix
// field 2: [32-byte Sui receiver], without 0x prefix
//
// The first two fields are used by E2E test to easily mock up the verifications against 'MessageContext' passed to the 'on_call'.
// A real-world app just encodes useful data that meets its needs, and it will be only the receiver 'suiAddress' in this case.
func (r *E2ERunner) SuiCreateExampleWACPayload(authorizedSender ethcommon.Address, suiAddress string) sui.CallPayload {
	// only the CCTX's coinType is needed, no additional type argument
	argumentTypes := []string{}
	objects := []string{
		r.SuiExample.GlobalConfigID.String(),
		r.SuiExample.PartnerID.String(),
		r.SuiExample.ClockID.String(),
	}

	message := make([]byte, 106)

	// field 0
	copy(message[:42], []byte(authorizedSender.Hex()))

	// field 1
	target, err := hex.DecodeString(r.SuiExample.PackageID.String()[2:])
	require.NoError(r, err)
	copy(message[42:74], target)

	// field 2
	receiver, err := hex.DecodeString(suiAddress[2:])
	require.NoError(r, err)
	copy(message[74:], receiver)

	return sui.NewCallPayload(argumentTypes, objects, message)
}

// TODO: https://github.com/zeta-chain/node/issues/4066
// remove legacy payload builder after re-enabling authenticated call
// SuiCreateExampleWACPayload creates a payload for on_call function in Sui the example package
// The example on_call function will just forward the withdrawn token to given 'suiAddress'
func (r *E2ERunner) SuiCreateExampleWACPayloadLegacy(suiAddress string) (sui.CallPayload, error) {
	// only the CCTX's coinType is needed, no additional arguments
	argumentTypes := []string{}
	objects := []string{
		r.SuiExample.GlobalConfigID.String(),
		r.SuiExample.PartnerID.String(),
		r.SuiExample.ClockID.String(),
	}

	// create the payload message from the sui address
	message, err := hex.DecodeString(suiAddress[2:]) // remove 0x prefix
	if err != nil {
		return sui.CallPayload{}, err
	}

	return sui.NewCallPayload(argumentTypes, objects, message), nil
}

// SuiCreateExampleWACPayload creates a payload that triggers a revert in the 'on_call'
// function in Sui the example package
func (r *E2ERunner) SuiCreateExampleWACPayloadForRevert() (sui.CallPayload, error) {
	// only the CCTX's coinType is needed, no additional arguments
	argumentTypes := []string{}
	objects := []string{
		r.SuiExample.GlobalConfigID.String(),
		r.SuiExample.PartnerID.String(),
		r.SuiExample.ClockID.String(),
	}

	// the 'on_call' method of the "connected" contract specifically throws an error
	// for this message to trigger "tx revert" test case
	message := []byte(onCallRevertMessage)

	return sui.NewCallPayload(argumentTypes, objects, message), nil
}

// SuiGetConnectedCalledCount reads the called_count from the GlobalConfig object in connected module
func (r *E2ERunner) SuiGetConnectedCalledCount() uint64 {
	// Get object data
	resp, err := r.Clients.Sui.SuiGetObject(r.Ctx, models.SuiGetObjectRequest{
		ObjectId: r.SuiExample.GlobalConfigID.String(),
		Options:  models.SuiObjectDataOptions{ShowContent: true},
	})

	require.NoError(r, err)
	require.Nil(r, resp.Error)
	require.NotNil(r, resp.Data)
	require.NotNil(r, resp.Data.Content)

	fields := resp.Data.Content.Fields

	// Extract called_count field from the object content
	rawValue, ok := fields["called_count"]
	require.True(r, ok, "missing called_count field")

	v, ok := rawValue.(string)
	require.True(r, ok, "want string, got %T for called_count", rawValue)

	// #nosec G115 always in range
	calledCount, err := strconv.ParseUint(v, 10, 64)
	require.NoError(r, err)

	return calledCount
}

// SuiMonitorCCTXByInboundHash monitors a CCTX by inbound hash until it gets mined
// This function wrapps WaitCctxMinedByInboundHash and prints additional logs needed in stress test
func (r *E2ERunner) SuiMonitorCCTXByInboundHash(inboundHash string, index int) (time.Duration, error) {
	startTime := time.Now()

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, inboundHash, r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		return 0, fmt.Errorf(
			"index %d: cctx failed with status %s, message %s, index %s",
			index,
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
			cctx.Index,
		)
	}

	timeElapsed := time.Since(startTime)
	r.Logger.Print("index %d: cctx succeeded in %s", index, timeElapsed.String())

	return timeElapsed, nil
}

// suiExecuteDeposit executes a deposit on the SUI contract
func (r *E2ERunner) suiExecuteDeposit(
	signer *sui.SignerSecp256k1,
	coinType string,
	coinObjectID string,
	receiver ethcommon.Address,
) models.SuiTransactionBlockResponse {
	// create the tx
	tx, err := r.Clients.Sui.MoveCall(r.Ctx, models.MoveCallRequest{
		Signer:          signer.Address(),
		PackageObjectId: r.SuiGateway.PackageID(),
		Module:          "gateway",
		Function:        "deposit",
		TypeArguments:   []any{coinType},
		Arguments:       []any{r.SuiGateway.ObjectID(), coinObjectID, receiver.Hex()},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)

	return r.suiExecuteTx(signer, tx)
}

// suiExecuteDepositAndCall executes a depositAndCall on the SUI contract
func (r *E2ERunner) suiExecuteDepositAndCall(
	signer *sui.SignerSecp256k1,
	coinType string,
	coinObjectID string,
	receiver ethcommon.Address,
	payload []byte,
) models.SuiTransactionBlockResponse {
	// create the tx
	tx, err := r.Clients.Sui.MoveCall(r.Ctx, models.MoveCallRequest{
		Signer:          signer.Address(),
		PackageObjectId: r.SuiGateway.PackageID(),
		Module:          "gateway",
		Function:        "deposit_and_call",
		TypeArguments:   []any{coinType},
		Arguments:       []any{r.SuiGateway.ObjectID(), coinObjectID, receiver.Hex(), payload},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)

	return r.suiExecuteTx(signer, tx)
}

// suiSplitUSDC splits USDC coin and obtain a USDC coin object with the wanted balance
func (r *E2ERunner) suiSplitUSDC(signer *sui.SignerSecp256k1, balance math.Uint) (objID string) {
	// find the coin to split
	originalCoin := r.suiFindCoinWithBalanceAbove(signer.Address(), balance, "0x"+r.SuiTokenCoinType)

	tx, err := r.Clients.Sui.SplitCoin(r.Ctx, models.SplitCoinRequest{
		Signer:       signer.Address(),
		CoinObjectId: originalCoin,
		SplitAmounts: []string{balance.String()},
		GasBudget:    "5000000000",
	})

	require.NoError(r, err)
	r.suiExecuteTx(signer, tx)

	// find the split coin
	return r.suiFindCoinWithBalance(signer.Address(), balance, "0x"+r.SuiTokenCoinType)
}

// suiSplitSUI splits SUI coin and obtain a SUI coin object with the wanted balance
func (r *E2ERunner) suiSplitSUI(signer *sui.SignerSecp256k1, balance math.Uint) (objID string) {
	// find the coin to split
	originalCoin := r.suiFindCoinWithBalanceAbove(signer.Address(), balance, string(sui.SUI))

	// split the coin using the PaySui API
	tx, err := r.Clients.Sui.PaySui(r.Ctx, models.PaySuiRequest{
		Signer:      signer.Address(),
		SuiObjectId: []string{originalCoin},
		Recipient:   []string{signer.Address()},
		Amount:      []string{balance.String()},
		GasBudget:   "5000000000",
	})
	require.NoError(r, err)

	r.suiExecuteTx(signer, tx)

	// find the split coin
	return r.suiFindCoinWithBalance(signer.Address(), balance, string(sui.SUI))
}

func (r *E2ERunner) suiFindCoinWithBalance(
	address string,
	balance math.Uint,
	coinType string,
) (coinID string) {
	return r.suiFindCoin(address, balance, coinType, func(a, b math.Uint) bool {
		return a.Equal(b)
	})
}

func (r *E2ERunner) suiFindCoinWithBalanceAbove(
	address string,
	balanceAbove math.Uint,
	coinType string,
) (coinID string) {
	return r.suiFindCoin(address, balanceAbove, coinType, func(a, b math.Uint) bool {
		return a.GTE(b)
	})
}

type compFunc func(a, b math.Uint) bool

func (r *E2ERunner) suiFindCoin(
	address string,
	balance math.Uint,
	coinType string,
	comp compFunc,
) (coinID string) {
	res, err := r.Clients.Sui.SuiXGetCoins(r.Ctx, models.SuiXGetCoinsRequest{
		Owner:    address,
		CoinType: coinType,
	})
	require.NoError(r, err)

	for _, data := range res.Data {
		coinBalance, err := math.ParseUint(data.Balance)
		require.NoError(r, err)

		if comp(coinBalance, balance) {
			return data.CoinObjectId
		}
	}

	require.FailNow(r, fmt.Sprintf("coin %s not found for address %s", coinType, address))
	return ""
}

func (r *E2ERunner) suiExecuteTx(
	signer *sui.SignerSecp256k1,
	tx models.TxnMetaData,
) models.SuiTransactionBlockResponse {
	// sign the tx
	signature, err := signer.SignTxBlock(tx)
	require.NoError(r, err, "sign transaction")

	resp, err := r.Clients.Sui.SuiExecuteTransactionBlock(r.Ctx, models.SuiExecuteTransactionBlockRequest{
		TxBytes:     tx.TxBytes,
		Signature:   []string{signature},
		Options:     models.SuiTransactionBlockOptions{ShowEffects: true},
		RequestType: "WaitForLocalExecution",
	})
	require.NoError(r, err)
	require.Equal(r, resp.Effects.Status.Status, client.TxStatusSuccess)

	return resp
}
