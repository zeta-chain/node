package snapshot

import (
	sdkmath "cosmossdk.io/math"
)

// zetaExp is the azeta -> ZETA exponent (18 decimals) for the human audit copy.
const zetaExp = 18

// DBHeader is the column order of snapshot_db.csv (Postgres COPY-ready).
var DBHeader = []string{
	"canonical_address", "class", "liquid", "staked", "unbonding",
	"total_azeta", "total_9dec", "claim_status",
}

// HumanHeader is the column order of the readable audit copy (ZETA units).
var HumanHeader = []string{
	"canonical_address", "class", "liquid_zeta", "staked_zeta", "unbonding_zeta",
	"total_zeta", "claim_status",
}

// SchemaSQL is the Postgres DDL matching snapshot_db.csv. azeta amounts are
// NUMERIC because raw azeta (up to ~1.7e27) overflows bigint; total_9dec is the
// folded 9-decimal amount and fits in bigint.
const SchemaSQL = `CREATE TABLE snapshot_balances (
    canonical_address TEXT PRIMARY KEY,
    class             TEXT    NOT NULL,
    liquid            NUMERIC NOT NULL,
    staked            NUMERIC NOT NULL,
    unbonding         NUMERIC NOT NULL,
    total_azeta       NUMERIC NOT NULL,
    total_9dec        BIGINT  NOT NULL,
    claim_status      TEXT    NOT NULL
);

CREATE INDEX idx_snapshot_balances_address ON snapshot_balances (canonical_address);
`

// DBRecords returns the header plus one row per claimable account and a final
// remainder row (class=remainder, routed to the multisig).
func (r *Result) DBRecords() [][]string {
	records := make([][]string, 0, len(r.Accounts)+2)
	records = append(records, DBHeader)
	for _, a := range r.Accounts {
		records = append(records, []string{
			a.Canonical,
			a.Class,
			a.Liquid.String(),
			a.Staked.String(),
			a.Unbonding.String(),
			a.Total().String(),
			a.Total9Dec().String(),
			claimStatusClaimable,
		})
	}
	records = append(records, []string{
		classRemainder,
		classRemainder,
		"0",
		"0",
		"0",
		r.Remainder.String(),
		r.Remainder.Quo(oneGZeta).String(),
		claimStatusMultisig,
	})
	return records
}

// HumanRecords mirrors DBRecords in ZETA units for auditing (not published).
func (r *Result) HumanRecords() [][]string {
	records := make([][]string, 0, len(r.Accounts)+2)
	records = append(records, HumanHeader)
	for _, a := range r.Accounts {
		records = append(records, []string{
			a.Canonical,
			a.Class,
			toZeta(a.Liquid),
			toZeta(a.Staked),
			toZeta(a.Unbonding),
			toZeta(a.Total()),
			claimStatusClaimable,
		})
	}
	records = append(records, []string{
		classRemainder,
		classRemainder,
		"0",
		"0",
		"0",
		toZeta(r.Remainder),
		claimStatusMultisig,
	})
	return records
}

// toZeta renders an azeta amount as a decimal ZETA string.
func toZeta(azeta sdkmath.Int) string {
	return sdkmath.LegacyNewDecFromIntWithPrec(azeta, zetaExp).String()
}
