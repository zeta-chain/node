package types

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/x/evm/statedb"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/zeta-chain/zetacore/common"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetSequence(ctx sdk.Context, addr sdk.AccAddress) (uint64, error)
	GetModuleAccount(ctx sdk.Context, name string) types.ModuleAccountI
	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	SetAccount(ctx sdk.Context, acc types.AccountI)
}

type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	IsSendEnabledCoin(ctx sdk.Context, coin sdk.Coin) bool
	BlockedAddr(addr sdk.AccAddress) bool
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata)
	HasSupply(ctx sdk.Context, denom string) bool
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type ObserverKeeper interface {
	SetObserverMapper(ctx sdk.Context, om *observertypes.ObserverMapper)
	GetObserverMapper(ctx sdk.Context, chain *common.Chain) (val observertypes.ObserverMapper, found bool)
	GetAllObserverMappers(ctx sdk.Context) (mappers []*observertypes.ObserverMapper)
	SetBallot(ctx sdk.Context, ballot *observertypes.Ballot)
	GetBallot(ctx sdk.Context, index string) (val observertypes.Ballot, found bool)
	GetAllBallots(ctx sdk.Context) (voters []*observertypes.Ballot)
	GetParams(ctx sdk.Context) (params observertypes.Params)
	GetCoreParamsByChainID(ctx sdk.Context, chainID int64) (params *common.CoreParams, found bool)
}

type EVMKeeper interface {
	ChainID() *big.Int
	GetBlockBloomTransient(ctx sdk.Context) *big.Int
	GetLogSizeTransient(ctx sdk.Context) uint64
	WithChainID(ctx sdk.Context)
	SetBlockBloomTransient(ctx sdk.Context, bloom *big.Int)
	SetLogSizeTransient(ctx sdk.Context, logSize uint64)
	EstimateGas(c context.Context, req *evmtypes.EthCallRequest) (*evmtypes.EstimateGasResponse, error)
	ApplyMessage(
		ctx sdk.Context,
		msg core.Message,
		tracer vm.EVMLogger,
		commit bool,
	) (*evmtypes.MsgEthereumTxResponse, error)
	GetAccount(ctx sdk.Context, addr ethcommon.Address) *statedb.Account
	SetAccount(ctx sdk.Context, addr ethcommon.Address, account statedb.Account) error
}
