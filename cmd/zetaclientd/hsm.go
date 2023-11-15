package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	keystone "github.com/zeta-chain/keystone/keys"
	"github.com/zeta-chain/zetacore/cmd"
	"github.com/zeta-chain/zetacore/common/cosmos"
	"github.com/zeta-chain/zetacore/zetaclient/hsm"
)

var HsmCmd = &cobra.Command{
	Use:   "hsm",
	Short: "Utility command to interact with hsm",
}

var GetHsmAddressCmd = &cobra.Command{
	Use:   "get-address",
	Short: "Get the address of a particular keypair by label",
	RunE:  GetHsmAddress,
}

var GenerateHsmKeyCmd = &cobra.Command{
	Use:   "gen-key",
	Short: "Generate keypair by label",
	RunE:  GenerateHsmKey,
}

type HsmArgs struct {
	label string
}

type HsmGenKeyArgs struct {
	algorithm int
}

var hsmArgs = HsmArgs{}
var hsmKeyGenArgs = HsmGenKeyArgs{}

func init() {
	RootCmd.AddCommand(HsmCmd)
	HsmCmd.AddCommand(GetHsmAddressCmd)
	HsmCmd.AddCommand(GenerateHsmKeyCmd)

	// HSM root arguments
	HsmCmd.PersistentFlags().StringVar(&hsmArgs.label, "key-label", "", "label used to identify key on HSM")

	// HSM key gen arguments
	GenerateHsmKeyCmd.Flags().IntVar(&hsmKeyGenArgs.algorithm, "algorithm", 0, "key algo; 0=SECP256K1, 1=SECP256R1, 2=ED25519")
}

func GetHsmAddress(_ *cobra.Command, _ []string) error {
	SetupConfigForTest()

	config, err := hsm.GetPKCS11Config()
	if err != nil {
		return err
	}
	_, pubKey, err := hsm.GetHSMAddress(config, hsmArgs.label)
	if err != nil {
		return err
	}

	address, err := cosmos.Bech32ifyAddressBytes(cmd.Bech32PrefixAccAddr, pubKey.Address().Bytes())
	if err != nil {
		return err
	}
	zetaPubKey, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
	if err != nil {
		return err
	}

	// Print formatted result
	fmt.Println("Address: ", address)
	fmt.Println("Public Key: ", zetaPubKey)
	fmt.Println("Label: ", hsmArgs.label)

	return nil
}

func GenerateHsmKey(_ *cobra.Command, _ []string) error {
	config, err := hsm.GetPKCS11Config()
	if err != nil {
		return err
	}
	if hsmKeyGenArgs.algorithm > 2 || hsmKeyGenArgs.algorithm < 0 {
		return errors.New("invalid algorithm selected")
	}
	algo := []keystone.KeygenAlgorithm{keystone.KEYGEN_SECP256K1, keystone.KEYGEN_SECP256R1, keystone.KEYGEN_ED25519}
	key, err := hsm.GenerateKey(hsmArgs.label, algo[hsmKeyGenArgs.algorithm], config)
	if err != nil {
		return err
	}

	// Print Generated key
	fmt.Println("Public Key: ", key.PubKey().String())
	fmt.Println("Label: ", hsmArgs.label)

	return nil
}
