// Copyright 2021 Evmos Foundation
// This file is part of Evmos' Ethermint library.
//
// The Ethermint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Ethermint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Ethermint library. If not, see https://github.com/evmos/ethermint/blob/main/LICENSE
package types

import (
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	abci "github.com/cometbft/cometbft/abci/types"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// EventFormat is the format version of the events.
//
// To fix the issue of tx exceeds block gas limit, we changed the event format in a breaking way.
// But to avoid forcing clients to re-sync from scatch, we make json-rpc logic to be compatible with both formats.
type EventFormat int

const (
	MessageType                    = "message"
	AmountType                     = "amount"
	SenderType                     = "sender"
	CosmosEVMTxType                = 88
	eventFormatUnknown EventFormat = iota

	// Event Format 1 (the format used before PR #1062):
	// ```
	// ethereum_tx(amount, ethereumTxHash, [txIndex, txGasUsed], txHash, [receipient], ethereumTxFailed)
	// tx_log(txLog, txLog, ...)
	// ethereum_tx(amount, ethereumTxHash, [txIndex, txGasUsed], txHash, [receipient], ethereumTxFailed)
	// tx_log(txLog, txLog, ...)
	// ...
	// ```
	eventFormat1

	// Event Format 2 (the format used after PR #1062):
	// ```
	// ethereum_tx(ethereumTxHash, txIndex)
	// ethereum_tx(ethereumTxHash, txIndex)
	// ...
	// ethereum_tx(amount, ethereumTxHash, txIndex, txGasUsed, txHash, [receipient], ethereumTxFailed)
	// tx_log(txLog, txLog, ...)
	// ethereum_tx(amount, ethereumTxHash, txIndex, txGasUsed, txHash, [receipient], ethereumTxFailed)
	// tx_log(txLog, txLog, ...)
	// ...
	// ```
	// If the transaction exceeds block gas limit, it only emits the first part.
	eventFormat2
)

// ParsedTx is the tx infos parsed from events.
type ParsedTx struct {
	// max uint32 means there is no sdk.Msg corresponding to eth tx
	MsgIndex int

	// the following fields are parsed from events

	Hash common.Hash
	// -1 means uninitialized
	EthTxIndex int32
	GasUsed    uint64
	Failed     bool
	// Additional cosmos EVM tx fields
	TxHash    string
	Type      uint64
	Amount    *big.Int
	Recipient common.Address
	Sender    common.Address
	Nonce     uint64
	Data      []byte
}

// NewParsedTx initialize a ParsedTx
func NewParsedTx(msgIndex int) ParsedTx {
	return ParsedTx{MsgIndex: msgIndex, EthTxIndex: -1}
}

// ParsedTxs is the tx infos parsed from eth tx events.
type ParsedTxs struct {
	// one item per message
	Txs []ParsedTx
	// map tx hash to msg index
	TxHashes map[common.Hash]int
}

// ParseTxResult parse eth tx infos from cosmos-sdk events.
// It supports two event formats, the formats are described in the comments of the format constants.
func ParseTxResult(result *abci.ResponseDeliverTx, tx sdk.Tx) (*ParsedTxs, error) {
	format := eventFormatUnknown
	// the index of current ethereum_tx event in format 1 or the second part of format 2
	eventIndex := -1

	p := &ParsedTxs{
		TxHashes: make(map[common.Hash]int),
	}
	prevEventType := ""
	for _, event := range result.Events {
		if event.Type != evmtypes.EventTypeEthereumTx &&
			(prevEventType != evmtypes.EventTypeEthereumTx || event.Type != MessageType) {
			continue
		}

		// Parse tendermint message after ethereum_tx event
		if prevEventType == evmtypes.EventTypeEthereumTx && event.Type == MessageType && eventIndex != -1 {
			err := fillTxAttributes(&p.Txs[eventIndex], event.Attributes)
			if err != nil {
				return nil, err
			}
		}

		if event.Type == MessageType {
			prevEventType = MessageType
			continue
		}

		if format == eventFormatUnknown {
			// discover the format version by inspect the first ethereum_tx event.
			if len(event.Attributes) > 2 {
				format = eventFormat1
			} else {
				format = eventFormat2
			}
		}

		if len(event.Attributes) == 2 {
			// the first part of format 2
			if err := p.newTx(event.Attributes); err != nil {
				return nil, err
			}
		} else {
			// format 1 or second part of format 2
			eventIndex++
			if format == eventFormat1 {
				// append tx
				if err := p.newTx(event.Attributes); err != nil {
					return nil, err
				}
			} else {
				// the second part of format 2, update tx fields
				if err := p.updateTx(eventIndex, event.Attributes); err != nil {
					return nil, err
				}
			}
		}

		prevEventType = evmtypes.EventTypeEthereumTx
	}

	// some old versions miss some events, fill it with tx result
	// txs with type CosmosEVMTxType will always emit GasUsed in events so no need to override for those
	if len(p.Txs) == 1 && p.Txs[0].Type != CosmosEVMTxType {
		// #nosec G115 always positive
		p.Txs[0].GasUsed = uint64(result.GasUsed)
	}

	// this could only happen if tx exceeds block gas limit
	if result.Code != 0 && tx != nil {
		for i := 0; i < len(p.Txs); i++ {
			p.Txs[i].Failed = true

			// replace gasUsed with gasLimit because that's what's actually deducted.
			gasLimit := tx.GetMsgs()[i].(*evmtypes.MsgEthereumTx).GetGas()
			p.Txs[i].GasUsed = gasLimit
		}
	}

	// fix msg indexes, because some eth txs indexed here don't have corresponding sdk.Msg
	currMsgIndex := 0
	for _, tx := range p.Txs {
		if tx.Type == CosmosEVMTxType {
			tx.MsgIndex = math.MaxUint32
			// todo: fix mapping as well
		} else {
			tx.MsgIndex = currMsgIndex
			currMsgIndex++
		}
	}
	return p, nil
}

// ParseTxIndexerResult parse tm tx result to a format compatible with the custom tx indexer.
func ParseTxIndexerResult(
	txResult *tmrpctypes.ResultTx,
	tx sdk.Tx,
	getter func(*ParsedTxs) *ParsedTx,
) (*ethermint.TxResult, *TxResultAdditionalFields, error) {
	txs, err := ParseTxResult(&txResult.TxResult, tx)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to parse tx events: block %d, index %d, %v",
			txResult.Height,
			txResult.Index,
			err,
		)
	}

	parsedTx := getter(txs)
	if parsedTx == nil {
		return nil, nil, fmt.Errorf(
			"ethereum tx not found in msgs: block %d, index %d",
			txResult.Height,
			txResult.Index,
		)
	}
	if parsedTx.Type == CosmosEVMTxType {
		return &ethermint.TxResult{
				Height:  txResult.Height,
				TxIndex: txResult.Index,
				// #nosec G115 always in range
				MsgIndex:          uint32(parsedTx.MsgIndex),
				EthTxIndex:        parsedTx.EthTxIndex,
				Failed:            parsedTx.Failed,
				GasUsed:           parsedTx.GasUsed,
				CumulativeGasUsed: txs.AccumulativeGasUsed(parsedTx.MsgIndex),
			}, &TxResultAdditionalFields{
				Value:     parsedTx.Amount,
				Hash:      parsedTx.Hash,
				TxHash:    parsedTx.TxHash,
				Type:      parsedTx.Type,
				Recipient: parsedTx.Recipient,
				Sender:    parsedTx.Sender,
				GasUsed:   parsedTx.GasUsed,
				Data:      parsedTx.Data,
				Nonce:     parsedTx.Nonce,
			}, nil
	}
	return &ethermint.TxResult{
		Height:  txResult.Height,
		TxIndex: txResult.Index,
		// #nosec G115 always in range
		MsgIndex:          uint32(parsedTx.MsgIndex),
		EthTxIndex:        parsedTx.EthTxIndex,
		Failed:            parsedTx.Failed,
		GasUsed:           parsedTx.GasUsed,
		CumulativeGasUsed: txs.AccumulativeGasUsed(parsedTx.MsgIndex),
	}, nil, nil
}

