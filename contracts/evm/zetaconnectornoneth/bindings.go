//go:generate sh -c "cat ZetaConnectorNonEth.json | jq .abi > ZetaConnectorNonEth.abi"
//go:generate sh -c "cat ZetaConnectorNonEth.json | jq .bytecode | tr -d '\"' > ZetaConnectorNonEth.bin"
//go:generate sh -c "abigen --abi ZetaConnectorNonEth.abi --bin ZetaConnectorNonEth.bin --pkg zetaconnectornoneth --type ZetaConnectorNonEth --out ZetaConnectorNonEth.go"

package zetaconnectornoneth
