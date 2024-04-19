//go:generate sh -c "solc --evm-version paris TestZetaConnectorZEVM.sol --combined-json abi,bin | jq '.contracts.\"TestZetaConnectorZEVM.sol:TestZetaConnectorZEVM\"'  > TestZetaConnectorZEVM.json"
//go:generate sh -c "cat TestZetaConnectorZEVM.json | jq .abi > TestZetaConnectorZEVM.abi"
//go:generate sh -c "cat TestZetaConnectorZEVM.json | jq .bin  | tr -d '\"'  > TestZetaConnectorZEVM.bin"
//go:generate sh -c "abigen --abi TestZetaConnectorZEVM.abi --bin TestZetaConnectorZEVM.bin --pkg testconnectorzevm --type TestZetaConnectorZEVM --out TestZetaConnectorZEVM.go"

package testconnectorzevm
