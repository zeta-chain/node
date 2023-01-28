//go:generate sh -c "cat USDT.json | jq .abi > USDT.abi"
//go:generate sh -c "cat USDT.json | jq .bytecode | tr -d '\"' > USDT.bin"
//go:generate sh -c "abigen --abi USDT.abi --bin USDT.bin --pkg erc20 --type USDT --out USDT.go"

package erc20
