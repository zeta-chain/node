package server

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBinaryInfo(t *testing.T) {
	tests := []struct {
		name           string
		upgradeVersion string
		goos           string
		goarch         string
		wantPlatform   string
		wantZetacored  string
		wantClient     string
	}{
		{
			name:           "v37 linux amd64",
			upgradeVersion: "v37",
			goos:           "linux",
			goarch:         "amd64",
			wantPlatform:   "linux/amd64",
			wantZetacored:  "https://github.com/zeta-chain/node/releases/download/v37/zetacored-ubuntu-22-amd64",
			wantClient:     "https://github.com/zeta-chain/node/releases/download/v37/zetaclientd-ubuntu-22-amd64",
		},
		{
			name:           "v37 darwin arm64",
			upgradeVersion: "v37",
			goos:           "darwin",
			goarch:         "arm64",
			wantPlatform:   "darwin/arm64",
			wantZetacored:  "https://github.com/zeta-chain/node/releases/download/v37/zetacored-ubuntu-22-arm64",
			wantClient:     "https://github.com/zeta-chain/node/releases/download/v37/zetaclientd-ubuntu-22-arm64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			infoBytes, err := getBinaryInfo(tt.upgradeVersion, tt.goos, tt.goarch)

			// Assert
			require.NoError(t, err)
			var result map[string]interface{}
			err = json.Unmarshal(infoBytes, &result)
			require.NoError(t, err)

			require.Contains(t, result, "binaries")
			binaries, ok := result["binaries"].(map[string]interface{})
			require.True(t, ok)

			zetacoredURL, ok := binaries[tt.wantPlatform].(string)
			require.True(t, ok)
			require.Equal(t, tt.wantZetacored, zetacoredURL)

			clientKey := "zetaclientd-" + tt.wantPlatform
			zetaclientdURL, ok := binaries[clientKey].(string)
			require.True(t, ok)
			require.Equal(t, tt.wantClient, zetaclientdURL)
		})
	}
}
