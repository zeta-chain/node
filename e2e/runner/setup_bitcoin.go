package runner

import (
	"encoding/hex"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
)

func (runner *E2ERunner) SetupBitcoinAccount(initNetwork bool) {
	runner.Logger.Print("⚙️ setting up Bitcoin account")
	startTime := time.Now()
	defer func() {
		runner.Logger.Print("✅ Bitcoin account setup in %s\n", time.Since(startTime))
	}()

	_, err := runner.BtcRPCClient.CreateWallet(runner.Name, rpcclient.WithCreateWalletBlank())
	if err != nil {
		if !strings.Contains(err.Error(), "Database already exists") {
			panic(err)
		}
	}

	runner.SetBtcAddress(runner.Name, true)

	if initNetwork {
		// import the TSS address
		err = runner.BtcRPCClient.ImportAddress(runner.BTCTSSAddress.EncodeAddress())
		if err != nil {
			panic(err)
		}

		// mine some blocks to get some BTC into the deployer address
		_, err = runner.BtcRPCClient.GenerateToAddress(101, runner.BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}

		_, err = runner.BtcRPCClient.GenerateToAddress(4, runner.BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}
	}
}

// GetBtcAddress returns the BTC address of the deployer from its EVM private key
func (runner *E2ERunner) GetBtcAddress() (string, string, error) {
	skBytes, err := hex.DecodeString(runner.DeployerPrivateKey)
	if err != nil {
		return "", "", err
	}

	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)
	privkeyWIF, err := btcutil.NewWIF(sk, runner.BitcoinParams, true)
	if err != nil {
		return "", "", err
	}

	address, err := btcutil.NewAddressWitnessPubKeyHash(
		btcutil.Hash160(privkeyWIF.SerializePubKey()),
		runner.BitcoinParams,
	)
	if err != nil {
		return "", "", err
	}

	// return the string representation of the address
	return address.EncodeAddress(), privkeyWIF.String(), nil
}

// SetBtcAddress imports the deployer's private key into the Bitcoin node
func (runner *E2ERunner) SetBtcAddress(name string, rescan bool) {
	skBytes, err := hex.DecodeString(runner.DeployerPrivateKey)
	if err != nil {
		panic(err)
	}

	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)
	privkeyWIF, err := btcutil.NewWIF(sk, runner.BitcoinParams, true)
	if err != nil {
		panic(err)
	}

	if rescan {
		err = runner.BtcRPCClient.ImportPrivKeyRescan(privkeyWIF, name, true)
		if err != nil {
			panic(err)
		}
	}

	runner.BTCDeployerAddress, err = btcutil.NewAddressWitnessPubKeyHash(
		btcutil.Hash160(privkeyWIF.PrivKey.PubKey().SerializeCompressed()),
		runner.BitcoinParams,
	)
	if err != nil {
		panic(err)
	}

	runner.Logger.Info("BTCDeployerAddress: %s", runner.BTCDeployerAddress.EncodeAddress())
}
