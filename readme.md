# ⚠️ Important Notice: Code Freeze for Security Audit ⚠️

Dear Contributors and Users,
We are committed to ensuring the highest standards of security and reliability in our project. To uphold this commitment, we are currently undergoing a thorough security audit conducted by Code4rena. More info: https://code4rena.com/contests/2023-11-zetachain.

During this period, we have instituted a code freeze on our public repository. This means there will be no new commits, merges, or major changes to the codebase until the audit is complete. This process is crucial to maintain the integrity and consistency of the code being audited.

The audit is scheduled from 20 Nov 9:00 PM GMT+1 to 18 Dec 9:00 PM GMT+1. We appreciate your patience and understanding during this vital phase of our project's development.

During this time, we encourage our community to review the current codebase and documentation. While we won't be merging new changes, we welcome your feedback, which we will address post-audit.

For any questions or concerns, feel free to reach out to us on our [Discord](https://discord.com/invite/zetachain).

Thank you for your cooperation and support!

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

## Creating a Release for Mainnet
Creating a release for mainnet is a straightforward process. Here are the steps to follow:

### Steps
 - Step 1. Open a Pull Request (PR): Begin by opening a PR from the release candidate branch (e.g., vx.x.x-rc) to the main branch.
 - Step 2. Testing and Validation: Allow the automated tests, including smoke tests, linting, and upgrade path testing, to run. Ensure that these tests pass successfully.
 - Step 3. Approval Process: Obtain the necessary approvals from relevant stakeholders or team members.
 - Step 4. Merging PR: Once all requirements have been met and the PR has received the required approvals, merge the PR. The automation will then be triggered to proceed with the release.

By following these steps, you can efficiently create a release for Mainnet, ensuring that the code has been thoroughly tested and validated before deployment.

## Creating a Release for Testnet
Creating a release for testnet is a straightforward process. Here are the steps to follow:

### Steps
 - Step 1. Create the release candidate tag with the following format (e.g., vx.x.x-rc) ex. v11.0.0-rc.
 - Step 2. Once a RC branch is created the automation will kickoff to build and upload the release and its binaries.

By following these steps, you can efficiently create a release candidate for testnet for QA and validation. In the future we will make this automatically deploy to testnet when a -rc branch is created. 
Currently, raising the proposal to deploy to testnet is a manual process via GitHub Action pipeline located in the infrastructure repo. 


## Creating a Hotfix Release
Creating a hotfix release is a straightforward process. Here are the steps to follow:

### Steps
 - Step 1. Execute pipeline: https://github.com/zeta-chain/node/actions/workflows/publish-release.yml 
 - Step 2. select branch when running pipeline manually your hotfix lives on.
 - Step 3. specify the version in the input field and run workflow. ex. vx.x.x-hotfix recommended.

Wheny ou execute with hotfix it will build and publish the binaries to the releases. 