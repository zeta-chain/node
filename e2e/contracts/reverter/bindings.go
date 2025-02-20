// Reverter
//go:generate sh -c "solc --evm-version paris Reverter.sol --combined-json abi,bin | jq '.contracts.\"Reverter.sol:Reverter\"'  > Reverter.json"
//go:generate sh -c "cat Reverter.json | jq .abi > Reverter.abi"
//go:generate sh -c "cat Reverter.json | jq .bin  | tr -d '\"'  > Reverter.bin"
//go:generate sh -c "abigen --abi Reverter.abi --bin Reverter.bin  --pkg reverter --type Reverter --out Reverter.go"

package reverter
