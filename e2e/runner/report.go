package runner

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"time"
)

// TestReport is a struct that contains the test report
type TestReport struct {
	Name     string
	Success  bool
	Time     time.Duration
	GasSpent AccountBalancesDiff
}

// TestReports is a slice of TestReport
type TestReports []TestReport

// String returns the string representation of the test report as a table
// it uses text/tabwriter to format the table
func (tr TestReports) String(prefix string) (string, error) {
	var b strings.Builder
	writer := tabwriter.NewWriter(&b, 0, 4, 4, ' ', 0)
	if _, err := fmt.Fprintln(writer, "Name\tSuccess\tTime\tSpent"); err != nil {
		return "", err
	}

	for _, report := range tr {
		spent := formatBalances(report.GasSpent)
		success := "✅"
		if !report.Success {
			success = "❌"
		}
		if _, err := fmt.Fprintf(writer, "%s%s\t%s\t%s\t%s\n", prefix, report.Name, success, report.Time, spent); err != nil {
			return "", err
		}
	}

	if err := writer.Flush(); err != nil {
		return "", err
	}
	return b.String(), nil
}

// PrintTestReports prints the test reports
func (runner *E2ERunner) PrintTestReports(tr TestReports) {
	runner.Logger.Print(" ---📈 E2E Test Report ---")
	table, err := tr.String("")
	if err != nil {
		runner.Logger.Print("Error rendering test report: %s", err)
	}
	runner.Logger.PrintNoPrefix(table)
}
