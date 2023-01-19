//go:generate sh -c "solc ERC20Custody.sol --combined-json abi,bin | jq '.contracts.\"ERC20Custody.sol:ERC20Custody\"'  > ERC20Custody.json"
//go:generate sh -c "cat ERC20Custody.json | jq .abi | abigen --abi - --pkg evm --type ERC20Custody --out ERC20Custody.go"

package evm

var _ = Connector{}
