package runner

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/utils"
)

// V2ETHDeposit calls Deposit of Gateway with gas token on EVM
func (r *E2ERunner) V2ETHDeposit(
	receiver ethcommon.Address,
	amount *big.Int,
	revertOptions gatewayevm.RevertOptions,
) *ethtypes.Transaction {
	// set the value of the transaction
	previousValue := r.EVMAuth.Value
	defer func() {
		r.EVMAuth.Value = previousValue
	}()
	r.EVMAuth.Value = amount

	tx, err := r.GatewayEVM.Deposit0(r.EVMAuth, receiver, revertOptions)
	require.NoError(r, err)

	logDepositInfoAndWaitForTxReceipt(r, tx, "eth_deposit")

	return tx
}

// V2ETHDepositAndCall calls DepositAndCall of Gateway with gas token on EVM
func (r *E2ERunner) V2ETHDepositAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewayevm.RevertOptions,
) *ethtypes.Transaction {
	// set the value of the transaction
	previousValue := r.EVMAuth.Value
	defer func() {
		r.EVMAuth.Value = previousValue
	}()
	r.EVMAuth.Value = amount

	tx, err := r.GatewayEVM.DepositAndCall(r.EVMAuth, receiver, payload, revertOptions)
	require.NoError(r, err)

	logDepositInfoAndWaitForTxReceipt(r, tx, "eth_deposit_and_call")

	return tx
}

// V2ERC20Deposit calls Deposit of Gateway with erc20 token on EVM
func (r *E2ERunner) V2ERC20Deposit(
	receiver ethcommon.Address,
	amount *big.Int,
	revertOptions gatewayevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayEVM.Deposit(r.EVMAuth, receiver, amount, r.ERC20Addr, revertOptions)
	require.NoError(r, err)

	logDepositInfoAndWaitForTxReceipt(r, tx, "erc20_deposit")

	return tx
}

// V2ERC20DepositAndCall calls DepositAndCall of Gateway with erc20 token on EVM
func (r *E2ERunner) V2ERC20DepositAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewayevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayEVM.DepositAndCall0(
		r.EVMAuth,
		receiver,
		amount,
		r.ERC20Addr,
		payload,
		revertOptions,
	)
	require.NoError(r, err)

	logDepositInfoAndWaitForTxReceipt(r, tx, "erc20_deposit_and_call")

	return tx
}

// V2EVMToZEMVCall calls Call of Gateway on EVM
func (r *E2ERunner) V2EVMToZEMVCall(
	receiver ethcommon.Address,
	payload []byte,
	revertOptions gatewayevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayEVM.Call(r.EVMAuth, receiver, payload, revertOptions)
	require.NoError(r, err)

	return tx
}

func logDepositInfoAndWaitForTxReceipt(
	r *E2ERunner,
	tx *ethtypes.Transaction,
	name string,
) {
	r.Logger.EVMTransaction(*tx, name)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt, name+" failed")

	r.Logger.EVMReceipt(*receipt, name)
	r.Logger.GatewayDeposit(r.GatewayEVM, *receipt, name)
}
