//go:generate sh -c "solc --evm-version cancun IStaking.sol --combined-json abi,bin --allow-paths .. | jq '.contracts.\"IStaking.sol:IStaking\"'  > IStaking.json"
//go:generate sh -c "cat IStaking.json | jq .abi > IStaking.abi"
//go:generate sh -c "cat IStaking.json | jq .bin  | tr -d '\"'  > IStaking.bin"
//go:generate sh -c "abigen --abi IStaking.abi --bin IStaking.bin --pkg istaking --type IStaking --out IStaking.go"

// Package istaking contains the bindings for the staking precompile interface
// It is used for E2E testing with the staking precompiles
// Note: for simplicity and to only validate that the precompile is enable, only the staking function is defined in the interface
package istaking
