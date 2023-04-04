package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	ckeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"io"
	"os"
)

func AddTssToGenesisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gentss [tssKeyName] [Password]",
		Short: "Add your tss address to the genesis file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}
			pubKeySet, address, err := GetPubKeySet(clientCtx, args[0], args[1])
			if err != nil {
				return err
			}

			zetaCrossChainGenState := types.GetGenesisStateFromAppState(cdc, appState)
			zetaCrossChainGenState.NodeAccountList = append(zetaCrossChainGenState.NodeAccountList, &types.NodeAccount{
				Creator:     address.String(),
				Index:       address.String(),
				NodeAddress: address,
				PubkeySet:   &pubKeySet,
				NodeStatus:  types.NodeStatus_Unknown,
			})

			zetaCrossChainStateBz, err := json.Marshal(zetaCrossChainGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal Observer List into Genesis File: %w", err)
			}

			if err != nil {
				return fmt.Errorf("failed to authz grants into Genesis File: %w", err)
			}
			appState[types.ModuleName] = zetaCrossChainStateBz
			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}
			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}
	return cmd
}

func GetPubKeySet(clientctx client.Context, tssAccountName, password string) (common.PubKeySet, sdk.AccAddress, error) {
	pubkeySet := common.PubKeySet{
		Secp256k1: "",
		Ed25519:   "",
	}
	address := sdk.AccAddress{}
	//kb, err := GetKeyringKeybase(keyringPath, tssAccountName, password)
	k, err := clientctx.Keyring.Key(tssAccountName)
	if err != nil {
		return pubkeySet, address, err
	}
	address, err = k.GetAddress()
	if err != nil {
		return pubkeySet, address, err
	}
	privKeyArmor, err := clientctx.Keyring.ExportPrivKeyArmor(tssAccountName, password)
	if err != nil {
		return pubkeySet, address, err
	}
	priKey, _, err := crypto.UnarmorDecryptPrivKey(privKeyArmor, password)
	if err != nil {
		return pubkeySet, address, fmt.Errorf("fail to unarmor private key: %w", err)
	}

	s, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, priKey.PubKey())
	if err != nil {
		return pubkeySet, address, err
	}
	pubkey, err := common.NewPubKey(s)
	if err != nil {
		return pubkeySet, address, err
	}
	pubkeySet.Secp256k1 = pubkey
	return pubkeySet, address, nil
}

func GetKeyringKeybase(chainHomeFolder, tssKeyName, password string) (ckeys.Keyring, error) {

	if len(tssKeyName) == 0 {
		return nil, fmt.Errorf("signer name is empty")
	}
	if len(password) == 0 {
		return nil, fmt.Errorf("password is empty")
	}

	buf := bytes.NewBufferString(password)
	// the library used by keyring is using ReadLine , which expect a new line
	buf.WriteByte('\n')
	buf.WriteString(password)
	buf.WriteByte('\n')
	kb, err := getKeybase(chainHomeFolder, buf)
	if err != nil {
		return nil, fmt.Errorf("fail to get keybase,err:%w", err)
	}
	oldStdIn := os.Stdin
	defer func() {
		os.Stdin = oldStdIn
	}()
	os.Stdin = nil
	_, err = kb.Key(tssKeyName)
	if err != nil {
		return nil, err
	}
	return kb, nil
}

// getKeybase will create an instance of Keybase
func getKeybase(zetaCoreHome string, reader io.Reader) (ckeys.Keyring, error) {
	cliDir := zetaCoreHome
	if len(zetaCoreHome) == 0 {
		return nil, fmt.Errorf("zetaCoreHome is empty")
	}
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	return ckeys.New(sdk.KeyringServiceName(), ckeys.BackendTest, cliDir, reader, cdc)
}
