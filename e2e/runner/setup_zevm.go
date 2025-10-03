package runner

import (
	"math/big"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/wzeta.sol"
	connectorzevm "github.com/zeta-chain/protocol-contracts/pkg/zetaconnectorzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/zrc20.sol"

	"github.com/zeta-chain/node/e2e/contracts/erc1967proxy"
	"github.com/zeta-chain/node/e2e/contracts/testdappv2"
	"github.com/zeta-chain/node/e2e/txserver"
	e2eutils "github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/contracts/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/zeta-chain/node/pkg/contracts/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// EmissionsPoolFunding represents the amount of ZETA to fund the emissions pool with
// This is the same value as used originally on mainnet (20M ZETA)
var EmissionsPoolFunding = big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(2e7))

// SetTSSAddresses set TSS addresses from information queried from ZetaChain
func (r *E2ERunner) SetTSSAddresses() error {
	btcChainID, err := chains.GetBTCChainIDFromChainParams(r.BitcoinParams)
	if err != nil {
		return err
	}

	res := &observertypes.QueryGetTssAddressResponse{}
	for i := 0; ; i++ {
		res, err = r.ObserverClient.GetTssAddress(r.Ctx, &observertypes.QueryGetTssAddressRequest{
			BitcoinChainId: btcChainID,
		})
		if err != nil {
			if i%10 == 0 {
				r.Logger.Info("ObserverClient.TSS error %s", err.Error())
				r.Logger.Info("TSS not ready yet, waiting for TSS to be appear in zetacore network...")
			}
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	tssAddress := ethcommon.HexToAddress(res.Eth)

	btcTSSAddress, err := btcutil.DecodeAddress(res.Btc, r.BitcoinParams)
	require.NoError(r, err)

	r.TSSAddress = tssAddress
	r.BTCTSSAddress = btcTSSAddress
	r.SuiTSSAddress = res.Sui

	return nil
}

// SetupZEVMZRC20s setup ZRC20 for the ZEVM
func (r *E2ERunner) SetupZEVMZRC20s(zrc20Deployment txserver.ZRC20Deployment) {
	r.Logger.Print("⚙️ deploying ZRC20s on ZEVM")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("System contract deployments took %s\n", time.Since(startTime))
	}()

	// deploy system contracts and ZRC20 contracts on ZetaChain
	deployedZRC20Addresses, err := r.ZetaTxServer.DeployZRC20s(
		zrc20Deployment,
		r.skipChainOperations,
	)
	require.NoError(r, err)

	// Set ERC20ZRC20Addr
	r.ERC20ZRC20Addr = deployedZRC20Addresses.ERC20ZRC20Addr
	r.ERC20ZRC20, err = zrc20.NewZRC20(r.ERC20ZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)

	// Set SPLZRC20Addr if set
	if deployedZRC20Addresses.SPLZRC20Addr != (ethcommon.Address{}) {
		r.SPLZRC20Addr = deployedZRC20Addresses.SPLZRC20Addr
		r.SPLZRC20, err = zrc20.NewZRC20(r.SPLZRC20Addr, r.ZEVMClient)
		require.NoError(r, err)
	}

	// set ZRC20 contracts
	r.SetupETHZRC20()
	r.SetupBTCZRC20()
	r.SetupSOLZRC20()
	r.SetupTONZRC20()
}

// SetupETHZRC20 sets up the ETH ZRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupETHZRC20() {
	ethZRC20Addr, err := r.SystemContract.GasCoinZRC20ByChainId(
		&bind.CallOpts{},
		big.NewInt(chains.GoerliLocalnet.ChainId),
	)
	require.NoError(r, err)
	require.NotEqual(r, ethcommon.Address{}, ethZRC20Addr, "eth zrc20 not found")

	r.ETHZRC20Addr = ethZRC20Addr
	ethZRC20, err := zrc20.NewZRC20(ethZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)

	r.ETHZRC20 = ethZRC20
}

// SetupBTCZRC20 sets up the BTC ZRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupBTCZRC20() {
	BTCZRC20Addr, err := r.SystemContract.GasCoinZRC20ByChainId(
		&bind.CallOpts{},
		big.NewInt(chains.BitcoinRegtest.ChainId),
	)
	require.NoError(r, err)
	r.BTCZRC20Addr = BTCZRC20Addr
	r.Logger.Info("BTCZRC20Addr: %s", BTCZRC20Addr.Hex())
	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)
	r.BTCZRC20 = BTCZRC20
}

