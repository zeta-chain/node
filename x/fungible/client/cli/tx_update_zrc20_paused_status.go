package cli

import (
	"strconv"
	"strings"

	cosmoserrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func CmdUpdateZRC20PausedStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-zrc20-paused-status [contractAddress1, contractAddress2, ...] [pausedStatus]",
		Short:   "Broadcast message UpdateZRC20PausedStatus",
		Example: `zetacored tx fungible update-zrc20-paused-status "0xece40cbB54d65282c4623f141c4a8a0bE7D6AdEc, 0xece40cbB54d65282c4623f141c4a8a0bEjgksncf" 0 `,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			contractAddressList := strings.Split(strings.TrimSpace(args[0]), ",")

			action, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}
			if (action != 0) && (action != 1) {
				return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid action (%d)", action)
			}

			pausedStatus := types.UpdatePausedStatusAction_PAUSE
			if action == 1 {
				pausedStatus = types.UpdatePausedStatusAction_UNPAUSE
			}

			msg := types.NewMsgUpdateZRC20PausedStatus(
				clientCtx.GetFromAddress().String(),
				contractAddressList,
				pausedStatus,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
