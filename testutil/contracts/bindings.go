//go:generate sh -c "solc --evm-version paris Example.sol --combined-json abi,bin | jq '.contracts.\"Example.sol:Example\"'  > Example.json"
//go:generate sh -c "cat Example.json | jq .abi > Example.abi"
//go:generate sh -c "cat Example.json | jq .bin  | tr -d '\"'  > Example.bin"
//go:generate sh -c "abigen --abi Example.abi --bin Example.bin  --pkg testcontracts --type Example --out Example.go"

package contracts
