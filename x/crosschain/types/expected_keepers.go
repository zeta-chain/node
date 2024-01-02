package types

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	eth "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/zetacore/common"

	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
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
	SetObserverMapper(ctx sdk.Context, om *observertypes.ObserverMapper)
	GetObserverMapper(ctx sdk.Context, chain *common.Chain) (val observertypes.ObserverMapper, found bool)
	GetAllObserverMappers(ctx sdk.Context) (mappers []*observertypes.ObserverMapper)
	SetBallot(ctx sdk.Context, ballot *observertypes.Ballot)
	GetBallot(ctx sdk.Context, index string) (val observertypes.Ballot, found bool)
	GetAllBallots(ctx sdk.Context) (voters []*observertypes.Ballot)
	GetParams(ctx sdk.Context) (params observertypes.Params)
	GetCoreParamsByChainID(ctx sdk.Context, chainID int64) (params *observertypes.CoreParams, found bool)
	GetNodeAccount(ctx sdk.Context, address string) (nodeAccount observertypes.NodeAccount, found bool)
	GetAllNodeAccount(ctx sdk.Context) (nodeAccounts []observertypes.NodeAccount)
	SetNodeAccount(ctx sdk.Context, nodeAccount observertypes.NodeAccount)
	IsInboundEnabled(ctx sdk.Context) (found bool)
	GetCrosschainFlags(ctx sdk.Context) (val observertypes.CrosschainFlags, found bool)
	GetKeygen(ctx sdk.Context) (val observertypes.Keygen, found bool)
	SetKeygen(ctx sdk.Context, keygen observertypes.Keygen)
	SetCrosschainFlags(ctx sdk.Context, crosschainFlags observertypes.CrosschainFlags)
	SetLastObserverCount(ctx sdk.Context, lbc *observertypes.LastObserverCount)
	AddVoteToBallot(ctx sdk.Context, ballot observertypes.Ballot, address string, observationType observertypes.VoteType) (observertypes.Ballot, error)
	CheckIfFinalizingVote(ctx sdk.Context, ballot observertypes.Ballot) (observertypes.Ballot, bool)
	IsAuthorized(ctx sdk.Context, address string, chain *common.Chain) bool
	FindBallot(ctx sdk.Context, index string, chain *common.Chain, observationType observertypes.ObservationType) (ballot observertypes.Ballot, isNew bool, err error)
	AddBallotToList(ctx sdk.Context, ballot observertypes.Ballot)
	GetBlockHeader(ctx sdk.Context, hash []byte) (val common.BlockHeader, found bool)
	CheckIfTssPubkeyHasBeenGenerated(ctx sdk.Context, tssPubkey string) (observertypes.TSS, bool)
	GetPreviousTSS(ctx sdk.Context) (val observertypes.TSS, found bool)
	GetAllTSS(ctx sdk.Context) (list []observertypes.TSS)
	GetTSS(ctx sdk.Context) (val observertypes.TSS, found bool)
	SetTSS(ctx sdk.Context, tss observertypes.TSS)
	SetTSSHistory(ctx sdk.Context, tss observertypes.TSS)
	GetTssAddress(goCtx context.Context, req *observertypes.QueryGetTssAddressRequest) (*observertypes.QueryGetTssAddressResponse, error)

	SetFundMigrator(ctx sdk.Context, fm observertypes.TssFundMigratorInfo)
	GetFundMigrator(ctx sdk.Context, chainID int64) (val observertypes.TssFundMigratorInfo, found bool)
	GetAllTssFundMigrators(ctx sdk.Context) (fms []observertypes.TssFundMigratorInfo)
	RemoveAllExistingMigrators(ctx sdk.Context)
	SetChainNonces(ctx sdk.Context, chainNonces observertypes.ChainNonces)
	GetChainNonces(ctx sdk.Context, index string) (val observertypes.ChainNonces, found bool)
	RemoveChainNonces(ctx sdk.Context, index string)
	GetAllChainNonces(ctx sdk.Context) (list []observertypes.ChainNonces)
	SetNonceToCctx(ctx sdk.Context, nonceToCctx observertypes.NonceToCctx)
	GetNonceToCctx(ctx sdk.Context, tss string, chainID int64, nonce int64) (val observertypes.NonceToCctx, found bool)
	RemoveNonceToCctx(ctx sdk.Context, cctx observertypes.NonceToCctx)
	GetAllPendingNonces(ctx sdk.Context) (list []observertypes.PendingNonces, err error)
	GetPendingNonces(ctx sdk.Context, tss string, chainID int64) (val observertypes.PendingNonces, found bool)
	SetPendingNonces(ctx sdk.Context, pendingNonces observertypes.PendingNonces)
	SetTssAndUpdateNonce(ctx sdk.Context, tss observertypes.TSS)
	RemoveFromPendingNonces(ctx sdk.Context, tss string, chainID int64, nonce int64)
	GetAllNonceToCctx(ctx sdk.Context) (list []observertypes.NonceToCctx)

	VoteOnInboundBallot(
		ctx sdk.Context,
		senderChainID int64,
		receiverChainID int64,
		coinType common.CoinType,
		voter string,
		ballotIndex string,
		inTxHash string,
	) (bool, error)
}

