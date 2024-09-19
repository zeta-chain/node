package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
)

func ConvertReceiveStatusToVoteType(status chains.ReceiveStatus) VoteType {
	switch status {
	case chains.ReceiveStatus_success:
		return VoteType_SuccessObservation
	case chains.ReceiveStatus_failed:
		return VoteType_FailureObservation
	default:
		return VoteType_NotYetVoted
	}
}

func ParseStringToObservationType(observationType string) ObservationType {
	c := ObservationType_value[observationType]
	return ObservationType(c)
}

func GetOperatorAddressFromAccAddress(accAddr string) (sdk.ValAddress, error) {
	accAddressBech32, err := sdk.AccAddressFromBech32(accAddr)
	if err != nil {
		return nil, err
	}
	valAddress := sdk.ValAddress(accAddressBech32)
	valAddressBech32, err := sdk.ValAddressFromBech32(valAddress.String())
	if err != nil {
		return nil, err
	}
	return valAddressBech32, nil
}

func GetAccAddressFromOperatorAddress(valAddress string) (sdk.AccAddress, error) {
	valAddressBech32, err := sdk.ValAddressFromBech32(valAddress)
	if err != nil {
		return nil, err
	}
	accAddress := sdk.AccAddress(valAddressBech32)
	accAddressBech32, err := sdk.AccAddressFromBech32(accAddress.String())
	if err != nil {
		return nil, err
	}
	return accAddressBech32, nil
}
