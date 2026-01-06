package types

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtrpcclient "github.com/cometbft/cometbft/rpc/client"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

// ExceedBlockGasLimitError defines the error message when tx execution exceeds the block gas limit.
// The tx fee is deducted in ante handler, so it shouldn't be ignored in JSON-RPC API.
const ExceedBlockGasLimitError = "out of gas in location: block gas meter; gasWanted:"

// StateDBCommitError defines the error message when commit after executing EVM transaction, for example
// transfer native token to a distribution module account 0x93354845030274cD4bf1686Abd60AB28EC52e1a7 using an evm type transaction
// note: the transfer amount cannot be set to 0, otherwise this problem will not be triggered
const StateDBCommitError = "failed to commit stateDB"

// RawTxToEthTx returns a evm MsgEthereum transaction from raw tx bytes.
func RawTxToEthTx(clientCtx client.Context, txBz cmttypes.Tx) ([]*evmtypes.MsgEthereumTx, error) {
	tx, err := clientCtx.TxConfig.TxDecoder()(txBz)
	if err != nil {
		return nil, errorsmod.Wrap(errortypes.ErrJSONUnmarshal, err.Error())
	}

	ethTxs := make([]*evmtypes.MsgEthereumTx, len(tx.GetMsgs()))
	for i, msg := range tx.GetMsgs() {
		ethTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return nil, fmt.Errorf("invalid message type %T, expected %T", msg, &evmtypes.MsgEthereumTx{})
		}
		ethTx.Hash = ethTx.AsTransaction().Hash().Hex()
		ethTxs[i] = ethTx
	}
	return ethTxs, nil
}

// EthHeaderFromTendermint is an util function that returns an Ethereum Header
// from a tendermint Header.
func EthHeaderFromTendermint(header cmttypes.Header, bloom ethtypes.Bloom, baseFee *big.Int) *ethtypes.Header {
	txHash := ethtypes.EmptyRootHash
	if len(header.DataHash) != 0 {
		txHash = common.BytesToHash(header.DataHash)
	}

	time := uint64(header.Time.UTC().Unix()) //#nosec G115 won't exceed uint64
	return &ethtypes.Header{
		ParentHash:  common.BytesToHash(header.LastBlockID.Hash.Bytes()),
		UncleHash:   ethtypes.EmptyUncleHash,
		Coinbase:    common.BytesToAddress(header.ProposerAddress),
		Root:        common.BytesToHash(header.AppHash),
		TxHash:      txHash,
		ReceiptHash: ethtypes.EmptyRootHash,
		Bloom:       bloom,
		Difficulty:  big.NewInt(0),
		Number:      big.NewInt(header.Height),
		GasLimit:    0,
		GasUsed:     0,
		Time:        time,
		Extra:       []byte{},
		MixDigest:   common.Hash{},
		Nonce:       ethtypes.BlockNonce{},
		BaseFee:     baseFee,
	}
}

// BlockMaxGasFromConsensusParams returns the gas limit for the current block from the chain consensus params.
func BlockMaxGasFromConsensusParams(goCtx context.Context, clientCtx client.Context, blockHeight int64) (int64, error) {
	tmrpcClient, ok := clientCtx.Client.(cmtrpcclient.Client)
	if !ok {
		panic("incorrect tm rpc client")
	}
	resConsParams, err := tmrpcClient.ConsensusParams(goCtx, &blockHeight)
	defaultGasLimit := int64(^uint32(0)) // #nosec G115
	if err != nil {
		return defaultGasLimit, err
	}

	gasLimit := resConsParams.ConsensusParams.Block.MaxGas
	if gasLimit == -1 {
		// Sets gas limit to max uint32 to not error with javascript dev tooling
		// This -1 value indicating no block gas limit is set to max uint64 with geth hexutils
		// which errors certain javascript dev tooling which only supports up to 53 bits
		gasLimit = defaultGasLimit
	}

	return gasLimit, nil
}

