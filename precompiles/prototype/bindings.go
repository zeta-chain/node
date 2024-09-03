//go:generate sh -c "solc IPrototype.sol --combined-json abi | jq '.contracts.\"IPrototype.sol:IPrototype\"'  > IPrototype.json"
//go:generate sh -c "cat IPrototype.json | jq .abi > IPrototype.abi"
//go:generate sh -c "abigen --abi IPrototype.abi  --pkg prototype --type IPrototype --out IPrototype.go"

package prototype

var _ Contract
