# Overview

The `observer` module keeps track of ballots for voting, a mapping between
chains and observer accounts, a list of supported connected chains, core
parameters (contract addresses, outbound transaction schedule interval, etc.),
observer parameters (ballot threshold, min observer delegation, etc.), and admin
policy parameters.

Ballots are used to vote on inbound and outbound transaction. The `observer`
module keeps create, read, update, and delete (CRUD) operations for ballots, as
well as helper functions to determine if a ballot has been finalized. The ballot
system is used by other modules, such as the `crosschain` module when observer
validators vote on transactions.

An observer validator is a validator that runs `zetaclient` alongside the
`zetacored` (the blockchain node) and is authorized to vote on inbound and
outbound cross-chain transactions.

A mapping between chains and observer accounts right now is set during genesis
and is used in the `crosschain` module to determine whether an observer
validator is authorized to vote on a transaction coming in/out of a specific
connected chain.
