package backend

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/pkg/errors"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"

	evmtypes "github.com/cosmos/evm/x/vm/types"
	rpctypes "github.com/zeta-chain/node/rpc/types"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RPCBlockFromCometBlock returns a JSON-RPC compatible Ethereum block from a
// given CometBFT block and its block result.
func (b *Backend) RPCHeaderFromCometBlock(
	resBlock *cmtrpctypes.ResultBlock,
	blockRes *cmtrpctypes.ResultBlockResults,
) (map[string]interface{}, error) {
	ethBlock, err := b.EthBlockFromCometBlock(resBlock, blockRes)
	if err != nil {
		return nil, fmt.Errorf("failed to get rpc block from comet block: %w", err)
	}

	return rpctypes.RPCMarshalHeader(ethBlock.Header(), resBlock.BlockID.Hash), nil
}

// RPCBlockFromCometBlock returns a JSON-RPC compatible Ethereum block from a
// given CometBFT block and its block result.
func (b *Backend) RPCBlockFromCometBlock( // TODO evm: refactor to include synthetic txs...
	resBlock *cmtrpctypes.ResultBlock,
	blockRes *cmtrpctypes.ResultBlockResults,
	fullTx bool,
) (map[string]interface{}, error) {
	msgs, _ := b.EthMsgsFromCometBlock(resBlock, blockRes)
	ethBlock, err := b.EthBlockFromCometBlock(resBlock, blockRes)
	if err != nil {
		return nil, fmt.Errorf("failed to get rpc block from comet block: %w", err)
	}

	return rpctypes.RPCMarshalBlock(ethBlock, resBlock, msgs, true, fullTx, b.ChainConfig())
}

// BlockNumberFromComet returns the BlockNumber from BlockNumberOrHash
func (b *Backend) BlockNumberFromComet(blockNrOrHash rpctypes.BlockNumberOrHash) (rpctypes.BlockNumber, error) {
	switch {
	case blockNrOrHash.BlockHash == nil && blockNrOrHash.BlockNumber == nil:
		return rpctypes.EthEarliestBlockNumber, fmt.Errorf("types BlockHash and BlockNumber cannot be both nil")
	case blockNrOrHash.BlockHash != nil:
		blockNumber, err := b.BlockNumberFromCometByHash(*blockNrOrHash.BlockHash)
		if err != nil {
			return rpctypes.EthEarliestBlockNumber, err
		}
		return rpctypes.NewBlockNumber(blockNumber), nil
	case blockNrOrHash.BlockNumber != nil:
		return *blockNrOrHash.BlockNumber, nil
	default:
		return rpctypes.EthEarliestBlockNumber, nil
	}
}

// BlockNumberFromCometByHash returns the block height of given block hash
func (b *Backend) BlockNumberFromCometByHash(blockHash common.Hash) (*big.Int, error) {
	resHeader, err := b.RPCClient.HeaderByHash(b.Ctx, blockHash.Bytes())
	if err != nil {
		return nil, err
	}

	if resHeader == nil || resHeader.Header == nil {
		return nil, errors.Errorf("header not found for hash %s", blockHash.Hex())
	}

	return big.NewInt(resHeader.Header.Height), nil
}

// RPCBlockFromCometBlock returns a JSON-RPC compatible Ethereum block from a
// given CometBFT block and its block result.
func (b *Backend) EthBlockFromCometBlock(
	resBlock *cmtrpctypes.ResultBlock,
	blockRes *cmtrpctypes.ResultBlockResults,
) (*ethtypes.Block, error) {
	cmtBlock := resBlock.Block

	// 1. get base fee
	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		// handle the error for pruned node.
		b.Logger.Error("failed to fetch Base Fee from prunned block. Check node prunning configuration", "height", cmtBlock.Height, "error", err)
	}

	// 2. get miner
	miner, err := b.MinerFromCometBlock(resBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to get miner(block proposer) address from comet block")
	}

	// 3. get block gasLimit
	ctx := rpctypes.ContextWithHeight(cmtBlock.Height)
	gasLimit, err := rpctypes.BlockMaxGasFromConsensusParams(ctx, b.ClientCtx, cmtBlock.Height)
	if err != nil {
		b.Logger.Error("failed to query consensus params", "error", err.Error())
	}

	// 4. create blockHeader without transactions, receipts, withdrawals, ...
	ethHeader := rpctypes.MakeHeader(cmtBlock.Header, gasLimit, miner, baseFee)

	// 5. get MsgEthereumTxs

	msgs, additionals := b.EthMsgsFromCometBlock(resBlock, blockRes)
	txs := make([]*ethtypes.Transaction, len(msgs))
	for i, ethMsg := range msgs {
		if additionals[i] == nil {
			txs = append(txs, ethMsg.AsTransaction())
		}
	}

	// 6. create ethBlock body with transactions
	body := &ethtypes.Body{
		Transactions: txs,
		Uncles:       []*ethtypes.Header{},
		Withdrawals:  []*ethtypes.Withdrawal{},
	}

	// 7. receipts
	// receipts, err := b.ReceiptsFromCometBlock(resBlock, blockRes, msgs)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get receipts from comet block: %w", err)
	// }

	// 8. Gas Used
	gasUsed := uint64(0)
	for _, txsResult := range blockRes.TxsResults {
		// workaround for cosmos-sdk bug. https://github.com/cosmos/cosmos-sdk/issues/10832
		if ShouldIgnoreGasUsed(txsResult) {
			// block gas limit has exceeded, other txs must have failed with same reason.
			break
		}
		gasUsed += uint64(txsResult.GetGasUsed()) // #nosec G115 -- checked for int overflow already
	}
	ethHeader.GasUsed = gasUsed

	// 9. create eth block
	ethBlock := ethtypes.NewBlock(ethHeader, body, nil, trie.NewStackTrie(nil))
	return ethBlock, nil
}

