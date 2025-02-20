// DappReverter
//go:generate sh -c "solc --evm-version paris DappReverter.sol --combined-json abi,bin | jq '.contracts.\"DappReverter.sol:DappReverter\"'  > DappReverter.json"
//go:generate sh -c "cat DappReverter.json | jq .abi > DappReverter.abi"
//go:generate sh -c "cat DappReverter.json | jq .bin  | tr -d '\"'  > DappReverter.bin"
//go:generate sh -c "abigen --abi DappReverter.abi --bin DappReverter.bin  --pkg dappreverter --type DappReverter --out DappReverter.go"

package dappreverter
