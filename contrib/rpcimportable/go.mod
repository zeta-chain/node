module github.com/zeta-chain/node/contrib/rpcimportable

go 1.22.5

// this go.mod should be empty when committed
// the go.sum should not be committed

replace (
    github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)

// uncomment this for local development/testing/debugging
// replace github.com/zeta-chain/node => ../..