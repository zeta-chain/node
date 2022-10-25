package types

import (
	"encoding/json"
	"github.com/zeta-chain/zetacore/common"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ParsefileToObserverMapper(fp string) ([]*ObserverMapper, error) {
	var observers []*ObserverMapper
	file, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}
	input, err := ioutil.ReadFile(file) // #nosec G402
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(input, &observers)
	if err != nil {
		return nil, err
	}
	return observers, nil
}

func (*ObserverMapper) Validate() bool {
	return true
}

func VerifyObserverMapper(obs []*ObserverMapper) bool {
	for _, mapper := range obs {
		ok := mapper.Validate()
		if !ok {
			return ok
		}
	}
	return true
}

func CheckReceiveStatus(status common.ReceiveStatus) error {
	switch status {
	case common.ReceiveStatus_Success:
		return nil
	case common.ReceiveStatus_Failed:
		return nil
	default:
		return ErrInvalidStatus
	}
}

func ConvertStringChaintoObservationChain(chain string) ObserverChain {
	commonChain, err := common.ParseChain(chain)
	if err != nil {
		return ObserverChain_Empty
	}
	switch commonChain {
	case common.ETHChain, common.Chain(strings.ToUpper(string(common.ETHChain))):
		return ObserverChain_Eth
	case common.BSCChain, common.Chain(strings.ToUpper(string(common.BSCChain))):
		return ObserverChain_Bsc
	case common.POLYGONChain, common.Chain(strings.ToUpper(string(common.POLYGONChain))):
		return ObserverChain_Polygon
	case common.GoerliChain, common.Chain(strings.ToUpper(string(common.GoerliChain))):
		return ObserverChain_Goerli
	case common.BSCTestnetChain:
		return ObserverChain_BscTestnet
	case common.MumbaiChain:
		return ObserverChain_Mumbai
	case common.ZETAChain:
		return ObserverChain_ZetaChain
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
