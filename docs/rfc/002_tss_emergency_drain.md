# RFC 002 — Emergency TSS native-fund drain (crosschain shutdown)

**Status:** Draft — for review
**Scope:** zetaclient only. No zetacore state or message changes (one existing query is reused).

## Motivation

We are shutting down crosschain traffic and need to move **all native funds held by the TSS — on every EVM chain and on Bitcoin** (ETH on Ethereum, the gas token on each EVM chain, BTC on Bitcoin) — to a safe wallet. The safe wallet is a plain address per chain family — **not** a TSS, and **not** derived from a keygen.

The built-in migration path can't do this:

- `MsgMigrateTssFunds` always migrates to the **latest registered TSS** (`newTss := tssList[len-1]`); the receiver is derived from a TSS pubkey and there is no destination field. It cannot target an arbitrary wallet.
- Making it target an arbitrary wallet — or adding any new admin message — is a zetacore change, i.e. a coordinated node upgrade. We want to avoid that.

The low-level primitive we need already exists and runs in production: `Signer.Sign(nil, to, amount, gas, nonce)` (`zetaclient/chains/evm/signer/signer.go:178`) builds a `LegacyTx` with empty calldata — a native transfer to any address — and signs it with the TSS. `signMigrateTssFundsCmd` (`signer_admin.go:139`) already calls it this way. The only reason it can't send to an arbitrary address today is that the *destination* is decided by the on-chain CCTX builder, not the signer.

So the drain is achievable entirely on the zetaclient side, provided we solve one problem: **coordinating N independent signers onto one identical transaction.**

## Goals / non-goals

**Goals**
- Drain ~100% of native TSS balance (minus fee + buffer) from every EVM chain and Bitcoin to a safe wallet per chain family.
- Zero zetacore changes — no new module state, no new message. (One existing read query, `GetBlockHeight`, is reused for timing.)
- Minimal, auditable zetaclient code. The client makes **no** decisions about the transaction.
- A phased path: localnet e2e → testnet → mainnet.

**Non-goals**
- Solana / TON / Sui native funds. This RFC covers EVM chains and Bitcoin only (the two families the e2e framework supports).
- A permanent, decentralized capability. This is a one-time operation; the trigger source is intentionally centralized and removed afterward (see Teardown).

## How TSS signing coordinates today (the constraint)

Every zetaclient signs independently, yet they produce one signature because all three inputs below are identical across nodes — and today all three come from the shared pending CCTX on zetacore:

1. **What to sign.** Each node reads the same CCTX, deterministically builds the same EVM tx (`to`, `amount`, `nonce`, `gasPrice`, `gasLimit`), and hashes it → identical digest.
2. **Ceremony matching.** go-tss pairs signers by `MsgID = hash(version + digests + signerPubKeys)` (`keysign/request.go:34`). Same digest → same ceremony. **`MsgID` does not include block height.**
3. **Leader election + timing.** `PickLeader(msgID, blockHeight, peers)` (`go-tss p2p/protocol.go:45`) picks the leader as the min of `hash(msgID, blockHeight, peerID)`. This means **`blockHeight` must also be identical across nodes**, or each node computes a different leader and the join-party never converges (fails at the 20s `KeySignTimeout` with no signature). Today `blockHeight` is derived from the shared zeta block via `KeysignHeight(chainID, zetaHeight)`, and nodes fire together because the scheduler triggers on `zetaHeight % scheduleInterval == 0`.

**Takeaway:** to drain without a CCTX, we must reproduce all three: an identical, fully-resolved tx (→ identical digest), an identical `blockHeight` value (→ same leader), and a near-simultaneous trigger (within the ~20s join window). Any single divergent byte and the ceremony silently fails to form.

## Design

Split the work so **all computation and every non-deterministic input live off-chain in a program we control**, and the zetaclient becomes a dumb signer:

