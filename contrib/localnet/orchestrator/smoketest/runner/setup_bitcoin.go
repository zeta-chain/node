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

func (sm *SmokeTestRunner) SetupBitcoinAccount(initNetwork bool) {
	sm.Logger.Print("⚙️ setting up Bitcoin account")
	startTime := time.Now()
	defer func() {
		sm.Logger.Print("✅ Bitcoin account setup in %s\n", time.Since(startTime))
	}()

	_, err := sm.BtcRPCClient.CreateWallet(sm.Name, rpcclient.WithCreateWalletBlank())
	if err != nil {
		if !strings.Contains(err.Error(), "Database already exists") {
			panic(err)
		}
	}

	sm.setBtcAddress()

	if initNetwork {
		// import the TSS address
		err = sm.BtcRPCClient.ImportAddress(sm.BTCTSSAddress.EncodeAddress())
		if err != nil {
			panic(err)
		}

		// mine some blocks to get some BTC into the deployer address
		_, err = sm.BtcRPCClient.GenerateToAddress(101, sm.BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}

		_, err = sm.BtcRPCClient.GenerateToAddress(4, sm.BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}
	}
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

//// getBtcAddress returns the BTC address of the deployer from its EVM private key
//func (sm *SmokeTestRunner) getBtcAddress() (string, error) {
//	skBytes, err := hex.DecodeString(sm.DeployerPrivateKey)
//	if err != nil {
//		return "", err
//	}
//
//	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)
//	privkeyWIF, err := btcutil.NewWIF(sk, &chaincfg.RegressionNetParams, true)
//	if err != nil {
//		return "", err
//	}
//
//	return privkeyWIF.PrivKey.PubKey()., nil
//}
