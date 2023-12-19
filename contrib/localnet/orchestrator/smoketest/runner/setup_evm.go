package runner

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	zetaeth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zeta.eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/erc20"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/testdapp"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

const (
	ContractsConfigFile = "contracts.toml"
)

// SetEVMContractsFromConfig set EVM contracts for smoke test from the config
func (sm *SmokeTestRunner) SetEVMContractsFromConfig() {
	conf, err := config.ReadConfig(ContractsConfigFile)
	if err != nil {
		panic(err)
	}

	// Set ZetaEthAddr
	sm.ZetaEthAddr = ethcommon.HexToAddress(conf.Contracts.EVM.ZetaEthAddress)
	sm.ZetaEth, err = zetaeth.NewZetaEth(sm.ZetaEthAddr, sm.GoerliClient)
	if err != nil {
		panic(err)
	}

	// Set ConnectorEthAddr
	sm.ConnectorEthAddr = ethcommon.HexToAddress(conf.Contracts.EVM.ConnectorEthAddr)
	sm.ConnectorEth, err = zetaconnectoreth.NewZetaConnectorEth(sm.ConnectorEthAddr, sm.GoerliClient)
	if err != nil {
		panic(err)
	}
}

// SetupEVM setup contracts on EVM for smoke test
func (sm *SmokeTestRunner) SetupEVM(contractsDeployed bool) {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()

	// TODO: put this logic outside of this function
	// We use this config to be consistent with the old implementation
	// https://github.com/zeta-chain/node-private/issues/41
	if contractsDeployed {
		sm.SetEVMContractsFromConfig()
		return
	}
	conf := config.DefaultConfig()

	utils.LoudPrintf("Deploy ZetaETH ConnectorETH ERC20Custody USDT\n")

	// fetch initial nonce to check if it get incremented correctly
	initialNonce, err := sm.GoerliClient.PendingNonceAt(context.Background(), sm.DeployerAddress)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deploying ZetaEth contract\n")
	zetaEthAddr, tx, ZetaEth, err := zetaeth.DeployZetaEth(sm.GoerliAuth, sm.GoerliClient, sm.DeployerAddress, big.NewInt(21_000_000_000))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	if receipt.Status != 1 {
		panic("ZetaEth deployment failed")
	}
	sm.ZetaEth = ZetaEth
	sm.ZetaEthAddr = zetaEthAddr
	conf.Contracts.EVM.ZetaEthAddress = zetaEthAddr.String()
	if err := utils.CheckNonce(sm.GoerliClient, sm.DeployerAddress, initialNonce+1); err != nil {
		panic(err)
	}
	fmt.Printf("ZetaEth contract address: %s, tx hash: %s\n", zetaEthAddr.Hex(), tx.Hash().Hex())

	fmt.Printf("Deploying ZetaConnectorEth contract\n")
	connectorEthAddr, tx, ConnectorEth, err := zetaconnectoreth.DeployZetaConnectorEth(
		sm.GoerliAuth,
		sm.GoerliClient,
		zetaEthAddr,
		sm.TSSAddress,
		sm.DeployerAddress,
		sm.DeployerAddress,
	)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	if receipt.Status != 1 {
		panic("ZetaConnectorEth deployment failed")
	}
	sm.ConnectorEth = ConnectorEth
	sm.ConnectorEthAddr = connectorEthAddr
	conf.Contracts.EVM.ConnectorEthAddr = connectorEthAddr.String()

	if err := utils.CheckNonce(sm.GoerliClient, sm.DeployerAddress, initialNonce+2); err != nil {
		panic(err)
	}
	fmt.Printf("ZetaConnectorEth contract address: %s, tx hash: %s\n", connectorEthAddr.Hex(), tx.Hash().Hex())

	fmt.Printf("Deploying ERC20Custody contract\n")
	erc20CustodyAddr, tx, ERC20Custody, err := erc20custody.DeployERC20Custody(
		sm.GoerliAuth,
		sm.GoerliClient,
		sm.DeployerAddress,
		sm.DeployerAddress,
		big.NewInt(0),
		big.NewInt(1e18),
		ethcommon.HexToAddress("0x"),
	)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	if receipt.Status != 1 {
		panic("ERC20Custody deployment failed")
	}
	sm.ERC20CustodyAddr = erc20CustodyAddr
	sm.ERC20Custody = ERC20Custody
	if err := utils.CheckNonce(sm.GoerliClient, sm.DeployerAddress, initialNonce+3); err != nil {
		panic(err)
	}
	fmt.Printf("ERC20Custody contract address: %s, tx hash: %s\n", erc20CustodyAddr.Hex(), tx.Hash().Hex())

	fmt.Printf("Deploying USDT contract\n")
	usdtAddr, tx, usdt, err := erc20.DeployUSDT(sm.GoerliAuth, sm.GoerliClient, "USDT", "USDT", 6)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	if receipt.Status != 1 {
		panic("USDT deployment failed")
	}
	sm.USDTERC20 = usdt
	sm.USDTERC20Addr = usdtAddr
	fmt.Printf("USDT contract address: %s, tx hash: %s\n", usdtAddr.Hex(), tx.Hash().Hex())

	fmt.Printf("Whitelist USDT\n")
	tx, err = ERC20Custody.Whitelist(sm.GoerliAuth, usdtAddr)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	if receipt.Status != 1 {
		panic("USDT whitelist failed")
	}

	fmt.Printf("Set TSS address\n")
	tx, err = ERC20Custody.UpdateTSSAddress(sm.GoerliAuth, sm.TSSAddress)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	if receipt.Status != 1 {
		panic("USDT update TSS address failed")
	}
	fmt.Printf("TSS set receipt tx hash: %s\n", tx.Hash().Hex())

	// deploy TestDApp contract
	sm.setupTestDapp()

	// save config containing contract addresses
	// TODO: put this logic outside of this function in a general config
	// We use this config to be consistent with the old implementation
	// https://github.com/zeta-chain/node-private/issues/41
	if err := config.WriteConfig(ContractsConfigFile, conf); err != nil {
		panic(err)
	}
}

// setupTestDapp deploys TestDApp contract
func (sm *SmokeTestRunner) setupTestDapp() {
	// deploy TestDApp contract
	appAddr, tx, _, err := testdapp.DeployTestDApp(sm.GoerliAuth, sm.GoerliClient, sm.ConnectorEthAddr, sm.ZetaEthAddr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("TestDApp contract address: %s, tx hash: %s\n", appAddr.Hex(), tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	if receipt.Status != 1 {
		panic("TestDApp deployment failed")
	}

	dapp, err := testdapp.NewTestDApp(receipt.ContractAddress, sm.GoerliClient)
	if err != nil {
		panic(err)
	}

	// check contract code
	code, err := sm.GoerliClient.CodeAt(context.Background(), receipt.ContractAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TestDApp contract code: len %d\n", len(code))
	if len(code) == 0 {
		panic("TestDApp contract code is empty")
	}

	// check connector deployed
	res, err := dapp.Connector(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if res != sm.ConnectorEthAddr {
		panic("mismatch of TestDApp connector address")
	}

	sm.TestDAppAddr = receipt.ContractAddress
}
