package runner

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

// V2ETHWithdraw calls Withdraw of Gateway with gas token on ZEVM
func (r *E2ERunner) V2ETHWithdraw(receiver ethcommon.Address, amount *big.Int) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.Withdraw(r.EVMAuth, receiver.Bytes(), amount, r.ETHZRC20Addr)
	require.NoError(r, err)

	return tx
}

// V2ETHWithdrawAndCall calls WithdrawAndCall of Gateway with gas token on ZEVM
func (r *E2ERunner) V2ETHWithdrawAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.WithdrawAndCall0(r.EVMAuth, receiver.Bytes(), amount, r.ETHZRC20Addr, payload)
	require.NoError(r, err)

	return tx
}

// V2ERC20Withdraw calls Withdraw of Gateway with erc20 token on ZEVM
func (r *E2ERunner) V2ERC20Withdraw(receiver ethcommon.Address, amount *big.Int) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.Withdraw(r.EVMAuth, receiver.Bytes(), amount, r.ERC20Addr)
	require.NoError(r, err)

	return tx
}

// V2ERC20WithdrawAndCall calls WithdrawAndCall of Gateway with erc20 token on ZEVM
func (r *E2ERunner) V2ERC20WithdrawAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.WithdrawAndCall0(r.EVMAuth, receiver.Bytes(), amount, r.ERC20Addr, payload)
	require.NoError(r, err)

	return tx
}

// V2ZEVMToEMVCall calls Call of Gateway on ZEVM
func (r *E2ERunner) V2ZEVMToEMVCall(receiver ethcommon.Address, payload []byte) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.Call(r.EVMAuth, receiver.Bytes(), payload)
	require.NoError(r, err)

	return tx
}