func (b *Backend) MinerFromCometBlock(
	resBlock *cmtrpctypes.ResultBlock,
) (common.Address, error) {
	cmtBlock := resBlock.Block

	req := &evmtypes.QueryValidatorAccountRequest{
		ConsAddress: sdk.ConsAddress(cmtBlock.Header.ProposerAddress).String(),
	}

	var validatorAccAddr sdk.AccAddress

	ctx := rpctypes.ContextWithHeight(cmtBlock.Height)
	res, err := b.QueryClient.ValidatorAccount(ctx, req)
	if err != nil {
		b.Logger.Debug(
			"failed to query validator operator address",
			"height", cmtBlock.Height,
			"cons-address", req.ConsAddress,
			"error", err.Error(),
		)
		// use zero address as the validator operator address
		validatorAccAddr = sdk.AccAddress(common.Address{}.Bytes())
	} else {
		validatorAccAddr, err = sdk.AccAddressFromBech32(res.AccountAddress)
		if err != nil {
			return common.Address{}, err
		}
	}

	return common.BytesToAddress(validatorAccAddr), nil
}

func (b *Backend) ReceiptsFromCometBlock(
	resBlock *cmtrpctypes.ResultBlock,
	blockRes *cmtrpctypes.ResultBlockResults,
	msgs []*evmtypes.MsgEthereumTx,
) ([]*ethtypes.Receipt, error) {
	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		// handle the error for pruned node.
		b.Logger.Error("failed to fetch Base Fee from prunned block. Check node prunning configuration", "height", resBlock.Block.Height, "error", err)
	}

	blockHash := common.BytesToHash(resBlock.BlockID.Hash)
	receipts := make([]*ethtypes.Receipt, len(msgs))
	cumulatedGasUsed := uint64(0)
	for i, ethMsg := range msgs {
		txResult, _, err := b.GetTxByEthHash(ethMsg.Hash()) // TODO evm
		if err != nil {
			return nil, fmt.Errorf("tx not found: hash=%s, error=%s", ethMsg.Hash(), err.Error())
		}

		cumulatedGasUsed += txResult.GasUsed

		var effectiveGasPrice *big.Int
		if baseFee != nil {
			effectiveGasPrice = rpctypes.EffectiveGasPrice(ethMsg.Raw.Transaction, baseFee)
		} else {
			effectiveGasPrice = ethMsg.Raw.GasFeeCap()
		}

		var status uint64
		if txResult.Failed {
			status = ethtypes.ReceiptStatusFailed
		} else {
			status = ethtypes.ReceiptStatusSuccessful
		}

		contractAddress := common.Address{}
		if ethMsg.Raw.To() == nil {
			contractAddress = crypto.CreateAddress(ethMsg.GetSender(), ethMsg.Raw.Nonce())
		}

		msgIndex := int(txResult.MsgIndex) // #nosec G115 -- checked for int overflow already
		logs, err := evmtypes.DecodeMsgLogs(
			blockRes.TxsResults[txResult.TxIndex].Data,
			msgIndex,
			uint64(resBlock.Block.Height), // #nosec G115 -- checked for int overflow already
		)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tx result to eth receipt: %w", err)
		}

		bloom := ethtypes.CreateBloom(&ethtypes.Receipt{Logs: logs})

		receipt := &ethtypes.Receipt{
			// Consensus fields: These fields are defined by the Yellow Paper
			Type:              ethMsg.Raw.Type(),
			PostState:         nil,
			Status:            status, // convert to 1=success, 0=failure
			CumulativeGasUsed: cumulatedGasUsed,
			Bloom:             bloom,
			Logs:              logs,

			// Implementation fields: These fields are added by geth when processing a transaction.
			TxHash:            ethMsg.Hash(),
			ContractAddress:   contractAddress,
			GasUsed:           txResult.GasUsed,
			EffectiveGasPrice: effectiveGasPrice,
			BlobGasUsed:       uint64(0),     // TODO: fill this field
			BlobGasPrice:      big.NewInt(0), // TODO: fill this field

			// Inclusion information: These fields provide information about the inclusion of the
			// transaction corresponding to this receipt.
			BlockHash:        blockHash,
			BlockNumber:      big.NewInt(resBlock.Block.Height),
			TransactionIndex: uint(txResult.EthTxIndex), // #nosec G115 -- checked for int overflow already
		}

		receipts[i] = receipt
	}

	return receipts, nil
}

