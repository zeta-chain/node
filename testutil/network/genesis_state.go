package network

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/assert"
	cmdcfg "github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func SetupZetaGenesisState(t *testing.T, genesisState map[string]json.RawMessage, codec codec.Codec, observerList []string) {

	// Cross-chain genesis state
	var crossChainGenesis types.GenesisState
	assert.NoError(t, codec.UnmarshalJSON(genesisState[types.ModuleName], &crossChainGenesis))
	nodeAccountList := make([]*observertypes.NodeAccount, len(observerList))
	for i, operator := range observerList {
		nodeAccountList[i] = &observertypes.NodeAccount{
			Operator:   operator,
			NodeStatus: observertypes.NodeStatus_Active,
		}
	}

	crossChainGenesis.Params.Enabled = true
	assert.NoError(t, crossChainGenesis.Validate())
	crossChainGenesisBz, err := codec.MarshalJSON(&crossChainGenesis)
	assert.NoError(t, err)

	// EVM genesis state
	var evmGenesisState evmtypes.GenesisState
	assert.NoError(t, codec.UnmarshalJSON(genesisState[evmtypes.ModuleName], &evmGenesisState))
	evmGenesisState.Params.EvmDenom = cmdcfg.BaseDenom
	assert.NoError(t, evmGenesisState.Validate())
	evmGenesisBz, err := codec.MarshalJSON(&evmGenesisState)
	assert.NoError(t, err)

	// Staking genesis state
	var stakingGenesisState stakingtypes.GenesisState
	assert.NoError(t, codec.UnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenesisState))
	stakingGenesisState.Params.BondDenom = cmdcfg.BaseDenom
	stakingGenesisStateBz, err := codec.MarshalJSON(&stakingGenesisState)
	assert.NoError(t, err)

	// Observer genesis state
	var observerGenesis observertypes.GenesisState
	assert.NoError(t, codec.UnmarshalJSON(genesisState[observertypes.ModuleName], &observerGenesis))
	observerMapper := make([]*observertypes.ObserverMapper, len(common.DefaultChainsList()))

	for i, chain := range common.DefaultChainsList() {
		observerMapper[i] = &observertypes.ObserverMapper{
			ObserverChain: chain,
			ObserverList:  observerList,
		}
	}
	observerGenesis.Observers = observerMapper
	observerGenesis.NodeAccountList = nodeAccountList
	observerGenesis.Keygen = &observertypes.Keygen{
		Status:         observertypes.KeygenStatus_PendingKeygen,
		GranteePubkeys: observerList,
		BlockNumber:    5,
	}
	assert.NoError(t, observerGenesis.Validate())
	observerGenesisBz, err := codec.MarshalJSON(&observerGenesis)
	assert.NoError(t, err)

	genesisState[types.ModuleName] = crossChainGenesisBz
	genesisState[stakingtypes.ModuleName] = stakingGenesisStateBz
	genesisState[observertypes.ModuleName] = observerGenesisBz
	genesisState[evmtypes.ModuleName] = evmGenesisBz
}

func AddObserverData(t *testing.T, genesisState map[string]json.RawMessage, codec codec.Codec, ballots []*observertypes.Ballot) *observertypes.GenesisState {
	state := observertypes.GenesisState{}
	assert.NoError(t, codec.UnmarshalJSON(genesisState[observertypes.ModuleName], &state))
	if len(ballots) > 0 {
		state.Ballots = ballots
	}
	//params := observertypes.DefaultParams()
	//params.BallotMaturityBlocks = 3
	state.Params.BallotMaturityBlocks = 3
	state.Keygen = &observertypes.Keygen{BlockNumber: 10, GranteePubkeys: []string{}}
	permissionFlags := &observertypes.PermissionFlags{}
	nullify.Fill(&permissionFlags)
	state.PermissionFlags = permissionFlags

	assert.NoError(t, state.Validate())

	buf, err := codec.MarshalJSON(&state)
	assert.NoError(t, err)
	genesisState[observertypes.ModuleName] = buf
	return &state
}
func AddCrosschainData(t *testing.T, n int, genesisState map[string]json.RawMessage, codec codec.Codec) *types.GenesisState {
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

	assert.NoError(t, state.Validate())

	buf, err := codec.MarshalJSON(&state)
	assert.NoError(t, err)
	genesisState[types.ModuleName] = buf
	return &state
}
