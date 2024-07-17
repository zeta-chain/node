package runner

import (
	"encoding/hex"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/require"
)

func (r *E2ERunner) AddTSSToNode() {
	r.Logger.Print("⚙️ add new tss to Bitcoin node")
	startTime := time.Now()
	defer func() {
		r.Logger.Print("✅ Bitcoin account setup in %s\n", time.Since(startTime))
	}()

	// import the TSS address
	err := r.BtcRPCClient.ImportAddress(r.BTCTSSAddress.EncodeAddress())
	require.NoError(r, err)

	// mine some blocks to get some BTC into the deployer address
	_, err = r.GenerateToAddressIfLocalBitcoin(101, r.BTCDeployerAddress)
	require.NoError(r, err)
}

func (r *E2ERunner) SetupBitcoinAccount(initNetwork bool) {
	r.Logger.Print("⚙️ setting up Bitcoin account")
	startTime := time.Now()
	defer func() {
		r.Logger.Print("✅ Bitcoin account setup in %s", time.Since(startTime))
	}()

	_, err := r.BtcRPCClient.CreateWallet(r.Name, rpcclient.WithCreateWalletBlank())
	if err != nil {
		require.ErrorContains(r, err, "Database already exists")
	}

	r.SetBtcAddress(r.Name, true)

	if initNetwork {
		// import the TSS address
		err = r.BtcRPCClient.ImportAddress(r.BTCTSSAddress.EncodeAddress())
		require.NoError(r, err)

		// mine some blocks to get some BTC into the deployer address
		_, err = r.GenerateToAddressIfLocalBitcoin(101, r.BTCDeployerAddress)
		require.NoError(r, err)

		_, err = r.GenerateToAddressIfLocalBitcoin(4, r.BTCDeployerAddress)
		require.NoError(r, err)
	}
}

// GetBtcAddress returns the BTC address of the deployer from its EVM private key
func (r *E2ERunner) GetBtcAddress() (string, string, error) {
	skBytes, err := hex.DecodeString(r.Account.RawPrivateKey.String())
	if err != nil {
		return "", "", err
	}

	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)
	privkeyWIF, err := btcutil.NewWIF(sk, r.BitcoinParams, true)
	if err != nil {
		return "", "", err
	}

	address, err := btcutil.NewAddressWitnessPubKeyHash(
		btcutil.Hash160(privkeyWIF.SerializePubKey()),
		r.BitcoinParams,
	)
	if err != nil {
		return "", "", err
	}

	// return the string representation of the address
	return address.EncodeAddress(), privkeyWIF.String(), nil
}

// SetBtcAddress imports the deployer's private key into the Bitcoin node
func (r *E2ERunner) SetBtcAddress(name string, rescan bool) {
	skBytes, err := hex.DecodeString(r.Account.RawPrivateKey.String())
	require.NoError(r, err)

	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)
	privkeyWIF, err := btcutil.NewWIF(sk, r.BitcoinParams, true)
	require.NoError(r, err)

	if rescan {
		err := r.BtcRPCClient.ImportPrivKeyRescan(privkeyWIF, name, true)
		require.NoError(r, err, "failed to execute ImportPrivKeyRescan")
	}

	r.BTCDeployerAddress, err = btcutil.NewAddressWitnessPubKeyHash(
		btcutil.Hash160(privkeyWIF.PrivKey.PubKey().SerializeCompressed()),
		r.BitcoinParams,
	)
	require.NoError(r, err)

	r.Logger.Info("BTCDeployerAddress: %s", r.BTCDeployerAddress.EncodeAddress())
}