```
┌─────────────────────────────┐        signed JSON         ┌──────────────────────────┐
│  Drain API (we host)        │  ───────────────────────▶  │  zetaclient (each node)  │
│  - operator supplies HEIGHT │   {height, per-chain        │  - poll + verify sig     │
│  - queries balance/gas/nonce│    fully-resolved tx body}  │  - at height: Sign+bcast │
│  - computes amount/fee      │                             │  - NO computation        │
└─────────────────────────────┘                            └──────────────────────────┘
```

### 1. Drain API — autonomous, operator supplies only the height

A small program (a new `zetatool` subcommand exposed over HTTP) that computes the **fully-resolved transaction body** for every EVM chain. The operator's **only** input is the trigger zeta-height; the API derives everything else itself:

- **Balance** — reuse the existing plumbing in `cmd/zetatool/cli/tss_balances.go` (`clients.GetEVMBalance` against the configured RPC per chain).
- **Gas price** — query **zetacore's median gas price** for the chain (the same value `MigrateFundCmdCCTX` uses), then apply the existing migration formula. This is a zetacore *query* on the API side only; the client never computes it.
- **Amount / fee** — reuse the exact production formula from `MigrateFundCmdCCTX` (`x/crosschain/types/cmd_cctxs.go:96-116`):
  ```
  gasLimit = 21_000                              (pkg/gas.EVMSend)
  gasPrice = medianGasPrice × 2.5               (TssMigrationGasMultiplierEVM)
  fee      = gasLimit × gasPrice + 2.1 gwei     (TSSMigrationBufferAmountEVM)
  amount   = balance − fee
  ```
  To keep a **single source of truth** for these constants, extract lines 96-116 into an exported pure function, e.g. `gas.ComputeEVMMigration(balance, medianGasPrice) (amount, gasPrice, gasLimit, error)`, and call it from both `MigrateFundCmdCCTX` and the new tool.
- **Nonce** — query the TSS account nonce once (`NonceAt(tssAddr)`). Because inbound is disabled and there are no pending outbounds during shutdown, the nonce is stable, so pinning it in the payload is safe and correct.

The API **signs** the payload with an operator/admin key; the zetaclient verifies against a baked-in public key.

#### Auto-updating amounts: draft → freeze → final

Amounts can't be committed too early — a gas spike between publish and mining would leave the tx underpriced and it would fail. But a payload that keeps changing also breaks signing: two clients polling seconds apart would sign different amounts → different digests → the ceremony never forms. Resolved with a **draft → final** model:

- The cron **recomputes and republishes continuously** (fresh median gas → fresh amounts) as **draft** payloads. Drafts are for monitoring only and **never trigger signing**.
- A window before the trigger height `H` (`H − K` blocks), the cron publishes exactly **one `final` payload** (marked `final: true`) and stops updating amounts.
- **Clients only ever sign a payload marked `final`.** The `final` is the single authoritative version everyone signs.
- Because TSS is **threshold, not unanimous**, a straggler that missed the `final` before `H` simply doesn't join — it doesn't break the ceremony, as long as ≥ threshold of clients hold the `final`.
- The existing **buffer** (`TSSMigrationBufferAmountEVM`) absorbs the small gas drift between freeze and the tx mining.

`K` is sized comfortably larger than the client poll interval so the `final` reliably propagates.

### 2. Payload — the fully-resolved, byte-final tx body

```json
{
  "trigger_zeta_height": 1234567,
  "seq": 12,
  "final": true,
  "evm_txs": [
    { "chain_id": 1,  "to": "0xSAFE", "nonce": 42, "amount": "…wei", "gas_price": "…wei", "gas_limit": 21000 }
  ],
  "btc_txs": [
    { "chain_id": 8332, "to": "bc1qSAFE", "output_sats": "…", "fee_sats": "…",
      "inputs": [ { "txid": "…", "vout": 0, "amount_sats": "…" }, { "txid": "…", "vout": 1, "amount_sats": "…" } ] }
  ],
  "signature": "0x…"
}
```

