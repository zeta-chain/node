# ZETA balance snapshot

Tooling to produce a per-address list of native ZETA balances from a ZetaChain
state export. It is the input to the ZETA to Solana migration mint, but the
compute stage is a plain, offline balance aggregator that anyone can run and
audit against testnet or mainnet.

This doc is the quick start. Grab the PR, run the two stages, poke at the output.

## What it produces

One CSV row per account with its native ZETA, split into `liquid / staked /
unbonding`, plus a single `remainder` row for everything not individually
attributed. Books balance to the wei against total supply.

## Design (two stages)

```
Stage 1  EXPORT     zetacored export        -> export.json      (contrib/snapshot, Python)
Stage 2  SNAPSHOT   zetatool snapshot       -> snapshot_db.csv  (Go, pure/offline)
```

- Stage 1 reads committed chain state and writes a genesis-format JSON. Chain
  specific, side-effecty, needs the matching binary.
- Stage 2 takes only `export.json` and computes balances. No network, fully
  deterministic, so the same code runs identically on testnet, devnet, mainnet.
  That is the boundary to build on: new attribution ideas go here without
  touching a node.

Key facts the design leans on:

- Native EVM balance IS the bank `azeta` balance (EVM StateDB is a view over
  bank), and `azeta` is 18 decimals == EVM wei 1:1. A bank export captures every
  EOA and every contract's native holdings.
- ERC-20 / WZETA per-holder balances live in EVM contract storage, not bank. A
  bank export shows only the aggregate native lump under a contract address.
- `--for-zero-height` folds staking rewards + validator commission into bank
  balances during export, so Stage 2 does not reimplement reward math. Trade-off:
  exported balances are higher than a live `balanceOf` by the withdrawn rewards.

Attribution model: attribute `liquid + staked + unbonding` to every account not
in the non-claimable set (module accounts, the zero address, and any pinned
contract). `remainder = total_supply - sum(attributed)` goes to a single row.
Nothing is dropped; the cap equals full supply.

## Prerequisites

- `python3`, `lz4`, `make`, Go toolchain.
- Disk: mainnet snapshot is ~35 GB compressed, ~39 GB extracted, and Stage 1
  copies the data dir once (~39 GB more). Budget ~120 GB free for mainnet.
- The binary must match the state version. Mainnet is currently on v36; a newer
  binary can fail to open un-upgraded state. Check the current version in the
  snapshot metadata:
  `curl -s https://snapshots.rpc.zetachain.com/mainnet/fullnode/latest.json`
  (`networkVersion` field). Then build the matching release:
  `git checkout v36.1.5 && make install-zetacore` (use a worktree so you keep your
  branch). Testnet was validated on v37.0.2.

## Stage 1: export

Testnet (downloads + exports latest):

```
make snapshot-export NETWORK=testnet SNAPSHOT_EXPORT_ARGS="--for-zero-height --output /tmp/testnet-export.json"
```

Mainnet:

```
make snapshot-export NETWORK=mainnet SNAPSHOT_EXPORT_ARGS="--for-zero-height --output /tmp/mainnet-export.json"
```

Or call the script directly for more control:

```
python3 contrib/snapshot/export_snapshot.py --network mainnet --for-zero-height \
  --output /tmp/mainnet-export.json [--skip-download] [--dry-run]
```

Useful flags:

- `--for-zero-height`  fold rewards + commission into bank. Required: Stage 2
                       refuses an export with unpaid rewards.
- `--modules`          default `bank,staking,distribution`. Stage 2 needs all three.
- `--height`           default `-1` (latest).
- `--skip-download`    reuse the cached snapshot in `~/zetacored_snapshot_<network>`.
- `--dry-run`          print every command, execute nothing.
- `--node-version`     only stamps a version string; it does NOT check out or build
                       that release. Build the matching binary yourself first (see
                       Prerequisites).

Gotcha handled by the script: `zetacored export` writes an early WARN line to
stdout ahead of the JSON, so the raw output is stripped to the first `{`.

## Stage 2: snapshot compute

```
go run ./cmd/zetatool snapshot --input /tmp/mainnet-export.json --chain-id 7000 --output /tmp/mainnet-snap
```

Flags:

- `--input`         path to `export.json` (required).
- `--output`        output directory (default `.`).
- `--chain-id`      codec chain id: testnet `7001`, mainnet `7000`.
- `--pin` / `--wzeta`  non-claimable contract addresses routed to the remainder
                    (repeatable or comma-separated). WZETA and known pools go here.
- `--pin-file`      same thing from a file, one address per line (`#` comments;
                    only the first field per line is read, so a TSV works).
                    Mainnet list is committed: `contrib/snapshot/pins_mainnet.txt`.
- `--supply-check`  fail on any output check (default true). It does not cover
                    the hard refusals in "Acceptance guarantees" below.

Mainnet with the committed pin list:

```
go run ./cmd/zetatool snapshot --input /tmp/mainnet-export.json --chain-id 7000 \
  --output /tmp/mainnet-snap --pin-file contrib/snapshot/pins_mainnet.txt
```

Outputs in `--output`:

- `snapshot_db.csv`   one row per attributed address + one `remainder` row.
  Columns: `canonical_address, class, liquid, staked, unbonding, total_azeta,
  total_9dec, claim_status`. Postgres COPY-ready (`schema.sql`), `azeta` as
  `NUMERIC` (raw azeta exceeds bigint). `total_9dec` = floored 18->9; SPL cap =
  sum of `total_9dec`.
