//go:generate sh -c "solc  TestBank.sol  --combined-json abi,bin | jq '.contracts.\"TestBank.sol:TestBank\"'  > TestBank.json"
//go:generate sh -c "cat TestBank.json | jq .abi > TestBank.abi"
//go:generate sh -c "cat TestBank.json | jq .bin  | tr -d '\"'  > TestBank.bin"
//go:generate sh -c "abigen --abi TestBank.abi --bin TestBank.bin  --pkg testbank --type TestBank --out TestBank.go"

package testbank

var _ TestBank
