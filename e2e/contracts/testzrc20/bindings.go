//go:generate sh -c "solc --evm-version paris TestZRC20.sol --combined-json abi,bin | jq '.contracts.\"TestZRC20.sol:TestZRC20\"'  > TestZRC20.json"
//go:generate sh -c "cat TestZRC20.json | jq .abi > TestZRC20.abi"
//go:generate sh -c "cat TestZRC20.json | jq .bin  | tr -d '\"'  > TestZRC20.bin"
//go:generate sh -c "abigen --abi TestZRC20.abi --bin TestZRC20.bin --pkg testzrc20 --type TestZRC20 --out TestZRC20.go"

package testzrc20
