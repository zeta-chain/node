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
replace (
    github.com/ethereum/go-ethereum => github.com/zeta-chain/go-ethereum v1.13.16-0.20241022183758-422c6ef93ccc
)

// uncomment this for local development/testing/debugging
// replace github.com/zeta-chain/node => ../..