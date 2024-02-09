package runner

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// AccountBalances is a struct that contains the balances of the accounts used in the E2E test
type AccountBalances struct {
	ZetaETH   *big.Int
	ZetaZETA  *big.Int
	ZetaWZETA *big.Int
	ZetaERC20 *big.Int
	ZetaBTC   *big.Int
	EvmETH    *big.Int
	EvmZETA   *big.Int
	EvmERC20  *big.Int
	BtcBTC    string
}

// AccountBalancesDiff is a struct that contains the difference in the balances of the accounts used in the E2E test
type AccountBalancesDiff struct {
	ETH   *big.Int
	ZETA  *big.Int
	ERC20 *big.Int
}

// GetAccountBalances returns the account balances of the accounts used in the E2E test
func (sm *E2ERunner) GetAccountBalances(skipBTC bool) (AccountBalances, error) {
	// zevm
	zetaZeta, err := sm.ZevmClient.BalanceAt(sm.Ctx, sm.DeployerAddress, nil)
	if err != nil {
		return AccountBalances{}, err
	}
	zetaWZeta, err := sm.WZeta.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		return AccountBalances{}, err
	}
	zetaEth, err := sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		return AccountBalances{}, err
	}
	zetaErc20, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		return AccountBalances{}, err
	}
	zetaBtc, err := sm.BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		return AccountBalances{}, err
	}

	// evm
	evmEth, err := sm.GoerliClient.BalanceAt(sm.Ctx, sm.DeployerAddress, nil)
	if err != nil {
		return AccountBalances{}, err
	}
	evmZeta, err := sm.ZetaEth.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		return AccountBalances{}, err
	}
	evmErc20, err := sm.USDTERC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		return AccountBalances{}, err
	}

	// bitcoin
	var BtcBTC string
	if !skipBTC {
		if BtcBTC, err = sm.GetBitcoinBalance(); err != nil {
			return AccountBalances{}, err
		}
	}

	return AccountBalances{
		ZetaETH:   zetaEth,
		ZetaZETA:  zetaZeta,
		ZetaWZETA: zetaWZeta,
		ZetaERC20: zetaErc20,
		ZetaBTC:   zetaBtc,
		EvmETH:    evmEth,
		EvmZETA:   evmZeta,
		EvmERC20:  evmErc20,
		BtcBTC:    BtcBTC,
	}, nil
}

// GetBitcoinBalance returns the spendable BTC balance of the BTC address
func (sm *E2ERunner) GetBitcoinBalance() (string, error) {
	addr, _, err := sm.GetBtcAddress()
	if err != nil {
		return "", fmt.Errorf("failed to get BTC address: %w", err)
	}

	address, err := btcutil.DecodeAddress(addr, sm.BitcoinParams)
	if err != nil {
		return "", fmt.Errorf("failed to decode BTC address: %w", err)
	}

	unspentList, err := sm.BtcRPCClient.ListUnspentMinMaxAddresses(1, 9999999, []btcutil.Address{address})
	if err != nil {
		return "", fmt.Errorf("failed to list unspent: %w", err)
	}

	// calculate total amount
	var totalAmount btcutil.Amount
	for _, unspent := range unspentList {
		if unspent.Spendable {
			totalAmount += btcutil.Amount(unspent.Amount * 1e8)
		}
	}

	return totalAmount.String(), nil
}

