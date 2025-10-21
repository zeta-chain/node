# **ZetaChain Git Workflow & Release Lifecycle**

## Repository

https://github.com/zeta-chain/node

## Overview

ZetaChain uses a dual-branch workflow inspired by GitFlow but adapted for our unique blockchain context. This workflow supports independent release cycles for ZetaCore (consensus-breaking changes) and ZetaClient (non-consensus-breaking features), enabling frequent ZetaClient updates without requiring network-wide consensus upgrades.

Unlike traditional GitFlow where all features eventually merge to a single integration branch, our workflow maintains separate tracks for consensus-breaking changes that require network coordination versus ZetaClient features that can be deployed independently.

### Core Branches

**`main`** - Consensus-breaking changes only

- All changes that affect network consensus
- Base branch for ZetaCore releases
- Examples: Protocol upgrades, consensus rule changes, state machine modifications

**`develop`** - Non-consensus-breaking changes

- ZetaClient features and improvements
- Base branch for ZetaClient-only releases
- Examples: Client-side optimizations, new features for connected chains that doesn’t require core protocol changes

### Release Branches

**`release/zetacore/vX`** - Created from `main`

- Major version releases (consensus-breaking)
- Follows semantic versioning: v36.0.0, v37.0.0, etc.

**`release/zetaclient/vY`** - Created from `develop` or `main`

- Independent version cycle: v1.0.0, v1.1.0, v2.0.0, etc.
- Source branch depends on release type:
    - From `develop`: Independent ZetaClient release
    - From `main`: Coordinated release with ZetaCore

## Versioning Strategy

### Independent Semantic Versioning

**ZetaCore**: `v36.0.0`

- Major version bumps for consensus-breaking changes
- Requires network-wide coordination and governance

**ZetaClient**: `v2.5.1`

- Follows independent release cycle
- Can iterate rapidly without network upgrades

## Release Process

### ZetaClient-Only Release

1. Create `release/zetaclient/vY` from `develop`
2. Run [ZetaClient release Github action](https://github.com/zeta-chain/node/actions/workflows/publish-release-zetaclient.yml)
3. No ZetaCore release required
4. Coordinate deployment of the release

### Consensus-Breaking Release

1. Merge `develop` → `main`
2. Create `release/zetacore/vX` from `main`
3. Create `release/zetaclient/vY` from `main` (coordinated release)
4. Run [ZetaClient release Github action](https://github.com/zeta-chain/node/actions/workflows/publish-release-zetaclient.yml)
5. Run [ZetaCore release Github action](https://github.com/zeta-chain/node/actions/workflows/publish-release-zetacore.yml)
6. Submit an upgrade governance proposal

### Key Principle

- **Every ZetaCore release includes a ZetaClient release**
- **ZetaClient releases can happen independently**

# Protocol Contracts Integration

## Repositories

- https://github.com/zeta-chain/protocol-contracts-evm
- https://github.com/zeta-chain/protocol-contracts-solana
- https://github.com/zeta-chain/protocol-contracts-ton
- https://github.com/zeta-chain/protocol-contracts-sui

## Overview

Protocol contracts repositories are tightly coupled with ZetaClient versions, as ZetaClient interacts directly with deployed smart contracts across different chains.

## Versioning Relationship

**ZetaClient ↔ Protocol Contracts Coupling**

- Each ZetaClient release corresponds to a specific set of protocol contracts
- Protocol contract changes often require ZetaClient updates to handle new interfaces
- Contract deployments must be coordinated with ZetaClient releases

Compatibility matrix can be found in the `VERSIONS.md` file at the root of the repository.

## Release Workflow

### Protocol Contracts Release Process

When a new protocol contracts release is required:

1. **Branch Creation**: Create `release/zetaclient/vY` in the protocol contracts repository
    - Branch name matches the corresponding ZetaClient release version
    - Example: ZetaClient v2.5.0 → protocol contracts `release/zetaclient/v2.5`
2. **Release**: Create a new release for the protocol contracts
    - Protocol contracts follow their own semantic versioning (e.g., v1.2.3, v1.3.0)
    - Release version is independent from ZetaClient version numbering
    - Example: ZetaClient v2.5.0 might correspond to protocol contracts v1.3.0
3. **Deployment**: Deploy contracts from the release branch to respective networks

# Notes on Consensus Breaking Changes

A **consensus-breaking change** is any modification that influences the on-chain state transition, potentially leading to inconsistency and consensus failure if nodes run different versions simultaneously. These changes require network-wide coordination through governance proposals, ensuring all validators upgrade to the new version at the same designated upgrade height.

## **Some Classification Guidelines**

### **Always Consensus-Breaking: `/x/**` directory**

- For security and safety, we treat **all changes in the `/x/` directory as consensus-breaking**
- This conservative approach prevents accidental consensus failures
- Note: In practice, some `/x/` changes may not affect consensus

### **Never Consensus-Breaking: ZetaClient & E2E**

- **`/zetaclient/**`**: Client-side logic that doesn't affect chain state
- **`/e2e/**`**: Testing infrastructure independent of consensus rules

### **Case-by-Case: `/pkg/**` directory**

- **Requires verification**: Changes in `/pkg/` may or may not be consensus-breaking
- **Assessment needed**: Determine where the modified package is used
- **If used by `/x/` modules**: Likely consensus-breaking
- **If used only by ZetaClient**: Likely non-consensus-breaking