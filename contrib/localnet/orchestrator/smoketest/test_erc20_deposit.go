package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/contracts/evm/erc20"
	"github.com/zeta-chain/zetacore/contracts/evm/erc20custody"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	"math/big"
	"time"
)

func TestERC20Deposit(goerliClient *ethclient.Client, zevmClient *ethclient.Client, cctxClient cctxtypes.QueryClient, fungibleClient fungibletypes.QueryClient) {
	LoudPrintf("Deposit USDT ERC20 into ZEVM\n")
	chainID, err := goerliClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	deployerPrivkey, err := crypto.HexToECDSA(DeployerPrivateKey)
	if err != nil {
		panic(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, chainID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Step 4: Deploying ERC20Custody contract\n")
	erc20CustodyAddr, tx, ERC20Custody, err := erc20custody.DeployERC20Custody(auth, goerliClient, DeployerAddress, DeployerAddress, big.NewInt(0), ethcommon.HexToAddress("0x"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("ERC20Custody contract address: %s, tx hash: %s\n", erc20CustodyAddr.Hex(), tx.Hash().Hex())
	time.Sleep(BLOCK)
	receipt, err := goerliClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("ERC20Custody contract receipt: contract address %s, status %d\n", receipt.ContractAddress, receipt.Status)
	if receipt.ContractAddress != ethcommon.HexToAddress(ERC20CustodyAddr) {
		panic("ERC20Custody contract address mismatch! check order of tx")
	}
	fmt.Printf("Step 5: Deploying USDT contract\n")
	usdtAddr, tx, _, err := erc20.DeployUSDT(auth, goerliClient, "USDT", "USDT", 6)
	if err != nil {
		panic(err)
	}
	fmt.Printf("USDT contract address: %s, tx hash: %s\n", usdtAddr.Hex(), tx.Hash().Hex())
	time.Sleep(BLOCK)
	receipt, err = goerliClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("USDT contract receipt: contract address %s, status %d\n", receipt.ContractAddress, receipt.Status)
	if receipt.ContractAddress != ethcommon.HexToAddress(USDTERC20Addr) {
		panic("USDT contract address mismatch! check order of tx")
	}
	fmt.Printf("Step 6: Whitelist USDT\n")
	tx, err = ERC20Custody.Whitelist(auth, usdtAddr)
	if err != nil {
		panic(err)
	}
	time.Sleep(BLOCK)
	receipt, err = goerliClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Whitelist receipt tx hash: %s\n", tx.Hash().Hex())

	fmt.Printf("Step 7: Set TSS address\n")
	tx, err = ERC20Custody.UpdateTSSAddress(auth, TSSAddress)
	if err != nil {
		panic(err)
	}
	time.Sleep(BLOCK)
	receipt, err = goerliClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("TSS set receipt tx hash: %s\n", tx.Hash().Hex())

	USDT, err := erc20.NewUSDT(usdtAddr, goerliClient)
	if err != nil {
		panic(err)
	}
	tx, err = USDT.Mint(auth, big.NewInt(1e10))
	if err != nil {
		panic(err)
	}
	time.Sleep(BLOCK)
	receipt, err = goerliClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Mint receipt tx hash: %s\n", tx.Hash().Hex())

	tx, err = USDT.Approve(auth, erc20CustodyAddr, big.NewInt(1e10))
	if err != nil {
		panic(err)
	}
	time.Sleep(BLOCK)
	receipt, err = goerliClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("USDT Approve receipt tx hash: %s\n", tx.Hash().Hex())

	tx, err = ERC20Custody.Deposit(auth, DeployerAddress.Bytes(), usdtAddr, big.NewInt(1e6), nil)
	if err != nil {
		panic(err)
	}
	time.Sleep(BLOCK)
	receipt, err = goerliClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deposit receipt tx hash: %s, status %d\n", receipt.TxHash.Hex(), receipt.Status)
	for _, log := range receipt.Logs {
		event, err := ERC20Custody.ParseDeposited(*log)
		if err != nil {
			continue
		}
		fmt.Printf("Deposited event: \n")
		fmt.Printf("  Recipient address: %x, \n", event.Recipient)
		fmt.Printf("  ERC20 address: %s, \n", event.Asset.Hex())
		fmt.Printf("  Amount: %d, \n", event.Amount)
		fmt.Printf("  Message: %x, \n", event.Message)
	}
	WaitCctxMinedByInTxHash(tx.Hash().Hex(), cctxClient)

}
