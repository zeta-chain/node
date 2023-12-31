package config

import (
	"context"

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
		return nil, err
	}
	// initialize client to send messages to ZetaChain
	zetaTxServer, err := txserver.NewZetaTxServer(
		conf.RPCs.ZetaCoreRPC,
		[]string{zetaUserName},
		[]string{zetaUserMnemonic},
		conf.ZetaChainID,
	)
	if err != nil {
		return nil, err
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
	return sm, nil
}
