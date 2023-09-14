//go:generate sh -c "solc --evm-version paris Vault.sol --combined-json abi,bin | jq '.contracts.\"Vault.sol:Vault\"'  > Vault.json"
//go:generate sh -c "cat Vault.json | jq .abi > Vault.abi"
//go:generate sh -c "cat Vault.json | jq .bin  | tr -d '\"'  > Vault.bin"
//go:generate sh -c "abigen --abi Vault.abi --bin Vault.bin  --pkg vault --type Vault --out Vault.go"

package vault
