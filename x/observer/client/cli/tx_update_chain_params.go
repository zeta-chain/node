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

func CmdUpdateChainParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-chain-params [chain-id] [client-params.json]",
		Short: "Broadcast message updateChainParams",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			argChainParams := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			var clientParams types.ChainParams
			file, err := filepath.Abs(argChainParams)
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

			msg := types.NewMsgUpdateChainParams(
				clientCtx.GetFromAddress().String(),
				&clientParams,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
