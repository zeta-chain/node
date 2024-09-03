package config

import (
	"context"
	"fmt"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/runner"
)

// RunnerFromConfig create test runner from config
func RunnerFromConfig(
	ctx context.Context,
	name string,
	ctxCancel context.CancelFunc,
	conf config.Config,
	account config.Account,
	logger *runner.Logger,
	opts ...runner.E2ERunnerOption,
) (*runner.E2ERunner, error) {
	// initialize all clients for E2E tests
	e2eClients, err := getClientsFromConfig(ctx, conf, account)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients from config: %w", err)
	}

	// initialize E2E test runner
	newRunner := runner.NewE2ERunner(
		ctx,
		name,
		ctxCancel,
		account,
		e2eClients,

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
	conf.Contracts.EVM.ZetaEthAddr = config.DoubleQuotedString(r.ZetaEthAddr.Hex())
	conf.Contracts.EVM.ConnectorEthAddr = config.DoubleQuotedString(r.ConnectorEthAddr.Hex())
	conf.Contracts.EVM.CustodyAddr = config.DoubleQuotedString(r.ERC20CustodyAddr.Hex())
	conf.Contracts.EVM.ERC20 = config.DoubleQuotedString(r.ERC20Addr.Hex())
	conf.Contracts.EVM.TestDappAddr = config.DoubleQuotedString(r.EvmTestDAppAddr.Hex())

	conf.Contracts.ZEVM.SystemContractAddr = config.DoubleQuotedString(r.SystemContractAddr.Hex())
	conf.Contracts.ZEVM.ETHZRC20Addr = config.DoubleQuotedString(r.ETHZRC20Addr.Hex())
	conf.Contracts.ZEVM.ERC20ZRC20Addr = config.DoubleQuotedString(r.ERC20ZRC20Addr.Hex())
	conf.Contracts.ZEVM.BTCZRC20Addr = config.DoubleQuotedString(r.BTCZRC20Addr.Hex())
	conf.Contracts.ZEVM.SOLZRC20Addr = config.DoubleQuotedString(r.SOLZRC20Addr.Hex())
	conf.Contracts.ZEVM.UniswapFactoryAddr = config.DoubleQuotedString(r.UniswapV2FactoryAddr.Hex())
	conf.Contracts.ZEVM.UniswapRouterAddr = config.DoubleQuotedString(r.UniswapV2RouterAddr.Hex())
	conf.Contracts.ZEVM.ConnectorZEVMAddr = config.DoubleQuotedString(r.ConnectorZEVMAddr.Hex())
	conf.Contracts.ZEVM.WZetaAddr = config.DoubleQuotedString(r.WZetaAddr.Hex())
	conf.Contracts.ZEVM.ZEVMSwapAppAddr = config.DoubleQuotedString(r.ZEVMSwapAppAddr.Hex())
	conf.Contracts.ZEVM.ContextAppAddr = config.DoubleQuotedString(r.ContextAppAddr.Hex())
	conf.Contracts.ZEVM.TestDappAddr = config.DoubleQuotedString(r.ZevmTestDAppAddr.Hex())

	// v2
	conf.Contracts.EVM.Gateway = config.DoubleQuotedString(r.GatewayEVMAddr.Hex())
	conf.Contracts.EVM.ERC20CustodyNew = config.DoubleQuotedString(r.ERC20CustodyV2Addr.Hex())
	conf.Contracts.EVM.TestDAppV2Addr = config.DoubleQuotedString(r.TestDAppV2EVMAddr.Hex())

	conf.Contracts.ZEVM.Gateway = config.DoubleQuotedString(r.GatewayZEVMAddr.Hex())
	conf.Contracts.ZEVM.TestDAppV2Addr = config.DoubleQuotedString(r.TestDAppV2ZEVMAddr.Hex())

	return conf
}
