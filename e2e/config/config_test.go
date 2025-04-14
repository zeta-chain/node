package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestConfigAdditionalAccountsSliceSynced ensures than if any new accounts are added the
// AccountsSlice() function is updated
func TestConfigAdditionalAccountsSliceSynced(t *testing.T) {
	additionalAccountsType := reflect.TypeOf(AdditionalAccounts{})
	numberOfAdditionalAccounts := additionalAccountsType.NumField()

	require.Greater(t, numberOfAdditionalAccounts, 0)

	conf := Config{}
	err := conf.GenerateKeys()
	require.NoError(t, err)

	additionalAccountsSlice := conf.AdditionalAccounts.AsSlice()
	require.Len(t, additionalAccountsSlice, numberOfAdditionalAccounts)

	for _, account := range additionalAccountsSlice {
		require.NoError(t, account.Validate())
	}
}

func TestConfigInvalidAccount(t *testing.T) {
	conf := Config{
		DefaultAccount: Account{
			RawEVMAddress: "asdf",
			RawPrivateKey: "asdf",
		},
	}
	err := conf.Validate()
	require.Error(t, err)
}
