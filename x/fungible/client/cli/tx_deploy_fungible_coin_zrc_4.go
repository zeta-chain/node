package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func CmdDeployFungibleCoinZRC4() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy-fungible-coin-zrc-4 [erc-20] [foreign-chain] [decimals] [name] [symbol] [coin-type] [gas-limit]",
		Short: "Broadcast message DeployFungibleCoinZRC20",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argERC20 := args[0]
			argForeignChain, err := strconv.ParseInt(args[1], 10, 32)
			if err != nil {
				return err
			}
			argDecimals, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return err
			}
			argName := args[3]
			argSymbol := args[4]
			argCoinType, err := strconv.ParseInt(args[5], 10, 32)
			if err != nil {
				return err
			}
			argGasLimit, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgDeployFungibleCoinZRC20(
				clientCtx.GetFromAddress().String(),
				argERC20,
				argForeignChain,
				// #nosec G115 parsed in range
				uint32(argDecimals),
				argName,
				argSymbol,
				coin.CoinType(argCoinType),
				argGasLimit,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
