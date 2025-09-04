package runner

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/pkg/coin"
	zetacrypto "github.com/zeta-chain/node/pkg/crypto"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const (
	// ZRC20InitialSupply is the initial supply of the ZRC20 token
	ZRC20SOLInitialSupply = 100000000

	// SolanaPDAInitialBalance is the initial balance (in lamports) of the gateway PDA account
	SolanaPDAInitialBalance = 1447680
)

type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Response struct {
	Amount Amount `json:"amount"`
}

func (r *E2ERunner) VerifyAccounting(testLegacy bool) {
	r.Logger.Print("Verifying accounting")
	r.checkETHTSSBalance()
	r.checkERC20TSSBalance()
	r.checkZETATSSBalance(testLegacy)
	r.CheckBTCTSSBalance()
	r.checkProtocolBalance()
}

func (r *E2ERunner) checkProtocolBalance() {
	balance := r.checkProtocolAddressBalance(config.BaseDenom)
	require.True(r, balance.IsZero())
}

func (r *E2ERunner) checkETHTSSBalance() {
	allTssAddress, err := r.ObserverClient.TssHistory(r.Ctx, &observertypes.QueryTssHistoryRequest{})
	require.NoError(r, err)

	tssTotalBalance := big.NewInt(0)

	for _, tssAddress := range allTssAddress.TssList {
		evmAddress, err := r.ObserverClient.GetTssAddressByFinalizedHeight(
			r.Ctx,
			&observertypes.QueryGetTssAddressByFinalizedHeightRequest{
				FinalizedZetaHeight: tssAddress.FinalizedZetaHeight,
			},
		)
		if err != nil {
			continue
		}

		tssBal, err := r.EVMClient.BalanceAt(r.Ctx, common.HexToAddress(evmAddress.Eth), nil)
		if err != nil {
			continue
		}
		tssTotalBalance.Add(tssTotalBalance, tssBal)
	}

	zrc20Supply, err := r.ETHZRC20.TotalSupply(&bind.CallOpts{})
	require.NoError(r, err)

	require.GreaterOrEqual(
		r,
		tssTotalBalance.Cmp(zrc20Supply),
		0,
		"ETH: TSS balance (%d) < ZRC20 TotalSupply (%d) ",
		tssTotalBalance,
		zrc20Supply,
	)
	r.Logger.Info("ETH: TSS balance (%d) >= ZRC20 TotalSupply (%d)", tssTotalBalance, zrc20Supply)
}

func (r *E2ERunner) CheckBTCTSSBalance() {
	allTssAddress, err := r.ObserverClient.TssHistory(r.Ctx, &observertypes.QueryTssHistoryRequest{})
	require.NoError(r, err)

	tssTotalBalance := float64(0)

	for _, tssAddress := range allTssAddress.TssList {
		btcTssAddress, err := zetacrypto.GetTSSAddrBTC(tssAddress.TssPubkey, r.BitcoinParams)
		if err != nil {
			continue
		}
		utxos, err := r.BtcRPCClient.ListUnspent(r.Ctx)
		if err != nil {
			continue
		}
		for _, utxo := range utxos {
			if utxo.Address == btcTssAddress {
				tssTotalBalance += utxo.Amount
			}
		}
	}

	zrc20Supply, err := r.BTCZRC20.TotalSupply(&bind.CallOpts{})
	require.NoError(r, err)

	// check the balance in TSS is greater than the total supply on ZetaChain
	// the amount minted to initialize the pool is subtracted from the total supply
	// #nosec G701 test - always in range
	require.GreaterOrEqual(r, int64(tssTotalBalance*1e8), zrc20Supply.Int64()-10000000,
		"BTC: TSS Balance (%d) < ZRC20 TotalSupply (%d)",
		int64(tssTotalBalance*1e8),
		zrc20Supply.Int64()-10000000,
	)
	// #nosec G115 test - always in range
	r.Logger.Info(
		"BTC: Balance (%d) >= ZRC20 TotalSupply (%d)",
		int64(tssTotalBalance*1e8),
		zrc20Supply.Int64()-10000000,
	)
}

