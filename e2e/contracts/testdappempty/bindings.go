//go:generate sh -c "solc TestDAppEmpty.sol --evm-version london --combined-json abi,bin | jq '.contracts.\"TestDAppEmpty.sol:TestDAppEmpty\"'  > TestDAppEmpty.json"
//go:generate sh -c "cat TestDAppEmpty.json | jq .abi > TestDAppEmpty.abi"
//go:generate sh -c "cat TestDAppEmpty.json | jq .bin  | tr -d '\"'  > TestDAppEmpty.bin"
//go:generate sh -c "abigen --abi TestDAppEmpty.abi --bin TestDAppEmpty.bin  --pkg testdappempty --type TestDAppEmpty --out TestDAppEmpty.go"

package testdappempty

var _ TestDAppEmpty
