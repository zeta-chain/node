package staking

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

func (c *Contract) AddStakeLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	staker common.Address,
	validator string,
	amount *big.Int,
) error {
	event := c.Abi().Events["Stake"]
	topics := []common.Hash{event.ID}
	valAddr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return err
	}

	// staker and validator are indexed event params
	topicsRes, err := abi.MakeTopics(
		[]interface{}{staker},
		[]interface{}{common.BytesToAddress(valAddr.Bytes())},
	)
	if err != nil {
		return err
	}

	for _, topic := range topicsRes {
		topics = append(topics, topic[0])
	}

	// amount is part of event data
	data, err := packAmount(amount)
	if err != nil {
		return err
	}

	stateDB.AddLog(&types.Log{
		Address:     c.Address(),
		Topics:      topics,
		Data:        data,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

func (c *Contract) AddUnstakeLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	staker common.Address,
	validator string,
	amount *big.Int,
) error {
	event := c.Abi().Events["Unstake"]
	topics := []common.Hash{event.ID}
	valAddr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return err
	}

	// staker and validator are indexed event params
	topicsRes, err := abi.MakeTopics(
		[]interface{}{staker},
		[]interface{}{common.BytesToAddress(valAddr.Bytes())},
	)
	if err != nil {
		return err
	}

	for _, topic := range topicsRes {
		topics = append(topics, topic[0])
	}

	// amount is part of event data
	data, err := packAmount(amount)
	if err != nil {
		return err
	}

	stateDB.AddLog(&types.Log{
		Address:     c.Address(),
		Topics:      topics,
		Data:        data,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

func (c *Contract) AddMoveStakeLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	staker common.Address,
	validatorSrc string,
	validatorDst string,
	amount *big.Int,
) error {
	event := c.Abi().Events["MoveStake"]
	topics := []common.Hash{event.ID}
	validatorSrcAddr, err := sdk.ValAddressFromBech32(validatorSrc)
	if err != nil {
		return err
	}

	validatorDstAddr, err := sdk.ValAddressFromBech32(validatorDst)
	if err != nil {
		return err
	}

	// staker and validators are indexed event params
	topicsRes, err := abi.MakeTopics(
		[]interface{}{staker},
		[]interface{}{common.BytesToAddress(validatorSrcAddr.Bytes())},
		[]interface{}{common.BytesToAddress(validatorDstAddr.Bytes())},
	)
	if err != nil {
		return err
	}

	for _, topic := range topicsRes {
		topics = append(topics, topic[0])
	}

	// amount is part of event data
	data, err := packAmount(amount)
	if err != nil {
		return err
	}

	stateDB.AddLog(&types.Log{
		Address:     c.Address(),
		Topics:      topics,
		Data:        data,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

// helper function to pack a uint256 amount
func packAmount(amount *big.Int) ([]byte, error) {
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
