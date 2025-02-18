//go:generate sh -c "solc --evm-version paris TestAbort.sol --combined-json abi,bin --allow-paths .. | jq '.contracts.\"TestAbort.sol:TestAbort\"'  > TestAbort.json"
//go:generate sh -c "cat TestAbort.json | jq .abi > TestAbort.abi"
//go:generate sh -c "cat TestAbort.json | jq .bin  | tr -d '\"'  > TestAbort.bin"
//go:generate sh -c "abigen --abi TestAbort.abi --bin TestAbort.bin --pkg testabort --type TestAbort --out TestAbort.go"

package testabort
