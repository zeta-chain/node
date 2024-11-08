//go:generate sh -c "solc TestGasConsumer.sol --evm-version london --combined-json abi,bin | jq '.contracts.\"TestGasConsumer.sol:TestGasConsumer\"'  > TestGasConsumer.json"
//go:generate sh -c "cat TestGasConsumer.json | jq .abi > TestGasConsumer.abi"
//go:generate sh -c "cat TestGasConsumer.json | jq .bin  | tr -d '\"'  > TestGasConsumer.bin"
//go:generate sh -c "abigen --abi TestGasConsumer.abi --bin TestGasConsumer.bin  --pkg testgasconsumer --type TestGasConsumer --out TestGasConsumer.go"

package testgasconsumer

var _ TestGasConsumer
