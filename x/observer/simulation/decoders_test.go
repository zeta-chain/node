package simulation_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/simulation"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestNewDecodeStore(t *testing.T) {
	k, _, _, _ := keepertest.ObserverKeeper(t)
	cdc := k.Codec()
	dec := simulation.NewDecodeStore(cdc)
	crosschainFlags := sample.CrosschainFlags()
	lastBlockObserverCount := sample.LastObserverCount(10)
	nodeAccount := sample.NodeAccount()
	keygen := sample.Keygen(t)

	ballotList := types.BallotListForHeight{
		Height:           10,
		BallotsIndexList: []string{sample.ZetaIndex(t)},
	}

	ballot := sample.Ballot(t, "sample")
	tss := sample.Tss()
	observerSet := sample.ObserverSet(10)
	chainParamsList := sample.ChainParamsList()
	//tssHistory := sample.TssList(10)
	tssFundMigrator := sample.TssFundsMigrator(chains.Ethereum.ChainId)
	pendingNonce := sample.PendingNoncesList(t, "index", 10)[0]
	chainNonces := sample.ChainNonces(chains.Ethereum.ChainId)
	nonceToCctx := sample.NonceToCCTX(t, "index")
	params := types.Params{BallotMaturityBlocks: 100}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.KeyPrefix(types.CrosschainFlagsKey), Value: cdc.MustMarshal(crosschainFlags)},
			{Key: types.KeyPrefix(types.LastBlockObserverCountKey), Value: cdc.MustMarshal(lastBlockObserverCount)},
			{Key: types.KeyPrefix(types.NodeAccountKey), Value: cdc.MustMarshal(nodeAccount)},
			{Key: types.KeyPrefix(types.KeygenKey), Value: cdc.MustMarshal(keygen)},
			{Key: types.KeyPrefix(types.BallotListKey), Value: cdc.MustMarshal(&ballotList)},
			{Key: types.KeyPrefix(types.VoterKey), Value: cdc.MustMarshal(ballot)},
			{Key: types.KeyPrefix(types.TSSKey), Value: cdc.MustMarshal(&tss)},
			{Key: types.KeyPrefix(types.ObserverSetKey), Value: cdc.MustMarshal(&observerSet)},
			{Key: types.KeyPrefix(types.AllChainParamsKey), Value: cdc.MustMarshal(&chainParamsList)},
			{Key: types.KeyPrefix(types.TSSHistoryKey), Value: cdc.MustMarshal(&tss)},
			{Key: types.KeyPrefix(types.TssFundMigratorKey), Value: cdc.MustMarshal(&tssFundMigrator)},
			{Key: types.KeyPrefix(types.PendingNoncesKeyPrefix), Value: cdc.MustMarshal(&pendingNonce)},
			{Key: types.KeyPrefix(types.ChainNoncesKey), Value: cdc.MustMarshal(&chainNonces)},
			{Key: types.KeyPrefix(types.NonceToCctxKeyPrefix), Value: cdc.MustMarshal(&nonceToCctx)},
			{Key: types.KeyPrefix(types.ParamsKey), Value: cdc.MustMarshal(&params)},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"CrosschainFlags", fmt.Sprintf("%v\n%v", *crosschainFlags, *crosschainFlags)},
		{"LastBlockObserverCount", fmt.Sprintf("%v\n%v", *lastBlockObserverCount, *lastBlockObserverCount)},
		{"NodeAccount", fmt.Sprintf("%v\n%v", *nodeAccount, *nodeAccount)},
		{"Keygen", fmt.Sprintf("%v\n%v", *keygen, *keygen)},
		{"BallotList", fmt.Sprintf("%v\n%v", ballotList, ballotList)},
		{"Ballot", fmt.Sprintf("%v\n%v", *ballot, *ballot)},
		{"TSS", fmt.Sprintf("%v\n%v", tss, tss)},
		{"TSSHistory", fmt.Sprintf("%v\n%v", tss, tss)},
		{"ObserverSet", fmt.Sprintf("%v\n%v", observerSet, observerSet)},
		{"ChainParamsList", fmt.Sprintf("%v\n%v", chainParamsList, chainParamsList)},
		{"TssFundMigrator", fmt.Sprintf("%v\n%v", tssFundMigrator, tssFundMigrator)},
		{"PendingNonces", fmt.Sprintf("%v\n%v", pendingNonce, pendingNonce)},
		{"ChainNonces", fmt.Sprintf("%v\n%v", chainNonces, chainNonces)},
		{"NonceToCctx", fmt.Sprintf("%v\n%v", nonceToCctx, nonceToCctx)},
		{"Params", fmt.Sprintf("%v\n%v", params, params)},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]))
		})
	}
}
