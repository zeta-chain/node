package local

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/e2e/config"
)

// GetConfig returns config from file from the command line flag
func GetConfig(cmd *cobra.Command) (config.Config, error) {
	configFile, err := cmd.Flags().GetString(FlagConfigFile)
	if err != nil {
		return config.Config{}, fmt.Errorf("--config is a required parameter")
	}

	configFile, err = filepath.Abs(configFile)
	if err != nil {
		return config.Config{}, err
	}

	return config.ReadConfig(configFile, true)
}

// OverwriteAccountData overwrites the account data in the config with the one from the file specified in the command line flag
// used for upgrade tests in case some accounts should be overwritten in zetae2e upgrade handler
func OverwriteAccountData(cmd *cobra.Command, conf *config.Config) error {
	configFile, err := cmd.Flags().GetString(flagAccountConfig)
	if err != nil || configFile == "" {
		return fmt.Errorf("--account-config is a required parameter to override account data")
	}
	configFile, err = filepath.Abs(configFile)
	if err != nil {
		return err
	}

	accountData, err := config.ReadConfig(configFile, false)
	if err != nil {
		return err
	}
	conf.DefaultAccount = accountData.DefaultAccount
	conf.AdditionalAccounts = accountData.AdditionalAccounts
	conf.PolicyAccounts = accountData.PolicyAccounts
	conf.ObserverRelayerAccounts = accountData.ObserverRelayerAccounts

	return nil
}
