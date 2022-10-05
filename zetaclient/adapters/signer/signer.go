package signer

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type Signer interface {
	Broadcast(tx *ethtypes.Transaction) error
	SignOutboundTx(sender ethcommon.Address, srcChainID *big.Int, to ethcommon.Address, amount *big.Int, gasLimit uint64, message []byte, sendHash [32]byte, nonce uint64, gasPrice *big.Int) (*ethtypes.Transaction, error)
	SignRevertTx(sender ethcommon.Address, srcChainID *big.Int, to []byte, toChainID *big.Int, amount *big.Int, gasLimit uint64, message []byte, sendHash [32]byte, nonce uint64, gasPrice *big.Int) (*ethtypes.Transaction, error)
	Chain() string
}
