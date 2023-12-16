package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/btcsuite/btcutil"
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
func (sm *SmokeTestRunner) SetTSSAddresses() {
	var err error
	res := &observertypes.QueryGetTssAddressResponse{}
	for {
		res, err = sm.ObserverClient.GetTssAddress(context.Background(), &observertypes.QueryGetTssAddressRequest{})
		if err != nil {
			fmt.Printf("cctxClient.TSS error %s\n", err.Error())
			fmt.Printf("TSS not ready yet, waiting for TSS to be appear in zetacore network...\n")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	tssAddress := ethcommon.HexToAddress(res.Eth)
	btcTSSAddress, err := btcutil.DecodeAddress(res.Btc, common.BitcoinRegnetParams)
	if err != nil {
		panic(err)
	}

	sm.TSSAddress = tssAddress
	sm.BTCTSSAddress = btcTSSAddress
}

// SetZEVMContracts set contracts for the ZEVM
func (sm *SmokeTestRunner) SetZEVMContracts() {
	// deploy system contracts and ZRC20 contracts on ZetaChain
	uniswapV2FactoryAddr, uniswapV2RouterAddr, usdtZRC20Addr, err := sm.ZetaTxServer.DeploySystemContractsAndZRC20(
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

	// query system contract address from the chain
	systemContractRes, err := sm.FungibleClient.SystemContract(
		context.Background(),
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
}

func (sm *SmokeTestRunner) SetupZEVMSwapApp() {
	zevmSwapAppAddr, tx, zevmSwapApp, err := zevmswap.DeployZEVMSwapApp(
		sm.ZevmAuth,
		sm.ZevmClient,
		sm.UniswapV2RouterAddr,
		sm.SystemContractAddr,
	)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status != 1 {
		panic("ZEVMSwapApp deployment failed")
	}
	fmt.Printf("ZEVMSwapApp contract address: %s, tx hash: %s\n", zevmSwapAppAddr.Hex(), tx.Hash().Hex())
	sm.ZEVMSwapAppAddr = zevmSwapAppAddr
	sm.ZEVMSwapApp = zevmSwapApp
}

func (sm *SmokeTestRunner) SetupContextApp() {
	contextAppAddr, tx, contextApp, err := contextapp.DeployContextApp(sm.ZevmAuth, sm.ZevmClient)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status != 1 {
		panic("ContextApp deployment failed")
	}
	fmt.Printf("ContextApp contract address: %s, tx hash: %s\n", contextAppAddr.Hex(), tx.Hash().Hex())
	sm.ContextAppAddr = contextAppAddr
	sm.ContextApp = contextApp
}

// SetCoreParams sets the core params with local Goerli and BtcRegtest chains enabled
func (sm *SmokeTestRunner) SetCoreParams() error {
	// set btc regtest  core params
	btcCoreParams := observertypes.GetDefaultBtcRegtestCoreParams()
	btcCoreParams.IsSupported = true
	if err := sm.ZetaTxServer.UpdateCoreParams(utils.FungibleAdminName, btcCoreParams); err != nil {
		return fmt.Errorf("failed to set core params for bitcoin: %s", err.Error())
	}

	// set goerli localnet core params
	goerliCoreParams := observertypes.GetDefaultGoerliLocalnetCoreParams()
	goerliCoreParams.IsSupported = true
	if err := sm.ZetaTxServer.UpdateCoreParams(utils.FungibleAdminName, goerliCoreParams); err != nil {
		return fmt.Errorf("failed to set core params for bitcoin: %s", err.Error())
	}

	return nil
}
