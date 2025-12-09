package txserver

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	evidencetypes "cosmossdk.io/x/evidence/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cenkalti/backoff/v4"
	abci "github.com/cometbft/cometbft/abci/types"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/evm/crypto/hd"
	cosmosevmtypes "github.com/cosmos/evm/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/samber/lo"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/retry"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/node/x/lightclient/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// SystemContractAddresses contains the addresses of the system contracts deployed
type SystemContractAddresses struct {
	UniswapV2FactoryAddr, UniswapV2RouterAddr, ZEVMConnectorAddr, WZETAAddr, ERC20zrc20Addr string
}

// ZRC20Deployment configures deployment of ZRC20 contracts
type ZRC20Deployment struct {
	ERC20Addr common.Address
	SPLAddr   *solana.PublicKey // if nil - no SPL ZRC20 is deployed
}

// ZRC20Addresses contains the addresses of deployed ZRC20 contracts
type ZRC20Addresses struct {
	ERC20ZRC20Addr common.Address
	SPLZRC20Addr   common.Address
}

// EmissionsPoolAddress is the address of the emissions pool
// This address is constant for all networks because it is derived from emissions name
const EmissionsPoolAddress = "zeta1w43fn2ze2wyhu5hfmegr6vp52c3dgn0srdgymy"

// ZetaTxServer is a ZetaChain tx server for E2E test
type ZetaTxServer struct {
	ctx               context.Context
	clientCtx         client.Context
	txFactory         tx.Factory
	name              []string
	mnemonic          []string
	address           []string
	blockTimeout      time.Duration
	authorityClient   authoritytypes.QueryClient
	validatorsKeyring keyring.Keyring
}

// NewZetaTxServer returns a new TxServer with provided account
func NewZetaTxServer(rpcAddr string, names []string, privateKeys []string, chainID string) (*ZetaTxServer, error) {
	ctx := context.Background()

	if len(names) == 0 {
		return nil, errors.New("no account provided")
	}

	if len(names) != len(privateKeys) {
		return nil, errors.New("invalid names and privateKeys")
	}

	// initialize rpc and check status
	rpc, err := rpchttp.New(rpcAddr, "/websocket")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rpc: %s", err.Error())
	}
	if _, err = rpc.Status(ctx); err != nil {
		return nil, fmt.Errorf("failed to query rpc: %s", err.Error())
	}

	zetachain, err := chains.ZetaChainFromCosmosChainID(chainID)
	if err != nil {
		return nil, err
	}

	// initialize codec
	cdc, reg := newCodec(uint64(zetachain.ChainId)) //#nosec G115 chain id won't exceed uint64

	// initialize keyring
	kr := keyring.NewInMemory(cdc, hd.EthSecp256k1Option())

	addresses := make([]string, 0, len(names))

	// create accounts
	for i := range names {
		err = kr.ImportPrivKeyHex(names[i], privateKeys[i], string(hd.EthSecp256k1Type))
		if err != nil {
			return nil, fmt.Errorf("failed to create account: %w", err)
		}
		r, err := kr.Key(names[i])
		if err != nil {
			return nil, fmt.Errorf("failed to get account key: %w", err)
		}
		accAddr, err := r.GetAddress()
		if err != nil {
			return nil, fmt.Errorf("failed to get account address: %w", err)
		}

		addresses = append(addresses, accAddr.String())
	}

	// initialize validators keyring. The validator keys have been copied to this location in the docker image.
	// Refer contrib/localnet/orchestrator/Dockerfile.fastbuild for more details
	validatorsKeyring, err := keyring.New(
		"validatorKeys",
		keyring.BackendTest,
		"/root/.zetacored/",
		nil,
		cdc,
		hd.EthSecp256k1Option(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator keyring: %w", err)
	}

	clientCtx := newContext(rpc, cdc, reg, kr, chainID)
	txf := newFactory(clientCtx)

	return &ZetaTxServer{
		ctx:               ctx,
		clientCtx:         clientCtx,
		txFactory:         txf,
		name:              names,
		address:           addresses,
		blockTimeout:      2 * time.Minute,
		validatorsKeyring: validatorsKeyring,
	}, nil
}

