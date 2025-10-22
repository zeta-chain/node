package runner

import (
	"fmt"
	"math/big"
	"time"

	"github.com/cenkalti/backoff/v4"
	query "github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/zrc20.sol"

	"github.com/zeta-chain/node/e2e/contracts/gatewayzevmcaller"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/retry"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

var defaultGasLimit = big.NewInt(250000)

// ApproveETHZRC20 approves ETH ZRC20 on EVM to a specific address
func (r *E2ERunner) ApproveETHZRC20(allowed ethcommon.Address) {
	r.approveZRC20(allowed, r.ETHZRC20)
}

// ApproveERC20ZRC20 approves ERC20 ZRC20 on EVM to a specific address
func (r *E2ERunner) ApproveERC20ZRC20(allowed ethcommon.Address) {
	r.approveZRC20(allowed, r.ERC20ZRC20)
}

// ApproveBTCZRC20 approves BTC ZRC20 on EVM to a specific address
func (r *E2ERunner) ApproveBTCZRC20(allowed ethcommon.Address) {
	r.approveZRC20(allowed, r.BTCZRC20)
}

// ApproveSOLZRC20 approves SOL ZRC20 on EVM to a specific address
func (r *E2ERunner) ApproveSOLZRC20(allowed ethcommon.Address) {
	r.approveZRC20(allowed, r.SOLZRC20)
}

// ApproveSPLZRC20 approves SPL ZRC20 on EVM to a specific address
func (r *E2ERunner) ApproveSPLZRC20(allowed ethcommon.Address) {
	r.approveZRC20(allowed, r.SPLZRC20)
}

// ApproveSUIZRC20 approves SUI ZRC20 on EVM to a specific address
func (r *E2ERunner) ApproveSUIZRC20(allowed ethcommon.Address) {
	r.approveZRC20(allowed, r.SUIZRC20)
}

// ApproveFungibleTokenZRC20 approves Sui fungible token ZRC20 on EVM to a specific address
func (r *E2ERunner) ApproveFungibleTokenZRC20(allowed ethcommon.Address) {
	r.approveZRC20(allowed, r.SuiTokenZRC20)
}

// ApproveTONZRC20 approves TON ZRC20 on EVM to a specific address
func (r *E2ERunner) ApproveTONZRC20(allowed ethcommon.Address) {
	r.approveZRC20(allowed, r.TONZRC20)
}

// approveZRC20 approves ZRC20 on EVM to a specific address
// check if allowance is zero before calling this method
// allow a high amount to avoid multiple approvals
func (r *E2ERunner) approveZRC20(allowed ethcommon.Address, zrc20 *zrc20.ZRC20) {
	allowance, err := zrc20.Allowance(&bind.CallOpts{}, r.Account.EVMAddress(), allowed)
	require.NoError(r, err)

	// approve 1M*1e18 if allowance is below 1k
	thousand := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000))
	if allowance.Cmp(thousand) < 0 {
		r.Logger.Info("Approving %s to %s", r.Account.EVMAddress().String(), allowed.String())
		tx, err := zrc20.Approve(r.ZEVMAuth, allowed, big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000000)))
		require.NoError(r, err)
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		require.True(r, receipt.Status == 1, "approval failed")
	}
}

