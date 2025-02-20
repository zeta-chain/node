// Depositor
//go:generate sh -c "solc --evm-version paris Depositor.sol --combined-json abi,bin | jq '.contracts.\"Depositor.sol:Depositor\"'  > Depositor.json"
//go:generate sh -c "cat Depositor.json | jq .abi > Depositor.abi"
//go:generate sh -c "cat Depositor.json | jq .bin  | tr -d '\"'  > Depositor.bin"
//go:generate sh -c "abigen --abi Depositor.abi --bin Depositor.bin  --pkg depositor --type Depositor --out Depositor.go"

package depositor
