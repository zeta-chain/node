package main

import (
	"encoding/json"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"io/ioutil"
	"path/filepath"
)

type ObserverInfoReader struct {
	SupportedChainsList       []int64                     `json:"SupportedChainsList"`
	ObserverAddress           string                      `json:"ObserverAddress"`
	ZetaClientGranteeAddress  string                      `json:"ZetaClientGranteeAddress"`
	TssSignerAddress          string                      `json:"TssSignerAddress"`
	StakingGranteeAddress     string                      `json:"StakingGranteeAddress"`
	StakingMaxTokens          string                      `json:"StakingMaxTokens"`
	StakingValidatorAllowList []string                    `json:"StakingValidatorAllowList"`
	SpendGranteeAddress       string                      `json:"SpendGranteeAddress"`
	SpendMaxTokens            string                      `json:"SpendMaxTokens"`
	GovGranteeAddress         string                      `json:"GovGranteeAddress"`
	NodeAccount               crosschaintypes.NodeAccount `json:"NodeAccount"`
}

func (o ObserverInfoReader) String() string {
	s, _ := json.MarshalIndent(o, "", "\t")
	return string(s)
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
