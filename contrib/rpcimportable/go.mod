module github.com/zeta-chain/node/contrib/rpcimportable

go 1.22.5

// this go.mod should be empty when committed
// the go.sum should not be committed

// this replacement is unavoidable until we upgrade cosmos sdk >=v0.50
// but we should not tolerate any other replacements
replace (
    github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)

// go-ethereum fork must be used as it removes incompatible pebbledb version
// TODO evm: verify if pebbledb issue is in cosmos go ethereum fork
replace (
    github.com/ethereum/go-ethereum => github.com/cosmos/go-ethereum v1.15.11-cosmos-0
)

replace (
	github.com/cosmos/evm => github.com/zeta-chain/evm v0.0.0-20250808111716-1882abec3ec9
)

// protocol-contracts was renamed to protocol-contracts-evm; redirect to old module path
replace (
	github.com/zeta-chain/protocol-contracts-evm => github.com/zeta-chain/protocol-contracts v0.0.0-20250909184950-6034c08e5870
)

// uncomment this for local development/testing/debugging
// replace github.com/zeta-chain/node => ../..