package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	sdkmath "cosmossdk.io/math"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/cmd/zetatool/snapshot"
)

const (
	flagSnapshotInput       = "input"
	flagSnapshotOutput      = "output"
	flagSnapshotChainID     = "chain-id"
	flagSnapshotPin         = "pin"
	flagSnapshotWZeta       = "wzeta"
	flagSnapshotSupplyCheck = "supply-check"

	fileSnapshotDB     = "snapshot_db.csv"
	fileSnapshotHuman  = "snapshot_human.csv"
	fileSnapshotSchema = "schema.sql"

	// column index of total_azeta in snapshot_db.csv, used to re-sum the file.
	dbColTotalAzeta = 5
)

// NewSnapshotCMD builds the `zetatool snapshot` command: it reads a
// `zetacored export` genesis and produces a per-address native-ZETA balance
// list for the ZETA->Solana migration mint.
func NewSnapshotCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Build a per-address native-ZETA balance snapshot for the Solana migration",
		Long: `Build the ZETA->Solana migration snapshot (Stage 2) from a zetacored export.

Reads the bank and staking genesis from an export produced with
--for-zero-height (so withdrawn staking rewards and validator commission are
already folded into liquid bank balances) and attributes native ZETA to each
claimable address. Module accounts and any pinned contracts (e.g. WZETA) are
routed to a single remainder that goes to the migration multisig.

Outputs into --output:
  snapshot_db.csv    Postgres COPY-ready (raw azeta + 9-decimal SPL amount)
  schema.sql         matching Postgres table DDL
  snapshot_human.csv readable audit copy in ZETA units (not published)`,
		RunE: runSnapshot,
	}

	cmd.Flags().String(flagSnapshotInput, "", "path to the zetacored export.json (required)")
	cmd.Flags().String(flagSnapshotOutput, ".", "output directory")
	cmd.Flags().Uint64(flagSnapshotChainID, 7001, "chain id for the codec (testnet=7001, mainnet=7000)")
	cmd.Flags().
		StringSlice(flagSnapshotPin, nil, "non-claimable addresses routed to the remainder (repeatable or comma-separated)")
	cmd.Flags().StringSlice(flagSnapshotWZeta, nil, "alias for --pin: non-claimable contract addresses (e.g. WZETA)")
	cmd.Flags().Bool(flagSnapshotSupplyCheck, true, "fail on any supply-check mismatch")

	if err := cmd.MarkFlagRequired(flagSnapshotInput); err != nil {
		panic(err)
	}
	return cmd
}

func runSnapshot(cmd *cobra.Command, _ []string) error {
	input, err := cmd.Flags().GetString(flagSnapshotInput)
	if err != nil {
		return err
	}
	outputDir, err := cmd.Flags().GetString(flagSnapshotOutput)
	if err != nil {
		return err
	}
	chainID, err := cmd.Flags().GetUint64(flagSnapshotChainID)
	if err != nil {
		return err
	}
	pins, err := cmd.Flags().GetStringSlice(flagSnapshotPin)
	if err != nil {
		return err
	}
	wzeta, err := cmd.Flags().GetStringSlice(flagSnapshotWZeta)
	if err != nil {
		return err
	}
	supplyCheck, err := cmd.Flags().GetBool(flagSnapshotSupplyCheck)
	if err != nil {
		return err
	}

	appState, err := readAppState(input)
	if err != nil {
		return err
	}

	encodingConfig := app.MakeEncodingConfig(chainID)
	gen, err := snapshot.ParseAppState(encodingConfig.Codec, appState)
	if err != nil {
		return err
	}

	cfg := snapshot.Config{Denom: snapshot.BaseDenom, Pins: append(pins, wzeta...)}
	res, err := snapshot.Compute(gen, cfg)
	if err != nil {
		return err
	}

	if err := writeOutputs(outputDir, res); err != nil {
		return err
	}

	if err := runChecks(gen, res, filepath.Join(outputDir, fileSnapshotDB), supplyCheck); err != nil {
		return err
	}

	printSummary(cmd, gen, res)
	return nil
}