// SetupSOLZRC20 sets up the SOL ZRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupSOLZRC20() {
	// set SOLZRC20 address by chain ID
	SOLZRC20Addr, err := r.SystemContract.GasCoinZRC20ByChainId(
		&bind.CallOpts{},
		big.NewInt(chains.SolanaLocalnet.ChainId),
	)
	require.NoError(r, err)

	// set SOLZRC20 address
	r.SOLZRC20Addr = SOLZRC20Addr
	r.Logger.Info("SOLZRC20Addr: %s", SOLZRC20Addr.Hex())

	// set SOLZRC20 contract
	SOLZRC20, err := zrc20.NewZRC20(SOLZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)
	r.SOLZRC20 = SOLZRC20
}

// SetupTONZRC20 sets up the TON ZRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupTONZRC20() {
	chainID := chains.TONLocalnet.ChainId

	// noop
	if r.skipChainOperations(chainID) {
		return
	}

	TONZRC20Addr, err := r.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(chainID))
	require.NoError(r, err)

	r.TONZRC20Addr = TONZRC20Addr
	r.Logger.Info("TON ZRC20 address: %s", TONZRC20Addr.Hex())

	TONZRC20, err := zrc20.NewZRC20(TONZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)

	r.TONZRC20 = TONZRC20
}

// SetupSUIZRC20 sets up the SUI ZRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupSUIZRC20() {
	chainID := chains.SuiLocalnet.ChainId

	// noop
	if r.skipChainOperations(chainID) {
		return
	}

	SUIZRC20Addr, err := r.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(chainID))
	require.NoError(r, err)

	r.SUIZRC20Addr = SUIZRC20Addr
	r.Logger.Info("SUI ZRC20 address: %s", SUIZRC20Addr.Hex())

	SUIZRC20, err := zrc20.NewZRC20(SUIZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)

	r.SUIZRC20 = SUIZRC20
}

// EnableHeaderVerification enables the header verification for the given chain IDs
func (r *E2ERunner) EnableHeaderVerification(chainIDList []int64) error {
	r.Logger.Print("⚙️ enabling verification flags for block headers")

	return r.ZetaTxServer.EnableHeaderVerification(e2eutils.AdminPolicyName, chainIDList)
}

