package sample

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/cosmos"
	zetacrypto "github.com/zeta-chain/node/pkg/crypto"
	"github.com/zeta-chain/node/pkg/ptr"
	"github.com/zeta-chain/node/x/observer/types"
)

func Ballot(t *testing.T, index string) *types.Ballot {
	r := newRandFromStringSeed(t, index)

	return &types.Ballot{
		BallotIdentifier:     index,
		VoterList:            []string{AccAddress(), AccAddress()},
		Votes:                []types.VoteType{types.VoteType_FailureObservation, types.VoteType_SuccessObservation},
		ObservationType:      types.ObservationType_EmptyObserverType,
		BallotThreshold:      sdkmath.LegacyNewDec(1),
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

func KeygenFromRand(r *rand.Rand) types.Keygen {
	pubkey := PubKey(r)
	return types.Keygen{
		Status:         types.KeygenStatus_KeyGenSuccess,
		GranteePubkeys: []string{pubkey.String()},
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

	fiftyPercent, err := sdkmath.LegacyNewDecFromStr("0.5")
	if err != nil {
		return nil
	}

	confirmationParams := ConfirmationParams(r)

	return &types.ChainParams{
		ChainId:                     chainID,
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
		MinObserverDelegation:       sdkmath.LegacyNewDec(r.Int63()),
		IsSupported:                 false,
		GatewayAddress:              EthAddress().String(),
		ConfirmationParams:          &confirmationParams,
	}
}

func ChainParamsFromRand(r *rand.Rand, chainID int64) *types.ChainParams {
	fiftyPercent := sdkmath.LegacyMustNewDecFromStr("0.5")
	return &types.ChainParams{
		ChainId:                     chainID,
		GasPriceTicker:              Uint64InRangeFromRand(r, 1, 300),
		InboundTicker:               Uint64InRangeFromRand(r, 1, 300),
		OutboundTicker:              Uint64InRangeFromRand(r, 1, 300),
		WatchUtxoTicker:             Uint64InRangeFromRand(r, 1, 300),
		ZetaTokenContractAddress:    EthAddressFromRand(r).String(),
		ConnectorContractAddress:    EthAddressFromRand(r).String(),
		Erc20CustodyContractAddress: EthAddressFromRand(r).String(),
		OutboundScheduleInterval:    Int64InRangeFromRand(r, 1, 100),
		OutboundScheduleLookahead:   Int64InRangeFromRand(r, 1, 500),
		BallotThreshold:             fiftyPercent,
		MinObserverDelegation:       sdkmath.LegacyNewDec(r.Int63()),
		IsSupported:                 true,
	}
}

func ChainParamsSupported(chainID int64) *types.ChainParams {
	cp := ChainParams(chainID)
	cp.IsSupported = true
	return cp
}

func ChainParamsList() (cpl types.ChainParamsList) {
	chainList := chains.ChainListByNetworkType(chains.NetworkType_privnet, []chains.Chain{})

	for _, chain := range chainList {
		cpl.ChainParams = append(cpl.ChainParams, ChainParams(chain.ChainId))
	}
	return
}

// TSSFromRand returns a random TSS,it uses the randomness provided as a parameter
func TSSFromRand(r *rand.Rand) (types.TSS, error) {
	pubKey := PubKey(r)
	spk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
	if err != nil {
		return types.TSS{}, err
	}
	pk, err := zetacrypto.NewPubKey(spk)
	if err != nil {
		return types.TSS{}, err
	}
	pubkeyString := pk.String()
	return types.TSS{
		TssPubkey:           pubkeyString,
		TssParticipantList:  []string{},
		OperatorAddressList: []string{},
		FinalizedZetaHeight: r.Int63(),
		KeyGenZetaHeight:    r.Int63(),
	}, nil
}

// TODO: rename to TSS
// https://github.com/zeta-chain/node/issues/3098
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

func ChainNonces(chainID int64) types.ChainNonces {
	r := newRandFromSeed(chainID)
	return types.ChainNonces{
		Creator:         AccAddress(),
		ChainId:         chainID,
		Nonce:           r.Uint64(),
		Signers:         []string{AccAddress(), AccAddress()},
		FinalizedHeight: r.Uint64(),
	}
}

func ChainNoncesList(n int) []types.ChainNonces {
	chainNoncesList := make([]types.ChainNonces, n)
	for i := 0; i < n; i++ {
		chainNoncesList[i] = ChainNonces(int64(i))
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

func BallotList(n int, observerSet []string) []types.Ballot {
	r := newRandFromSeed(int64(n))
	ballotList := make([]types.Ballot, n)

	for i := 0; i < n; i++ {
		identifier := crypto.Keccak256Hash(fmt.Appendf(nil, "%d-%d-%d", r.Int63(), r.Int63(), r.Int63()))
		ballotList[i] = types.Ballot{
			BallotIdentifier:     identifier.Hex(),
			VoterList:            observerSet,
			Votes:                VotesSuccessOnly(len(observerSet)),
			ObservationType:      types.ObservationType_InboundTx,
			BallotThreshold:      sdkmath.LegacyOneDec(),
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

func GasPriceIncreaseFlags() types.GasPriceIncreaseFlags {
	return types.GasPriceIncreaseFlags{
		EpochLength:             1,
		RetryInterval:           1,
		GasPriceIncreasePercent: 1,
		MaxPendingCctxs:         100,
		RetryIntervalBTC:        2,
	}
}

func GasPriceIncreaseFlagsFromRand(r *rand.Rand) types.GasPriceIncreaseFlags {
	minValue := 1
	maxValue := 100
	return types.GasPriceIncreaseFlags{
		EpochLength:             int64(r.Intn(maxValue-minValue) + minValue),
		RetryInterval:           time.Duration(r.Intn(maxValue-minValue) + minValue),
		GasPriceIncreasePercent: 1,
		MaxPendingCctxs:         100,
		RetryIntervalBTC:        time.Duration(r.Intn(maxValue-minValue) + minValue),
	}
}

func OperationalFlags() types.OperationalFlags {
	return types.OperationalFlags{
		RestartHeight:         1,
		SignerBlockTimeOffset: ptr.Ptr(time.Second),
	}
}

func ConfirmationParams(r *rand.Rand) types.ConfirmationParams {
	randInboundCount := Uint64InRangeFromRand(r, 2, 200)
	randOutboundCount := Uint64InRangeFromRand(r, 2, 200)

	return types.ConfirmationParams{
		SafeInboundCount: randInboundCount,
		// enabled fast inbound confirmation count should be less than safe count
		FastInboundCount:  Uint64InRange(1, randInboundCount-1),
		SafeOutboundCount: randOutboundCount,
		// enabled fast outbound confirmation count should be less than safe count
		FastOutboundCount: Uint64InRange(1, randOutboundCount-1),
	}
}
