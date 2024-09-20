//go:generate sh -c "solc IBank.sol --combined-json abi | jq '.contracts.\"IBank.sol:IBank\"'  > IBank.json"
//go:generate sh -c "cat IBank.json | jq .abi > IBank.abi"
//go:generate sh -c "abigen --abi IBank.abi  --pkg bank --type IBank --out IBank.gen.go"

package bank

var _ Contract
