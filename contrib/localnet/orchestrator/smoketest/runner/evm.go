package runner

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/ethereum"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (sm *SmokeTestRunner) DepositERC20() {
	sm.Logger.Print("⏳ depositing ERC20 into ZEVM")
	startTime := time.Now()
	defer func() {
		sm.Logger.Print("✅ ERC20 deposited in %s", time.Since(startTime))
	}()

	initialBal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}

	txHash := sm.DepositERC20WithAmountAndMessage(big.NewInt(1e18), []byte{})
	utils.WaitCctxMinedByInTxHash(txHash.Hex(), sm.CctxClient, sm.Logger)

	// checking balance diff
	bal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}

	diff := big.NewInt(0)
	diff.Sub(bal, initialBal)
	if diff.Int64() != 1e18 {
		panic("balance is not correct")
	}
	sm.Logger.Info("balance of deployer on USDT ZRC20: %d", bal)
}

func (sm *SmokeTestRunner) DepositERC20WithAmountAndMessage(amount *big.Int, msg []byte) ethcommon.Hash {
	USDT := sm.USDTERC20
	tx, err := USDT.Mint(sm.GoerliAuth, amount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	sm.Logger.Info("Mint receipt tx hash: %s", tx.Hash().Hex())

	tx, err = USDT.Approve(sm.GoerliAuth, sm.ERC20CustodyAddr, amount)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	sm.Logger.Info("USDT Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = sm.ERC20Custody.Deposit(sm.GoerliAuth, sm.DeployerAddress.Bytes(), sm.USDTERC20Addr, amount, msg)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
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
func (sm *SmokeTestRunner) DepositEther() {
	sm.Logger.Print("⏳ depositing Ethers into ZEVM")
	startTime := time.Now()
	defer func() {
		sm.Logger.Print("✅ Ethers deposited in %s", time.Since(startTime))
	}()

	ethZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.GoerliLocalnetChain().ChainId))
	if err != nil {
		panic(err)
	}
	if (ethZRC20Addr == ethcommon.Address{}) {
		panic("eth zrc20 not found")
	}
	sm.ETHZRC20Addr = ethZRC20Addr
	sm.Logger.Info("eth zrc20 address: %s", ethZRC20Addr.String())
	ethZRC20, err := zrc20.NewZRC20(ethZRC20Addr, sm.ZevmClient)
	if err != nil {
		panic(err)
	}
	sm.ETHZRC20 = ethZRC20
	initialBalance, err := ethZRC20.BalanceOf(nil, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	signedTx, err := sm.SendEther(sm.TSSAddress, value, nil)
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("GOERLI tx sent: %s; to %s, nonce %d", signedTx.Hash().String(), signedTx.To().Hex(), signedTx.Nonce())
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, signedTx, sm.Logger)
	sm.Logger.Info("GOERLI tx receipt: %d", receipt.Status)
	sm.Logger.Info("  tx hash: %s", receipt.TxHash.String())
	sm.Logger.Info("  to: %s", signedTx.To().String())
	sm.Logger.Info("  value: %d", signedTx.Value())
	sm.Logger.Info("  block num: %d", receipt.BlockNumber)

	{
		sm.Logger.InfoLoud("Merkle Proof\n")
		txHash := receipt.TxHash
		blockHash := receipt.BlockHash

		// #nosec G701 smoketest - always in range
		txIndex := int(receipt.TransactionIndex)

		block, err := sm.GoerliClient.BlockByHash(context.Background(), blockHash)
		if err != nil {
			panic(err)
		}
		i := 0
		for {
			if i > 20 {
				panic("block header not found")
			}
			_, err := sm.ObserverClient.GetBlockHeaderByHash(context.Background(), &observertypes.QueryGetBlockHeaderByHashRequest{
				BlockHash: blockHash.Bytes(),
			})
			if err != nil {
				sm.Logger.Info("WARN: block header not found; retrying... error: %s", err.Error())
				time.Sleep(5 * time.Second)
			} else {
				sm.Logger.Info("OK: block header found")
				break
			}
			i++
		}

		trie := ethereum.NewTrie(block.Transactions())
		if trie.Hash() != block.Header().TxHash {
			panic("tx root hash & block tx root mismatch")
		}
		txProof, err := trie.GenerateProof(txIndex)
		if err != nil {
			panic("error generating txProof")
		}
		val, err := txProof.Verify(block.TxHash(), txIndex)
		if err != nil {
			panic("error verifying txProof")
		}
		var txx ethtypes.Transaction
		err = txx.UnmarshalBinary(val)
		if err != nil {
			panic("error unmarshalling txProof'd tx")
		}
		res, err := sm.ObserverClient.Prove(context.Background(), &observertypes.QueryProveRequest{
			BlockHash: blockHash.Hex(),
			TxIndex:   int64(txIndex),
			TxHash:    txHash.Hex(),
			Proof:     common.NewEthereumProof(txProof),
			ChainId:   common.GoerliLocalnetChain().ChainId,
		})
		if err != nil {
			panic(err)
		}
		if !res.Valid {
			panic("txProof invalid") // FIXME: don't do this in production
		}
		sm.Logger.Info("OK: txProof verified")
	}

	c := make(chan any)
	sm.WG.Add(1)
	go func() {
		defer sm.WG.Done()
		cctx := utils.WaitCctxMinedByInTxHash(signedTx.Hash().Hex(), sm.CctxClient, sm.Logger)
		if cctx.CctxStatus.Status != types.CctxStatus_OutboundMined {
			panic(fmt.Sprintf("expected cctx status to be mined; got %s, message: %s",
				cctx.CctxStatus.Status.String(),
				cctx.CctxStatus.StatusMessage),
			)
		}
		c <- 0
	}()
	sm.WG.Add(1)
	go func() {
		defer sm.WG.Done()
		<-c

		currentBalance, err := ethZRC20.BalanceOf(nil, sm.DeployerAddress)
		if err != nil {
			panic(err)
		}
		diff := big.NewInt(0)
		diff.Sub(currentBalance, initialBalance)
		sm.Logger.Info("eth zrc20 balance: %s", currentBalance.String())
		if diff.Cmp(value) != 0 {
			sm.Logger.Info("eth zrc20 bal wanted %d, got %d", value, diff)
			panic("bal mismatch")
		}

	}()
	sm.WG.Wait()
}

// SendEther sends ethers to the TSS on Goerli
func (sm *SmokeTestRunner) SendEther(_ ethcommon.Address, value *big.Int, data []byte) (*ethtypes.Transaction, error) {
	goerliClient := sm.GoerliClient

	nonce, err := goerliClient.PendingNonceAt(context.Background(), sm.DeployerAddress)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(30000) // in units
	gasPrice, err := goerliClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	tx := ethtypes.NewTransaction(nonce, sm.TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := goerliClient.NetworkID(context.Background())
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
	err = goerliClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
