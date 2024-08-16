package os

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// ExpandHomeDir expands a leading tilde in the path to the home directory of the current user.
// ~someuser/tmp will not be expanded.
func ExpandHomeDir(p string) (string, error) {
	if p == "~" ||
		strings.HasPrefix(p, "~/") ||
		strings.HasPrefix(p, "~\\") {
		usr, err := user.Current()
		if err != nil {
			return p, err
		}

		p = filepath.Join(usr.HomeDir, p[1:])
	}
	return filepath.Clean(p), nil
}

// FileExists checks if a file exists.
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
