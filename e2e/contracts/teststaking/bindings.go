//go:generate sh -c "solc  TestStaking.sol  --combined-json abi,bin | jq '.contracts.\"TestStaking.sol:TestStaking\"'  > TestStaking.json"
//go:generate sh -c "cat TestStaking.json | jq .abi > TestStaking.abi"
//go:generate sh -c "cat TestStaking.json | jq .bin  | tr -d '\"'  > TestStaking.bin"
//go:generate sh -c "abigen --abi TestStaking.abi --bin TestStaking.bin  --pkg teststaking --type TestStaking --out TestStaking.go"

package teststaking

var _ TestStaking
