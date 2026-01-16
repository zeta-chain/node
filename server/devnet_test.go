package server

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBinaryInfo(t *testing.T) {

	tests := []struct {
		name            string
		upgradeVersion  string
		goos            string
		goarch          string
		wantPlatform    string
		wantZetacored   string
		wantZetaClientd string
	}{
		{
			name:           "v37 local",
			upgradeVersion: "v37",
			goos:           runtime.GOOS,
			goarch:         runtime.GOARCH,
			wantPlatform:   runtime.GOOS + "/" + runtime.GOARCH,
			wantZetacored: fmt.Sprintf(
				"https://github.com/zeta-chain/node/releases/download/v37/zetacored-%s-%s",
				runtime.GOOS,
				runtime.GOARCH,
			),
			wantZetaClientd: fmt.Sprintf(
				"https://github.com/zeta-chain/node/releases/download/v37/zetaclientd-%s-%s",
				runtime.GOOS,
				runtime.GOARCH,
			),
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
			require.Equal(t, tt.wantZetaClientd, zetaclientdURL)
		})
	}
}
