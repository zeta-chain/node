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

// Depositor
//go:generate sh -c "solc --evm-version paris Depositor.sol --combined-json abi,bin | jq '.contracts.\"Depositor.sol:Depositor\"'  > Depositor.json"
//go:generate sh -c "cat Depositor.json | jq .abi > Depositor.abi"
//go:generate sh -c "cat Depositor.json | jq .bin  | tr -d '\"'  > Depositor.bin"
//go:generate sh -c "abigen --abi Depositor.abi --bin Depositor.bin  --pkg contracts --type Depositor --out Depositor.go"

package contracts
