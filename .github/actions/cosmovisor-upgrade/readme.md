# Cosmosvisor Upgrade Proposal

Creates a cosmovisor upgrade proposal and then issues a vote on all validators

## Example 
      - name: Cosmovisor Upgrade Proposal and Vote
        uses:  ./.github/actions/cosmovisor-upgrade
        with:
          RELEASE_NAME: "$(echo ${{ github.ref_name }} | tr '//' '-')"
          VERSION: "0.0.1"
          DESCRIPTION: "Upgrade Description Goes Here"
          ZETACORED_CHECKSUM: ""
          ZETACORED_URL: "https://${{ env.S3_BUCKET_NAME }}.s3.amazonaws.com/develop/zetacored"
          ZETACLIENTD_CHECKSUM: ""
          ZETACLIENTD_URL: "https://${{ env.S3_BUCKET_NAME }}.s3.amazonaws.com/develop/zetaclientd"
          CHAIN_ID: ${{ github.event.inputs.CHAIN_ID }}
          API_ENDPOINT: "https://api.${{ github.event.inputs.ENVIRONMENT }}.zetachain.com"
          UPGRADE_BLOCK_HEIGHT: 999999

## Functions
These are bash functions used by the action. 

## Prerequestices

### AWS Authentication 
You must authenticate to AWS before calling this action

