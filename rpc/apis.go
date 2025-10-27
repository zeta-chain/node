package rpc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"

	evmmempool "github.com/cosmos/evm/mempool"
	servertypes "github.com/cosmos/evm/server/types"
	"github.com/zeta-chain/node/rpc/backend"
	"github.com/zeta-chain/node/rpc/namespaces/ethereum/debug"
	"github.com/zeta-chain/node/rpc/namespaces/ethereum/eth"
	"github.com/zeta-chain/node/rpc/namespaces/ethereum/eth/filters"
	"github.com/zeta-chain/node/rpc/namespaces/ethereum/miner"
	"github.com/zeta-chain/node/rpc/namespaces/ethereum/net"
	"github.com/zeta-chain/node/rpc/namespaces/ethereum/personal"
	"github.com/zeta-chain/node/rpc/namespaces/ethereum/txpool"
	"github.com/zeta-chain/node/rpc/namespaces/ethereum/web3"
	"github.com/zeta-chain/node/rpc/stream"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
)

// RPC namespaces and API version
const (
	// Cosmos namespaces

	CosmosNamespace = "cosmos"

	// Ethereum namespaces

	Web3Namespace     = "web3"
	EthNamespace      = "eth"
	PersonalNamespace = "personal"
	NetNamespace      = "net"
	TxPoolNamespace   = "txpool"
	DebugNamespace    = "debug"
	MinerNamespace    = "miner"

	apiVersion = "1.0"
)

// APICreator creates the JSON-RPC API implementations.
type APICreator = func(
	ctx *server.Context,
	clientCtx client.Context,
	stream *stream.RPCStream,
	allowUnprotectedTxs bool,
	indexer servertypes.EVMTxIndexer,
	mempool *evmmempool.ExperimentalEVMMempool,
) []rpc.API

// apiCreators defines the JSON-RPC API namespaces.
var apiCreators map[string]APICreator

func init() {
	apiCreators = map[string]APICreator{
		EthNamespace: func(ctx *server.Context,
			clientCtx client.Context,
			stream *stream.RPCStream,
			allowUnprotectedTxs bool,
			indexer servertypes.EVMTxIndexer,
			mempool *evmmempool.ExperimentalEVMMempool,
		) []rpc.API {
			evmBackend := backend.NewBackend(ctx, ctx.Logger, clientCtx, allowUnprotectedTxs, indexer, mempool)
			return []rpc.API{
				{
					Namespace: EthNamespace,
					Version:   apiVersion,
					Service:   eth.NewPublicAPI(ctx.Logger, evmBackend),
					Public:    true,
				},
				{
					Namespace: EthNamespace,
					Version:   apiVersion,
					Service:   filters.NewPublicAPI(ctx.Logger, clientCtx, stream, evmBackend),
					Public:    true,
				},
			}
		},
		Web3Namespace: func(*server.Context, client.Context, *stream.RPCStream, bool, servertypes.EVMTxIndexer, *evmmempool.ExperimentalEVMMempool) []rpc.API {
			return []rpc.API{
				{
					Namespace: Web3Namespace,
					Version:   apiVersion,
					Service:   web3.NewPublicAPI(),
					Public:    true,
				},
			}
		},
		NetNamespace: func(ctx *server.Context, clientCtx client.Context, _ *stream.RPCStream, _ bool, _ servertypes.EVMTxIndexer, _ *evmmempool.ExperimentalEVMMempool) []rpc.API {
			return []rpc.API{
				{
					Namespace: NetNamespace,
					Version:   apiVersion,
					Service:   net.NewPublicAPI(ctx, clientCtx),
					Public:    true,
				},
			}
		},
		PersonalNamespace: func(ctx *server.Context,
			clientCtx client.Context,
			_ *stream.RPCStream,
			allowUnprotectedTxs bool,
			indexer servertypes.EVMTxIndexer,
			mempool *evmmempool.ExperimentalEVMMempool,
		) []rpc.API {
			evmBackend := backend.NewBackend(ctx, ctx.Logger, clientCtx, allowUnprotectedTxs, indexer, mempool)
			return []rpc.API{
				{
					Namespace: PersonalNamespace,
					Version:   apiVersion,
					Service:   personal.NewAPI(ctx.Logger, evmBackend),
					Public:    false,
				},
			}
		},
		TxPoolNamespace: func(ctx *server.Context,
			clientCtx client.Context,
			_ *stream.RPCStream,
			allowUnprotectedTxs bool,
			indexer servertypes.EVMTxIndexer,
			mempool *evmmempool.ExperimentalEVMMempool,
		) []rpc.API {
			evmBackend := backend.NewBackend(ctx, ctx.Logger, clientCtx, allowUnprotectedTxs, indexer, mempool)
			return []rpc.API{
				{
					Namespace: TxPoolNamespace,
					Version:   apiVersion,
					Service:   txpool.NewPublicAPI(ctx.Logger, evmBackend),
					Public:    true,
				},
			}
		},
		DebugNamespace: func(ctx *server.Context,
			clientCtx client.Context,
			_ *stream.RPCStream,
			allowUnprotectedTxs bool,
			indexer servertypes.EVMTxIndexer,
			mempool *evmmempool.ExperimentalEVMMempool,
		) []rpc.API {
			evmBackend := backend.NewBackend(ctx, ctx.Logger, clientCtx, allowUnprotectedTxs, indexer, mempool)
			return []rpc.API{
				{
					Namespace: DebugNamespace,
					Version:   apiVersion,
					Service:   debug.NewAPI(ctx, evmBackend, evmBackend.GetConfig().JSONRPC.EnableProfiling),
					Public:    true,
				},
			}
		},
		MinerNamespace: func(ctx *server.Context,
			clientCtx client.Context,
			_ *stream.RPCStream,
			allowUnprotectedTxs bool,
			indexer servertypes.EVMTxIndexer,
			mempool *evmmempool.ExperimentalEVMMempool,
		) []rpc.API {
			evmBackend := backend.NewBackend(ctx, ctx.Logger, clientCtx, allowUnprotectedTxs, indexer, mempool)
			return []rpc.API{
				{
					Namespace: MinerNamespace,
					Version:   apiVersion,
					Service:   miner.NewPrivateAPI(ctx, evmBackend),
					Public:    false,
				},
			}
		},
	}
}

// GetRPCAPIs returns the list of all APIs
func GetRPCAPIs(ctx *server.Context,
	clientCtx client.Context,
	stream *stream.RPCStream,
	allowUnprotectedTxs bool,
	indexer servertypes.EVMTxIndexer,
	selectedAPIs []string,
	mempool *evmmempool.ExperimentalEVMMempool,
) []rpc.API {
	var apis []rpc.API

	for _, ns := range selectedAPIs {
		if creator, ok := apiCreators[ns]; ok {
			apis = append(apis, creator(ctx, clientCtx, stream, allowUnprotectedTxs, indexer, mempool)...)
		} else {
			ctx.Logger.Error("invalid namespace value", "namespace", ns)
		}
	}

	return apis
}

// RegisterAPINamespace registers a new API namespace with the API creator.
// This function fails if the namespace is already registered.
func RegisterAPINamespace(ns string, creator APICreator) error {
	if _, ok := apiCreators[ns]; ok {
		return fmt.Errorf("duplicated api namespace %s", ns)
	}
	apiCreators[ns] = creator
	return nil
}
