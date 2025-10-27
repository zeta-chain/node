package types

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtrpcclient "github.com/cometbft/cometbft/rpc/client"
	cmtrpccore "github.com/cometbft/cometbft/rpc/core/types"
	cmttypes "github.com/cometbft/cometbft/types"

	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
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
		ethTxs[i] = ethTx
	}
	return ethTxs, nil
}

// EthHeaderFromComet is an util function that returns an Ethereum Header
// from a CometBFT Header.
//
// TODO: Remove this function.
// Currently, this function is only used in rpc/stream package for websocket api.
// But there are many missing fields in returned eth header.
// When kv_indexer is improved and we can get eth header from indexer, we can remove this function.
func EthHeaderFromComet(header cmttypes.Header, bloom ethtypes.Bloom, baseFee *big.Int) *ethtypes.Header {
	txHash := ethtypes.EmptyRootHash
	if len(header.DataHash) != 0 {
		txHash = common.BytesToHash(header.DataHash)
	}

	time := uint64(header.Time.UTC().Unix()) //nolint:gosec // G115 // won't exceed uint64
	return &ethtypes.Header{
		ParentHash:  common.BytesToHash(header.LastBlockID.Hash.Bytes()),
		UncleHash:   ethtypes.EmptyUncleHash,
		Coinbase:    common.BytesToAddress(header.ProposerAddress),
		Root:        common.BytesToHash(header.AppHash),
		TxHash:      txHash,
		ReceiptHash: ethtypes.EmptyReceiptsHash,
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

		// In chains that use cosmos-sdk and cometbft,
		// these fields are irrelevant.
		WithdrawalsHash:  &ethtypes.EmptyWithdrawalsHash, // EIP-4895: Beacon chain push withdrawals as operations
		BlobGasUsed:      new(uint64),                    // EIP-4844: Shard Blob Transactions
		ExcessBlobGas:    new(uint64),                    // EIP-4844: Shard Blob Transactions
		ParentBeaconRoot: &ethtypes.EmptyRootHash,        // EIP-4788: Beacon block root in the EVM
		RequestsHash:     &ethtypes.EmptyRequestsHash,    // EIP-7685: General purpose execution layer requests
	}
}

// BlockMaxGasFromConsensusParams returns the gas limit for the current block from the chain consensus params.
func BlockMaxGasFromConsensusParams(goCtx context.Context, clientCtx client.Context, blockHeight int64) (int64, error) {
	cmtrpcclient, ok := clientCtx.Client.(cmtrpcclient.Client)
	if !ok {
		panic("incorrect tm rpc client")
	}
	resConsParams, err := cmtrpcclient.ConsensusParams(goCtx, &blockHeight)
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

// MakeHeader make initial ethereum header based on cometbft header.
//
// This method refers to chainMaker.makeHeader method of go-ethereum v1.16.3
// (https://github.com/ethereum/go-ethereum/blob/d818a9af7bd5919808df78f31580f59382c53150/core/chain_makers.go#L596-L623)
func MakeHeader(
	cmtHeader cmttypes.Header, gasLimit int64,
	validatorAddr common.Address, baseFee *big.Int,
) *ethtypes.Header {
	header := &ethtypes.Header{
		Root:       common.BytesToHash(hexutil.Bytes(cmtHeader.AppHash)),
		ParentHash: common.BytesToHash(cmtHeader.LastBlockID.Hash.Bytes()),
		Coinbase:   validatorAddr,
		Difficulty: big.NewInt(0),
		GasLimit:   uint64(gasLimit), //nolint:gosec // G115 // gas limit won't exceed uint64
		Number:     big.NewInt(cmtHeader.Height),
		Time:       uint64(cmtHeader.Time.UTC().Unix()), //nolint:gosec // G115 // timestamp won't exceed uint64
	}

	if evmtypes.GetEthChainConfig().IsLondon(header.Number) {
		header.BaseFee = baseFee
	}
	if evmtypes.GetEthChainConfig().IsCancun(header.Number, header.Time) {
		header.ExcessBlobGas = new(uint64)
		header.BlobGasUsed = new(uint64)
		header.ParentBeaconRoot = new(common.Hash)
	}
	if evmtypes.GetEthChainConfig().IsPrague(header.Number, header.Time) {
		header.RequestsHash = &ethtypes.EmptyRequestsHash
	}
	return header
}

// NewTransactionFromMsg returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func NewTransactionFromMsg(
	msg *evmtypes.MsgEthereumTx,
	blockHash common.Hash,
	blockNumber, blockTime, index uint64,
	baseFee *big.Int,
	config *ethparams.ChainConfig,
	txAdditional *TxResultAdditionalFields,
) (*RPCTransaction, error) {
	if txAdditional != nil {
		return NewRPCTransactionFromIncompleteMsg(msg, blockHash, blockNumber, index, baseFee, config.ChainID, txAdditional)
	}
	return NewRPCTransaction(msg.AsTransaction(), blockHash, blockNumber, blockTime, index, baseFee, config), nil
}

