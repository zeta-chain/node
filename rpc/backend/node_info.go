package backend

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	cmttypes "github.com/cometbft/cometbft/types"

	"github.com/cosmos/evm/crypto/ethsecp256k1"
	"github.com/cosmos/evm/testutil/constants"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	rpctypes "github.com/zeta-chain/node/rpc/types"
	"github.com/zeta-chain/node/server/config"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdkconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// Accounts returns the list of accounts available to this node.
func (b *Backend) Accounts() ([]common.Address, error) {
	addresses := make([]common.Address, 0) // return [] instead of nil if empty

	if !b.Cfg.JSONRPC.AllowInsecureUnlock {
		b.Logger.Debug("account unlock with HTTP access is forbidden")
		return addresses, fmt.Errorf("account unlock with HTTP access is forbidden")
	}

	infos, err := b.ClientCtx.Keyring.List()
	if err != nil {
		return addresses, err
	}

	for _, info := range infos {
		pubKey, err := info.GetPubKey()
		if err != nil {
			return nil, err
		}
		addressBytes := pubKey.Address().Bytes()
		addresses = append(addresses, common.BytesToAddress(addressBytes))
	}

	return addresses, nil
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronize from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (b *Backend) Syncing() (interface{}, error) {
	status, err := b.ClientCtx.Client.Status(b.Ctx)
	if err != nil {
		return false, err
	}

	if !status.SyncInfo.CatchingUp {
		return false, nil
	}

	return map[string]interface{}{
		"startingBlock": hexutil.Uint64(status.SyncInfo.EarliestBlockHeight), //nolint:gosec // G115 // won't exceed uint64
		"currentBlock":  hexutil.Uint64(status.SyncInfo.LatestBlockHeight),   //nolint:gosec // G115 // won't exceed uint64
		// "highestBlock":  nil, // NA
		// "pulledStates":  nil, // NA
		// "knownStates":   nil, // NA
	}, nil
}

// SetEtherbase sets the etherbase of the miner
func (b *Backend) SetEtherbase(etherbase common.Address) bool {
	if !b.Cfg.JSONRPC.AllowInsecureUnlock {
		b.Logger.Debug("account unlock with HTTP access is forbidden")
		return false
	}

	delAddr, err := b.GetCoinbase()
	if err != nil {
		b.Logger.Debug("failed to get coinbase address", "error", err.Error())
		return false
	}

	withdrawAddr := sdk.AccAddress(etherbase.Bytes())
	msg := distributiontypes.NewMsgSetWithdrawAddress(delAddr, withdrawAddr)

	// Assemble transaction from fields
	builder, ok := b.ClientCtx.TxConfig.NewTxBuilder().(authtx.ExtensionOptionsTxBuilder)
	if !ok {
		b.Logger.Debug("clientCtx.TxConfig.NewTxBuilder returns unsupported builder")
		return false
	}

	err = builder.SetMsgs(msg)
	if err != nil {
		b.Logger.Error("builder.SetMsgs failed", "error", err.Error())
		return false
	}

	// Fetch minimum gas price to calculate fees using the configuration.
	minGasPrices := b.Cfg.GetMinGasPrices()
	if len(minGasPrices) == 0 || minGasPrices.Empty() {
		b.Logger.Debug("the minimum fee is not set")
		return false
	}
	minGasPriceValue := minGasPrices[0].Amount
	denom := minGasPrices[0].Denom

	delCommonAddr := common.BytesToAddress(delAddr.Bytes())
	nonce, err := b.GetTransactionCount(delCommonAddr, rpctypes.EthPendingBlockNumber)
	if err != nil {
		b.Logger.Debug("failed to get nonce", "error", err.Error())
		return false
	}

	txFactory := tx.Factory{}
	txFactory = txFactory.
		WithChainID(b.ClientCtx.ChainID).
		WithKeybase(b.ClientCtx.Keyring).
		WithTxConfig(b.ClientCtx.TxConfig).
		WithSequence(uint64(*nonce)).
		WithGasAdjustment(1.25)

	_, gas, err := tx.CalculateGas(b.ClientCtx, txFactory, msg)
	if err != nil {
		b.Logger.Debug("failed to calculate gas", "error", err.Error())
		return false
	}

	txFactory = txFactory.WithGas(gas)

	value := new(big.Int).SetUint64(gas * minGasPriceValue.Ceil().TruncateInt().Uint64())
	fees := sdk.Coins{sdk.NewCoin(denom, sdkmath.NewIntFromBigInt(value))}
	builder.SetFeeAmount(fees)
	builder.SetGasLimit(gas)

	keyInfo, err := b.ClientCtx.Keyring.KeyByAddress(delAddr)
	if err != nil {
		b.Logger.Debug("failed to get the wallet address using the keyring", "error", err.Error())
		return false
	}

	if err := tx.Sign(b.ClientCtx.CmdContext, txFactory, keyInfo.Name, builder, false); err != nil {
		b.Logger.Debug("failed to sign tx", "error", err.Error())
		return false
	}

	// Encode transaction by default Tx encoder
	txEncoder := b.ClientCtx.TxConfig.TxEncoder()
	txBytes, err := txEncoder(builder.GetTx())
	if err != nil {
		b.Logger.Debug("failed to encode eth tx using default encoder", "error", err.Error())
		return false
	}

	tmHash := common.BytesToHash(cmttypes.Tx(txBytes).Hash())

	// Broadcast transaction in sync mode (default)
	// NOTE: If error is encountered on the node, the broadcast will not return an error
	syncCtx := b.ClientCtx.WithBroadcastMode(flags.BroadcastSync)
	rsp, err := syncCtx.BroadcastTx(txBytes)
	if rsp != nil && rsp.Code != 0 {
		err = errorsmod.ABCIError(rsp.Codespace, rsp.Code, rsp.RawLog)
	}
	if err != nil {
		b.Logger.Debug("failed to broadcast tx", "error", err.Error())
		return false
	}

	b.Logger.Debug("broadcasted tx to set miner withdraw address (etherbase)", "hash", tmHash.String())
	return true
}

// ImportRawKey armors and encrypts a given raw hex encoded ECDSA key and stores it into the key directory.
// The name of the key will have the format "personal_<length-keys>", where <length-keys> is the total number of
// keys stored on the keyring.
//
// NOTE: The key will be both armored and encrypted using the same passphrase.
func (b *Backend) ImportRawKey(privkey, password string) (common.Address, error) {
	priv, err := crypto.HexToECDSA(privkey)
	if err != nil {
		return common.Address{}, err
	}

	privKey := &ethsecp256k1.PrivKey{Key: crypto.FromECDSA(priv)}

	addr := sdk.AccAddress(privKey.PubKey().Address().Bytes())
	ethereumAddr := common.BytesToAddress(addr)

	// return if the key has already been imported
	if _, err := b.ClientCtx.Keyring.KeyByAddress(addr); err == nil {
		return ethereumAddr, nil
	}

	// ignore error as we only care about the length of the list
	list, _ := b.ClientCtx.Keyring.List() // #nosec G703
	privKeyName := fmt.Sprintf("personal_%d", len(list))

	armor := sdkcrypto.EncryptArmorPrivKey(privKey, password, ethsecp256k1.KeyType)

	if err := b.ClientCtx.Keyring.ImportPrivKey(privKeyName, armor, password); err != nil {
		return common.Address{}, err
	}

	b.Logger.Info("key successfully imported", "name", privKeyName, "address", ethereumAddr.String())

	return ethereumAddr, nil
}

// ListAccounts will return a list of addresses for accounts this node manages.
func (b *Backend) ListAccounts() ([]common.Address, error) {
	addrs := []common.Address{}

	if !b.Cfg.JSONRPC.AllowInsecureUnlock {
		b.Logger.Debug("account unlock with HTTP access is forbidden")
		return addrs, fmt.Errorf("account unlock with HTTP access is forbidden")
	}

	list, err := b.ClientCtx.Keyring.List()
	if err != nil {
		return nil, err
	}

	for _, info := range list {
		pubKey, err := info.GetPubKey()
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, common.BytesToAddress(pubKey.Address()))
	}

	return addrs, nil
}

