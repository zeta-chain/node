package logs

// A group of predefined field keys and module names for zetaclient logs
const (
	// field keys
	FieldModule = "module"
	FieldMethod = "method"
	FieldChain  = "chain"
	FieldNonce  = "nonce"
	FieldTx     = "tx"
	FieldCctx   = "cctx"

	// module names
	ModNameInbound  = "inbound"
	ModNameOutbound = "outbound"
	ModNameGasPrice = "gasprice"
	ModNameHeaders  = "headers"
)
