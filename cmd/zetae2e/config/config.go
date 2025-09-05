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
	ctxCancel context.CancelCauseFunc,
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
	// copy contracts from deployer runner
	conf.Contracts.Solana.GatewayProgramID = config.DoubleQuotedString(r.GatewayProgram.String())
	conf.Contracts.Solana.SPLAddr = config.DoubleQuotedString(r.SPLAddr.String())
	conf.Contracts.Solana.ConnectedProgramID = config.DoubleQuotedString(r.ConnectedProgram.String())
	conf.Contracts.Solana.ConnectedSPLProgramID = config.DoubleQuotedString(r.ConnectedSPLProgram.String())

	conf.Contracts.TON.GatewayAccountID = config.DoubleQuotedString(r.TONGateway.ToRaw())

	if r.SuiGateway != nil {
		conf.Contracts.Sui.GatewayPackageID = config.DoubleQuotedString(r.SuiGateway.PackageID())
		conf.Contracts.Sui.GatewayObjectID = config.DoubleQuotedString(r.SuiGateway.ObjectID())
	}
	conf.Contracts.Sui.GatewayUpgradeCap = config.DoubleQuotedString(r.SuiGatewayUpgradeCap)

	conf.Contracts.Sui.FungibleTokenCoinType = config.DoubleQuotedString(r.SuiTokenCoinType)
	conf.Contracts.Sui.FungibleTokenTreasuryCap = config.DoubleQuotedString(r.SuiTokenTreasuryCap)
	conf.Contracts.Sui.Example = r.SuiExample

	conf.Contracts.EVM.ZetaEthAddr = config.DoubleQuotedString(r.ZetaEthAddr.Hex())
	conf.Contracts.EVM.ConnectorEthAddr = config.DoubleQuotedString(r.ConnectorEthAddr.Hex())
	conf.Contracts.EVM.CustodyAddr = config.DoubleQuotedString(r.ERC20CustodyAddr.Hex())
	conf.Contracts.EVM.ERC20 = config.DoubleQuotedString(r.ERC20Addr.Hex())
	conf.Contracts.EVM.TestDappAddr = config.DoubleQuotedString(r.EvmTestDAppAddr.Hex())
	conf.Contracts.EVM.Gateway = config.DoubleQuotedString(r.GatewayEVMAddr.Hex())
	conf.Contracts.EVM.TestDAppV2Addr = config.DoubleQuotedString(r.TestDAppV2EVMAddr.Hex())
	conf.Contracts.EVM.ConnectorNative = config.DoubleQuotedString(r.ConnectorNativeAddr.Hex())

	conf.Contracts.ZEVM.SystemContractAddr = config.DoubleQuotedString(r.SystemContractAddr.Hex())
	conf.Contracts.ZEVM.ETHZRC20Addr = config.DoubleQuotedString(r.ETHZRC20Addr.Hex())
	conf.Contracts.ZEVM.ERC20ZRC20Addr = config.DoubleQuotedString(r.ERC20ZRC20Addr.Hex())
	conf.Contracts.ZEVM.BTCZRC20Addr = config.DoubleQuotedString(r.BTCZRC20Addr.Hex())
	conf.Contracts.ZEVM.SOLZRC20Addr = config.DoubleQuotedString(r.SOLZRC20Addr.Hex())
	conf.Contracts.ZEVM.SPLZRC20Addr = config.DoubleQuotedString(r.SPLZRC20Addr.Hex())
	conf.Contracts.ZEVM.TONZRC20Addr = config.DoubleQuotedString(r.TONZRC20Addr.Hex())
	conf.Contracts.ZEVM.SUIZRC20Addr = config.DoubleQuotedString(r.SUIZRC20Addr.Hex())
	conf.Contracts.ZEVM.SuiTokenZRC20Addr = config.DoubleQuotedString(r.SuiTokenZRC20Addr.Hex())

	conf.Contracts.ZEVM.UniswapFactoryAddr = config.DoubleQuotedString(r.UniswapV2FactoryAddr.Hex())
	conf.Contracts.ZEVM.UniswapRouterAddr = config.DoubleQuotedString(r.UniswapV2RouterAddr.Hex())
	conf.Contracts.ZEVM.ConnectorZEVMAddr = config.DoubleQuotedString(r.ConnectorZEVMAddr.Hex())
	conf.Contracts.ZEVM.WZetaAddr = config.DoubleQuotedString(r.WZetaAddr.Hex())
	conf.Contracts.ZEVM.ZEVMSwapAppAddr = config.DoubleQuotedString(r.ZEVMSwapAppAddr.Hex())
	conf.Contracts.ZEVM.TestDappAddr = config.DoubleQuotedString(r.ZevmTestDAppAddr.Hex())
	conf.Contracts.ZEVM.Gateway = config.DoubleQuotedString(r.GatewayZEVMAddr.Hex())
	conf.Contracts.ZEVM.TestDAppV2Addr = config.DoubleQuotedString(r.TestDAppV2ZEVMAddr.Hex())
	conf.Contracts.ZEVM.CoreRegistry = config.DoubleQuotedString(r.CoreRegistryAddr.Hex())

	return conf
}
