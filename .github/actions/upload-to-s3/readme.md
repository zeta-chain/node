# Uploade Files To S3


## Example 

```
      - name: upload-files-to-s3
        uses: ./.github/actions/upload-to-s3
        with:
          bucket-name: ${{ env.S3_BUCKET_PATH }}
          release-name: "$(echo ${{ github.ref_name }} | tr '//' '-')"
          git-sha: ${{ github.sha }}
          files: |
            zetacored
            zetaclientd
            cosmovisor
```