// PrintAccountBalances shows the account balances of the accounts used in the E2E test
// Note: USDT is mentioned as erc20 here because we want to show the balance of any erc20 contract
func (sm *E2ERunner) PrintAccountBalances(balances AccountBalances) {
	sm.Logger.Print(" ---💰 Account info %s ---", sm.DeployerAddress.Hex())

	// zevm
	sm.Logger.Print("ZetaChain:")
	sm.Logger.Print("* ZETA balance:  %s", balances.ZetaZETA.String())
	sm.Logger.Print("* WZETA balance: %s", balances.ZetaWZETA.String())
	sm.Logger.Print("* ETH balance:   %s", balances.ZetaETH.String())
	sm.Logger.Print("* ERC20 balance: %s", balances.ZetaERC20.String())
	sm.Logger.Print("* BTC balance:   %s", balances.ZetaBTC.String())

	// evm
	sm.Logger.Print("EVM:")
	sm.Logger.Print("* ZETA balance:  %s", balances.EvmZETA.String())
	sm.Logger.Print("* ETH balance:   %s", balances.EvmETH.String())
	sm.Logger.Print("* ERC20 balance: %s", balances.EvmERC20.String())

	// bitcoin
	sm.Logger.Print("Bitcoin:")
	sm.Logger.Print("* BTC balance: %s", balances.BtcBTC)

	return
}

// PrintTotalDiff shows the difference in the account balances of the accounts used in the e2e test from two balances structs
func (sm *E2ERunner) PrintTotalDiff(accoutBalancesDiff AccountBalancesDiff) {
	sm.Logger.Print(" ---💰 Total gas spent ---")

	// show the value only if it is not zero
	if accoutBalancesDiff.ZETA.Cmp(big.NewInt(0)) != 0 {
		sm.Logger.Print("* ZETA spent:  %s", accoutBalancesDiff.ZETA.String())
	}
	if accoutBalancesDiff.ETH.Cmp(big.NewInt(0)) != 0 {
		sm.Logger.Print("* ETH spent:   %s", accoutBalancesDiff.ETH.String())
	}
	if accoutBalancesDiff.ERC20.Cmp(big.NewInt(0)) != 0 {
		sm.Logger.Print("* ERC20 spent: %s", accoutBalancesDiff.ERC20.String())
	}
}

// GetAccountBalancesDiff returns the difference in the account balances of the accounts used in the E2E test
func GetAccountBalancesDiff(balancesBefore, balancesAfter AccountBalances) AccountBalancesDiff {
	balancesBeforeZeta := big.NewInt(0).Add(balancesBefore.ZetaZETA, balancesBefore.EvmZETA)
	balancesBeforeEth := big.NewInt(0).Add(balancesBefore.ZetaETH, balancesBefore.EvmETH)
	balancesBeforeErc20 := big.NewInt(0).Add(balancesBefore.ZetaERC20, balancesBefore.EvmERC20)

	balancesAfterZeta := big.NewInt(0).Add(balancesAfter.ZetaZETA, balancesAfter.EvmZETA)
	balancesAfterEth := big.NewInt(0).Add(balancesAfter.ZetaETH, balancesAfter.EvmETH)
	balancesAfterErc20 := big.NewInt(0).Add(balancesAfter.ZetaERC20, balancesAfter.EvmERC20)

	diffZeta := big.NewInt(0).Sub(balancesBeforeZeta, balancesAfterZeta)
	diffEth := big.NewInt(0).Sub(balancesBeforeEth, balancesAfterEth)
	diffErc20 := big.NewInt(0).Sub(balancesBeforeErc20, balancesAfterErc20)

	return AccountBalancesDiff{
		ETH:   diffEth,
		ZETA:  diffZeta,
		ERC20: diffErc20,
	}
}

// formatBalances formats the AccountBalancesDiff into a one-liner string
func formatBalances(balances AccountBalancesDiff) string {
	parts := []string{}
	if balances.ETH != nil && balances.ETH.Cmp(big.NewInt(0)) > 0 {
		parts = append(parts, fmt.Sprintf("ETH:%s", balances.ETH.String()))
	}
	if balances.ZETA != nil && balances.ZETA.Cmp(big.NewInt(0)) > 0 {
		parts = append(parts, fmt.Sprintf("ZETA:%s", balances.ZETA.String()))
	}
	if balances.ERC20 != nil && balances.ERC20.Cmp(big.NewInt(0)) > 0 {
		parts = append(parts, fmt.Sprintf("ERC20:%s", balances.ERC20.String()))
	}
	return strings.Join(parts, ",")
}
