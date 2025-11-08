package types

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/tracing"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"

	"github.com/cosmos/evm/x/vm/statedb"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// Copied the Account and StorageResult types since they are registered under an
// internal pkg on geth.

// AccountResult struct for account proof
type AccountResult struct {
	Address      common.Address  `json:"address"`
	AccountProof []string        `json:"accountProof"`
	Balance      *hexutil.Big    `json:"balance"`
	CodeHash     common.Hash     `json:"codeHash"`
	Nonce        hexutil.Uint64  `json:"nonce"`
	StorageHash  common.Hash     `json:"storageHash"`
	StorageProof []StorageResult `json:"storageProof"`
}

// StorageResult defines the format for storage proof return
type StorageResult struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"`
	Proof []string     `json:"proof"`
}

// TxResultAdditionalFields allows to return additional cosmos EVM txs fields
type TxResultAdditionalFields struct {
	Value     *big.Int       `json:"amount"`
	Hash      common.Hash    `json:"hash"`
	TxHash    string         `json:"txHash"`
	Type      uint64         `json:"type"`
	Recipient common.Address `json:"recipient"`
	Sender    common.Address `json:"sender"`
	GasUsed   uint64         `json:"gasUsed"`
	GasLimit  *uint64        `json:"gasLimit"`
	Nonce     uint64         `json:"nonce"`
	Data      []byte         `json:"data"`
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	BlockHash           *common.Hash                    `json:"blockHash"`
	BlockNumber         *hexutil.Big                    `json:"blockNumber"`
	From                common.Address                  `json:"from"`
	Gas                 hexutil.Uint64                  `json:"gas"`
	GasPrice            *hexutil.Big                    `json:"gasPrice"`
	GasFeeCap           *hexutil.Big                    `json:"maxFeePerGas,omitempty"`
	GasTipCap           *hexutil.Big                    `json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerBlobGas    *hexutil.Big                    `json:"maxFeePerBlobGas,omitempty"`
	Hash                common.Hash                     `json:"hash"`
	Input               hexutil.Bytes                   `json:"input"`
	Nonce               hexutil.Uint64                  `json:"nonce"`
	To                  *common.Address                 `json:"to"`
	TransactionIndex    *hexutil.Uint64                 `json:"transactionIndex"`
	Value               *hexutil.Big                    `json:"value"`
	Type                hexutil.Uint64                  `json:"type"`
	Accesses            *ethtypes.AccessList            `json:"accessList,omitempty"`
	ChainID             *hexutil.Big                    `json:"chainId,omitempty"`
	BlobVersionedHashes []common.Hash                   `json:"blobVersionedHashes,omitempty"`
	AuthorizationList   []ethtypes.SetCodeAuthorization `json:"authorizationList,omitempty"`
	V                   *hexutil.Big                    `json:"v"`
	R                   *hexutil.Big                    `json:"r"`
	S                   *hexutil.Big                    `json:"s"`
	YParity             *hexutil.Uint64                 `json:"yParity,omitempty"`
}

// StateOverride is the collection of overridden accounts.
type StateOverride map[common.Address]OverrideAccount

func (diff *StateOverride) has(address common.Address) bool {
	_, ok := (*diff)[address]
	return ok
}

// Apply overrides the fields of specified accounts into the given state.
func (diff *StateOverride) Apply(db *statedb.StateDB, precompiles vm.PrecompiledContracts) error {
	if db == nil || diff == nil {
		return nil
	}
	// Tracks destinations of precompiles that were moved.
	dirtyAddrs := make(map[common.Address]struct{})
	for addr, account := range *diff {
		// If a precompile was moved to this address already, it can't be overridden.
		if _, ok := dirtyAddrs[addr]; ok {
			return fmt.Errorf("account %s has already been overridden by a precompile", addr.Hex())
		}
		p, isPrecompile := precompiles[addr]
		// The MoveTo feature makes it possible to move a precompile
		// code to another address. If the target address is another precompile
		// the code for the latter is lost for this session.
		// Note the destination account is not cleared upon move.
		if account.MovePrecompileTo != nil {
			if !isPrecompile {
				return fmt.Errorf("account %s is not a precompile", addr.Hex())
			}
			// Refuse to move a precompile to an address that has been
			// or will be overridden.
			if diff.has(*account.MovePrecompileTo) {
				return fmt.Errorf("account %s is already overridden", account.MovePrecompileTo.Hex())
			}
			precompiles[*account.MovePrecompileTo] = p
			dirtyAddrs[*account.MovePrecompileTo] = struct{}{}
		}
		if isPrecompile {
			delete(precompiles, addr)
		}
		// Override account nonce.
		if account.Nonce != nil {
			db.SetNonce(addr, uint64(*account.Nonce), tracing.NonceChangeUnspecified)
		}
		// Override account(contract) code.
		if account.Code != nil {
			db.SetCode(addr, *account.Code)
		}
		// Override account balance.
		if account.Balance != nil && *account.Balance != nil {
			u256Balance, _ := uint256.FromBig((*big.Int)(*account.Balance))
			db.SetBalance(addr, u256Balance, tracing.BalanceChangeUnspecified)
		}
		if account.State != nil && account.StateDiff != nil {
			return fmt.Errorf("account %s has both 'state' and 'stateDiff'", addr.Hex())
		}
		// Replace entire state if caller requires.
		if account.State != nil {
			db.SetStorage(addr, *account.State)
		}
		// Apply state diff into specified accounts.
		if account.StateDiff != nil {
			for key, value := range *account.StateDiff {
				db.SetState(addr, key, value)
			}
		}
	}

	// Now finalize the changes. Finalize is normally performed between transactions.
	// By using finalize, the overrides are semantically behaving as
	// if they were created in a transaction just before the tracing occur.
	db.Finalise(false)
	return nil
}

// OverrideAccount indicates the overriding fields of account during the execution of
// a message call.
// Note, state and stateDiff can't be specified at the same time. If state is
// set, message execution will only use the data in the given state. Otherwise
// if statDiff is set, all diff will be applied first and then execute the call
// message.
type OverrideAccount struct {
	Nonce            *hexutil.Uint64              `json:"nonce"`
	Code             *hexutil.Bytes               `json:"code"`
	Balance          **hexutil.Big                `json:"balance"`
	State            *map[common.Hash]common.Hash `json:"state"`
	StateDiff        *map[common.Hash]common.Hash `json:"stateDiff"`
	MovePrecompileTo *common.Address              `json:"movePrecompileToAddress"`
}

type FeeHistoryResult struct {
	OldestBlock      *hexutil.Big     `json:"oldestBlock"`
	Reward           [][]*hexutil.Big `json:"reward,omitempty"`
	BaseFee          []*hexutil.Big   `json:"baseFeePerGas,omitempty"`
	GasUsedRatio     []float64        `json:"gasUsedRatio"`
	BlobBaseFee      []*hexutil.Big   `json:"baseFeePerBlobGas,omitempty"`
	BlobGasUsedRatio []float64        `json:"blobGasUsedRatio,omitempty"`
}

// SignTransactionResult represents a RLP encoded signed transaction.
type SignTransactionResult struct {
	Raw hexutil.Bytes         `json:"raw"`
	Tx  *ethtypes.Transaction `json:"tx"`
}

type OneFeeHistory struct {
	BaseFee, NextBaseFee         *big.Int   // base fee for each block
	Reward                       []*big.Int // each element of the array will have the tip provided to miners for the percentile given
	GasUsedRatio                 float64    // the ratio of gas used to the gas limit for each block
	BlobBaseFee, NextBlobBaseFee *big.Int   // blob base fee for each block
	BlobGasUsedRatio             float64    // the ratio of blob gas used to the blob gas limit for each block
}

// AccessListResult represents the access list and gas used for a transaction
type AccessListResult struct {
	AccessList *ethtypes.AccessList `json:"accessList"`
	GasUsed    *hexutil.Uint64      `json:"gasUsed"`
	Error      string               `json:"error,omitempty"`
}

// Embedded TraceConfig type to store raw JSON data of config in custom field
type TraceConfig struct {
	evmtypes.TraceConfig
	TracerConfig json.RawMessage `json:"tracerConfig"`
}
