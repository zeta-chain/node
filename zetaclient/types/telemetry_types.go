package types

// Status type for telemetry. More fields can be added as needed
type Status struct {
	BTCNumberOfUTXOs int `json:"btc_number_of_utxos"`
}
