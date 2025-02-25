package config

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("FromSource_path", func(t *testing.T) {
		// ARRANGE
		tmpDir := t.TempDir()
		defer os.RemoveAll(tmpDir)

		configPath := filepath.Join(tmpDir, "config.json")
		require.NoError(t, os.WriteFile(configPath, sampleConfigBytes(t), 0644))

		// ACT
		cfg, err := FromSource(ctx, configPath)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, cfg)
	})

	t.Run("FromSource_url", func(t *testing.T) {
		// ARRANGE
		cfgBytes := sampleConfigBytes(t)

		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Write(cfgBytes)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		url := fmt.Sprintf("%s/config.json", server.URL)

		// ACT
		cfg, err := FromSource(ctx, url)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, cfg)
	})
}

func sampleConfigBytes(t *testing.T) []byte {
	// https://ton.org/testnet-global.config.json
	const sampleConfig = `{
		"liteservers": [
			{
				"ip": 822907680,
				"port": 27842,
				"provided":"Beavis",
				"id": {
					"@type": "pub.ed25519",
					"key": "sU7QavX2F964iI9oToP9gffQpCQIoOLppeqL/pdPvpM="
				}
			},
			{
				"ip": 1091956407,
				"port": 16351,
				"id": {
					"@type": "pub.ed25519",
					"key": "Mf/JGvcWAvcrN3oheze8RF/ps6p7oL6ifrIzFmGQFQ8="
				}
			}
		]
	}`

	return []byte(sampleConfig)
}
