package config

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go/rpc"
	"google.golang.org/grpc"

	"github.com/zeta-chain/zetacore/e2e/config"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// getClientsFromConfig get clients from config
func getClientsFromConfig(ctx context.Context, conf config.Config, account config.Account) (
	*rpcclient.Client,
	*rpc.Client,
	*ethclient.Client,
	*bind.TransactOpts,
	crosschaintypes.QueryClient,
	fungibletypes.QueryClient,
	authtypes.QueryClient,
	banktypes.QueryClient,
	observertypes.QueryClient,
	lightclienttypes.QueryClient,
	*ethclient.Client,
	*bind.TransactOpts,
	error,
) {
	solanaClient := rpc.New(conf.RPCs.SolanaRPC)
	if solanaClient == nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to get solana client")
	}
	btcRPCClient, err := getBtcClient(conf.RPCs.Bitcoin)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to get btc client: %w", err)
	}
	evmClient, evmAuth, err := getEVMClient(ctx, conf.RPCs.EVM, account)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to get evm client: %w", err)
	}
	cctxClient, fungibleClient, authClient, bankClient, observerClient, lightclientClient, err := getZetaClients(
		conf.RPCs.ZetaCoreGRPC,
	)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to get zeta clients: %w", err)
	}
	zevmClient, zevmAuth, err := getEVMClient(ctx, conf.RPCs.Zevm, account)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to get zevm client: %w", err)
	}
	return btcRPCClient,
		solanaClient,
		evmClient,
		evmAuth,
		cctxClient,
		fungibleClient,
		authClient,
		bankClient,
		observerClient,
		lightclientClient,
		zevmClient,
		zevmAuth,
		nil
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

// getZetaClients get zeta clients
func getZetaClients(rpc string) (
	crosschaintypes.QueryClient,
	fungibletypes.QueryClient,
	authtypes.QueryClient,
	banktypes.QueryClient,
	observertypes.QueryClient,
	lightclienttypes.QueryClient,
	error,
) {
	grpcConn, err := grpc.Dial(rpc, grpc.WithInsecure())
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	cctxClient := crosschaintypes.NewQueryClient(grpcConn)
	fungibleClient := fungibletypes.NewQueryClient(grpcConn)
	authClient := authtypes.NewQueryClient(grpcConn)
	bankClient := banktypes.NewQueryClient(grpcConn)
	observerClient := observertypes.NewQueryClient(grpcConn)
	lightclientClient := lightclienttypes.NewQueryClient(grpcConn)

	return cctxClient, fungibleClient, authClient, bankClient, observerClient, lightclientClient, nil
}
