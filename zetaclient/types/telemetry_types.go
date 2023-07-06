package types

// Status type for telemetry. More fields can be added as needed
type Status struct {
	BTCNextNonce             int `json:"btc_next_nonce"`
	BTCNumberOfUTXOs         int `json:"btc_number_of_utxos"`
	BTCNumberOfFilteredUTXOs int `json:"btc_number_of_filtered_utxos"`
}
