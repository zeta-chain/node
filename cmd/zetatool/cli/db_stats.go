package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

// moduleStats represents statistics for a single module
type moduleStats struct {
	Count        int `json:"count"`        // Number of key-value pairs
	TotalKeySize int `json:"totalKeySize"` // Total size of all keys in bytes
	TotalValSize int `json:"totalValSize"` // Total size of all values in bytes
	AvgKeySize   int `json:"avgKeySize"`   // Average key size in bytes
	AvgValSize   int `json:"avgValSize"`   // Average value size in bytes
}

// databaseStats represents the overall database statistics
type databaseStats struct {
	TotalKeys    int                     `json:"totalKeys"`
	TotalKeySize int                     `json:"totalKeySize"`
	TotalValSize int                     `json:"totalValSize"`
	Modules      map[string]*moduleStats `json:"modules"`
}

// OutputFormat represents the supported output formats
type OutputFormat string

const (
	formatTable OutputFormat = "table"
	formatJSON  OutputFormat = "json"
)

var moduleRe = regexp.MustCompile(`s\/k:(\w+)\/`)

// NewApplicationDBStatsCMD creates a new cobra command for database statistics
func NewApplicationDBStatsCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db-stats",
		Short: "Show database statistics for the application",
		Long: `Show detailed statistics about the application database including:
- Module-wise statistics (key count, sizes)
- Total database size
- Average key and value sizes per module

The output can be formatted as a table (default) or JSON.`,
		RunE: runStatsCommand,
	}

	cmd.Flags().String("dbpath", "", "Path to the application DB directory (required)")
	cmd.Flags().String("format", string(formatTable), "Output format (table|json)")

	if err := cmd.MarkFlagRequired("dbpath"); err != nil {
		fmt.Println("Error marking flag as required")
		os.Exit(1)
	}

	return cmd
}

// runStatsCommand is the main entry point for the stats command
func runStatsCommand(cmd *cobra.Command, _ []string) error {
	dbPath, err := cmd.Flags().GetString("dbpath")
	if err != nil {
		return errors.Wrap(err, "failed to get dbpath")
	}

	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return errors.Wrap(err, "failed to get format")
	}

	// Open database
	db, err := openDatabase(dbPath)
	if err != nil {
		return errors.Wrap(err, "failed to open db")
	}
	defer db.Close()

	// Collect statistics
	stats, err := collectStats(db)
	if err != nil {
		return errors.Wrap(err, "failed to collect stats")
	}

	// Display results
	return displayStats(stats, OutputFormat(format))
}

// openDatabase opens the database in read-only mode
func openDatabase(dbPath string) (dbm.DB, error) {
	db, err := dbm.NewDB("application", dbm.PebbleDBBackend, dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// collectStats iterates through the database and collects statistics
func collectStats(db dbm.DB) (*databaseStats, error) {
	stats := &databaseStats{
		Modules: make(map[string]*moduleStats),
	}

	iter, err := db.Iterator(nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create iterator")
	}
	defer iter.Close()

	if err := iter.Error(); err != nil {
		return nil, errors.Wrap(err, "iteration error")
	}

	for ; iter.Valid(); iter.Next() {
		keySize := len(iter.Key())
		valSize := len(iter.Value())

		stats.TotalKeys++
		stats.TotalKeySize += keySize
		stats.TotalValSize += valSize

		moduleName := getModuleName(string(iter.Key()))
		updateModuleStats(stats.Modules, moduleName, keySize, valSize)
	}

	// Calculate averages for each module
	for _, moduleStats := range stats.Modules {
		if moduleStats.Count > 0 {
			moduleStats.AvgKeySize = moduleStats.TotalKeySize / moduleStats.Count
			moduleStats.AvgValSize = moduleStats.TotalValSize / moduleStats.Count
		}
	}

	return stats, nil
}

// getModuleName extracts the module name from a key
func getModuleName(key string) string {
	if strings.HasPrefix(key, "s/k:") {
		tokens := moduleRe.FindStringSubmatch(key)
		if len(tokens) > 1 {
			return tokens[1]
		}
	}
	return "misc"
}

// updateModuleStats updates statistics for a specific module
func updateModuleStats(modules map[string]*moduleStats, moduleName string, keySize, valSize int) {
	if modules[moduleName] == nil {
		modules[moduleName] = &moduleStats{}
	}

	stats := modules[moduleName]
	stats.Count++
	stats.TotalKeySize += keySize
	stats.TotalValSize += valSize
}

// displayStats renders the statistics in the specified format
func displayStats(stats *databaseStats, format OutputFormat) error {
	switch format {
	case formatTable:
		displayTable(stats)
	case formatJSON:
		return displayJSON(stats)
	default:
		return errors.New("unsupported format")
	}
	return nil
}

// displayTable renders the statistics in a formatted table
func displayTable(stats *databaseStats) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(
		table.Row{"Module", "Avg Key Size", "Avg Value Size", "Total Key Size", "Total Value Size", "Total Key Pairs"},
	)

	modules := maps.Keys(stats.Modules)
	sortSlice(modules)

	for _, m := range modules {
		s := stats.Modules[m]
		t.AppendRow([]interface{}{
			m,
			formatBytes(s.AvgKeySize),
			formatBytes(s.AvgValSize),
			formatBytes(s.TotalKeySize),
			formatBytes(s.TotalValSize),
			s.Count,
		})
	}

	t.AppendFooter(
		table.Row{
			"Total",
			"",
			"",
			formatBytes(stats.TotalKeySize),
			formatBytes(stats.TotalValSize),
			stats.TotalKeys,
		},
	)
	t.Render()
}

// displayJSON renders the statistics in JSON format
func displayJSON(stats *databaseStats) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(stats)
}

// sortSlice sorts a slice of ordered values
func sortSlice[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

// formatBytes formats a byte count into a human-readable string
func formatBytes(b int) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := int64(b) / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
