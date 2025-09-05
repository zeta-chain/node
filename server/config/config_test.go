package config_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	serverconfig "github.com/cosmos/evm/server/config"
	"github.com/cosmos/evm/testutil/constants"
)

func TestDefaultConfig(t *testing.T) {
	cfg := serverconfig.DefaultConfig()
	require.False(t, cfg.JSONRPC.Enable)
	require.Equal(t, cfg.JSONRPC.Address, serverconfig.DefaultJSONRPCAddress)
	require.Equal(t, cfg.JSONRPC.WsAddress, serverconfig.DefaultJSONRPCWsAddress)
}

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name    string
		args    func() *viper.Viper
		want    func() serverconfig.Config
		wantErr bool
	}{
		{
			"test unmarshal embedded structs",
			func() *viper.Viper {
				v := viper.New()
				v.Set("minimum-gas-prices", fmt.Sprintf("100%s", constants.ExampleAttoDenom))
				return v
			},
			func() serverconfig.Config {
				cfg := serverconfig.DefaultConfig()
				cfg.MinGasPrices = fmt.Sprintf("100%s", constants.ExampleAttoDenom)
				return *cfg
			},
			false,
		},
		{
			"test unmarshal EVMConfig",
			func() *viper.Viper {
				v := viper.New()
				v.Set("evm.tracer", "struct")
				return v
			},
			func() serverconfig.Config {
				cfg := serverconfig.DefaultConfig()
				require.NotEqual(t, "struct", cfg.EVM.Tracer)
				cfg.EVM.Tracer = "struct"
				return *cfg
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serverconfig.GetConfig(tt.args())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want()) {
				t.Errorf("GetConfig() got = %v, want %v", got, tt.want())
			}
		})
	}
}
