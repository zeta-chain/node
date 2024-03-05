//go:generate sh -c "solc --evm-version paris ZEVMSwapApp.sol --combined-json abi,bin --allow-paths .. | jq '.contracts.\"ZEVMSwapApp.sol:ZEVMSwapApp\"'  > ZEVMSwapApp.json"
//go:generate sh -c "cat ZEVMSwapApp.json | jq .abi > ZEVMSwapApp.abi"
//go:generate sh -c "cat ZEVMSwapApp.json | jq .bin  | tr -d '\"'  > ZEVMSwapApp.bin"
//go:generate sh -c "abigen --abi ZEVMSwapApp.abi --bin ZEVMSwapApp.bin --pkg zevmswap --type ZEVMSwapApp --out ZEVMSwapApp.go"

package zevmswap
