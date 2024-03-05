//go:generate sh -c "solc TestDApp.sol --combined-json abi,bin | jq '.contracts.\"TestDApp.sol:TestDApp\"'  > TestDApp.json"
//go:generate sh -c "cat TestDApp.json | jq .abi > TestDApp.abi"
//go:generate sh -c "cat TestDApp.json | jq .bin  | tr -d '\"'  > TestDApp.bin"
//go:generate sh -c "abigen --abi TestDApp.abi --bin TestDApp.bin  --pkg testdapp --type TestDApp --out TestDApp.go"

package testdapp

var _ TestDApp
