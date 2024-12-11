package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

// Simulation operation weights constants
// Operation weights are used by the simulation program to simulate the weight of different operations.
// This decides what percentage of a certain type of operation is part of a block.
// Based on the weights assigned in the cosmos sdk modules , 100 seems to the max weight used , and therefore guarantees that at least one operation of that type is present in a block.
// TODO Add more details to comment based on what the number represents in terms of percentage of operations in a block
// https://github.com/zeta-chain/node/issues/3100
const (
	// #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgEnableCCTX                  = "op_weight_msg_enable_crosschain_flags"
	OpWeightMsgTypeMsgDisableCCTX                 = "op_weight_msg_disable_crosschain_flags"
	OpWeightMsgTypeMsgVoteTSS                     = "op_weight_msg_vote_tss"
	OpWeightMsgTypeMsgUpdateKeygen                = "op_weight_msg_update_keygen"
	OpWeightMsgTypeMsgUpdateObserver              = "op_weight_msg_update_observer"
	OpWeightMsgTypeMsgUpdateChainParams           = "op_weight_msg_update_chain_params"
	OpWeightMsgTypeMsgRemoveChainParams           = "op_weight_msg_remove_chain_params"
	OpWeightMsgTypeMsgResetChainNonces            = "op_weight_msg_reset_chain_nonces"
	OpWeightMsgTypeMsgUpdateGasPriceIncreaseFlags = "op_weight_msg_update_gas_price_increase_flags"
	OpWeightMsgTypeMsgAddObserver                 = "op_weight_msg_add_observer"

	// DefaultWeightMsgTypeMsgEnableCCTX We use a high weight for this operation
	// to ensure that it is present in the block more number of times than any operation that changes the validator set
	// Arrived at this number based on the weights used in the cosmos sdk staking module and through some trial and error
	DefaultWeightMsgTypeMsgEnableCCTX                  = 3650
	DefaultWeightMsgTypeMsgDisableCCTX                 = 10
	DefaultWeightMsgTypeMsgVoteTSS                     = 10
	DefaultWeightMsgTypeMsgUpdateKeygen                = 10
	DefaultWeightMsgTypeMsgUpdateObserver              = 10
	DefaultWeightMsgTypeMsgUpdateChainParams           = 10
	DefaultWeightMsgTypeMsgRemoveChainParams           = 10
	DefaultWeightMsgTypeMsgResetChainNonces            = 5
	DefaultWeightMsgTypeMsgUpdateGasPriceIncreaseFlags = 10
	DefaultWeightMsgTypeMsgAddObserver                 = 5
)

