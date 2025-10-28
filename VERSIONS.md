# ZetaChain Versions Compatibility Matrix

> **Note:** This matrix tracks major versions only. For specific patch versions, refer to individual release notes.

> **Note:** ZetaClient v36 was released together with ZetaCore v36.
> Starting with v37, ZetaClient and ZetaCore use independent release cycles (ZetaClient v37 is available, ZetaCore v37 is not yet released).
> For details, see [Release Lifecycle Documentation](docs/development/RELEASE_LIFECYCLE.md).

## Current Versions

| Component      | Mainnet         | Testnet               | Development (not live) |
|----------------|-----------------|-----------------------|------------------------|
| ZetaCore       | [v36][v36]      | [v36][v36]            | [v36][v36]             |
| ZetaClient     | [v36][v36]      | [v37][zetaclient-v37] | [v37][zetaclient-v37]  |
| EVM Gateway    | [v14][evm-v14]  | [v15][evm-v15]        | [v15][evm-v15]         |
| Solana Gateway | [v5][solana-v5] | [v6][solana-v6]       | [v6][solana-v6]        |
| Sui Gateway    | [v1][sui-v1]    | [v2][sui-v2]          | [v2][sui-v2]           |
| TON Gateway    | [v2][ton-v2]    | [v2][ton-v2]          | [v2][ton-v2]           |

## Compatibility Table

### ZetaCore v36

| ZetaClient            | EVM Gateway    | Solana Gateway  | Sui Gateway  | TON Gateway  |
|-----------------------|----------------|-----------------|--------------|--------------|
| [v37][zetaclient-v37] | [v15][evm-v15] | [v6][solana-v6] | [v2][sui-v2] | [v2][ton-v2] | 
| [v36][v36]            | [v14][evm-v14] | [v5][solana-v5] | [v1][sui-v1] | [v2][ton-v2] | 

---

*Last updated: 2025-10-21*

[v36]: https://github.com/zeta-chain/node/releases/tag/v36.0.4
[zetaclient-v37]: https://github.com/zeta-chain/node/releases/tag/zetaclient_v37.0.0
[evm-v14]: https://github.com/zeta-chain/protocol-contracts-evm/releases/tag/v14.0.1
[evm-v15]: https://github.com/zeta-chain/protocol-contracts-evm/releases/tag/v15.0.0
[solana-v5]: https://github.com/zeta-chain/protocol-contracts-solana/releases/tag/v5.0.0
[solana-v6]: https://github.com/zeta-chain/protocol-contracts-solana/releases/tag/v6.0.0
[sui-v1]: https://github.com/zeta-chain/protocol-contracts-sui/releases/tag/v1.0.0
[sui-v2]: https://github.com/zeta-chain/protocol-contracts-sui/releases/tag/v2.0.0
[ton-v2]: https://github.com/zeta-chain/protocol-contracts-ton/releases/tag/v2.0.0