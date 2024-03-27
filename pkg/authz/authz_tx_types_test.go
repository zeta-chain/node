package authz

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTxTypeString(t *testing.T) {
	tests := []struct {
		name           string
		txType         TxType
		expectedString string
	}{
		{"InboundVoter", InboundVoter, "InboundVoter"},
		{"OutboundVoter", OutboundVoter, "OutboundVoter"},
		{"NonceVoter", NonceVoter, "NonceVoter"},
		{"GasPriceVoter", GasPriceVoter, "GasPriceVoter"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedString, test.txType.String())
		})
	}
}

func TestKeyTypeString(t *testing.T) {
	tests := []struct {
		name           string
		keyType        KeyType
		expectedString string
	}{
		{"TssSignerKey", TssSignerKey, "tss_signer"},
		{"ValidatorGranteeKey", ValidatorGranteeKey, "validator_grantee"},
		{"ZetaClientGranteeKey", ZetaClientGranteeKey, "zetaclient_grantee"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedString, test.keyType.String())
		})
	}
}

func TestGetAllKeyTypes(t *testing.T) {
	expectedKeys := []KeyType{ValidatorGranteeKey, ZetaClientGranteeKey, TssSignerKey}
	actualKeys := GetAllKeyTypes()

	require.ElementsMatch(t, expectedKeys, actualKeys)
}
