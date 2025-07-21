package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

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

func (r *E2ERunner) CheckZRC20BalanceAndSupply() {
	r.Logger.Print("Checking ZRC20 Balance vs. Supply")

	err := r.checkETHTSSBalance()
	require.NoError(r, err, "ETH balance check failed")

	err = r.checkERC20TSSBalance()
	require.NoError(r, err, "ERC20 balance check failed")

	r.checkZetaTSSBalance()

	err = r.CheckBTCTSSBalance()
	require.NoError(r, err, "BTC balance check failed")
}

func (r *E2ERunner) checkETHTSSBalance() error {
	allTssAddress, err := r.ObserverClient.TssHistory(r.Ctx, &observertypes.QueryTssHistoryRequest{})
	if err != nil {
		return err
	}

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
	if err != nil {
		return err
	}
	if tssTotalBalance.Cmp(zrc20Supply) < 0 {
		return fmt.Errorf("ETH: TSS balance (%d) < ZRC20 TotalSupply (%d) ", tssTotalBalance, zrc20Supply)
	}
	r.Logger.Info("ETH: TSS balance (%d) >= ZRC20 TotalSupply (%d)", tssTotalBalance, zrc20Supply)
	return nil
}

func (r *E2ERunner) CheckBTCTSSBalance() error {
	allTssAddress, err := r.ObserverClient.TssHistory(r.Ctx, &observertypes.QueryTssHistoryRequest{})
	if err != nil {
		return err
	}

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
	if err != nil {
		return err
	}

	// check the balance in TSS is greater than the total supply on ZetaChain
	// the amount minted to initialize the pool is subtracted from the total supply
	// #nosec G701 test - always in range
	if int64(tssTotalBalance*1e8) < (zrc20Supply.Int64() - 10000000) {
		// #nosec G701 test - always in range
		return fmt.Errorf(
			"BTC: TSS Balance (%d) < ZRC20 TotalSupply (%d)",
			int64(tssTotalBalance*1e8),
			zrc20Supply.Int64()-10000000,
		)
	}
	// #nosec G115 test - always in range
	r.Logger.Info(
		"BTC: Balance (%d) >= ZRC20 TotalSupply (%d)",
		int64(tssTotalBalance*1e8),
		zrc20Supply.Int64()-10000000,
	)

	return nil
}

// CheckSolanaTSSBalance compares the gateway PDA balance with the total supply of the SOL ZRC20 on ZetaChain
func (r *E2ERunner) CheckSolanaTSSBalance() error {
	zrc20Supply, err := r.SOLZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}

	// get PDA received amount
	pda := r.ComputePdaAddress()
	balance, err := r.SolanaClient.GetBalance(r.Ctx, pda, rpc.CommitmentConfirmed)
	require.NoError(r, err)
	pdaReceivedAmount := balance.Value - SolanaPDAInitialBalance

	// the SOL balance in gateway PDA must not be less than the total supply on ZetaChain
	// the amount minted to initialize the pool is subtracted from the total supply
	// #nosec G115 test - always in range
	if pdaReceivedAmount < (zrc20Supply.Uint64() - ZRC20SOLInitialSupply) {
		// #nosec G115 test - always in range
		return fmt.Errorf(
			"SOL: Gateway PDA Received (%d) < ZRC20 TotalSupply (%d)",
			pdaReceivedAmount,
			zrc20Supply.Uint64()-ZRC20SOLInitialSupply,
		)
	}
	// #nosec G115 test - always in range
	r.Logger.Info(
		"SOL: Gateway PDA Received (%d) >= ZRC20 TotalSupply (%d)",
		pdaReceivedAmount,
		zrc20Supply.Int64()-ZRC20SOLInitialSupply,
	)

	return nil
}

// CheckSUITSSBalance checks the TSS balance on Sui against the ZRC20 total supply
func (r *E2ERunner) CheckSUITSSBalance() error {
	gatewayBalance, err := r.SuiGetGatewaySUIBalance()
	if err != nil {
		return fmt.Errorf("failed to get SUI balance for Sui gateway: %w", err)
	}

	zrc20Supply, err := r.SUIZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}

	// subtract value from the gas stability pool because of the artificial minting bug
	// TODO: remove on the chain upgrade to v33
	gasStabiltiyPoolBalance, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, fungibletypes.GasStabilityPoolAddressEVM())
	if err != nil {
		return fmt.Errorf("failed to get SUI gas stability pool balance: %w", err)
	}
	zrc20Supply = zrc20Supply.Sub(zrc20Supply, gasStabiltiyPoolBalance)

	if gatewayBalance.Cmp(zrc20Supply) < 0 {
		return fmt.Errorf("SUI: TSS balance (%d) < ZRC20 TotalSupply (%d) ", gatewayBalance, zrc20Supply)
	}
	r.Logger.Info("SUI: TSS balance (%d) >= ZRC20 TotalSupply (%d)", gatewayBalance, zrc20Supply)
	return nil
}

func (r *E2ERunner) checkERC20TSSBalance() error {
	custodyBalance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.ERC20CustodyAddr)
	if err != nil {
		return err
	}

	erc20zrc20Supply, err := r.ERC20ZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}
	if custodyBalance.Cmp(erc20zrc20Supply) < 0 {
		return fmt.Errorf("ERC20: custody balance (%d) < ZRC20 TotalSupply (%d) ", custodyBalance, erc20zrc20Supply)
	}
	r.Logger.Info("ERC20: TSS balance (%d) >= ERC20 ZRC20 TotalSupply (%d)", custodyBalance, erc20zrc20Supply)
	return nil
}

func (r *E2ERunner) checkZetaTSSBalance() {
	zetaLockedLegacyConnector, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.ConnectorEthAddr)
	require.NoError(r, err, "BalanceOf failed for legacy connector")

	zetaLockedConnectorNative, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.ConnectorNativeAddr)
	require.NoError(r, err, "BalanceOf failed for new connector")

	zetaLocked := big.NewInt(0).Add(zetaLockedLegacyConnector, zetaLockedConnectorNative)

	resp, err := http.Get("http://zetacore0:1317/cosmos/bank/v1beta1/supply/by_denom?denom=azeta")
	require.NoError(r, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(r, err)
	var result Response
	err = json.Unmarshal(body, &result)
	require.NoError(r, err)
	zetaSupply, _ := big.NewInt(0).SetString(result.Amount.Amount, 10)
	if zetaLocked.Cmp(zetaSupply) < 0 {
		r.Logger.Info("ZETA: TSS balance (%d) < ZRC20 TotalSupply (%d)", zetaLocked, zetaSupply)
	} else {
		r.Logger.Info("ZETA: TSS balance (%d) >= ZRC20 TotalSupply (%d)", zetaLocked, zetaSupply)
	}
}
