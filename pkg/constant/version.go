package constant

import "strings"

var (
	Name       = ""
	Version    = ""
	CommitHash = ""
	BuildTime  = ""
)

const versionWhenBuiltWithoutMake = "v0.0.0+nomake"

// GetNormalizedVersion applies NormalizeVersion to the Version constant
// and ensures that it is semver compliant
func GetNormalizedVersion() string {
	version := Version
	if version == "" {
		version = versionWhenBuiltWithoutMake
	}
	return NormalizeVersion(version)
}

// NormalizeVersion ensures that the version string starts with a v
func NormalizeVersion(version string) string {
	return ensurePrefix(version, "v")
}

func ensurePrefix(s, prefix string) string {
	if !strings.HasPrefix(s, prefix) {
		return prefix + s
	}
	return s
}
