# **ZetaChain Git Workflow & Release Lifecycle**

## Repository

https://github.com/zeta-chain/node

## Overview

ZetaChain uses a dual-branch workflow inspired by GitFlow but adapted for our unique blockchain context. This workflow supports independent release cycles for ZetaCore (consensus-breaking changes) and ZetaClient (non-consensus-breaking features), enabling frequent ZetaClient updates without requiring network-wide consensus upgrades.

Unlike traditional GitFlow where all features eventually merge to a single integration branch, our workflow uses `develop` as the primary integration branch for all changes, while `main` represents the production-ready state for consensus-breaking releases.
### Core Branches

**`develop`** - Primary integration branch

- All changes merge here: both consensus-breaking and non-consensus-breaking
- Base branch for ZetaClient releases
- Consensus-breaking changes can safely be merged here since they only affect the network when a ZetaCore release is deployed
- Examples: Protocol upgrades, consensus rule changes, ZetaClient features, client-side optimizations

**`main`** - Production-ready consensus state

- Only updated when preparing a ZetaCore release
- Represents the current consensus-breaking version running on the network
- Created by merging `develop` → `main` when ready for a consensus upgrade

### Release Branches

**`release/zetacore/vX`** - Created from `main`

- Major version releases (consensus-breaking)
- Follows semantic versioning: v36.0.0, v37.0.0, etc.

**`release/zetaclient/vY`** - Created from `develop` or `main`

- Independent version cycle: v1.0.0, v1.1.0, v2.0.0, etc.
- Source branch depends on release type:
    - From `develop`: Independent ZetaClient release
    - From `main`: Coordinated release with ZetaCore

## Protocol Contracts Integration

### Repositories

- https://github.com/zeta-chain/protocol-contracts-evm
- https://github.com/zeta-chain/protocol-contracts-solana
- https://github.com/zeta-chain/protocol-contracts-ton
- https://github.com/zeta-chain/protocol-contracts-sui

### Overview

Protocol contracts repositories are related to ZetaClient versions, as ZetaClient interacts with deployed smart contracts across different chains. Each protocol contracts repository maintains its own independent versioning.

### Versioning Relationship

**Independent Versioning**

- Protocol contracts follow their own semantic versioning independent from ZetaCore and ZetaClient
- Each repository (EVM, Solana, TON, Sui) has its own version lifecycle
- Contract changes may correspond with ZetaClient updates to handle new interfaces

## Versioning Strategy

### Independent Semantic Versioning

**ZetaCore**: `v36.0.0`

- Major version bumps for consensus-breaking changes
- Requires network-wide coordination and governance

**ZetaClient**: `v2.5.1`

- Follows independent release cycle
- Can iterate rapidly without network upgrades

**Protocol Contracts**: `v1.2.3`

- Each protocol contracts repository (EVM, Solana, TON, Sui) follows its own independent semantic versioning


## Release Process

### ZetaClient-Only Release

1. Create `release/zetaclient/vY` from `develop`
2. Run [ZetaClient release Github action](https://github.com/zeta-chain/node/actions/workflows/publish-release-zetaclient.yml)
3. No ZetaCore release required
4. Coordinate deployment of the release

### Consensus-Breaking Release

1. Merge `develop` → `main` (bringing all accumulated changes including consensus-breaking ones)
2. Create `release/zetacore/vX` from `main`
3. Create `release/zetaclient/vY` from `main` (coordinated release)
4. Run [ZetaClient release Github action](https://github.com/zeta-chain/node/actions/workflows/publish-release-zetaclient.yml)
5. Run [ZetaCore release Github action](https://github.com/zeta-chain/node/actions/workflows/publish-release-zetacore.yml)
6. Submit an upgrade governance proposal

### Protocol Contracts Release Process

When a new protocol contracts release is required:

1. **Branch Creation**: Create `release/vX` in the protocol contracts repository
    - Uses the protocol contracts' own versioning scheme
    - Example: `release/v1` for protocol contracts v1.0.0
    - Any subsequent patches or minor updates use the same branch (e.g., v1.0.1, v1.1.0)
2. **Release**: Create a new release for the protocol contracts
    - Protocol contracts follow their own semantic versioning (e.g., v1.2.3, v1.3.0)
    - Version is independent from ZetaCore and ZetaClient versions

### Compatibility Tracking

The `VERSIONS.md` file in the node repository maintains the compatibility matrix between:
- ZetaCore versions
- ZetaClient versions
- Protocol contract versions (EVM, Solana, TON, Sui)

This ensures clear visibility of which component versions work together.

### Key Principle

- **Every ZetaCore release includes a ZetaClient release**
- **ZetaClient releases can happen independently**
- **Consensus-breaking changes in `develop` don't affect the network until a ZetaCore release is deployed**

## Notes on Consensus-Breaking Changes

A **consensus-breaking change** is any modification that influences the on-chain state transition, potentially leading to inconsistency and consensus failure if nodes run different versions simultaneously. These changes require network-wide coordination through governance proposals, ensuring all validators upgrade to the new version at the same designated upgrade height.

### **Some Classification Guidelines**

#### **Always Consensus-Breaking: `/x/**` directory**

- For security and safety, we treat **all changes in the `/x/` directory as consensus-breaking**
- This conservative approach prevents accidental consensus failures
- Note: In practice, some `/x/` changes may not affect consensus

#### **Never Consensus-Breaking: ZetaClient & E2E**

- **`/zetaclient/**`**: Client-side logic that doesn't affect chain state
- **`/e2e/**`**: Testing infrastructure independent of consensus rules

#### **Case-by-Case: `/pkg/**` directory**

- **Requires verification**: Changes in `/pkg/` may or may not be consensus-breaking
- **Assessment needed**: Determine where the modified package is used
- **If used by `/x/` modules**: Likely consensus-breaking
- **If used only by ZetaClient**: Likely non-consensus-breaking