// UpdateKeyring updates the keyring being used by the ZetaTxServer
// This returns a copy of the ZetaTxServer with the updated keyring
func (zts *ZetaTxServer) UpdateKeyring(kr keyring.Keyring) ZetaTxServer {
	ztsUpdated := *zts
	ztsUpdated.clientCtx = ztsUpdated.clientCtx.WithKeyring(kr)
	ztsUpdated.txFactory = newFactory(ztsUpdated.clientCtx)
	return ztsUpdated
}

// GetAccountName returns the account name from the given index
// returns empty string if index is out of bound, caller should handle error
func (zts *ZetaTxServer) GetAccountName(index int) string {
	if index >= len(zts.name) {
		return ""
	}
	return zts.name[index]
}

// GetCodec returns the codec for the ZetaTxServer
func (zts *ZetaTxServer) GetCodec() codec.Codec {
	return zts.clientCtx.Codec
}

// GetAccountAddress returns the account address from the given index
// returns empty string if index is out of bound, error should be handled by caller
func (zts *ZetaTxServer) GetAccountAddress(index int) string {
	if index >= len(zts.address) {
		return ""
	}
	return zts.address[index]
}

// GetAccountAddressFromName returns the account address from the given name
func (zts *ZetaTxServer) GetAccountAddressFromName(name string) (string, error) {
	acc, err := zts.clientCtx.Keyring.Key(name)
	if err != nil {
		return "", err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

// MustGetAccountAddressFromName returns the account address from the given name.It panics on error
// and should be used in tests only
func (zts *ZetaTxServer) MustGetAccountAddressFromName(name string) string {
	acc, err := zts.clientCtx.Keyring.Key(name)
	if err != nil {
		panic(err)
	}
	addr, err := acc.GetAddress()
	if err != nil {
		panic(err)
	}
	return addr.String()
}

// GetValidatorsKeyring returns the validators keyring
func (zts *ZetaTxServer) GetValidatorsKeyring() keyring.Keyring {
	return zts.validatorsKeyring
}

// GetAllAccountAddress returns all account addresses
func (zts *ZetaTxServer) GetAllAccountAddress() []string {
	return zts.address
}

// GetAccountMnemonic returns the account name from the given index
// returns empty string if index is out of bound, error should be handled by caller
func (zts *ZetaTxServer) GetAccountMnemonic(index int) string {
	if index >= len(zts.mnemonic) {
		return ""
	}
	return zts.mnemonic[index]
}

// GasLimit returns the gas limit of the tx factory
func (zts *ZetaTxServer) GasLimit() uint64 {
	return zts.txFactory.Gas()
}

// SetGasLimit sets the gas limit of the tx factory
func (zts *ZetaTxServer) SetGasLimit(gas uint64) {
	zts.txFactory = zts.txFactory.WithGas(gas)
}

// Fee returns the fee of the tx factory
func (zts *ZetaTxServer) Fee() sdktypes.Coins {
	return zts.txFactory.Fees()
}

// SetFee sets the fee of the tx factory
func (zts *ZetaTxServer) SetFee(fee sdktypes.Coins) {
	zts.txFactory = zts.txFactory.WithFees(fee.String())
}

// BroadcastTx broadcasts a tx to ZetaChain with the provided msg from the account
// and waiting for blockTime for tx to be included in the block
func (zts *ZetaTxServer) BroadcastTx(account string, msgs ...sdktypes.Msg) (*sdktypes.TxResponse, error) {
	bo := backoff.NewConstantBackOff(5 * time.Second)
	boWithMaxRetries := backoff.WithMaxRetries(bo, 5)

	return retry.DoTypedWithBackoff(func() (*sdktypes.TxResponse, error) {
		// Find number and sequence and set it
		acc, err := zts.clientCtx.Keyring.Key(account)
		if err != nil {
			return nil, err
		}

		addr, err := acc.GetAddress()
		if err != nil {
			return nil, err
		}

		accountNumber, accountSeq, err := zts.clientCtx.AccountRetriever.GetAccountNumberSequence(zts.clientCtx, addr)
		if err != nil {
			return nil, retry.Retry(err)
		}

		zts.txFactory = zts.txFactory.WithAccountNumber(accountNumber).WithSequence(accountSeq)

		txBuilder, err := zts.txFactory.BuildUnsignedTx(msgs...)
		if err != nil {
			return nil, retry.Retry(err)
		}

		txBuilder.SetGasLimit(zts.txFactory.Gas())
		txBuilder.SetFeeAmount(zts.txFactory.Fees())

		// Sign tx
		err = tx.Sign(zts.ctx, zts.txFactory, account, txBuilder, true)
		if err != nil {
			return nil, retry.Retry(err)
		}

		txBytes, err := zts.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
		if err != nil {
			return nil, retry.Retry(err)
		}

		result, err := broadcastWithBlockTimeout(zts, txBytes)
		if err != nil {
			return nil, retry.Retry(err)
		}

		return result, nil
	}, boWithMaxRetries)
}

func broadcastWithBlockTimeout(zts *ZetaTxServer, txBytes []byte) (*sdktypes.TxResponse, error) {
	res, err := zts.clientCtx.BroadcastTx(txBytes)
	if err != nil {
		if res == nil {
			return nil, err
		}
		return &sdktypes.TxResponse{
			Code:      res.Code,
			Codespace: res.Codespace,
			TxHash:    res.TxHash,
		}, err
	}
	if res.Code != 0 {
		return res, fmt.Errorf("broadcast failed: %s", res.RawLog)
	}

	exitAfter := time.After(zts.blockTimeout)
	hash, err := hex.DecodeString(res.TxHash)
	if err != nil {
		return nil, err
	}
	for {
		select {
		case <-exitAfter:
			return nil, fmt.Errorf("timed out after waiting for tx to get included in the block: %d", zts.blockTimeout)
		case <-time.After(time.Millisecond * 100):
			resTx, err := zts.clientCtx.Client.Tx(zts.ctx, hash, false)
			if err == nil {
				return mkTxResult(zts.ctx, zts.clientCtx, resTx)
			}
		}
	}
}

func mkTxResult(
	ctx context.Context,
	clientCtx client.Context,
	resTx *coretypes.ResultTx,
) (*sdktypes.TxResponse, error) {
	if resTx.TxResult.Code != 0 {
		return nil, fmt.Errorf("tx failed: %s", resTx.TxResult.Log)
	}
	txb, err := clientCtx.TxConfig.TxDecoder()(resTx.Tx)
	if err != nil {
		return nil, err
	}
	p, ok := txb.(intoAny)
	if !ok {
		return nil, fmt.Errorf("expecting a type implementing intoAny, got: %T", txb)
	}
	resBlock, err := clientCtx.Client.Block(ctx, &resTx.Height)
	if err != nil {
		return nil, err
	}
	return sdktypes.NewResponseResultTx(resTx, p.AsAny(), resBlock.Block.Time.Format(time.RFC3339)), nil
}

type intoAny interface {
	AsAny() *codectypes.Any
}

// UpdateChainParams updates the chain params for the given chain ID
// It takes into account that the required policy get updated in v28 operational -> admin
func (zts *ZetaTxServer) UpdateChainParams(chainParams *observertypes.ChainParams) error {
	_, err := zts.BroadcastTx(utils.AdminPolicyName, observertypes.NewMsgUpdateChainParams(
		zts.MustGetAccountAddressFromName(utils.AdminPolicyName),
		chainParams,
	))
	return err
}

// EnableHeaderVerification enables the header verification for the given chain IDs
func (zts *ZetaTxServer) EnableHeaderVerification(account string, chainIDList []int64) error {
	// retrieve account
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return err
	}

	_, err = zts.BroadcastTx(account, lightclienttypes.NewMsgEnableHeaderVerification(
		addr.String(),
		chainIDList,
	))
	return err
}

// UpdateGatewayAddress updates the gateway address
func (zts *ZetaTxServer) UpdateGatewayAddress(account, gatewayAddr string) error {
	// retrieve account
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return err
	}

	_, err = zts.BroadcastTx(account, fungibletypes.NewMsgUpdateGatewayContract(
		addr.String(),
		gatewayAddr,
	))

	return err
}

