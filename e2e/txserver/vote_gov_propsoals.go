package txserver

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

func (zts ZetaTxServer) CreateGovProposal() {
	validatorKeys, err := zts.validatorKeys.List()
	if err != nil {
		panic(err)
	}
	newZts := zts.UpdateKeyring(zts.validatorKeys)

	depositor := validatorKeys[0]
	path := "/work/emissions_change.json"
	fromAddress, err := depositor.GetAddress()
	if err != nil {
		panic(err)
	}

	proposal, msgs, deposit, err := parseSubmitProposal(zts.clientCtx.Codec, path)
	if err != nil {
		panic(err)
	}
	msg, err := govv1.NewMsgSubmitProposal(msgs, deposit, fromAddress.String(), proposal.Metadata, proposal.Title, proposal.Summary, proposal.Expedited)
	if err != nil {
		panic(err)
	}

	res, err := newZts.BroadcastTx(depositor.Name, msg)
	if err != nil {
		panic(err)
	}
	proposalId := uint64(0)

	for _, event := range res.Events {
		if event.Type == types.EventTypeSubmitProposal {
			for _, attr := range event.Attributes {
				if string(attr.Key) == types.AttributeKeyProposalID {
					id, err := strconv.ParseUint(attr.Value, 10, 64)
					if err != nil {
						fmt.Printf("Error converting proposal ID to uint64: %v\n", err)
						continue
					}
					proposalId = id
				}
			}
		}
	}

	err = zts.VoteGovProposals(proposalId)
	if err != nil {
		panic(err)
	}

}

func (zts ZetaTxServer) VoteGovProposals(proposalId uint64) error {
	validatorKeys, err := zts.validatorKeys.List()
	if err != nil {
		return fmt.Errorf("failed to list validator keys: %w", err)
	}

	newZts := zts.UpdateKeyring(zts.validatorKeys)

	for _, key := range validatorKeys {
		address, err := key.GetAddress()
		if err != nil {
			panic(err)
		}
		// Create the message
		msg := govv1.NewMsgVote(address, proposalId, govv1.VoteOption_VOTE_OPTION_YES, "vote")

		res, err := newZts.BroadcastTx(key.Name, msg)
		if err != nil {
			panic(err)
		}
		fmt.Printf("VoteGovProposals: %s\n", res.TxHash)
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
