package logs

const (
	// FieldModule is the field key for the module name in logs
	FieldModule = "module"

	// FieldMethod is the field key for the method name in logs
	FieldMethod = "method"

	// FieldChain is the field key for the chain ID in logs
	FieldChain = "chain"

	// FieldNonce is the field key for the nonce in logs
	FieldNonce = "nonce"

	// FieldTx is the field key for the transaction hash in logs
	FieldTx = "tx"

	// FieldCctx is the field key for the cctx index in logs
	FieldCctx = "cctx"

	// ModNameInbound is the module name for inbound logs
	ModNameInbound = "inbound"

	// ModNameOutbound is the module name for outbound logs
	ModNameOutbound = "outbound"

	// ModNameGasPrice is the module name for gas price logs
	ModNameGasPrice = "gasprice"

	// ModNameHeaders is the module name for block headers logs
	ModNameHeaders = "headers"
)
