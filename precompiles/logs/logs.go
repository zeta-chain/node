package logs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

type Argument struct {
	Type  string
	Value interface{}
}

// AddLog adds log to stateDB
func AddLog(ctx sdk.Context, precompileAddr common.Address, stateDB vm.StateDB, topics []common.Hash, data []byte) {
	stateDB.AddLog(&types.Log{
		Address: precompileAddr,
		Topics:  topics,
		Data:    data,
		// #nosec G115 block height always positive
		BlockNumber: uint64(ctx.BlockHeight()),
	})
}

// MakeTopics creates topics for log as wrapper around geth abi.MakeTopics function
func MakeTopics(event abi.Event, query ...[]interface{}) ([]common.Hash, error) {
	topics := []common.Hash{event.ID}

	topicsRes, err := abi.MakeTopics(
		query...,
	)
	if err != nil {
		return nil, err
	}

	for _, topic := range topicsRes {
		topics = append(topics, topic[0])
	}

	return topics, nil
}

// PackArguments packs an arbitrary number of logs.Arguments as non-indexed data.
// When packing data, make sure the Argument are passed in the same order as the event definition.
func PackArguments(args []Argument) ([]byte, error) {
	types := abi.Arguments{}
	toPack := []interface{}{}

	for _, arg := range args {
		abiType, err := abi.NewType(arg.Type, "", nil)
		if err != nil {
			return nil, err
		}

		types = append(types, abi.Argument{
			Type: abiType,
		})

		toPack = append(toPack, arg.Value)
	}

	data, err := types.Pack(toPack...)
	if err != nil {
		return nil, err
	}

	return data, nil
}
