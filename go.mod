module github.com/zeta-chain/zetacore

go 1.16

require (
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/cosmos/cosmos-sdk v0.45.4
	github.com/cosmos/ibc-go/v2 v2.2.0
	github.com/ethereum/go-ethereum v1.10.16
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-retryablehttp v0.5.3
	github.com/ignite-hq/cli v0.20.4
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/prometheus/client_golang v1.12.1
	github.com/rs/zerolog v1.23.0
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.1
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/tendermint/tendermint v0.34.19
	github.com/tendermint/tm-db v0.6.7
	gitlab.com/thorchain/tss/go-tss v1.5.1-0.20220209042552-9900e94275ab
	google.golang.org/genproto v0.0.0-20220719170305-83ca9fad585f
	google.golang.org/grpc v1.48.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
)

require github.com/tendermint/spm v0.1.9

require (
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/lib/pq v1.10.6
	github.com/mattn/go-sqlite3 v1.14.9
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/sys v0.0.0-20220610221304-9f5ed59c137d // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/binance-chain/tss-lib => github.com/zeta-chain/tss-lib v0.1.3-0.20220721223335-682db65bafd1

replace gitlab.com/thorchain/tss/go-tss => github.com/zeta-chain/go-tss v1.5.2-0.20220721223537-74acdb1e3abf

replace github.com/agl/ed25519 => github.com/binance-chain/edwards25519 v0.0.0-20200305024217-f36fc4b53d43
