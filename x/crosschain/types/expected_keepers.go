package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	eth "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/zetacore/common"

	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

type StakingKeeper interface {
	GetAllValidators(ctx sdk.Context) (validators []stakingtypes.Validator)
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) types.ModuleAccountI

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, types.ModuleAccountI)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
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
	GetNodeAccount(ctx sdk.Context, address string) (nodeAccount zetaObserverTypes.NodeAccount, found bool)
	GetAllNodeAccount(ctx sdk.Context) (nodeAccounts []zetaObserverTypes.NodeAccount)
	SetNodeAccount(ctx sdk.Context, nodeAccount zetaObserverTypes.NodeAccount)
	IsInboundAllowed(ctx sdk.Context) (found bool)
	GetKeygen(ctx sdk.Context) (val zetaObserverTypes.Keygen, found bool)
	SetKeygen(ctx sdk.Context, keygen zetaObserverTypes.Keygen)
	SetPermissionFlags(ctx sdk.Context, permissionFlags zetaObserverTypes.PermissionFlags)
	SetLastObserverCount(ctx sdk.Context, lbc *zetaObserverTypes.LastObserverCount)
	AddVoteToBallot(
		ctx sdk.Context,
		ballot zetaObserverTypes.Ballot,
		address string,
		observationType zetaObserverTypes.VoteType,
	) (zetaObserverTypes.Ballot, error)
	CheckIfFinalizingVote(ctx sdk.Context, ballot zetaObserverTypes.Ballot) (zetaObserverTypes.Ballot, bool)
	IsAuthorized(ctx sdk.Context, address string, chain *common.Chain) (bool, error)
	FindBallot(
		ctx sdk.Context,
		index string,
		chain *common.Chain,
		observationType zetaObserverTypes.ObservationType,
	) (ballot zetaObserverTypes.Ballot, isNew bool, err error)
}

type FungibleKeeper interface {
	GetForeignCoins(ctx sdk.Context, zrc20Addr string) (val fungibletypes.ForeignCoins, found bool)
	GetAllForeignCoins(ctx sdk.Context) (list []fungibletypes.ForeignCoins)
	SetForeignCoins(ctx sdk.Context, foreignCoins fungibletypes.ForeignCoins)
	GetAllForeignCoinsForChain(ctx sdk.Context, foreignChainID int64) (list []fungibletypes.ForeignCoins)
	GetSystemContract(ctx sdk.Context) (val fungibletypes.SystemContract, found bool)
	QuerySystemContractGasCoinZRC20(ctx sdk.Context, chainID *big.Int) (eth.Address, error)
	QueryUniswapv2RouterGetAmountsIn(ctx sdk.Context, amountOut *big.Int, outZRC4 eth.Address) (*big.Int, error)
	SetGasPrice(ctx sdk.Context, chainID *big.Int, gasPrice *big.Int) (uint64, error)
	DepositCoinZeta(ctx sdk.Context, to eth.Address, amount *big.Int) error
	ZRC20DepositAndCallContract(
		ctx sdk.Context,
		from []byte,
		to eth.Address,
		amount *big.Int,
		senderChain *common.Chain,
		message string,
		contract eth.Address,
		data []byte,
		coinType common.CoinType,
		asset string,
	) (*evmtypes.MsgEthereumTxResponse, error)
	CallUniswapv2RouterSwapExactETHForToken(
		ctx sdk.Context,
		sender eth.Address,
		to eth.Address,
		amountIn *big.Int,
		outZRC4 eth.Address,
	) ([]*big.Int, error)
	CallZRC20Burn(ctx sdk.Context, sender eth.Address, zrc20address eth.Address, amount *big.Int) error
	DeployZRC20Contract(
		ctx sdk.Context,
		name, symbol string,
		decimals uint8,
		chainId int64,
		coinType common.CoinType,
		erc20Contract string,
		gasLimit *big.Int,
	) (eth.Address, error)
}
