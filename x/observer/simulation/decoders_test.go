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
			{Key: types.KeyPrefix(types.TSSHistoryKey), Value: cdc.MustMarshal(&tss)},
			{Key: types.KeyPrefix(types.ObserverSetKey), Value: cdc.MustMarshal(&observerSet)},
			{Key: types.KeyPrefix(types.AllChainParamsKey), Value: cdc.MustMarshal(&chainParamsList)},
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
		{
			"CrosschainFlags",
			fmt.Sprintf("key %s value A %v value B %v", types.CrosschainFlagsKey, *crosschainFlags, *crosschainFlags),
		},
		{
			"LastBlockObserverCount",
			fmt.Sprintf(
				"key %s value A %v value B %v",
				types.LastBlockObserverCountKey,
				*lastBlockObserverCount,
				*lastBlockObserverCount,
			),
		},
		{"NodeAccount", fmt.Sprintf("key %s value A %v value B %v", types.NodeAccountKey, *nodeAccount, *nodeAccount)},
		{"Keygen", fmt.Sprintf("key %s value A %v value B %v", types.KeygenKey, *keygen, *keygen)},
		{"BallotList", fmt.Sprintf("key %s value A %v value B %v", types.BallotListKey, ballotList, ballotList)},
		{"Ballot", fmt.Sprintf("key %s value A %v value B %v", types.VoterKey, *ballot, *ballot)},
		{"TSS", fmt.Sprintf("key %s value A %v value B %v", types.TSSKey, tss, tss)},
		{"TSSHistory", fmt.Sprintf("key %s value A %v value B %v", types.TSSHistoryKey, tss, tss)},
		{"ObserverSet", fmt.Sprintf("key %s value A %v value B %v", types.ObserverSetKey, observerSet, observerSet)},
		{
			"ChainParamsList",
			fmt.Sprintf("key %s value A %v value B %v", types.AllChainParamsKey, chainParamsList, chainParamsList),
		},
		{
			"TssFundMigrator",
			fmt.Sprintf("key %s value A %v value B %v", types.TssFundMigratorKey, tssFundMigrator, tssFundMigrator),
		},
		{
			"PendingNonces",
			fmt.Sprintf("key %s value A %v value B %v", types.PendingNoncesKeyPrefix, pendingNonce, pendingNonce),
		},
		{"ChainNonces", fmt.Sprintf("key %s value A %v value B %v", types.ChainNoncesKey, chainNonces, chainNonces)},
		{
			"NonceToCctx",
			fmt.Sprintf("key %s value A %v value B %v", types.NonceToCctxKeyPrefix, nonceToCctx, nonceToCctx),
		},
		{"Params", fmt.Sprintf("key %s value A %v value B %v", types.ParamsKey, params, params)},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]))
		})
	}
}
