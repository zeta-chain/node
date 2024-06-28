package config

import (
	"context"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// RunnerFromConfig create test runner from config
func RunnerFromConfig(
	ctx context.Context,
	name string,
	ctxCancel context.CancelFunc,
	conf config.Config,
	evmUserAddr ethcommon.Address,
	evmUserPrivKey string,
	logger *runner.Logger,
	opts ...runner.E2ERunnerOption,
) (*runner.E2ERunner, error) {
	// initialize clients
	btcRPCClient,
		solanaClient,
		evmClient,
		evmAuth,
		cctxClient,
		fungibleClient,
		authClient,
		bankClient,
		observerClient,
		lightClient,
		zevmClient,
		zevmAuth,
		err := getClientsFromConfig(ctx, conf, evmUserPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients from config: %w", err)
	}

	// initialize E2E test runner
	newRunner := runner.NewE2ERunner(
		ctx,
		name,
		ctxCancel,
		evmUserAddr,
		evmUserPrivKey,
		evmClient,
		zevmClient,
		cctxClient,
		fungibleClient,
		authClient,
		bankClient,
		observerClient,
		lightClient,
		evmAuth,
		zevmAuth,
		btcRPCClient,
		solanaClient,

		logger,
		opts...,
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
func ExportContractsFromRunner(r *runner.E2ERunner, conf config.Config) config.Config {
	conf.Contracts.Solana.GatewayProgramID = r.GatewayProgram.String()

	// copy contracts from deployer runner
	conf.Contracts.EVM.ZetaEthAddress = r.ZetaEthAddr.Hex()
	conf.Contracts.EVM.ConnectorEthAddr = r.ConnectorEthAddr.Hex()
	conf.Contracts.EVM.CustodyAddr = r.ERC20CustodyAddr.Hex()
	conf.Contracts.EVM.ERC20 = r.ERC20Addr.Hex()
	conf.Contracts.EVM.TestDappAddr = r.EvmTestDAppAddr.Hex()

	conf.Contracts.ZEVM.SystemContractAddr = r.SystemContractAddr.Hex()
	conf.Contracts.ZEVM.ETHZRC20Addr = r.ETHZRC20Addr.Hex()
	conf.Contracts.ZEVM.ERC20ZRC20Addr = r.ERC20ZRC20Addr.Hex()
	conf.Contracts.ZEVM.BTCZRC20Addr = r.BTCZRC20Addr.Hex()
	conf.Contracts.ZEVM.UniswapFactoryAddr = r.UniswapV2FactoryAddr.Hex()
	conf.Contracts.ZEVM.UniswapRouterAddr = r.UniswapV2RouterAddr.Hex()
	conf.Contracts.ZEVM.ConnectorZEVMAddr = r.ConnectorZEVMAddr.Hex()
	conf.Contracts.ZEVM.WZetaAddr = r.WZetaAddr.Hex()
	conf.Contracts.ZEVM.ZEVMSwapAppAddr = r.ZEVMSwapAppAddr.Hex()
	conf.Contracts.ZEVM.ContextAppAddr = r.ContextAppAddr.Hex()
	conf.Contracts.ZEVM.TestDappAddr = r.ZevmTestDAppAddr.Hex()

	return conf
}
