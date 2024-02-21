package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Response struct {
	Amount Amount `json:"amount"`
}

func (runner *E2ERunner) CheckZRC20ReserveAndSupply() error {
	runner.Logger.Info("Checking ZRC20 Reserve and Supply")
	if err := runner.checkEthTSSBalance(); err != nil {
		return err
	}
	if err := runner.checkUsdtTSSBalance(); err != nil {
		return err
	}
	return runner.checkZetaTSSBalance()
}

func (runner *E2ERunner) checkEthTSSBalance() error {
	tssBal, err := runner.GoerliClient.BalanceAt(runner.Ctx, runner.TSSAddress, nil)
	if err != nil {
		return err
	}
	zrc20Supply, err := runner.ETHZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}
	if tssBal.Cmp(zrc20Supply) < 0 {
		return fmt.Errorf("ETH: TSS balance (%d) < ZRC20 TotalSupply (%d) ", tssBal, zrc20Supply)
	}
	runner.Logger.Info("ETH: TSS balance (%d) >= ZRC20 TotalSupply (%d)", tssBal, zrc20Supply)
	return nil
}

func (runner *E2ERunner) CheckBtcTSSBalance() error {
	utxos, err := runner.BtcRPCClient.ListUnspent()
	if err != nil {
		return err
	}
	var btcBalance float64
	for _, utxo := range utxos {
		if utxo.Address == runner.BTCTSSAddress.EncodeAddress() {
			btcBalance += utxo.Amount
		}
	}

	zrc20Supply, err := runner.BTCZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}

	// check the balance in TSS is greater than the total supply on ZetaChain
	// the amount minted to initialize the pool is subtracted from the total supply
	// #nosec G701 test - always in range
	if int64(btcBalance*1e8) < (zrc20Supply.Int64() - 10000000) {
		// #nosec G701 test - always in range
		return fmt.Errorf(
			"BTC: TSS Balance (%d) < ZRC20 TotalSupply (%d)",
			int64(btcBalance*1e8),
			zrc20Supply.Int64()-10000000,
		)
	}
	// #nosec G701 test - always in range
	runner.Logger.Info("BTC: Balance (%d) >= ZRC20 TotalSupply (%d)", int64(btcBalance*1e8), zrc20Supply.Int64()-10000000)

	return nil
}

func (runner *E2ERunner) checkUsdtTSSBalance() error {
	usdtBal, err := runner.USDTERC20.BalanceOf(&bind.CallOpts{}, runner.ERC20CustodyAddr)
	if err != nil {
		return err
	}
	zrc20Supply, err := runner.USDTZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}
	if usdtBal.Cmp(zrc20Supply) < 0 {
		return fmt.Errorf("USDT: TSS balance (%d) < ZRC20 TotalSupply (%d) ", usdtBal, zrc20Supply)
	}
	runner.Logger.Info("USDT: TSS balance (%d) >= ZRC20 TotalSupply (%d)", usdtBal, zrc20Supply)
	return nil
}

func (runner *E2ERunner) checkZetaTSSBalance() error {
	zetaLocked, err := runner.ConnectorEth.GetLockedAmount(&bind.CallOpts{})
	if err != nil {
		return err
	}
	resp, err := http.Get("http://zetacore0:1317/cosmos/bank/v1beta1/supply/by_denom?denom=azeta")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	zetaSupply, _ := big.NewInt(0).SetString(result.Amount.Amount, 10)
	if zetaLocked.Cmp(zetaSupply) < 0 {
		runner.Logger.Info(fmt.Sprintf("ZETA: TSS balance (%d) < ZRC20 TotalSupply (%d)", zetaLocked, zetaSupply))
	} else {
		runner.Logger.Info("ZETA: TSS balance (%d) >= ZRC20 TotalSupply (%d)", zetaLocked, zetaSupply)
	}
	return nil
}
