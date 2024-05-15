import requests
import time
import logging
import os
from github import Github
from datetime import datetime, timezone, timedelta
import sys

# Setup logger
logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)
ch = logging.StreamHandler()
ch.setLevel(logging.DEBUG)
formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
ch.setFormatter(formatter)
logger.addHandler(ch)

STATUS_ENDPOINT = os.environ["STATUS_ENDPOINT"]
PROPOSALS_ENDPOINT = os.environ["PROPOSAL_ENDPOINT"]


def get_current_block_height():
    response = requests.get(STATUS_ENDPOINT)
    response.raise_for_status()
    data = response.json()
    return int(data['result']['sync_info']['latest_block_height'])


def get_proposals():
    response = requests.get(PROPOSALS_ENDPOINT)
    response.raise_for_status()
    data = response.json()
    return [proposal['id'] for proposal in data['proposals']]


def get_proposal_details(proposal_id):
    url = f"{PROPOSALS_ENDPOINT}/{proposal_id}"
    response = requests.get(url)
    response.raise_for_status()
    data = response.json()
    return data


def is_older_than_one_week(submit_time_str):
    submit_time_str = submit_time_str.split('.')[0] + 'Z'
    submit_time = datetime.strptime(submit_time_str, '%Y-%m-%dT%H:%M:%S%z')
    current_time = datetime.now(timezone.utc)
    return current_time - submit_time > timedelta(weeks=1)


def monitor_block_height(proposal_height):
    max_checks = int(os.environ["MAX_WAIT_FOR_PROCESSING_BLOCKS_CHECK"])
    checks = 0
    while checks < max_checks:
        current_height = get_current_block_height()
        if current_height >= proposal_height:
            for _ in range(50):
                prev_height = current_height
                time.sleep(3)
                current_height = get_current_block_height()
                if current_height > prev_height:
                    logger.info("Block height is moving. Network is processing blocks.")
                    return True
            logger.warning("Network is not processing blocks.")
            return False
        checks += 1
        time.sleep(3)
    logger.warning("Max wait time reached. Proposal height not reached.")
    return False


def update_github_release(proposal_title):
    github_token = os.environ["GITHUB_TOKEN"]
    g = Github(github_token)
    repo = g.get_repo("zeta-chain/node")
    releases = repo.get_releases()
    for release in releases:
        if release.title == proposal_title and release.prerelease:
            release.update_release(release.title, release.body, draft=False, prerelease=False)
            logger.info(f"Updated GitHub release '{proposal_title}' from pre-release to release.")
            return
    logger.warning(f"No matching GitHub pre-release found for title '{proposal_title}'.")


def main():
    current_block_height = get_current_block_height()
    logger.info(f"Current Block Height: {current_block_height}")
    proposals_retrieved = get_proposals()
    proposals = {}

    for proposal_id in proposals_retrieved:
        proposal_details = get_proposal_details(proposal_id)

        submit_time_str = proposal_details["proposal"]["submit_time"]
        if is_older_than_one_week(submit_time_str):
            logger.info(f"Proposal {proposal_id} is older than one week. Skipping.")
            continue

        for message in proposal_details["proposal"]["messages"]:
            if "content" not in message:
                continue
            proposal_type = message["content"]["@type"]
            logger.info(f"id: {proposal_id}, proposal type: {proposal_type}")
            if "plan" not in message["content"]:
                continue
            if proposal_details["proposal"]["status"] != "PROPOSAL_STATUS_PASSED":
                logger.info(f'Proposal did not pass: {proposal_details["proposal"]["status"]}')
                continue
            if 'SoftwareUpgradeProposal' in str(proposal_type) or 'MsgSoftwareUpgrade' in str(proposal_type):
                proposals[proposal_id] = {
                    "proposal_height": int(message["content"]["plan"]["height"]),
                    "proposal_title": message["content"]["title"]
                }
            break
    if len(proposals) <= 0:
        logger.info("No proposals found within the timeframe.")
        sys.exit(0)
    for proposal_id, proposal_data in proposals.items():
        if current_block_height >= proposal_data["proposal_height"]:
            logger.info(f"Proposal {proposal_id} height {proposal_data['proposal_height']} has been reached.")
            update_github_release(proposal_data["proposal_title"])
        else:
            logger.info(
                f"Waiting for proposal {proposal_id} height {proposal_data['proposal_height']} to be reached. Current height: {current_block_height}")
            if monitor_block_height(proposal_data['proposal_height']):
                logger.info(
                    f"Proposal {proposal_id} height {proposal_data['proposal_height']} has been reached and network is processing blocks.")
                update_github_release(proposal_data["proposal_title"])
            else:
                logger.warning(
                    f"Failed to reach proposal {proposal_id} height {proposal_data['proposal_height']} or network is not processing blocks.")


if __name__ == "__main__":
    main()