package main

import (
	"bufio"

	"github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	evmclient "github.com/cosmos/evm/client"
	clientkeys "github.com/cosmos/evm/client/keys"
	"github.com/cosmos/evm/crypto/hd"
	"github.com/spf13/cobra"
)

// TODO evm: unsupported algo, had to fork this command

// KeyCommands registers a subtree of commands to interact with
// local private key storage.
//
// The defaultToEthKeys boolean indicates whether to use
// Ethereum compatible keys by default or stick with the Cosmos SDK default.
func KeyCommands(defaultNodeHome string, defaultToEthKeys bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage your application's keys",
		Long: `Keyring management commands. These keys may be in any format supported by the
CometBFT crypto library and can be used by light-clients, full nodes, or any other application
that needs to sign with a private key.

The keyring supports the following backends:

    os          Uses the operating system's default credentials store.
    file        Uses encrypted file-based keystore within the app's configuration directory.
                This keyring will request a password each time it is accessed, which may occur
                multiple times in a single command resulting in repeated password prompts.
    kwallet     Uses KDE Wallet Manager as a credentials management application.
    pass        Uses the pass command line utility to store and retrieve keys.
    test        Stores keys insecurely to disk. It does not prompt for a password to be unlocked
                and it should be use only for testing purposes.

kwallet and pass backends depend on external tools. Refer to their respective documentation for more
information:
    KWallet     https://github.com/KDE/kwallet
    pass        https://www.passwordstore.org/

The pass backend requires GnuPG: https://gnupg.org/
`,
	}

	// support adding Ethereum supported keys
	addCmd := keys.AddKeyCommand()

	if defaultToEthKeys {
		ethSecp256k1Str := string(hd.EthSecp256k1Type)

		// update the default signing algorithm value to "eth_secp256k1"
		algoFlag := addCmd.Flag(flags.FlagKeyType)
		algoFlag.DefValue = ethSecp256k1Str
		err := algoFlag.Value.Set(ethSecp256k1Str)
		if err != nil {
			panic(err)
		}
	}

	addCmd.RunE = runAddCmd

	cmd.AddCommand(
		keys.MnemonicKeyCommand(),
		addCmd,
		keys.ExportKeyCommand(),
		keys.ImportKeyCommand(),
		keys.ListKeysCmd(),
		keys.ListKeyTypesCmd(),
		keys.ShowKeysCmd(),
		keys.DeleteKeyCommand(),
		keys.RenameKeyCommand(),
		keys.ParseKeyStringCommand(),
		keys.MigrateCommand(),
		flags.LineBreak,
		evmclient.UnsafeExportEthKeyCommand(),
		evmclient.UnsafeImportKeyCommand(),
	)

	cmd.PersistentFlags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.PersistentFlags().
		String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
	cmd.PersistentFlags().String(flags.FlagKeyringBackend, keyring.BackendOS, "Select keyring's backend (os|file|test)")
	cmd.PersistentFlags().String(cli.OutputFlag, "text", "Output format (text|json)")
	return cmd
}

func runAddCmd(cmd *cobra.Command, args []string) error {
	clientCtx := client.GetClientContextFromCmd(cmd).WithKeyringOptions(hd.EthSecp256k1Option())
	clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
	if err != nil {
		return err
	}
	buf := bufio.NewReader(clientCtx.Input)
	return clientkeys.RunAddCmd(clientCtx, cmd, args, buf)
}
