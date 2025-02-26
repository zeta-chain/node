package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/zeta-chain/node/e2e/txserver"
)

type ExecuteProposalSequence string

const (
	StartOfE2E ExecuteProposalSequence = "start"
	EndOfE2E   ExecuteProposalSequence = "end"
)

// CreateGovProposals creates and votes on proposals from the given directory
// The directory should contain JSON files with the correct proposal format
// The directories currently used are:
// - /contrib/orchestrator/proposals_e2e_start
// - /contrib/orchestrator/proposals_e2e_end
func (r *E2ERunner) CreateGovProposals(sequence ExecuteProposalSequence) error {
	validatorsKeyring := r.ZetaTxServer.GetValidatorsKeyring()

	// List keys to verify that the keyring is working
	validatorKeys, err := validatorsKeyring.List()
	if err != nil {
		return fmt.Errorf("failed to list validator keys: %w", err)
	}
	if len(validatorKeys) == 0 {
		return fmt.Errorf("no validator keys found")
	}
	zetaTxServer := r.ZetaTxServer
	newZts := zetaTxServer.UpdateKeyring(validatorsKeyring)

	// Use the first validator as the depositor for all proposals
	depositor := validatorKeys[0]
	fromAddress, err := depositor.GetAddress()
	if err != nil {
		return fmt.Errorf("failed to get address for depositor: %w", err)
	}

	// Read the proposals directory
	proposalsDir := filepath.Join(fmt.Sprintf("/work/proposals_e2e_%s/", sequence))
	files, err := os.ReadDir(proposalsDir)
	if err != nil {
		return fmt.Errorf("failed to read proposals directory: %w", err)
	}
	// Process each file
	for _, file := range files {
		// Skip non-JSON files
		if !strings.HasSuffix(file.Name(), ".json") || file.IsDir() {
			continue
		}

		proposalPath := filepath.Join(proposalsDir, file.Name())

		// Parse the proposal file
		// Returning an error here as all proposals added to the directory should be valid
		parsedProposal, msgs, deposit, err := parseSubmitProposal(r.ZetaTxServer.GetCodec(), proposalPath)
		if err != nil {
			return fmt.Errorf("failed to parse proposal file %s: %w", file.Name(), err)
		}
		r.Logger.Print("executing proposal : file name: %s title: %s", file.Name(), parsedProposal.Title)

		// Create the proposal message
		msg, err := govv1.NewMsgSubmitProposal(
			msgs,
			deposit,
			fromAddress.String(),
			parsedProposal.Metadata,
			parsedProposal.Title,
			parsedProposal.Summary,
			parsedProposal.Expedited,
		)
		if err != nil {
			return fmt.Errorf("failed to create proposal message %s: %w", file.Name(), err)
		}

		// Broadcast the transaction
		res, err := newZts.BroadcastTx(depositor.Name, msg)
		if err != nil {
			return fmt.Errorf("failed to broadcast transaction for proposal %s: %w", file.Name(), err)
		}

		// Extract the proposal ID
		proposalID := uint64(0)
		for _, event := range res.Events {
			if event.Type == types.EventTypeSubmitProposal {
				for _, attr := range event.Attributes {
					if attr.Key == types.AttributeKeyProposalID {
						id, err := strconv.ParseUint(attr.Value, 10, 64)
						if err != nil {
							return err
						}
						proposalID = id
						break
					}
				}
			}
		}

		// First proposal ID is always 1
		if proposalID == 0 {
			return fmt.Errorf("failed to extract proposal ID from transaction for proposal %s", file.Name())
		}

		// Vote on the proposal
		err = voteGovProposals(proposalID, newZts, validatorKeys)
		if err != nil {
			return fmt.Errorf("failed to vote on proposal %d: %w", proposalID, err)
		}
	}
	return nil
}

// Vote yes on the proposal from both the validators
func voteGovProposals(proposalId uint64, zts txserver.ZetaTxServer, validatorKeys []*keyring.Record) error {
	for _, key := range validatorKeys {
		address, err := key.GetAddress()
		if err != nil {
			return fmt.Errorf("failed to get address for key %s: %w", key.Name, err)
		}
		// Create the message
		msg := govv1.NewMsgVote(address, proposalId, govv1.VoteOption_VOTE_OPTION_YES, "vote")

		_, err = zts.BroadcastTx(key.Name, msg)
		if err != nil {
			return fmt.Errorf("failed to broadcast transaction for vote on proposal %d: %w", proposalId, err)
		}
	}
	return nil
}

type proposal struct {
	// Msgs defines an array of sdk.Msgs proto-JSON-encoded as Anys.
	Messages  []json.RawMessage `json:"messages,omitempty"`
	Metadata  string            `json:"metadata"`
	Deposit   string            `json:"deposit"`
	Title     string            `json:"title"`
	Summary   string            `json:"summary"`
	Expedited bool              `json:"expedited"`
}

// parseSubmitProposal reads and parses the proposal.
func parseSubmitProposal(cdc codec.Codec, path string) (proposal, []sdk.Msg, sdk.Coins, error) {
	var proposal proposal

	contents, err := os.ReadFile(path)
	if err != nil {
		return proposal, nil, nil, err
	}

	err = json.Unmarshal(contents, &proposal)
	if err != nil {
		return proposal, nil, nil, err
	}

	msgs := make([]sdk.Msg, len(proposal.Messages))
	for i, anyJSON := range proposal.Messages {
		var msg sdk.Msg
		err := cdc.UnmarshalInterfaceJSON(anyJSON, &msg)
		if err != nil {
			return proposal, nil, nil, err
		}

		msgs[i] = msg
	}

	deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
	if err != nil {
		return proposal, nil, nil, err
	}

	return proposal, msgs, deposit, nil
}
