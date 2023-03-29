package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

var _ = strconv.Itoa(0)

func CmdUpdateClientParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-client-params [chain-id] [client-params.json]",
		Short: "Broadcast message updateClientParams",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChainId := args[0]
			argClientParams := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			chainid, err := strconv.ParseInt(argChainId, 10, 64)
			if err != nil {
				return err
			}
			var clientParams types.ClientParams
			file, err := filepath.Abs(argClientParams)
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

			msg := types.NewMsgUpdateClientParams(
				clientCtx.GetFromAddress().String(),
				chainid,
				&clientParams,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
