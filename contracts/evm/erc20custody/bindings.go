//go:generate sh -c "solc  ERC20Custody.sol --combined-json abi,bin | jq '.contracts.\"ERC20Custody.sol:ERC20Custody\"'  > ERC20Custody.json"
//go:generate sh -c "cat ERC20Custody.json | jq .abi > ERC20Custody.abi"
//go:generate sh -c "cat ERC20Custody.json | jq .bin | tr -d '\"' > ERC20Custody.bin"
//go:generate sh -c "abigen --abi ERC20Custody.abi --bin ERC20Custody.bin --pkg erc20custody --type ERC20Custody --out ERC20Custody.go"

package erc20custody

var _ = ERC20Custody{}
