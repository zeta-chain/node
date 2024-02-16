# ZetaChain

ZetaChain is an EVM-compatible L1 blockchain that enables omnichain, generic
smart contracts and messaging between any blockchain.

## Prerequisites

- [Go](https://golang.org/doc/install) 1.20
- [Docker](https://docs.docker.com/install/) and
  [Docker Compose](https://docs.docker.com/compose/install/) (optional, for
  running tests locally)
- [buf](https://buf.build/) (optional, for processing protocol buffer files)
- [jq](https://stedolan.github.io/jq/download/) (optional, for running scripts)

## Components of ZetaChain

ZetaChain is built with [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), a
modular framework for building blockchain and
[Ethermint](https://github.com/evmos/ethermint), a module that implements
EVM-compatibility.

- [zeta-node](https://github.com/zeta-chain/zeta-node) (this repository)
  contains the source code for the ZetaChain node (`zetacored`) and the
  ZetaChain client (`zetaclientd`).
- [protocol-contracts](https://github.com/zeta-chain/protocol-contracts)
  contains the source code for the Solidity smart contracts that implement the
  core functionality of ZetaChain.

## Building the zetacored/zetaclientd binaries
For the Athens 3 testnet, clone this repository, checkout the latest release tag, and type the following command to build the binaries:
```
make install-testnet
```
to build. 

This command will install the `zetacored` and `zetaclientd` binaries in your
`$GOPATH/bin` directory.

Verify that the version of the binaries match the release tag.  
```
zetacored version
zetaclientd version
```

## Making changes to the source code

After making changes to any of the protocol buffer files, run the following
command to generate the Go files:

```
make proto
```

This command will use `buf` to generate the Go files from the protocol buffer
files and move them to the correct directories inside `x/`. It will also
generate an OpenAPI spec.

### Generate documentation

To generate the documentation, run the following command:

```
make specs
```

This command will run a script to update the modules' documentation. The script
uses static code analysis to read the protocol buffer files and identify all
Cosmos SDK messages. It then searches the source code for the corresponding
message handler functions and retrieves the documentation for those functions.
Finally, it creates a `messages.md` file for each module, which contains the
documentation for all the messages in that module.

## Running tests

To check that the source code is working as expected, refer to the manual on how
to [run the smoke test](./contrib/localnet/README.md).

## Community

[Twitter](https://twitter.com/zetablockchain) |
[Discord](https://discord.com/invite/zetachain) |
[Telegram](https://t.me/zetachainofficial) | [Website](https://zetachain.com)


## Creating a Release Candidate
Creating a release candidate for testing is a straightforward process. Here are the steps to follow:

### Steps
 - Step 1. Create the release candidate tag with the following format (e.g., vx.x.x-rc) ex. v11.0.0-rc.
 - Step 2. Once a RC branch is created the automation will kickoff to build and upload the release and its binaries.

By following these steps, you can efficiently create a release candidate for QA and validation. In the future we will make this automatically deploy to a testnet when a -rc branch is created. 
Currently, raising the proposal to deploy to testnet is a manual process via GovOps repo. 

## Creating a Release / Hotfix Release

To create a release simply execute the publish-release workflow and follow the steps below.

### Steps
 - Step 1. Go to this pipeline: https://github.com/zeta-chain/node/actions/workflows/publish-release.yml
 - Step 2. Select the dropdown branch / tag you want to create the release with.
 - Step 3. In the version input, include the version of your release. Note. The major version must match what is in the upgrade handler.
 - Step 4. Select if you want to skip the tests by checking the checkbox for skip tests.
 - Step 5. Once the testing steps pass it will create a Github Issue. This Github Issue needes to be approved by one of the approvers: kingpinXD,lumtis,brewmaster012

Once the release is approved the pipeline will continue and will publish the releases with the title / version you specified in the user input.