// DeploySystemContracts deploys the system contracts
// returns the addresses of uniswap factory, router
func (zts *ZetaTxServer) DeploySystemContracts(
	accountOperational, accountAdmin string,
) (SystemContractAddresses, error) {
	// retrieve account
	accOperational, err := zts.clientCtx.Keyring.Key(accountOperational)
	if err != nil {
		return SystemContractAddresses{}, err
	}
	addrOperational, err := accOperational.GetAddress()
	if err != nil {
		return SystemContractAddresses{}, err
	}
	accAdmin, err := zts.clientCtx.Keyring.Key(accountAdmin)
	if err != nil {
		return SystemContractAddresses{}, err
	}
	addrAdmin, err := accAdmin.GetAddress()
	if err != nil {
		return SystemContractAddresses{}, err
	}

	// deploy new system contracts
	res, err := zts.BroadcastTx(accountOperational, fungibletypes.NewMsgDeploySystemContracts(addrOperational.String()))
	if err != nil {
		return SystemContractAddresses{}, fmt.Errorf("failed to deploy system contracts: %s", err.Error())
	}

	deployedEvent, ok := EventOfType[*fungibletypes.EventSystemContractsDeployed](res.Events)
	if !ok {
		return SystemContractAddresses{}, fmt.Errorf("no EventSystemContractsDeployed in %s", res.TxHash)
	}

	// get system contract
	_, err = zts.BroadcastTx(
		accountAdmin,
		fungibletypes.NewMsgUpdateSystemContract(addrAdmin.String(), deployedEvent.SystemContract),
	)
	if err != nil {
		return SystemContractAddresses{}, fmt.Errorf("failed to set system contract: %s", err.Error())
	}

	return SystemContractAddresses{
		UniswapV2FactoryAddr: deployedEvent.UniswapV2Factory,
		UniswapV2RouterAddr:  deployedEvent.UniswapV2Router,
		ZEVMConnectorAddr:    deployedEvent.ConnectorZevm,
		WZETAAddr:            deployedEvent.Wzeta,
	}, nil
}