// readAppState reads the export file and returns its decoded app_state map.
func readAppState(path string) (map[string]json.RawMessage, error) {
	raw, err := os.ReadFile(path) // #nosec G304 -- operator-supplied export path
	if err != nil {
		return nil, fmt.Errorf("read input %q: %w", path, err)
	}
	var doc struct {
		AppState map[string]json.RawMessage `json:"app_state"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("unmarshal export json: %w", err)
	}
	if doc.AppState == nil {
		return nil, fmt.Errorf("export json has no app_state")
	}
	return doc.AppState, nil
}

func writeOutputs(dir string, res *snapshot.Result) error {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("create output dir %q: %w", dir, err)
	}
	if err := writeCSV(filepath.Join(dir, fileSnapshotDB), res.DBRecords()); err != nil {
		return err
	}
	if err := writeCSV(filepath.Join(dir, fileSnapshotHuman), res.HumanRecords()); err != nil {
		return err
	}
	schemaPath := filepath.Join(dir, fileSnapshotSchema)
	if err := os.WriteFile(schemaPath, []byte(snapshot.SchemaSQL), 0o600); err != nil {
		return fmt.Errorf("write %q: %w", schemaPath, err)
	}
	return nil
}

func writeCSV(path string, records [][]string) error {
	f, err := os.Create(path) // #nosec G304 -- operator-supplied output dir
	if err != nil {
		return fmt.Errorf("create %q: %w", path, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.WriteAll(records); err != nil {
		return fmt.Errorf("write %q: %w", path, err)
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("flush %q: %w", path, err)
	}
	return nil
}

// runChecks enforces the three supply invariants. When supplyCheck is false the
// mismatches are reported but do not fail the command.
func runChecks(gen *snapshot.Genesis, res *snapshot.Result, dbPath string, supplyCheck bool) error {
	var problems []string

	// (1) every exported bank azeta balance is accounted for
	if bankSum := snapshot.SumBankDenom(gen, snapshot.BaseDenom); !bankSum.Equal(res.Supply) {
		problems = append(problems, fmt.Sprintf("bank balances sum %s != supply %s", bankSum, res.Supply))
	}

	// (2) attributed + remainder == supply (remainder >= 0 is enforced in Compute)
	if got := res.AttributedTotal().Add(res.Remainder); !got.Equal(res.Supply) {
		problems = append(problems, fmt.Sprintf("attributed+remainder %s != supply %s", got, res.Supply))
	}

	// (3) the written CSV re-sums to the supply
	csvSum, err := sumCSVColumn(dbPath, dbColTotalAzeta)
	if err != nil {
		return err
	}
	if !csvSum.Equal(res.Supply) {
		problems = append(problems, fmt.Sprintf("csv total_azeta sum %s != supply %s", csvSum, res.Supply))
	}

	if len(problems) == 0 {
		return nil
	}
	for _, p := range problems {
		fmt.Fprintf(os.Stderr, "supply check failed: %s\n", p)
	}
	if supplyCheck {
		return fmt.Errorf("%d supply check(s) failed", len(problems))
	}
	return nil
}

// sumCSVColumn re-reads a written CSV and sums an integer column (skips header).
func sumCSVColumn(path string, col int) (sdkmath.Int, error) {
	f, err := os.Open(path) // #nosec G304 -- path we just wrote
	if err != nil {
		return sdkmath.Int{}, fmt.Errorf("reopen %q: %w", path, err)
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return sdkmath.Int{}, fmt.Errorf("read %q: %w", path, err)
	}
	total := sdkmath.ZeroInt()
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if col >= len(row) {
			return sdkmath.Int{}, fmt.Errorf("%q row %d has no column %d", path, i, col)
		}
		v, ok := sdkmath.NewIntFromString(row[col])
		if !ok {
			return sdkmath.Int{}, fmt.Errorf("%q row %d: bad integer %q", path, i, row[col])
		}
		total = total.Add(v)
	}
	return total, nil
}

func printSummary(cmd *cobra.Command, gen *snapshot.Genesis, res *snapshot.Result) {
	out := cmd.OutOrStdout()
	fmt.Fprintln(out, "\nZETA->Solana migration snapshot")
	fmt.Fprintf(out, "  bank balances:   %d\n", len(gen.Bank.Balances))
	fmt.Fprintf(out, "  validators:      %d\n", len(gen.Staking.Validators))
	fmt.Fprintf(out, "  claimable accts: %d\n", len(res.Accounts))
	if res.NonStandard > 0 {
		fmt.Fprintf(out, "  non-std swept:   %d holder(s), %s azeta -> remainder\n",
			res.NonStandard, res.NonStandardAzeta)
	}
	fmt.Fprintln(out, "  buckets (azeta):")
	fmt.Fprintf(out, "    liquid:        %s\n", res.TotalLiquid)
	fmt.Fprintf(out, "    staked:        %s\n", res.TotalStaked)
	fmt.Fprintf(out, "    unbonding:     %s\n", res.TotalUnbonding)
	fmt.Fprintf(out, "    remainder:     %s\n", res.Remainder)
	fmt.Fprintf(out, "  total supply:    %s azeta\n", res.Supply)
	fmt.Fprintf(out, "  SPL cap (9dec):  %s\n", res.SPLCap())
}