Every field the tx hash / sighash depends on is present and final. The zetaclient computes nothing — for EVM it copies values into a `LegacyTx`; for BTC it builds the exact `wire.MsgTx` from the **pinned inputs and output**. `trigger_zeta_height` doubles as the leader-election `blockHeight` (fed to `PickLeader`) **and** the sync clock. `seq` is a monotonic version for observability; `final` gates whether the payload may trigger signing (drafts are `false`). A TSS with >20 UTXOs produces **multiple `btc_txs`** (disjoint ≤20-input sweeps).

### 3. zetaclient — minimal poller + signer

New code is small and self-contained:

1. Poll the drain API URL on an interval; verify the payload signature against the **baked-in** public key. Reject on failure. Only a payload with `final: true` is eligible to trigger signing.
2. Watch the zeta head via the existing `GetBlockHeight` query. Fire only when `H ≤ current < H + window` — an **"ignored if missed"** guard mirroring `OperationalFlags.RestartHeight` (`maintenance/shutdown_listener.go`), so a late-starting node doesn't drain into a changed world.
3. For each tx: **assert `to == hardcoded safe address`** (abort otherwise), then call the existing `Signer.Sign(nil, to, amount, gas, nonce)` with the payload values (feed `H` as the keysign height), and broadcast via the EVM client's `SendTransaction`.

This reuses the existing native-transfer signer; the genuinely new EVM signing code is close to zero.

### 4. Bitcoin specifics

BTC needs more than EVM because it's UTXO-based, and the determinism-critical parts must move into the (API-produced) payload:

- **Pin the UTXO set.** `SignWithdrawTx`'s internal `SelectUTXOs` runs against each node's own BTC RPC view and is non-deterministic across nodes. So the **API selects the UTXOs once** and pins the exact outpoints + amounts in `btc_txs[].inputs`. Clients never re-select — they build the tx from the pinned inputs, so every node's per-input sighashes are identical.
- **Don't reuse `SignWithdrawTx`.** Build the exact `wire.MsgTx` (pinned inputs → single output to the receiver, minus fee; **no change output, no nonce-mark** — this is a sweep, not a tracked withdrawal) and sign with the lower-level `SignTx` (`zetaclient/chains/bitcoin/signer/sign.go:227`), which computes witness sighashes and calls `TSS().SignBatch`. That's the real TSS-signing primitive, fed fully-deterministic input.
- **20-input cap → multiple sweeps.** `MaxNoOfInputsPerTx = 20`. The API partitions the UTXO set into disjoint groups of ≤20 and emits one **independent** sweep tx per group (they spend disjoint inputs, so no chaining and they can be signed/broadcast in parallel). Each tx forms its own go-tss ceremony (distinct sighashes → distinct `MsgID`).
- **Fee/amount.** `fee = feeRate × txSize` per tx (reuse `common.EstimateOutboundSize`); amount logic reuses the existing BTC migration calc. All pinned in the payload.
- **Receiver** is a hardcoded per-network **BTC** address (separate from the EVM one); client asserts the payload `to` matches.

The draft → freeze → final model applies identically: the UTXO set is stable under quiescence, and the `final` pins it.

## Security model

Leaving zetacore means giving up on-chain consensus over the trigger. The following restore most of that trust:

- **Hardcoded receiver is the anchor.** The safe address is compiled into the client, and the client asserts the payload's `to` matches it. A compromised/spoofed/MITM'd API can at worst change *when* the drain fires — it can **never redirect funds**. Money can only ever go to the compile-time address.
- **Signed payload (off-chain verification).** A dedicated signing keypair is generated; the **public key is compiled into the client** as a constant (reviewed in the PR, shipped in the release) — the same trust model as the hardcoded receiver, no blockchain lookup involved. The cron holds the private key and signs each `final` payload; the client verifies locally with `verify(baked_in_pubkey, bytes, sig)` (secp256k1, reusing existing crypto). This stops a compromised host from firing early or injecting griefing gas/amounts. It does **not** affect fund destination — that's already fixed by the hardcoded receiver.
- **Per-operator arming (optional).** A local `--drain-armed` flag so no external party can trigger without each operator opting in — defense in depth.
- **Quiescence assumption.** Nonce determinism relies on inbound being disabled and no pending outbounds. Disable inbound (`MsgDisableCCTX`, emergency policy) and drain pending nonces before arming, exactly as the existing TSS-migration runbook does.

