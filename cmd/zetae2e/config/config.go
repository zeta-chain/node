package config

import (
	"context"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/txserver"
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
) (*runner.SmokeTestRunner, error) {
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

	// initialize smoke test runner
	sm := runner.NewSmokeTestRunner(
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
	err = setContractsFromConfig(sm, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to set contracts from config: %w", err)
	}

	// set bitcoin params
	chainParams, err := conf.RPCs.Bitcoin.Params.GetParams()
	if err != nil {
		return nil, fmt.Errorf("failed to get bitcoin params: %w", err)
	}
	sm.BitcoinParams = &chainParams

	return sm, err
}

// ExportContractsFromRunner export contracts from the runner to config using a source config
func ExportContractsFromRunner(sm *runner.SmokeTestRunner, conf config.Config) config.Config {
	// copy contracts from deployer runner
	conf.Contracts.EVM.ZetaEthAddress = sm.ZetaEthAddr.Hex()
	conf.Contracts.EVM.ConnectorEthAddr = sm.ConnectorEthAddr.Hex()
	conf.Contracts.EVM.CustodyAddr = sm.ERC20CustodyAddr.Hex()
	conf.Contracts.EVM.USDT = sm.USDTERC20Addr.Hex()

	conf.Contracts.ZEVM.SystemContractAddr = sm.SystemContractAddr.Hex()
	conf.Contracts.ZEVM.ETHZRC20Addr = sm.ETHZRC20Addr.Hex()
	conf.Contracts.ZEVM.USDTZRC20Addr = sm.USDTZRC20Addr.Hex()
	conf.Contracts.ZEVM.BTCZRC20Addr = sm.BTCZRC20Addr.Hex()
	conf.Contracts.ZEVM.UniswapFactoryAddr = sm.UniswapV2FactoryAddr.Hex()
	conf.Contracts.ZEVM.UniswapRouterAddr = sm.UniswapV2RouterAddr.Hex()
	conf.Contracts.ZEVM.ZEVMSwapAppAddr = sm.ZEVMSwapAppAddr.Hex()
	conf.Contracts.ZEVM.ContextAppAddr = sm.ContextAppAddr.Hex()
	conf.Contracts.ZEVM.TestDappAddr = sm.TestDAppAddr.Hex()

	return conf
}
