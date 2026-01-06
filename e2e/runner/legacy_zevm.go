package runner

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.eth.sol"
	connectorzevm "github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnectorzevm.sol"

	"github.com/zeta-chain/node/e2e/utils"
)

// LegacyDepositZeta deposits ZETA on ZetaChain from the ZETA smart contract on EVM using legacy protocol contracts
func (r *E2ERunner) LegacyDepositZeta() ethcommon.Hash {
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(100)) // 100 Zeta

	return r.LegacyDepositZetaWithAmount(r.EVMAddress(), amount)
}

// LegacyDepositZetaWithAmountAndPayload deposits ZETA on ZetaChain from the ZETA smart contract on EVM with the specified amount and payload using legacy protocol contracts
func (r *E2ERunner) LegacyDepositZetaWithAmountAndPayload(
	to ethcommon.Address,
	amount *big.Int,
	payload []byte,
) ethcommon.Hash {
	tx, err := r.ZetaEth.Approve(r.EVMAuth, r.ConnectorEthAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "approve")
	r.requireTxSuccessful(receipt, "approve tx failed")

	// query the chain ID using zevm client
	zetaChainID, err := r.ZEVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	paused, err := r.ConnectorEth.Paused(&bind.CallOpts{})
	require.NoError(r, err)
	require.False(r, paused, "ZetaConnectorEth is paused, cannot send ZETA")

	tx, err = r.ConnectorEth.Send(r.EVMAuth, zetaconnectoreth.ZetaInterfacesSendInput{
		// TODO: allow user to specify destination chain id
		// https://github.com/zeta-chain/node-private/issues/41
		DestinationChainId:  zetaChainID,
		DestinationAddress:  to.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             payload,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	require.NoError(r, err)

	r.Logger.Info("Send tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "send")
	r.requireTxSuccessful(receipt, "send tx failed")

	r.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := r.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			r.Logger.Info("    Connector: %s", r.ConnectorEthAddr.String())
			r.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			r.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			r.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			r.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
			r.Logger.Info("    Block Num: %d", log.BlockNumber)
		}
	}

	return tx.Hash()
}

// LegacyDepositZetaWithAmount deposits ZETA on ZetaChain from the ZETA smart contract on EVM with the specified amount using legacy protocol contracts
func (r *E2ERunner) LegacyDepositZetaWithAmount(to ethcommon.Address, amount *big.Int) ethcommon.Hash {
	tx, err := r.ZetaEth.Approve(r.EVMAuth, r.ConnectorEthAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "approve")
	r.requireTxSuccessful(receipt, "approve tx failed")

	// query the chain ID using zevm client
	zetaChainID, err := r.ZEVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	paused, err := r.ConnectorEth.Paused(&bind.CallOpts{})
	require.NoError(r, err)
	require.False(r, paused, "ZetaConnectorEth is paused, cannot send ZETA")

	tx, err = r.ConnectorEth.Send(r.EVMAuth, zetaconnectoreth.ZetaInterfacesSendInput{
		// TODO: allow user to specify destination chain id
		// https://github.com/zeta-chain/node-private/issues/41
		DestinationChainId:  zetaChainID,
		DestinationAddress:  to.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	require.NoError(r, err)

	r.Logger.Info("Send tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "send")
	r.requireTxSuccessful(receipt, "send tx failed")

	r.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := r.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			r.Logger.Info("    Connector: %s", r.ConnectorEthAddr.String())
			r.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			r.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			r.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			r.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
			r.Logger.Info("    Block Num: %d", log.BlockNumber)
		}
	}

	return tx.Hash()
}

// LegacyDepositAndApproveWZeta deposits and approves WZETA on ZetaChain from the ZETA smart contract on ZEVM using legacy protocol contracts
func (r *E2ERunner) LegacyDepositAndApproveWZeta(amount *big.Int) {
	r.ZEVMAuth.Value = amount
	tx, err := r.WZeta.Deposit(r.ZEVMAuth)
	require.NoError(r, err)

	r.ZEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("wzeta deposit tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta deposit")
	r.requireTxSuccessful(receipt, "deposit failed")

	tx, err = r.WZeta.Approve(r.ZEVMAuth, r.ConnectorZEVMAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("wzeta approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta approve")
	r.requireTxSuccessful(receipt, "approve failed, logs: %+v", receipt.Logs)
}

// DepositWZeta deposits WZETA on ZetaChain
func (r *E2ERunner) DepositWZeta(amount *big.Int) {
	r.ZEVMAuth.Value = amount
	tx, err := r.WZeta.Deposit(r.ZEVMAuth)
	require.NoError(r, err)

	r.ZEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("wzeta deposit tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta deposit")
	r.requireTxSuccessful(receipt, "deposit failed")
}

// LegacyWithdrawZeta withdraws ZETA from ZetaChain to the ZETA smart contract on EVM using legacy protocol contracts
// waitReceipt specifies whether to wait for the tx receipt and check if the tx was successful
func (r *E2ERunner) LegacyWithdrawZeta(amount *big.Int, waitReceipt bool) *ethtypes.Transaction {
	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	tx, err := r.ConnectorZEVM.Send(r.ZEVMAuth, connectorzevm.ZetaInterfacesSendInput{
		DestinationChainId:  chainID,
		DestinationAddress:  r.EVMAddress().Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	require.NoError(r, err)

	r.Logger.Info("send tx hash: %s", tx.Hash().Hex())

	if waitReceipt {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.Logger.EVMReceipt(*receipt, "send")
		r.requireTxSuccessful(receipt, "send failed, logs: %+v", receipt.Logs)

		r.Logger.Info("  Logs:")
		for _, log := range receipt.Logs {
			sentLog, err := r.ConnectorZEVM.ParseZetaSent(*log)
			if err == nil {
				r.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
				r.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
				r.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
				r.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
			}
		}
	}

	return tx
}

// LegacyWithdrawEther withdraws Ether from ZetaChain to the ZETA smart contract on EVM using legacy protocol contracts
func (r *E2ERunner) LegacyWithdrawEther(amount *big.Int) *ethtypes.Transaction {
	// withdraw
	tx, err := r.ETHZRC20.Withdraw(r.ZEVMAuth, r.EVMAddress().Bytes(), amount)
	require.NoError(r, err)

	r.Logger.EVMTransaction(tx, "withdraw")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt, "withdraw failed")

	r.Logger.EVMReceipt(*receipt, "withdraw")
	r.Logger.ZRC20Withdrawal(r.ETHZRC20, *receipt, "withdraw")

	return tx
}

// LegacyWithdrawERC20 withdraws an ERC20 token from ZetaChain to the ZETA smart contract on EVM using legacy protocol contracts
func (r *E2ERunner) LegacyWithdrawERC20(amount *big.Int) *ethtypes.Transaction {
	tx, err := r.ERC20ZRC20.Withdraw(r.ZEVMAuth, r.EVMAddress().Bytes(), amount)
	require.NoError(r, err)

	r.Logger.EVMTransaction(tx, "withdraw")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)
	for _, log := range receipt.Logs {
		event, err := r.ERC20ZRC20.ParseWithdrawal(*log)
		if err != nil {
			continue
		}
		r.Logger.Info(
			"  logs: from %s, to %x, value %d, gasfee %d",
			event.From.Hex(),
			event.To,
			event.Value,
			event.GasFee,
		)
	}

	return tx
}
