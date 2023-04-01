package main

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// AddObserverAccountCmd Deprecated : Use AddObserverAccountsCmd instead
func AddrConversionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addr-conversion [zeta address]",
		Short: "convert a zeta1xxx address to validator operator address zetavaloper1xxx",
		Long: `
read a zeta1xxx or zetavaloper1xxx address and convert it to the other type;
it always outputs two lines; the first line is the zeta1xxx address, the second line is the zetavaloper1xxx address
			`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := sdk.AccAddressFromBech32(args[0])
			if err == nil {
				valAddr := sdk.ValAddress(addr.Bytes())
				fmt.Printf("%s\n", addr.String())
				fmt.Printf("%s\n", valAddr.String())
				return nil
			}
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err == nil {
				addr := sdk.AccAddress(valAddr.Bytes())
				fmt.Printf("%s\n", addr.String())
				fmt.Printf("%s\n", valAddr.String())
				return nil
			}
			return fmt.Errorf("invalid address: %s", args[0])
		},
	}
	return cmd
}
