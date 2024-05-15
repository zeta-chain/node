//go:generate sh -c "solc --evm-version paris TestDAppNoRevert.sol --combined-json abi,bin | jq '.contracts.\"TestDAppNoRevert.sol:TestDAppNoRevert\"'  > TestDAppNoRevert.json"
//go:generate sh -c "cat TestDAppNoRevert.json | jq .abi > TestDAppNoRevert.abi"
//go:generate sh -c "cat TestDAppNoRevert.json | jq .bin  | tr -d '\"'  > TestDAppNoRevert.bin"
//go:generate sh -c "abigen --abi TestDAppNoRevert.abi --bin TestDAppNoRevert.bin  --pkg testdappnorevert --type TestDAppNoRevert --out TestDAppNoRevert.go"

package testdappnorevert

var _ TestDAppNoRevert
