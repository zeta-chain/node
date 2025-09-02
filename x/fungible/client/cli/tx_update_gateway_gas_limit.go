package cli

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/fungible/types"
)

func CmdUpdateGatewayGasLimit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-gateway-gas-limit [gas-limit]",
		Short: "Broadcast message UpdateGatewayGasLimit to update the gateway gas limit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasLimit, ok := sdkmath.NewIntFromString(args[0])
			if !ok {
				return types.ErrInvalidGasLimit
			}

			msg := types.NewMsgUpdateGatewayGasLimit(clientCtx.GetFromAddress().String(), gasLimit)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
