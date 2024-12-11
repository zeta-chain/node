# Zetaclient TSS Overview

(Threshold Signature Scheme)

This package wraps the go-tss library, providing a high-level API for signing arbitrary digests using TSS.
The underlying go-tss library relies on tss-lib.

## What is a Digest?

A digest is simply a byte slice (`[]byte`), typically representing a transaction hash or other cryptographic input.
The API allows secure signing of these digests in a distributed manner.

## Architecture Overview

This is the approximate structure of the TSS implementation within Zetaclient:

```text
zetaclientd(
    tss.Service(
        gotss.Server(libp2p()) 
    )
)
```

## Package Structure

- `setup.go`: Initializes the go-tss TSS server and the **Service** wrapper of this package.
- `keygen.go`: Manages the key generation ceremony, creating keys used by TSS.
- `service.go`: Implements the **Service** struct, offering methods for signing and verifying digests.
- Other Files: Utilities and supporting tools for TSS operations.

## Links

- `go-tss`: https://github.com/zeta-chain/go-tss
- `tss-lib`: https://github.com/zeta-chain/tss-lib