type FungibleKeeper interface {
	GetForeignCoins(ctx sdk.Context, zrc20Addr string) (val fungibletypes.ForeignCoins, found bool)
	GetAllForeignCoins(ctx sdk.Context) (list []fungibletypes.ForeignCoins)
	SetForeignCoins(ctx sdk.Context, foreignCoins fungibletypes.ForeignCoins)
	GetAllForeignCoinsForChain(ctx sdk.Context, foreignChainID int64) (list []fungibletypes.ForeignCoins)
	GetForeignCoinFromAsset(ctx sdk.Context, asset string, chainID int64) (fungibletypes.ForeignCoins, bool)
	GetGasCoinForForeignCoin(ctx sdk.Context, chainID int64) (fungibletypes.ForeignCoins, bool)
	GetSystemContract(ctx sdk.Context) (val fungibletypes.SystemContract, found bool)
	QuerySystemContractGasCoinZRC20(ctx sdk.Context, chainID *big.Int) (eth.Address, error)
	GetUniswapV2Router02Address(ctx sdk.Context) (eth.Address, error)
	QueryUniswapV2RouterGetZetaAmountsIn(ctx sdk.Context, amountOut *big.Int, outZRC4 eth.Address) (*big.Int, error)
	QueryUniswapV2RouterGetZRC4AmountsIn(ctx sdk.Context, amountOut *big.Int, inZRC4 eth.Address) (*big.Int, error)
	QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn(ctx sdk.Context, amountOut *big.Int, inZRC4, outZRC4 eth.Address) (*big.Int, error)
	QueryGasLimit(ctx sdk.Context, contract eth.Address) (*big.Int, error)
	QueryProtocolFlatFee(ctx sdk.Context, contract eth.Address) (*big.Int, error)
	SetGasPrice(ctx sdk.Context, chainID *big.Int, gasPrice *big.Int) (uint64, error)
	DepositCoinZeta(ctx sdk.Context, to eth.Address, amount *big.Int) error
	DepositZRC20(
		ctx sdk.Context,
		contract eth.Address,
		to eth.Address,
		amount *big.Int,
	) (*evmtypes.MsgEthereumTxResponse, error)
	ZRC20DepositAndCallContract(
		ctx sdk.Context,
		from []byte,
		to eth.Address,
		amount *big.Int,
		senderChainID int64,
		data []byte,
		coinType common.CoinType,
		asset string,
	) (*evmtypes.MsgEthereumTxResponse, bool, error)
	CallUniswapV2RouterSwapExactTokensForTokens(
		ctx sdk.Context,
		sender eth.Address,
		to eth.Address,
		amountIn *big.Int,
		inZRC4,
		outZRC4 eth.Address,
		noEthereumTxEvent bool,
	) (ret []*big.Int, err error)
	CallUniswapV2RouterSwapExactETHForToken(
		ctx sdk.Context,
		sender eth.Address,
		to eth.Address,
		amountIn *big.Int,
		outZRC4 eth.Address,
		noEthereumTxEvent bool,
	) ([]*big.Int, error)
	CallZRC20Burn(ctx sdk.Context, sender eth.Address, zrc20address eth.Address, amount *big.Int, noEthereumTxEvent bool) error
	CallZRC20Approve(
		ctx sdk.Context,
		owner eth.Address,
		zrc20address eth.Address,
		spender eth.Address,
		amount *big.Int,
		noEthereumTxEvent bool,
	) error
	DeployZRC20Contract(
		ctx sdk.Context,
		name, symbol string,
		decimals uint8,
		chainID int64,
		coinType common.CoinType,
		erc20Contract string,
		gasLimit *big.Int,
	) (eth.Address, error)
	FundGasStabilityPool(ctx sdk.Context, chainID int64, amount *big.Int) error
	WithdrawFromGasStabilityPool(ctx sdk.Context, chainID int64, amount *big.Int) error
}