// NewTransactionFromData returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
//
// This method refers to internal package method of go-ethereum v1.16.3 - newRPCTransaction
// (https://github.com/ethereum/go-ethereum/blob/d818a9af7bd5919808df78f31580f59382c53150/internal/ethapi/api.go#L991-L1081)
func NewRPCTransaction(
	tx *ethtypes.Transaction,
	blockHash common.Hash,
	blockNumber,
	blockTime uint64,
	index uint64,
	baseFee *big.Int,
	config *ethparams.ChainConfig,
) *RPCTransaction {
	// Determine the signer. For replay-protected transactions, use the most permissive
	// signer, because we assume that signers are backwards-compatible with old
	// transactions. For non-protected transactions, the frontier signer is used
	// because the latest signer will reject the unprotected transactions.
	signer := ethtypes.MakeSigner(config, new(big.Int).SetUint64(blockNumber), blockTime)
	from, _ := ethtypes.Sender(signer, tx)
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
	}
	if blockHash != (common.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}
	switch tx.Type() {
	case ethtypes.LegacyTxType:
		// if a legacy transaction has an EIP-155 chain id, include it explicitly
		if id := tx.ChainId(); id.Sign() != 0 {
			result.ChainID = (*hexutil.Big)(id)
		}

	case ethtypes.AccessListTxType:
		al := tx.AccessList()
		yparity := hexutil.Uint64(v.Sign()) //nolint:gosec // G115
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.YParity = &yparity

	case ethtypes.DynamicFeeTxType:
		al := tx.AccessList()
		yparity := hexutil.Uint64(v.Sign()) //nolint:gosec // G115
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.YParity = &yparity
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		// if the transaction has been mined, compute the effective gas price
		if baseFee != nil && blockHash != (common.Hash{}) {
			// price = min(tip, gasFeeCap - baseFee) + baseFee
			price := new(big.Int).Add(tx.GasTipCap(), baseFee)
			if price.Cmp(tx.GasFeeCap()) > 0 {
				price = tx.GasFeeCap()
			}
			result.GasPrice = (*hexutil.Big)(price)
		} else {
			result.GasPrice = (*hexutil.Big)(tx.GasFeeCap())
		}

	case ethtypes.BlobTxType:
		al := tx.AccessList()
		yparity := hexutil.Uint64(v.Sign()) //nolint:gosec
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.YParity = &yparity
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		// if the transaction has been mined, compute the effective gas price
		if baseFee != nil && blockHash != (common.Hash{}) {
			result.GasPrice = (*hexutil.Big)(effectiveGasPrice(tx, baseFee))
		} else {
			result.GasPrice = (*hexutil.Big)(tx.GasFeeCap())
		}
		result.MaxFeePerBlobGas = (*hexutil.Big)(tx.BlobGasFeeCap())
		result.BlobVersionedHashes = tx.BlobHashes()

	case ethtypes.SetCodeTxType:
		al := tx.AccessList()
		yparity := hexutil.Uint64(v.Sign()) //nolint:gosec
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.YParity = &yparity
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		// if the transaction has been mined, compute the effective gas price
		if baseFee != nil && blockHash != (common.Hash{}) {
			result.GasPrice = (*hexutil.Big)(effectiveGasPrice(tx, baseFee))
		} else {
			result.GasPrice = (*hexutil.Big)(tx.GasFeeCap())
		}
		result.AuthorizationList = tx.SetCodeAuthorizations()
	}

	return result
}

