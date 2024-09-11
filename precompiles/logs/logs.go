package logs

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

// AddLog adds log to stateDB
func AddLog(ctx sdk.Context, precompileAddr common.Address, stateDB vm.StateDB, topics []common.Hash, data []byte) {
	stateDB.AddLog(&types.Log{
		Address:     precompileAddr,
		Topics:      topics,
		Data:        data,
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

// PackBigInt is a helper function to pack a uint256 amount
func PackBigInt(amount *big.Int) ([]byte, error) {
	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, err
	}

	arguments := abi.Arguments{
		{
			Type: uint256Type,
		},
	}

	return arguments.Pack(amount)
}
