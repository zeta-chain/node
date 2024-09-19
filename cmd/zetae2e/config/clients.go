package config

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go/rpc"
	ton "github.com/tonkeeper/tongo/liteapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/runner"
	tonrunner "github.com/zeta-chain/node/e2e/runner/ton"
	"github.com/zeta-chain/node/pkg/retry"
	zetacore_rpc "github.com/zeta-chain/node/pkg/rpc"
	tonconfig "github.com/zeta-chain/node/zetaclient/chains/ton"
)

// getClientsFromConfig get clients from config
func getClientsFromConfig(ctx context.Context, conf config.Config, account config.Account) (runner.Clients, error) {
	btcRPCClient, err := getBtcClient(conf.RPCs.Bitcoin)
	if err != nil {
		return runner.Clients{}, fmt.Errorf("failed to get btc client: %w", err)
	}

	evmClient, evmAuth, err := getEVMClient(ctx, conf.RPCs.EVM, account)
	if err != nil {
		return runner.Clients{}, fmt.Errorf("failed to get evm client: %w", err)
	}

	var solanaClient *rpc.Client
	if conf.RPCs.Solana != "" {
		if solanaClient = rpc.New(conf.RPCs.Solana); solanaClient == nil {
			return runner.Clients{}, fmt.Errorf("failed to get solana client")
		}
	}

	var tonClient *tonrunner.Client
	if conf.RPCs.TONSidecarURL != "" {
		c, err := getTONClient(ctx, conf.RPCs.TONSidecarURL)
		if err != nil {
			return runner.Clients{}, fmt.Errorf("failed to get ton client: %w", err)
		}
		tonClient = c
	}

	zetaCoreClients, err := GetZetacoreClient(conf)
	if err != nil {
		return runner.Clients{}, fmt.Errorf("failed to get zetacore client: %w", err)
	}

	zevmClient, zevmAuth, err := getEVMClient(ctx, conf.RPCs.Zevm, account)
	if err != nil {
		return runner.Clients{}, fmt.Errorf("failed to get zevm client: %w", err)
	}

	return runner.Clients{
		Zetacore: zetaCoreClients,
		BtcRPC:   btcRPCClient,
		Solana:   solanaClient,
		TON:      tonClient,
		Evm:      evmClient,
		EvmAuth:  evmAuth,
		Zevm:     zevmClient,
		ZevmAuth: zevmAuth,
	}, nil
}

// getBtcClient get btc client
func getBtcClient(rpcConf config.BitcoinRPC) (*rpcclient.Client, error) {
	var param string
	switch rpcConf.Params {
	case config.Regnet:
	case config.Testnet3:
		param = "testnet3"
	case config.Mainnet:
		param = "mainnet"
	default:
		return nil, fmt.Errorf("invalid bitcoin params %s", rpcConf.Params)
	}

	connCfg := &rpcclient.ConnConfig{
		Host:         rpcConf.Host,
		User:         rpcConf.User,
		Pass:         rpcConf.Pass,
		HTTPPostMode: rpcConf.HTTPPostMode,
		DisableTLS:   rpcConf.DisableTLS,
		Params:       param,
	}
	return rpcclient.New(connCfg, nil)
}

// getEVMClient get evm client
func getEVMClient(
	ctx context.Context,
	rpc string,
	account config.Account,
) (*ethclient.Client, *bind.TransactOpts, error) {
	evmClient, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial evm client: %w", err)
	}

	chainid, err := evmClient.ChainID(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chain id: %w", err)
	}
	privKey, err := account.PrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get deployer privkey: %w", err)
	}
	evmAuth, err := bind.NewKeyedTransactorWithChainID(privKey, chainid)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get keyed transactor: %w", err)
	}

	return evmClient, evmAuth, nil
}

func getTONClient(ctx context.Context, sidecarURL string) (*tonrunner.Client, error) {
	if sidecarURL == "" {
		return nil, fmt.Errorf("sidecar URL is empty")
	}

	sidecar := tonrunner.NewSidecarClient(sidecarURL)

	// It might take some time to bootstrap the sidecar
	cfg, err := retry.DoTypedWithRetry(
		func() (*tonconfig.GlobalConfigurationFile, error) {
			return tonconfig.ConfigFromURL(ctx, sidecar.LiteServerURL())
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get ton config: %w", err)
	}

	client, err := ton.NewClient(ton.WithConfigurationFile(*cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to create ton client: %w", err)
	}

	return &tonrunner.Client{
		Client:        client,
		SidecarClient: sidecar,
	}, nil
}

func GetZetacoreClient(conf config.Config) (zetacore_rpc.Clients, error) {
	if conf.RPCs.ZetaCoreGRPC != "" {
		return zetacore_rpc.NewGRPCClients(
			conf.RPCs.ZetaCoreGRPC,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	}
	if conf.RPCs.ZetaCoreRPC != "" {
		return zetacore_rpc.NewCometBFTClients(conf.RPCs.ZetaCoreRPC)
	}
	return zetacore_rpc.Clients{}, fmt.Errorf("no ZetaCore gRPC or RPC specified")
}
