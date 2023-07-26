package network

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethcfg "github.com/evmos/ethermint/cmd/config"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/app"
	cmdcfg "github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"testing"
)

func Setconfig() {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	ethcfg.SetBip44CoinType(config)
	// Make sure address is compatible with ethereum
	config.SetAddressVerifier(app.VerifyAddressFormat)
	//config.Seal()
}

func SetupZetaGenesisState(t *testing.T, genesisState map[string]json.RawMessage, codec codec.Codec, observerList []string) map[string]json.RawMessage {

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
	return genesisState
}
