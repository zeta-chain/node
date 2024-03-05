#!/usr/bin/env python3
import hashlib
import logging
import os
import requests
import sys

# Constants defining the binary name, version, expected checksum, download URL, and installation path
BINARY_NAME = "cosmovisor"
BINARY_VERSION = os.getenv("COSMOVISOR_VERSION")  # Get the cosmovisor version from environment variable
EXPECTED_CHECKSUM = os.getenv("COSMOVISOR_CHECKSUM")  # Get the expected checksum from environment variable
BINARY_URL = f"https://binary-pickup.zetachain.com/cosmovisor-{BINARY_VERSION}-linux-amd64"  # Construct the binary download URL
INSTALL_PATH = f"/usr/local/bin/{BINARY_NAME}"  # Define the installation path for the binary

# Check if necessary environment variables are set; exit if not
if not BINARY_VERSION or not EXPECTED_CHECKSUM:
    logging.error("Environment variables COSMOVISOR_VERSION and COSMOVISOR_CHECKSUM must be set.")
    sys.exit(1)

# Configure logging to both stdout and a file
logging.basicConfig(
    level=logging.INFO,  # Set logging level to INFO
    format="%(levelname)s: %(message)s",  # Define log message format
    handlers=[
        logging.StreamHandler(sys.stdout),  # Log to stdout
        logging.FileHandler("/var/log/update_cosmovisor.log")  # Log to a file
    ]
)


# Function to calculate the SHA-256 checksum of the downloaded binary
def calculate_checksum(file_path):
    sha256 = hashlib.sha256()  # Create a new SHA-256 hash object
    with open(file_path, "rb") as f:  # Open the binary file in binary read mode
        for byte_block in iter(lambda: f.read(4096),
                               b""):  # Read the file in chunks to avoid loading it all into memory
            sha256.update(byte_block)  # Update the hash object with the chunk
    return sha256.hexdigest()  # Return the hexadecimal digest of the hash object


# Function to download the binary and update it if the checksum matches
def download_and_update_binary():
    try:
        response = requests.get(BINARY_URL)  # Attempt to download the binary
        response.raise_for_status()  # Check if the download was successful, raises exception on failure
        logging.info("Binary downloaded successfully.")
    except requests.exceptions.RequestException as e:
        logging.error(f"Failed to download the binary: {e}")  # Log any error during download
        sys.exit(1)  # Exit the script on download failure

    with open(INSTALL_PATH, "wb") as f:  # Open the installation path file in binary write mode
        f.write(response.content)  # Write the downloaded binary content to the file

    actual_checksum = calculate_checksum(INSTALL_PATH)  # Calculate the checksum of the downloaded binary
    if actual_checksum == EXPECTED_CHECKSUM:  # Compare the actual checksum with the expected checksum
        logging.info("Cosmovisor binary checksum verified.")  # Log success if checksums match
        os.chmod(INSTALL_PATH, 0o755)  # Make the binary executable
        logging.info("Cosmovisor binary updated successfully.")
    else:
        logging.error(
            "Checksums do not match. Possible corrupted download. Deleting the downloaded binary.")  # Log failure if checksums do not match
        os.remove(INSTALL_PATH)  # Remove the potentially corrupted binary
        sys.exit(1)  # Exit the script due to checksum mismatch


# Main script execution starts here
logging.info(
    f"Downloading the {BINARY_NAME} binary (version {BINARY_VERSION})...")  # Log the start of the download process
download_and_update_binary()  # Call the function to download and update the binary
