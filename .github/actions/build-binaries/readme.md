# Build ZetaChain Binaries 

## Example 
```
      - name: Build Binaries
        uses: ./.github/actions/build-binaries
        with:
          run-tests: ${{ env.GITHUB_REF_NAME != 'develop' }}
          build-indexer: false
          go-version: '1.20'
```