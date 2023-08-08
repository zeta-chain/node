package types

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/zeta-chain/zetacore/common"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetSequence(ctx sdk.Context, addr sdk.AccAddress) (uint64, error)
	GetModuleAccount(ctx sdk.Context, name string) types.ModuleAccountI
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

type ZetaObserverKeeper interface {
	SetObserverMapper(ctx sdk.Context, om *zetaObserverTypes.ObserverMapper)
	GetObserverMapper(ctx sdk.Context, chain *common.Chain) (val zetaObserverTypes.ObserverMapper, found bool)
	GetAllObserverMappers(ctx sdk.Context) (mappers []*zetaObserverTypes.ObserverMapper)
	SetBallot(ctx sdk.Context, ballot *zetaObserverTypes.Ballot)
	GetBallot(ctx sdk.Context, index string) (val zetaObserverTypes.Ballot, found bool)
	GetAllBallots(ctx sdk.Context) (voters []*zetaObserverTypes.Ballot)
	GetParams(ctx sdk.Context) (params zetaObserverTypes.Params)
	GetCoreParamsByChainID(ctx sdk.Context, chainID int64) (params *zetaObserverTypes.CoreParams, found bool)
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
}
