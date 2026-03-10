package cli

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/term"
)

// newTableWriter creates a table.Writer configured to fit the terminal width.
// If the terminal width cannot be detected, no width limit is applied.
func newTableWriter() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && width > 0 {
		t.Style().Size.WidthMax = width
	}

	return t
}
