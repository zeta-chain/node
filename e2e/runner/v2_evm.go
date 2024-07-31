package runner

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

// V2ETHDeposit calls Deposit of Gateway with gas token on EVM
func (r *E2ERunner) V2ETHDeposit(receiver ethcommon.Address, amount *big.Int) *ethtypes.Transaction {
	// set the value of the transaction
	previousValue := r.EVMAuth.Value
	defer func() {
		r.EVMAuth.Value = previousValue
	}()
	r.EVMAuth.Value = amount

	tx, err := r.GatewayEVM.Deposit(r.EVMAuth, receiver)
	require.NoError(r, err)

	return tx
}

// V2ETHDepositAndCall calls DepositAndCall of Gateway with gas token on EVM
func (r *E2ERunner) V2ETHDepositAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
) *ethtypes.Transaction {
	// set the value of the transaction
	previousValue := r.EVMAuth.Value
	defer func() {
		r.EVMAuth.Value = previousValue
	}()
	r.EVMAuth.Value = amount

	tx, err := r.GatewayEVM.DepositAndCall(r.EVMAuth, receiver, payload)
	require.NoError(r, err)

	return tx
}

// V2ERC20Deposit calls Deposit of Gateway with erc20 token on EVM
func (r *E2ERunner) V2ERC20Deposit(receiver ethcommon.Address, amount *big.Int) *ethtypes.Transaction {
	tx, err := r.GatewayEVM.Deposit0(r.EVMAuth, receiver, amount, r.ERC20Addr)
	require.NoError(r, err)

	return tx
}

// V2ERC20DepositAndCall calls DepositAndCall of Gateway with erc20 token on EVM
func (r *E2ERunner) V2ERC20DepositAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
) *ethtypes.Transaction {
	tx, err := r.GatewayEVM.DepositAndCall0(r.EVMAuth, receiver, amount, r.ERC20Addr, payload)
	require.NoError(r, err)

	return tx
}

// V2EVMToZEMVCall calls Call of Gateway on EVM
func (r *E2ERunner) V2EVMToZEMVCall(receiver ethcommon.Address, payload []byte) *ethtypes.Transaction {
	tx, err := r.GatewayEVM.Call(r.EVMAuth, receiver, payload)
	require.NoError(r, err)

	return tx
}
