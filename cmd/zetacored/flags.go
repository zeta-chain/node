package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var KeyAddCommand = []string{"key", "add"}

const (
	HDPathFlag     = "hd-path"
	HDPathEthereum = "m/44'/60'/0'/0/0"
)

// SetEthereumHDPath sets the default HD path to Ethereum's
func SetEthereumHDPath(cmd *cobra.Command) error {
	return ReplaceFlag(cmd, KeyAddCommand, HDPathFlag, HDPathEthereum)
}

// ReplaceFlag replaces the default value of a flag of a sub-command
func ReplaceFlag(cmd *cobra.Command, subCommand []string, flagName, newDefaultValue string) error {
	// Find the sub-command
	c, _, err := cmd.Find(subCommand)
	if err != nil {
		return fmt.Errorf("failed to find %v sub-command: %v", subCommand, err)
	}

	// Get the flag from the sub-command
	f := c.Flags().Lookup(flagName)
	if f == nil {
		return fmt.Errorf("%s flag not found in %v sub-command", flagName, subCommand)
	}

	// Set the default value for the flag
	f.DefValue = newDefaultValue
	if err := f.Value.Set(newDefaultValue); err != nil {
		return fmt.Errorf("failed to set the value of %s flag: %v", flagName, err)
	}

	return nil
}
