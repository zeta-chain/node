//go:generate sh -c "cat ZetaConnectorEth.json | jq .abi > ZetaConnectorEth.abi"
//go:generate sh -c "cat ZetaConnectorEth.json | jq .bytecode | tr -d '\"' > ZetaConnectorEth.bin"
//go:generate sh -c "abigen --abi ZetaConnectorEth.abi --bin ZetaConnectorEth.bin --pkg zetaconnectoreth --type ZetaConnectorEth --out ZetaConnectorEth.go"

package zetaconnectoreth
