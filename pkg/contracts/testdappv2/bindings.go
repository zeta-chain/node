//go:generate sh -c "solc TestDAppV2.sol --combined-json abi,bin | jq '.contracts.\"TestDAppV2.sol:TestDAppV2\"'  > TestDAppV2.json"
//go:generate sh -c "cat TestDAppV2.json | jq .abi > TestDAppV2.abi"
//go:generate sh -c "cat TestDAppV2.json | jq .bin  | tr -d '\"'  > TestDAppV2.bin"
//go:generate sh -c "abigen --abi TestDAppV2.abi --bin TestDAppV2.bin  --pkg testdappv2 --type TestDAppV2 --out TestDAppV2.go"

package testdappv2

var _ TestDAppV2