// FormatBlock creates an ethereum block from a tendermint header and ethereum-formatted
// transactions.
func FormatBlock(
	header cmttypes.Header, size int, gasLimit int64,
	gasUsed *big.Int, transactions []interface{}, bloom ethtypes.Bloom,
	validatorAddr common.Address, baseFee *big.Int,
) map[string]interface{} {
	var transactionsRoot common.Hash
	if len(transactions) == 0 {
		transactionsRoot = ethtypes.EmptyRootHash
	} else {
		transactionsRoot = common.BytesToHash(header.DataHash)
	}

	result := map[string]interface{}{
		"number":           hexutil.Uint64(header.Height), //#nosec G115 won't exceed uint64
		"hash":             hexutil.Bytes(header.Hash()),
		"parentHash":       common.BytesToHash(header.LastBlockID.Hash.Bytes()),
		"nonce":            ethtypes.BlockNonce{},   // PoW specific
		"sha3Uncles":       ethtypes.EmptyUncleHash, // No uncles in Tendermint
		"logsBloom":        bloom,
		"stateRoot":        hexutil.Bytes(header.AppHash),
		"miner":            validatorAddr,
		"mixHash":          common.Hash{},
		"difficulty":       (*hexutil.Big)(big.NewInt(0)),
		"extraData":        "0x",
		"size":             hexutil.Uint64(size),     //#nosec G115 size won't exceed uint64
		"gasLimit":         hexutil.Uint64(gasLimit), //#nosec G115 gas limit won't exceed uint64
		"gasUsed":          (*hexutil.Big)(gasUsed),
		"timestamp":        hexutil.Uint64(header.Time.Unix()), //#nosec G115 won't exceed uint64
		"transactionsRoot": transactionsRoot,
		"receiptsRoot":     ethtypes.EmptyRootHash,

		"uncles":          []common.Hash{},
		"transactions":    transactions,
		"totalDifficulty": (*hexutil.Big)(big.NewInt(0)),
	}

	if baseFee != nil {
		result["baseFeePerGas"] = (*hexutil.Big)(baseFee)
	}

	return result
}

// NewTransactionFromMsg returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func NewTransactionFromMsg(
	msg *evmtypes.MsgEthereumTx,
	blockHash common.Hash,
	blockNumber, index uint64,
	baseFee *big.Int,
	chainID *big.Int,
	txAdditional *TxResultAdditionalFields,
) (*RPCTransaction, error) {
	if txAdditional != nil {
		return NewRPCTransactionFromIncompleteMsg(msg, blockHash, blockNumber, index, baseFee, chainID, txAdditional)
	}
	return NewRPCTransaction(msg, blockHash, blockNumber, index, baseFee, chainID)
}

// NewRPCTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func NewRPCTransaction(
	msg *evmtypes.MsgEthereumTx,
	blockHash common.Hash,
	blockNumber,
	index uint64,
	baseFee,
	chainID *big.Int,
) (*RPCTransaction, error) {
	tx := msg.AsTransaction()
	// Determine the signer. For replay-protected transactions, use the most permissive
	// signer, because we assume that signers are backwards-compatible with old
	// transactions. For non-protected transactions, the frontier signer is used
	// because the latest signer will reject the unprotected transactions.
	var signer ethtypes.Signer
	if tx.Protected() {
		signer = ethtypes.LatestSignerForChainID(tx.ChainId())
	} else {
		signer = ethtypes.FrontierSigner{}
	}
	from, err := msg.GetSenderLegacy(signer)
	if err != nil {
		return nil, err
	}
	v, r, s := tx.RawSignatureValues()
	result := &RPCTransaction{
		Type:     hexutil.Uint64(tx.Type()),
		From:     from,
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Hash:     tx.Hash(),
		Input:    hexutil.Bytes(tx.Data()),
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Value()),
		V:        (*hexutil.Big)(v),
		R:        (*hexutil.Big)(r),
		S:        (*hexutil.Big)(s),
		ChainID:  (*hexutil.Big)(chainID),
	}
	if blockHash != (common.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}
	switch tx.Type() {
	case ethtypes.AccessListTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
	case ethtypes.DynamicFeeTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		if blockHash != (common.Hash{}) {
			result.GasPrice = (*hexutil.Big)(EffectiveGasPrice(tx, baseFee))
		} else {
			// For pending transactions, use gasFeeCap as placeholder
			result.GasPrice = (*hexutil.Big)(tx.GasFeeCap())
		}
	}

	return result, nil
}

