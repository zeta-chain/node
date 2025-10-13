# Gov Precompile

The Gov precompile provides an EVM interface to the Cosmos SDK governance module,
enabling smart contracts to interact with on-chain governance proposals, voting, and deposits.

## Address

The precompile is available at the fixed address: `0x0000000000000000000000000000000000000805`

## Interface

### Data Structures

```solidity
enum VoteOption {
    Unspecified,    // 0 - No-op vote option
    Yes,            // 1 - Yes vote
    Abstain,        // 2 - Abstain vote
    No,             // 3 - No vote
    NoWithVeto      // 4 - No with veto vote
}

struct WeightedVoteOption {
    VoteOption option;
    string weight;      // Decimal string representation (e.g., "0.5")
}

struct ProposalData {
    uint64 id;
    string[] messages;  // Proposal messages in JSON format
    uint32 status;
    TallyResultData finalTallyResult;
    uint64 submitTime;
    uint64 depositEndTime;
    Coin[] totalDeposit;
    uint64 votingStartTime;
    uint64 votingEndTime;
    string metadata;
    string title;
    string summary;
    address proposer;
}

struct TallyResultData {
    string yes;
    string abstain;
    string no;
    string noWithVeto;
}
```

### Transaction Methods

```solidity
// Submit a new governance proposal
function submitProposal(
    address proposer,
    bytes calldata jsonProposal,
    Coin[] calldata deposit
) external returns (uint64 proposalId);

// Cancel an existing proposal
function cancelProposal(
    address proposer,
    uint64 proposalId
) external returns (bool success);

// Add deposit to a proposal
function deposit(
    address depositor,
    uint64 proposalId,
    Coin[] calldata amount
) external returns (bool success);

// Submit a simple vote
function vote(
    address voter,
    uint64 proposalId,
    VoteOption option,
    string memory metadata
) external returns (bool success);

// Submit a weighted vote
function voteWeighted(
    address voter,
    uint64 proposalId,
    WeightedVoteOption[] calldata options,
    string memory metadata
) external returns (bool success);
```

### Query Methods

```solidity
// Get a specific vote
function getVote(
    uint64 proposalId,
    address voter
) external view returns (WeightedVote memory vote);

// Get all votes for a proposal
function getVotes(
    uint64 proposalId,
    PageRequest calldata pagination
) external view returns (WeightedVote[] memory votes, PageResponse memory pageResponse);

// Get a specific deposit
function getDeposit(
    uint64 proposalId,
    address depositor
) external view returns (DepositData memory deposit);

// Get all deposits for a proposal
function getDeposits(
    uint64 proposalId,
    PageRequest calldata pagination
) external view returns (DepositData[] memory deposits, PageResponse memory pageResponse);

// Get tally results
function getTallyResult(
    uint64 proposalId
) external view returns (TallyResultData memory tallyResult);

// Get proposal details
function getProposal(
    uint64 proposalId
) external view returns (ProposalData memory proposal);

// Get proposals by status/voter/depositor
function getProposals(
    uint32 proposalStatus,
    address voter,
    address depositor,
    PageRequest calldata pagination
) external view returns (ProposalData[] memory proposals, PageResponse memory pageResponse);

// Get governance parameters
function getParams() external view returns (Params memory params);

// Get constitution
function getConstitution() external view returns (string memory constitution);
```

## Gas Costs

Gas costs are calculated dynamically based on the method and the Cosmos SDK operations performed.
The precompile uses the standard gas configuration for key-value operations.

## Implementation Details

### Proposal Submission

- Proposals are submitted in JSON format following Cosmos SDK proposal message structure
- The proposer must be the transaction sender
- Initial deposits can be included with the proposal
- Returns the newly created proposal ID

### Voting Mechanism

- **Simple voting**: Single vote option with full voting power
- **Weighted voting**: Multiple options with specified weights (must sum to 1.0)
- Votes are recorded in the Cosmos SDK governance module
- Metadata can be attached to votes for additional context

### Deposit Handling

- Deposits use the native token balance handler
- Deposits must meet minimum requirements defined in governance parameters
- The depositor must be the transaction sender

### Query Operations

- All queries are read-only and don't consume significant gas
- Pagination is supported for large result sets
- Proposal filtering by status, voter, or depositor

## Events

```solidity
event SubmitProposal(address indexed proposer, uint64 proposalId);
event CancelProposal(address indexed proposer, uint64 proposalId);
event Deposit(address indexed depositor, uint64 proposalId, Coin[] amount);
event Vote(address indexed voter, uint64 proposalId, uint8 option);
event VoteWeighted(address indexed voter, uint64 proposalId, WeightedVoteOption[] options);
```

## Security Considerations

1. **Sender Verification**: All transactions verify that the message sender matches the specified address parameter
2. **Balance Handling**: Uses the balance handler for proper native token management
3. **Permission Checks**: Inherits all permission checks from the Cosmos SDK governance module

## Usage Example

```solidity
IGov gov = IGov(GOV_PRECOMPILE_ADDRESS);

// Submit a proposal with initial deposit
Coin[] memory initialDeposit = new Coin[](1);
initialDeposit[0] = Coin({denom: "aevmos", amount: 1000000000000000000}); // 1 token

bytes memory proposalJSON = '{"messages":[...],"metadata":"...","title":"...","summary":"..."}';
uint64 proposalId = gov.submitProposal(msg.sender, proposalJSON, initialDeposit);

// Vote on the proposal
gov.vote(msg.sender, proposalId, VoteOption.Yes, "Supporting this proposal");

// Add additional deposit
Coin[] memory additionalDeposit = new Coin[](1);
additionalDeposit[0] = Coin({denom: "aevmos", amount: 500000000000000000}); // 0.5 token
gov.deposit(msg.sender, proposalId, additionalDeposit);

// Query proposal status
ProposalData memory proposal = gov.getProposal(proposalId);
```

## Integration Notes

- The precompile integrates directly with the Cosmos SDK governance module
- All governance parameters and rules apply
- Proposal JSON must be properly formatted according to Cosmos SDK standards
- Vote weights in weighted voting must be decimal strings that sum to "1.0"

