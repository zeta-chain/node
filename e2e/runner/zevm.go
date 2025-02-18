package runner

import (
	"fmt"
	"math/big"
	"time"

	"github.com/cenkalti/backoff/v4"
	query "github.com/cosmos/cosmos-sdk/types/query"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/gatewayzevmcaller"
	"github.com/zeta-chain/node/pkg/retry"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

var gasLimit = big.NewInt(250000)

// ETHWithdraw calls Withdraw of Gateway with gas token on ZEVM
func (r *E2ERunner) ETHWithdraw(
	receiver ethcommon.Address,
	amount *big.Int,
	revertOptions gatewayzevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.Withdraw(
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
	tx, err := r.GatewayZEVM.WithdrawAndCall0(
		r.ZEVMAuth,
		receiver.Bytes(),
		amount,
		r.ETHZRC20Addr,
		payload,
		gatewayzevm.CallOptions{GasLimit: gasLimit, IsArbitraryCall: true},
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
) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.WithdrawAndCall0(
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
			GasLimit:        gasLimit,
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
	tx, err := r.GatewayZEVM.Withdraw(
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

	tx, err := r.GatewayZEVM.WithdrawAndCall0(
		r.ZEVMAuth,
		receiver.Bytes(),
		amount,
		r.ERC20ZRC20Addr,
		payload,
		gatewayzevm.CallOptions{GasLimit: gasLimit, IsArbitraryCall: true},
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
) *ethtypes.Transaction {
	// this function take more gas than default 500k
	// so we need to increase the gas limit
	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10000000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	tx, err := r.GatewayZEVM.WithdrawAndCall0(
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
		gatewayzevm.CallOptions{GasLimit: gasLimit, IsArbitraryCall: true},
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
			GasLimit:        gasLimit,
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

	bo := backoff.NewConstantBackOff(time.Second * 5)
	boWithMaxRetries := backoff.WithMaxRetries(bo, 10)
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
		return fmt.Errorf("waiting for %d blocks, current height %d", n, height.Height)
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
			Pagination: &query.PageRequest{Reverse: true},
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