// NewAccount will create a new account and returns the address for the new account.
func (b *Backend) NewMnemonic(uid string,
	_ keyring.Language,
	hdPath,
	bip39Passphrase string,
	algo keyring.SignatureAlgo,
) (*keyring.Record, error) {
	info, _, err := b.ClientCtx.Keyring.NewMnemonic(uid, keyring.English, hdPath, bip39Passphrase, algo)
	if err != nil {
		return nil, err
	}
	return info, err
}

// SetGasPrice sets the minimum accepted gas price for the miner.
// NOTE: this function accepts only integers to have the same interface than go-eth
// to use float values, the gas prices must be configured using the configuration file
func (b *Backend) SetGasPrice(gasPrice hexutil.Big) bool {
	appConf, err := config.GetConfig(b.ClientCtx.Viper)
	if err != nil {
		b.Logger.Debug("could not get the server config", "error", err.Error())
		return false
	}
	c := b.GenerateMinGasCoin(gasPrice, appConf)

	appConf.SetMinGasPrices(sdk.DecCoins{c})
	sdkconfig.WriteConfigFile(b.ClientCtx.Viper.ConfigFileUsed(), appConf)
	b.Logger.Info("Your configuration file was modified. Please RESTART your node.", "gas-price", c.String())
	return true
}

