package runner

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
)

// LegacyDepositERC20 sends ERC20 into ZEVM using legacy protocol contracts
func (r *E2ERunner) LegacyDepositERC20() ethcommon.Hash {
	r.Logger.Print("⏳ depositing ERC20 into ZEVM")

	oneHundred := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(100))
	return r.LegacyDepositERC20WithAmountAndMessage(r.EVMAddress(), oneHundred, []byte{})
}

// LegacyDepositERC20WithAmountAndMessage sends ERC20 into ZEVM using legacy protocol contracts
func (r *E2ERunner) LegacyDepositERC20WithAmountAndMessage(
	to ethcommon.Address,
	amount *big.Int,
	msg []byte,
) ethcommon.Hash {
	// reset allowance, necessary for USDT
	tx, err := r.ERC20.Approve(r.EVMAuth, r.ERC20CustodyAddr, big.NewInt(0))
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt)

	r.Logger.Info("ERC20 Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = r.ERC20.Approve(r.EVMAuth, r.ERC20CustodyAddr, amount)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt)

	r.Logger.Info("ERC20 Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = r.ERC20Custody.Deposit(r.EVMAuth, to.Bytes(), r.ERC20Addr, amount, msg)
	require.NoError(r, err)

	r.Logger.Info("TX: %v", tx)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt)

	r.Logger.Info("Deposit receipt tx hash: %s, status %d", receipt.TxHash.Hex(), receipt.Status)
	for _, log := range receipt.Logs {
		event, err := r.ERC20Custody.ParseDeposited(*log)
		if err != nil {
			continue
		}
		r.Logger.Info("Deposited event:")
		r.Logger.Info("  Recipient address: %x", event.Recipient)
		r.Logger.Info("  ERC20 address: %s", event.Asset.Hex())
		r.Logger.Info("  Amount: %d", event.Amount)
		r.Logger.Info("  Message: %x", event.Message)
	}
	return tx.Hash()
}

// LegacyDepositEther sends Ethers into ZEVM using legacy protocol contracts using legacy protocol contracts
func (r *E2ERunner) LegacyDepositEther() ethcommon.Hash {
	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(100)) // 100 eth
	return r.LegacyDepositEtherWithAmount(amount)
}

// LegacyDepositEtherWithAmount sends Ethers into ZEVM
func (r *E2ERunner) LegacyDepositEtherWithAmount(amount *big.Int) ethcommon.Hash {
	r.Logger.Print("⏳ depositing Ethers into ZEVM")

	signedTx, err := r.LegacySendEther(r.TSSAddress, amount, nil)
	require.NoError(r, err)

	r.Logger.EVMTransaction(signedTx, "send to TSS")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt, "deposit failed")

	r.Logger.EVMReceipt(*receipt, "send to TSS")

	return signedTx.Hash()
}

// LegacySendEther sends ethers to the TSS on EVM using legacy protocol contracts
func (r *E2ERunner) LegacySendEther(_ ethcommon.Address, value *big.Int, data []byte) (*ethtypes.Transaction, error) {
	evmClient := r.EVMClient

	nonce, err := evmClient.PendingNonceAt(r.Ctx, r.EVMAddress())
	if err != nil {
		return nil, err
	}

	gasLimit := uint64(30000) // in units
	gasPrice, err := evmClient.SuggestGasPrice(r.Ctx)
	if err != nil {
		return nil, err
	}

	tx := ethtypes.NewTransaction(nonce, r.TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := evmClient.ChainID(r.Ctx)
	if err != nil {
		return nil, err
	}

	deployerPrivkey, err := r.Account.PrivateKey()
	if err != nil {
		return nil, err
	}

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		return nil, err
	}
	err = evmClient.SendTransaction(r.Ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
