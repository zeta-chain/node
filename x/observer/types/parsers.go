package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"io/ioutil"
	"path/filepath"
)

type ObserverMapperReader struct {
	Index             string   `json:"index"`
	ObserverChainName string   `json:"observerChainName"`
	ObserverChainID   int64    `json:"observerChainId"`
	ObservationType   string   `json:"observationType"`
	ObserverList      []string `json:"observerList"`
}

func ParsefileToObserverMapper(fp string) ([]*ObserverMapper, error) {
	var observers []ObserverMapperReader
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

	observerMappers := make([]*ObserverMapper, len(observers))
	for i, readerValue := range observers {
		chain := &common.Chain{
			ChainName: common.ParseStringToObserverChain(readerValue.ObserverChainName),
			ChainId:   readerValue.ObserverChainID,
		}
		observationType := ParseStringToObservationType(readerValue.ObservationType)
		if observationType == 0 || chain.ChainName == 0 {
			return nil, errors.Wrap(ErrUnableToParseMapper, fmt.Sprintf("Chain %s | ObserVation %s", readerValue.ObserverChainName, readerValue.ObservationType))
		}
		observerMappers[i] = &ObserverMapper{
			Index:           readerValue.Index,
			ObserverChain:   chain,
			ObservationType: ParseStringToObservationType(readerValue.ObservationType),
			ObserverList:    readerValue.ObserverList,
		}
	}
	return observerMappers, nil
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
