package runner

import (
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

// WaitForTxReceiptOnEvm waits for a tx receipt on EVM
func (sm *SmokeTestRunner) WaitForTxReceiptOnEvm(tx *ethtypes.Transaction) {
	defer func() {
		sm.Unlock()
	}()
	sm.Lock()

	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger)
	if receipt.Status != 1 {
		panic("tx failed")
	}
}

// MintUSDTOnEvm mints USDT on EVM
// amountUSDT is a multiple of 1e18
func (sm *SmokeTestRunner) MintUSDTOnEvm(amountUSDT int64) {
	defer func() {
		sm.Unlock()
	}()
	sm.Lock()

	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(amountUSDT))

	tx, err := sm.USDTERC20.Mint(sm.GoerliAuth, amount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger)
	if receipt.Status == 0 {
		panic("mint failed")
	}
	sm.Logger.Info("Mint receipt tx hash: %s", tx.Hash().Hex())
}

// SendUSDTOnEvm sends USDT to an address on EVM
// this allows the USDT contract deployer to funds other accounts on EVM
// amountUSDT is a multiple of 1e18
func (sm *SmokeTestRunner) SendUSDTOnEvm(address ethcommon.Address, amountUSDT int64) *ethtypes.Transaction {
	// the deployer might be sending USDT in different goroutines
	defer func() {
		sm.Unlock()
	}()
	sm.Lock()

	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(amountUSDT))

	// transfer
	tx, err := sm.USDTERC20.Transfer(sm.GoerliAuth, address, amount)
	if err != nil {
		panic(err)
	}
	return tx
}

func (sm *SmokeTestRunner) DepositERC20() ethcommon.Hash {
	sm.Logger.Print("⏳ depositing ERC20 into ZEVM")
	startTime := time.Now()
	defer func() {
		sm.Logger.Print("✅ ERC20 deposited in %s", time.Since(startTime))
	}()

	return sm.DepositERC20WithAmountAndMessage(big.NewInt(1e18), []byte{})
}

func (sm *SmokeTestRunner) DepositERC20WithAmountAndMessage(amount *big.Int, msg []byte) ethcommon.Hash {
	tx, err := sm.USDTERC20.Approve(sm.GoerliAuth, sm.ERC20CustodyAddr, amount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	sm.Logger.Info("USDT Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = sm.ERC20Custody.Deposit(sm.GoerliAuth, sm.DeployerAddress.Bytes(), sm.USDTERC20Addr, amount, msg)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger)
	if receipt.Status == 0 {
		panic("deposit failed")
	}
	sm.Logger.Info("Deposit receipt tx hash: %s, status %d", receipt.TxHash.Hex(), receipt.Status)
	for _, log := range receipt.Logs {
		event, err := sm.ERC20Custody.ParseDeposited(*log)
		if err != nil {
			continue
		}
		sm.Logger.Info("Deposited event:")
		sm.Logger.Info("  Recipient address: %x", event.Recipient)
		sm.Logger.Info("  ERC20 address: %s", event.Asset.Hex())
		sm.Logger.Info("  Amount: %d", event.Amount)
		sm.Logger.Info("  Message: %x", event.Message)
	}
	sm.Logger.Info("gas limit %d", sm.ZevmAuth.GasLimit)
	return tx.Hash()
}

// DepositEther sends Ethers into ZEVM
func (sm *SmokeTestRunner) DepositEther() ethcommon.Hash {
	sm.Logger.Print("⏳ depositing Ethers into ZEVM")
	startTime := time.Now()
	defer func() {
		sm.Logger.Print("✅ Ethers deposited in %s", time.Since(startTime))
	}()

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	signedTx, err := sm.SendEther(sm.TSSAddress, value, nil)
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("GOERLI tx sent: %s; to %s, nonce %d", signedTx.Hash().String(), signedTx.To().Hex(), signedTx.Nonce())
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, signedTx, sm.Logger)
	if receipt.Status == 0 {
		panic("deposit failed")
	}
	sm.Logger.Info("GOERLI tx receipt: %d", receipt.Status)
	sm.Logger.Info("  tx hash: %s", receipt.TxHash.String())
	sm.Logger.Info("  to: %s", signedTx.To().String())
	sm.Logger.Info("  value: %d", signedTx.Value())
	sm.Logger.Info("  block num: %d", receipt.BlockNumber)

	return signedTx.Hash()
}

// SendEther sends ethers to the TSS on Goerli
func (sm *SmokeTestRunner) SendEther(_ ethcommon.Address, value *big.Int, data []byte) (*ethtypes.Transaction, error) {
	goerliClient := sm.GoerliClient

	nonce, err := goerliClient.PendingNonceAt(sm.Ctx, sm.DeployerAddress)
	if err != nil {
		return nil, err
	}

	gasLimit := uint64(30000) // in units
	gasPrice, err := goerliClient.SuggestGasPrice(sm.Ctx)
	if err != nil {
		return nil, err
	}

	tx := ethtypes.NewTransaction(nonce, sm.TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := goerliClient.NetworkID(sm.Ctx)
	if err != nil {
		return nil, err
	}

	deployerPrivkey, err := crypto.HexToECDSA(sm.DeployerPrivateKey)
	if err != nil {
		return nil, err
	}

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		return nil, err
	}
	err = goerliClient.SendTransaction(sm.Ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
