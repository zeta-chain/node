package runner

import (
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/evm/erc20custody.sol"
	zetaeth "github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/evm/zeta.eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/evm/zetaconnector.eth.sol"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/contracts/erc20"
	"github.com/zeta-chain/node/e2e/contracts/testdapp"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const (
	ContractsConfigFile = "contracts.toml"
)

// SetEVMContractsFromConfig set EVM contracts for e2e test from the config
func (r *E2ERunner) SetEVMContractsFromConfig() {
	conf, err := config.ReadConfig(ContractsConfigFile)
	require.NoError(r, err)

	// Set ZetaEthAddr
	r.ZetaEthAddr = ethcommon.HexToAddress(conf.Contracts.EVM.ZetaEthAddr.String())
	r.ZetaEth, err = zetaeth.NewZetaEth(r.ZetaEthAddr, r.EVMClient)
	require.NoError(r, err)

	// Set ConnectorEthAddr
	r.ConnectorEthAddr = ethcommon.HexToAddress(conf.Contracts.EVM.ConnectorEthAddr.String())
	r.ConnectorEth, err = zetaconnectoreth.NewZetaConnectorEth(r.ConnectorEthAddr, r.EVMClient)
	require.NoError(r, err)
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
	txDonation, err := r.SendEther(r.TSSAddress, big.NewInt(101000000000000000), []byte(constant.DonationMessage))
	require.NoError(r, err)

	r.Logger.Info("Deploying ZetaEth contract")
	zetaEthAddr, txZetaEth, ZetaEth, err := zetaeth.DeployZetaEth(
		r.EVMAuth,
		r.EVMClient,
		r.EVMAddress(),
		big.NewInt(21_000_000_000),
	)
	require.NoError(r, err)

	r.ZetaEth = ZetaEth
	r.ZetaEthAddr = zetaEthAddr
	conf.Contracts.EVM.ZetaEthAddr = config.DoubleQuotedString(zetaEthAddr.String())
	r.Logger.Info("ZetaEth contract address: %s, tx hash: %s", zetaEthAddr.Hex(), zetaEthAddr.Hash().Hex())

	r.Logger.Info("Deploying ZetaConnectorEth contract")
	connectorEthAddr, txConnector, ConnectorEth, err := zetaconnectoreth.DeployZetaConnectorEth(
		r.EVMAuth,
		r.EVMClient,
		zetaEthAddr,
		r.TSSAddress,
		r.EVMAddress(),
		r.EVMAddress(),
	)
	require.NoError(r, err)

	r.ConnectorEth = ConnectorEth
	r.ConnectorEthAddr = connectorEthAddr
	conf.Contracts.EVM.ConnectorEthAddr = config.DoubleQuotedString(connectorEthAddr.String())

	r.Logger.Info(
		"ZetaConnectorEth contract address: %s, tx hash: %s",
		connectorEthAddr.Hex(),
		txConnector.Hash().Hex(),
	)

	r.Logger.Info("Deploying ERC20Custody contract")
	erc20CustodyAddr, txCustody, ERC20Custody, err := erc20custody.DeployERC20Custody(
		r.EVMAuth,
		r.EVMClient,
		r.EVMAddress(),
		r.EVMAddress(),
		big.NewInt(0),
		big.NewInt(1e18),
		ethcommon.HexToAddress("0x"),
	)
	require.NoError(r, err)

	r.ERC20CustodyAddr = erc20CustodyAddr
	r.ERC20Custody = ERC20Custody
	r.Logger.Info("ERC20Custody contract address: %s, tx hash: %s", erc20CustodyAddr.Hex(), txCustody.Hash().Hex())

	r.Logger.Info("Deploying ERC20 contract")
	erc20Addr, txERC20, erc20, err := erc20.DeployERC20(r.EVMAuth, r.EVMClient, "TESTERC20", "TESTERC20", 6)
	require.NoError(r, err)

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
	require.NoError(r, err)

	r.EvmTestDAppAddr = appAddr
	r.Logger.Info("TestDApp contract address: %s, tx hash: %s", appAddr.Hex(), txApp.Hash().Hex())

	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage)
	}

	// check contract deployment receipt
	ensureTxReceipt(txDonation, "EVM donation tx failed")
	ensureTxReceipt(txZetaEth, "ZetaEth deployment failed")
	ensureTxReceipt(txConnector, "ZetaConnectorEth deployment failed")
	ensureTxReceipt(txCustody, "ERC20Custody deployment failed")
	ensureTxReceipt(txERC20, "ERC20 deployment failed")
	ensureTxReceipt(txApp, "TestDApp deployment failed")

	// initialize custody contract
	r.Logger.Info("Whitelist ERC20")
	if whitelistERC20 {
		txWhitelist, err := ERC20Custody.Whitelist(r.EVMAuth, erc20Addr)
		require.NoError(r, err)
		ensureTxReceipt(txWhitelist, "ERC20 whitelist failed")
	}

	r.Logger.Info("Set TSS address")
	txCustody, err = ERC20Custody.UpdateTSSAddress(r.EVMAuth, r.TSSAddress)
	require.NoError(r, err)

	ensureTxReceipt(txCustody, "ERC20 update TSS address failed")

	r.Logger.Info("TSS set receipt tx hash: %s", txCustody.Hash().Hex())

	// save config containing contract addresses
	// TODO: put this logic outside of this function in a general config
	// We use this config to be consistent with the old implementation
	// https://github.com/zeta-chain/node-private/issues/41
	require.NoError(r, config.WriteConfig(ContractsConfigFile, conf))

	// chain params will need to be updated if they do not match the default params
	// this be required if the deployer account changes
	currentChainParamsRes, err := r.ObserverClient.GetChainParamsForChain(
		r.Ctx,
		&observertypes.QueryGetChainParamsForChainRequest{
			ChainId: chains.GoerliLocalnet.ChainId,
		},
	)
	require.NoError(r, err, "failed to get chain params for chain %d", chains.GoerliLocalnet.ChainId)

	chainParams := currentChainParamsRes.ChainParams
	chainParams.Erc20CustodyContractAddress = r.ERC20CustodyAddr.Hex()
	chainParams.ConnectorContractAddress = r.ConnectorEthAddr.Hex()
	chainParams.ZetaTokenContractAddress = r.ZetaEthAddr.Hex()

	_, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, observertypes.NewMsgUpdateChainParams(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		chainParams,
	))

	require.NoError(r, err, "failed to update chain params")
	r.Logger.Print("🔄 updated chain params")
}
