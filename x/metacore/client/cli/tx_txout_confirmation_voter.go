package cli

import (
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"strconv"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

var _ = strconv.Itoa(0)

func CmdTxoutConfirmationVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txout-confirmation-voter [txoutId] [txHash] [mMint] [destinationAsset] [destinationAmount] [toAddress] [blockHeight]",
		Short: "Broadcast message TxoutConfirmationVoter",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsTxoutId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argsTxHash := args[1]
			argsMMint,err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			argsDestinationAsset := args[3]
			argsDestinationAmount, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}
			toAddress := args[5]

			argsBlockHeight, err := cast.ToUint64E(args[6])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgTxoutConfirmationVoter(clientCtx.GetFromAddress().String(), argsTxoutId, argsTxHash, argsMMint, argsDestinationAsset, argsDestinationAmount, toAddress, argsBlockHeight)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
