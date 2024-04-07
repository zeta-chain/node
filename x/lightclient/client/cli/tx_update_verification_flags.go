package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

func CmdUpdateVerificationFlags() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-verification-flags [eth-type-chain-enabled] [btc-type-chain-enabled]",
		Short: "Update verification flags",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			argEthEnabled, err := strconv.ParseBool(args[0])
			if err != nil {
				return err
			}
			arsBtcEnabled, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgUpdateVerificationFlags(clientCtx.GetFromAddress().String(), argEthEnabled, arsBtcEnabled)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
