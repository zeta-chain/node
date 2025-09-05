package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	mempool "github.com/cosmos/cosmos-sdk/types/mempool"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

var _ mempool.SignerExtractionAdapter = EthSignerExtractionAdapter{}

// EthSignerExtractionAdapter is the default implementation of SignerExtractionAdapter. It extracts the signers
// from a cosmos-sdk tx via GetSignaturesV2.
type EthSignerExtractionAdapter struct {
	fallback mempool.SignerExtractionAdapter
}

// NewEthSignerExtractionAdapter constructs a new EthSignerExtractionAdapter instance
func NewEthSignerExtractionAdapter(fallback mempool.SignerExtractionAdapter) EthSignerExtractionAdapter {
	return EthSignerExtractionAdapter{fallback}
}

// GetSigners implements the Adapter interface
// NOTE: only the first item is used by the mempool
func (s EthSignerExtractionAdapter) GetSigners(tx sdk.Tx) ([]mempool.SignerData, error) {
	if txWithExtensions, ok := tx.(authante.HasExtensionOptionsTx); ok {
		opts := txWithExtensions.GetExtensionOptions()
		if len(opts) > 0 && opts[0].GetTypeUrl() == "/cosmos.evm.vm.v1.ExtensionOptionsEthereumTx" {
			for _, msg := range tx.GetMsgs() {
				if ethMsg, ok := msg.(*evmtypes.MsgEthereumTx); ok {
					return []mempool.SignerData{
						mempool.NewSignerData(
							ethMsg.GetFrom(),
							ethMsg.AsTransaction().Nonce(),
						),
					}, nil
				}
			}
		}
	}

	return s.fallback.GetSigners(tx)
}
