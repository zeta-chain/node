module github.com/zeta-chain/node/contrib/rpcimportable

go 1.22.5

require github.com/zeta-chain/node v0.0.0-20240903163921-74f1ab59c658 // indirect

replace (
    github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)

// uncomment this for local testing/development
// replace github.com/zeta-chain/node => ../..