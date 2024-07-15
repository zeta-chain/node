package runner

import (
	"math/big"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/wzeta.sol"
	connectorzevm "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zetaconnectorzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-periphery/contracts/uniswapv2router02.sol"

	"github.com/zeta-chain/zetacore/e2e/contracts/contextapp"
	"github.com/zeta-chain/zetacore/e2e/contracts/testdapp"
	"github.com/zeta-chain/zetacore/e2e/contracts/zevmswap"
	"github.com/zeta-chain/zetacore/e2e/txserver"
	e2eutils "github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/chains"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// EmissionsPoolFunding represents the amount of ZETA to fund the emissions pool with
// This is the same value as used originally on mainnet (20M ZETA)
var EmissionsPoolFunding = big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(2e7))

// SetTSSAddresses set TSS addresses from information queried from ZetaChain
func (r *E2ERunner) SetTSSAddresses() error {
	r.Logger.Print("⚙️ setting up TSS address")

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

	return nil
}

// SetZEVMContracts set contracts for the ZEVM
func (r *E2ERunner) SetZEVMContracts() {
	r.Logger.Print("⚙️ deploying system contracts and ZRC20s on ZEVM")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("System contract deployments took %s\n", time.Since(startTime))
	}()

	// deploy system contracts and ZRC20 contracts on ZetaChain
	uniswapV2FactoryAddr, uniswapV2RouterAddr, zevmConnectorAddr, wzetaAddr, erc20zrc20Addr, err := r.ZetaTxServer.DeploySystemContractsAndZRC20(
		e2eutils.OperationalPolicyName,
		r.ERC20Addr.Hex(),
	)
	require.NoError(r, err)

	// Set ERC20ZRC20Addr
	r.ERC20ZRC20Addr = ethcommon.HexToAddress(erc20zrc20Addr)
	r.ERC20ZRC20, err = zrc20.NewZRC20(r.ERC20ZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)

	// UniswapV2FactoryAddr
	r.UniswapV2FactoryAddr = ethcommon.HexToAddress(uniswapV2FactoryAddr)
	r.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(r.UniswapV2FactoryAddr, r.ZEVMClient)
	require.NoError(r, err)

	// UniswapV2RouterAddr
	r.UniswapV2RouterAddr = ethcommon.HexToAddress(uniswapV2RouterAddr)
	r.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(r.UniswapV2RouterAddr, r.ZEVMClient)
	require.NoError(r, err)

	// ZevmConnectorAddr
	r.ConnectorZEVMAddr = ethcommon.HexToAddress(zevmConnectorAddr)
	r.ConnectorZEVM, err = connectorzevm.NewZetaConnectorZEVM(r.ConnectorZEVMAddr, r.ZEVMClient)
	require.NoError(r, err)

	// WZetaAddr
	r.WZetaAddr = ethcommon.HexToAddress(wzetaAddr)
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

	// set ZRC20 contracts
	r.SetupETHZRC20()
	r.SetupBTCZRC20()

	// deploy TestDApp contract on zEVM
	appAddr, txApp, _, err := testdapp.DeployTestDApp(
		r.ZEVMAuth,
		r.ZEVMClient,
		r.ConnectorZEVMAddr,
		r.WZetaAddr,
	)
	require.NoError(r, err)

	r.ZevmTestDAppAddr = appAddr
	r.Logger.Info("TestDApp Zevm contract address: %s, tx hash: %s", appAddr.Hex(), txApp.Hash().Hex())

	// deploy ZEVMSwapApp and ContextApp
	zevmSwapAppAddr, txZEVMSwapApp, zevmSwapApp, err := zevmswap.DeployZEVMSwapApp(
		r.ZEVMAuth,
		r.ZEVMClient,
		r.UniswapV2RouterAddr,
		r.SystemContractAddr,
	)
	require.NoError(r, err)

	contextAppAddr, txContextApp, contextApp, err := contextapp.DeployContextApp(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	receipt := e2eutils.MustWaitForTxReceipt(
		r.Ctx,
		r.ZEVMClient,
		txZEVMSwapApp,
		r.Logger,
		r.ReceiptTimeout,
	)
	r.requireTxSuccessful(receipt, "ZEVMSwapApp deployment failed")

	r.ZEVMSwapAppAddr = zevmSwapAppAddr
	r.ZEVMSwapApp = zevmSwapApp

	receipt = e2eutils.MustWaitForTxReceipt(
		r.Ctx,
		r.ZEVMClient,
		txContextApp,
		r.Logger,
		r.ReceiptTimeout,
	)
	r.requireTxSuccessful(receipt, "ContextApp deployment failed")

	r.ContextAppAddr = contextAppAddr
	r.ContextApp = contextApp
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

// EnableHeaderVerification enables the header verification for the given chain IDs
func (r *E2ERunner) EnableHeaderVerification(chainIDList []int64) error {
	r.Logger.Print("⚙️ enabling verification flags for block headers")

	return r.ZetaTxServer.EnableHeaderVerification(e2eutils.OperationalPolicyName, chainIDList)
}

// FundEmissionsPool funds the emissions pool on ZetaChain with the same value as used originally on mainnet (20M ZETA)
func (r *E2ERunner) FundEmissionsPool() error {
	r.Logger.Print("⚙️ funding the emissions pool on ZetaChain with 20M ZETA (%s)", txserver.EmissionsPoolAddress)

	return r.ZetaTxServer.FundEmissionsPool(e2eutils.OperationalPolicyName, EmissionsPoolFunding)
}
