//go:generate sh -c "cat ZetaEth.json | jq .abi > ZetaEth.abi"
//go:generate sh -c "cat ZetaEth.json | jq .bytecode | tr -d '\"' > ZetaEth.bin"
//go:generate sh -c "abigen --abi ZetaEth.abi --bin ZetaEth.bin --pkg ZetaEth --type ZetaEth --out ZetaEth.go"

package ZetaEth

import (
	_ "embed"
)

var _ = ZetaEth{}