// ParseTxIndexerResult parse tm tx result to a format compatible with the custom tx indexer.
func ParseTxBlockResult(
	txResult *abci.ResponseDeliverTx,
	tx sdk.Tx,
	txIndex int,
	height int64,
) (*ethermint.TxResult, *TxResultAdditionalFields, error) {
	txs, err := ParseTxResult(txResult, tx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse tx events: block %d, index %d, %v", height, txIndex, err)
	}

	if len(txs.Txs) == 0 {
		return nil, nil, fmt.Errorf("ethereum tx not found in msgs: block %d, index %d", height, txIndex)
	}
	parsedTx := txs.Txs[0]
	if parsedTx.Type == CosmosEVMTxType {
		return &ethermint.TxResult{
				Height: height,
				// #nosec G115 always in range
				TxIndex: uint32(txIndex),
				// #nosec G115 always in range
				MsgIndex:          uint32(parsedTx.MsgIndex),
				EthTxIndex:        parsedTx.EthTxIndex,
				Failed:            parsedTx.Failed,
				GasUsed:           parsedTx.GasUsed,
				CumulativeGasUsed: txs.AccumulativeGasUsed(parsedTx.MsgIndex),
			}, &TxResultAdditionalFields{
				Value:     parsedTx.Amount,
				Hash:      parsedTx.Hash,
				TxHash:    parsedTx.TxHash,
				Type:      parsedTx.Type,
				Recipient: parsedTx.Recipient,
				Sender:    parsedTx.Sender,
				GasUsed:   parsedTx.GasUsed,
				Data:      parsedTx.Data,
				Nonce:     parsedTx.Nonce,
			}, nil
	}
	return &ethermint.TxResult{
		Height: height,
		// #nosec G115 always in range
		TxIndex: uint32(txIndex),
		// #nosec G115 always in range
		MsgIndex:          uint32(parsedTx.MsgIndex),
		EthTxIndex:        parsedTx.EthTxIndex,
		Failed:            parsedTx.Failed,
		GasUsed:           parsedTx.GasUsed,
		CumulativeGasUsed: txs.AccumulativeGasUsed(parsedTx.MsgIndex),
	}, nil, nil
}

