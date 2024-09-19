//go:generate sh -c "solc IStaking.sol --combined-json abi | jq '.contracts.\"IStaking.sol:IStaking\"'  > IStaking.json"
//go:generate sh -c "cat IStaking.json | jq .abi > IStaking.abi"
//go:generate sh -c "abigen --abi IStaking.abi  --pkg staking --type IStaking --out IStaking.go"

package staking

var _ Contract
