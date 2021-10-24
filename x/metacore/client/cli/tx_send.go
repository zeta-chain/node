package cli

import (
	"github.com/spf13/cobra"

	"github.com/spf13/cast"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

func CmdCreateSend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-send [index] [sender] [senderChain] [receiver] [receiverChain] [mBurnt] [mMint] [message] [inTxHash] [inBlockHeight] [outTxHash] [outBlockHeight]",
		Short: "Create a new Send",
		Args:  cobra.ExactArgs(12),
		RunE: func(cmd *cobra.Command, args []string) error {
			index := args[0]
			argsSender, err := cast.ToStringE(args[1])
			if err != nil {
				return err
			}
			argsSenderChain, err := cast.ToStringE(args[2])
			if err != nil {
				return err
			}
			argsReceiver, err := cast.ToStringE(args[3])
			if err != nil {
				return err
			}
			argsReceiverChain, err := cast.ToStringE(args[4])
			if err != nil {
				return err
			}
			argsMBurnt, err := cast.ToStringE(args[5])
			if err != nil {
				return err
			}
			argsMMint, err := cast.ToStringE(args[6])
			if err != nil {
				return err
			}
			argsMessage, err := cast.ToStringE(args[7])
			if err != nil {
				return err
			}
			argsInTxHash, err := cast.ToStringE(args[8])
			if err != nil {
				return err
			}
			argsInBlockHeight, err := cast.ToUint64E(args[9])
			if err != nil {
				return err
			}
			argsOutTxHash, err := cast.ToStringE(args[10])
			if err != nil {
				return err
			}
			argsOutBlockHeight, err := cast.ToUint64E(args[11])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateSend(clientCtx.GetFromAddress().String(), index, argsSender, argsSenderChain, argsReceiver, argsReceiverChain, argsMBurnt, argsMMint, argsMessage, argsInTxHash, argsInBlockHeight, argsOutTxHash, argsOutBlockHeight)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
