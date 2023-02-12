package testdapp

//go:generate sh -c "solc TestDApp.sol --combined-json abi,bin | jq '.contracts.\"TestDApp.sol:TestDApp\"'  > TestDApp.json"
//go:generate sh -c "cat TestDApp.json | jq .abi | abigen --abi - --pkg testdapp --type TestDApp --out TestDApp.go"

var _ TestDApp
