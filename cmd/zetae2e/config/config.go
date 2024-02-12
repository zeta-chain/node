package config

import (
	"context"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/txserver"
)

// RunnerFromConfig create test runner from config
func RunnerFromConfig(
	ctx context.Context,
	name string,
	ctxCancel context.CancelFunc,
	conf config.Config,
	evmUserAddr ethcommon.Address,
	evmUserPrivKey string,
	zetaUserName string,
	zetaUserMnemonic string,
	logger *runner.Logger,
) (*runner.E2ERunner, error) {
	// initialize clients
	btcRPCClient,
		goerliClient,
		goerliAuth,
		cctxClient,
		fungibleClient,
		authClient,
		bankClient,
		observerClient,
		zevmClient,
		zevmAuth,
		err := getClientsFromConfig(ctx, conf, evmUserPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients from config: %w", err)
	}
	// initialize client to send messages to ZetaChain
	zetaTxServer, err := txserver.NewZetaTxServer(
		conf.RPCs.ZetaCoreRPC,
		[]string{zetaUserName},
		[]string{zetaUserMnemonic},
		conf.ZetaChainID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ZetaChain tx server: %w", err)
	}

	// initialize E2E test runner
	newRunner := runner.NewE2ERunner(
		ctx,
		name,
		ctxCancel,
		evmUserAddr,
		evmUserPrivKey,
		zetaUserMnemonic,
		goerliClient,
		zevmClient,
		cctxClient,
		zetaTxServer,
		fungibleClient,
		authClient,
		bankClient,
		observerClient,
		goerliAuth,
		zevmAuth,
		btcRPCClient,
		logger,
	)

	// set contracts
	err = setContractsFromConfig(newRunner, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to set contracts from config: %w", err)
	}

	// set bitcoin params
	chainParams, err := conf.RPCs.Bitcoin.Params.GetParams()
	if err != nil {
		return nil, fmt.Errorf("failed to get bitcoin params: %w", err)
	}
	newRunner.BitcoinParams = &chainParams

	return newRunner, err
}

// ExportContractsFromRunner export contracts from the runner to config using a source config
func ExportContractsFromRunner(oldRunner *runner.E2ERunner, conf config.Config) config.Config {
	// copy contracts from deployer runner
	conf.Contracts.EVM.ZetaEthAddress = oldRunner.ZetaEthAddr.Hex()
	conf.Contracts.EVM.ConnectorEthAddr = oldRunner.ConnectorEthAddr.Hex()
	conf.Contracts.EVM.CustodyAddr = oldRunner.ERC20CustodyAddr.Hex()
	conf.Contracts.EVM.USDT = oldRunner.USDTERC20Addr.Hex()

	conf.Contracts.ZEVM.SystemContractAddr = oldRunner.SystemContractAddr.Hex()
	conf.Contracts.ZEVM.ETHZRC20Addr = oldRunner.ETHZRC20Addr.Hex()
	conf.Contracts.ZEVM.USDTZRC20Addr = oldRunner.USDTZRC20Addr.Hex()
	conf.Contracts.ZEVM.BTCZRC20Addr = oldRunner.BTCZRC20Addr.Hex()
	conf.Contracts.ZEVM.UniswapFactoryAddr = oldRunner.UniswapV2FactoryAddr.Hex()
	conf.Contracts.ZEVM.UniswapRouterAddr = oldRunner.UniswapV2RouterAddr.Hex()
	conf.Contracts.ZEVM.ZEVMSwapAppAddr = oldRunner.ZEVMSwapAppAddr.Hex()
	conf.Contracts.ZEVM.ContextAppAddr = oldRunner.ContextAppAddr.Hex()
	conf.Contracts.ZEVM.TestDappAddr = oldRunner.TestDAppAddr.Hex()

	return conf
}
