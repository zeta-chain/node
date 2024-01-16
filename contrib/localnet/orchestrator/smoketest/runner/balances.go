package runner

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
	"strings"
)

// AccountBalances is a struct that contains the balances of the accounts used in the smoke test
type AccountBalances struct {
	ZetaETH   *big.Int
	ZetaZETA  *big.Int
	ZetaERC20 *big.Int
	EvmETH    *big.Int
	EvmZETA   *big.Int
	EvmERC20  *big.Int
}

// AccountBalancesDiff is a struct that contains the difference in the balances of the accounts used in the smoke test
type AccountBalancesDiff struct {
	ETH   *big.Int
	ZETA  *big.Int
	ERC20 *big.Int
}

// GetAccountBalances returns the account balances of the accounts used in the smoke test
func (sm *SmokeTestRunner) GetAccountBalances() (AccountBalances, error) {
	// zevm
	zetaZeta, err := sm.ZevmClient.BalanceAt(sm.Ctx, sm.DeployerAddress, nil)
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

	// evm
	ethEth, err := sm.GoerliClient.BalanceAt(sm.Ctx, sm.DeployerAddress, nil)
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

	return AccountBalances{
		ZetaETH:   zetaEth,
		ZetaZETA:  zetaZeta,
		ZetaERC20: zetaErc20,
		EvmETH:    ethEth,
		EvmZETA:   evmZeta,
		EvmERC20:  evmErc20,
	}, nil
}

// PrintAccountBalances shows the account balances of the accounts used in the smoke test
// Note: USDT is mentioned as erc20 here because we want to show the balance of any erc20 contract
func (sm *SmokeTestRunner) PrintAccountBalances(balances AccountBalances) {
	sm.Logger.Print(" ---💰 Account info %s ---", sm.DeployerAddress.Hex())

	// zevm
	sm.Logger.Print("ZetaChain:")
	sm.Logger.Print("* ZETA balance:  %s", balances.ZetaZETA.String())
	sm.Logger.Print("* ETH balance:   %s", balances.ZetaETH.String())
	sm.Logger.Print("* ERC20 balance: %s", balances.ZetaERC20.String())

	// evm
	sm.Logger.Print("Ethereum:")
	sm.Logger.Print("* ZETA balance:  %s", balances.EvmZETA.String())
	sm.Logger.Print("* ETH balance:   %s", balances.EvmETH.String())
	sm.Logger.Print("* ERC20 balance: %s", balances.EvmERC20.String())

	return
}

// PrintTotalDiff shows the difference in the account balances of the accounts used in the e2e test from two balances structs
func (sm *SmokeTestRunner) PrintTotalDiff(accoutBalancesDiff AccountBalancesDiff) {
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

// GetAccountBalancesDiff returns the difference in the account balances of the accounts used in the smoke test
func GetAccountBalancesDiff(balancesBefore, balancesAfter AccountBalances) AccountBalancesDiff {
	balancesBeforeZeta := big.NewInt(0).Add(balancesBefore.ZetaETH, balancesBefore.ZetaZETA)
	balancesBeforeEth := big.NewInt(0).Add(balancesBefore.EvmETH, balancesBefore.EvmZETA)
	balancesBeforeErc20 := big.NewInt(0).Add(balancesBefore.ZetaERC20, balancesBefore.EvmERC20)

	balancesAfterZeta := big.NewInt(0).Add(balancesAfter.ZetaETH, balancesAfter.ZetaZETA)
	balancesAfterEth := big.NewInt(0).Add(balancesAfter.EvmETH, balancesAfter.EvmZETA)
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
