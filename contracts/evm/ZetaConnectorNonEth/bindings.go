//go:generate sh -c "cat ZetaConnectorNonEth.json | jq .abi > ZetaConnectorNonEth.abi"
//go:generate sh -c "cat ZetaConnectorNonEth.json | jq .bytecode | tr -d '\"' > ZetaConnectorNonEth.bin"
//go:generate sh -c "abigen --abi ZetaConnectorNonEth.abi --bin ZetaConnectorNonEth.bin --pkg ZetaConnectorNonEth --type ZetaConnectorNonEth --out ZetaConnectorNonEth.go"

package ZetaConnectorNonEth

import (
	_ "embed"
)
