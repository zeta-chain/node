package mocks

import (
	ckeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var TestKeyringPair cryptotypes.PrivKey

func init() {
	TestKeyringPair = cryptotypes.PrivKey(secp256k1.GenPrivKey())
}

type Keyring struct {
}

var _ ckeys.Keyring = Keyring{}

func NewKeyring() ckeys.Keyring {
	return Keyring{}
}

func (m Keyring) Backend() string {
	return ""
}

func (m Keyring) List() ([]*ckeys.Record, error) {
	return nil, nil
}

func (m Keyring) SupportedAlgorithms() (ckeys.SigningAlgoList, ckeys.SigningAlgoList) {
	return nil, nil
}

func (m Keyring) Key(_ string) (*ckeys.Record, error) {
	return ckeys.NewLocalRecord("", TestKeyringPair, TestKeyringPair.PubKey())
}

func (m Keyring) KeyByAddress(_ sdk.Address) (*ckeys.Record, error) {
	return nil, nil
}

func (m Keyring) Delete(_ string) error {
	return nil
}

func (m Keyring) DeleteByAddress(_ sdk.Address) error {
	return nil
}

func (m Keyring) Rename(_ string, _ string) error {
	return nil
}

func (m Keyring) NewMnemonic(_ string, _ ckeys.Language, _, _ string, _ ckeys.SignatureAlgo) (*ckeys.Record, string, error) {
	return nil, "", nil
}

func (m Keyring) NewAccount(_, _, _, _ string, _ ckeys.SignatureAlgo) (*ckeys.Record, error) {
	return nil, nil
}

func (m Keyring) SaveLedgerKey(_ string, _ ckeys.SignatureAlgo, _ string, _, _, _ uint32) (*ckeys.Record, error) {
	return nil, nil
}

func (m Keyring) SaveOfflineKey(_ string, _ cryptotypes.PubKey) (*ckeys.Record, error) {
	return nil, nil
}

func (m Keyring) SaveMultisig(_ string, _ cryptotypes.PubKey) (*ckeys.Record, error) {
	return nil, nil
}

func (m Keyring) Sign(_ string, _ []byte) ([]byte, cryptotypes.PubKey, error) {
	return nil, nil, nil
}

func (m Keyring) SignByAddress(_ sdk.Address, _ []byte) ([]byte, cryptotypes.PubKey, error) {
	return nil, nil, nil
}

func (m Keyring) ImportPrivKey(_, _, _ string) error {
	return nil
}

func (m Keyring) ImportPrivKeyHex(_, _, _ string) error {
	return nil
}

func (m Keyring) ImportPubKey(_ string, _ string) error {
	return nil
}

func (m Keyring) ExportPubKeyArmor(_ string) (string, error) {
	return "", nil
}

func (m Keyring) ExportPubKeyArmorByAddress(_ sdk.Address) (string, error) {
	return "", nil
}

func (m Keyring) ExportPrivKeyArmor(_, _ string) (armor string, err error) {
	return "", nil
}

func (m Keyring) ExportPrivKeyArmorByAddress(_ sdk.Address, _ string) (armor string, err error) {
	return "", nil
}

func (m Keyring) MigrateAll() ([]*ckeys.Record, error) {
	return nil, nil
}
