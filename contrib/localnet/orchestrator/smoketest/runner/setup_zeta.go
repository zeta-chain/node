package runner

import (
	"math/big"
	"time"

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
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/contextapp"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/zevmswap"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// SetTSSAddresses set TSS addresses from information queried from ZetaChain
func (sm *SmokeTestRunner) SetTSSAddresses() error {
	sm.Logger.Print("⚙️ setting up TSS address")

	btcChainID, err := common.GetBTCChainIDFromChainParams(sm.BitcoinParams)
	if err != nil {
		return err
	}

	res := &observertypes.QueryGetTssAddressResponse{}
	for i := 0; ; i++ {
		res, err = sm.ObserverClient.GetTssAddress(sm.Ctx, &observertypes.QueryGetTssAddressRequest{
			BitcoinChainId: btcChainID,
		})
		if err != nil {
			if i%10 == 0 {
				sm.Logger.Info("ObserverClient.TSS error %s", err.Error())
				sm.Logger.Info("TSS not ready yet, waiting for TSS to be appear in zetacore network...")
			}
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	tssAddress := ethcommon.HexToAddress(res.Eth)

	btcTSSAddress, err := btcutil.DecodeAddress(res.Btc, sm.BitcoinParams)
	if err != nil {
		panic(err)
	}

	sm.TSSAddress = tssAddress
	sm.BTCTSSAddress = btcTSSAddress

	return nil
}

// SetZEVMContracts set contracts for the ZEVM
func (sm *SmokeTestRunner) SetZEVMContracts() {
	sm.Logger.Print("⚙️ deploying system contracts and ZRC20s on ZEVM")
	startTime := time.Now()
	defer func() {
		sm.Logger.Info("System contract deployments took %s\n", time.Since(startTime))
	}()

	// deploy system contracts and ZRC20 contracts on ZetaChain
	uniswapV2FactoryAddr, uniswapV2RouterAddr, zevmConnectorAddr, wzetaAddr, usdtZRC20Addr, err := sm.ZetaTxServer.DeploySystemContractsAndZRC20(
		utils.FungibleAdminName,
		sm.USDTERC20Addr.Hex(),
	)
	if err != nil {
		panic(err)
	}

	// Set USDTZRC20Addr
	sm.USDTZRC20Addr = ethcommon.HexToAddress(usdtZRC20Addr)
	sm.USDTZRC20, err = zrc20.NewZRC20(sm.USDTZRC20Addr, sm.ZevmClient)
	if err != nil {
		panic(err)
	}

	// UniswapV2FactoryAddr
	sm.UniswapV2FactoryAddr = ethcommon.HexToAddress(uniswapV2FactoryAddr)
	sm.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(sm.UniswapV2FactoryAddr, sm.ZevmClient)
	if err != nil {
		panic(err)
	}

	// UniswapV2RouterAddr
	sm.UniswapV2RouterAddr = ethcommon.HexToAddress(uniswapV2RouterAddr)
	sm.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(sm.UniswapV2RouterAddr, sm.ZevmClient)
	if err != nil {
		panic(err)
	}

	// ZevmConnectorAddr
	sm.ConnectorZEVMAddr = ethcommon.HexToAddress(zevmConnectorAddr)
	sm.ConnectorZEVM, err = connectorzevm.NewZetaConnectorZEVM(sm.ConnectorZEVMAddr, sm.ZevmClient)
	if err != nil {
		panic(err)
	}

	// WZetaAddr
	sm.WZetaAddr = ethcommon.HexToAddress(wzetaAddr)
	sm.WZeta, err = wzeta.NewWETH9(sm.WZetaAddr, sm.ZevmClient)
	if err != nil {
		panic(err)
	}

	// query system contract address from the chain
	systemContractRes, err := sm.FungibleClient.SystemContract(
		sm.Ctx,
		&fungibletypes.QueryGetSystemContractRequest{},
	)
	if err != nil {
		panic(err)
	}
	systemContractAddr := ethcommon.HexToAddress(systemContractRes.SystemContract.SystemContract)

	SystemContract, err := systemcontract.NewSystemContract(
		systemContractAddr,
		sm.ZevmClient,
	)
	if err != nil {
		panic(err)
	}

	sm.SystemContract = SystemContract
	sm.SystemContractAddr = systemContractAddr

	// set ZRC20 contracts
	sm.SetupETHZRC20()
	sm.SetupBTCZRC20()

	// deploy ZEVMSwapApp and ContextApp
	zevmSwapAppAddr, txZEVMSwapApp, zevmSwapApp, err := zevmswap.DeployZEVMSwapApp(
		sm.ZevmAuth,
		sm.ZevmClient,
		sm.UniswapV2RouterAddr,
		sm.SystemContractAddr,
	)
	if err != nil {
		panic(err)
	}

	contextAppAddr, txContextApp, contextApp, err := contextapp.DeployContextApp(sm.ZevmAuth, sm.ZevmClient)
	if err != nil {
		panic(err)
	}

	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, txZEVMSwapApp, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("ZEVMSwapApp deployment failed")
	}
	sm.ZEVMSwapAppAddr = zevmSwapAppAddr
	sm.ZEVMSwapApp = zevmSwapApp

	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, txContextApp, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("ContextApp deployment failed")
	}
	sm.ContextAppAddr = contextAppAddr
	sm.ContextApp = contextApp
}

func (sm *SmokeTestRunner) SetupETHZRC20() {
	ethZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.GoerliLocalnetChain().ChainId))
	if err != nil {
		panic(err)
	}
	if (ethZRC20Addr == ethcommon.Address{}) {
		panic("eth zrc20 not found")
	}
	sm.ETHZRC20Addr = ethZRC20Addr
	ethZRC20, err := zrc20.NewZRC20(ethZRC20Addr, sm.ZevmClient)
	if err != nil {
		panic(err)
	}
	sm.ETHZRC20 = ethZRC20
}

func (sm *SmokeTestRunner) SetupBTCZRC20() {
	BTCZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain().ChainId))
	if err != nil {
		panic(err)
	}
	sm.BTCZRC20Addr = BTCZRC20Addr
	sm.Logger.Info("BTCZRC20Addr: %s", BTCZRC20Addr.Hex())
	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, sm.ZevmClient)
	if err != nil {
		panic(err)
	}
	sm.BTCZRC20 = BTCZRC20
}