func (b *Backend) GenerateMinGasCoin(gasPrice hexutil.Big, appConf config.Config) sdk.DecCoin {
	var unit string
	minGasPrices := appConf.GetMinGasPrices()

	// fetch the base denom from the sdk Config in case it's not currently defined on the node config
	if len(minGasPrices) == 0 || minGasPrices.Empty() {
		unit = evmtypes.GetEVMCoinDenom()
	} else {
		unit = minGasPrices[0].Denom
	}

	// The provided gasPrice has 18 decimals.
	// We need to update to the denom's real precision
	scaledAmt := evmtypes.ConvertBigIntFrom18DecimalsToLegacyDec(gasPrice.ToInt())
	c := sdk.DecCoin{Denom: unit, Amount: scaledAmt}

	return c
}

// UnprotectedAllowed returns the node configuration value for allowing
// unprotected transactions (i.e not replay-protected)
func (b Backend) UnprotectedAllowed() bool {
	return b.AllowUnprotectedTxs
}

// RPCGasCap is the global gas cap for eth-call variants.
func (b *Backend) RPCGasCap() uint64 {
	return b.Cfg.JSONRPC.GasCap
}

// RPCEVMTimeout is the global evm timeout for eth-call variants.
func (b *Backend) RPCEVMTimeout() time.Duration {
	return b.Cfg.JSONRPC.EVMTimeout
}

// RPCGasCap is the global gas cap for eth-call variants.
func (b *Backend) RPCTxFeeCap() float64 {
	return b.Cfg.JSONRPC.TxFeeCap
}

// RPCFilterCap is the limit for total number of filters that can be created
func (b *Backend) RPCFilterCap() int32 {
	return b.Cfg.JSONRPC.FilterCap
}

// RPCFeeHistoryCap is the limit for total number of blocks that can be fetched
func (b *Backend) RPCFeeHistoryCap() int32 {
	return b.Cfg.JSONRPC.FeeHistoryCap
}

// RPCLogsCap defines the max number of results can be returned from single `eth_getLogs` query.
func (b *Backend) RPCLogsCap() int32 {
	return b.Cfg.JSONRPC.LogsCap
}

// RPCBlockRangeCap defines the max block range allowed for `eth_getLogs` query.
func (b *Backend) RPCBlockRangeCap() int32 {
	return b.Cfg.JSONRPC.BlockRangeCap
}

// RPCMinGasPrice returns the minimum gas price for a transaction obtained from
// the node config. If set value is 0, it will default to 20.
func (b *Backend) RPCMinGasPrice() *big.Int {
	baseDenom := evmtypes.GetEVMCoinDenom()

	minGasPrice := b.Cfg.GetMinGasPrices()
	amt := minGasPrice.AmountOf(baseDenom)
	if amt.IsNil() || amt.IsZero() {
		return big.NewInt(constants.DefaultGasPrice)
	}

	return evmtypes.ConvertAmountTo18DecimalsLegacy(amt).TruncateInt().BigInt()
}