// SetupZEVMProtocolContracts setup protocol contracts for the ZEVM
func (r *E2ERunner) SetupZEVMProtocolContracts() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := e2eutils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage+" tx hash: "+tx.Hash().Hex())
	}

	r.Logger.Print("⚙️ setting up ZEVM protocol contracts")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("ZEVM protocol contracts took %s\n", time.Since(startTime))
	}()

	// deploy system contracts and ZRC20 contracts on ZetaChain
	addresses, err := r.ZetaTxServer.DeploySystemContracts(
		e2eutils.OperationalPolicyName,
		e2eutils.AdminPolicyName,
	)
	require.NoError(r, err)

	// UniswapV2FactoryAddr
	r.UniswapV2FactoryAddr = ethcommon.HexToAddress(addresses.UniswapV2FactoryAddr)
	r.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(r.UniswapV2FactoryAddr, r.ZEVMClient)
	require.NoError(r, err)

	// UniswapV2RouterAddr
	r.UniswapV2RouterAddr = ethcommon.HexToAddress(addresses.UniswapV2RouterAddr)
	r.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(r.UniswapV2RouterAddr, r.ZEVMClient)
	require.NoError(r, err)

	// ZevmConnectorAddr
	r.ConnectorZEVMAddr = ethcommon.HexToAddress(addresses.ZEVMConnectorAddr)
	r.ConnectorZEVM, err = connectorzevm.NewZetaConnectorZEVM(r.ConnectorZEVMAddr, r.ZEVMClient)
	require.NoError(r, err)

	// WZetaAddr
	r.WZetaAddr = ethcommon.HexToAddress(addresses.WZETAAddr)
	r.WZeta, err = wzeta.NewWETH9(r.WZetaAddr, r.ZEVMClient)
	require.NoError(r, err)

	// query system contract address from the chain
	systemContractRes, err := r.FungibleClient.SystemContract(
		r.Ctx,
		&fungibletypes.QueryGetSystemContractRequest{},
	)
	require.NoError(r, err)

	systemContractAddr := ethcommon.HexToAddress(systemContractRes.SystemContract.SystemContract)
	systemContract, err := systemcontract.NewSystemContract(
		systemContractAddr,
		r.ZEVMClient,
	)
	require.NoError(r, err)

	r.SystemContract = systemContract
	r.SystemContractAddr = systemContractAddr

	r.Logger.Info("Deploying Gateway ZEVM")
	gatewayZEVMAddr, txGateway, _, err := gatewayzevm.DeployGatewayZEVM(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	ensureTxReceipt(txGateway, "Gateway deployment failed")

	gatewayZEVMABI, err := gatewayzevm.GatewayZEVMMetaData.GetAbi()
	require.NoError(r, err)

	// Encode the initializer data
	initializerData, err := gatewayZEVMABI.Pack("initialize", r.WZetaAddr, r.Account.EVMAddress())
	require.NoError(r, err)

	// Deploy the proxy contract
	r.Logger.Info(
		"Deploying proxy with %s and %s, address: %s",
		r.WZetaAddr.Hex(),
		r.Account.EVMAddress().Hex(),
		gatewayZEVMAddr.Hex(),
	)
	proxyAddress, txProxy, _, err := erc1967proxy.DeployERC1967Proxy(
		r.ZEVMAuth,
		r.ZEVMClient,
		gatewayZEVMAddr,
		initializerData,
	)
	require.NoError(r, err)

	r.GatewayZEVMAddr = proxyAddress
	r.GatewayZEVM, err = gatewayzevm.NewGatewayZEVM(proxyAddress, r.ZEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Gateway ZEVM contract address: %s, tx hash: %s", gatewayZEVMAddr.Hex(), txGateway.Hash().Hex())

	// Set the gateway address in the protocol
	err = r.ZetaTxServer.UpdateGatewayAddress(e2eutils.AdminPolicyName, r.GatewayZEVMAddr.Hex())
	require.NoError(r, err)
	ensureTxReceipt(txProxy, "Gateway proxy deployment failed")

	r.SetupZEVMTestDappV2(ensureTxReceipt)
}

func (r *E2ERunner) SetupZEVMTestDappV2(ensureTxReceipt func(tx *ethtypes.Transaction, failMessage string)) {
	// deploy test dapp v2
	testDAppV2Addr, txTestDAppV2, _, err := testdappv2.DeployTestDAppV2(
		r.ZEVMAuth,
		r.ZEVMClient,
		true,
		r.GatewayEVMAddr,
	)
	require.NoError(r, err)
	ensureTxReceipt(txTestDAppV2, "TestDAppV2 deployment failed")

	r.TestDAppV2ZEVMAddr = testDAppV2Addr
	r.TestDAppV2ZEVM, err = testdappv2.NewTestDAppV2(testDAppV2Addr, r.ZEVMClient)
	require.NoError(r, err)

	// check isZetaChain is true
	isZetaChain, err := r.TestDAppV2ZEVM.IsZetaChain(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, isZetaChain)
}

// UpdateProtocolContractsInChainParams update the erc20 custody contract and gateway address in the chain params
// TODO: should be used for all protocol contracts including the ZETA connector
// https://github.com/zeta-chain/node/issues/3257
func (r *E2ERunner) UpdateProtocolContractsInChainParams(testLegacy bool) {
	res, err := r.ObserverClient.GetChainParams(r.Ctx, &observertypes.QueryGetChainParamsRequest{})
	require.NoError(r, err)

	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	// find old chain params
	var (
		chainParams *observertypes.ChainParams
		found       bool
	)
	for _, cp := range res.ChainParams.ChainParams {
		if cp.ChainId == evmChainID.Int64() {
			chainParams = cp
			found = true
			break
		}
	}
	require.True(r, found, "Chain params not found for chain id %d", evmChainID)

	// update with the new ERC20 custody contract address
	chainParams.Erc20CustodyContractAddress = r.ERC20CustodyAddr.Hex()

	// update with the new gateway address
	chainParams.GatewayAddress = r.GatewayEVMAddr.Hex()

	//  update with the new connector address only if not running legacy tests
	// when running legacy tests the connector address is set by the LegacySetupEVM function
	if !testLegacy {
		chainParams.ConnectorContractAddress = r.ConnectorNativeAddr.Hex()
	}

	// update the chain params
	err = r.ZetaTxServer.UpdateChainParams(chainParams)
	require.NoError(r, err)
}

// DeployTestDAppV2ZEVM deploys the test DApp V2 contract
func (r *E2ERunner) DeployTestDAppV2ZEVM() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := e2eutils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage+" tx hash: "+tx.Hash().Hex())
	}

	testDAppV2Addr, txTestDAppV2, _, err := testdappv2.DeployTestDAppV2(
		r.ZEVMAuth,
		r.ZEVMClient,
		true,
		r.GatewayEVMAddr,
	)
	require.NoError(r, err)
	ensureTxReceipt(txTestDAppV2, "TestDAppV2 deployment failed")

	r.TestDAppV2ZEVMAddr = testDAppV2Addr
	r.TestDAppV2ZEVM, err = testdappv2.NewTestDAppV2(testDAppV2Addr, r.ZEVMClient)
	require.NoError(r, err)

	isZetaChain, err := r.TestDAppV2ZEVM.IsZetaChain(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, isZetaChain)
}
