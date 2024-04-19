package types

// DefaultVerificationFlags returns the default verification flags.
// By default, everything disabled.
func DefaultVerificationFlags() VerificationFlags {
	return VerificationFlags{
		EthTypeChainEnabled: false,
		BtcTypeChainEnabled: false,
	}
}
