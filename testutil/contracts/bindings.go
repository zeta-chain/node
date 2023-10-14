// Example
//go:generate sh -c "solc --evm-version paris Example.sol --combined-json abi,bin | jq '.contracts.\"Example.sol:Example\"'  > Example.json"
//go:generate sh -c "cat Example.json | jq .abi > Example.abi"
//go:generate sh -c "cat Example.json | jq .bin  | tr -d '\"'  > Example.bin"
//go:generate sh -c "abigen --abi Example.abi --bin Example.bin  --pkg contracts --type Example --out Example.go"

// Reverter
//go:generate sh -c "solc --evm-version paris Reverter.sol --combined-json abi,bin | jq '.contracts.\"Reverter.sol:Reverter\"'  > Reverter.json"
//go:generate sh -c "cat Reverter.json | jq .abi > Reverter.abi"
//go:generate sh -c "cat Reverter.json | jq .bin  | tr -d '\"'  > Reverter.bin"
//go:generate sh -c "abigen --abi Reverter.abi --bin Reverter.bin  --pkg contracts --type Reverter --out Reverter.go"

package contracts
