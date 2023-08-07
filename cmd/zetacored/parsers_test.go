//go:build PRIVNET
// +build PRIVNET

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
	"github.com/zeta-chain/zetacore/app"

	//"os"
	"testing"
)

func TestParsefileToObserverMapper(t *testing.T) {
	file := "tmp.json"
	defer func(t *testing.T, fp string) {
		err := os.RemoveAll(fp)
		assert.NoError(t, err)
	}(t, file)
	app.SetConfig()
	createObserverList(file)
	obsListReadFromFile, err := ParsefileToObserverDetails(file)
	assert.NoError(t, err)
	for _, obs := range obsListReadFromFile {
		assert.Equal(t, obs.ZetaClientGranteeAddress, sdk.AccAddress(crypto.AddressHash([]byte("ObserverGranteeAddress"))).String())
	}
}

func createObserverList(fp string) {
	var listReader []ObserverInfoReader
	//listChainID := []int64{common.GoerliLocalNetChain().ChainId, common.BtcRegtestChain().ChainId, common.ZetaChain().ChainId}
	commonGrantAddress := sdk.AccAddress(crypto.AddressHash([]byte("ObserverGranteeAddress")))
	observerAddress := sdk.AccAddress(crypto.AddressHash([]byte("ObserverAddress")))
	validatorAddress := sdk.ValAddress(crypto.AddressHash([]byte("ValidatorAddress")))
	info := ObserverInfoReader{
		ObserverAddress:           observerAddress.String(),
		ZetaClientGranteeAddress:  commonGrantAddress.String(),
		StakingGranteeAddress:     commonGrantAddress.String(),
		StakingMaxTokens:          "100000000",
		StakingValidatorAllowList: []string{validatorAddress.String()},
		SpendMaxTokens:            "100000000",
		GovGranteeAddress:         commonGrantAddress.String(),
		ZetaClientGranteePubKey:   "zetapub1addwnpepqggtjvkmj6apcqr6ynyc5edxf2mpf5fxp2d3kwupemxtfwvg6gm7qv79fw0",
	}
	listReader = append(listReader, info)

	file, _ := json.MarshalIndent(listReader, "", " ")
	_ = ioutil.WriteFile(fp, file, 0600)
}
