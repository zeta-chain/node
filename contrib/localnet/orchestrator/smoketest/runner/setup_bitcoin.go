package runner

import (
	"encoding/hex"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
)

func (sm *SmokeTestRunner) SetupBitcoinAccount() {
	sm.Logger.Print("⚙️ setting up Bitcoin account")
	startTime := time.Now()
	defer func() {
		sm.Logger.Info("Bitcoin account setup took %s\n", time.Since(startTime))
	}()

	btc := sm.BtcRPCClient
	_, err := btc.CreateWallet(sm.Name, rpcclient.WithCreateWalletBlank())
	if err != nil {
		if !strings.Contains(err.Error(), "Database already exists") {
			panic(err)
		}
	}

	sm.setBtcAddress()

	err = btc.ImportAddress(sm.BTCTSSAddress.EncodeAddress())
	if err != nil {
		panic(err)
	}

	// mine some blocks to get some BTC into the deployer address
	_, err = btc.GenerateToAddress(101, sm.BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}

	//bal, err := btc.GetBalance("*")
	//if err != nil {
	//	panic(err)
	//}
	_, err = btc.GenerateToAddress(4, sm.BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	//bal, err = btc.GetBalance("*")
	//if err != nil {
	//	panic(err)
	//}
	//sm.Logger.Info("balance: %f", bal.ToBTC())
	//
	//bals, err := btc.GetBalances()
	//if err != nil {
	//	panic(err)
	//}
	//sm.Logger.Info("balances: ")
	//sm.Logger.Info("  mine (Deployer): %+v\n", bals.Mine)
	//if bals.WatchOnly != nil {
	//	sm.Logger.Info("  watchonly (TSSAddress): %+v", bals.WatchOnly)
	//}
	//sm.Logger.Info("  TSS Address: %s", sm.BTCTSSAddress.EncodeAddress())
	//go func() {
	//	// keep bitcoin chain going
	//	for {
	//		_, err = btc.GenerateToAddress(4, sm.BTCDeployerAddress, nil)
	//		if err != nil {
	//			sm.Logger.Info(err.Error())
	//		}
	//		time.Sleep(5 * time.Second)
	//	}
	//}()
}

// setBtcAddress
func (sm *SmokeTestRunner) setBtcAddress() {
	skBytes, err := hex.DecodeString(sm.DeployerPrivateKey)
	if err != nil {
		panic(err)
	}

	// TODO: support non regtest chain
	// https://github.com/zeta-chain/node/issues/1482
	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)
	privkeyWIF, err := btcutil.NewWIF(sk, &chaincfg.RegressionNetParams, true)
	if err != nil {
		panic(err)
	}

	err = sm.BtcRPCClient.ImportPrivKeyRescan(privkeyWIF, sm.Name, true)
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

	sm.Logger.Info("BTCDeployerAddress: %s", sm.BTCDeployerAddress.EncodeAddress())
}
