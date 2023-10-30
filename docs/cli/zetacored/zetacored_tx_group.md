# tx group

Group transaction subcommands

```
zetacored tx group [flags]
```

### Options

```
  -h, --help   help for group
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) 
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tx](zetacored_tx.md)	 - Transactions subcommands
* [zetacored tx group create-group](zetacored_tx_group_create-group.md)	 - Create a group which is an aggregation of member accounts with associated weights and an administrator account.
* [zetacored tx group create-group-policy](zetacored_tx_group_create-group-policy.md)	 - Create a group policy which is an account associated with a group and a decision policy. Note, the '--from' flag is ignored as it is implied from [admin].
* [zetacored tx group create-group-with-policy](zetacored_tx_group_create-group-with-policy.md)	 - Create a group with policy which is an aggregation of member accounts with associated weights, an administrator account and decision policy.
* [zetacored tx group draft-proposal](zetacored_tx_group_draft-proposal.md)	 - Generate a draft proposal json file. The generated proposal json contains only one message (skeleton).
* [zetacored tx group exec](zetacored_tx_group_exec.md)	 - Execute a proposal
* [zetacored tx group leave-group](zetacored_tx_group_leave-group.md)	 - Remove member from the group
* [zetacored tx group submit-proposal](zetacored_tx_group_submit-proposal.md)	 - Submit a new proposal
* [zetacored tx group update-group-admin](zetacored_tx_group_update-group-admin.md)	 - Update a group's admin
* [zetacored tx group update-group-members](zetacored_tx_group_update-group-members.md)	 - Update a group's members. Set a member's weight to "0" to delete it.
* [zetacored tx group update-group-metadata](zetacored_tx_group_update-group-metadata.md)	 - Update a group's metadata
* [zetacored tx group update-group-policy-admin](zetacored_tx_group_update-group-policy-admin.md)	 - Update a group policy admin
* [zetacored tx group update-group-policy-decision-policy](zetacored_tx_group_update-group-policy-decision-policy.md)	 - Update a group policy's decision policy
* [zetacored tx group update-group-policy-metadata](zetacored_tx_group_update-group-policy-metadata.md)	 - Update a group policy metadata
* [zetacored tx group vote](zetacored_tx_group_vote.md)	 - Vote on a proposal
* [zetacored tx group withdraw-proposal](zetacored_tx_group_withdraw-proposal.md)	 - Withdraw a submitted proposal

