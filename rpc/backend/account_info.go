package backend

import (
	"fmt"
	"math"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	rpctypes "github.com/zeta-chain/node/rpc/types"
)

// GetCode returns the contract code at the given address and block number.
func (b *Backend) GetCode(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (hexutil.Bytes, error) {
	blockNum, err := b.BlockNumberFromTendermint(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	req := &evmtypes.QueryCodeRequest{
		Address: address.String(),
	}

	res, err := b.QueryClient.Code(rpctypes.ContextWithHeight(blockNum.Int64()), req)
	if err != nil {
		return nil, err
	}

	return res.Code, nil
}

// GetProof returns an account object with proof and any storage proofs
func (b *Backend) GetProof(
	address common.Address,
	storageKeys []string,
	blockNrOrHash rpctypes.BlockNumberOrHash,
) (*rpctypes.AccountResult, error) {
	blockNum, err := b.BlockNumberFromTendermint(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	height := int64(blockNum)

	_, err = b.TendermintBlockByNumber(blockNum)
	if err != nil {
		// the error message imitates geth behavior
		return nil, errors.New("header not found")
	}

	// if the height is equal to zero, meaning the query condition of the block is either "pending" or "latest"
	if height == 0 {
		bn, err := b.BlockNumber()
		if err != nil {
			return nil, err
		}

		if bn > math.MaxInt64 {
			return nil, fmt.Errorf("not able to query block number greater than MaxInt64")
		}

		height = int64(bn) //#nosec G115 -- checked for int overflow already
	}

	ctx := rpctypes.ContextWithHeight(height)
	clientCtx := b.ClientCtx.WithHeight(height)

	// query storage proofs
	storageProofs := make([]rpctypes.StorageResult, len(storageKeys))

	for i, key := range storageKeys {
		hexKey := common.HexToHash(key)
		valueBz, proof, err := b.QueryClient.GetProof(
			clientCtx,
			evmtypes.StoreKey,
			evmtypes.StateKey(address, hexKey.Bytes()),
		)
		if err != nil {
			return nil, err
		}

		storageProofs[i] = rpctypes.StorageResult{
			Key:   key,
			Value: (*hexutil.Big)(new(big.Int).SetBytes(valueBz)),
			Proof: GetHexProofs(proof),
		}
	}

	// query EVM account
	req := &evmtypes.QueryAccountRequest{
		Address: address.String(),
	}

	res, err := b.QueryClient.Account(ctx, req)
	if err != nil {
		return nil, err
	}

	// query account proofs
	accountKey := bytes.HexBytes(append(authtypes.AddressStoreKeyPrefix, address.Bytes()...))
	_, proof, err := b.QueryClient.GetProof(clientCtx, authtypes.StoreKey, accountKey)
	if err != nil {
		return nil, err
	}

	balance, ok := sdkmath.NewIntFromString(res.Balance)
	if !ok {
		return nil, errors.New("invalid balance")
	}

	return &rpctypes.AccountResult{
		Address:      address,
		AccountProof: GetHexProofs(proof),
		Balance:      (*hexutil.Big)(balance.BigInt()),
		CodeHash:     common.HexToHash(res.CodeHash),
		Nonce:        hexutil.Uint64(res.Nonce),
		StorageHash:  common.Hash{}, // NOTE: Cosmos EVM doesn't have a storage hash. TODO: implement?
		StorageProof: storageProofs,
	}, nil
}

// GetStorageAt returns the contract storage at the given address, block number, and key.
func (b *Backend) GetStorageAt(
	address common.Address,
	key string,
	blockNrOrHash rpctypes.BlockNumberOrHash,
) (hexutil.Bytes, error) {
	blockNum, err := b.BlockNumberFromTendermint(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	req := &evmtypes.QueryStorageRequest{
		Address: address.String(),
		Key:     key,
	}

	res, err := b.QueryClient.Storage(rpctypes.ContextWithHeight(blockNum.Int64()), req)
	if err != nil {
		return nil, err
	}

	value := common.HexToHash(res.Value)
	return value.Bytes(), nil
}

// GetBalance returns the provided account's *spendable* balance up to the provided block number.
func (b *Backend) GetBalance(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (*hexutil.Big, error) {
	blockNum, err := b.BlockNumberFromTendermint(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	req := &evmtypes.QueryBalanceRequest{
		Address: address.String(),
	}

	_, err = b.TendermintBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	res, err := b.QueryClient.Balance(rpctypes.ContextWithHeight(blockNum.Int64()), req)
	if err != nil {
		return nil, err
	}

	val, ok := sdkmath.NewIntFromString(res.Balance)
	if !ok {
		return nil, errors.New("invalid balance")
	}

	// balance can only be negative in case of pruned node
	if val.IsNegative() {
		return nil, errors.New("couldn't fetch balance. Node state is pruned")
	}

	return (*hexutil.Big)(val.BigInt()), nil
}

// GetTransactionCount returns the number of transactions at the given address up to the given block number.
func (b *Backend) GetTransactionCount(address common.Address, blockNum rpctypes.BlockNumber) (*hexutil.Uint64, error) {
	n := hexutil.Uint64(0)
	bn, err := b.BlockNumber()
	if err != nil {
		return &n, err
	}
	height := blockNum.Int64()

	currentHeight := int64(bn) //#nosec G115 -- checked for int overflow already
	if height > currentHeight {
		return &n, errorsmod.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"cannot query with height in the future (current: %d, queried: %d); please provide a valid height",
			currentHeight, height,
		)
	}
	// Get nonce (sequence) from account
	from := sdk.AccAddress(address.Bytes())
	accRet := b.ClientCtx.AccountRetriever

	err = accRet.EnsureExists(b.ClientCtx, from)
	if err != nil {
		// account doesn't exist yet, return 0
		return &n, nil
	}

	includePending := blockNum == rpctypes.EthPendingBlockNumber
	nonce, err := b.getAccountNonce(address, includePending, blockNum.Int64(), b.Logger)
	if err != nil {
		return nil, err
	}

	n = hexutil.Uint64(nonce)
	return &n, nil
}
