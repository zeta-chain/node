package cli

import (
	"strconv"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/zetacore/x/observer/types"
)

func CmdUpdateObserver() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-observer [old-observer-address] [new-observer-address] [update-reason]",
		Short: "Broadcast message add-observer",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			updateReasonInt, err := strconv.ParseInt(args[2], 10, 32)
			if err != nil {
				return err
			}
			// #nosec G115 parsed in range
			updateReason, err := parseUpdateReason(int32(updateReasonInt))
			if err != nil {
				return err
			}
			msg := types.NewMsgUpdateObserver(
				clientCtx.GetFromAddress().String(),
				args[0],
				args[1],
				updateReason,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func parseUpdateReason(i int32) (types.ObserverUpdateReason, error) {
	if _, ok := types.ObserverUpdateReason_name[i]; ok {
		switch i {
		case 1:
			return types.ObserverUpdateReason_Tombstoned, nil
		case 2:
			return types.ObserverUpdateReason_AdminUpdate, nil
		}
	}
	return types.ObserverUpdateReason_Tombstoned, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid update reason")
}
