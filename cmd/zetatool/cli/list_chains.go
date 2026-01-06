package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/cmd/zetatool/common"
)

func NewListChainsCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "list-chains",
		Short: "List all available chains",
		Long:  `List all available chains with their names and chain IDs.`,
		Args:  cobra.NoArgs,
		RunE:  listChains,
	}
}

func listChains(_ *cobra.Command, _ []string) error {
	fmt.Print(common.ListAvailableChains())
	return nil
}
