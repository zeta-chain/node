package main

import (
	"context"
	"fmt"
	"os"

	ecdsakeygen "github.com/bnb-chain/tss-lib/ecdsa/keygen"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/cmd"
	"github.com/zeta-chain/node/pkg/constant"
)

// globalOptions defines the global options for all commands.
type globalOptions struct {
	ZetacoreHome string
}

var (
	RootCmd = &cobra.Command{
		Use:   "zetaclientd",
		Short: "ZetaClient CLI",
	}
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "prints version",
		Run:   func(_ *cobra.Command, _ []string) { fmt.Print(constant.Version) },
	}
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Start ZetaClient Observer",
		RunE:  Start,
	}
	InitializeConfigCmd = &cobra.Command{
		Use:     "init-config",
		Aliases: []string{"init"},
		Short:   "Initialize Zetaclient Configuration file",
		RunE:    InitializeConfig,
	}
)

var (
	preParams  *ecdsakeygen.LocalPreParams
	globalOpts globalOptions
)

func main() {
	ctx := context.Background()

	if err := RootCmd.ExecuteContext(ctx); err != nil {
		fmt.Printf("Error: %s. Exit code 1\n", err)
		os.Exit(1)
	}
}

func init() {
	cmd.SetupCosmosConfig()

	// Setup options
	setupGlobalOptions()
	setupInitializeConfigOptions()

	// Define commands
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(StartCmd)
	RootCmd.AddCommand(InitializeConfigCmd)
}

func setupGlobalOptions() {
	globals := RootCmd.PersistentFlags()

	globals.StringVar(&globalOpts.ZetacoreHome, tmcli.HomeFlag, app.DefaultNodeHome, "home path")
	// add more options here (e.g. verbosity, etc...)
}
