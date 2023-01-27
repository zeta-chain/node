//go:generate sh -c "cat ZetaNonEth.json | jq .abi > ZetaNonEth.abi"
//go:generate sh -c "cat ZetaNonEth.json | jq .bytecode | tr -d '\"' > ZetaNonEth.bin"
//go:generate sh -c "abigen --abi ZetaNonEth.abi --bin ZetaNonEth.bin --pkg zetanoneth --type ZetaNonEth --out ZetaNonEth.go"

package zetanoneth

var _ = ZetaNonEth{}
