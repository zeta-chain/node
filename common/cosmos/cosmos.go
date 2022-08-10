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
	StoreKey   = sdk.StoreKey
	Querier    = sdk.Querier
	TxResponse = sdk.TxResponse
	Account    = authtypes.AccountI
)

var _ sdk.Address = AccAddress{}
