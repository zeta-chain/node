package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gagliardetto/solana-go"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/e2e/config"
)

func NewPopulateAddressesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "populate-addresses [config-file] ",
		Short: "Derive addresses from the configured private keys and populate in the config file",
		RunE:  runPopulateAddresses,
		Args:  cobra.ExactArgs(1),
	}
	return cmd
}

func runPopulateAddresses(_ *cobra.Command, args []string) error {
	// do not validate the config on load as it will not have addresses populated yet
	conf, err := config.ReadConfig(args[0], false)
	if err != nil {
		return err
	}

	// evm address
	privKey, err := conf.DefaultAccount.PrivateKey()
	if err != nil {
		return fmt.Errorf("decoding private key: %w", err)
	}
	evmAddress := crypto.PubkeyToAddress(privKey.PublicKey)
	bech32Address, err := bech32.ConvertAndEncode("zeta", evmAddress.Bytes())
	if err != nil {
		return fmt.Errorf("bech32 convert and encode: %w", err)
	}
	conf.DefaultAccount.RawEVMAddress = config.DoubleQuotedString(evmAddress.String())
	conf.DefaultAccount.RawBech32Address = config.DoubleQuotedString(bech32Address)

	// solana address
	if conf.DefaultAccount.SolanaPrivateKey != "" {
		sPrivKey, err := solana.PrivateKeyFromBase58(conf.DefaultAccount.SolanaPrivateKey.String())
		if err != nil {
			return fmt.Errorf("decoding Solana private key: %w", err)
		}
		sAddress := sPrivKey.PublicKey().String()
		conf.DefaultAccount.SolanaAddress = config.DoubleQuotedString(sAddress)
	}

	err = config.WriteConfig(args[0], conf)
	if err != nil {
		return err
	}
	return nil
}
