package network

import (
	"cosmossdk.io/math"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/assert"
	cmdcfg "github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"strconv"
	"testing"
)

func SetupZetaGenesisState(t *testing.T, genesisState map[string]json.RawMessage, codec codec.Codec, observerList []string) {

	// Cross-chain genesis state
	var crossChainGenesis types.GenesisState
	assert.NoError(t, codec.UnmarshalJSON(genesisState[types.ModuleName], &crossChainGenesis))
	nodeAccountList := make([]*observerTypes.NodeAccount, len(observerList))
	for i, operator := range observerList {
		nodeAccountList[i] = &observerTypes.NodeAccount{
			Operator:   operator,
			NodeStatus: observerTypes.NodeStatus_Active,
		}
	}

	crossChainGenesis.Params.Enabled = true
	crossChainGenesisBz, err := codec.MarshalJSON(&crossChainGenesis)
	assert.NoError(t, err)

	// Cross-chain genesis state
	var evmGenesisState evmtypes.GenesisState
	assert.NoError(t, codec.UnmarshalJSON(genesisState[evmtypes.ModuleName], &evmGenesisState))
	evmGenesisState.Params.EvmDenom = cmdcfg.BaseDenom
	evmGenesisBz, err := codec.MarshalJSON(&evmGenesisState)
	assert.NoError(t, err)

	// Staking genesis state
	var stakingGenesisState stakingtypes.GenesisState
	assert.NoError(t, codec.UnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenesisState))
	stakingGenesisState.Params.BondDenom = cmdcfg.BaseDenom
	stakingGenesisStateBz, err := codec.MarshalJSON(&stakingGenesisState)
	assert.NoError(t, err)

	// Observer genesis state
	var observerGenesis observerTypes.GenesisState
	assert.NoError(t, codec.UnmarshalJSON(genesisState[observerTypes.ModuleName], &observerGenesis))
	observerMapper := make([]*observerTypes.ObserverMapper, len(common.DefaultChainsList()))

	for i, chain := range common.DefaultChainsList() {
		observerMapper[i] = &observerTypes.ObserverMapper{
			ObserverChain: chain,
			ObserverList:  observerList,
		}
	}
	observerGenesis.Observers = observerMapper
	observerGenesis.NodeAccountList = nodeAccountList
	observerGenesis.Keygen = &observerTypes.Keygen{
		Status:         observerTypes.KeygenStatus_PendingKeygen,
		GranteePubkeys: observerList,
		BlockNumber:    5,
	}
	observerGenesisBz, err := codec.MarshalJSON(&observerGenesis)
	assert.NoError(t, err)

	genesisState[types.ModuleName] = crossChainGenesisBz
	genesisState[stakingtypes.ModuleName] = stakingGenesisStateBz
	genesisState[observerTypes.ModuleName] = observerGenesisBz
	genesisState[evmtypes.ModuleName] = evmGenesisBz
}

func AddObserverData(n int, t *testing.T, genesisState map[string]json.RawMessage, codec codec.Codec, ballots []*observerTypes.Ballot) *observerTypes.GenesisState {
	state := observerTypes.GenesisState{}
	assert.NoError(t, codec.UnmarshalJSON(genesisState[observerTypes.ModuleName], &state))
	//for i := 0; i < n; i++ {
	//	state.NodeAccountList = append(state.NodeAccountList, &observerTypes.NodeAccount{Operator: strconv.Itoa(i), GranteeAddress: "signer"})
	//}
	if len(ballots) > 0 {
		state.Ballots = ballots
	}
	state.Keygen = &observerTypes.Keygen{BlockNumber: 10, GranteePubkeys: []string{}}
	permissionFlags := &observerTypes.PermissionFlags{}
	nullify.Fill(&permissionFlags)
	state.PermissionFlags = permissionFlags

	buf, err := codec.MarshalJSON(&state)
	assert.NoError(t, err)
	genesisState[observerTypes.ModuleName] = buf
	return &state
}
func AddCrossChainData(n int, t *testing.T, genesisState map[string]json.RawMessage, codec codec.Codec) *types.GenesisState {
	state := types.GenesisState{}
	assert.NoError(t, codec.UnmarshalJSON(genesisState[types.ModuleName], &state))
	// TODO : Fix add EVM balance to deploy contracts
	for i := 0; i < n; i++ {
		state.CrossChainTxs = append(state.CrossChainTxs, &types.CrossChainTx{
			Creator: "ANY",
			Index:   strconv.Itoa(i),
			CctxStatus: &types.Status{
				Status:              types.CctxStatus_PendingInbound,
				StatusMessage:       "",
				LastUpdateTimestamp: 0,
			},
			InboundTxParams:  &types.InboundTxParams{InboundTxObservedHash: fmt.Sprintf("Hash-%d", i), Amount: math.OneUint()},
			OutboundTxParams: []*types.OutboundTxParams{},
			ZetaFees:         math.OneUint()},
		)
	}
	for i := 0; i < n; i++ {
		state.ChainNoncesList = append(state.ChainNoncesList, &types.ChainNonces{Creator: "ANY", Index: strconv.Itoa(i), Signers: []string{}})
	}
	for i := 0; i < n; i++ {
		state.GasPriceList = append(state.GasPriceList, &types.GasPrice{Creator: "ANY", ChainId: int64(i), Index: strconv.Itoa(i), Prices: []uint64{}, BlockNums: []uint64{}, Signers: []string{}})
	}
	for i := 0; i < n; i++ {
		state.LastBlockHeightList = append(state.LastBlockHeightList, &types.LastBlockHeight{Creator: "ANY", Index: strconv.Itoa(i)})
	}

	state.Tss = &types.TSS{
		TssPubkey:           "tssPubkey",
		TssParticipantList:  []string{"tssParticipantList"},
		OperatorAddressList: []string{"operatorAddressList"},
		FinalizedZetaHeight: 1,
		KeyGenZetaHeight:    1,
	}
	for i := 0; i < n; i++ {
		outTxTracker := types.OutTxTracker{
			Index:   fmt.Sprintf("%d-%d", i, i),
			ChainId: int64(i),
			Nonce:   uint64(i),
		}
		nullify.Fill(&outTxTracker)
		state.OutTxTrackerList = append(state.OutTxTrackerList, outTxTracker)
	}

	for i := 0; i < n; i++ {
		inTxHashToCctx := types.InTxHashToCctx{
			InTxHash: strconv.Itoa(i),
		}
		nullify.Fill(&inTxHashToCctx)
		state.InTxHashToCctxList = append(state.InTxHashToCctxList, inTxHashToCctx)
	}

	buf, err := codec.MarshalJSON(&state)
	assert.NoError(t, err)
	genesisState[types.ModuleName] = buf
	return &state
}
