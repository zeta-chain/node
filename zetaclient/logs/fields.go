package logs

// A group of predefined field keys and module names for zetaclient logs
const (
	FieldModule      = "module"
	FieldChain       = "chain"
	FieldNetwork     = "network"
	FieldBlock       = "block"
	FieldNonce       = "nonce"
	FieldTx          = "tx"
	FieldCctxIndex   = "cctx_index"
	FieldZetaTx      = "zeta_tx"
	FieldOutboundID  = "outbound_id"
	FieldBallotIndex = "ballot_index"
	FieldCoinType    = "coin_type"

	// chain specific
	FieldBtcTxid = "txid"
)

// module names
const (
	ModNameOrchestrator   = "orchestrator"
	ModNameInbound        = "inbound"
	ModNameOutbound       = "outbound"
	ModNameSigner         = "signer"
	ModNameGasPrice       = "gas_price"
	ModNameBtcClient      = "btc_client"
	ModNameZetaCoreClient = "zetacore_client"

	ModNameTssHealthCheck = "tss_healthcheck"
	ModNameTssKeyGen      = "tss_keygen"
	ModNameTssKeySign     = "tss_keysign"
	ModNameTssService     = "tss_service"
	ModNameTssSetup       = "tss_setup"
)
