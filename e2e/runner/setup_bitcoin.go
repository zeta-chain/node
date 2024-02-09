package runner

import (
	"encoding/hex"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
)

func (sm *E2ERunner) SetupBitcoinAccount(initNetwork bool) {
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

	sm.SetBtcAddress(sm.Name, true)

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

// GetBtcAddress returns the BTC address of the deployer from its EVM private key
func (sm *E2ERunner) GetBtcAddress() (string, string, error) {
	skBytes, err := hex.DecodeString(sm.DeployerPrivateKey)
	if err != nil {
		return "", "", err
	}

	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)
	privkeyWIF, err := btcutil.NewWIF(sk, sm.BitcoinParams, true)
	if err != nil {
		return "", "", err
	}

	address, err := btcutil.NewAddressWitnessPubKeyHash(
		btcutil.Hash160(privkeyWIF.SerializePubKey()),
		sm.BitcoinParams,
	)
	if err != nil {
		return "", "", err
	}

	// return the string representation of the address
	return address.EncodeAddress(), privkeyWIF.String(), nil
}

// SetBtcAddress imports the deployer's private key into the Bitcoin node
func (sm *E2ERunner) SetBtcAddress(name string, rescan bool) {
	skBytes, err := hex.DecodeString(sm.DeployerPrivateKey)
	if err != nil {
		panic(err)
	}

	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)
	privkeyWIF, err := btcutil.NewWIF(sk, sm.BitcoinParams, true)
	if err != nil {
		panic(err)
	}

	if rescan {
		err = sm.BtcRPCClient.ImportPrivKeyRescan(privkeyWIF, name, true)
		if err != nil {
			panic(err)
		}
	}

	sm.BTCDeployerAddress, err = btcutil.NewAddressWitnessPubKeyHash(
		btcutil.Hash160(privkeyWIF.PrivKey.PubKey().SerializeCompressed()),
		sm.BitcoinParams,
	)
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("BTCDeployerAddress: %s", sm.BTCDeployerAddress.EncodeAddress())
}
