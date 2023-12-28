package local

import (
	"context"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/txserver"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc"
)

// getConfig returns config from file from the command line flag
func getConfig(cmd *cobra.Command) (config.Config, error) {
	configFile, err := cmd.Flags().GetString(flagConfigFile)
	if err != nil {
		return config.Config{}, err
	}

	// use default config if no config file is specified
	if configFile == "" {
		return config.DefaultConfig(), nil
	}

	return config.ReadConfig(configFile)
}

// setCosmosConfig set account prefix to zeta
func setCosmosConfig() {
	cosmosConf := sdk.GetConfig()
	cosmosConf.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cosmosConf.Seal()
}

// initTestRunner initializes a runner for smoke tests
// it creates a runner with an account and copy contracts from deployer runner
func initTestRunner(
	name string,
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	userAddress ethcommon.Address,
	userPrivKey string,
	logger *runner.Logger,
) (*runner.SmokeTestRunner, error) {
	// initialize runner for smoke test
	testRunner, err := runnerFromConfig(name, conf, userAddress, userPrivKey, logger)
	if err != nil {
		return nil, err
	}

	// copy contracts from deployer runner
	if err := testRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}

	return testRunner, nil
}

// runnerFromConfig create test runner from config
func runnerFromConfig(
	name string,
	conf config.Config,
	userAddr ethcommon.Address,
	userPrivKey string,
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
		err := getClientsFromConfig(conf, userPrivKey)
	if err != nil {
		return nil, err
	}
	// initialize client to send messages to ZetaChain
	zetaTxServer, err := txserver.NewZetaTxServer(
		conf.RPCs.ZetaCoreRPC,
		[]string{utils.FungibleAdminName},
		[]string{FungibleAdminMnemonic},
		conf.ZetaChainID,
	)
	if err != nil {
		return nil, err
	}

	// initialize smoke test runner
	sm := runner.NewSmokeTestRunner(
		name,
		userAddr,
		userPrivKey,
		FungibleAdminMnemonic,
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

// getClientsFromConfig get clients from config
func getClientsFromConfig(conf config.Config, evmPrivKey string) (
	*rpcclient.Client,
	*ethclient.Client,
	*bind.TransactOpts,
	crosschaintypes.QueryClient,
	fungibletypes.QueryClient,
	authtypes.QueryClient,
	banktypes.QueryClient,
	observertypes.QueryClient,
	*ethclient.Client,
	*bind.TransactOpts,
	error,
) {
	btcRPCClient, err := getBtcClient(conf.RPCs.Bitcoin)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	goerliClient, goerliAuth, err := getEVMClient(conf.RPCs.EVM, evmPrivKey)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	cctxClient, fungibleClient, authClient, bankClient, observerClient, err := getZetaClients(conf.RPCs.ZetaCoreGRPC)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	zevmClient, zevmAuth, err := getEVMClient(conf.RPCs.Zevm, evmPrivKey)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	return btcRPCClient,
		goerliClient,
		goerliAuth,
		cctxClient,
		fungibleClient,
		authClient,
		bankClient,
		observerClient,
		zevmClient,
		zevmAuth,
		nil
}

// getBtcClient get btc client
func getBtcClient(rpc string) (*rpcclient.Client, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         rpc,
		User:         "smoketest",
		Pass:         "123",
		HTTPPostMode: true,
		DisableTLS:   true,
		Params:       "testnet3",
	}
	return rpcclient.New(connCfg, nil)
}

// getEVMClient get goerli client
func getEVMClient(rpc, privKey string) (*ethclient.Client, *bind.TransactOpts, error) {
	evmClient, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, nil, err
	}

	chainid, err := evmClient.ChainID(context.Background())
	if err != nil {
		return nil, nil, err
	}
	deployerPrivkey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, nil, err
	}
	evmAuth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, chainid)
	if err != nil {
		return nil, nil, err
	}

	return evmClient, evmAuth, nil
}

// getZetaClients get zeta clients
func getZetaClients(rpc string) (
	crosschaintypes.QueryClient,
	fungibletypes.QueryClient,
	authtypes.QueryClient,
	banktypes.QueryClient,
	observertypes.QueryClient,
	error,
) {
	grpcConn, err := grpc.Dial(rpc, grpc.WithInsecure())
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	cctxClient := crosschaintypes.NewQueryClient(grpcConn)
	fungibleClient := fungibletypes.NewQueryClient(grpcConn)
	authClient := authtypes.NewQueryClient(grpcConn)
	bankClient := banktypes.NewQueryClient(grpcConn)
	observerClient := observertypes.NewQueryClient(grpcConn)

	return cctxClient, fungibleClient, authClient, bankClient, observerClient, nil
}

// waitKeygenHeight waits for keygen height
func waitKeygenHeight(
	cctxClient crosschaintypes.QueryClient,
	logger *runner.Logger,
) {
	// wait for keygen to be completed. ~ height 30
	keygenHeight := int64(60)
	logger.Print("⏳ wait height %v for keygen to be completed", keygenHeight)
	for {
		time.Sleep(2 * time.Second)
		response, err := cctxClient.LastZetaHeight(context.Background(), &crosschaintypes.QueryLastZetaHeightRequest{})
		if err != nil {
			logger.Error("cctxClient.LastZetaHeight error: %s", err)
			continue
		}
		if response.Height >= keygenHeight {
			break
		}
		logger.Info("Last ZetaHeight: %d", response.Height)
	}
}
