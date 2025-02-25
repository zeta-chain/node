// Dapp
//
//go:generate sh -c "solc --evm-version paris Dapp.sol --combined-json abi,bin | jq '.contracts.\"Dapp.sol:Dapp\"'  > Dapp.json"
//go:generate sh -c "cat Dapp.json | jq .abi > Dapp.abi"
//go:generate sh -c "cat Dapp.json | jq .bin  | tr -d '\"'  > Dapp.bin"
//go:generate sh -c "abigen --abi Dapp.abi --bin Dapp.bin  --pkg dapp --type Dapp --out Dapp.go"

package dapp
