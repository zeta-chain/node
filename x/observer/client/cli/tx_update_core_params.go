package cli

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func CmdUpdateCoreParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-client-params [chain-id] [client-params.json]",
		Short: "Broadcast message updateClientParams",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			argCoreParams := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			var clientParams types.CoreParams
			file, err := filepath.Abs(argCoreParams)
			if err != nil {
				return err
			}
			file = filepath.Clean(file)
			input, err := os.ReadFile(file) // #nosec G304
			if err != nil {
				return err
			}
			err = json.Unmarshal(input, &clientParams)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateCoreParams(
				clientCtx.GetFromAddress().String(),
				&clientParams,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
