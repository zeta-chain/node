package runner

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go/rpc"
	zetacore_rpc "github.com/zeta-chain/zetacore/pkg/rpc"
)

// Clients contains all the RPC clients and gRPC clients for E2E tests
type Clients struct {
	Zetacore zetacore_rpc.Clients

	// the RPC clients for external chains in the localnet
	BtcRPC  *rpcclient.Client
	Solana  *rpc.Client
	Evm     *ethclient.Client
	EvmAuth *bind.TransactOpts

	// the RPC clients for ZetaChain
	Zevm     *ethclient.Client
	ZevmAuth *bind.TransactOpts
}
