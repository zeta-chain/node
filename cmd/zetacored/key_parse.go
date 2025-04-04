package main

import (
	"strings"

	"github.com/spf13/cobra"
	
	sdkkeys "github.com/cosmos/cosmos-sdk/client/keys"
)

func ParseKeyCommand() *cobra.Command {
	cmd := sdkkeys.ParseKeyStringCommand()
	origRunE := cmd.RunE
	
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 && strings.HasPrefix(args[0], "0x") {
			args[0] = args[0][2:]
		}
		
		return origRunE(cmd, args)
	}
	
	return cmd
}