// DeployZRC20s deploys the ZRC20 contracts
// returns the addresses of erc20 and spl zrc20
func (zts *ZetaTxServer) DeployZRC20s(
	zrc20Deployment ZRC20Deployment,
	skipChain func(chainID int64) bool,
) (*ZRC20Addresses, error) {
	// retrieve account
	accOperational, err := zts.clientCtx.Keyring.Key(utils.OperationalPolicyName)
	if err != nil {
		return nil, err
	}
	addrOperational, err := accOperational.GetAddress()
	if err != nil {
		return nil, err
	}
	accAdmin, err := zts.clientCtx.Keyring.Key(utils.AdminPolicyName)
	if err != nil {
		return nil, err
	}
	addrAdmin, err := accAdmin.GetAddress()
	if err != nil {
		return nil, err
	}

	// authorization for deploying new ZRC20 has changed from accountOperational to accountAdmin in v19
	// we use this query to check the current authorization for the message
	// if pre v19 the query is not implement and authorization is operational
	deployerAccount := utils.AdminPolicyName
	deployerAddr := addrAdmin.String()
	authorization, preV19, err := zts.fetchMessagePermissions(&fungibletypes.MsgDeployFungibleCoinZRC20{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch message permissions: %s", err.Error())
	}
	if preV19 || authorization == authoritytypes.PolicyType_groupOperational {
		deployerAccount = utils.OperationalPolicyName
		deployerAddr = addrOperational.String()
	}

	deployMsgs := []*fungibletypes.MsgDeployFungibleCoinZRC20{
		fungibletypes.NewMsgDeployFungibleCoinZRC20(
			deployerAddr,
			"",
			chains.GoerliLocalnet.ChainId,
			18,
			"ETH",
			"gETH",
			coin.CoinType_Gas,
			100000,
			nil,
		),
		fungibletypes.NewMsgDeployFungibleCoinZRC20(
			deployerAddr,
			"",
			chains.BitcoinRegtest.ChainId,
			8,
			"BTC",
			"tBTC",
			coin.CoinType_Gas,
			100000,
			nil,
		),
		fungibletypes.NewMsgDeployFungibleCoinZRC20(
			deployerAddr,
			"",
			chains.SolanaLocalnet.ChainId,
			9,
			"Solana",
			"SOL",
			coin.CoinType_Gas,
			100000,
			nil,
		),
		fungibletypes.NewMsgDeployFungibleCoinZRC20(
			deployerAddr,
			"",
			chains.TONLocalnet.ChainId,
			9,
			"TON",
			"TON",
			coin.CoinType_Gas,
			100_000,
			nil,
		),
		fungibletypes.NewMsgDeployFungibleCoinZRC20(
			deployerAddr,
			zrc20Deployment.ERC20Addr.Hex(),
			chains.GoerliLocalnet.ChainId,
			6,
			"USDT",
			"USDT",
			coin.CoinType_ERC20,
			100000,
			nil,
		),
	}

	if zrc20Deployment.SPLAddr != nil {
		deployMsgs = append(deployMsgs, fungibletypes.NewMsgDeployFungibleCoinZRC20(
			deployerAddr,
			zrc20Deployment.SPLAddr.String(),
			chains.SolanaLocalnet.ChainId,
			9,
			"USDT",
			"USDT",
			coin.CoinType_ERC20,
			100000,
			nil,
		))
	}

	// apply skipChain filter and convert to sdk.Msg
	deployMsgsIface := lo.FilterMap(
		deployMsgs,
		func(msg *fungibletypes.MsgDeployFungibleCoinZRC20, _ int) (sdktypes.Msg, bool) {
			if skipChain(msg.ForeignChainId) {
				return nil, false
			}
			return msg, true
		},
	)

	// update the gas limit as this is a gas intensive operation
	previousGasLimit := zts.GasLimit()
	zts.SetGasLimit(20_000_000)
	previousFee := zts.Fee()
	parsedFees, err := sdktypes.ParseCoinsNormalized("200000000000000000azeta")
	if err != nil {
		return nil, fmt.Errorf("parse fee: %w", err)
	}
	zts.SetFee(parsedFees)

	res, err := zts.BroadcastTx(deployerAccount, deployMsgsIface...)
	if err != nil {
		return nil, fmt.Errorf("deploy zrc20s: %w", err)
	}

	// restore previous gas limit
	zts.SetGasLimit(previousGasLimit)
	zts.SetFee(previousFee)

	deployedEvents, ok := EventsOfType[*fungibletypes.EventZRC20Deployed](res.Events)
	if !ok {
		return nil, fmt.Errorf("no EventZRC20Deployed in %s", res.TxHash)
	}

	zrc20Addrs := lo.Map(deployedEvents, func(ev *fungibletypes.EventZRC20Deployed, _ int) string {
		return ev.Contract
	})

	err = zts.InitializeLiquidityCaps(zrc20Addrs...)
	if err != nil {
		return nil, fmt.Errorf("initialize liquidity cap: %w", err)
	}

	// find erc20 zrc20
	erc20zrc20, ok := lo.Find(deployedEvents, func(ev *fungibletypes.EventZRC20Deployed) bool {
		return ev.ChainId == chains.GoerliLocalnet.ChainId && ev.CoinType == coin.CoinType_ERC20
	})
	if !ok {
		return nil, fmt.Errorf("unable to find erc20 zrc20")
	}

	// find spl zrc20
	splzrc20Addr := common.Address{}
	if zrc20Deployment.SPLAddr != nil {
		splzrc20, ok := lo.Find(deployedEvents, func(ev *fungibletypes.EventZRC20Deployed) bool {
			return ev.ChainId == chains.SolanaLocalnet.ChainId && ev.CoinType == coin.CoinType_ERC20
		})
		if !ok {
			return nil, fmt.Errorf("unable to find spl zrc20")
		}

		splzrc20Addr = common.HexToAddress(splzrc20.Contract)
	}

	return &ZRC20Addresses{
		ERC20ZRC20Addr: common.HexToAddress(erc20zrc20.Contract),
		SPLZRC20Addr:   splzrc20Addr,
	}, nil
}

// FundEmissionsPool funds the emissions pool with the given amount
func (zts *ZetaTxServer) FundEmissionsPool(account string, amount *big.Int) error {
	// retrieve account
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return err
	}

	// retrieve account address
	emissionPoolAccAddr, err := sdktypes.AccAddressFromBech32(EmissionsPoolAddress)
	if err != nil {
		return err
	}

	// convert amount
	amountInt := sdkmath.NewIntFromBigInt(amount)

	// fund emissions pool
	_, err = zts.BroadcastTx(account, banktypes.NewMsgSend(
		addr,
		emissionPoolAccAddr,
		sdktypes.NewCoins(sdktypes.NewCoin(config.BaseDenom, amountInt)),
	))
	return err
}

