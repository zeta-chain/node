package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"io/ioutil"
	"path/filepath"
)

type ObserverInfoReader struct {
	SupportedChainsList       []int64  `json:"SupportedChainsList"`
	ObserverAddress           string   `json:"ObserverAddress"`
	ZetaClientGranteeAddress  string   `json:"ZetaClientGranteeAddress"`
	StakingGranteeAddress     string   `json:"StakingGranteeAddress"`
	StakingMaxTokens          string   `json:"StakingMaxTokens"`
	StakingValidatorAllowList []string `json:"StakingValidatorAllowList"`
	SpendGranteeAddress       string   `json:"SpendGranteeAddress"`
	SpendMaxTokens            string   `json:"SpendMaxTokens"`
	GovGranteeAddress         string   `json:"GovGranteeAddress"`
}

func ParsefileToObserverDetails(fp string) ([]ObserverInfoReader, error) {
	var observers []ObserverInfoReader
	file, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}
	file = filepath.Clean(file)
	input, err := ioutil.ReadFile(file) // #nosec G304
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(input, &observers)
	if err != nil {
		return nil, err
	}
	return observers, nil
}

func ConvertReceiveStatusToVoteType(status common.ReceiveStatus) VoteType {
	switch status {
	case common.ReceiveStatus_Success:
		return VoteType_SuccessObservation
	case common.ReceiveStatus_Failed:
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
