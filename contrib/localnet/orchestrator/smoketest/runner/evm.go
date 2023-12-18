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
	"github.com/zeta-chain/zetacore/zetaclient"
)

func (sm *SmokeTestRunner) DepositERC20(amount *big.Int, msg []byte) ethcommon.Hash {
	USDT := sm.USDTERC20
	tx, err := USDT.Mint(sm.GoerliAuth, amount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	fmt.Printf("Mint receipt tx hash: %s\n", tx.Hash().Hex())

	tx, err = USDT.Approve(sm.GoerliAuth, sm.ERC20CustodyAddr, amount)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	fmt.Printf("USDT Approve receipt tx hash: %s\n", tx.Hash().Hex())

	tx, err = sm.ERC20Custody.Deposit(sm.GoerliAuth, sm.DeployerAddress.Bytes(), sm.USDTERC20Addr, amount, msg)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	if receipt.Status == 0 {
		panic("deposit failed")
	}
	fmt.Printf("Deposit receipt tx hash: %s, status %d\n", receipt.TxHash.Hex(), receipt.Status)
	for _, log := range receipt.Logs {
		event, err := sm.ERC20Custody.ParseDeposited(*log)
		if err != nil {
			continue
		}
		fmt.Printf("Deposited event: \n")
		fmt.Printf("  Recipient address: %x, \n", event.Recipient)
		fmt.Printf("  ERC20 address: %s, \n", event.Asset.Hex())
		fmt.Printf("  Amount: %d, \n", event.Amount)
		fmt.Printf("  Message: %x, \n", event.Message)
	}
	fmt.Printf("gas limit %d\n", sm.ZevmAuth.GasLimit)
	return tx.Hash()
}

// DepositEtherIntoZRC20 sends Ethers into ZEVM
func (sm *SmokeTestRunner) DepositEtherIntoZRC20() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	goerliClient := sm.GoerliClient
	utils.LoudPrintf("Deposit Ether into ZEVM\n")
	bn, err := goerliClient.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI block number: %d\n", bn)
	bal, err := goerliClient.BalanceAt(context.Background(), sm.DeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI deployer balance: %s\n", bal.String())

	systemContract := sm.SystemContract
	if err != nil {
		panic(err)
	}
	ethZRC20Addr, err := systemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.GoerliLocalnetChain().ChainId))
	if err != nil {
		panic(err)
	}
	if (ethZRC20Addr == ethcommon.Address{}) {
		panic("eth zrc20 not found")
	}
	sm.ETHZRC20Addr = ethZRC20Addr
	fmt.Printf("eth zrc20 address: %s\n", ethZRC20Addr.String())
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

	fmt.Printf("GOERLI tx sent: %s; to %s, nonce %d\n", signedTx.Hash().String(), signedTx.To().Hex(), signedTx.Nonce())
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, signedTx)
	fmt.Printf("GOERLI tx receipt: %d\n", receipt.Status)
	fmt.Printf("  tx hash: %s\n", receipt.TxHash.String())
	fmt.Printf("  to: %s\n", signedTx.To().String())
	fmt.Printf("  value: %d\n", signedTx.Value())
	fmt.Printf("  block num: %d\n", receipt.BlockNumber)

	{
		utils.LoudPrintf("Merkle Proof\n")
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
				fmt.Printf("WARN: block header not found; retrying... error: %s \n", err.Error())
				time.Sleep(5 * time.Second)
			} else {
				fmt.Printf("OK: block header found\n")
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
		fmt.Printf("OK: txProof verified\n")
	}

	{
		tx, err := sm.SendEther(sm.TSSAddress, big.NewInt(101000000000000000), []byte(zetaclient.DonationMessage))
		if err != nil {
			panic(err)
		}
		receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
		fmt.Printf("GOERLI donation tx receipt: %d\n", receipt.Status)
	}

	c := make(chan any)
	sm.WG.Add(1)
	go func() {
		defer sm.WG.Done()
		cctx := utils.WaitCctxMinedByInTxHash(signedTx.Hash().Hex(), sm.CctxClient)
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
		fmt.Printf("eth zrc20 balance: %s\n", currentBalance.String())
		if diff.Cmp(value) != 0 {
			fmt.Printf("eth zrc20 bal wanted %d, got %d\n", value, diff)
			panic("bal mismatch")
		}

	}()
	sm.WG.Wait()
}

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
