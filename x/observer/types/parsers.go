package types

import (
	"encoding/json"
	"github.com/zeta-chain/zetacore/common"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type ObserverMapperReader struct {
	Index           string   `json:"index"`
	ObserverChain   string   `json:"observerChain"`
	ObservationType string   `json:"observationType"`
	ObserverList    []string `json:"observerList"`
}

func ParsefileToObserverMapper(fp string) ([]*ObserverMapper, error) {
	var observers []ObserverMapperReader
	file, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}
	input, err := ioutil.ReadFile(file) // #nosec G402 G304
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(input, &observers)
	if err != nil {
		return nil, err
	}

	observerMappers := make([]*ObserverMapper, len(observers))
	for i, readerValue := range observers {
		observerMappers[i] = &ObserverMapper{
			Index:           readerValue.Index,
			ObserverChain:   ParseStringToObserverChain(readerValue.ObserverChain),
			ObservationType: ParseStringToObservationType(readerValue.ObservationType),
			ObserverList:    readerValue.ObserverList,
		}
	}
	return observerMappers, nil
}

func ParseCommonChaintoObservationChain(chain string) ObserverChain {
	commonChain := common.Chain(chain)
	switch commonChain {
	// Mainnet Chains
	case common.ZETAChain, common.Chain(strings.ToUpper(string(common.ZETAChain))):
		return ObserverChain_ZetaChain
	case common.ETHChain, common.Chain(strings.ToUpper(string(common.ETHChain))):
		return ObserverChain_Eth
	case common.BSCChain, common.Chain(strings.ToUpper(string(common.BSCChain))):
		return ObserverChain_BscMainnet
	case common.POLYGONChain, common.Chain(strings.ToUpper(string(common.POLYGONChain))):
		return ObserverChain_Polygon
	// Testnet Chains
	case common.MumbaiChain, common.Chain(strings.ToUpper(string(common.MumbaiChain))):
		return ObserverChain_Mumbai
	case common.BaobabChain, common.Chain(strings.ToUpper(string(common.BaobabChain))):
		return ObserverChain_Baobab

	case common.GoerliChain, common.Chain(strings.ToUpper(string(common.GoerliChain))):
		return ObserverChain_Goerli
	case common.BSCTestnetChain, common.Chain(strings.ToUpper(string(common.BSCTestnetChain))):
		return ObserverChain_BscTestnet
	case common.BTCTestnetChain, common.Chain(strings.ToUpper(string(common.BTCTestnetChain))):
		return ObserverChain_BtcTestnet
	}

	return ObserverChain_Empty
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

func ParseStringToObserverChain(chain string) ObserverChain {
	c := ObserverChain_value[chain]
	return ObserverChain(c)
}

func ParseStringToObservationType(observationType string) ObservationType {
	c := ObservationType_value[observationType]
	return ObservationType(c)
}
