package cosmos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32/legacybech32" // nolint
	se "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const DefaultCoinDecimals = 18

var (
	KeyringServiceName      = sdk.KeyringServiceName
	NewRoute                = sdk.NewRoute
	NewKVStoreKeys          = sdk.NewKVStoreKeys
	NewUint                 = sdk.NewUint
	NewInt                  = sdk.NewInt
	NewDec                  = sdk.NewDec
	ZeroDec                 = sdk.ZeroDec
	NewCoin                 = sdk.NewCoin
	NewCoins                = sdk.NewCoins
	ParseCoins              = sdk.ParseCoinsNormalized
	NewDecWithPrec          = sdk.NewDecWithPrec
	NewDecFromBigInt        = sdk.NewDecFromBigInt
	NewIntFromBigInt        = sdk.NewIntFromBigInt
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
	NewContext              = sdk.NewContext
	Bech32ifyAddressBytes   = sdk.Bech32ifyAddressBytes
	GetPubKeyFromBech32     = legacybech32.UnmarshalPubKey
	Bech32ifyPubKey         = legacybech32.MarshalPubKey
	Bech32PubKeyTypeConsPub = legacybech32.ConsPK
	Bech32PubKeyTypeAccPub  = legacybech32.AccPK
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
	Querier    = sdk.Querier
	TxResponse = sdk.TxResponse
	Account    = authtypes.AccountI
)

var _ sdk.Address = AccAddress{}
