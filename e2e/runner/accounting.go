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

func (r *E2ERunner) CheckZRC20ReserveAndSupply() error {
	r.Logger.Info("Checking ZRC20 Reserve and Supply")
	if err := r.checkEthTSSBalance(); err != nil {
		return err
	}
	if err := r.checkERC20TSSBalance(); err != nil {
		return err
	}
	return r.checkZetaTSSBalance()
}

func (r *E2ERunner) checkEthTSSBalance() error {
	tssBal, err := r.EVMClient.BalanceAt(r.Ctx, r.TSSAddress, nil)
	if err != nil {
		return err
	}
	zrc20Supply, err := r.ETHZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}
	if tssBal.Cmp(zrc20Supply) < 0 {
		return fmt.Errorf("ETH: TSS balance (%d) < ZRC20 TotalSupply (%d) ", tssBal, zrc20Supply)
	}
	r.Logger.Info("ETH: TSS balance (%d) >= ZRC20 TotalSupply (%d)", tssBal, zrc20Supply)
	return nil
}

func (r *E2ERunner) CheckBtcTSSBalance() error {
	utxos, err := r.BtcRPCClient.ListUnspent()
	if err != nil {
		return err
	}
	var btcBalance float64
	for _, utxo := range utxos {
		if utxo.Address == r.BTCTSSAddress.EncodeAddress() {
			btcBalance += utxo.Amount
		}
	}

	zrc20Supply, err := r.BTCZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}

	// check the balance in TSS is greater than the total supply on ZetaChain
	// the amount minted to initialize the pool is subtracted from the total supply
	// #nosec G115 test - always in range
	if int64(btcBalance*1e8) < (zrc20Supply.Int64() - 10000000) {
		// #nosec G115 test - always in range
		return fmt.Errorf(
			"BTC: TSS Balance (%d) < ZRC20 TotalSupply (%d)",
			int64(btcBalance*1e8),
			zrc20Supply.Int64()-10000000,
		)
	}
	// #nosec G115 test - always in range
	r.Logger.Info(
		"BTC: Balance (%d) >= ZRC20 TotalSupply (%d)",
		int64(btcBalance*1e8),
		zrc20Supply.Int64()-10000000,
	)

	return nil
}

func (r *E2ERunner) checkERC20TSSBalance() error {
	erc20Balance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.ERC20CustodyAddr)
	if err != nil {
		return err
	}
	erc20zrc20Supply, err := r.ERC20ZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}
	if erc20Balance.Cmp(erc20zrc20Supply) < 0 {
		return fmt.Errorf("ERC20: TSS balance (%d) < ZRC20 TotalSupply (%d) ", erc20Balance, erc20zrc20Supply)
	}
	r.Logger.Info("ERC20: TSS balance (%d) >= ERC20 ZRC20 TotalSupply (%d)", erc20Balance, erc20zrc20Supply)
	return nil
}

func (r *E2ERunner) checkZetaTSSBalance() error {
	zetaLocked, err := r.ConnectorEth.GetLockedAmount(&bind.CallOpts{})
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
		r.Logger.Info(fmt.Sprintf("ZETA: TSS balance (%d) < ZRC20 TotalSupply (%d)", zetaLocked, zetaSupply))
	} else {
		r.Logger.Info("ZETA: TSS balance (%d) >= ZRC20 TotalSupply (%d)", zetaLocked, zetaSupply)
	}
	return nil
}
