package parsers

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ObserverInfoReader struct {
	IsObserver                string   `json:"IsObserver"`
	ObserverAddress           string   `json:"ObserverAddress"`
	ZetaClientGranteeAddress  string   `json:"ZetaClientGranteeAddress,omitempty"`
	StakingGranteeAddress     string   `json:"StakingGranteeAddress,omitempty"`
	StakingMaxTokens          string   `json:"StakingMaxTokens,omitempty"`
	StakingValidatorAllowList []string `json:"StakingValidatorAllowList,omitempty"`
	SpendGranteeAddress       string   `json:"SpendGranteeAddress,omitempty"`
	SpendMaxTokens            string   `json:"SpendMaxTokens,omitempty"`
	GovGranteeAddress         string   `json:"GovGranteeAddress,omitempty"`
	ZetaClientGranteePubKey   string   `json:"ZetaClientGranteePubKey,omitempty"`
}

func (o ObserverInfoReader) String() string {
	s, err := json.MarshalIndent(o, "", "\t")
	if err != nil {
		return ""
	}
	return string(s)
}

func ParsefileToObserverDetails(fp string) ([]ObserverInfoReader, error) {
	var observers []ObserverInfoReader
	file, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}
	file = filepath.Clean(file)
	input, err := os.ReadFile(file) // #nosec G304
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(input, &observers)
	if err != nil {
		return nil, err
	}
	return observers, nil
}
