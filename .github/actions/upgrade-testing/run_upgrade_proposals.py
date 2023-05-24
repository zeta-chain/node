import os
import json
import time
from libraries.zetaops import GithubBinaryDownload
from libraries.zetaops import Utilities
from libraries.zetaops import Logger
import sys
import re

logger = Logger()
logger.log.info("**************************Initiate GitHub Binary Downloader**************************")
binary_downloader = GithubBinaryDownload(os.environ["GITHUB_TOKEN"], os.environ["GITHUB_OWNER"], os.environ["GITHUB_REPO"])

logger.log.info("Initiate Utilities")
command_runner = Utilities(os.environ["GO_PATH"])
command_runner.logger = logger.log
command_runner.NODE = os.environ["NODE"]
command_runner.MONIKER = os.environ["MONIKER"]
command_runner.CHAIN_ID = os.environ["CHAIN_ID"]

git_tags = command_runner.run_command("git tag --list > git_tags && cat git_tags && rm -rf git_tags").split("\n")
p = re.compile(r'[a-z][0-9].[0-9].[0-9]*')
tag_list = []
for tag in git_tags:
    if p.match(tag):
        logger.log.info(tag)
        tag_list.append(tag)

if len(tag_list) == 0 or len(tag_list) == 1:
    sys.exit(0)

tag_list.sort()

upgrades_json = open("upgrades.json", "r").read()
upgrades_json = json.loads(upgrades_json)
binary_download_list = []
first = True
non_concensus_upgrades = []
# for tag in tag_list:
#     if first:
#         first_major_version = tag.split(".")[0]
#         first_minor_version = tag.split(".")[1]
#         first = False
#     else:
#         major_version = tag.split(".")[0]
#         minor_version = tag.split(".")[1]
#         #Essentially check the last known major and minor version to determine if it was a concensus breaking version change.
#         if major_version == first_major_version and minor_version != first_minor_version:
#             non_concensus_upgrades.append(tag)
#         first_major_version = tag.split(".")[0]
#         first_minor_version = tag.split(".")[1]
#     binary_download_list.append([f"{tag}", f"zetacored-{os.environ['BINARY_NAME_SUFFIX']}"])

binary_download_list.append(["v1.1.6", "zetacored-ubuntu-22-amd64"],["v1.1.7", "zetacored-ubuntu-22-amd64"])

logger.log.info(binary_download_list)
os.environ["STARTING_VERSION"] = tag_list[0]
os.environ["END_VERSION"] = tag_list[len(tag_list)-1]
logger.log.info(f"Starting Version: {os.environ['STARTING_VERSION']}")
logger.log.info(f"End Version Version: {os.environ['END_VERSION']}")

tag_list.pop(0)

upgrades_json["upgrade_sleep_time"] = os.environ["UPGRADES_SLEEP_TIME"]
upgrades_json["binary_versions"] = binary_download_list
upgrades_json["upgrade_versions"] = tag_list
upgrades_json_write = open("upgrades.json", "w")
logger.log.info(upgrades_json)
upgrades_json_write.write(json.dumps(upgrades_json))
upgrades_json_write.close()

logger.log.info("**************************Generate Wallet For Test**************************")
command_runner.generate_wallet()
command_runner.load_key()
logger.log.info("")

logger.log.info("**************************Download Github Binaries from upgrades.json**************************")
binary_downloader.download_testing_binaries()
logger.log.info("")

logger.log.info("**************************Build Docker Image**************************")
command_runner.build_docker_image(os.environ["DOCKER_FILE_LOCATION"])
logger.log.info("")

logger.log.info("**************************Start Docker Container and Sleep for 60 Seconds**************************")
command_runner.start_docker_container(os.environ["GAS_PRICES"],
                               os.environ["DAEMON_HOME"],
                               os.environ["DAEMON_NAME"],
                               os.environ["DENOM"],
                               os.environ["DAEMON_ALLOW_DOWNLOAD_BINARIES"],
                               os.environ["DAEMON_RESTART_AFTER_UPGRADE"],
                               os.environ["EXTERNAL_IP"],
                               os.environ["STARTING_VERSION"],
                               os.environ["PROPOSAL_TIME_SECONDS"],
                               os.environ["LOG_LEVEL"],
                               os.environ["UNSAFE_SKIP_BACKUP"],
                               os.environ["CLIENT_DAEMON_NAME"],
                               os.environ["CLIENT_DAEMON_ARGS"],
                               os.environ["CLIENT_SKIP_UPGRADE"],
                               os.environ["CLIENT_START_PROCESS"])