// CheckSolanaTSSBalance compares the gateway PDA balance with the total supply of the SOL ZRC20 on ZetaChain
func (r *E2ERunner) CheckSolanaTSSBalance() {
	zrc20Supply, err := r.SOLZRC20.TotalSupply(&bind.CallOpts{})
	require.NoError(r, err)

	// get PDA received amount
	pda := r.ComputePdaAddress()
	balance, err := r.SolanaClient.GetBalance(r.Ctx, pda, rpc.CommitmentConfirmed)
	require.NoError(r, err)
	pdaReceivedAmount := balance.Value - SolanaPDAInitialBalance

	// the SOL balance in gateway PDA must not be less than the total supply on ZetaChain
	// the amount minted to initialize the pool is subtracted from the total supply
	// #nosec G115 test - always in range
	require.GreaterOrEqual(r, pdaReceivedAmount, zrc20Supply.Uint64()-ZRC20SOLInitialSupply,
		"SOL: Gateway PDA Received (%d) < ZRC20 TotalSupply (%d)",
		pdaReceivedAmount,
		zrc20Supply.Uint64()-ZRC20SOLInitialSupply,
	)
	// #nosec G115 test - always in range
	r.Logger.Info(
		"SOL: Gateway PDA Received (%d) >= ZRC20 TotalSupply (%d)",
		pdaReceivedAmount,
		zrc20Supply.Int64()-ZRC20SOLInitialSupply,
	)
}

// CheckSUITSSBalance checks the TSS balance on Sui against the ZRC20 total supply
func (r *E2ERunner) CheckSUITSSBalance() {
	gatewayBalance, err := r.SuiGetGatewaySUIBalance()
	require.NoError(r, err, "failed to get SUI balance for Sui gateway: %w", err)

	zrc20Supply, err := r.SUIZRC20.TotalSupply(&bind.CallOpts{})
	require.NoError(r, err)

	// Subtract 0.1 SUI to take in consideration the 0.1 SUI minted in the gas pool
	// TODO: a proper implementation is to implement a donate method in Sui Contract and use it to donate 0.1 SUI
	// https://github.com/zeta-chain/protocol-contracts-sui/issues/58
	zrc20Supply = zrc20Supply.Sub(zrc20Supply, big.NewInt(100000000))

	require.GreaterOrEqual(
		r,
		gatewayBalance.Cmp(zrc20Supply),
		0,
		"SUI: TSS balance (%d) < ZRC20 TotalSupply (%d) ",
		gatewayBalance,
		zrc20Supply,
	)
	r.Logger.Info("SUI: TSS balance (%d) >= ZRC20 TotalSupply (%d)", gatewayBalance, zrc20Supply)
}

func (r *E2ERunner) checkERC20TSSBalance() {
	custodyBalance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.ERC20CustodyAddr)
	require.NoError(r, err)

	erc20zrc20Supply, err := r.ERC20ZRC20.TotalSupply(&bind.CallOpts{})
	require.NoError(r, err)

	require.GreaterOrEqual(
		r,
		custodyBalance.Cmp(erc20zrc20Supply),
		0,
		"ERC20: custody balance (%d) < ZRC20 TotalSupply (%d) ",
		custodyBalance,
		erc20zrc20Supply,
	)
	r.Logger.Info("ERC20: TSS balance (%d) >= ERC20 ZRC20 TotalSupply (%d)", custodyBalance, erc20zrc20Supply)
}

func (r *E2ERunner) checkZETATSSBalance(testLegacy bool) {
	// zetaMintedPoolSetup is the amount of Zeta minted to add liquidy to a gas token pool when setting it up
	zetaMintedPoolSetup := r.fetchZetaMintedByGasPoolCreations()
	zetaTokensMintedDuringSetup := r.fetchTokensMintedAtGenesis().Add(zetaMintedPoolSetup)

	zetaLockedLegacyConnector, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.ConnectorEthAddr)
	require.NoError(r, err, "BalanceOf failed for legacy connector")
	zetaLockedConnectorNative, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.ConnectorNativeAddr)
	require.NoError(r, err, "BalanceOf failed for new connector")

	zetaSupply := r.fetchZetaSupply()
	// Subtract the amount of Zeta minted during setup from the total supply to fetch the number of zeta tokens created by e2e tests
	zetaSupply = zetaSupply.Sub(zetaTokensMintedDuringSetup)
	abortedAmount := r.fetchAbortedAmount()
	zetaMinted := zetaSupply.Add(abortedAmount)

	// Fetch the total ZETA locked in the legacy and native connectors, remove 1 ZETA if this is a legacy test
	// This is to account for the tokens from the test TestZetaDepositRestricted
	// Related issue : https://github.com/zeta-chain/node/issues/4057
	zetaLocked := sdkmath.NewIntFromBigInt(big.NewInt(0).Add(zetaLockedLegacyConnector, zetaLockedConnectorNative))
	if testLegacy {
		oneZeta := sdkmath.NewInt(1e18)
		zetaLocked = zetaLocked.Sub(oneZeta)
	}

	require.True(
		r,
		zetaMinted.Equal(zetaLocked),
		"ZETA: Connector balance (%s) != ZETA Minted (%s) [ ZETA TotalSupply (%s) + AbortedAmount (%s) ]",
		zetaLocked.String(),
		zetaMinted.String(),
		zetaSupply.String(),
		abortedAmount.String(),
	)
}

