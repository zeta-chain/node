package cli

import (
	"errors"
	"fmt"
	"strconv"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdVoteOutbound() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-outbound [sendHash] [outboundHash] [outBlockHeight] [outGasUsed] [outEffectiveGasPrice] [outEffectiveGasLimit] [valueReceived] [Status] [chain] [outTXNonce] [coinType]",
		Short: "Broadcast message to vote an outbound",
		Args:  cobra.ExactArgs(11),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsSendHash := args[0]
			argsOutboundHash := args[1]

			argsOutBlockHeight, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			argsOutGasUsed, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			argsOutEffectiveGasPrice, ok := math.NewIntFromString(args[4])
			if !ok {
				return errors.New("invalid effective gas price, enter 0 if unused")
			}

			argsOutEffectiveGasLimit, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return err
			}

			argsMMint := args[6]

			status, err := chains.ReceiveStatusFromString(args[7])
			if err != nil {
				return err
			}

			chain, err := strconv.ParseInt(args[8], 10, 64)
			if err != nil {
				return err
			}

			outTxNonce, err := strconv.ParseUint(args[9], 10, 64)
			if err != nil {
				return err
			}

			coinType, ok := coin.CoinType_value[args[10]]
			if !ok {
				return fmt.Errorf("wrong coin type %s", args[10])
			}
			argsCoinType := coin.CoinType(coinType)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgVoteOutbound(
				clientCtx.GetFromAddress().String(),
				argsSendHash,
				argsOutboundHash,
				argsOutBlockHeight,
				argsOutGasUsed,
				argsOutEffectiveGasPrice,
				argsOutEffectiveGasLimit,
				math.NewUintFromString(argsMMint),
				status,
				chain,
				outTxNonce,
				argsCoinType,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
