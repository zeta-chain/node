// Withdrawer
//go:generate sh -c "solc --evm-version paris Withdrawer.sol --combined-json abi,bin | jq '.contracts.\"Withdrawer.sol:Withdrawer\"'  > Withdrawer.json"
//go:generate sh -c "cat Withdrawer.json | jq .abi > Withdrawer.abi"
//go:generate sh -c "cat Withdrawer.json | jq .bin  | tr -d '\"'  > Withdrawer.bin"
//go:generate sh -c "abigen --abi Withdrawer.abi --bin Withdrawer.bin  --pkg withdrawer --type Withdrawer --out Withdrawer.go"

package withdrawer
