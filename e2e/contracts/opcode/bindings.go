// Example
//go:generate sh -c "solc --version && solc --evm-version cancun Opcode.sol --combined-json abi,bin | jq '.contracts.\"Opcode.sol:Opcode\"' > Opcode.json"
//go:generate sh -c "cat Opcode.json | jq .abi > Opcode.abi"
//go:generate sh -c "cat Opcode.json | jq .bin  | tr -d '\"'  > Opcode.bin"
//go:generate sh -c "abigen --abi Opcode.abi --bin Opcode.bin  --pkg opcode --type Opcode --out Opcode.go"

package opcode
