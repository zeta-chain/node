package logs

// A group of predefined field keys and module names for zetaclient logs
const (
	// field keys
	FieldModule       = "module"
	FieldMethod       = "method"
	FieldChain        = "chain"
	FieldChainNetwork = "chain_network"
	FieldNonce        = "nonce"
	FieldTx           = "tx"
	FieldOutboundID   = "outbound_id"
	FieldCctx         = "cctx"
	FieldZetaTx       = "zeta_tx"
	FieldBallot       = "ballot"
	FieldCoinType     = "coin_type"

	// module names
	ModNameInbound  = "inbound"
	ModNameOutbound = "outbound"
	ModNameGasPrice = "gasprice"
	ModNameHeaders  = "headers"
)
