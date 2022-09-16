package cli

import (
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/common"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

var _ = strconv.Itoa(0)

func CmdSendVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-voter [sender] [senderChain] [receiver] [receiverChain] [mBurnt] [mMint] [message] [inTxHash] [inBlockHeight] [coinType]",
		Short: "Broadcast message sendVoter",
		Args:  cobra.ExactArgs(10),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsSender := (args[0])
			argsSenderChain := (args[1])
			argsReceiver := (args[2])
			argsReceiverChain := (args[3])
			argsMBurnt := (args[4])
			argsMMint := (args[5])
			argsMessage := (args[6])
			argsInTxHash := (args[7])
			argsInBlockHeight, err := strconv.ParseInt(args[8], 10, 64)
			if err != nil {
				return err
			}
			argsCoinType, err := strconv.ParseInt(args[9], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSendVoter(clientCtx.GetFromAddress().String(), (argsSender), (argsSenderChain), (argsReceiver), (argsReceiverChain), (argsMBurnt), (argsMMint), (argsMessage), (argsInTxHash), uint64(argsInBlockHeight), 250_000, common.CoinType(argsCoinType))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