// newTx parse a new tx from events, called during parsing.
func (p *ParsedTxs) newTx(attrs []abci.EventAttribute) error {
	msgIndex := len(p.Txs)
	tx := NewParsedTx(msgIndex)
	if err := fillTxAttributes(&tx, attrs); err != nil {
		return err
	}
	p.Txs = append(p.Txs, tx)
	p.TxHashes[tx.Hash] = msgIndex
	return nil
}

// updateTx updates an exiting tx from events, called during parsing.
// In event format 2, we update the tx with the attributes of the second `ethereum_tx` event,
// Due to bug https://github.com/evmos/ethermint/issues/1175, the first `ethereum_tx` event may emit incorrect tx hash,
// so we prefer the second event and override the first one.
func (p *ParsedTxs) updateTx(eventIndex int, attrs []abci.EventAttribute) error {
	tx := NewParsedTx(eventIndex)
	if err := fillTxAttributes(&tx, attrs); err != nil {
		return err
	}
	if tx.Hash != p.Txs[eventIndex].Hash {
		// if hash is different, index the new one too
		p.TxHashes[tx.Hash] = eventIndex
	}
	// override the tx because the second event is more trustworthy
	p.Txs[eventIndex] = tx
	return nil
}

// GetTxByHash find ParsedTx by tx hash, returns nil if not exists.
func (p *ParsedTxs) GetTxByHash(hash common.Hash) *ParsedTx {
	if idx, ok := p.TxHashes[hash]; ok {
		return &p.Txs[idx]
	}
	return nil
}

