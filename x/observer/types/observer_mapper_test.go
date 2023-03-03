//go:build PRIVNET
// +build PRIVNET

package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/common"
	"io/ioutil"
	"os"
	"testing"
)

func TestParsefileToObserverMapper(t *testing.T) {
	file := "tmp.json"
	defer func(t *testing.T, fp string) {
		err := os.RemoveAll(fp)
		assert.NoError(t, err)
	}(t, file)
	expectedList := createObserverList(file)
	obsListReadFromFile, err := ParsefileToObserverMapper(file)
	assert.NoError(t, err)
	assert.Equal(t, expectedList, obsListReadFromFile)
}

func createObserverList(fp string) (list []*ObserverMapper) {
	list = append(append(append(list, CreateObserverMapperList(1, common.GoerliLocalNetChain())...),
		CreateObserverMapperList(1, common.BtcRegtestChain())...),
		CreateObserverMapperList(1, common.ZetaChain())...)
	listReader := make([]ObserverMapperReader, len(list))
	for i, mapper := range list {
		listReader[i] = ObserverMapperReader{
			Index:             mapper.Index,
			ObserverChainName: mapper.ObserverChain.ChainName.String(),
			ObserverChainID:   mapper.ObserverChain.ChainId,
			ObserverList:      mapper.ObserverList,
		}
	}
	file, _ := json.MarshalIndent(listReader, "", " ")
	_ = ioutil.WriteFile(fp, file, 0600)
	return list
}