func (zts *ZetaTxServer) WithdrawAllEmissions(withdrawAmount sdkmath.Int, account, observer string) error {
	// retrieve account
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return fmt.Errorf("failed to get withdrawer account: %w", err)
	}
	withdrawerAddress, err := acc.GetAddress()
	if err != nil {
		return fmt.Errorf("failed to get withdrawer account address: %w", err)
	}

	msg := emissionstypes.MsgWithdrawEmission{
		Creator: observer,
		Amount:  withdrawAmount,
	}

	authzMessage := authz.NewMsgExec(withdrawerAddress, []sdktypes.Msg{&msg})

	_, err = zts.BroadcastTx(account, &authzMessage)
	return err
}

// UpdateKeygen sets a new keygen height . The new height is the current height + 30
func (zts *ZetaTxServer) UpdateKeygen(height int64) error {
	keygenHeight := height + 30
	_, err := zts.BroadcastTx(zts.GetAccountName(0), observertypes.NewMsgUpdateKeygen(
		zts.GetAccountAddress(0),
		keygenHeight,
	))
	return err
}

func (zts *ZetaTxServer) RemoveObserver(observer string) error {
	account := zts.MustGetAccountAddressFromName(utils.AdminPolicyName)
	_, err := zts.BroadcastTx(utils.AdminPolicyName,
		observertypes.NewMsgRemoveObserver(
			account,
			observer,
		))
	return err
}

