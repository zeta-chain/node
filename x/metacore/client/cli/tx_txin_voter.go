package cli

import (
	"github.com/spf13/cobra"

	"github.com/spf13/cast"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

func CmdCreateTxinVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-txin-voter [index] [txHash] [sourceAsset] [sourceAmount] [mBurnt] [destinationAsset] [fromAddress] [toAddress] [blockHeight] [signer] [signature]",
		Short: "Create a new TxinVoter",
		Args:  cobra.ExactArgs(11),
		RunE: func(cmd *cobra.Command, args []string) error {
			index := args[0]
			argsTxHash, err := cast.ToStringE(args[1])
			if err != nil {
				return err
			}
			argsSourceAsset, err := cast.ToStringE(args[2])
			if err != nil {
				return err
			}
			argsSourceAmount, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}
			argsMBurnt, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}
			argsDestinationAsset, err := cast.ToStringE(args[5])
			if err != nil {
				return err
			}
			argsFromAddress, err := cast.ToStringE(args[6])
			if err != nil {
				return err
			}
			argsToAddress, err := cast.ToStringE(args[7])
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

			_ = index // index is set to "TxHash-Creator" automatically by NewMsgCreateTxinVoter
			msg := types.NewMsgCreateTxinVoter(clientCtx.GetFromAddress().String(), argsTxHash, argsSourceAsset, argsSourceAmount, argsMBurnt, argsDestinationAsset, argsFromAddress, argsToAddress, argsBlockHeight)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