// ETHWithdraw calls Withdraw of Gateway with gas token on ZEVM
func (r *E2ERunner) ETHWithdraw(
	receiver ethcommon.Address,
	amount *big.Int,
	revertOptions gatewayzevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.Withdraw0(
		r.ZEVMAuth,
		receiver.Bytes(),
		amount,
		r.ETHZRC20Addr,
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ETHWithdrawAndArbitraryCall calls WithdrawAndCall of Gateway with gas token on ZEVM using arbitrary call
func (r *E2ERunner) ETHWithdrawAndArbitraryCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewayzevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.WithdrawAndCall(
		r.ZEVMAuth,
		receiver.Bytes(),
		amount,
		r.ETHZRC20Addr,
		payload,
		gatewayzevm.CallOptions{GasLimit: defaultGasLimit, IsArbitraryCall: true},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ETHWithdrawAndCall calls WithdrawAndCall of Gateway with gas token on ZEVM using authenticated call
func (r *E2ERunner) ETHWithdrawAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewayzevm.RevertOptions,
	gasLimit *big.Int,
) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.WithdrawAndCall(
		r.ZEVMAuth,
		receiver.Bytes(),
		amount,
		r.ETHZRC20Addr,
		payload,
		gatewayzevm.CallOptions{
			IsArbitraryCall: false,
			GasLimit:        gasLimit,
		},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ETHWithdrawAndCallThroughContract calls WithdrawAndCall of Gateway with gas token on ZEVM using authenticated call
// through contract
func (r *E2ERunner) ETHWithdrawAndCallThroughContract(
	gatewayZEVMCaller *gatewayzevmcaller.GatewayZEVMCaller,
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewayzevmcaller.RevertOptions,
) *ethtypes.Transaction {
	tx, err := gatewayZEVMCaller.WithdrawAndCallGatewayZEVM(
		r.ZEVMAuth,
		receiver.Bytes(),
		amount,
		r.ETHZRC20Addr,
		payload,
		gatewayzevmcaller.CallOptions{
			IsArbitraryCall: false,
			GasLimit:        defaultGasLimit,
		},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ERC20Withdraw calls Withdraw of Gateway with erc20 token on ZEVM
func (r *E2ERunner) ERC20Withdraw(
	receiver ethcommon.Address,
	amount *big.Int,
	revertOptions gatewayzevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.Withdraw0(
		r.ZEVMAuth,
		receiver.Bytes(),
		amount,
		r.ERC20ZRC20Addr,
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ERC20WithdrawAndArbitraryCall calls WithdrawAndCall of Gateway with erc20 token on ZEVM using arbitrary call
func (r *E2ERunner) ERC20WithdrawAndArbitraryCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewayzevm.RevertOptions,
) *ethtypes.Transaction {
	// this function take more gas than default 500k
	// so we need to increase the gas limit
	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10000000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	tx, err := r.GatewayZEVM.WithdrawAndCall(
		r.ZEVMAuth,
		receiver.Bytes(),
		amount,
		r.ERC20ZRC20Addr,
		payload,
		gatewayzevm.CallOptions{GasLimit: defaultGasLimit, IsArbitraryCall: true},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ERC20WithdrawAndCall calls WithdrawAndCall of Gateway with erc20 token on ZEVM using authenticated call
func (r *E2ERunner) ERC20WithdrawAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewayzevm.RevertOptions,
	gasLimit *big.Int,
) *ethtypes.Transaction {
	// this function take more gas than default 500k
	// so we need to increase the gas limit
	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10000000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	tx, err := r.GatewayZEVM.WithdrawAndCall(
		r.ZEVMAuth,
		receiver.Bytes(),
		amount,
		r.ERC20ZRC20Addr,
		payload,
		gatewayzevm.CallOptions{GasLimit: gasLimit, IsArbitraryCall: false},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ZEVMToEMVArbitraryCall calls Call of Gateway on ZEVM using arbitrary call
func (r *E2ERunner) ZEVMToEMVArbitraryCall(
	receiver ethcommon.Address,
	payload []byte,
	revertOptions gatewayzevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.Call(
		r.ZEVMAuth,
		receiver.Bytes(),
		r.ETHZRC20Addr,
		payload,
		gatewayzevm.CallOptions{GasLimit: defaultGasLimit, IsArbitraryCall: true},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ZEVMToEMVCall calls authenticated Call of Gateway on ZEVM using authenticated call
func (r *E2ERunner) ZEVMToEMVCall(
	receiver ethcommon.Address,
	payload []byte,
	revertOptions gatewayzevm.RevertOptions,
	gasLimit *big.Int,
) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.Call(
		r.ZEVMAuth,
		receiver.Bytes(),
		r.ETHZRC20Addr,
		payload,
		gatewayzevm.CallOptions{
			GasLimit:        gasLimit,
			IsArbitraryCall: false,
		},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ZEVMToEMVCallThroughContract calls authenticated Call of Gateway on ZEVM through contract using authenticated call
func (r *E2ERunner) ZEVMToEMVCallThroughContract(
	gatewayZEVMCaller *gatewayzevmcaller.GatewayZEVMCaller,
	receiver ethcommon.Address,
	payload []byte,
	revertOptions gatewayzevmcaller.RevertOptions,
) *ethtypes.Transaction {
	tx, err := gatewayZEVMCaller.CallGatewayZEVM(
		r.ZEVMAuth,
		receiver.Bytes(),
		r.ETHZRC20Addr,
		payload,
		gatewayzevmcaller.CallOptions{
			GasLimit:        defaultGasLimit,
			IsArbitraryCall: false,
		},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// WaitForBlocks waits for a specific number of blocks to be generated
// The parameter n is the number of blocks to wait for
func (r *E2ERunner) WaitForBlocks(n int64) {
	height, err := r.CctxClient.LastZetaHeight(r.Ctx, &types.QueryLastZetaHeightRequest{})
	if err != nil {
		return
	}
	call := func() error {
		return retry.Retry(r.waitForBlock(height.Height + n))
	}
	retryBuffer := uint64(20)
	bo := backoff.NewConstantBackOff(time.Second * 6)
	// #nosec G115 always in range
	boWithMaxRetries := backoff.WithMaxRetries(bo, uint64(n)+retryBuffer)
	err = retry.DoWithBackoff(call, boWithMaxRetries)
	require.NoError(r, err, "failed to wait for %d blocks", n)
}

// WaitForTSSGeneration waits for a specific number of TSS to be generated
// The parameter n is the number of TSS to wait for
func (r *E2ERunner) WaitForTSSGeneration(tssNumber int64) {
	call := func() error {
		return retry.Retry(r.checkNumberOfTSSGenerated(tssNumber))
	}
	bo := backoff.NewConstantBackOff(time.Second * 5)
	boWithMaxRetries := backoff.WithMaxRetries(bo, 10)
	err := retry.DoWithBackoff(call, boWithMaxRetries)
	require.NoError(r, err, "failed to wait for %d tss generation", tssNumber)
}

// checkNumberOfTSSGenerated checks the number of TSS generated
// if the number of tss is less that the `tssNumber` provided we return an error
func (r *E2ERunner) checkNumberOfTSSGenerated(tssNumber int64) error {
	tssList, err := r.ObserverClient.TssHistory(r.Ctx, &observertypes.QueryTssHistoryRequest{})
	if err != nil {
		return err
	}
	if int64(len(tssList.TssList)) < tssNumber {
		return fmt.Errorf("waiting for %d tss generation, number of TSS :%d", tssNumber, len(tssList.TssList))
	}
	return nil
}

func (r *E2ERunner) waitForBlock(n int64) error {
	height, err := r.CctxClient.LastZetaHeight(r.Ctx, &types.QueryLastZetaHeightRequest{})
	if err != nil {
		return err
	}
	if height.Height < n {
		return fmt.Errorf("waiting for height: %d, current height: %d", n, height.Height)
	}
	return nil
}

// WaitForTxReceiptOnZEVM waits for a tx receipt on ZEVM
func (r *E2ERunner) WaitForTxReceiptOnZEVM(tx *ethtypes.Transaction) {
	r.Lock()
	defer r.Unlock()

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt)
}

// WaitForMinedCCTX waits for a cctx to be mined from a tx
func (r *E2ERunner) WaitForMinedCCTX(txHash ethcommon.Hash) {
	r.Lock()
	defer r.Unlock()

	cctx := utils.WaitCctxMinedByInboundHash(
		r.Ctx,
		txHash.Hex(),
		r.CctxClient,
		r.Logger,
		r.CctxTimeout,
	)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_OutboundMined)
}

// WaitForMinedCCTXFromIndex waits for a cctx to be mined from its index
func (r *E2ERunner) WaitForMinedCCTXFromIndex(index string) *types.CrossChainTx {
	return r.waitForMinedCCTXFromIndex(index, types.CctxStatus_OutboundMined)
}

func (r *E2ERunner) waitForMinedCCTXFromIndex(index string, status types.CctxStatus) *types.CrossChainTx {
	r.Lock()
	defer r.Unlock()

	cctx := utils.WaitCCTXMinedByIndex(r.Ctx, index, r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, status)

	return cctx
}

// WaitForSpecificCCTX scans for cctx by filters and ensures it's mined
func (r *E2ERunner) WaitForSpecificCCTX(
	filter func(*types.CrossChainTx) bool,
	status types.CctxStatus,
	timeout time.Duration,
) *types.CrossChainTx {
	var (
		ctx      = r.Ctx
		start    = time.Now()
		reqQuery = &types.QueryAllCctxRequest{
			Pagination: &query.PageRequest{
				Limit:   50,
				Reverse: false,
			},
		}
	)

	for time.Since(start) < timeout {
		res, err := r.CctxClient.CctxAll(ctx, reqQuery)
		require.NoError(r, err)

		for i := range res.CrossChainTx {
			tx := res.CrossChainTx[i]
			if filter(tx) {
				return r.waitForMinedCCTXFromIndex(tx.Index, status)
			}
		}

		time.Sleep(time.Second)
	}

	r.Logger.Error("WaitForSpecificCCTX: No CCTX found. Timed out")
	r.FailNow()

	return nil
}

// skipChainOperations checks if the chain operations should be skipped for E2E
func (r *E2ERunner) skipChainOperations(chainID int64) bool {
	skip := r.IsRunningUpgrade() && chains.IsTONChain(chainID, nil)

	if skip {
		r.Logger.Print("Skipping chain operations for chain %d", chainID)
	}

	return skip
}

// AddInboundTracker adds an inbound tracker from the tx hash
func (r *E2ERunner) AddInboundTracker(coinType coin.CoinType, txHash string) {
	require.NotNil(r, r.ZetaTxServer)

	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	msg := types.NewMsgAddInboundTracker(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.EmergencyPolicyName),
		chainID.Int64(),
		coinType,
		txHash,
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.EmergencyPolicyName, msg)
	require.NoError(r, err)
}

// UpdateGatewayGasLimit updates the gateway gas limit used by the fungible module for ZEVM calls
func (r *E2ERunner) UpdateGatewayGasLimit(newGasLimit uint64) {
	msgUpdateGatewayGasLimit := fungibletypes.NewMsgUpdateGatewayGasLimit(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		newGasLimit,
	)
	_, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msgUpdateGatewayGasLimit)
	require.NoError(r, err)

	r.WaitForBlocks(1)

	// Verify that the gas limit has been updated
	systemContract, err := r.FungibleClient.SystemContract(r.Ctx, &fungibletypes.QueryGetSystemContractRequest{})
	require.NoError(r, err)
	require.Equal(r, newGasLimit, systemContract.SystemContract.GatewayGasLimit)
}
