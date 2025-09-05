package types

import (
	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/pkg/chains"
	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/pkg/crypto"
)

// NewEmptyRevertOptions initializes a new empty RevertOptions
func NewEmptyRevertOptions() RevertOptions {
	return RevertOptions{
		RevertGasLimit: sdkmath.ZeroUint(), // default to 0 instead of nil
	}
}

// NewRevertOptionsFromZEVM initializes a new RevertOptions from a gatewayzevm.RevertOptions
func NewRevertOptionsFromZEVM(revertOptions gatewayzevm.RevertOptions) RevertOptions {
	revertGasLimit := sdkmath.ZeroUint()
	if revertOptions.OnRevertGasLimit != nil {
		revertGasLimit = sdkmath.NewUintFromBigInt(revertOptions.OnRevertGasLimit)
	}

	return RevertOptions{
		RevertAddress:  revertOptions.RevertAddress.Hex(),
		CallOnRevert:   revertOptions.CallOnRevert,
		AbortAddress:   revertOptions.AbortAddress.Hex(),
		RevertMessage:  revertOptions.RevertMessage,
		RevertGasLimit: revertGasLimit,
	}
}

// NewRevertOptionsFromEVM initializes a new RevertOptions from a gatewayevm.RevertOptions
func NewRevertOptionsFromEVM(revertOptions gatewayevm.RevertOptions) RevertOptions {
	revertGasLimit := sdkmath.ZeroUint()
	if revertOptions.OnRevertGasLimit != nil {
		revertGasLimit = sdkmath.NewUintFromBigInt(revertOptions.OnRevertGasLimit)
	}

	return RevertOptions{
		RevertAddress:  revertOptions.RevertAddress.Hex(),
		CallOnRevert:   revertOptions.CallOnRevert,
		AbortAddress:   revertOptions.AbortAddress.Hex(),
		RevertMessage:  revertOptions.RevertMessage,
		RevertGasLimit: revertGasLimit,
	}
}

// NewRevertOptionsFromSOL initializes a new RevertOptions from a solana.RevertOptions
func NewRevertOptionsFromSOL(revertOptions solanacontracts.RevertOptions) RevertOptions {
	revertGasLimit := sdkmath.ZeroUint()
	if revertOptions.OnRevertGasLimit != 0 {
		revertGasLimit = sdkmath.Uint(sdkmath.NewIntFromUint64(revertOptions.OnRevertGasLimit))
	}

	return RevertOptions{
		RevertAddress:  revertOptions.RevertAddress.String(),
		AbortAddress:   revertOptions.AbortAddress.Hex(),
		CallOnRevert:   revertOptions.CallOnRevert,
		RevertMessage:  revertOptions.RevertMessage,
		RevertGasLimit: revertGasLimit,
	}
}

// ToZEVMRevertOptions converts the RevertOptions to a gatewayzevm.RevertOptions
func (r RevertOptions) ToZEVMRevertOptions() gatewayzevm.RevertOptions {
	return gatewayzevm.RevertOptions{
		RevertAddress: ethcommon.HexToAddress(r.RevertAddress),
		CallOnRevert:  r.CallOnRevert,
		AbortAddress:  ethcommon.HexToAddress(r.AbortAddress),
		RevertMessage: r.RevertMessage,
	}
}

// ToEVMRevertOptions converts the RevertOptions to a gatewayevm.RevertOptions
func (r RevertOptions) ToEVMRevertOptions() gatewayevm.RevertOptions {
	return gatewayevm.RevertOptions{
		RevertAddress: ethcommon.HexToAddress(r.RevertAddress),
		CallOnRevert:  r.CallOnRevert,
		AbortAddress:  ethcommon.HexToAddress(r.AbortAddress),
		RevertMessage: r.RevertMessage,
	}
}

// GetEVMRevertAddress returns the EVM revert address
// if the revert address is not a valid address, it returns false
func (r RevertOptions) GetEVMRevertAddress() (ethcommon.Address, bool) {
	addr := ethcommon.HexToAddress(r.RevertAddress)
	return addr, !crypto.IsEmptyAddress(addr)
}

// GetSOLRevertAddress returns the SOL revert address
// if the revert address is not a valid address, it returns false
func (r RevertOptions) GetSOLRevertAddress() (solana.PublicKey, bool) {
	addr, err := solana.PublicKeyFromBase58(r.RevertAddress)
	return addr, err == nil
}

// GetBTCRevertAddress validates and returns the BTC revert address
func (r RevertOptions) GetBTCRevertAddress(chainID int64) (string, bool) {
	btcAddress, err := chains.DecodeBtcAddress(r.RevertAddress, chainID)
	if err != nil {
		return "", false
	}
	if !chains.IsBtcAddressSupported(btcAddress) {
		return "", false
	}

	return btcAddress.EncodeAddress(), true
}

// GetEVMAbortAddress returns the EVM abort address
// if the abort address is not a valid address, it returns false
func (r RevertOptions) GetEVMAbortAddress() (ethcommon.Address, bool) {
	addr := ethcommon.HexToAddress(r.AbortAddress)
	return addr, !crypto.IsEmptyAddress(addr)
}
