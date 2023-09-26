package sample

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
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

func ObserverMapper(t *testing.T, index string) *types.ObserverMapper {
	r := newRandFromStringSeed(t, index)

	return &types.ObserverMapper{
		Index:         index,
		ObserverChain: Chain(r.Int63()),
		ObserverList:  []string{AccAddress(), AccAddress()},
	}
}

func NodeAccount() *types.NodeAccount {
	operator := AccAddress()

	return &types.NodeAccount{
		Operator:       operator,
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

func CoreParams(chainID int64) *types.CoreParams {
	r := newRandFromSeed(chainID)

	return &types.CoreParams{
		ChainId:                     chainID,
		ConfirmationCount:           r.Uint64(),
		GasPriceTicker:              r.Uint64(),
		InTxTicker:                  r.Uint64(),
		OutTxTicker:                 r.Uint64(),
		WatchUtxoTicker:             r.Uint64(),
		ZetaTokenContractAddress:    EthAddress().String(),
		ConnectorContractAddress:    EthAddress().String(),
		Erc20CustodyContractAddress: EthAddress().String(),
		OutboundTxScheduleInterval:  r.Int63(),
		OutboundTxScheduleLookahead: r.Int63(),
	}
}

func CoreParamsList() (cpl types.CoreParamsList) {
	chainList := common.DefaultChainsList()

	for _, chain := range chainList {
		cpl.CoreParams = append(cpl.CoreParams, CoreParams(chain.ChainId))
	}
	return
}
