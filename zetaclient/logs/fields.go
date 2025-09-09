package logs

// A group of predefined field keys and module names for zetaclient logs
const (
	// field keys
	FieldModule           = "module"
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
)

const (
	// module names
	ModNameOrchestrator   = "orchestrator"
	ModNameInbound        = "inbound"
	ModNameOutbound       = "outbound"
	ModNameSigner         = "signer"
	ModNameGasPrice       = "gas_price"
	ModNameClientBTC      = "btc_client"
	ModNameClientZetaCore = "zetacore_client"
	// TODO: This seems excessive.
	//       Replace for ModNameTss = "tss", then use FieldMethod for the suffixes?
	ModNameTssHealthCheck = "tss_healthcheck"
	ModNameTssKeyGen      = "tss_keygen"
	ModNameTssKeySign     = "tss_keysign"
	ModNameTssService     = "tss_service"
	ModNameTssSetup       = "tss_setup"
)
