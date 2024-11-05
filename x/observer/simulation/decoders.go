package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/zeta-chain/node/x/observer/types"
)

// TODO Add comments in this pr to explain the purpose of the function
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.CrosschainFlagsKey)):
			var crosschainFlagsA, crosschainFlagsB types.CrosschainFlags
			cdc.MustUnmarshal(kvA.Value, &crosschainFlagsA)
			cdc.MustUnmarshal(kvB.Value, &crosschainFlagsB)
			return fmt.Sprintf("%v\n%v", crosschainFlagsA, crosschainFlagsB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.LastBlockObserverCountKey)):
			var lastBlockObserverCountA, lastBlockObserverCountB types.LastObserverCount
			cdc.MustUnmarshal(kvA.Value, &lastBlockObserverCountA)
			cdc.MustUnmarshal(kvB.Value, &lastBlockObserverCountB)
			return fmt.Sprintf("%v\n%v", lastBlockObserverCountA, lastBlockObserverCountB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.NodeAccountKey)):
			var nodeAccountA, nodeAccountB types.NodeAccount
			cdc.MustUnmarshal(kvA.Value, &nodeAccountA)
			cdc.MustUnmarshal(kvB.Value, &nodeAccountB)
			return fmt.Sprintf("%v\n%v", nodeAccountA, nodeAccountB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.KeygenKey)):
			var keygenA, keygenB types.Keygen
			cdc.MustUnmarshal(kvA.Value, &keygenA)
			cdc.MustUnmarshal(kvB.Value, &keygenB)
			return fmt.Sprintf("%v\n%v", keygenA, keygenB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.BallotListKey)):
			var ballotListA, ballotListB types.BallotListForHeight
			cdc.MustUnmarshal(kvA.Value, &ballotListA)
			cdc.MustUnmarshal(kvB.Value, &ballotListB)
			return fmt.Sprintf("%v\n%v", ballotListA, ballotListB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.VoterKey)):
			var voterA, voterB types.Ballot
			cdc.MustUnmarshal(kvA.Value, &voterA)
			cdc.MustUnmarshal(kvB.Value, &voterB)
			return fmt.Sprintf("%v\n%v", voterA, voterB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.TSSKey)):
			var tssA, tssB types.TSS
			cdc.MustUnmarshal(kvA.Value, &tssA)
			cdc.MustUnmarshal(kvB.Value, &tssB)
			return fmt.Sprintf("%v\n%v", tssA, tssB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.ObserverSetKey)):
			var observerSetA, observerSetB types.ObserverSet
			cdc.MustUnmarshal(kvA.Value, &observerSetA)
			cdc.MustUnmarshal(kvB.Value, &observerSetB)
			return fmt.Sprintf("%v\n%v", observerSetA, observerSetB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.AllChainParamsKey)):
			var allChainParamsA, allChainParamsB types.ChainParamsList
			cdc.MustUnmarshal(kvA.Value, &allChainParamsA)
			cdc.MustUnmarshal(kvB.Value, &allChainParamsB)
			return fmt.Sprintf("%v\n%v", allChainParamsA, allChainParamsB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.TSSHistoryKey)):
			var tssHistoryA, tssHistoryB []types.TSS
			return fmt.Sprintf("%v\n%v", tssHistoryA, tssHistoryB)
			//cdc.MustUnmarshal(kvA.Value, &tssHistoryA)
			//cdc.MustUnmarshal(kvB.Value, &tssHistoryB)
			//return fmt.Sprintf("%v\n%v", tssHistoryA, tssHistoryB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.TssFundMigratorKey)):
			var tssFundMigratorA, tssFundMigratorB types.TssFundMigratorInfo
			cdc.MustUnmarshal(kvA.Value, &tssFundMigratorA)
			cdc.MustUnmarshal(kvB.Value, &tssFundMigratorB)
			return fmt.Sprintf("%v\n%v", tssFundMigratorA, tssFundMigratorB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.PendingNoncesKeyPrefix)):
			var pendingNoncesA, pendingNoncesB types.PendingNonces
			cdc.MustUnmarshal(kvA.Value, &pendingNoncesA)
			cdc.MustUnmarshal(kvB.Value, &pendingNoncesB)
			return fmt.Sprintf("%v\n%v", pendingNoncesA, pendingNoncesB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.ChainNoncesKey)):
			var chainNoncesA, chainNoncesB types.ChainNonces
			cdc.MustUnmarshal(kvA.Value, &chainNoncesA)
			cdc.MustUnmarshal(kvB.Value, &chainNoncesB)
			return fmt.Sprintf("%v\n%v", chainNoncesA, chainNoncesB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.NonceToCctxKeyPrefix)):
			var nonceToCctxA, nonceToCctxB types.NonceToCctx
			cdc.MustUnmarshal(kvA.Value, &nonceToCctxA)
			cdc.MustUnmarshal(kvB.Value, &nonceToCctxB)
			return fmt.Sprintf("%v\n%v", nonceToCctxA, nonceToCctxB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.ParamsKey)):
			var paramsA, paramsB types.Params
			cdc.MustUnmarshal(kvA.Value, &paramsA)
			cdc.MustUnmarshal(kvB.Value, &paramsB)
			return fmt.Sprintf("%v\n%v", paramsA, paramsB)
		default:
			panic(fmt.Sprintf("invalid observer key prefix %X", kvA.Key))
		}
	}
}
