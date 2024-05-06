package sample

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/cosmos"
	zetacrypto "github.com/zeta-chain/zetacore/pkg/crypto"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func Ballot(t *testing.T, index string) *types.Ballot {
	r := newRandFromStringSeed(t, index)

	return &types.Ballot{
		Index:                index,
		BallotIdentifier:     StringRandom(r, 32),
		VoterList:            []string{AccAddress(), AccAddress()},
		Votes:                []types.VoteType{types.VoteType_FailureObservation, types.VoteType_SuccessObservation},
		ObservationType:      types.ObservationType_EmptyObserverType,
		BallotThreshold:      sdk.NewDec(1),
		BallotStatus:         types.BallotStatus_BallotInProgress,
		BallotCreationHeight: r.Int63(),
	}
}

func ObserverSet(n int) types.ObserverSet {
	observerList := make([]string, n)
	for i := 0; i < n; i++ {
		observerList[i] = AccAddress()
	}

	return types.ObserverSet{
		ObserverList: observerList,
	}
}

func NodeAccount() *types.NodeAccount {
	return &types.NodeAccount{
		Operator:       AccAddress(),
		GranteeAddress: AccAddress(),
		GranteePubkey:  PubKeySet(),
		NodeStatus:     types.NodeStatus_Active,
	}
}

func CrosschainFlags() *types.CrosschainFlags {
	return &types.CrosschainFlags{
		IsInboundEnabled:  true,
		IsOutboundEnabled: true,
	}
}

func Keygen(t *testing.T) *types.Keygen {
	pubKey := ed25519.GenPrivKey().PubKey().String()
	r := newRandFromStringSeed(t, pubKey)

	return &types.Keygen{
		Status:         types.KeygenStatus_KeyGenSuccess,
		GranteePubkeys: []string{pubKey},
		BlockNumber:    r.Int63(),
	}
}

func LastObserverCount(lastChangeHeight int64) *types.LastObserverCount {
	r := newRandFromSeed(lastChangeHeight)

	return &types.LastObserverCount{
		Count:            r.Uint64(),
		LastChangeHeight: lastChangeHeight,
	}
}

func ChainParams(chainID int64) *types.ChainParams {
	r := newRandFromSeed(chainID)

	fiftyPercent, err := sdk.NewDecFromStr("0.5")
	if err != nil {
		return nil
	}

	return &types.ChainParams{
		ChainId:           chainID,
		ConfirmationCount: r.Uint64(),

		GasPriceTicker:              Uint64InRange(1, 300),
		InboundTicker:               Uint64InRange(1, 300),
		OutboundTicker:              Uint64InRange(1, 300),
		WatchUtxoTicker:             Uint64InRange(1, 300),
		ZetaTokenContractAddress:    EthAddress().String(),
		ConnectorContractAddress:    EthAddress().String(),
		Erc20CustodyContractAddress: EthAddress().String(),
		OutboundScheduleInterval:    Int64InRange(1, 100),
		OutboundScheduleLookahead:   Int64InRange(1, 500),
		BallotThreshold:             fiftyPercent,
		MinObserverDelegation:       sdk.NewDec(r.Int63()),
		IsSupported:                 false,
	}
}

func ChainParamsSupported(chainID int64) *types.ChainParams {
	cp := ChainParams(chainID)
	cp.IsSupported = true
	return cp
}

func ChainParamsList() (cpl types.ChainParamsList) {
	chainList := chains.ChainListByNetworkType(chains.NetworkType_privnet)

	for _, chain := range chainList {
		cpl.ChainParams = append(cpl.ChainParams, ChainParams(chain.ChainId))
	}
	return
}

func Tss() types.TSS {
	_, pubKey, _ := testdata.KeyTestPubAddr()
	spk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
	if err != nil {
		panic(err)
	}
	pk, err := zetacrypto.NewPubKey(spk)
	if err != nil {
		panic(err)
	}
	pubkeyString := pk.String()
	return types.TSS{
		TssPubkey:           pubkeyString,
		FinalizedZetaHeight: 1000,
		KeyGenZetaHeight:    1000,
	}
}

func TssList(n int) (tssList []types.TSS) {
	for i := 0; i < n; i++ {
		tss := Tss()
		tss.FinalizedZetaHeight = tss.FinalizedZetaHeight + int64(i)
		tss.KeyGenZetaHeight = tss.KeyGenZetaHeight + int64(i)
		tssList = append(tssList, tss)
	}
	return
}

