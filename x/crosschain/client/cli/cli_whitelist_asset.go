package cli

import (
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func CmdWhitelistAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist-asset [assetAddress] [chainID] [name] [symbol] [decimals] [gasLimit] [liquidityCap]",
		Short: "Add a new asset token to whitelist",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			assetAddress := args[0]
			chainID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			name := args[2]
			symbol := args[3]
			decimals, err := strconv.ParseUint(args[4], 10, 32)
			if err != nil {
				return err
			}
			if decimals > 128 {
				return fmt.Errorf("decimals must be less than 128")
			}

			gasLimit, err := strconv.ParseInt(args[5], 10, 64)
			if err != nil {
				return err
			}

			liquidityCap := sdkmath.NewUintFromString(args[6])

			msg := types.NewMsgWhitelistAsset(
				clientCtx.GetFromAddress().String(),
				assetAddress,
				chainID,
				name,
				symbol,
				// #nosec G115 always in range
				uint32(decimals),
				gasLimit,
				liquidityCap,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
