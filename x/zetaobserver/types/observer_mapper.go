package types

import (
	"encoding/json"
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
