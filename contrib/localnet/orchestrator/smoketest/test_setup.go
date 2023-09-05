//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/testdapp"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	erc20custody "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	zetaeth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zeta.eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	zrc20 "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	uniswapv2factory "github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/erc20"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

const (
	ContractsConfigFile = "contracts.toml"
)

type Contracts struct {
	ZetaEthAddress   string
	ConnectorEthAddr string
}

func (sm *SmokeTest) TestSetupZetaTokenAndConnectorAndZEVMContracts() {
	contracts := Contracts{}
	if localTestArgs.contractsDeployed {
		sm.setContracts()
		return
	}

	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()

	goerliClient := sm.goerliClient
	auth := sm.getDeployerAuth()

	LoudPrintf("Deploy ZetaETH ConnectorETH ERC20Custody USDT\n")

	initialNonce, err := goerliClient.PendingNonceAt(context.Background(), DeployerAddress)
	if err != nil {
		panic(err)
	}

	contractsDeployed := false
	res, err := sm.fungibleClient.ForeignCoinsAll(context.Background(), &fungibletypes.QueryAllForeignCoinsRequest{})
	if err != nil {
		panic(err)
	}
	if len(res.ForeignCoins) > 0 {
		contractsDeployed = true
	}

	if err := CheckNonce(goerliClient, DeployerAddress, initialNonce); err != nil {
		panic(err)
	}
	zetaEthAddr, tx, ZetaEth, err := zetaeth.DeployZetaEth(auth, goerliClient, DeployerAddress, big.NewInt(21_000_000_000))
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZetaEth contract address: %s, tx hash: %s\n", zetaEthAddr.Hex(), tx.Hash().Hex())
	receipt := MustWaitForTxReceipt(goerliClient, tx)
	fmt.Printf("ZetaEth contract receipt: contract address %s, status %d\n", receipt.ContractAddress, receipt.Status)
	sm.ZetaEth = ZetaEth
	sm.ZetaEthAddr = zetaEthAddr
	contracts.ZetaEthAddress = zetaEthAddr.String()

	if err := CheckNonce(goerliClient, DeployerAddress, initialNonce+1); err != nil {
		panic(err)
	}
	connectorEthAddr, tx, ConnectorEth, err := zetaconnectoreth.DeployZetaConnectorEth(auth, goerliClient, zetaEthAddr,
		TSSAddress, DeployerAddress, DeployerAddress)
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(goerliClient, tx)
	fmt.Printf("ZetaConnectorEth contract address: %s, tx hash: %s\n", connectorEthAddr.Hex(), tx.Hash().Hex())
	fmt.Printf("ZetaConnectorEth contract receipt: contract address %s, status %d\n", receipt.ContractAddress, receipt.Status)
	sm.ConnectorEth = ConnectorEth
	sm.ConnectorEthAddr = connectorEthAddr
	contracts.ConnectorEthAddr = connectorEthAddr.String()

	fungibleClient := sm.fungibleClient
	fmt.Printf("Deploying ERC20Custody contract\n")
	if err := CheckNonce(goerliClient, DeployerAddress, initialNonce+2); err != nil {
		panic(err)
	}
	erc20CustodyAddr, tx, ERC20Custody, err := erc20custody.DeployERC20Custody(auth, goerliClient, DeployerAddress, DeployerAddress, big.NewInt(0), big.NewInt(1e18), ethcommon.HexToAddress("0x"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("ERC20Custody contract address: %s, tx hash: %s\n", erc20CustodyAddr.Hex(), tx.Hash().Hex())
	receipt = MustWaitForTxReceipt(goerliClient, tx)
	fmt.Printf("ERC20Custody contract receipt: contract address %s, status %d\n", receipt.ContractAddress, receipt.Status)
	if erc20CustodyAddr != ethcommon.HexToAddress(ERC20CustodyAddr) {
		panic("ERC20Custody contract address mismatch! check order of tx")
	}

	sm.ERC20CustodyAddr = erc20CustodyAddr
	sm.ERC20Custody = ERC20Custody

	fmt.Printf("Deploying USDT contract\n")
	if err := CheckNonce(goerliClient, DeployerAddress, initialNonce+3); err != nil {
		panic(err)
	}
	usdtAddr, tx, _, err := erc20.DeployUSDT(auth, goerliClient, "USDT", "USDT", 6)
	if err != nil {
		panic(err)
	}
	fmt.Printf("USDT contract address: %s, tx hash: %s\n", usdtAddr.Hex(), tx.Hash().Hex())
	receipt = MustWaitForTxReceipt(goerliClient, tx)
	fmt.Printf("USDT contract receipt: contract address %s, status %d\n", receipt.ContractAddress, receipt.Status)

	if !contractsDeployed {
		if receipt.ContractAddress != ethcommon.HexToAddress(USDTERC20Addr) {
			panic("USDT contract address mismatch! check order of tx")
		}
	}

	fmt.Printf("Step 6: Whitelist USDT\n")
	tx, err = ERC20Custody.Whitelist(auth, usdtAddr)
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(goerliClient, tx)
	fmt.Printf("Whitelist receipt tx hash: %s\n", tx.Hash().Hex())

	fmt.Printf("Step 7: Set TSS address\n")
	tx, err = ERC20Custody.UpdateTSSAddress(auth, TSSAddress)
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(goerliClient, tx)
	fmt.Printf("TSS set receipt tx hash: %s\n", tx.Hash().Hex())

	fmt.Printf("Checking foreign coins...\n")
	res, err = fungibleClient.ForeignCoinsAll(context.Background(), &fungibletypes.QueryAllForeignCoinsRequest{})
	if err != nil {
		panic(err)
	}
	found := false
	zrc20addr := ""
	for _, fcoin := range res.ForeignCoins {
		if fcoin.Asset == USDTERC20Addr {
			found = true
			zrc20addr = fcoin.Zrc20ContractAddress
		}
	}
	if !found {
		fmt.Printf("foreign coins: %v", res.ForeignCoins)
		panic(fmt.Sprintf("fungible module does not have foreign coin that represent USDT ERC20 %s", usdtAddr))
	}
	fmt.Printf("USDT ZRC20 Address: %s\n", zrc20addr)
	if !contractsDeployed {
		if HexToAddress(zrc20addr) != HexToAddress(USDTZRC20Addr) {
			panic("mismatch of foreign coin USDT ZRC20 and the USDTZRC20Addr constant in smoketest")
		}
	}

	sm.USDTZRC20Addr = ethcommon.HexToAddress(zrc20addr)
	sm.USDTZRC20, err = zrc20.NewZRC20(sm.USDTZRC20Addr, sm.zevmClient)
	if err != nil {
		panic(err)
	}

	USDT, err := erc20.NewUSDT(usdtAddr, goerliClient)
	if err != nil {
		panic(err)
	}
	sm.USDTERC20 = USDT
	sm.USDTERC20Addr = usdtAddr
	sm.UniswapV2FactoryAddr = ethcommon.HexToAddress(UniswapV2FactoryAddr)
	sm.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(sm.UniswapV2FactoryAddr, sm.zevmClient)
	if err != nil {
		panic(err)
	}
	sm.UniswapV2RouterAddr = ethcommon.HexToAddress(UniswapV2RouterAddr)
	sm.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(sm.UniswapV2RouterAddr, sm.zevmClient)
	if err != nil {
		panic(err)
	}

	fmt.Printf("UniswapV2FactoryAddr: %s, UniswapV2RouterAddr: %s", sm.UniswapV2FactoryAddr.String(), sm.UniswapV2RouterAddr.String())

	// deploy TestDApp contract
	//auth.GasLimit = 1_000_000
	sm.setupTestDapp(auth)

	// Save contract addresses to toml file
	b, err := toml.Marshal(contracts)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(ContractsConfigFile, b, 0666)
	if err != nil {
		panic(err)
	}
}

// Set existing deployed contracts
func (sm *SmokeTest) setContracts() {
	err := error(nil)
	var contracts Contracts

	// Read contracts toml file
	b, err := os.ReadFile(ContractsConfigFile)
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(b, &contracts)
	if err != nil {
		panic(err)
	}

	//Set ZetaEthAddr
	sm.ZetaEthAddr = ethcommon.HexToAddress(contracts.ZetaEthAddress)
	fmt.Println("Connector Eth address: ", contracts.ZetaEthAddress)
	sm.ZetaEth, err = zetaeth.NewZetaEth(sm.ZetaEthAddr, sm.goerliClient)
	if err != nil {
		panic(err)
	}

	//Set ConnectorEthAddr
	sm.ConnectorEthAddr = ethcommon.HexToAddress(contracts.ConnectorEthAddr)
	sm.ConnectorEth, err = zetaconnectoreth.NewZetaConnectorEth(sm.ConnectorEthAddr, sm.goerliClient)
	if err != nil {
		panic(err)
	}

	//Set ERC20CustodyAddr
	sm.ERC20CustodyAddr = ethcommon.HexToAddress(ERC20CustodyAddr)
	sm.ERC20Custody, err = erc20custody.NewERC20Custody(sm.ERC20CustodyAddr, sm.goerliClient)
	if err != nil {
		panic(err)
	}

	//Set USDTERC20Addr
	sm.USDTERC20Addr = ethcommon.HexToAddress(USDTERC20Addr)
	sm.USDTERC20, err = erc20.NewUSDT(sm.USDTERC20Addr, sm.goerliClient)
	if err != nil {
		panic(err)
	}

	//Set USDTZRC20Addr
	sm.USDTZRC20Addr = ethcommon.HexToAddress(USDTZRC20Addr)
	sm.USDTZRC20, err = zrc20.NewZRC20(sm.USDTZRC20Addr, sm.zevmClient)
	if err != nil {
		panic(err)
	}

	//UniswapV2FactoryAddr
	sm.UniswapV2FactoryAddr = ethcommon.HexToAddress(UniswapV2FactoryAddr)
	sm.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(sm.UniswapV2FactoryAddr, sm.zevmClient)
	if err != nil {
		panic(err)
	}

	//UniswapV2RouterAddr
	sm.UniswapV2RouterAddr = ethcommon.HexToAddress(UniswapV2RouterAddr)
	sm.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(sm.UniswapV2RouterAddr, sm.zevmClient)
	if err != nil {
		panic(err)
	}

	sm.setupTestDapp(sm.getDeployerAuth())
}

func (sm *SmokeTest) setupTestDapp(auth *bind.TransactOpts) {
	// deploy TestDApp contract
	//auth.GasLimit = 1_000_000
	appAddr, tx, _, err := testdapp.DeployTestDApp(auth, sm.goerliClient, sm.ConnectorEthAddr, sm.ZetaEthAddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TestDApp contract address: %s, tx hash: %s\n", appAddr.Hex(), tx.Hash().Hex())
	receipt := MustWaitForTxReceipt(sm.goerliClient, tx)
	fmt.Printf("TestDApp contract receipt: contract address %s, status %d; used gas %d\n", receipt.ContractAddress, receipt.Status, receipt.GasUsed)
	dapp, err := testdapp.NewTestDApp(receipt.ContractAddress, sm.goerliClient)
	if err != nil {
		panic(err)
	}
	{
		code, err := sm.goerliClient.CodeAt(context.Background(), receipt.ContractAddress, nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("TestDApp contract code: len %d\n", len(code))
		if len(code) == 0 {
			panic("TestDApp contract code is empty")
		}
		res, err := dapp.Connector(&bind.CallOpts{})
		if err != nil {
			panic(err)
		}
		if res != sm.ConnectorEthAddr {
			panic("mismatch of TestDApp connector address")
		}
	}
	sm.TestDAppAddr = receipt.ContractAddress
}

func (sm *SmokeTest) getDeployerAuth() *bind.TransactOpts {
	chainid, err := sm.goerliClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	deployerPrivkey, err := crypto.HexToECDSA(DeployerPrivateKey)
	if err != nil {
		panic(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, chainid)
	if err != nil {
		panic(err)
	}
	return auth
}