// WeightedOperations for observer module
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, k keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgTypeMsgEnableCCTX                  int
		weightMsgTypeMsgDisableCCTX                 int
		weightMsgTypeMsgVoteTSS                     int
		weightMsgTypeMsgUpdateKeygen                int
		weightMsgTypeMsgUpdateObserver              int
		weightMsgTypeMsgUpdateChainParams           int
		weightMsgTypeMsgRemoveChainParams           int
		weightMsgTypeMsgResetChainNonces            int
		weightMsgTypeMsgUpdateGasPriceIncreaseFlags int
		weightMsgTypeMsgAddObserver                 int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgTypeMsgEnableCCTX, &weightMsgTypeMsgEnableCCTX, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgEnableCCTX = DefaultWeightMsgTypeMsgEnableCCTX
		})

	appParams.GetOrGenerate(cdc, OpWeightMsgTypeMsgDisableCCTX, &weightMsgTypeMsgDisableCCTX, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgDisableCCTX = DefaultWeightMsgTypeMsgDisableCCTX
		})

	appParams.GetOrGenerate(cdc, OpWeightMsgTypeMsgVoteTSS, &weightMsgTypeMsgVoteTSS, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgVoteTSS = DefaultWeightMsgTypeMsgVoteTSS
		})

	appParams.GetOrGenerate(cdc, OpWeightMsgTypeMsgUpdateKeygen, &weightMsgTypeMsgUpdateKeygen, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgUpdateKeygen = DefaultWeightMsgTypeMsgUpdateKeygen
		})

	appParams.GetOrGenerate(cdc, OpWeightMsgTypeMsgUpdateObserver, &weightMsgTypeMsgUpdateObserver, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgUpdateObserver = DefaultWeightMsgTypeMsgUpdateObserver
		})

	appParams.GetOrGenerate(cdc, OpWeightMsgTypeMsgUpdateChainParams, &weightMsgTypeMsgUpdateChainParams, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgUpdateChainParams = DefaultWeightMsgTypeMsgUpdateChainParams
		})

	appParams.GetOrGenerate(cdc, OpWeightMsgTypeMsgRemoveChainParams, &weightMsgTypeMsgRemoveChainParams, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgRemoveChainParams = DefaultWeightMsgTypeMsgRemoveChainParams
		})

	appParams.GetOrGenerate(cdc, OpWeightMsgTypeMsgResetChainNonces, &weightMsgTypeMsgResetChainNonces, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgResetChainNonces = DefaultWeightMsgTypeMsgResetChainNonces
		})

	appParams.GetOrGenerate(
		cdc,
		OpWeightMsgTypeMsgUpdateGasPriceIncreaseFlags,
		&weightMsgTypeMsgUpdateGasPriceIncreaseFlags,
		nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgUpdateGasPriceIncreaseFlags = DefaultWeightMsgTypeMsgUpdateGasPriceIncreaseFlags
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgTypeMsgAddObserver, &weightMsgTypeMsgAddObserver, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgAddObserver = DefaultWeightMsgTypeMsgAddObserver
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgTypeMsgEnableCCTX,
			SimulateMsgEnableCCTX(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgDisableCCTX,
			SimulateMsgDisableCCTX(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgUpdateKeygen,
			SimulateMsgUpdateKeygen(k),
		),
		//
		simulation.NewWeightedOperation(
			weightMsgTypeMsgUpdateChainParams,
			SimulateMsgUpdateChainParams(k),
		),
		//
		//simulation.NewWeightedOperation(
		//	weightMsgTypeMsgRemoveChainParams,
		//	SimulateMsgRemoveChainParams(k),
		//),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgResetChainNonces,
			SimulateMsgResetChainNonces(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgUpdateGasPriceIncreaseFlags,
			SimulateMsgUpdateGasPriceIncreaseFlags(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgAddObserver,
			SimulateMsgUpdateObserver(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgAddObserver,
			SimulateMsgAddObserverNodeAccount(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgAddObserver,
			SimulateMsgAddObserver(k),
		),
	}
}

func GetPolicyAccount(ctx sdk.Context, k types.AuthorityKeeper, accounts []simtypes.Account) (simtypes.Account, error) {
	policies, found := k.GetPolicies(ctx)
	if !found {
		return simtypes.Account{}, fmt.Errorf("policies object not found")
	}
	if len(policies.Items) == 0 {
		return simtypes.Account{}, fmt.Errorf("no policies found")
	}

	admin := policies.Items[0].Address
	address, err := types.GetOperatorAddressFromAccAddress(admin)
	if err != nil {
		return simtypes.Account{}, err
	}
	simAccount, found := simtypes.FindAccount(accounts, address)
	if !found {
		return simtypes.Account{}, fmt.Errorf("admin account not found in list of simulation accounts")
	}
	return simAccount, nil
}

func GetExternalChain(ctx sdk.Context, k keeper.Keeper, r *rand.Rand, retryCount int) (chains.Chain, error) {
	supportedChains := k.GetSupportedChains(ctx)
	if len(supportedChains) == 0 {
		return chains.Chain{}, fmt.Errorf("no supported chains found")
	}
	// remove zeta chain from the supported chains
	for i := 0; i < retryCount; i++ {
		c := supportedChains[r.Intn(len(supportedChains))]
		if !c.IsZetaChain() {
			return c, nil
		}
	}
	return chains.Chain{}, fmt.Errorf("no external chain found")
}

// GetRandomAccountAndObserver returns a random account and the associated observer address
func GetRandomAccountAndObserver(
	r *rand.Rand,
	ctx sdk.Context,
	k keeper.Keeper,
	accounts []simtypes.Account,
) (simtypes.Account, string, []string, error) {
	observerList := []string{}
	observers, found := k.GetObserverSet(ctx)
	if !found {
		return simtypes.Account{}, "", observerList, fmt.Errorf("observer set not found")
	}

	observerList = observers.ObserverList

	if len(observers.ObserverList) == 0 {
		return simtypes.Account{}, "", observerList, fmt.Errorf("no observers present in observer set found")
	}

	randomObserver := ""
	foundObserver := false
	for i := 0; i < 10; i++ {
		randomObserver = GetRandomObserver(r, observers.ObserverList)
		_, foundNodeAccount := k.GetNodeAccount(ctx, randomObserver)
		if !foundNodeAccount {
			continue
		}
		ok := k.IsNonTombstonedObserver(ctx, randomObserver)
		if ok {
			foundObserver = true
			break
		}
	}

	if !foundObserver {
		return simtypes.Account{}, "", observerList, fmt.Errorf("no observer found")
	}

	simAccount, err := GetObserverAccount(randomObserver, accounts)
	if err != nil {
		return simtypes.Account{}, "", observerList, err
	}
	return simAccount, randomObserver, observerList, nil
}

func GetRandomObserver(r *rand.Rand, observerList []string) string {
	idx := r.Intn(len(observerList))
	return observerList[idx]
}

// GetObserverAccount returns the account associated with the observer address from the list of accounts provided
// GetObserverAccount can fail if all the observers are removed from the observer set ,this can happen
//if the other modules create transactions which affect the validator
//and triggers any of the staking hooks defined in the observer modules

func GetObserverAccount(observerAddress string, accounts []simtypes.Account) (simtypes.Account, error) {
	operatorAddress, err := types.GetOperatorAddressFromAccAddress(observerAddress)
	if err != nil {
		return simtypes.Account{}, fmt.Errorf("validator not found for observer ")
	}

	simAccount, found := simtypes.FindAccount(accounts, operatorAddress)
	if !found {
		return simtypes.Account{}, fmt.Errorf("operator account not found")
	}
	return simAccount, nil
}
