package logs

// A group of predefined field keys and module names for zetaclient logs
const (
	// field keys
	FieldModule       = "module"
	FieldChain        = "chain"
	FieldChainNetwork = "chain_network"
	FieldBlock        = "block"
	FieldNonce        = "nonce"
	FieldTx           = "tx"
	FieldCctx         = "cctx"
	FieldZetaTx       = "zeta_tx"
	FieldOutboundID   = "outbound_id"
	FieldBallot       = "ballot"
	FieldCoinType     = "coin_type"

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
