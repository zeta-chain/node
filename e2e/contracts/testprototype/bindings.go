//go:generate sh -c "solc  TestPrototype.sol  --combined-json abi,bin | jq '.contracts.\"TestPrototype.sol:TestPrototype\"'  > TestPrototype.json"
//go:generate sh -c "cat TestPrototype.json | jq .abi > TestPrototype.abi"
//go:generate sh -c "cat TestPrototype.json | jq .bin  | tr -d '\"'  > TestPrototype.bin"
//go:generate sh -c "abigen --abi TestPrototype.abi --bin TestPrototype.bin  --pkg testprototype --type TestPrototype --out TestPrototype.go"

package testprototype

var _ TestPrototype