// SetAuthorityClient sets the authority client
func (zts *ZetaTxServer) SetAuthorityClient(authorityClient authoritytypes.QueryClient) {
	zts.authorityClient = authorityClient
}

// InitializeLiquidityCaps initializes the liquidity cap for the given coin with a large value
func (zts *ZetaTxServer) InitializeLiquidityCaps(zrc20s ...string) error {
	liquidityCap := sdkmath.NewUint(1e18).MulUint64(1e12)

	msgs := make([]sdktypes.Msg, len(zrc20s))
	for i, zrc20 := range zrc20s {
		msgs[i] = fungibletypes.NewMsgUpdateZRC20LiquidityCap(
			zts.MustGetAccountAddressFromName(utils.OperationalPolicyName),
			zrc20,
			liquidityCap,
		)
	}
	_, err := zts.BroadcastTx(utils.OperationalPolicyName, msgs...)
	return err
}

// fetchMessagePermissions fetches the message permissions for a given message
// return a bool preV19 to indicate the node is preV19 and the query doesn't exist
func (zts *ZetaTxServer) fetchMessagePermissions(msg sdktypes.Msg) (authoritytypes.PolicyType, bool, error) {
	msgURL := sdktypes.MsgTypeURL(msg)

	res, err := zts.authorityClient.Authorization(zts.ctx, &authoritytypes.QueryAuthorizationRequest{
		MsgUrl: msgURL,
	})

	// check if error is unknown method
	if err != nil {
		if strings.Contains(err.Error(), "unknown method") {
			return authoritytypes.PolicyType_groupOperational, true, nil
		}
		return authoritytypes.PolicyType_groupOperational, false, err
	}

	return res.Authorization.AuthorizedPolicy, false, nil
}

