package cli

import (
	"github.com/spf13/cobra"

	"github.com/spf13/cast"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

func CmdCreateSendVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-send-voter [index] [sender] [senderChainId] [receiver] [receiverChainId] [mBurnt] [message] [txHash] [blockHeight]",
		Short: "Create a new SendVoter",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {
			index := args[0]
			argsSender, err := cast.ToStringE(args[1])
			if err != nil {
				return err
			}
			argsSenderChainId, err := cast.ToStringE(args[2])
			if err != nil {
				return err
			}
			argsReceiver, err := cast.ToStringE(args[3])
			if err != nil {
				return err
			}
			argsReceiverChainId, err := cast.ToStringE(args[4])
			if err != nil {
				return err
			}
			argsMBurnt, err := cast.ToStringE(args[5])
			if err != nil {
				return err
			}
			argsMessage, err := cast.ToStringE(args[6])
			if err != nil {
				return err
			}
			argsTxHash, err := cast.ToStringE(args[7])
			if err != nil {
				return err
			}
			argsBlockHeight, err := cast.ToUint64E(args[8])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateSendVoter(clientCtx.GetFromAddress().String(), index, argsSender, argsSenderChainId, argsReceiver, argsReceiverChainId, argsMBurnt, argsMessage, argsTxHash, argsBlockHeight)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
