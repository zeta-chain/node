package chains

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// CosmosToEthChainID converts a Cosmos chain ID to an Ethereum chain ID
// parse value between _ and -
// e.g. cosmoshub_400-1 -> 400
func CosmosToEthChainID(chainID string) (int64, error) {
	// extract the substring
	extracted, err := extractBetweenUnderscoreAndDash(chainID)
	if err != nil {
		return 0, fmt.Errorf("can't convert cosmos to ethereum chain id: %w", err)
	}

	// convert to int64
	ethChainID, err := strconv.ParseInt(extracted, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("can't convert cosmos to ethereum chain id: %w", err)
	}

	return ethChainID, nil
}

func extractBetweenUnderscoreAndDash(s string) (string, error) {
	// Find the position of the underscore and dash
	underscoreIndex := strings.Index(s, "_")
	dashIndex := strings.Index(s, "-")

	// Check if both characters are found and in the correct order
	if underscoreIndex == -1 || dashIndex == -1 || underscoreIndex > dashIndex {
		return "", errors.New("value does not contain underscore followed by dash")
	}

	// Extract and return the substring
	return s[underscoreIndex+1 : dashIndex], nil
}
