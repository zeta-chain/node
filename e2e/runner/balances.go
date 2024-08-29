package runner

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/zrc20.sol"
)

var errNilZRC20 = errors.New("zrc20 contract is nil")

// AccountBalances is a struct that contains the balances of the accounts used in the E2E test
type AccountBalances struct {
	ZetaETH   *big.Int
	ZetaZETA  *big.Int
	ZetaWZETA *big.Int
	ZetaERC20 *big.Int
	ZetaBTC   *big.Int
	ZetaSOL   *big.Int
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

func (r *E2ERunner) getZRC20BalanceSafe(z *zrc20.ZRC20) (*big.Int, error) {
	if z == nil {
		return new(big.Int), errNilZRC20
	}
	return z.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
}

// GetAccountBalances returns the account balances of the accounts used in the E2E test
func (r *E2ERunner) GetAccountBalances(skipBTC bool) (AccountBalances, error) {
	// zevm
	zetaZeta, err := r.ZEVMClient.BalanceAt(r.Ctx, r.EVMAddress(), nil)
	if err != nil {
		return AccountBalances{}, err
	}
	zetaWZeta, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	if err != nil {
		return AccountBalances{}, err
	}
	zetaEth, err := r.getZRC20BalanceSafe(r.ETHZRC20)
	if err != nil {
		return AccountBalances{}, err
	}
	zetaErc20, err := r.getZRC20BalanceSafe(r.ERC20ZRC20)
	if err != nil {
		return AccountBalances{}, err
	}
	zetaBtc, err := r.getZRC20BalanceSafe(r.BTCZRC20)
	if err != nil {
		return AccountBalances{}, err
	}
	zetaSol, err := r.getZRC20BalanceSafe(r.SOLZRC20)
	if err != nil {
		r.Logger.Error("get SOL balance: %v", err)
	}

	// evm
	evmEth, err := r.EVMClient.BalanceAt(r.Ctx, r.EVMAddress(), nil)
	if err != nil {
		return AccountBalances{}, err
	}
	evmZeta, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	if err != nil {
		return AccountBalances{}, err
	}
	evmErc20, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	if err != nil {
		return AccountBalances{}, err
	}

	// bitcoin
	var BtcBTC string
	if !skipBTC {
		if BtcBTC, err = r.GetBitcoinBalance(); err != nil {
			return AccountBalances{}, err
		}
	}

	return AccountBalances{
		ZetaETH:   zetaEth,
		ZetaZETA:  zetaZeta,
		ZetaWZETA: zetaWZeta,
		ZetaERC20: zetaErc20,
		ZetaBTC:   zetaBtc,
		ZetaSOL:   zetaSol,
		EvmETH:    evmEth,
		EvmZETA:   evmZeta,
		EvmERC20:  evmErc20,
		BtcBTC:    BtcBTC,
	}, nil
}

// GetBitcoinBalance returns the spendable BTC balance of the BTC address
func (r *E2ERunner) GetBitcoinBalance() (string, error) {
	addr, _, err := r.GetBtcAddress()
	if err != nil {
		return "", fmt.Errorf("failed to get BTC address: %w", err)
	}

	address, err := btcutil.DecodeAddress(addr, r.BitcoinParams)
	if err != nil {
		return "", fmt.Errorf("failed to decode BTC address: %w", err)
	}

	total, err := r.GetBitcoinBalanceByAddress(address)
	if err != nil {
		return "", err
	}

	return total.String(), nil
}

// GetBitcoinBalanceByAddress get btc balance by address.
func (r *E2ERunner) GetBitcoinBalanceByAddress(address btcutil.Address) (btcutil.Amount, error) {
	unspentList, err := r.BtcRPCClient.ListUnspentMinMaxAddresses(1, 9999999, []btcutil.Address{address})
	if err != nil {
		return 0, errors.Wrap(err, "failed to list unspent")
	}

	var total btcutil.Amount
	for _, unspent := range unspentList {
		if unspent.Spendable {
			total += btcutil.Amount(unspent.Amount * 1e8)
		}
	}

	return total, nil
}

// PrintAccountBalances shows the account balances of the accounts used in the E2E test
// Note: USDT is mentioned as erc20 here because we want to show the balance of any erc20 contract
func (r *E2ERunner) PrintAccountBalances(balances AccountBalances) {
	r.Logger.Print(" ---💰 Account info %s ---", r.EVMAddress().Hex())

	// zevm
	r.Logger.Print("ZetaChain:")
	r.Logger.Print("* ZETA balance:  %s", balances.ZetaZETA.String())
	r.Logger.Print("* WZETA balance: %s", balances.ZetaWZETA.String())
	r.Logger.Print("* ETH balance:   %s", balances.ZetaETH.String())
	r.Logger.Print("* ERC20 balance: %s", balances.ZetaERC20.String())
	r.Logger.Print("* BTC balance:   %s", balances.ZetaBTC.String())

	// evm
	r.Logger.Print("EVM:")
	r.Logger.Print("* ZETA balance:  %s", balances.EvmZETA.String())
	r.Logger.Print("* ETH balance:   %s", balances.EvmETH.String())
	r.Logger.Print("* ERC20 balance: %s", balances.EvmERC20.String())

	// bitcoin
	r.Logger.Print("Bitcoin:")
	r.Logger.Print("* BTC balance: %s", balances.BtcBTC)

	// solana
	r.Logger.Print("Solana:")
	r.Logger.Print("* SOL balance: %s", balances.ZetaSOL.String())
}

// PrintTotalDiff shows the difference in the account balances of the accounts used in the e2e test from two balances structs
func (r *E2ERunner) PrintTotalDiff(accoutBalancesDiff AccountBalancesDiff) {
	r.Logger.Print(" ---💰 Total gas spent ---")

	// show the value only if it is not zero
	if accoutBalancesDiff.ZETA.Cmp(big.NewInt(0)) != 0 {
		r.Logger.Print("* ZETA spent:  %s", accoutBalancesDiff.ZETA.String())
	}
	if accoutBalancesDiff.ETH.Cmp(big.NewInt(0)) != 0 {
		r.Logger.Print("* ETH spent:   %s", accoutBalancesDiff.ETH.String())
	}
	if accoutBalancesDiff.ERC20.Cmp(big.NewInt(0)) != 0 {
		r.Logger.Print("* ERC20 spent: %s", accoutBalancesDiff.ERC20.String())
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
