package parsers

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/zeta-chain/node/pkg/sdkconfig/default"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestParsefileToObserverMapper(t *testing.T) {
	file := "tmp.json"
	defer func(t *testing.T, fp string) {
		err := os.RemoveAll(fp)
		require.NoError(t, err)
	}(t, file)

	observerAddress := sample.AccAddress()
	commonGrantAddress := sample.AccAddress()
	validatorAddress := sample.AccAddress()

	createObserverList(file, observerAddress, commonGrantAddress, validatorAddress)
	obsListReadFromFile, err := ParsefileToObserverDetails(file)
	require.NoError(t, err)
	for _, obs := range obsListReadFromFile {
		require.Equal(
			t,
			obs.ObserverAddress,
			observerAddress,
		)
		require.Equal(
			t,
			obs.ZetaClientGranteeAddress,
			commonGrantAddress,
		)
	}
}

func createObserverList(fp string, observerAddress, commonGrantAddress, validatorAddress string) {
	var listReader []ObserverInfoReader
	info := ObserverInfoReader{
		ObserverAddress:           observerAddress,
		ZetaClientGranteeAddress:  commonGrantAddress,
		StakingGranteeAddress:     commonGrantAddress,
		StakingMaxTokens:          "100000000",
		StakingValidatorAllowList: []string{validatorAddress},
		SpendMaxTokens:            "100000000",
		GovGranteeAddress:         commonGrantAddress,
		ZetaClientGranteePubKey:   "zetapub1addwnpepqggtjvkmj6apcqr6ynyc5edxf2mpf5fxp2d3kwupemxtfwvg6gm7qv79fw0",
	}
	listReader = append(listReader, info)

	file, _ := json.MarshalIndent(listReader, "", " ")
	_ = os.WriteFile(fp, file, 0600)
}