// NewRPCTransactionFromIncompleteMsg returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func NewRPCTransactionFromIncompleteMsg(
	msg *evmtypes.MsgEthereumTx, blockHash common.Hash, blockNumber, index uint64, baseFee *big.Int,
	chainID *big.Int, txAdditional *TxResultAdditionalFields,
) (*RPCTransaction, error) {
	to := &common.Address{}
	*to = txAdditional.Recipient
	// for transactions before v31 this value was mistakenly used for Gas field
	gas := txAdditional.GasUsed
	if txAdditional.GasLimit != nil {
		gas = *txAdditional.GasLimit
	}
	result := &RPCTransaction{
		Type:     hexutil.Uint64(txAdditional.Type),
		From:     common.BytesToAddress(msg.From),
		Gas:      hexutil.Uint64(gas),
		GasPrice: (*hexutil.Big)(baseFee),
		Hash:     common.HexToHash(msg.Hash),
		Input:    txAdditional.Data,
		Nonce:    hexutil.Uint64(txAdditional.Nonce), // TODO: get nonce for "from" from evm
		To:       to,
		Value:    (*hexutil.Big)(txAdditional.Value),
		V:        (*hexutil.Big)(big.NewInt(0)),
		R:        (*hexutil.Big)(big.NewInt(0)),
		S:        (*hexutil.Big)(big.NewInt(0)),
		ChainID:  (*hexutil.Big)(chainID),
	}
	if blockHash != (common.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}
	return result, nil
}

// BaseFeeFromEvents parses the feemarket basefee from cosmos events
func BaseFeeFromEvents(events []abci.Event) *big.Int {
	for _, event := range events {
		if event.Type != evmtypes.EventTypeFeeMarket {
			continue
		}

		for _, attr := range event.Attributes {
			if attr.Key == evmtypes.AttributeKeyBaseFee {
				result, success := sdkmath.NewIntFromString(attr.Value)
				if success {
					return result.BigInt()
				}

				return nil
			}
		}
	}
	return nil
}

// EffectiveGasPrice computes the transaction gas fee, based on the given baseFee value.
//
// For EIP-1559 transactions:
//
//	price = min(gasTipCap + baseFee, gasFeeCap)
//
// For legacy transactions:
//
//	price = gasPrice
//
// This method is based on go-ethereum's internal effectiveGasPrice calculation.
// (https://github.com/ethereum/go-ethereum/blob/d818a9af7bd5919808df78f31580f59382c53150/internal/ethapi/api.go#L1083-L1093)
func EffectiveGasPrice(tx *ethtypes.Transaction, baseFee *big.Int) *big.Int {
	if tx == nil {
		return big.NewInt(0)
	}

	// For legacy (0x00) and access list (0x01) transactions, return the gas price directly
	// For EIP-1559 (0x02), Blob (0x03), and SetCode (0x04) transactions, calculate effective gas price
	switch tx.Type() {
	case ethtypes.LegacyTxType, ethtypes.AccessListTxType:
		return tx.GasPrice()
	case ethtypes.DynamicFeeTxType, ethtypes.BlobTxType, ethtypes.SetCodeTxType:
		if baseFee == nil {
			return tx.GasFeeCap()
		}
		price := new(big.Int).Add(tx.GasTipCap(), baseFee)
		if price.Cmp(tx.GasFeeCap()) > 0 {
			return tx.GasFeeCap()
		}
		return price
	default:
		return tx.GasPrice()
	}
}

// CheckTxFee is an internal function used to check whether the fee of
// the given transaction is _reasonable_(under the minimum cap).
func CheckTxFee(gasPrice *big.Int, gas uint64, minCap float64) error {
	// Short circuit if there is no cap for transaction fee at all.
	if minCap == 0 {
		return nil
	}
	// Return an error if gasPrice is nil
	if gasPrice == nil {
		return errors.New("gasprice is nil")
	}

	totalfee := new(big.Float).SetInt(new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gas)))
	// 1 token in atto units (1e18)
	oneToken := new(big.Float).SetInt(big.NewInt(params.Ether))
	// quo = rounded(x/y)
	feeEth := new(big.Float).Quo(totalfee, oneToken)
	// no need to check error from parsing
	feeFloat, _ := feeEth.Float64()
	if feeFloat > minCap {
		return fmt.Errorf("tx fee (%.2f ether) exceeds the configured cap (%.2f ether)", feeFloat, minCap)
	}
	return nil
}

// TxExceedBlockGasLimit returns true if the tx exceeds block gas limit.
func TxExceedBlockGasLimit(res *abci.ExecTxResult) bool {
	return strings.Contains(res.Log, ExceedBlockGasLimitError)
}

// TxStateDBCommitError returns true if the evm tx commit error.
func TxStateDBCommitError(res *abci.ExecTxResult) bool {
	return strings.Contains(res.Log, StateDBCommitError)
}

// TxSucessOrExpectedFailure returns true if the transaction was successful
// or if it failed with an ExceedBlockGasLimit error or TxStateDBCommitError error
func TxSucessOrExpectedFailure(res *abci.ExecTxResult) bool {
	return res.Code == 0 || TxExceedBlockGasLimit(res) || TxStateDBCommitError(res)
}