func TssFundsMigrator(chainID int64) types.TssFundMigratorInfo {
	return types.TssFundMigratorInfo{
		ChainId:            chainID,
		MigrationCctxIndex: "sampleIndex",
	}
}

func BlameRecord(t *testing.T, index string) types.Blame {
	r := newRandFromStringSeed(t, index)
	return types.Blame{
		Index:         fmt.Sprintf("%d-%s", r.Int63(), index),
		FailureReason: "sample failure reason",
		Nodes:         nil,
	}
}
func BlameRecordsList(t *testing.T, n int) []types.Blame {
	blameList := make([]types.Blame, n)
	for i := 0; i < n; i++ {
		blameList[i] = BlameRecord(t, fmt.Sprintf("%d", i))
	}
	return blameList
}

func ChainNonces(t *testing.T, index string) types.ChainNonces {
	r := newRandFromStringSeed(t, index)
	return types.ChainNonces{
		Creator:         AccAddress(),
		Index:           index,
		ChainId:         r.Int63(),
		Nonce:           r.Uint64(),
		Signers:         []string{AccAddress(), AccAddress()},
		FinalizedHeight: r.Uint64(),
	}
}

func ChainNoncesList(t *testing.T, n int) []types.ChainNonces {
	chainNoncesList := make([]types.ChainNonces, n)
	for i := 0; i < n; i++ {
		chainNoncesList[i] = ChainNonces(t, fmt.Sprintf("%d", i))
	}
	return chainNoncesList
}

func PendingNoncesList(t *testing.T, index string, count int) []types.PendingNonces {
	r := newRandFromStringSeed(t, index)
	nonceLow := r.Int63()
	list := make([]types.PendingNonces, count)
	for i := 0; i < count; i++ {
		list[i] = types.PendingNonces{
			ChainId:   int64(i),
			NonceLow:  nonceLow,
			NonceHigh: nonceLow + r.Int63(),
			Tss:       StringRandom(r, 32),
		}
	}
	return list
}

func NonceToCctxList(t *testing.T, index string, count int) []types.NonceToCctx {
	r := newRandFromStringSeed(t, index)
	list := make([]types.NonceToCctx, count)
	for i := 0; i < count; i++ {
		list[i] = types.NonceToCctx{
			ChainId:   int64(i),
			Nonce:     r.Int63(),
			CctxIndex: StringRandom(r, 32),
		}
	}
	return list
}

func LegacyObserverMapper(t *testing.T, index string, observerList []string) *types.ObserverMapper {
	r := newRandFromStringSeed(t, index)

	return &types.ObserverMapper{
		Index:         index,
		ObserverChain: Chain(r.Int63()),
		ObserverList:  observerList,
	}
}

func LegacyObserverMapperList(t *testing.T, n int, index string) []*types.ObserverMapper {
	r := newRandFromStringSeed(t, index)
	observerList := []string{AccAddress(), AccAddress()}
	observerMapperList := make([]*types.ObserverMapper, n)
	for i := 0; i < n; i++ {
		observerMapperList[i] = LegacyObserverMapper(t, fmt.Sprintf("%d-%s", r.Int63(), index), observerList)
	}
	return observerMapperList
}

func BallotList(n int, observerSet []string) []types.Ballot {
	r := newRandFromSeed(int64(n))
	ballotList := make([]types.Ballot, n)

	for i := 0; i < n; i++ {
		identifier := crypto.Keccak256Hash([]byte(fmt.Sprintf("%d-%d-%d", r.Int63(), r.Int63(), r.Int63())))
		ballotList[i] = types.Ballot{
			Index:                identifier.Hex(),
			BallotIdentifier:     identifier.Hex(),
			VoterList:            observerSet,
			Votes:                VotesSuccessOnly(len(observerSet)),
			ObservationType:      types.ObservationType_InBoundTx,
			BallotThreshold:      sdk.OneDec(),
			BallotStatus:         types.BallotStatus_BallotFinalized_SuccessObservation,
			BallotCreationHeight: 0,
		}
	}
	return ballotList
}

func VotesSuccessOnly(voteCount int) []types.VoteType {
	votes := make([]types.VoteType, voteCount)
	for i := 0; i < voteCount; i++ {
		votes[i] = types.VoteType_SuccessObservation
	}
	return votes
}

func NonceToCCTX(t *testing.T, seed string) types.NonceToCctx {
	r := newRandFromStringSeed(t, seed)
	return types.NonceToCctx{
		ChainId:   r.Int63(),
		Nonce:     r.Int63(),
		CctxIndex: StringRandom(r, 64),
		Tss:       Tss().TssPubkey,
	}
}
