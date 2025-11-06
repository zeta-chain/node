## Building Test Release Binaries

Build and upload test release binaries to Cloudflare R2 using the GitHub Actions workflow. **Note:** Only contributors with write access can trigger this workflow.

### How to Build

1. Navigate to the [Actions](https://github.com/zeta-chain/node/actions/workflows/release-test-build.yml) tab
2. Select "Release Test Binaries to Cloudflare R2" workflow and click "Run workflow"
3. Optionally provide:
   - **Version**: Custom release version (e.g., `v1.0.0-test`). If not provided, auto-generates as `{branch-name}-{date}-{short-hash}`
   - **Ref**: Branch, tag, or commit SHA to build from. If not provided, uses the default branch

The workflow builds both `zetacored` and `zetaclientd` binaries for all supported platforms and uploads them to Cloudflare R2. Completion notification with release ID and download location is provided in the workflow summary and Slack.

### Example Release IDs

- `develop-20241215-abc1234` (auto-generated)
- `v1.0.0-test` (custom version)
