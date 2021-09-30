package cosmos

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/input"
	ckeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	se "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/hashicorp/go-multierror"
)

const DefaultCoinDecimals = 8

var (
	KeyringServiceName      = sdk.KeyringServiceName
	NewRoute                = sdk.NewRoute
	NewKVStoreKeys          = sdk.NewKVStoreKeys
	NewUint                 = sdk.NewUint
	ParseUint               = sdk.ParseUint
	NewInt                  = sdk.NewInt
	NewDec                  = sdk.NewDec
	ZeroUint                = sdk.ZeroUint
	ZeroDec                 = sdk.ZeroDec
	OneUint                 = sdk.OneUint
	NewCoin                 = sdk.NewCoin
	NewCoins                = sdk.NewCoins
	ParseCoins              = sdk.ParseCoinsNormalized
	NewDecWithPrec          = sdk.NewDecWithPrec
	NewDecFromBigInt        = sdk.NewDecFromBigInt
	NewIntFromBigInt        = sdk.NewIntFromBigInt
	NewUintFromBigInt       = sdk.NewUintFromBigInt
	AccAddressFromBech32    = sdk.AccAddressFromBech32
	VerifyAddressFormat     = sdk.VerifyAddressFormat
	GetFromBech32           = sdk.GetFromBech32
	NewAttribute            = sdk.NewAttribute
	NewDecFromStr           = sdk.NewDecFromStr
	GetConfig               = sdk.GetConfig
	NewEvent                = sdk.NewEvent
	RegisterCodec           = sdk.RegisterLegacyAminoCodec
	NewEventManager         = sdk.NewEventManager
	EventTypeMessage        = sdk.EventTypeMessage
	AttributeKeyModule      = sdk.AttributeKeyModule
	KVStorePrefixIterator   = sdk.KVStorePrefixIterator
	NewKVStoreKey           = sdk.NewKVStoreKey
	NewTransientStoreKey    = sdk.NewTransientStoreKey
	StoreTypeTransient      = sdk.StoreTypeTransient
	StoreTypeIAVL           = sdk.StoreTypeIAVL
	NewContext              = sdk.NewContext
	GetPubKeyFromBech32     = sdk.GetPubKeyFromBech32
	Bech32ifyPubKey         = sdk.Bech32ifyPubKey
	Bech32PubKeyTypeConsPub = sdk.Bech32PubKeyTypeConsPub
	Bech32PubKeyTypeAccPub  = sdk.Bech32PubKeyTypeAccPub
	Wrapf                   = se.Wrapf
	MustSortJSON            = sdk.MustSortJSON
	CodeUnauthorized        = uint32(4)
	CodeInsufficientFunds   = uint32(5)
)

type (
	Context    = sdk.Context
	Route      = sdk.Route
	Uint       = sdk.Uint
	Coin       = sdk.Coin
	Coins      = sdk.Coins
	AccAddress = sdk.AccAddress
	Attribute  = sdk.Attribute
	Result     = sdk.Result
	Event      = sdk.Event
	Events     = sdk.Events
	Dec        = sdk.Dec
	Msg        = sdk.Msg
	Iterator   = sdk.Iterator
	Handler    = sdk.Handler
	StoreKey   = sdk.StoreKey
	Querier    = sdk.Querier
	TxResponse = sdk.TxResponse
	Account    = authtypes.AccountI
)

var _ sdk.Address = AccAddress{}

func ErrUnknownRequest(msg string) error {
	return se.Wrap(se.ErrUnknownRequest, msg)
}

func ErrInvalidAddress(addr string) error {
	return se.Wrap(se.ErrInvalidAddress, addr)
}

func ErrInvalidCoins(msg string) error {
	return se.Wrap(se.ErrInvalidCoins, msg)
}

func ErrUnauthorized(msg string) error {
	return se.Wrap(se.ErrUnauthorized, msg)
}

func ErrInsufficientCoins(err error, msg string) error {
	return se.Wrap(multierror.Append(se.ErrInsufficientFunds, err), msg)
}

/*
func SetupThorchainForTest(c *C) (config.ClientConfiguration, ckeys.Info, ckeys.Keybase) {
	thorchain.SetupConfigForTest()
	cfg := config.ClientConfiguration{
		ChainID:         "thorchain",
		ChainHost:       "localhost",
		SignerName:      "bob",
		SignerPasswd:    "password",
		ChainHomeFolder: ".",
	}
	kb := ckeys.NewInMemory()
	info, _, err := kb.NewMnemonic(cfg.SignerName, ckeys.English, cfg.SignerPasswd, hd.Secp256k1)
	c.Assert(err, IsNil)
	return cfg, info, kb
}
*/

// RoundToDecimal round the given amt to the desire decimals
func RoundToDecimal(amt Uint, dec int64) Uint {
	if dec != 0 && dec < DefaultCoinDecimals {
		prec := DefaultCoinDecimals - dec
		if prec == 0 { // sanity check
			return amt
		}
		precisionAdjust := sdk.NewUintFromBigInt(big.NewInt(0).Exp(big.NewInt(10), big.NewInt(prec), nil))
		amt = amt.Quo(precisionAdjust).Mul(precisionAdjust)
	}
	return amt
}

// KeybaseStore to store keys
type KeybaseStore struct {
	Keybase      ckeys.Keyring
	SignerName   string
	SignerPasswd string
}

func SignerCreds() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	username, _ := input.GetString("Enter Signer name:", reader)
	password, _ := input.GetPassword("Enter Signer password:", reader)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}

// GetKeybase will create an instance of Keybase
func GetKeybase(thorchainHome string) (KeybaseStore, error) {
	username, password := SignerCreds()
	buf := bytes.NewBufferString(password)
	// the library used by keyring is using ReadLine , which expect a new line
	buf.WriteByte('\n')

	cliDir := thorchainHome
	if len(thorchainHome) == 0 {
		usr, err := user.Current()
		if err != nil {
			return KeybaseStore{}, fmt.Errorf("fail to get current user,err:%w", err)
		}
		cliDir = filepath.Join(usr.HomeDir, ".thornode")
	}

	kb, err := ckeys.New(KeyringServiceName(), ckeys.BackendFile, cliDir, buf)
	return KeybaseStore{
		SignerName:   username,
		SignerPasswd: password,
		Keybase:      kb,
	}, err
}
