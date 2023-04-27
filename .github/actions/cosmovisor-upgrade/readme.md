# Cosmosvisor Upgrade Proposal

Creates a cosmovisor upgrade proposal and then issues a vote on all validators

## Example 
      - name: Cosmovisor Upgrade Proposal and Vote
        uses:  ./.github/actions/cosmovisor-upgrade
        with:
          UPGRADE_NAME: ${{ github.event.inputs.UPGRADE_NAME }}
          DESCRIPTION: "Upgrade Description Goes Here"
          CHAIN_ID: ${{ env.CHAIN_ID }}
          ZETACORED_CHECKSUM: "1234567" #SHA256 
          ZETACORED_URL: "https://${{ env.S3_BUCKET_NAME }}.s3.amazonaws.com/builds/zeta-node/develop/zetacored"
          ZETACLIENTD_CHECKSUM: "1234567" #SHA256 
          ZETACLIENTD_URL: "https://${{ env.S3_BUCKET_NAME }}.s3.amazonaws.com/builds/zeta-node/develop/zetaclientd"
          CHAIN_ID: ${{ github.event.inputs.CHAIN_ID }}
          API_ENDPOINT: "https://api.${{ github.event.inputs.ENVIRONMENT }}.zetachain.com"
          UPGRADE_BLOCK_HEIGHT: 999999
## Functions File
These are bash functions used by the action. 

## Prerequestices

### AWS Authentication 
You must authenticate to AWS before calling this action