// newCodec returns the codec for msg server
func newCodec(evmChainID uint64) (*codec.ProtoCodec, codectypes.InterfaceRegistry) {
	encodingConfig := app.MakeEncodingConfig(evmChainID)
	interfaceRegistry := encodingConfig.InterfaceRegistry
	cdc := codec.NewProtoCodec(interfaceRegistry)

	sdktypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	authz.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	stakingtypes.RegisterInterfaces(interfaceRegistry)
	slashingtypes.RegisterInterfaces(interfaceRegistry)
	upgradetypes.RegisterInterfaces(interfaceRegistry)
	distrtypes.RegisterInterfaces(interfaceRegistry)
	evidencetypes.RegisterInterfaces(interfaceRegistry)
	crisistypes.RegisterInterfaces(interfaceRegistry)
	evmtypes.RegisterInterfaces(interfaceRegistry)
	cosmosevmtypes.RegisterInterfaces(interfaceRegistry)
	crosschaintypes.RegisterInterfaces(interfaceRegistry)
	emissionstypes.RegisterInterfaces(interfaceRegistry)
	fungibletypes.RegisterInterfaces(interfaceRegistry)
	observertypes.RegisterInterfaces(interfaceRegistry)
	lightclienttypes.RegisterInterfaces(interfaceRegistry)
	authoritytypes.RegisterInterfaces(interfaceRegistry)

	return cdc, interfaceRegistry
}

// newContext returns the client context for msg server
func newContext(
	rpc *rpchttp.HTTP,
	cdc *codec.ProtoCodec,
	reg codectypes.InterfaceRegistry,
	kr keyring.Keyring,
	chainID string,
) client.Context {
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	return client.Context{}.
		WithChainID(chainID).
		WithInterfaceRegistry(reg).
		WithCodec(cdc).
		WithTxConfig(txConfig).
		WithLegacyAmino(codec.NewLegacyAmino()).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithBroadcastMode(flags.BroadcastSync).
		WithClient(rpc).
		WithSkipConfirmation(true).
		WithFromName("creator").
		WithFromAddress(sdktypes.AccAddress{}).
		WithKeyring(kr).
		WithAccountRetriever(authtypes.AccountRetriever{})
}

// newFactory returns the tx factory for msg server
func newFactory(clientCtx client.Context) tx.Factory {
	return tx.Factory{}.
		WithChainID(clientCtx.ChainID).
		WithKeybase(clientCtx.Keyring).
		WithGas(10000000).
		WithGasAdjustment(1).
		WithSignMode(signing.SignMode_SIGN_MODE_UNSPECIFIED).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithTxConfig(clientCtx.TxConfig).
		WithFees("100000000000000000azeta")
}

// EventsOfType gets events of a specified type
func EventsOfType[T proto.Message](events []abci.Event) ([]T, bool) {
	var filteredEvents []T
	for _, ev := range events {
		pEvent, err := sdktypes.ParseTypedEvent(ev)
		if err != nil {
			continue
		}
		if typedEvent, ok := pEvent.(T); ok {
			filteredEvents = append(filteredEvents, typedEvent)
		}
	}
	return filteredEvents, len(filteredEvents) > 0
}

// EventOfType gets one event of a specific type
func EventOfType[T proto.Message](events []abci.Event) (T, bool) {
	var event T
	for _, ev := range events {
		pEvent, err := sdktypes.ParseTypedEvent(ev)
		if err != nil {
			continue
		}
		if typedEvent, ok := pEvent.(T); ok {
			return typedEvent, true
		}
	}
	return event, false
}
