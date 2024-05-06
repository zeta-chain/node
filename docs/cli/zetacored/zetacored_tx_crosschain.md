# tx crosschain

crosschain transactions subcommands

```
zetacored tx crosschain [flags]
```

### Options

```
  -h, --help   help for crosschain
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tx](zetacored_tx.md)	 - Transactions subcommands
* [zetacored tx crosschain abort-stuck-cctx](zetacored_tx_crosschain_abort-stuck-cctx.md)	 - abort a stuck CCTX
* [zetacored tx crosschain add-inbound-tracker](zetacored_tx_crosschain_add-inbound-tracker.md)	 - Add a in-tx-tracker 
				Use 0:Zeta,1:Gas,2:ERC20
* [zetacored tx crosschain add-outbound-tracker](zetacored_tx_crosschain_add-outbound-tracker.md)	 - Add a outbound-tracker
* [zetacored tx crosschain migrate-tss-funds](zetacored_tx_crosschain_migrate-tss-funds.md)	 - Migrate TSS funds to the latest TSS address
* [zetacored tx crosschain refund-aborted](zetacored_tx_crosschain_refund-aborted.md)	 - Refund an aborted tx , the refund address is optional, if not provided, the refund will be sent to the sender/tx origin of the cctx.
* [zetacored tx crosschain remove-outbound-tracker](zetacored_tx_crosschain_remove-outbound-tracker.md)	 - Remove a outbound-tracker
* [zetacored tx crosschain update-tss-address](zetacored_tx_crosschain_update-tss-address.md)	 - Create a new TSSVoter
* [zetacored tx crosschain vote-gas-price](zetacored_tx_crosschain_vote-gas-price.md)	 - Broadcast message to vote gas price
* [zetacored tx crosschain vote-inbound](zetacored_tx_crosschain_vote-inbound.md)	 - Broadcast message to vote an inbound
* [zetacored tx crosschain vote-outbound](zetacored_tx_crosschain_vote-outbound.md)	 - Broadcast message to vote an outbound
* [zetacored tx crosschain whitelist-erc20](zetacored_tx_crosschain_whitelist-erc20.md)	 - Add a new erc20 token to whitelist

