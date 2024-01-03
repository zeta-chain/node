package runner

import (
	"math/big"
	"time"

	"github.com/zeta-chain/zetacore/zetaclient"

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
	sm.Logger.Print("⚙️ setting up Goerli network")
	startTime := time.Now()
	defer func() {
		sm.Logger.Info("EVM setup took %s\n", time.Since(startTime))
	}()

	// TODO: put this logic outside of this function
	// We use this config to be consistent with the old implementation
	// https://github.com/zeta-chain/node-private/issues/41
	if contractsDeployed {
		sm.SetEVMContractsFromConfig()
		return
	}
	conf := config.DefaultConfig()

	sm.Logger.InfoLoud("Deploy ZetaETH ConnectorETH ERC20Custody USDT\n")

	// donate to the TSS address to avoid account errors because deploying gas token ZRC20 will automatically mint
	// gas token on ZetaChain to initialize the pool
	txDonation, err := sm.SendEther(sm.TSSAddress, big.NewInt(101000000000000000), []byte(zetaclient.DonationMessage))
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("Deploying ZetaEth contract")
	zetaEthAddr, txZetaEth, ZetaEth, err := zetaeth.DeployZetaEth(
		sm.GoerliAuth,
		sm.GoerliClient,
		sm.DeployerAddress,
		big.NewInt(21_000_000_000),
	)
	if err != nil {
		panic(err)
	}
	sm.ZetaEth = ZetaEth
	sm.ZetaEthAddr = zetaEthAddr
	conf.Contracts.EVM.ZetaEthAddress = zetaEthAddr.String()
	sm.Logger.Info("ZetaEth contract address: %s, tx hash: %s", zetaEthAddr.Hex(), zetaEthAddr.Hash().Hex())

	sm.Logger.Info("Deploying ZetaConnectorEth contract")
	connectorEthAddr, txConnector, ConnectorEth, err := zetaconnectoreth.DeployZetaConnectorEth(
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
	sm.ConnectorEth = ConnectorEth
	sm.ConnectorEthAddr = connectorEthAddr
	conf.Contracts.EVM.ConnectorEthAddr = connectorEthAddr.String()

	sm.Logger.Info("ZetaConnectorEth contract address: %s, tx hash: %s", connectorEthAddr.Hex(), txConnector.Hash().Hex())

	sm.Logger.Info("Deploying ERC20Custody contract")
	erc20CustodyAddr, txCustody, ERC20Custody, err := erc20custody.DeployERC20Custody(
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
	sm.ERC20CustodyAddr = erc20CustodyAddr
	sm.ERC20Custody = ERC20Custody
	sm.Logger.Info("ERC20Custody contract address: %s, tx hash: %s", erc20CustodyAddr.Hex(), txCustody.Hash().Hex())

	sm.Logger.Info("Deploying USDT contract")
	usdtAddr, txUSDT, usdt, err := erc20.DeployUSDT(sm.GoerliAuth, sm.GoerliClient, "USDT", "USDT", 6)
	if err != nil {
		panic(err)
	}
	sm.USDTERC20 = usdt
	sm.USDTERC20Addr = usdtAddr
	sm.Logger.Info("USDT contract address: %s, tx hash: %s", usdtAddr.Hex(), txUSDT.Hash().Hex())

	// deploy TestDApp contract
	appAddr, txApp, _, err := testdapp.DeployTestDApp(sm.GoerliAuth, sm.GoerliClient, sm.ConnectorEthAddr, sm.ZetaEthAddr)
	if err != nil {
		panic(err)
	}
	sm.TestDAppAddr = appAddr
	sm.Logger.Info("TestDApp contract address: %s, tx hash: %s", appAddr.Hex(), txApp.Hash().Hex())

	// check contract deployment receipt
	if receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, txDonation, sm.Logger, sm.ReceiptTimeout); receipt.Status != 1 {
		panic("GOERLI donation tx failed")
	}
	if receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, txZetaEth, sm.Logger, sm.ReceiptTimeout); receipt.Status != 1 {
		panic("ZetaEth deployment failed")
	}
	if receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, txConnector, sm.Logger, sm.ReceiptTimeout); receipt.Status != 1 {
		panic("ZetaConnectorEth deployment failed")
	}
	if receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, txCustody, sm.Logger, sm.ReceiptTimeout); receipt.Status != 1 {
		panic("ERC20Custody deployment failed")
	}
	if receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, txUSDT, sm.Logger, sm.ReceiptTimeout); receipt.Status != 1 {
		panic("USDT deployment failed")
	}
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, txApp, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("TestDApp deployment failed")
	}

	// initialize custody contract
	sm.Logger.Info("Whitelist USDT")
	txWhitelist, err := ERC20Custody.Whitelist(sm.GoerliAuth, usdtAddr)
	if err != nil {
		panic(err)
	}
	if receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, txWhitelist, sm.Logger, sm.ReceiptTimeout); receipt.Status != 1 {
		panic("USDT whitelist failed")
	}

	sm.Logger.Info("Set TSS address")
	txCustody, err = ERC20Custody.UpdateTSSAddress(sm.GoerliAuth, sm.TSSAddress)
	if err != nil {
		panic(err)
	}
	if receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, txCustody, sm.Logger, sm.ReceiptTimeout); receipt.Status != 1 {
		panic("USDT update TSS address failed")
	}
	sm.Logger.Info("TSS set receipt tx hash: %s", txCustody.Hash().Hex())

	// save config containing contract addresses
	// TODO: put this logic outside of this function in a general config
	// We use this config to be consistent with the old implementation
	// https://github.com/zeta-chain/node-private/issues/41
	if err := config.WriteConfig(ContractsConfigFile, conf); err != nil {
		panic(err)
	}
}
