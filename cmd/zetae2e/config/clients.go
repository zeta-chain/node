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
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zeta-chain/zetacore/e2e/config"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// E2EClients contains all the RPC clients and gRPC clients for E2E tests
type E2EClients struct {
	// the RPC clients for external chains in the localnet
	BtcRPCClient *rpcclient.Client
	SolanaClient *rpc.Client
	EvmClient    *ethclient.Client
	EvmAuth      *bind.TransactOpts

	// the gRPC clients for ZetaChain
	AuthorityClient authoritytypes.QueryClient
	CctxClient      crosschaintypes.QueryClient
	FungibleClient  fungibletypes.QueryClient
	AuthClient      authtypes.QueryClient
	BankClient      banktypes.QueryClient
	ObserverClient  observertypes.QueryClient
	LightClient     lightclienttypes.QueryClient

	// the RPC clients for ZetaChain
	ZevmClient *ethclient.Client
	ZevmAuth   *bind.TransactOpts
}

// zetaChainClients contains all the RPC clients and gRPC clients for ZetaChain
type zetaChainClients struct {
	AuthorityClient authoritytypes.QueryClient
	CctxClient      crosschaintypes.QueryClient
	FungibleClient  fungibletypes.QueryClient
	AuthClient      authtypes.QueryClient
	BankClient      banktypes.QueryClient
	ObserverClient  observertypes.QueryClient
	LightClient     lightclienttypes.QueryClient
}

// getClientsFromConfig get clients from config
func getClientsFromConfig(ctx context.Context, conf config.Config, account config.Account) (
	E2EClients,
	error,
) {
	var solanaClient *rpc.Client
	if conf.RPCs.Solana != "" {
		if solanaClient = rpc.New(conf.RPCs.Solana); solanaClient == nil {
			return E2EClients{}, fmt.Errorf("failed to get solana client")
		}
	}
	btcRPCClient, err := getBtcClient(conf.RPCs.Bitcoin)
	if err != nil {
		return E2EClients{}, fmt.Errorf("failed to get btc client: %w", err)
	}
	evmClient, evmAuth, err := getEVMClient(ctx, conf.RPCs.EVM, account)
	if err != nil {
		return E2EClients{}, fmt.Errorf("failed to get evm client: %w", err)
	}
	zetaChainClients, err := getZetaClients(
		conf.RPCs.ZetaCoreGRPC,
	)
	if err != nil {
		return E2EClients{}, fmt.Errorf("failed to get zeta clients: %w", err)
	}
	zevmClient, zevmAuth, err := getEVMClient(ctx, conf.RPCs.Zevm, account)
	if err != nil {
		return E2EClients{}, fmt.Errorf("failed to get zevm client: %w", err)
	}

	return E2EClients{
		BtcRPCClient:    btcRPCClient,
		SolanaClient:    solanaClient,
		EvmClient:       evmClient,
		EvmAuth:         evmAuth,
		AuthorityClient: zetaChainClients.AuthorityClient,
		CctxClient:      zetaChainClients.CctxClient,
		FungibleClient:  zetaChainClients.FungibleClient,
		AuthClient:      zetaChainClients.AuthClient,
		BankClient:      zetaChainClients.BankClient,
		ObserverClient:  zetaChainClients.ObserverClient,
		LightClient:     zetaChainClients.LightClient,
		ZevmClient:      zevmClient,
		ZevmAuth:        zevmAuth,
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

// getZetaClients get zeta clients
func getZetaClients(rpc string) (
	zetaChainClients,
	error,
) {
	grpcConn, err := grpc.Dial(rpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return zetaChainClients{}, err
	}

	authorityClient := authoritytypes.NewQueryClient(grpcConn)
	cctxClient := crosschaintypes.NewQueryClient(grpcConn)
	fungibleClient := fungibletypes.NewQueryClient(grpcConn)
	authClient := authtypes.NewQueryClient(grpcConn)
	bankClient := banktypes.NewQueryClient(grpcConn)
	observerClient := observertypes.NewQueryClient(grpcConn)
	lightclientClient := lightclienttypes.NewQueryClient(grpcConn)

	return zetaChainClients{
		AuthorityClient: authorityClient,
		CctxClient:      cctxClient,
		FungibleClient:  fungibleClient,
		AuthClient:      authClient,
		BankClient:      bankClient,
		ObserverClient:  observerClient,
		LightClient:     lightclientClient,
	}, nil
}
