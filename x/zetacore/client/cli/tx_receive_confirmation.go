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

func CmdReceiveConfirmation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "receive-confirmation [sendHash] [outTxHash] [outBlockHeight] [mMint]",
		Short: "Broadcast message receiveConfirmation",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsSendHash := (args[0])
			argsOutTxHash := (args[1])
			argsOutBlockHeight, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}
			argsMMint := (args[3])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgReceiveConfirmation(clientCtx.GetFromAddress().String(), (argsSendHash), (argsOutTxHash), uint64(argsOutBlockHeight), (argsMMint), common.ReceiveStatus_Success, "ETH")
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
