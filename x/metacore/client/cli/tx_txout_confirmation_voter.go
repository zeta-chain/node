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
		Use:   "txout-confirmation-voter [txoutId] [txHash] [mMint] [destinationAsset] [destinationAmount] [blockHeight]",
		Short: "Broadcast message TxoutConfirmationVoter",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsTxoutId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argsTxHash := string(args[1])
			argsMMint,err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			argsDestinationAsset := string(args[3])
			argsDestinationAmount, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}
			argsBlockHeight, err := cast.ToUint64E(args[5])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgTxoutConfirmationVoter(clientCtx.GetFromAddress().String(), argsTxoutId, string(argsTxHash), argsMMint, string(argsDestinationAsset), (argsDestinationAmount), (argsBlockHeight))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