// EthMsgsFromCometBlock returns all real and synthetic MsgEthereumTxs from a
// Coment block. It also ensures consistency over the correct txs indexes
// across RPC endpoints
func (b *Backend) EthMsgsFromCometBlock(
	resBlock *tmrpctypes.ResultBlock,
	blockRes *tmrpctypes.ResultBlockResults,
) ([]*evmtypes.MsgEthereumTx, []*rpctypes.TxResultAdditionalFields) {
	var ethMsgs []*evmtypes.MsgEthereumTx
	var txsAdditional []*rpctypes.TxResultAdditionalFields
	block := resBlock.Block

	txResults := blockRes.TxsResults
	for i, tx := range block.Txs {
		// Check if tx exists on EVM by cross checking with blockResults:
		//  - Include unsuccessful tx that exceeds block gas limit
		//  - Include unsuccessful tx that failed when committing changes to stateDB
		//  - Exclude unsuccessful tx with any other error but ExceedBlockGasLimit
		if !rpctypes.TxSucessOrExpectedFailure(txResults[i]) {
			b.Logger.Debug("invalid tx result code", "cosmos-hash", hexutil.Encode(tx.Hash()))
			continue
		}

		tx, err := b.ClientCtx.TxConfig.TxDecoder()(tx)
		// assumption is that if regular ethermint msg is found in tx
		// there should not be synthetic one as well
		shouldCheckForSyntheticTx := true
		// if tx can be decoded, try to find MsgEthereumTx inside
		if err == nil {
			for _, msg := range tx.GetMsgs() {
				ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
				if ok {
					shouldCheckForSyntheticTx = false
					// ethMsg.Hash = ethMsg.AsTransaction().Hash().Hex() TODO: will this deserialize correctly?
					ethMsgs = append(ethMsgs, ethMsg)
					txsAdditional = append(txsAdditional, nil)
				}
			}
		} else {
			b.Logger.Debug("failed to decode transaction in block", "height", block.Height, "error", err.Error())
		}

		// if tx can not be decoded or MsgEthereumTx was not found, try to parse it from block results
		if shouldCheckForSyntheticTx {
			ethMsg, additional := b.parseSyntheticTxFromBlockResults(txResults, i, tx, block)
			if ethMsg != nil {
				ethMsgs = append(ethMsgs, ethMsg)
				txsAdditional = append(txsAdditional, additional)
			}
		}
	}
	return ethMsgs, txsAdditional
}

// BlockBloom query block bloom filter from block results
func (b *Backend) BlockBloomFromCometBlock(blockRes *cmtrpctypes.ResultBlockResults) (ethtypes.Bloom, error) {
	for _, event := range blockRes.FinalizeBlockEvents {
		if event.Type != evmtypes.EventTypeBlockBloom {
			continue
		}

		for _, attr := range event.Attributes {
			if attr.Key == evmtypes.AttributeKeyEthereumBloom {
				return ethtypes.BytesToBloom([]byte(attr.Value)), nil
			}
		}
	}
	return ethtypes.Bloom{}, errors.New("block bloom event is not found")
}

func (b *Backend) parseSyntheticTxFromBlockResults(
	txResults []*abci.ExecTxResult,
	i int,
	tx sdk.Tx,
	block *types.Block,
) (*evmtypes.MsgEthereumTx, *rpctypes.TxResultAdditionalFields) {
	res, additional, err := rpctypes.ParseTxBlockResult(txResults[i], tx, i, block.Height)
	// just skip tx if it can not be parsed, so remaining txs from the block are parsed
	if err != nil {
		b.Logger.Error(err.Error())
		return nil, nil
	}
	if additional == nil || res == nil {
		b.Logger.Debug("synthetic ethereum tx not found in msgs: block %d, index %d", block.Height, i)
		return nil, nil
	}
	return b.parseSyntethicTxFromAdditionalFields(additional), additional
}

func (b *Backend) parseSyntethicTxFromAdditionalFields(
	additional *rpctypes.TxResultAdditionalFields,
) *evmtypes.MsgEthereumTx {
	recipient := additional.Recipient
	// for transactions before v31 this value was mistakenly used for Gas field
	gas := additional.GasUsed
	if additional.GasLimit != nil {
		gas = *additional.GasLimit
	}
	t := ethtypes.NewTx(&ethtypes.LegacyTx{
		Nonce:    additional.Nonce,
		Data:     additional.Data,
		Gas:      gas,
		To:       &recipient,
		GasPrice: nil,
		Value:    additional.Value,
		V:        big.NewInt(0),
		R:        big.NewInt(0),
		S:        big.NewInt(0),
	})
	ethMsg := &evmtypes.MsgEthereumTx{}
	ethMsg.FromEthereumTx(t)

	// ethMsg.Hash = additional.Hash.Hex() // TODO evm check
	ethMsg.From = additional.Sender.Bytes()
	return ethMsg
}
