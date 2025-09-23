package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/zeta-chain/node/x/observer/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding observer types.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.CrosschainFlagsKey)):
			var crosschainFlagsA, crosschainFlagsB types.CrosschainFlags
			cdc.MustUnmarshal(kvA.Value, &crosschainFlagsA)
			cdc.MustUnmarshal(kvB.Value, &crosschainFlagsB)
			return fmt.Sprintf(
				"key %s value A %v value B %v",
				types.CrosschainFlagsKey,
				crosschainFlagsA,
				crosschainFlagsB,
			)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.LastBlockObserverCountKey)):
			var lastBlockObserverCountA, lastBlockObserverCountB types.LastObserverCount
			cdc.MustUnmarshal(kvA.Value, &lastBlockObserverCountA)
			cdc.MustUnmarshal(kvB.Value, &lastBlockObserverCountB)
			return fmt.Sprintf(
				"key %s value A %v value B %v",
				types.LastBlockObserverCountKey,
				lastBlockObserverCountA,
				lastBlockObserverCountB,
			)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.NodeAccountKey)):
			var nodeAccountA, nodeAccountB types.NodeAccount
			cdc.MustUnmarshal(kvA.Value, &nodeAccountA)
			cdc.MustUnmarshal(kvB.Value, &nodeAccountB)
			return fmt.Sprintf("key %s value A %v value B %v", types.NodeAccountKey, nodeAccountA, nodeAccountB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.KeygenKey)):
			var keygenA, keygenB types.Keygen
			cdc.MustUnmarshal(kvA.Value, &keygenA)
			cdc.MustUnmarshal(kvB.Value, &keygenB)
			return fmt.Sprintf("key %s value A %v value B %v", types.KeygenKey, keygenA, keygenB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.BallotListKey)):
			var ballotListA, ballotListB types.BallotListForHeight
			cdc.MustUnmarshal(kvA.Value, &ballotListA)
			cdc.MustUnmarshal(kvB.Value, &ballotListB)
			return fmt.Sprintf("key %s value A %v value B %v", types.BallotListKey, ballotListA, ballotListB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.VoterKey)):
			var voterA, voterB types.Ballot
			cdc.MustUnmarshal(kvA.Value, &voterA)
			cdc.MustUnmarshal(kvB.Value, &voterB)
			return fmt.Sprintf("key %s value A %v value B %v", types.VoterKey, voterA, voterB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.TSSKey)):
			var tssA, tssB types.TSS
			cdc.MustUnmarshal(kvA.Value, &tssA)
			cdc.MustUnmarshal(kvB.Value, &tssB)
			return fmt.Sprintf("key %s value A %v value B %v", types.TSSKey, tssA, tssB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.ObserverSetKey)):
			var observerSetA, observerSetB types.ObserverSet
			cdc.MustUnmarshal(kvA.Value, &observerSetA)
			cdc.MustUnmarshal(kvB.Value, &observerSetB)
			return fmt.Sprintf("key %s value A %v value B %v", types.ObserverSetKey, observerSetA, observerSetB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.AllChainParamsKey)):
			var allChainParamsA, allChainParamsB types.ChainParamsList
			cdc.MustUnmarshal(kvA.Value, &allChainParamsA)
			cdc.MustUnmarshal(kvB.Value, &allChainParamsB)
			return fmt.Sprintf(
				"key %s value A %v value B %v",
				types.AllChainParamsKey,
				allChainParamsA,
				allChainParamsB,
			)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.TSSHistoryKey)):
			var tssHistoryA, tssHistoryB types.TSS
			cdc.MustUnmarshal(kvA.Value, &tssHistoryA)
			cdc.MustUnmarshal(kvB.Value, &tssHistoryB)
			return fmt.Sprintf("key %s value A %v value B %v", types.TSSHistoryKey, tssHistoryA, tssHistoryB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.TssFundMigratorKey)):
			var tssFundMigratorA, tssFundMigratorB types.TssFundMigratorInfo
			cdc.MustUnmarshal(kvA.Value, &tssFundMigratorA)
			cdc.MustUnmarshal(kvB.Value, &tssFundMigratorB)
			return fmt.Sprintf(
				"key %s value A %v value B %v",
				types.TssFundMigratorKey,
				tssFundMigratorA,
				tssFundMigratorB,
			)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.PendingNoncesKeyPrefix)):
			var pendingNoncesA, pendingNoncesB types.PendingNonces
			cdc.MustUnmarshal(kvA.Value, &pendingNoncesA)
			cdc.MustUnmarshal(kvB.Value, &pendingNoncesB)
			return fmt.Sprintf(
				"key %s value A %v value B %v",
				types.PendingNoncesKeyPrefix,
				pendingNoncesA,
				pendingNoncesB,
			)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.ChainNoncesKey)):
			var chainNoncesA, chainNoncesB types.ChainNonces
			cdc.MustUnmarshal(kvA.Value, &chainNoncesA)
			cdc.MustUnmarshal(kvB.Value, &chainNoncesB)
			return fmt.Sprintf("key %s value A %v value B %v", types.ChainNoncesKey, chainNoncesA, chainNoncesB)
		case bytes.Equal(kvA.Key, types.KeyPrefix(types.NonceToCctxKeyPrefix)):
			var nonceToCctxA, nonceToCctxB types.NonceToCctx
			cdc.MustUnmarshal(kvA.Value, &nonceToCctxA)
			cdc.MustUnmarshal(kvB.Value, &nonceToCctxB)
			return fmt.Sprintf("key %s value A %v value B %v", types.NonceToCctxKeyPrefix, nonceToCctxA, nonceToCctxB)
		default:
			panic(
				fmt.Sprintf(
					"invalid observer key prefix %X (first 8 bytes: %X)",
					kvA.Key[:1],
					kvA.Key[:min(8, len(kvA.Key))],
				),
			)
		}
	}
}
