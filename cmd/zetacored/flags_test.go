package main_test

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	zetacore "github.com/zeta-chain/zetacore/cmd/zetacored"
)

// alwaysErrorValue allows to test f.Value.Set failure
type alwaysErrorValue struct{}

func (a *alwaysErrorValue) Set(string) error { return errors.New("error") }
func (a *alwaysErrorValue) String() string   { return "" }
func (a *alwaysErrorValue) Type() string     { return "string" }

func TestReplaceFlag(t *testing.T) {
	// Setting up a mock command structure
	rootCmd := &cobra.Command{Use: "app"}

	fooCmd := &cobra.Command{Use: "foo"}
	barCmd := &cobra.Command{Use: "bar"}

	barCmd.Flags().String("baz", "old", "Bar")
	barCmd.Flags().Var(&alwaysErrorValue{}, "error", "Always fails to set")

	fooCmd.AddCommand(barCmd)
	rootCmd.AddCommand(fooCmd)

	tests := []struct {
		name            string
		cmd             *cobra.Command
		subCommand      []string
		flagName        string
		newDefaultValue string
		wantErr         bool
		expectedValue   string
	}{
		{
			name:            "Replace valid flag",
			cmd:             rootCmd,
			subCommand:      []string{"foo", "bar"},
			flagName:        "baz",
			newDefaultValue: "new",
			wantErr:         false,
			expectedValue:   "new",
		},
		{
			name:            "Sub-command not found",
			cmd:             rootCmd,
			subCommand:      []string{"key", "nonexistent"},
			flagName:        "baz",
			newDefaultValue: "new",
			wantErr:         true,
		},
		{
			name:            "Flag not found",
			cmd:             rootCmd,
			subCommand:      []string{"foo", "bar"},
			flagName:        "nonexistent",
			newDefaultValue: "new",
			wantErr:         true,
		},
		{
			name:            "Flag value cannot be set",
			cmd:             rootCmd,
			subCommand:      []string{"foo", "bar"},
			flagName:        "error",
			newDefaultValue: "new",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := zetacore.ReplaceFlag(tt.cmd, tt.subCommand, tt.flagName, tt.newDefaultValue)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Check if the value was replaced correctly
				c, _, _ := tt.cmd.Find(tt.subCommand)
				f := c.Flags().Lookup(tt.flagName)
				assert.Equal(t, tt.expectedValue, f.DefValue)
			}
		})
	}
}
