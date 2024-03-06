import requests
import os
import json
import logging
import sys

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

# Parse JSON from an environment variable to get binary download information
info = json.loads(os.environ["DOWNLOAD_BINARIES"])

try:
    # Iterate over binaries to download
    for binary in info["binaries"]:
        download_link = binary["download_url"]
        binary_location = f'{os.environ["DAEMON_HOME"]}/{binary["binary_location"]}'
        binary_directory = os.path.dirname(binary_location)
        # Log download link
        logger.log.info(f"DOWNLOAD LINK: {download_link}")
        split_download_link = download_link.split("/")
        # Log split download link parts
        logger.log.info(f"SPLIT DOWNLOAD LINK: {split_download_link}")
        # Extract binary name and version from the download link
        binary_name = download_link.split("/")[8]
        # Check if binary already exists
        logger.log.info(f"CHECKING / DOWNLOADING {binary_location}")

        if os.path.exists(binary_location):
            # If binary exists, log and do nothing
            logger.log.info(f"BINARY EXISTS ALREADY: {binary_location}")
        else:
            # If binary doesn't exist, download and save it
            logger.log.info("BINARY DOES NOT EXIST.")
            os.makedirs(binary_directory, exist_ok=True)
            response = requests.get(download_link)
            if response.status_code == 200:
                with open(binary_location, "wb") as f:
                    f.write(response.content)
                os.chmod(binary_location, 0o755)
                logger.log.info("BINARY DOWNLOADED SUCCESSFULLY.")
            else:
                logger.log.info("FAILED TO DOWNLOAD BINARY. Status code:", response.status_code)
    logger.log.info("BINARIES DOWNLOAD FINISHED...")
except Exception as e:
    logger.log.error(str(e))
