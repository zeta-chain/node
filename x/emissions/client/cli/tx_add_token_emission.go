package cli

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

var _ = strconv.Itoa(0)

func CmdAddTokenEmission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-token-emission [category] [amount]",
		Short: "Broadcast message add_token_emission",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argCategory := types.ParseStringToEmissionCategory(args[0])
			argAmount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("Unable to parse INT")
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddTokenEmission(
				clientCtx.GetFromAddress().String(),
				argCategory,
				argAmount,
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