time.sleep(10)
logger.log.info("**************************DOCKER PS**************************")
command_runner.docker_ps()
logger.log.info("")

logger.log.info("**************************CHECK CONTAINER ID**************************")
if not command_runner.CONTAINER_ID:
    logger.log.error(f"Container didn't start. No Container ID: {command_runner.CONTAINER_ID}")
    sys.exit(1)
logger.log.info("")

logger.log.info("**************************CHECK DEBUG IF SET SHOW LOGS**************************")
if "DEBUG_UPGRADES" in os.environ:
    if os.environ["DEBUG_UPGRADES"] != "false":
        command_runner.get_docker_container_logs()
logger.log.info("")

logger.log.info("**************************start upgrade process, open upgrades.json and read what upgrades to start.**************************")
UPGRADE_DATA = json.loads(open("upgrades.json", "r").read())
logger.log.info("")

for version in UPGRADE_DATA["upgrade_versions"]:
    logger.log.info(f"**************************starting upgrade for version: {version}**************************")
    VERSION=version
    BLOCK_TIME_SECONDS = int(os.environ["BLOCK_TIME_SECONDS"])
    PROPOSAL_TIME_SECONDS = int(os.environ["PROPOSAL_TIME_SECONDS"])
    UPGRADE_INFO = '{}'

    if version not in non_concensus_upgrades:
        logger.log.info("**************************raise governance proposal**************************")
        GOVERNANCE_TX_HASH = command_runner.raise_governance_proposal(VERSION, BLOCK_TIME_SECONDS, PROPOSAL_TIME_SECONDS, UPGRADE_INFO)[0]

        logger.log.info("**************************sleep for 10 seconds to allow the proposal to show up on the network**************************")
        time.sleep(10)

        logger.log.info("**************************get proposal id**************************")
        PROPOSAL_ID = command_runner.get_proposal_id()
        logger.log.info(PROPOSAL_ID)

        logger.log.info(f"raise governance vote on proposal id: {PROPOSAL_ID}")
        vote_output = command_runner.raise_governance_vote(PROPOSAL_ID)
        logger.log.info(f"""**************************UPGRADE INFO**************************
            MONIKER: {command_runner.MONIKER}
            NODE: {command_runner.NODE}
            PROPOSAL_ID: {PROPOSAL_ID}
            VERSION: {VERSION}
            UPGRADE_HEIGHT: {command_runner.UPGRADE_HEIGHT}
            UPGRADE_INFO: {UPGRADE_INFO}
            CHAIN_ID: {command_runner.CHAIN_ID}
            LATEST_BLOCK: {command_runner.CURRENT_HEIGHT}
        **************************UPGRADE INFO**************************""")
        time.sleep(int(UPGRADE_DATA["upgrade_sleep_time"]))
        command_runner.docker_ps()
    else:
        print(version, "this version will be done as non-consensus breaking upgrade")

logger.log.info("Check docker process is still running for debug purposes.")
command_runner.docker_ps()
time.sleep(30)
if command_runner.version_check(os.environ["END_VERSION"]):
    logger.log.info("**************************Version is what was expected.**************************")
    current_block = command_runner.current_block()
    logger.log.info("**************************Check to see if chain is still processing blocks.**************************")
    time.sleep(10)
    end_block = command_runner.current_block()
    if abs(end_block - current_block) > 0:
        logger.log.info("**************************chain still processing blocks upgrade path looks good**************************")
        logger.log.info("**************************kill running docker containers and cleanup.**************************")
        command_runner.kill_docker_containers()
        sys.exit(0)
    else:
        logger.log.info("**************************Chain doesn't seem to be processign blocks upgrade path was a failure.**************************")
        logger.log.info("**************************kill running docker containers and cleanup.**************************")
        command_runner.get_docker_container_logs()
        command_runner.kill_docker_containers()
        sys.exit(1)
else:
    logger.log.info("**************************Version didn't match what was expected.**************************")
    logger.log.info("**************************kill running docker containers and cleanup.**************************")
    command_runner.get_docker_container_logs()
    command_runner.kill_docker_containers()
    sys.exit(1)

