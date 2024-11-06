package runner

import (
	"encoding/hex"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
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

// SetupBitcoinAccounts sets up the TSS account and deployer account
func (r *E2ERunner) SetupBitcoinAccounts(createWallet bool) {
	r.Logger.Info("⚙️ setting up Bitcoin account")
	startTime := time.Now()
	defer func() {
		r.Logger.Print("✅ Bitcoin account setup in %s", time.Since(startTime))
	}()

	// setup deployer account
	r.SetupBtcAddress(r.Name, createWallet)

	// import the TSS address to index TSS utxos and txs
	err := r.BtcRPCClient.ImportAddress(r.BTCTSSAddress.EncodeAddress())
	require.NoError(r, err)
	r.Logger.Info("⚙️ imported BTC TSSAddress: %s", r.BTCTSSAddress.EncodeAddress())

	// import deployer address to index deployer utxos and txs
	err = r.BtcRPCClient.ImportAddress(r.BTCDeployerAddress.EncodeAddress())
	require.NoError(r, err)
	r.Logger.Info("⚙️ imported BTCDeployerAddress: %s", r.BTCDeployerAddress.EncodeAddress())
}

// GetBtcAddress returns the BTC address of the deployer and private key in WIF format
func (r *E2ERunner) GetBtcAddress() (*btcutil.AddressWitnessPubKeyHash, *btcutil.WIF) {
	// load configured private key
	skBytes, err := hex.DecodeString(r.Account.RawPrivateKey.String())
	require.NoError(r, err)

	// create private key in WIF format
	sk, _ := btcec.PrivKeyFromBytes(skBytes)
	privkeyWIF, err := btcutil.NewWIF(sk, r.BitcoinParams, true)
	require.NoError(r, err)

	// derive address from private key
	address, err := btcutil.NewAddressWitnessPubKeyHash(
		btcutil.Hash160(privkeyWIF.SerializePubKey()),
		r.BitcoinParams,
	)
	require.NoError(r, err)

	return address, privkeyWIF
}

// SetupBtcAddress setups the deployer Bitcoin address
func (r *E2ERunner) SetupBtcAddress(name string, setupWallet bool) {
	// set the deployer address
	address, privkeyWIF := r.GetBtcAddress()
	r.BTCDeployerAddress = address

	r.Logger.Info("BTCDeployerAddress: %s, %v", r.BTCDeployerAddress.EncodeAddress(), setupWallet)

	// import the deployer private key as a Bitcoin node wallet
	if setupWallet {
		_, err := r.BtcRPCClient.CreateWallet(r.Name, rpcclient.WithCreateWalletBlank())
		if err != nil {
			require.ErrorContains(r, err, "Database already exists")
		}

		err = r.BtcRPCClient.ImportPrivKeyRescan(privkeyWIF, name, true)
		require.NoError(r, err, "failed to execute ImportPrivKeyRescan")
	}
}
