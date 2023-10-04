package cli

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func CmdUpdateZRC20LiquidityCap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-zrc20-liquidity-cap [zrc20] [liquidity-cap]",
		Short: "Broadcast message UpdateZRC20LiquidityCap",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			newCap := math.NewUintFromString(args[1])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fmt.Printf("CLI address: %s\n", clientCtx.GetFromAddress().String())
			msg := types.NewMsgUpdateZRC20LiquidityCap(
				clientCtx.GetFromAddress().String(),
				args[0],
				newCap,
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
