package runner

import (
	"math/big"
	"time"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/contracts/erc20"
	"github.com/zeta-chain/zetacore/e2e/contracts/testdapp"
	"github.com/zeta-chain/zetacore/e2e/utils"

	"github.com/zeta-chain/zetacore/zetaclient"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	zetaeth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zeta.eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
)

const (
	ContractsConfigFile = "contracts.toml"
)

// SetEVMContractsFromConfig set EVM contracts for e2e test from the config
func (runner *E2ERunner) SetEVMContractsFromConfig() {
	conf, err := config.ReadConfig(ContractsConfigFile)
	if err != nil {
		panic(err)
	}

	// Set ZetaEthAddr
	runner.ZetaEthAddr = ethcommon.HexToAddress(conf.Contracts.EVM.ZetaEthAddress)
	runner.ZetaEth, err = zetaeth.NewZetaEth(runner.ZetaEthAddr, runner.GoerliClient)
	if err != nil {
		panic(err)
	}

	// Set ConnectorEthAddr
	runner.ConnectorEthAddr = ethcommon.HexToAddress(conf.Contracts.EVM.ConnectorEthAddr)
	runner.ConnectorEth, err = zetaconnectoreth.NewZetaConnectorEth(runner.ConnectorEthAddr, runner.GoerliClient)
	if err != nil {
		panic(err)
	}
}

// SetupEVM setup contracts on EVM for e2e test
func (runner *E2ERunner) SetupEVM(contractsDeployed bool) {
	runner.Logger.Print("⚙️ setting up Goerli network")
	startTime := time.Now()
	defer func() {
		runner.Logger.Info("EVM setup took %s\n", time.Since(startTime))
	}()

	// TODO: put this logic outside of this function
	// We use this config to be consistent with the old implementation
	// https://github.com/zeta-chain/node-private/issues/41
	if contractsDeployed {
		runner.SetEVMContractsFromConfig()
		return
	}
	conf := config.DefaultConfig()

	runner.Logger.InfoLoud("Deploy ZetaETH ConnectorETH ERC20Custody USDT\n")

	// donate to the TSS address to avoid account errors because deploying gas token ZRC20 will automatically mint
	// gas token on ZetaChain to initialize the pool
	txDonation, err := runner.SendEther(runner.TSSAddress, big.NewInt(101000000000000000), []byte(zetaclient.DonationMessage))
	if err != nil {
		panic(err)
	}

	runner.Logger.Info("Deploying ZetaEth contract")
	zetaEthAddr, txZetaEth, ZetaEth, err := zetaeth.DeployZetaEth(
		runner.GoerliAuth,
		runner.GoerliClient,
		runner.DeployerAddress,
		big.NewInt(21_000_000_000),
	)
	if err != nil {
		panic(err)
	}
	runner.ZetaEth = ZetaEth
	runner.ZetaEthAddr = zetaEthAddr
	conf.Contracts.EVM.ZetaEthAddress = zetaEthAddr.String()
	runner.Logger.Info("ZetaEth contract address: %s, tx hash: %s", zetaEthAddr.Hex(), zetaEthAddr.Hash().Hex())

	runner.Logger.Info("Deploying ZetaConnectorEth contract")
	connectorEthAddr, txConnector, ConnectorEth, err := zetaconnectoreth.DeployZetaConnectorEth(
		runner.GoerliAuth,
		runner.GoerliClient,
		zetaEthAddr,
		runner.TSSAddress,
		runner.DeployerAddress,
		runner.DeployerAddress,
	)
	if err != nil {
		panic(err)
	}
	runner.ConnectorEth = ConnectorEth
	runner.ConnectorEthAddr = connectorEthAddr
	conf.Contracts.EVM.ConnectorEthAddr = connectorEthAddr.String()

	runner.Logger.Info("ZetaConnectorEth contract address: %s, tx hash: %s", connectorEthAddr.Hex(), txConnector.Hash().Hex())

	runner.Logger.Info("Deploying ERC20Custody contract")
	erc20CustodyAddr, txCustody, ERC20Custody, err := erc20custody.DeployERC20Custody(
		runner.GoerliAuth,
		runner.GoerliClient,
		runner.DeployerAddress,
		runner.DeployerAddress,
		big.NewInt(0),
		big.NewInt(1e18),
		ethcommon.HexToAddress("0x"),
	)
	if err != nil {
		panic(err)
	}
	runner.ERC20CustodyAddr = erc20CustodyAddr
	runner.ERC20Custody = ERC20Custody
	runner.Logger.Info("ERC20Custody contract address: %s, tx hash: %s", erc20CustodyAddr.Hex(), txCustody.Hash().Hex())

	runner.Logger.Info("Deploying USDT contract")
	usdtAddr, txUSDT, usdt, err := erc20.DeployUSDT(runner.GoerliAuth, runner.GoerliClient, "USDT", "USDT", 6)
	if err != nil {
		panic(err)
	}
	runner.USDTERC20 = usdt
	runner.USDTERC20Addr = usdtAddr
	runner.Logger.Info("USDT contract address: %s, tx hash: %s", usdtAddr.Hex(), txUSDT.Hash().Hex())

	// deploy TestDApp contract
	appAddr, txApp, _, err := testdapp.DeployTestDApp(runner.GoerliAuth, runner.GoerliClient, runner.ConnectorEthAddr, runner.ZetaEthAddr)
	if err != nil {
		panic(err)
	}
	runner.TestDAppAddr = appAddr
	runner.Logger.Info("TestDApp contract address: %s, tx hash: %s", appAddr.Hex(), txApp.Hash().Hex())

	// check contract deployment receipt
	if receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, txDonation, runner.Logger, runner.ReceiptTimeout); receipt.Status != 1 {
		panic("GOERLI donation tx failed")
	}
	if receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, txZetaEth, runner.Logger, runner.ReceiptTimeout); receipt.Status != 1 {
		panic("ZetaEth deployment failed")
	}
	if receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, txConnector, runner.Logger, runner.ReceiptTimeout); receipt.Status != 1 {
		panic("ZetaConnectorEth deployment failed")
	}
	if receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, txCustody, runner.Logger, runner.ReceiptTimeout); receipt.Status != 1 {
		panic("ERC20Custody deployment failed")
	}
	if receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, txUSDT, runner.Logger, runner.ReceiptTimeout); receipt.Status != 1 {
		panic("USDT deployment failed")
	}
	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, txApp, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("TestDApp deployment failed")
	}

	// initialize custody contract
	runner.Logger.Info("Whitelist USDT")
	txWhitelist, err := ERC20Custody.Whitelist(runner.GoerliAuth, usdtAddr)
	if err != nil {
		panic(err)
	}
	if receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, txWhitelist, runner.Logger, runner.ReceiptTimeout); receipt.Status != 1 {
		panic("USDT whitelist failed")
	}

	runner.Logger.Info("Set TSS address")
	txCustody, err = ERC20Custody.UpdateTSSAddress(runner.GoerliAuth, runner.TSSAddress)
	if err != nil {
		panic(err)
	}
	if receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.GoerliClient, txCustody, runner.Logger, runner.ReceiptTimeout); receipt.Status != 1 {
		panic("USDT update TSS address failed")
	}
	runner.Logger.Info("TSS set receipt tx hash: %s", txCustody.Hash().Hex())

	// save config containing contract addresses
	// TODO: put this logic outside of this function in a general config
	// We use this config to be consistent with the old implementation
	// https://github.com/zeta-chain/node-private/issues/41
	if err := config.WriteConfig(ContractsConfigFile, conf); err != nil {
		panic(err)
	}
}
