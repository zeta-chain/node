package cli

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/zetacore/x/authority/types"
)

func CmdUpdatePolices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-policies [policies-json-file]",
		Short: "Update the policies",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			policies, err := readPoliciesFromFile(args[0])
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

// readPoliciesFromFile read the policies from the file using os package and unmarshal it into the policies variable
func readPoliciesFromFile(filePath string) (types.Policies, error) {
	var policies types.Policies
	policiesBytes, err := os.ReadFile(filePath)
	if err != nil {
		return policies, err
	}

	err = policies.Unmarshal(policiesBytes)
	return policies, err
}
