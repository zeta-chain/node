package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestParsefileToObserverMapper(t *testing.T) {
	file := "tmp.json"
	defer func(t *testing.T, fp string) {
		err := os.Remove(fp)
		assert.NoError(t, err)
	}(t, file)
	expectedList := createObserverList(file)
	obsListReadFromFile, err := ParsefileToObserverMapper(file)
	assert.NoError(t, err)
	assert.Equal(t, expectedList, obsListReadFromFile)
}

func createObserverList(fp string) (list []*ObserverMapper) {
	list = append(append(append(list, CreateObserverMapperList(1, ObserverChain_Eth, ObservationType_InBoundTx)...),
		CreateObserverMapperList(1, ObserverChain_Bsc, ObservationType_InBoundTx)...),
		CreateObserverMapperList(1, ObserverChain_Polygon, ObservationType_OutBoundTx)...)

	file, _ := json.MarshalIndent(list, "", " ")
	_ = ioutil.WriteFile(fp, file, 0600)
	return list
}
