package backend

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/cosmos/evm/mempool"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
)

// SendTransaction sends transaction based on received args using Node's key to sign it
func (b *Backend) SendTransaction(args evmtypes.TransactionArgs) (common.Hash, error) {
	// Look up the wallet containing the requested signer
	if !b.Cfg.JSONRPC.AllowInsecureUnlock {
		b.Logger.Debug("account unlock with HTTP access is forbidden")
		return common.Hash{}, fmt.Errorf("account unlock with HTTP access is forbidden")
	}

	_, err := b.ClientCtx.Keyring.KeyByAddress(sdk.AccAddress(args.GetFrom().Bytes()))
	if err != nil {
		b.Logger.Error("failed to find key in keyring", "address", args.GetFrom(), "error", err.Error())
		return common.Hash{}, fmt.Errorf("failed to find key in the node's keyring; %s; %s", keystore.ErrNoMatch, err.Error())
	}

	if args.ChainID != nil && (b.EvmChainID).Cmp((*big.Int)(args.ChainID)) != 0 {
		return common.Hash{}, fmt.Errorf("chainId does not match node's (have=%v, want=%v)", args.ChainID, (*hexutil.Big)(b.EvmChainID))
	}

	args, err = b.SetTxDefaults(args)
	if err != nil {
		return common.Hash{}, err
	}

	bn, err := b.BlockNumber()
	if err != nil {
		b.Logger.Debug("failed to fetch latest block number", "error", err.Error())
		return common.Hash{}, err
	}

	header, err := b.CurrentHeader()
	if err != nil {
		return common.Hash{}, err
	}

	signer := ethtypes.MakeSigner(b.ChainConfig(), new(big.Int).SetUint64(uint64(bn)), header.Time)

	// LegacyTx derives EvmChainID from the signature. To make sure the msg.ValidateBasic makes
	// the corresponding EvmChainID validation, we need to sign the transaction before calling it

	// Sign transaction
	msg := evmtypes.NewTxFromArgs(&args)
	if err := msg.Sign(signer, b.ClientCtx.Keyring); err != nil {
		b.Logger.Debug("failed to sign tx", "error", err.Error())
		return common.Hash{}, err
	}

	if err := msg.ValidateBasic(); err != nil {
		b.Logger.Debug("tx failed basic validation", "error", err.Error())
		return common.Hash{}, err
	}

	baseDenom := evmtypes.GetEVMCoinDenom()

	// Assemble transaction from fields
	tx, err := msg.BuildTx(b.ClientCtx.TxConfig.NewTxBuilder(), baseDenom)
	if err != nil {
		b.Logger.Error("build cosmos tx failed", "error", err.Error())
		return common.Hash{}, err
	}

	// Encode transaction by default Tx encoder
	txEncoder := b.ClientCtx.TxConfig.TxEncoder()
	txBytes, err := txEncoder(tx)
	if err != nil {
		b.Logger.Error("failed to encode eth tx using default encoder", "error", err.Error())
		return common.Hash{}, err
	}

	ethTx := msg.AsTransaction()

	// check the local node config in case unprotected txs are disabled
	if !b.UnprotectedAllowed() && !ethTx.Protected() {
		// Ensure only eip155 signed transactions are submitted if EIP155Required is set.
		return common.Hash{}, errors.New("only replay-protected (EIP-155) transactions allowed over RPC")
	}

	txHash := ethTx.Hash()

	// Broadcast transaction in sync mode (default)
	// NOTE: If error is encountered on the node, the broadcast will not return an error
	syncCtx := b.ClientCtx.WithBroadcastMode(flags.BroadcastSync)
	rsp, err := syncCtx.BroadcastTx(txBytes)
	if rsp != nil && rsp.Code != 0 {
		err = errorsmod.ABCIError(rsp.Codespace, rsp.Code, rsp.RawLog)
	}
	if err != nil {
		// Check if this is a nonce gap error that was successfully queued
		if b.Mempool != nil && strings.Contains(err.Error(), mempool.ErrNonceGap.Error()) {
			// Transaction was successfully queued due to nonce gap, return success to client
			b.Logger.Debug("transaction queued due to nonce gap", "hash", txHash.Hex())
			return txHash, nil
		}
		b.Logger.Error("failed to broadcast tx", "error", err.Error())
		return txHash, err
	}

	// Return transaction hash
	return txHash, nil
}

// Sign signs the provided data using the private key of address via Geth's signature standard.
func (b *Backend) Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	from := sdk.AccAddress(address.Bytes())

	_, err := b.ClientCtx.Keyring.KeyByAddress(from)
	if err != nil {
		b.Logger.Error("failed to find key in keyring", "address", address.String())
		return nil, fmt.Errorf("%s; %s", keystore.ErrNoMatch, err.Error())
	}

	// Sign the requested hash with the wallet
	signature, _, err := b.ClientCtx.Keyring.SignByAddress(from, data, signingtypes.SignMode_SIGN_MODE_TEXTUAL)
	if err != nil {
		b.Logger.Error("keyring.SignByAddress failed", "address", address.Hex())
		return nil, err
	}

	signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return signature, nil
}

// SignTypedData signs EIP-712 conformant typed data
func (b *Backend) SignTypedData(address common.Address, typedData apitypes.TypedData) (hexutil.Bytes, error) {
	from := sdk.AccAddress(address.Bytes())

	_, err := b.ClientCtx.Keyring.KeyByAddress(from)
	if err != nil {
		b.Logger.Error("failed to find key in keyring", "address", address.String())
		return nil, fmt.Errorf("%s; %s", keystore.ErrNoMatch, err.Error())
	}

	sigHash, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return nil, err
	}

	// Sign the requested hash with the wallet
	signature, _, err := b.ClientCtx.Keyring.SignByAddress(from, sigHash, signingtypes.SignMode_SIGN_MODE_TEXTUAL)
	if err != nil {
		b.Logger.Error("keyring.SignByAddress failed", "address", address.Hex())
		return nil, err
	}

	signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return signature, nil
}
