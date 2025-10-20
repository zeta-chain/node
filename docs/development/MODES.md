# ZetaClient modes of execution
ZetaClient can execute in one of three modes:
    `standard`, `dry`, or `chaos`.

## Standard
Represents the **standard** mode of execution for the ZetaClient.
A standard observer-signer
    observes transactions from ZetaChain,
    signs them,
    and relays them to the appropriate connected chains.
Symmetrically,
    it observes transactions from the connected chains
    and relays them to ZetaChain.

## Dry
Represents the **read-only** execution mode for the ZetaClient.
A dry observer-signer observes the transactions from ZetaChain
    and the connected chains
    without signing them,
    or otherwise mutating the state of the ZetaChain
    or the state of the connected chains.

### Implementation details
We implemented dry-mode by adding flag checks in specific parts of the code.
For example,
    the `ZetaRepo` checks if the client is in dry-mode
    before voting on an inbound.
Similarly,
    the signers of the connected chains check the client-mode flag
    before signing and broadcasting transactions.

In order to ensure that we never call mutating methods while in dry mode,
    we also wrap the client interfaces
    of ZetaCore, TSS, and of the connected chains
    with "dry" structures.
These dry wrappers override the mutating functions of their interfaces
    with methods that panic when called.
Naturally,
    the dry-wrappers are redundant
    and ZetaClient should never panic in such manner;
    their addition is part of a defense-in-depth strategy.

## Chaos
Represents the **chaos-testing** execution mode for the ZetaClient.
A observer-signer in chaos mode works as if in standard mode,
    but function calls that interact with outside resources
    (e.g. ZetaChain, connected chains, TSS, and other nodes)
    may intentionally fail.
We use chaos mode to replicate unstable environments for testing.

