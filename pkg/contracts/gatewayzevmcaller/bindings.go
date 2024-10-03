//go:generate sh -c "solc GatewayZEVMCaller.sol --combined-json abi,bin | jq '.contracts.\"GatewayZEVMCaller.sol:GatewayZEVMCaller\"'  > GatewayZEVMCaller.json"
//go:generate sh -c "cat GatewayZEVMCaller.json | jq .abi > GatewayZEVMCaller.abi"
//go:generate sh -c "cat GatewayZEVMCaller.json | jq .bin  | tr -d '\"'  > GatewayZEVMCaller.bin"
//go:generate sh -c "abigen --abi GatewayZEVMCaller.abi --bin GatewayZEVMCaller.bin  --pkg gatewayzevmcaller --type GatewayZEVMCaller --out GatewayZEVMCaller.go"

package gatewayzevmcaller

var _ GatewayZEVMCaller
