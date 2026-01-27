package config

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/runner"
	tonrunner "github.com/zeta-chain/node/e2e/runner/ton"
	"github.com/zeta-chain/node/pkg/chains"
	zetacore_rpc "github.com/zeta-chain/node/pkg/rpc"
	btcclient "github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	zetaclientconfig "github.com/zeta-chain/node/zetaclient/config"
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
	if conf.RPCs.TON != "" {
		tonClient = tonrunner.NewClient(conf.RPCs.TON)
	}

	var suiClient sui.ISuiAPI
	if conf.RPCs.Sui != "" {
		suiClient = sui.NewSuiClient(conf.RPCs.Sui)
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
		Zetacore:          zetaCoreClients,
		BtcRPC:            btcRPCClient,
		Solana:            solanaClient,
		TON:               tonClient,
		Sui:               suiClient,
		Evm:               evmClient,
		EvmAuth:           evmAuth,
		Zevm:              zevmClient,
		ZevmAuth:          zevmAuth,
		ZetaclientMetrics: &runner.MetricsClient{URL: conf.RPCs.ZetaclientMetrics},
	}, nil
}

// getBtcClient get btc client
func getBtcClient(e2eConfig config.BitcoinRPC) (*btcclient.Client, error) {
	cfg := zetaclientconfig.BTCConfig{
		RPCUsername: e2eConfig.User,
		RPCPassword: e2eConfig.Pass,
		RPCHost:     e2eConfig.Host,
		RPCParams:   string(e2eConfig.Params),
	}

	var chain chains.Chain
	switch e2eConfig.Params {
	case config.Regnet:
		chain = chains.BitcoinRegtest
	case config.Testnet3:
		chain = chains.BitcoinTestnet
	case config.Signet:
		chain = chains.BitcoinSignetTestnet
	case config.Mainnet:
		chain = chains.BitcoinMainnet
	default:
		return nil, fmt.Errorf("invalid bitcoin params %s", e2eConfig.Params)
	}

	return btcclient.New(cfg, chain.ChainId, zerolog.Nop())
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