// NewRPCTransactionFromIncompleteMsg returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func NewRPCTransactionFromIncompleteMsg(
	msg *evmtypes.MsgEthereumTx, blockHash common.Hash, blockNumber, index uint64, baseFee *big.Int,
	chainID *big.Int, txAdditional *TxResultAdditionalFields,
) (*RPCTransaction, error) {
	tx := msg.AsTransaction()

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
		Hash:     tx.Hash(), // TODO: evm, check if this is correct, and what to do with added fields
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

// NewRPCPendingTransaction returns a pending transaction that will serialize to the RPC representation
func NewRPCPendingTransaction(tx *ethtypes.Transaction, current *ethtypes.Header, config *ethparams.ChainConfig) *RPCTransaction {
	var (
		baseFee     *big.Int
		blockNumber = uint64(0)
		blockTime   = uint64(0)
	)
	if current != nil {
		baseFee = eip1559.CalcBaseFee(config, current)
		blockNumber = current.Number.Uint64()
		blockTime = current.Time
	}
	return NewRPCTransaction(tx, common.Hash{}, blockNumber, blockTime, 0, baseFee, config)
}

// effectiveGasPrice computes the transaction gas fee, based on the given basefee value.
//
//	price = min(gasTipCap + baseFee, gasFeeCap)
func effectiveGasPrice(tx *ethtypes.Transaction, baseFee *big.Int) *big.Int {
	fee := tx.GasTipCap()
	fee = fee.Add(fee, baseFee)
	if tx.GasFeeCapIntCmp(fee) < 0 {
		return tx.GasFeeCap()
	}
	return fee
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

// CheckTxFee is an internal function used to check whether the fee of
// the given transaction is _reasonable_(under the minimum cap).
func CheckTxFee(gasPrice *big.Int, gas uint64, minCap float64) error {
	// Short circuit if there is no cap for transaction fee at all.
	if minCap == 0 {
		return nil
	}
	totalfee := new(big.Float).SetInt(new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gas)))
	// 1 token in atto units (1e18)
	oneToken := new(big.Float).SetInt(big.NewInt(ethparams.Ether))
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

// CalcBaseFee calculates the basefee of the header.
func CalcBaseFee(config *ethparams.ChainConfig, parent *ethtypes.Header, p feemarkettypes.Params) (*big.Int, error) {
	// If the current block is the first EIP-1559 block, return the InitialBaseFee.
	if !config.IsLondon(parent.Number) {
		return new(big.Int).SetUint64(ethparams.InitialBaseFee), nil
	}
	if p.ElasticityMultiplier == 0 {
		return nil, errors.New("ElasticityMultiplier cannot be 0 as it's checked in the params validation")
	}
	parentGasTarget := parent.GasLimit / uint64(p.ElasticityMultiplier)

	factor := evmtypes.GetEVMCoinDecimals().ConversionFactor()
	minGasPrice := p.MinGasPrice.Mul(sdkmath.LegacyNewDecFromInt(factor))
	return feemarkettypes.CalcGasBaseFee(
		parent.GasUsed, parentGasTarget, uint64(p.BaseFeeChangeDenominator),
		sdkmath.LegacyNewDecFromBigInt(parent.BaseFee), sdkmath.LegacyOneDec(), minGasPrice,
	).TruncateInt().BigInt(), nil
}

// RPCMarshalHeader converts the given header to the RPC output .
//
// This method refers to internal package method of go-ethereum v1.16.3 - RPCMarshalHeader
// (https://github.com/ethereum/go-ethereum/blob/d818a9af7bd5919808df78f31580f59382c53150/internal/ethapi/api.go#L888-L927)
// but it uses the cometbft Header to get the block hash.
func RPCMarshalHeader(head *ethtypes.Header, blockHash []byte) map[string]interface{} {
	result := map[string]interface{}{
		"number":           (*hexutil.Big)(head.Number),
		"hash":             hexutil.Bytes(blockHash), // use cometbft header hash
		"parentHash":       head.ParentHash,
		"nonce":            head.Nonce,
		"mixHash":          head.MixDigest,
		"sha3Uncles":       head.UncleHash,
		"logsBloom":        head.Bloom,
		"stateRoot":        head.Root,
		"miner":            head.Coinbase,
		"difficulty":       (*hexutil.Big)(head.Difficulty),
		"extraData":        hexutil.Bytes(head.Extra),
		"gasLimit":         hexutil.Uint64(head.GasLimit),
		"gasUsed":          (*hexutil.Big)(big.NewInt(int64(head.GasUsed))), //nolint:gosec // G115
		"timestamp":        hexutil.Uint64(head.Time),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}
	if head.BaseFee != nil {
		result["baseFeePerGas"] = (*hexutil.Big)(head.BaseFee)
	}
	if head.WithdrawalsHash != nil {
		result["withdrawalsRoot"] = head.WithdrawalsHash
	}
	if head.BlobGasUsed != nil {
		result["blobGasUsed"] = hexutil.Uint64(*head.BlobGasUsed)
	}
	if head.ExcessBlobGas != nil {
		result["excessBlobGas"] = hexutil.Uint64(*head.ExcessBlobGas)
	}
	if head.ParentBeaconRoot != nil {
		result["parentBeaconBlockRoot"] = head.ParentBeaconRoot
	}
	if head.RequestsHash != nil {
		result["requestsHash"] = head.RequestsHash
	}
	return result
}

// RPCMarshalBlock converts the given block to the RPC output which depends on fullTx. If inclTx is true transactions are
// returned. When fullTx is true the returned block contains full transaction details, otherwise it will only contain
// transaction hashes.
//
// This method refers to go-ethereum v1.16.3 internal package method - RPCMarshalBlock
// (https://github.com/ethereum/go-ethereum/blob/d818a9af7bd5919808df78f31580f59382c53150/internal/ethapi/api.go#L929-L962)
func RPCMarshalBlock(block *ethtypes.Block, cmtBlock *cmtrpccore.ResultBlock, msgs []*evmtypes.MsgEthereumTx, inclTx bool, fullTx bool, config *ethparams.ChainConfig) (map[string]interface{}, error) {
	blockHash := cmtBlock.BlockID.Hash.Bytes()
	fields := RPCMarshalHeader(block.Header(), blockHash)
	fields["size"] = hexutil.Uint64(block.Size())

	if inclTx {
		formatTx := func(idx int, tx *ethtypes.Transaction) interface{} {
			return tx.Hash()
		}
		if fullTx {
			formatTx = func(idx int, _ *ethtypes.Transaction) interface{} {
				txIdx := uint64(idx) //nolint:gosec // G115
				return newRPCTransactionFromBlockIndex(block, common.BytesToHash(blockHash), txIdx, config)
			}
		}
		txs := block.Transactions()
		transactions := make([]interface{}, len(txs))
		for i, tx := range txs {
			transactions[i] = formatTx(i, tx)
		}
		fields["transactions"] = transactions
	}
	uncles := block.Uncles()
	uncleHashes := make([]common.Hash, len(uncles))
	for i, uncle := range uncles {
		uncleHashes[i] = uncle.Hash()
	}
	fields["uncles"] = uncleHashes
	if block.Withdrawals() != nil {
		fields["withdrawals"] = block.Withdrawals()
	}
	return fields, nil
}

// newRPCTransactionFromBlockIndex returns a transaction that will serialize to the RPC representation.
func newRPCTransactionFromBlockIndex(b *ethtypes.Block, blockHash common.Hash, index uint64, config *ethparams.ChainConfig) *RPCTransaction {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	return NewRPCTransaction(txs[index], blockHash, b.NumberU64(), b.Time(), index, b.BaseFee(), config)
}

// RPCMarshalReceipt marshals a transaction receipt into a JSON object.
//
// This method refers to go-ethereum v1.16.3 internal package method marshalReceipt
// (https://github.com/ethereum/go-ethereum/blob/d818a9af7bd5919808df78f31580f59382c53150/internal/ethapi/api.go#L1478-L1518)
func RPCMarshalReceipt(receipt *ethtypes.Receipt, tx *ethtypes.Transaction, from common.Address) (map[string]interface{}, error) {
	fields := map[string]interface{}{
		"blockHash":         receipt.BlockHash,
		"blockNumber":       hexutil.Uint64(receipt.BlockNumber.Uint64()),
		"transactionHash":   tx.Hash(),
		"transactionIndex":  hexutil.Uint64(receipt.TransactionIndex),
		"from":              from,
		"to":                tx.To(),
		"gasUsed":           hexutil.Uint64(receipt.GasUsed),
		"cumulativeGasUsed": hexutil.Uint64(receipt.CumulativeGasUsed),
		"contractAddress":   nil,
		"logs":              receipt.Logs,
		"logsBloom":         receipt.Bloom,
		"type":              hexutil.Uint(tx.Type()),
		"effectiveGasPrice": (*hexutil.Big)(receipt.EffectiveGasPrice),
	}

	// Assign receipt status or post state.
	if len(receipt.PostState) > 0 {
		fields["root"] = hexutil.Bytes(receipt.PostState)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}
	if receipt.Logs == nil {
		fields["logs"] = []*ethtypes.Log{}
	}

	if tx.Type() == ethtypes.BlobTxType {
		fields["blobGasUsed"] = hexutil.Uint64(receipt.BlobGasUsed)
		fields["blobGasPrice"] = (*hexutil.Big)(receipt.BlobGasPrice)
	}

	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if receipt.ContractAddress != (common.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress
	}
	return fields, nil
}

// EffectiveGasPrice computes the transaction gas fee, based on the given basefee value.
//
// price = min(gasTipCap + baseFee, gasFeeCap)
//
// This method refers to go-ethereum v1.16.3 internal package method, effectiveGasPrice.
// (https://github.com/ethereum/go-ethereum/blob/d818a9af7bd5919808df78f31580f59382c53150/internal/ethapi/api.go#L1083-L1093)
func EffectiveGasPrice(tx *ethtypes.Transaction, baseFee *big.Int) *big.Int {
	if tx == nil {
		return big.NewInt(0)
	}
	if baseFee == nil {
		return tx.GasFeeCap()
	}

	fee := tx.GasTipCap()
	fee = fee.Add(fee, baseFee)
	if tx.GasFeeCapIntCmp(fee) < 0 {
		return tx.GasFeeCap()
	}

	return fee
}