// GetTxByMsgIndex returns ParsedTx by msg index
func (p *ParsedTxs) GetTxByMsgIndex(i int) *ParsedTx {
	if i < 0 || i >= len(p.Txs) {
		return nil
	}
	return &p.Txs[i]
}

// GetTxByTxIndex returns ParsedTx by tx index
func (p *ParsedTxs) GetTxByTxIndex(txIndex int) *ParsedTx {
	if len(p.Txs) == 0 {
		return nil
	}
	// assuming the `EthTxIndex` increase continuously,
	// convert TxIndex to MsgIndex by subtract the begin TxIndex.
	msgIndex := txIndex - int(p.Txs[0].EthTxIndex)
	// GetTxByMsgIndex will check the bound
	return p.GetTxByMsgIndex(msgIndex)
}

// AccumulativeGasUsed calculates the accumulated gas used within the batch of txs
func (p *ParsedTxs) AccumulativeGasUsed(msgIndex int) (result uint64) {
	for i := 0; i <= msgIndex; i++ {
		result += p.Txs[i].GasUsed
	}
	return result
}

// fillTxAttribute parse attributes by name, less efficient than hardcode the index, but more stable against event
// format changes.
func fillTxAttribute(tx *ParsedTx, key, value string) error {
	switch key {
	case evmtypes.AttributeKeyEthereumTxHash:
		tx.Hash = common.HexToHash(value)
	case evmtypes.AttributeKeyTxIndex:
		txIndex, err := strconv.ParseUint(value, 10, 31)
		if err != nil {
			return err
		}
		// #nosec G115 always in range
		tx.EthTxIndex = int32(txIndex)
	case evmtypes.AttributeKeyTxGasUsed:
		gasUsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		tx.GasUsed = gasUsed
	case evmtypes.AttributeKeyEthereumTxFailed:
		tx.Failed = len(value) > 0
	case SenderType:
		tx.Sender = common.HexToAddress(value)
	case evmtypes.AttributeKeyRecipient:
		tx.Recipient = common.HexToAddress(value)
	case evmtypes.AttributeKeyTxHash:
		tx.TxHash = value
	case evmtypes.AttributeKeyTxType:
		txType, err := strconv.ParseUint(value, 10, 31)
		if err != nil {
			return err
		}
		tx.Type = txType
	case AmountType:
		var success bool
		tx.Amount, success = big.NewInt(0).SetString(value, 10)
		if !success {
			return nil
		}
	case evmtypes.AttributeKeyTxNonce:
		nonce, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		tx.Nonce = nonce

	case evmtypes.AttributeKeyTxData:
		hexBytes, err := hexutil.Decode(value)
		if err != nil {
			return err
		}
		tx.Data = hexBytes
	}
	return nil
}

func fillTxAttributes(tx *ParsedTx, attrs []abci.EventAttribute) error {
	// before cosmos upgrade to 0.47, attributes are base64 encoded
	// purpose of this is to support older txs as well
	isLegacyAttrs := isLegacyAttrEncoding(attrs)
	for _, attr := range attrs {
		if isLegacyAttrs {
			// only decode if value can be decoded
			// (error should not happen because at this point it is determined it is legacy attr)
			decKey, err := base64.StdEncoding.DecodeString(attr.Key)
			if err != nil {
				return err
			}
			attr.Key = string(decKey)
			decValue, err := base64.StdEncoding.DecodeString(attr.Value)
			if err != nil {
				return err
			}
			attr.Value = string(decValue)
		}

		if err := fillTxAttribute(tx, attr.Key, attr.Value); err != nil {
			return err
		}
	}
	return nil
}

func isLegacyAttrEncoding(attrs []abci.EventAttribute) bool {
	for _, attr := range attrs {
		if strings.Contains(attr.Key, "==") || strings.Contains(attr.Value, "==") {
			return true
		}
	}

	return false
}
