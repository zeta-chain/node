//go:generate sh -c "solc TestGatewayZEVMCaller.sol --combined-json abi,bin | jq '.contracts.\"TestGatewayZEVMCaller.sol:TestGatewayZEVMCaller\"'  > TestGatewayZEVMCaller.json"
//go:generate sh -c "cat TestGatewayZEVMCaller.json | jq .abi > TestGatewayZEVMCaller.abi"
//go:generate sh -c "cat TestGatewayZEVMCaller.json | jq .bin  | tr -d '\"'  > TestGatewayZEVMCaller.bin"
//go:generate sh -c "abigen --abi TestGatewayZEVMCaller.abi --bin TestGatewayZEVMCaller.bin  --pkg testgatewayzevmcaller --type TestGatewayZEVMCaller --out TestGatewayZEVMCaller.go"

package testgatewayzevmcaller

var _ TestGatewayZEVMCaller
