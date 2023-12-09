package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	sdkmath "cosmossdk.io/math"
)

func CmdUpdateZRC20WithdrawFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-zrc20-withdraw-fee [contractAddress] [newWithdrawFee] [newGasLimit]",
		Short: "Broadcast message UpdateZRC20WithdrawFee",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			contractAddress := args[0]

			newWithdrawFee := sdkmath.NewUintFromString(args[1])

			newGasLimit := sdkmath.NewUintFromString(args[2])
			
			msg := types.NewMsgUpdateZRC20WithdrawFee(
				clientCtx.GetFromAddress().String(),
				contractAddress,
				newWithdrawFee,
				newGasLimit,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
