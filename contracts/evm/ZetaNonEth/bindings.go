//go:generate sh -c "cat ZetaNonEth.json | jq .abi > ZetaNonEth.abi"
//go:generate sh -c "cat ZetaNonEth.json | jq .bytecode | tr -d '\"' > ZetaNonEth.bin"
//go:generate sh -c "abigen --abi ZetaNonEth.abi --bin ZetaNonEth.bin --pkg ZetaNonEth --type ZetaNonEth --out ZetaNonEth.go"

package ZetaNonEth

import (
	_ "embed"
)

var _ = ZetaNonEth{}
