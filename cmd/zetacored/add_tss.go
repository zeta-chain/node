package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func AddTssToGenesisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-node-account [tssKeyName] [Password] [operatorAddress]",
		Short: "Get you node account",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			tssKeyName := args[0]
			operatorAddress := args[2]
			pubKeySet, err := GetPubKeySet(clientCtx, tssKeyName, args[1])
			if err != nil {
				return err
			}
			sdk.MustAccAddressFromBech32(operatorAddress)

			k, err := clientCtx.Keyring.Key(tssKeyName)
			if err != nil {
				return err
			}

			address, err := k.GetAddress()
			if err != nil {
				return err
			}

			nodeAccount := types.NodeAccount{
				Creator:          args[2],
				TssSignerAddress: address.String(),
				PubkeySet:        &pubKeySet,
			}
			info := ObserverInfoReader{NodeAccount: nodeAccount}
			fmt.Println(info.String())
			return nil
		},
	}
	return cmd
}

func GetPubKeySet(clientctx client.Context, tssAccountName, password string) (common.PubKeySet, error) {
	pubkeySet := common.PubKeySet{
		Secp256k1: "",
		Ed25519:   "",
	}
	//kb, err := GetKeyringKeybase(keyringPath, tssAccountName, password)
	privKeyArmor, err := clientctx.Keyring.ExportPrivKeyArmor(tssAccountName, password)
	if err != nil {
		return pubkeySet, err
	}
	priKey, _, err := crypto.UnarmorDecryptPrivKey(privKeyArmor, password)
	if err != nil {
		return pubkeySet, fmt.Errorf("fail to unarmor private key: %w", err)
	}

	s, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, priKey.PubKey())
	if err != nil {
		return pubkeySet, err
	}
	pubkey, err := common.NewPubKey(s)
	if err != nil {
		return pubkeySet, err
	}
	pubkeySet.Secp256k1 = pubkey
	return pubkeySet, nil
}
