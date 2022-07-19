package zetaclient

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/crypto"
	ckeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	metachainCliFolderPath = `.meta-chaind`
)

// Keys manages all the keys used by metachain
type Keys struct {
	signerName string
	password   string // TODO this is a bad way , need to fix it
	kb         ckeys.Keyring
}

// NewKeysWithKeybase create a new instance of Keys
func NewKeysWithKeybase(kb ckeys.Keyring, name, password string) *Keys {
	return &Keys{
		signerName: name,
		password:   password,
		kb:         kb,
	}
}

// GetKeyringKeybase return keyring and key info
func GetKeyringKeybase(chainHomeFolder, signerName, password string) (ckeys.Keyring, ckeys.Info, error) {
	if len(signerName) == 0 {
		return nil, nil, fmt.Errorf("signer name is empty")
	}
	if len(password) == 0 {
		return nil, nil, fmt.Errorf("password is empty")
	}

	buf := bytes.NewBufferString(password)
	// the library used by keyring is using ReadLine , which expect a new line
	buf.WriteByte('\n')
	buf.WriteString(password)
	buf.WriteByte('\n')
	//fmt.Printf("password buf: %s\n", buf)
	kb, err := getKeybase(chainHomeFolder, buf)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to get keybase,err:%w", err)
	}
	// the keyring library which used by cosmos sdk , will use interactive terminal if it detect it has one
	// this will temporary trick it think there is no interactive terminal, thus will read the password from the buffer provided
	oldStdIn := os.Stdin
	defer func() {
		os.Stdin = oldStdIn
	}()
	os.Stdin = nil
	//fmt.Println("signer: ", signerName)
	si, err := kb.Key(signerName)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to get signer info(%s): %w;", signerName, err)
	}
	return kb, si, nil
}

// getKeybase will create an instance of Keybase
func getKeybase(metacoreHome string, reader io.Reader) (ckeys.Keyring, error) {
	cliDir := metacoreHome
	if len(metacoreHome) == 0 {
		usr, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("fail to get current user,err:%w", err)
		}
		cliDir = filepath.Join(usr.HomeDir, metachainCliFolderPath)
	}
	//FIXME: BackendTest is used for convenient testing with Starport generated accouts.
	// Change to BackendFile with password!
	return ckeys.New(sdk.KeyringServiceName(), ckeys.BackendTest, cliDir, reader)
}

// GetSignerInfo return signer info
func (k *Keys) GetSignerInfo() ckeys.Info {
	info, err := k.kb.Key(k.signerName)
	if err != nil {
		panic(err)
	}
	return info
}

// GetPrivateKey return the private key
func (k *Keys) GetPrivateKey() (cryptotypes.PrivKey, error) {
	// return k.kb.ExportPrivateKeyObject(k.signerName)
	privKeyArmor, err := k.kb.ExportPrivKeyArmor(k.signerName, k.password)
	if err != nil {
		return nil, err
	}
	priKey, _, err := crypto.UnarmorDecryptPrivKey(privKeyArmor, k.password)
	if err != nil {
		return nil, fmt.Errorf("fail to unarmor private key: %w", err)
	}
	return priKey, nil
}

// GetKeybase return the keybase
func (k *Keys) GetKeybase() ckeys.Keyring {
	return k.kb
}

func (k *Keys) GetPubKeySet() (common.PubKeySet, error) {
	pubkeySet := common.PubKeySet{
		Secp256k1: "",
		Ed25519:   "",
	}
	pK, err := k.GetPrivateKey()
	if err != nil {
		return pubkeySet, err
	}

	s, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pK.PubKey())
	if err != nil {
		return pubkeySet, ErrBech32ifyPubKey
	}
	pubkey, err := common.NewPubKey(s)
	if err != nil {
		return pubkeySet, errors.Wrap(ErrNewPubKey, fmt.Sprintf("Pubkey %s", pK.PubKey().String()))
	}
	pubkeySet.Secp256k1 = pubkey
	return pubkeySet, nil
}
