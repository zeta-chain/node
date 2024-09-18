package observer

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteapi"
)

// todo tmp (will be resolved automatically)
// taken from ton:8000/lite-client.json
const configRaw = `{"@type":"config.global","dht":{"@type":"dht.config.global","k":3,"a":3,"static_nodes":
{"@type":"dht.nodes","nodes":[]}},"liteservers":[{"id":{"key":"+DjLFqH/N5jO1ZO8PYVYU6a6e7EnnsF0GWFsteE+qy8=","@type":
"pub.ed25519"},"port":4443,"ip":2130706433}],"validator":{"@type":"validator.config.global","zero_state":
{"workchain":-1,"shard":-9223372036854775808,"seqno":0,"root_hash":"rR8EFZNlyj3rfYlMyQC8gT0A6ghDrbKe4aMmodiNw6I=",
"file_hash":"fT2hXGv1OF7XDhraoAELrYz6wX3ue16QpSoWTiPrUAE="},"init_block":{"workchain":-1,"shard":-9223372036854775808,
"seqno":0,"root_hash":"rR8EFZNlyj3rfYlMyQC8gT0A6ghDrbKe4aMmodiNw6I=",
"file_hash":"fT2hXGv1OF7XDhraoAELrYz6wX3ue16QpSoWTiPrUAE="}}}`

func TestObserver(t *testing.T) {
	t.Skip("skip test")

	ctx := context.Background()

	cfg, err := config.ParseConfig(strings.NewReader(configRaw))
	require.NoError(t, err)

	client, err := liteapi.NewClient(liteapi.WithConfigurationFile(*cfg))
	require.NoError(t, err)

	res, err := client.GetMasterchainInfo(ctx)
	require.NoError(t, err)

	// Outputs:
	// {
	//          "Last": {
	//            "Workchain": 4294967295,
	//            "Shard": 9223372036854775808,
	//            "Seqno": 915,
	//            "RootHash": "2e9e312c5bd3b7b96d23ce1342ac76e5486012c9aac44781c2c25dbc55f5c8ad",
	//            "FileHash": "d3745319bfaeebb168d9db6bb5b4752b6b28ab9041735c81d4a02fc820040851"
	//          },
	//          "StateRootHash": "02538fb9dc802004012285a90a7af9ba279706e2deea9ca635decd80e94a7045",
	//          "Init": {
	//            "Workchain": 4294967295,
	//            "RootHash": "ad1f04159365ca3deb7d894cc900bc813d00ea0843adb29ee1a326a1d88dc3a2",
	//            "FileHash": "7d3da15c6bf5385ed70e1adaa0010bad8cfac17dee7b5e90a52a164e23eb5001"
	//          }
	//        }
	t.Logf("Masterchain info")
	logJSON(t, res)
}

func logJSON(t *testing.T, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)

	t.Log(string(b))
}
