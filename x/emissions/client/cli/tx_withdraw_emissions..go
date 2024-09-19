package cli

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/emissions/types"
)

func CmdWithdrawEmission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-emission [amount]",
		Short: "create a new withdrawEmission",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsAmount, ok := sdkmath.NewIntFromString(args[0])
			if !ok {
				return errors.New("invalid amount")
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgWithdrawEmissions(clientCtx.GetFromAddress().String(), argsAmount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
