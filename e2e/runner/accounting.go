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
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"

	zetacrypto "github.com/zeta-chain/node/pkg/crypto"
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
	r.Logger.Print("Checking ZRC20 Balance vs Supply")
	r.checkETHTSSBalance()
	r.checkERC20TSSBalance()
	r.checkZetaTSSBalance(testLegacy)
	r.CheckBTCTSSBalance()
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

	require.GreaterOrEqual(r, tssTotalBalance.Cmp(zrc20Supply), 0, "ETH: TSS balance (%d) < ZRC20 TotalSupply (%d) ", tssTotalBalance, zrc20Supply)
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

	// subtract value from the gas stability pool because of the artificial minting bug
	// TODO: remove on the chain upgrade to v33
	// https://github.com/zeta-chain/node/issues/4034
	gasStabiltiyPoolBalance, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, fungibletypes.GasStabilityPoolAddressEVM())
	require.NoError(r, err, "failed to get SUI gas stability pool balance: %w", err)
	zrc20Supply = zrc20Supply.Sub(zrc20Supply, gasStabiltiyPoolBalance)

	// Subtract 0.1 SUI to take in consideration the 0.1 SUI minted in the gas pool
	// TODO: a proper implementation is to implement a donate method in Sui Contract and use it to donate 0.1 SUI
	// https://github.com/zeta-chain/protocol-contracts-sui/issues/58
	zrc20Supply = zrc20Supply.Sub(zrc20Supply, big.NewInt(100000000))

	require.GreaterOrEqual(r, gatewayBalance.Cmp(zrc20Supply), 0, "SUI: TSS balance (%d) < ZRC20 TotalSupply (%d) ", gatewayBalance, zrc20Supply)
	r.Logger.Info("SUI: TSS balance (%d) >= ZRC20 TotalSupply (%d)", gatewayBalance, zrc20Supply)
}

func (r *E2ERunner) checkERC20TSSBalance() {
	custodyBalance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.ERC20CustodyAddr)
	require.NoError(r, err)

	erc20zrc20Supply, err := r.ERC20ZRC20.TotalSupply(&bind.CallOpts{})
	require.NoError(r, err)

	require.GreaterOrEqual(r, custodyBalance.Cmp(erc20zrc20Supply), 0, "ERC20: custody balance (%d) < ZRC20 TotalSupply (%d) ", custodyBalance, erc20zrc20Supply)
	r.Logger.Info("ERC20: TSS balance (%d) >= ERC20 ZRC20 TotalSupply (%d)", custodyBalance, erc20zrc20Supply)
}

func (r *E2ERunner) checkZetaTSSBalance(testLegacy bool) {
	unknownNumber, ok := sdkmath.NewIntFromString("400000000000000000")
	require.True(r, ok, "failed to parse unknown number for zeta tokens minted during setup")
	zetaTokensMintedDuringSetup := r.GetTokensMintedAtGenesis().Add(unknownNumber)

	zetaLockedLegacyConnector, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.ConnectorEthAddr)
	require.NoError(r, err, "BalanceOf failed for legacy connector")

	zetaLockedConnectorNative, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.ConnectorNativeAddr)
	require.NoError(r, err, "BalanceOf failed for new connector")

	zetaSupply := r.FetchZetaSupply()
	zetaSupply = zetaSupply.Sub(zetaTokensMintedDuringSetup)
	abortedAmount := r.FetchAbortedAmount()

	zetaMinted := zetaSupply.Add(abortedAmount)
	zetaLocked := sdkmath.NewIntFromBigInt(big.NewInt(0).Add(zetaLockedLegacyConnector, zetaLockedConnectorNative))
	if testLegacy {
		oneZeta := sdkmath.NewInt(1e18)
		zetaLocked = zetaLocked.Sub(oneZeta)
	}

	require.True(r, zetaMinted.Equal(zetaLocked), "ZETA: Connector balance (%s) != ZETA TotalSupply (%s) + AbortedAmount (%d)", zetaLocked.String(), zetaSupply.String(), abortedAmount.String())
}

func (r *E2ERunner) GetTokensMintedAtGenesis() sdkmath.Int {
	genesisFilePath := "/root/.zetacored/data/genesis.json"
	_, genesis, err := genutiltypes.GenesisStateFromGenFile(genesisFilePath)
	require.NoError(r, err, "failed to get genesis state from file: %s", genesisFilePath)

	appState, err := genutiltypes.GenesisStateFromAppGenesis(genesis)
	require.NoError(r, err, "failed to get app genesis state from genesis")

	bankStateBz, ok := appState[banktypes.ModuleName]
	require.True(r, ok, "bank genesis state is missing")
	cdc := app.MakeEncodingConfig().Codec

	bankState := new(banktypes.GenesisState)
	err = cdc.UnmarshalJSON(bankStateBz, bankState)
	require.NoError(r, err)

	return bankState.Supply.AmountOf(config.BaseDenom)
}

func (r *E2ERunner) FetchZetaSupply() sdkmath.Int {
	res, err := r.Clients.Zetacore.Bank.SupplyOf(r.Ctx, &banktypes.QuerySupplyOfRequest{
		Denom: config.BaseDenom,
	})
	require.NoError(r, err)
	require.NotNil(r, res)
	r.Logger.Print("ZetaSupply: %s", res.Amount.Amount.String())
	return res.Amount.Amount
}

func (r *E2ERunner) FetchAbortedAmount() sdkmath.Int {
	res, err := r.Clients.Zetacore.Crosschain.ZetaAccounting(r.Ctx, &crosschaintypes.QueryZetaAccountingRequest{})
	require.NoError(r, err)
	require.NotNil(r, res)

	abortedAmount, ok := sdkmath.NewIntFromString(res.GetAbortedZetaAmount())
	require.True(r, ok, "failed to parse aborted ZETA amount")
	return abortedAmount
}
