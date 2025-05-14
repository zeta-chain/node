# Zetachain x Sui Gateway

This package is only used for gateway upgrade test.

The `sui client upgrade` command requires the presence of the whole `gateway` package for re-publishing,
so we have a minimized copy of sui gateway project in here to help the upgrading process.


The source code is copied from [protocol-contracts-sui](https://github.com/zeta-chain/protocol-contracts-sui) with one single additional function added:

```
// upgraded returns true to indicate gateway has been upgraded
entry fun upgraded(): bool {
    true
}
```