- `schema.sql`        matching Postgres DDL.
- `snapshot_human.csv`  readable ZETA-unit copy for eyeballing.

The command prints bucket totals and runs four checks (sum export == supply; sum
attributed + remainder == total; re-sum CSV == total; every pin matched a holder),
non-zero exit on any failure.

Pin accounting: the summary prints one line per pin with the azeta it kept out of
attribution, plus a subtotal. A pin that matched no holder in the export reports
zero and fails the run, which is what catches a typo'd address or a pin file
pointed at the wrong network. `--supply-check=false` downgrades these four checks,
including this one, to a stderr warning.

The summary also splits the remainder by category (module accounts, zero
address, non-20-byte holders, pins, share-rounding dust) and prints the azeta
lost to the 18->9 decimal floor.

## Acceptance guarantees

Before writing anything, Stage 2 refuses an export it cannot fully account for.
These are hard errors. `--supply-check=false` does not turn them off.

- The export has no `distribution` module.
- The staking `bond_denom` is not `azeta`.
- A validator is listed twice, has an invalid status, or has negative tokens or
  delegator shares.
- A delegation has negative shares, points at a validator not in the export, or
  repeats a (delegator, validator) pair.
- A validator's delegation shares do not sum exactly to its `delegator_shares`.
  This is what catches a dropped delegation.
- A validator's tokens exceed what its delegators are credited by more than one
  azeta per delegation (rounding).
- An unbonding record repeats a (delegator, validator) pair, or has a negative
  entry.
- The bonded pool balance is not exactly the bonded validators' tokens.
- The not-bonded pool balance is not exactly the unbonding and unbonded
  validators' tokens plus every unbonding entry. This is what catches a dropped
  unbonding record.
- Any validator has outstanding rewards, current rewards or accumulated
  commission. That means the export was not taken with `--for-zero-height`.
- The distribution module balance is not exactly the community pool.
- A pin is one of the staking pools.
- The remainder is not exactly the sum of its categories.

What these checks do NOT prove, and the operator still has to show for the
final export:

- The emissions collection address and the contract inventory are complete.
  Both go to the foundation treasury, so they need checking once by hand.
- The remainder reaches the treasury, once.
- The export came from the agreed chain and height.

## Look at the output

DuckDB queries the CSV in place (handles the big NUMERICs; `brew install duckdb`):

```
duckdb -box -c "SELECT class, count(*) n, round(sum(total_azeta)/1e18,2) zeta \
  FROM read_csv_auto('/tmp/mainnet-snap/snapshot_db.csv', types={'total_azeta':'DECIMAL(38,0)'}) \
  GROUP BY class ORDER BY zeta DESC"
```

Browser UI: `duckdb -ui` then `CREATE TABLE s AS SELECT * FROM read_csv_auto(...)`.

Or load `schema.sql` + `\copy` the CSV into Postgres and use any GUI.

## Contract vs EOA (the WZETA / pool question)

Accounts export as cosmos `BaseAccount` with no code hash, so you CANNOT tell a
contract from an EOA in the export. Detect contracts out-of-band with
`eth_getCode` against an EVM RPC (contracts return bytecode, EOAs return `0x`),
then feed the contract addresses to Stage 2 via `--pin`. Contracts hold native
that no one can claim on Solana, so pinning routes it to the remainder.

Measured on mainnet (top 500 holders, H17619782): the 8 largest holders are EOAs
(~1.0B ZETA, claimable); 16 are contracts holding 9.35M ZETA (~0.55%), dominated
by WZETA `0x5f0b1a82749cb4e2278ec87f8bf6b618dc71a8bf` (6.6M) and the Accumulated
Finance liquid-staking / lending stack. WZETA is 1:1 backed native, so the native
backing for all DEX/LP WZETA positions sits under that one contract.

Those 16 addresses are committed as `contrib/snapshot/pins_mainnet.txt` and wired
in with `--pin-file` (see Stage 2). The file is mainnet only; a testnet run needs
no pins.

## Try new ideas here

- Change the non-claimable set: `--pin <addr>,<addr>` to route more contracts to
  the remainder.
- Swap the attribution rule (e.g. exclude unbonding, or split rewards out) in the
  pure package `cmd/zetatool/snapshot/` and re-run against the same `export.json`.
- Point Stage 2 at a devnet export to test on a small state.

## Verification status

- Covered by unit tests: pin-file parsing (comments, blanks, TSV columns),
  per-pin azeta accounting across liquid + staked + unbonding, and the
  unmatched-pin failure.
- NOT run yet: the end-to-end mainnet run with `pins_mainnet.txt`. Expected when
  it happens — claimable accounts down by 16, remainder 41.59M -> ~50.94M ZETA,
  total supply unchanged at 1,707,671,363.36 ZETA. Deferred, so treat those
  numbers as a prediction, not a result.

## Open questions / parked

- Full contract detection: the top-500 contract list is committed and wired via
  `--pin-file`, which captures ~all contract value. What is still parked is the
  exhaustive all-accounts `eth_getCode` sweep against a local node, the
  pre-migration step that proves nothing was missed.
- Squads multisig address for the remainder recipient (pending tokenomics/legal).
- Snapshot height selection for the real migration (freeze inbound, drain
  outbounds, zero emissions, halt, export). Separate workstream.
- Per-holder WZETA / DEX / LP attribution and the claim service / Solana program
  are out of scope for the snapshot.
