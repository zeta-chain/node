//go:generate sh -c "solc  TestDistribute.sol  --combined-json abi,bin | jq '.contracts.\"TestDistribute.sol:TestDistribute\"'  > TestDistribute.json"
//go:generate sh -c "cat TestDistribute.json | jq .abi > TestDistribute.abi"
//go:generate sh -c "cat TestDistribute.json | jq .bin  | tr -d '\"'  > TestDistribute.bin"
//go:generate sh -c "abigen --abi TestDistribute.abi --bin TestDistribute.bin  --pkg testdistribute --type TestDistribute --out TestDistribute.go"

package testdistribute

var _ TestDistribute
