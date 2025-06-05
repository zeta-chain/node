package cli

import (
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func CmdMigrateConnectorFunds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-connector-funds [chain-id] [new-connector-address] [amount]",
		Short: "Migrate connector funds to a new connector address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get chain ID
			chainID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			// Get new connector address
			newConnectorAddress := args[1]

			// Get amount
			amount, err := sdkmath.ParseUint(args[2])
			if err != nil {
				return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid amount: %s", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMigrateConnectorFunds(
				clientCtx.GetFromAddress().String(),
				chainID,
				newConnectorAddress,
				amount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
