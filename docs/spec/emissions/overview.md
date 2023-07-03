# Overview

The `emissions` module is responsible for orchestrating rewards distribution for
observers, validators and TSS signers. Currently, it only distributes rewards to
validators every block. The undistributed amount for TSS and observers is stored
in their respective pools.

The distribution of rewards is implemented in the begin blocker.

The module keeps track of parameters used for calculating rewards:

- Maximum bond factor
- Minimum bond factor
- Average block time
- Target bond ratio
- Validator emission percentage
- Observer emission percentage
- TSS Signer emission percentage
- Duration factor constant