// fetchZetaMintedByGasPoolCreations retrieves the total amount of ZETA tokens minted by gas pool creations.
// https://github.com/zeta-chain/node/blob/854040f044d198a0453a7b9245d544debd9da055/x/fungible/keeper/gas_coin_and_pool.go#L88
func (r *E2ERunner) fetchZetaMintedByGasPoolCreations() sdkmath.Int {
	zetaPerGasPool := sdkmath.NewInt(1e17) // 0.1 ZETA per gas pool creation
	// Get the protocol address balance
	res, err := r.Clients.Zetacore.Fungible.ForeignCoinsAll(r.Ctx, &fungibletypes.QueryAllForeignCoinsRequest{})
	require.NoError(r, err)
	require.NotNil(r, res)

	gasCoinCount := 0
	for _, foreignCoin := range res.ForeignCoins {
		if foreignCoin.CoinType == coin.CoinType_Gas {
			gasCoinCount++
		}
	}
	if res.ForeignCoins == nil || gasCoinCount == 0 {
		return sdkmath.ZeroInt()
	}
	return zetaPerGasPool.Mul(sdkmath.NewInt(int64(gasCoinCount)))
}

// fetchTokensMintedAtGenesis retrieves the total amount of ZETA tokens minttend from the genesis file.
func (r *E2ERunner) fetchTokensMintedAtGenesis() sdkmath.Int {
	genesisFilePath := "/root/.zetacored/data/genesis.json"
	_, genesis, err := genutiltypes.GenesisStateFromGenFile(genesisFilePath)
	require.NoError(r, err, "failed to get genesis state from file: %s", genesisFilePath)

	appState, err := genutiltypes.GenesisStateFromAppGenesis(genesis)
	require.NoError(r, err, "failed to get app genesis state from genesis")

	bankStateBz, ok := appState[banktypes.ModuleName]
	require.True(r, ok, "bank genesis state is missing")

	zevmChainID, err := r.ZEVMClient.ChainID(r.Ctx)
	require.NoError(r, err, "failed to get ZetaChain ID from ZEVM client")

	r.Logger.Info("ZetaChain ID: %d", zevmChainID.Uint64())
	cdc := app.MakeEncodingConfig(zevmChainID.Uint64()).Codec

	bankState := new(banktypes.GenesisState)
	err = cdc.UnmarshalJSON(bankStateBz, bankState)
	require.NoError(r, err)

	return bankState.Supply.AmountOf(config.BaseDenom)
}

// fetchZetaSupply retrieves the total supply of ZETA tokens from the Zetacore bank module.
func (r *E2ERunner) fetchZetaSupply() sdkmath.Int {
	res, err := r.Clients.Zetacore.Bank.SupplyOf(r.Ctx, &banktypes.QuerySupplyOfRequest{
		Denom: config.BaseDenom,
	})
	require.NoError(r, err)
	require.NotNil(r, res)

	return res.Amount.Amount
}

// fetchAbortedAmount retrieves the total amount of aborted ZETA tokens from the Zetacore crosschain module.
func (r *E2ERunner) fetchAbortedAmount() sdkmath.Int {
	res, err := r.Clients.Zetacore.Crosschain.ZetaAccounting(r.Ctx, &crosschaintypes.QueryZetaAccountingRequest{})
	require.NoError(r, err)
	require.NotNil(r, res)

	abortedAmount, ok := sdkmath.NewIntFromString(res.GetAbortedZetaAmount())
	require.True(r, ok, "failed to parse aborted ZETA amount")
	return abortedAmount
}
func (r *E2ERunner) checkProtocolAddressBalance(denom string) sdkmath.Int {
	res, err := r.BankClient.Balance(r.Ctx, &banktypes.QueryBalanceRequest{
		Address: fungibletypes.ModuleAddress.String(),
		Denom:   denom,
	})
	require.NoError(r, err, "failed to get protocol address balance")
	return res.GetBalance().Amount
}
