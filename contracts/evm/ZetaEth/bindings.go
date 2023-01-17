//go:generate sh -c "cat ZetaEth.json | jq .abi > ZetaEth.abi"
//go:generate sh -c "cat ZetaEth.json | jq .bytecode | tr -d '\"' > ZetaEth.bin"
//go:generate sh -c "abigen --abi ZetaEth.abi --bin ZetaEth.bin --pkg zetaeth --type ZetaEth --out ZetaEth.go"

package zetaeth

var _ = ZetaEth{}
