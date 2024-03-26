package stub

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

type mockKeyring struct {
}

var _ ckeys.Keyring = mockKeyring{}

func NewMockKeyring() ckeys.Keyring {
	return mockKeyring{}
}

func (m mockKeyring) Backend() string {
	return ""
}

func (m mockKeyring) List() ([]*ckeys.Record, error) {
	return nil, nil
}

func (m mockKeyring) SupportedAlgorithms() (ckeys.SigningAlgoList, ckeys.SigningAlgoList) {
	return nil, nil
}

func (m mockKeyring) Key(_ string) (*ckeys.Record, error) {
	record, _ := ckeys.NewLocalRecord("", TestKeyringPair, TestKeyringPair.PubKey())
	return record, nil
}

func (m mockKeyring) KeyByAddress(_ sdk.Address) (*ckeys.Record, error) {
	return nil, nil
}

func (m mockKeyring) Delete(_ string) error {
	return nil
}

func (m mockKeyring) DeleteByAddress(_ sdk.Address) error {
	return nil
}

func (m mockKeyring) Rename(_ string, _ string) error {
	return nil
}

func (m mockKeyring) NewMnemonic(_ string, _ ckeys.Language, _, _ string, _ ckeys.SignatureAlgo) (*ckeys.Record, string, error) {
	return nil, "", nil
}

func (m mockKeyring) NewAccount(_, _, _, _ string, _ ckeys.SignatureAlgo) (*ckeys.Record, error) {
	return nil, nil
}

func (m mockKeyring) SaveLedgerKey(_ string, _ ckeys.SignatureAlgo, _ string, _, _, _ uint32) (*ckeys.Record, error) {
	return nil, nil
}

func (m mockKeyring) SaveOfflineKey(_ string, _ cryptotypes.PubKey) (*ckeys.Record, error) {
	return nil, nil
}

func (m mockKeyring) SaveMultisig(_ string, _ cryptotypes.PubKey) (*ckeys.Record, error) {
	return nil, nil
}

func (m mockKeyring) Sign(_ string, _ []byte) ([]byte, cryptotypes.PubKey, error) {
	return nil, nil, nil
}

func (m mockKeyring) SignByAddress(_ sdk.Address, _ []byte) ([]byte, cryptotypes.PubKey, error) {
	return nil, nil, nil
}

func (m mockKeyring) ImportPrivKey(_, _, _ string) error {
	return nil
}

func (m mockKeyring) ImportPubKey(_ string, _ string) error {
	return nil
}

func (m mockKeyring) ExportPubKeyArmor(_ string) (string, error) {
	return "", nil
}

func (m mockKeyring) ExportPubKeyArmorByAddress(_ sdk.Address) (string, error) {
	return "", nil
}

func (m mockKeyring) ExportPrivKeyArmor(_, _ string) (armor string, err error) {
	return "", nil
}

func (m mockKeyring) ExportPrivKeyArmorByAddress(_ sdk.Address, _ string) (armor string, err error) {
	return "", nil
}

func (m mockKeyring) MigrateAll() ([]*ckeys.Record, error) {
	return nil, nil
}
