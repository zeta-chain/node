package ante

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/evm/crypto/ethsecp256k1"
	legacysecp256k1 "github.com/cosmos/evm/legacy/ethsecp256k1"
)

const (
	secp256k1VerifyCost uint64 = 21_000
)

var _ authante.SignatureVerificationGasConsumer = DefaultSigVerificationGasConsumer

// DefaultSigVerificationGasConsumer is the default implementation of SignatureVerificationGasConsumer. It consumes gas
// for signature verification based upon the public key type. The cost is fetched from the given params and is matched
// by the concrete type.
func DefaultSigVerificationGasConsumer(
	meter storetypes.GasMeter,
	sig signing.SignatureV2,
	params authtypes.Params,
) error {
	pubkey := sig.PubKey
	switch pubkey := pubkey.(type) {
	case *ethsecp256k1.PubKey:
		meter.ConsumeGas(secp256k1VerifyCost, "ante verify: eth_secp256k1")
		return nil
	case *legacysecp256k1.PubKey:
		meter.ConsumeGas(secp256k1VerifyCost, "ante verify: legacy eth_secp256k1")
		return nil
	case multisig.PubKey:
		// Multisig keys
		multiSig, ok := sig.Data.(*signing.MultiSignatureData)
		if !ok {
			return fmt.Errorf("expected %T, got, %T", &signing.MultiSignatureData{}, sig.Data)
		}
		return ConsumeMultiSignatureVerificationGas(meter, multiSig, pubkey, params, sig.Sequence)
	default:
		return authante.DefaultSigVerificationGasConsumer(meter, sig, params)
	}
}

// ConsumeMultiSignatureVerificationGas consumes gas from a GasMeter for verifying a multisig pubkey signature
func ConsumeMultiSignatureVerificationGas(
	meter storetypes.GasMeter,
	sig *signing.MultiSignatureData,
	pubkey multisig.PubKey,
	params authtypes.Params,
	accSeq uint64,
) error {
	size := sig.BitArray.Count()
	sigIndex := 0

	for i := 0; i < size; i++ {
		if !sig.BitArray.GetIndex(i) {
			continue
		}

		sigV2 := signing.SignatureV2{
			PubKey:   pubkey.GetPubKeys()[i],
			Data:     sig.Signatures[sigIndex],
			Sequence: accSeq,
		}

		err := DefaultSigVerificationGasConsumer(meter, sigV2, params)
		if err != nil {
			return err
		}

		sigIndex++
	}

	return nil
}
