module github.com/Meta-Protocol/metacore

go 1.16

require (
	github.com/binance-chain/tss-lib v1.3.2 // indirect
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/cosmos/cosmos-sdk v0.42.6
	github.com/ethereum/go-ethereum v1.10.9
	github.com/gogo/protobuf v1.3.3
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-retryablehttp v0.5.3
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.13.0 // indirect
	github.com/prometheus/client_golang v1.10.0
	github.com/rs/zerolog v1.21.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/tendermint/spm v0.1.4
	github.com/tendermint/tendermint v0.34.11
	github.com/tendermint/tm-db v0.6.4
	gitlab.com/thorchain/tss/go-tss v1.4.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/sys v0.0.0-20210923061019-b8560ed6a9b7 // indirect
	google.golang.org/genproto v0.0.0-20211118181313-81c1377c94b1
	google.golang.org/grpc v1.42.0
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/binance-chain/tss-lib => gitlab.com/thorchain/tss/tss-lib v0.0.0-20201118045712-70b2cb4bf916
