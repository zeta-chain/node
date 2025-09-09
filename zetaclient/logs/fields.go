package logs

// A group of predefined field keys and module names for zetaclient logs
const (
	// field keys
	FieldModule           = "module"
	FieldMethod           = "method"
	FieldChain            = "chain"
	FieldChainNetwork     = "chain_network"
	FieldNonce            = "nonce"
	FieldTracker          = "tracker_id"
	FieldTx               = "tx"
	FieldOutboundID       = "outbound_id"
	FieldBlock            = "block"
	FieldCctx             = "cctx"
	FieldZetaTx           = "zeta_tx"
	FieldBallot           = "ballot"
	FieldCoinType         = "coin_type"
	FieldConfirmationMode = "confirmation_mode"

	// chain specific field keys
	FieldBtcTxid = "txid"

	// module names
	ModNameInbound  = "inbound"
	ModNameOutbound = "outbound"
	ModNameGasPrice = "gasprice"
)
