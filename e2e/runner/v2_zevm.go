package runner

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

// V2ETHWithdraw calls Withdraw of Gateway with gas token on ZEVM
func (r *E2ERunner) V2ETHWithdraw(receiver ethcommon.Address, amount *big.Int) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.Withdraw(r.ZEVMAuth, receiver.Bytes(), amount, r.ETHZRC20Addr)
	require.NoError(r, err)

	return tx
}

// V2ETHWithdrawAndCall calls WithdrawAndCall of Gateway with gas token on ZEVM
func (r *E2ERunner) V2ETHWithdrawAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.WithdrawAndCall(r.ZEVMAuth, receiver.Bytes(), amount, r.ETHZRC20Addr, payload)
	require.NoError(r, err)

	return tx
}

// V2ERC20Withdraw calls Withdraw of Gateway with erc20 token on ZEVM
func (r *E2ERunner) V2ERC20Withdraw(receiver ethcommon.Address, amount *big.Int) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.Withdraw(r.ZEVMAuth, receiver.Bytes(), amount, r.ERC20ZRC20Addr)
	require.NoError(r, err)

	return tx
}

// V2ERC20WithdrawAndCall calls WithdrawAndCall of Gateway with erc20 token on ZEVM
func (r *E2ERunner) V2ERC20WithdrawAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
) *ethtypes.Transaction {
	tx, err := r.GatewayZEVM.WithdrawAndCall(r.ZEVMAuth, receiver.Bytes(), amount, r.ERC20ZRC20Addr, payload)
	require.NoError(r, err)

	return tx
}

// V2ZEVMToEMVCall calls Call of Gateway on ZEVM
func (r *E2ERunner) V2ZEVMToEMVCall(receiver ethcommon.Address, payload []byte) *ethtypes.Transaction {
	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	tx, err := r.GatewayZEVM.Call(r.ZEVMAuth, receiver.Bytes(), chainID, payload)
	require.NoError(r, err)

	return tx
}
