package runner

import (
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	zetaeth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zeta.eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/contracts/erc20"
	"github.com/zeta-chain/zetacore/e2e/contracts/testdapp"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/constant"
)

const (
	ContractsConfigFile = "contracts.toml"
)

// SetEVMContractsFromConfig set EVM contracts for e2e test from the config
func (r *E2ERunner) SetEVMContractsFromConfig() {
	conf, err := config.ReadConfig(ContractsConfigFile)
	if err != nil {
		panic(err)
	}

	// Set ZetaEthAddr
	r.ZetaEthAddr = ethcommon.HexToAddress(conf.Contracts.EVM.ZetaEthAddress)
	r.ZetaEth, err = zetaeth.NewZetaEth(r.ZetaEthAddr, r.EVMClient)
	if err != nil {
		panic(err)
	}

	// Set ConnectorEthAddr
	r.ConnectorEthAddr = ethcommon.HexToAddress(conf.Contracts.EVM.ConnectorEthAddr)
	r.ConnectorEth, err = zetaconnectoreth.NewZetaConnectorEth(r.ConnectorEthAddr, r.EVMClient)
	if err != nil {
		panic(err)
	}
}

// SetupEVM setup contracts on EVM for e2e test
func (r *E2ERunner) SetupEVM(contractsDeployed bool, whitelistERC20 bool) {
	r.Logger.Print("⚙️ setting up EVM network")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("EVM setup took %s\n", time.Since(startTime))
	}()

	// TODO: put this logic outside of this function
	// We use this config to be consistent with the old implementation
	// https://github.com/zeta-chain/node-private/issues/41
	if contractsDeployed {
		r.SetEVMContractsFromConfig()
		return
	}
	conf := config.DefaultConfig()

	r.Logger.InfoLoud("Deploy ZetaETH ConnectorETH ERC20Custody ERC20\n")

	// donate to the TSS address to avoid account errors because deploying gas token ZRC20 will automatically mint
	// gas token on ZetaChain to initialize the pool
	txDonation, err := r.SendEther(
		r.TSSAddress,
		big.NewInt(101000000000000000),
		[]byte(constant.DonationMessage),
	)
	if err != nil {
		panic(err)
	}

	r.Logger.Info("Deploying ZetaEth contract")
	zetaEthAddr, txZetaEth, ZetaEth, err := zetaeth.DeployZetaEth(
		r.EVMAuth,
		r.EVMClient,
		r.DeployerAddress,
		big.NewInt(21_000_000_000),
	)
	if err != nil {
		panic(err)
	}
	r.ZetaEth = ZetaEth
	r.ZetaEthAddr = zetaEthAddr
	conf.Contracts.EVM.ZetaEthAddress = zetaEthAddr.String()
	r.Logger.Info("ZetaEth contract address: %s, tx hash: %s", zetaEthAddr.Hex(), zetaEthAddr.Hash().Hex())

	r.Logger.Info("Deploying ZetaConnectorEth contract")
	connectorEthAddr, txConnector, ConnectorEth, err := zetaconnectoreth.DeployZetaConnectorEth(
		r.EVMAuth,
		r.EVMClient,
		zetaEthAddr,
		r.TSSAddress,
		r.DeployerAddress,
		r.DeployerAddress,
	)
	if err != nil {
		panic(err)
	}
	r.ConnectorEth = ConnectorEth
	r.ConnectorEthAddr = connectorEthAddr
	conf.Contracts.EVM.ConnectorEthAddr = connectorEthAddr.String()

	r.Logger.Info(
		"ZetaConnectorEth contract address: %s, tx hash: %s",
		connectorEthAddr.Hex(),
		txConnector.Hash().Hex(),
	)

	r.Logger.Info("Deploying ERC20Custody contract")
	erc20CustodyAddr, txCustody, ERC20Custody, err := erc20custody.DeployERC20Custody(
		r.EVMAuth,
		r.EVMClient,
		r.DeployerAddress,
		r.DeployerAddress,
		big.NewInt(0),
		big.NewInt(1e18),
		ethcommon.HexToAddress("0x"),
	)
	if err != nil {
		panic(err)
	}
	r.ERC20CustodyAddr = erc20CustodyAddr
	r.ERC20Custody = ERC20Custody
	r.Logger.Info("ERC20Custody contract address: %s, tx hash: %s", erc20CustodyAddr.Hex(), txCustody.Hash().Hex())

	r.Logger.Info("Deploying ERC20 contract")
	erc20Addr, txERC20, erc20, err := erc20.DeployERC20(r.EVMAuth, r.EVMClient, "TESTERC20", "TESTERC20", 6)
	if err != nil {
		panic(err)
	}
	r.ERC20 = erc20
	r.ERC20Addr = erc20Addr
	r.Logger.Info("ERC20 contract address: %s, tx hash: %s", erc20Addr.Hex(), txERC20.Hash().Hex())

	// deploy TestDApp contract
	appAddr, txApp, _, err := testdapp.DeployTestDApp(
		r.EVMAuth,
		r.EVMClient,
		r.ConnectorEthAddr,
		r.ZetaEthAddr,
	)
	if err != nil {
		panic(err)
	}
	r.EvmTestDAppAddr = appAddr
	r.Logger.Info("TestDApp contract address: %s, tx hash: %s", appAddr.Hex(), txApp.Hash().Hex())

	// check contract deployment receipt
	if receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, txDonation, r.Logger, r.ReceiptTimeout); receipt.Status != 1 {
		panic("EVM donation tx failed")
	}
	if receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, txZetaEth, r.Logger, r.ReceiptTimeout); receipt.Status != 1 {
		panic("ZetaEth deployment failed")
	}
	if receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, txConnector, r.Logger, r.ReceiptTimeout); receipt.Status != 1 {
		panic("ZetaConnectorEth deployment failed")
	}
	if receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, txCustody, r.Logger, r.ReceiptTimeout); receipt.Status != 1 {
		panic("ERC20Custody deployment failed")
	}
	if receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, txERC20, r.Logger, r.ReceiptTimeout); receipt.Status != 1 {
		panic("ERC20 deployment failed")
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, txApp, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("TestDApp deployment failed")
	}

	// initialize custody contract
	r.Logger.Info("Whitelist ERC20")
	if whitelistERC20 {
		txWhitelist, err := ERC20Custody.Whitelist(r.EVMAuth, erc20Addr)
		if err != nil {
			panic(err)
		}
		if receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, txWhitelist, r.Logger, r.ReceiptTimeout); receipt.Status != 1 {
			panic("ERC20 whitelist failed")
		}
	}

	r.Logger.Info("Set TSS address")
	txCustody, err = ERC20Custody.UpdateTSSAddress(r.EVMAuth, r.TSSAddress)
	if err != nil {
		panic(err)
	}
	if receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, txCustody, r.Logger, r.ReceiptTimeout); receipt.Status != 1 {
		panic("ERC20 update TSS address failed")
	}
	r.Logger.Info("TSS set receipt tx hash: %s", txCustody.Hash().Hex())

	// save config containing contract addresses
	// TODO: put this logic outside of this function in a general config
	// We use this config to be consistent with the old implementation
	// https://github.com/zeta-chain/node-private/issues/41
	if err := config.WriteConfig(ContractsConfigFile, conf); err != nil {
		panic(err)
	}
}
