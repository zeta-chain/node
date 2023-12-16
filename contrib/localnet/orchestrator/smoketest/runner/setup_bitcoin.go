package runner

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

func (sm *SmokeTestRunner) SetupBitcoin() {
	utils.LoudPrintf("Setup Bitcoin\n")
	startTime := time.Now()
	defer func() {
		fmt.Printf("Bitcoin setup took %s\n", time.Since(startTime))
	}()

	btc := sm.BtcRPCClient
	_, err := btc.CreateWallet("smoketest", rpcclient.WithCreateWalletBlank())
	if err != nil {
		if !strings.Contains(err.Error(), "Database already exists") {
			panic(err)
		}
	}

	skBytes, err := hex.DecodeString(sm.DeployerPrivateKey)
	if err != nil {
		panic(err)
	}
	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)
	privkeyWIF, err := btcutil.NewWIF(sk, &chaincfg.RegressionNetParams, true)
	if err != nil {
		panic(err)
	}
	err = btc.ImportPrivKeyRescan(privkeyWIF, "deployer", true)
	if err != nil {
		panic(err)
	}
	sm.BTCDeployerAddress, err = btcutil.NewAddressWitnessPubKeyHash(
		btcutil.Hash160(privkeyWIF.PrivKey.PubKey().SerializeCompressed()),
		&chaincfg.RegressionNetParams,
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("BTCDeployerAddress: %s\n", sm.BTCDeployerAddress.EncodeAddress())

	err = btc.ImportAddress(sm.BTCTSSAddress.EncodeAddress())
	if err != nil {
		panic(err)
	}
	_, err = btc.GenerateToAddress(101, sm.BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	bal, err := btc.GetBalance("*")
	if err != nil {
		panic(err)
	}
	_, err = btc.GenerateToAddress(4, sm.BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	bal, err = btc.GetBalance("*")
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance: %f\n", bal.ToBTC())

	bals, err := btc.GetBalances()
	if err != nil {
		panic(err)
	}
	fmt.Printf("balances: \n")
	fmt.Printf("  mine (Deployer): %+v\n", bals.Mine)
	if bals.WatchOnly != nil {
		fmt.Printf("  watchonly (TSSAddress): %+v\n", bals.WatchOnly)
	}
	fmt.Printf("  TSS Address: %s\n", sm.BTCTSSAddress.EncodeAddress())
	go func() {
		// keep bitcoin chain going
		for {
			_, err = btc.GenerateToAddress(4, sm.BTCDeployerAddress, nil)
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(5 * time.Second)
		}
	}()
}
