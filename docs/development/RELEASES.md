## Releases

NOTE: This release process will be deprecated as part of the new gitflow process. It will be updated as part of https://github.com/zeta-chain/node/issues/4263.

To start a new major release, begin by creating and pushing a `release/` branch based on `develop`.

<details>
<summary>Example Commands</summary>

```bash
git fetch
git checkout -b release/v15 origin/develop
git push origin release/v15
```

</details>

Most changes should first be merged into `develop` then backported to the release branch via a PR.

<details>
<summary>Example Commands to Backport a Change</summary>

```bash
git fetch
git checkout -b my-backport-branch origin/release/v15
git cherry-pick <commit SHA from develop>
git push origin my-backport-branch
```

</details>

### Creating a Release Candidate
You can use github actions to create a release candidate:
1) Create the release candidate tag with the following format (e.g., vx.x.x-rc) ex. v11.0.0-rc.
2) Push the tag and the automation will take care of the rest

You may create the RC tag directly off `develop` if a release branch has not been created yet. You should use the release branch if it exists and has diverged from develop.

By following these steps, you can efficiently create a release candidate for QA and validation. In the future we will make this automatically deploy to a testnet when a -rc branch is created.
Currently, raising the proposal to deploy to testnet is a manual process via GovOps repo.

### Creating a Release / Hotfix Release

To create a release simply execute the publish-release workflow and follow these steps:

1) Go to this pipeline: https://github.com/zeta-chain/node/actions/workflows/publish-release.yml
2) Select the release branch.
3) In the version input, include the version of your release.
- The major version must match what is in the upgrade handler
- The version should look like this: `v15.0.0`
4) Select if you want to skip the tests by checking the checkbox for skip tests.
5) Once the testing steps pass it will create a Github Issue. This Github Issue needes to be approved by one of the approvers: kingpinXD,lumtis,brewmaster012

Once the release is approved the pipeline will continue and will publish the releases with the title / version you specified in the user input.
