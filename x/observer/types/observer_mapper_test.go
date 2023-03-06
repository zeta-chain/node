//go:build PRIVNET
// +build PRIVNET

package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
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
	createObserverList(file)
	obsListReadFromFile, err := ParsefileToObserverDetails(file)
	assert.NoError(t, err)
	for _, obs := range obsListReadFromFile {
		assert.Equal(t, obs.ObserverGranteeAddress, sdk.AccAddress(crypto.AddressHash([]byte("ObserverGranteeAddress"))).String())
	}
}

func createObserverList(fp string) {
	//list = append(append(append(list, CreateObserverMapperList(1, common.GoerliLocalNetChain())...),
	//	CreateObserverMapperList(1, common.BtcRegtestChain())...),
	//	CreateObserverMapperList(1, common.ZetaChain())...)
	var listReader []ObserverInfoReader
	listChainID := []int64{common.GoerliLocalNetChain().ChainId, common.BtcRegtestChain().ChainId, common.ZetaChain().ChainId}
	commonGrantAddress := sdk.AccAddress(crypto.AddressHash([]byte("ObserverGranteeAddress")))
	observerAddress := sdk.AccAddress(crypto.AddressHash([]byte("ObserverAddress")))
	info := ObserverInfoReader{
		SupportedChainsList:    listChainID,
		ObserverAddress:        observerAddress.String(),
		ObserverGranteeAddress: commonGrantAddress.String(),
	}
	listReader = append(listReader, info)

	file, _ := json.MarshalIndent(listReader, "", " ")
	_ = ioutil.WriteFile(fp, file, 0600)
}
