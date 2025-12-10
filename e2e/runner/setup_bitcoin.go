package runner

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"
)

func (r *E2ERunner) AddTSSToNode() {
	startTime := time.Now()
	defer func() {
		r.Logger.Print("✅ Bitcoin account setup in %s\n", time.Since(startTime))
	}()

	// import the TSS address
	err := r.BtcRPCClient.ImportAddress(r.Ctx, r.BTCTSSAddress.EncodeAddress())
	require.NoError(r, err)

	address, _ := r.GetBtcKeypair()

	// mine some blocks to get some BTC into the deployer address
	_, err = r.GenerateToAddressIfLocalBitcoin(101, address)
	require.NoError(r, err)
}

// SetupBitcoinAccounts sets up the TSS account and deployer account
func (r *E2ERunner) SetupBitcoinAccounts(createWallet bool) {
	r.Logger.Info("⚙️ setting up Bitcoin account")
	startTime := time.Now()
	defer func() {
		r.Logger.Print("✅ Bitcoin account setup in %s", time.Since(startTime))
	}()

	// setup deployer address
	r.SetupBtcAddress(createWallet)

	// import the TSS address to index TSS utxos and transactions
	err := r.BtcRPCClient.ImportAddress(r.Ctx, r.BTCTSSAddress.EncodeAddress())
	require.NoError(r, err)
	r.Logger.Info("⚙️ imported BTC TSSAddress: %s", r.BTCTSSAddress.EncodeAddress())
}

// GetBtcKeypair returns the BTC address of the runner account and private key in WIF format
func (r *E2ERunner) GetBtcKeypair() (*btcutil.AddressWitnessPubKeyHash, *btcutil.WIF) {
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

// GetBtcAddress returns the BTC address of the runner account
func (r *E2ERunner) GetBtcAddress() *btcutil.AddressWitnessPubKeyHash {
	address, _ := r.GetBtcKeypair()
	return address
}

// SetupBtcAddress setups the deployer Bitcoin address
func (r *E2ERunner) SetupBtcAddress(createWallet bool) {
	// set the deployer address
	address, _ := r.GetBtcKeypair()

	r.Logger.Info("BTCDeployerAddress: %s, %v", address.EncodeAddress(), createWallet)

	// import the deployer private key as a Bitcoin node wallet
	if createWallet {
		// we must use a raw request as the rpcclient does not expose the
		// descriptors arg which must be set to false
		// https://github.com/btcsuite/btcd/issues/2179
		// https://developer.bitcoin.org/reference/rpc/createwallet.html
		args := []interface{}{
			r.Name, // wallet_name
			true,   // disable_private_keys
			true,   // blank
			"",     // passphrase
			false,  // avoid_reuse
			false,  // descriptors
			true,   // load_on_startup
		}
		argsRawMsg := []json.RawMessage{}
		for _, arg := range args {
			encodedArg, err := json.Marshal(arg)
			require.NoError(r, err)
			argsRawMsg = append(argsRawMsg, encodedArg)
		}
		_, err := r.BtcRPCClient.RawRequest(r.Ctx, "createwallet", argsRawMsg)
		if err != nil {
			require.ErrorContains(r, err, "Database already exists")
		}
	}

	// import account address to index utxos and transactions
	err := r.BtcRPCClient.ImportAddress(r.Ctx, address.EncodeAddress())
	require.NoError(r, err)
	r.Logger.Info("⚙️ imported address: %s", address.EncodeAddress())
}
