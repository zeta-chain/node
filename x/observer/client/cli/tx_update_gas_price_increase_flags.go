package cli

import (
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/zetacore/x/observer/types"
)

func CmdUpdateGasPriceIncreaseFlags() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-gas-price-increase-flags [epochLength] [retryInterval] [gasPriceIncreasePercent] [gasPriceIncreaseMax] [maxPendingCctxs]",
		Short: "Enable inbound and outbound cross-chain flags",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			epochLength, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			retryInterval, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			gasPriceIncreasePercent, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return err
			}
			gasPriceIncreaseMax, err := strconv.ParseUint(args[3], 10, 32)
			if err != nil {
				return err
			}
			maxPendingCctxs, err := strconv.ParseUint(args[4], 10, 32)
			if err != nil {
				return err
			}
			gasPriceIncreaseFlags := types.GasPriceIncreaseFlags{
				epochLength,
				time.Duration(retryInterval),
				uint32(gasPriceIncreasePercent),
				uint32(gasPriceIncreaseMax),
				uint32(maxPendingCctxs)}
			msg := types.NewMsgUpdateGasPriceIncreaseFlags(clientCtx.GetFromAddress().String(), gasPriceIncreaseFlags)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
