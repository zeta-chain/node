package cli

import (
	"fmt"
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
		Use:   "receive-confirmation [sendHash] [outTxHash] [outBlockHeight] [mMint] [Status] [chain] [outTXNonce]",
		Short: "Broadcast message receiveConfirmation",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsSendHash := (args[0])
			argsOutTxHash := (args[1])
			argsOutBlockHeight, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}
			argsMMint := (args[3])
			var status common.ReceiveStatus
			if args[4] == "0" {
				status = common.ReceiveStatus_Success
			} else if args[4] == "1" {
				status = common.ReceiveStatus_Failed
			} else {
				return fmt.Errorf("wrong status")
			}
			chain := args[5]
			outTxNonce, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgReceiveConfirmation(clientCtx.GetFromAddress().String(), (argsSendHash), (argsOutTxHash), uint64(argsOutBlockHeight), (argsMMint), status, chain, uint64(outTxNonce))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
