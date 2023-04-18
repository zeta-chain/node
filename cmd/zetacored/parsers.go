package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type ObserverInfoReader struct {
	ObserverAddress           string   `json:"ObserverAddress"`
	ZetaClientGranteeAddress  string   `json:"ZetaClientGranteeAddress"`
	StakingGranteeAddress     string   `json:"StakingGranteeAddress,omitempty"`
	StakingMaxTokens          string   `json:"StakingMaxTokens,omitempty"`
	StakingValidatorAllowList []string `json:"StakingValidatorAllowList,omitempty"`
	SpendGranteeAddress       string   `json:"SpendGranteeAddress,omitempty"`
	SpendMaxTokens            string   `json:"SpendMaxTokens,omitempty"`
	GovGranteeAddress         string   `json:"GovGranteeAddress,omitempty"`
	ZetaClientGranteePubKey   string   `json:"ZetaClientGranteePubKey"`
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
