import re
import requests
import os
import json
import logging
import sys
import shutil


# Logger class for easier logging setup
class Logger:
    def __init__(self):
        self.log = logging.getLogger()
        self.log.setLevel(logging.INFO)
        self.handler = logging.StreamHandler(sys.stdout)
        self.handler.setLevel(logging.DEBUG)
        self.formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
        self.handler.setFormatter(self.formatter)
        self.log.addHandler(self.handler)


# Initialize logger instance
logger = Logger()

# Define the path where upgrades will be stored, using an environment variable for the base path
upgrade_path = f'{os.environ["DAEMON_HOME"]}/cosmovisor/upgrades/'


# Function to find the latest patch version of a binary based on major and optional minor version
def find_latest_patch_version(major_version, minor_version=None):
    # Define a regex pattern to match version directories
    pattern = re.compile(
        f"v{major_version}\.{minor_version}\.(\d+)" if minor_version else f"v{major_version}\.0\.(\d+)")
    # List directories that match the version pattern
    versions = [folder for folder in os.listdir(upgrade_path) if pattern.match(folder)]
    if versions:
        try:
            # Find the maximum version, assuming it's the latest patch
            latest_patch_version = max(versions)
            # Return the path to the binary of the latest patch version
            return os.path.join(upgrade_path, latest_patch_version, "bin", "zetacored")
        except ValueError as e:
            logger.log.error(f"Error finding latest patch version: {e}")
            return None
    return None


# Function to replace an old binary with a new one
def replace_binary(source, target):
    try:
        # Log deletion of old binary
        if os.path.exists(target):
            logger.log.info(f"Deleted old binary: {target}")
            os.remove(target)
        # Copy the new binary to the target location
        shutil.copy(source, target)
        logger.log.info(f"Binary replaced: {target} -> {source}")
    except Exception as e:
        logger.log.error(f"Error replacing binary: {e}")


# Parse JSON from an environment variable to get binary download information
info = json.loads(os.environ["DOWNLOAD_BINARIES"])

try:
    # Iterate over binaries to download
    for binary in info["binaries"]:
        download_link = binary
        # Log download link
        logger.log.info(f"DOWNLOAD LINK: {download_link}")
        split_download_link = download_link.split("/")
        # Log split download link parts
        logger.log.info(f"SPLIT DOWNLOAD LINK: {split_download_link}")
        # Extract binary name and version from the download link
        binary_name = download_link.split("/")[8]
        version = download_link.split("/")[7]
        formatted_version = re.search(r'v\d{1,2}\.\d{1,2}\.\d{1,2}', version).group()
        end_binary_name = os.environ["DAEMON_NAME"]
        # Define the directory path where the binary will be stored
        directory_path = f"{os.environ['DAEMON_HOME']}/{os.environ['VISOR_NAME']}/upgrades/{formatted_version}/bin"
        # Check if binary already exists
        logger.log.info(f"CHECKING / DOWNLOADING {directory_path}/{end_binary_name}")

        if os.path.exists(f"{directory_path}/{end_binary_name}"):
            # If binary exists, log and do nothing
            logger.log.info(f"BINARY EXISTS ALREADY: {directory_path}/{end_binary_name}")
        else:
            # If binary doesn't exist, download and save it
            logger.log.info("BINARY DOES NOT EXIST.")
            os.makedirs(directory_path, exist_ok=True)
            response = requests.get(download_link)
            if response.status_code == 200:
                with open(f"{directory_path}/{end_binary_name}", "wb") as f:
                    f.write(response.content)
                os.chmod(f"{directory_path}/{end_binary_name}", 0o755)
                logger.log.info("BINARY DOWNLOADED SUCCESSFULLY.")
            else:
                logger.log.info("FAILED TO DOWNLOAD BINARY. Status code:", response.status_code)

    logger.log.info("BINARIES DOWNLOAD FINISHED...")

    # Start the process of upgrading binaries to the latest patch version
    # versions = set()
    # logger.log.info("UPGRADING BINARIES WITH LATEST PATCH VERSION UPGRADE")
    # # Collect versions of all binaries
    # for folder in os.listdir(upgrade_path):
    #     match = re.match(r'v(\d+)\.(\d+)\.(\d+)', folder)
    #     if match:
    #         versions.add(match.groups())
    #
    # # For each version, find and replace with the latest patch version if applicable
    # for major_version, minor_version, patch_version in versions:
    #     logger.log.info(f"BINARY VERSION: v{major_version}.{minor_version}.{patch_version}")
    #     latest_patch_version_path = find_latest_patch_version(major_version, minor_version)
    #     if latest_patch_version_path:
    #         logger.log.info(f"LATEST PATCH VERSION: {latest_patch_version_path}")
    #         symlink_path = os.path.join(upgrade_path, f"v{major_version}.{minor_version}.{patch_version}", "bin",
    #                                     "zetacored")
    #         logger.log.info(f"UPDATING BINARY: {symlink_path} TO: {latest_patch_version_path}")
    #         replace_binary(latest_patch_version_path, symlink_path)
    #     else:
    #         logger.log.info(f"NO PATCH UPDATE FOR v{major_version}.{minor_version}")

except Exception as e:
    logger.log.error(str(e))