## Teardown (one-time capability)

This leaves an external URL that can trigger a TSS transfer. It must not live permanently in production:

- Gate the drain poller behind a build tag / dedicated release that operators run **only** for the drain window.
- After the drain completes on all chains, upgrade operators off that build and take the endpoint down.

## Testing plan

- **e2e (localnet).** Run the drain API **locally** and point the two localnet zetaclients at it. Model on `e2e/e2etests/test_migrate_tss.go` (which already covers both EVM and BTC migration): disable inbound, generate receiver addresses (EVM + BTC), publish the signed payload, wait for the txs to mine, assert TSS balances → ~0 and the receivers increased. For BTC, fetch UTXOs via the runner's BTC RPC (not mempool.space, which lacks regtest) and follow the existing 20-UTXO / multi-round pattern. This exercises the **real** 2-node ceremony and identical-digest/sighash coordination — not a mock signer. **No new TSS keygen and no `MsgUpdateTssAddress`** — drain only.
- **testnet (Athens).** Host the real endpoint, drain testnet gas to a throwaway wallet. Validates real operator clock skew, real gas prices, and the signature-verify path end to end.
- **mainnet.** Operators run the drain build, disable inbound and drain pending nonces, operator publishes the signed payload with the agreed height, fire once per chain, then upgrade away and take the endpoint down.

## Resolved decisions

- **Distribution:** signed JSON published to an **object store (S3/GCS) by a cron job** that reruns the tool each interval (drafts), then a single `final`. No always-on service. Clients poll the URL. Host is trusted only for availability/timing, never destination.
- **Payload integrity:** **signed**, verified off-chain against a **baked-in secp256k1 pubkey** (see Security).
- **Receiver:** **hardcoded per network** (localnet/testnet/mainnet constant map); client selects by its own network and asserts the payload `to` matches.
- **Gas price:** **zetacore median** gas price (matches on-chain migration), queried API-side.
- **Amount freshness:** draft → freeze → single `final`; buffer absorbs post-freeze drift.

## Remaining tuning parameters (proposed defaults — confirm during implementation)

- **Freeze window `K`** — blocks between the `final` publish and `H`. Proposed: large enough to cover a few client poll intervals (e.g. tens of zeta blocks). Must be ≫ poll interval.
- **Client poll interval** — proposed on the order of seconds–low tens of seconds.
- **Trigger `window`** — blocks after `H` during which a node may still fire before giving up ("ignored if missed"). Proposed: small (a handful of blocks).
- **Signing key custody** — who holds the cron's private key and how the matching pubkey lands in the release build (operator process, not a code decision).

## Files touched (summary)

- `pkg/gas/` (or `x/crosschain/types`) — extract `ComputeEVMMigration` and a BTC amount helper (shared with `MigrateFundCmdCCTX`); refactor the on-chain builder to call them. (No behavior change on the existing path.)
- `cmd/zetatool/cli/` — new subcommand emitting the signed, fully-resolved JSON payload for EVM + BTC (reuses `tss_balances.go` plumbing; BTC UTXO selection/partitioning pinned here).
- `zetaclient/` — new drain poller + EVM signer (reuse `Signer.Sign`) + BTC sweep signer (build `wire.MsgTx` from pinned inputs, reuse `SignTx` + `Broadcast`); per-network EVM & BTC receiver assertion; behind a build tag / opt-in flag.
- `e2e/e2etests/` — new drain test (EVM + BTC) with a local API.
