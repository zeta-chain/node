package cli

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/authority/types"
)

func CmdUpdatePolicies() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-policies [policies-json-file]",
		Short: "Update policies to values provided in the JSON file.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			policies, err := ReadPoliciesFromFile(os.DirFS("."), args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdatePolicies(
				clientCtx.GetFromAddress().String(),
				policies,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ReadPoliciesFromFile read the policies from the file using os package and unmarshal it into the policies variable
func ReadPoliciesFromFile(fsys fs.FS, filePath string) (types.Policies, error) {
	var policies types.Policies
	policiesBytes, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		return policies, fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(policiesBytes, &policies)
	return policies, err
}
