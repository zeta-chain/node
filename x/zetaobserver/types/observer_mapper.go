package types

import (
	"encoding/json"
	"github.com/zeta-chain/zetacore/common"
	"io/ioutil"
	"path/filepath"
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

func ConvertStringChaintoObservationChain(chain string) ObserverChain {
	commonChain := common.Chain(chain)
	switch commonChain {
	case common.ETHChain:
		return ObserverChain_EthChainObserver
	case common.BSCChain:
		return ObserverChain_BscChainObserver
	case common.POLYGONChain:
		return ObserverChain_PolygonChainObserver
	}
	return ObserverChain_EmptyObserver
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
