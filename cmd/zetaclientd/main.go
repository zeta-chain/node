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
	TSSCmd = &cobra.Command{
		Use:   "tss",
		Short: "TSS commands",
	}
	TSSEncryptCmd = &cobra.Command{
		Use:   "encrypt [file-path] [secret-key]",
		Short: "Utility command to encrypt existing tss key-share file",
		Args:  cobra.ExactArgs(2),
		RunE:  TSSEncryptFile,
	}
	TSSGeneratePreParamsCmd = &cobra.Command{
		Use:   "tss gen-pre-params [path]",
		Short: "Generate pre parameters for TSS",
		Args:  cobra.ExactArgs(1),
		RunE:  TSSGeneratePreParams,
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

	RootCmd.AddCommand(TSSCmd)
	TSSCmd.AddCommand(TSSEncryptCmd)
	TSSCmd.AddCommand(TSSGeneratePreParamsCmd)
}

func setupGlobalOptions() {
	globals := RootCmd.PersistentFlags()

	globals.StringVar(&globalOpts.ZetacoreHome, tmcli.HomeFlag, app.DefaultNodeHome, "home path")
	// add more options here (e.g. verbosity, etc...)
}
