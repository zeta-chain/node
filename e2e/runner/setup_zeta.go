package runner

import (
	"math/big"
	"time"

	"github.com/zeta-chain/zetacore/e2e/contracts/contextapp"
	"github.com/zeta-chain/zetacore/e2e/contracts/zevmswap"
	utils2 "github.com/zeta-chain/zetacore/e2e/utils"

	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/connectorzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/wzeta.sol"

	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/zeta-chain/zetacore/common"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// SetTSSAddresses set TSS addresses from information queried from ZetaChain
func (runner *E2ERunner) SetTSSAddresses() error {
	runner.Logger.Print("⚙️ setting up TSS address")

	btcChainID, err := common.GetBTCChainIDFromChainParams(runner.BitcoinParams)
	if err != nil {
		return err
	}

	res := &observertypes.QueryGetTssAddressResponse{}
	for i := 0; ; i++ {
		res, err = runner.ObserverClient.GetTssAddress(runner.Ctx, &observertypes.QueryGetTssAddressRequest{
			BitcoinChainId: btcChainID,
		})
		if err != nil {
			if i%10 == 0 {
				runner.Logger.Info("ObserverClient.TSS error %s", err.Error())
				runner.Logger.Info("TSS not ready yet, waiting for TSS to be appear in zetacore network...")
			}
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	tssAddress := ethcommon.HexToAddress(res.Eth)

	btcTSSAddress, err := btcutil.DecodeAddress(res.Btc, runner.BitcoinParams)
	if err != nil {
		panic(err)
	}

	runner.TSSAddress = tssAddress
	runner.BTCTSSAddress = btcTSSAddress

	return nil
}

// SetZEVMContracts set contracts for the ZEVM
func (runner *E2ERunner) SetZEVMContracts() {
	runner.Logger.Print("⚙️ deploying system contracts and ZRC20s on ZEVM")
	startTime := time.Now()
	defer func() {
		runner.Logger.Info("System contract deployments took %s\n", time.Since(startTime))
	}()

	// deploy system contracts and ZRC20 contracts on ZetaChain
	uniswapV2FactoryAddr, uniswapV2RouterAddr, zevmConnectorAddr, wzetaAddr, usdtZRC20Addr, err := runner.ZetaTxServer.DeploySystemContractsAndZRC20(
		utils2.FungibleAdminName,
		runner.USDTERC20Addr.Hex(),
	)
	if err != nil {
		panic(err)
	}

	// Set USDTZRC20Addr
	runner.USDTZRC20Addr = ethcommon.HexToAddress(usdtZRC20Addr)
	runner.USDTZRC20, err = zrc20.NewZRC20(runner.USDTZRC20Addr, runner.ZevmClient)
	if err != nil {
		panic(err)
	}

	// UniswapV2FactoryAddr
	runner.UniswapV2FactoryAddr = ethcommon.HexToAddress(uniswapV2FactoryAddr)
	runner.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(runner.UniswapV2FactoryAddr, runner.ZevmClient)
	if err != nil {
		panic(err)
	}

	// UniswapV2RouterAddr
	runner.UniswapV2RouterAddr = ethcommon.HexToAddress(uniswapV2RouterAddr)
	runner.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(runner.UniswapV2RouterAddr, runner.ZevmClient)
	if err != nil {
		panic(err)
	}

	// ZevmConnectorAddr
	runner.ConnectorZEVMAddr = ethcommon.HexToAddress(zevmConnectorAddr)
	runner.ConnectorZEVM, err = connectorzevm.NewZetaConnectorZEVM(runner.ConnectorZEVMAddr, runner.ZevmClient)
	if err != nil {
		panic(err)
	}

	// WZetaAddr
	runner.WZetaAddr = ethcommon.HexToAddress(wzetaAddr)
	runner.WZeta, err = wzeta.NewWETH9(runner.WZetaAddr, runner.ZevmClient)
	if err != nil {
		panic(err)
	}

	// query system contract address from the chain
	systemContractRes, err := runner.FungibleClient.SystemContract(
		runner.Ctx,
		&fungibletypes.QueryGetSystemContractRequest{},
	)
	if err != nil {
		panic(err)
	}
	systemContractAddr := ethcommon.HexToAddress(systemContractRes.SystemContract.SystemContract)

	SystemContract, err := systemcontract.NewSystemContract(
		systemContractAddr,
		runner.ZevmClient,
	)
	if err != nil {
		panic(err)
	}

	runner.SystemContract = SystemContract
	runner.SystemContractAddr = systemContractAddr

	// set ZRC20 contracts
	runner.SetupETHZRC20()
	runner.SetupBTCZRC20()

	// deploy ZEVMSwapApp and ContextApp
	zevmSwapAppAddr, txZEVMSwapApp, zevmSwapApp, err := zevmswap.DeployZEVMSwapApp(
		runner.ZevmAuth,
		runner.ZevmClient,
		runner.UniswapV2RouterAddr,
		runner.SystemContractAddr,
	)
	if err != nil {
		panic(err)
	}

	contextAppAddr, txContextApp, contextApp, err := contextapp.DeployContextApp(runner.ZevmAuth, runner.ZevmClient)
	if err != nil {
		panic(err)
	}

	receipt := utils2.MustWaitForTxReceipt(runner.Ctx, runner.ZevmClient, txZEVMSwapApp, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("ZEVMSwapApp deployment failed")
	}
	runner.ZEVMSwapAppAddr = zevmSwapAppAddr
	runner.ZEVMSwapApp = zevmSwapApp

	receipt = utils2.MustWaitForTxReceipt(runner.Ctx, runner.ZevmClient, txContextApp, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("ContextApp deployment failed")
	}
	runner.ContextAppAddr = contextAppAddr
	runner.ContextApp = contextApp
}

func (runner *E2ERunner) SetupETHZRC20() {
	ethZRC20Addr, err := runner.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.GoerliLocalnetChain().ChainId))
	if err != nil {
		panic(err)
	}
	if (ethZRC20Addr == ethcommon.Address{}) {
		panic("eth zrc20 not found")
	}
	runner.ETHZRC20Addr = ethZRC20Addr
	ethZRC20, err := zrc20.NewZRC20(ethZRC20Addr, runner.ZevmClient)
	if err != nil {
		panic(err)
	}
	runner.ETHZRC20 = ethZRC20
}

func (runner *E2ERunner) SetupBTCZRC20() {
	BTCZRC20Addr, err := runner.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain().ChainId))
	if err != nil {
		panic(err)
	}
	runner.BTCZRC20Addr = BTCZRC20Addr
	runner.Logger.Info("BTCZRC20Addr: %s", BTCZRC20Addr.Hex())
	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, runner.ZevmClient)
	if err != nil {
		panic(err)
	}
	runner.BTCZRC20 = BTCZRC20
}
