# filterdeposit

Filter missing inbound deposits

### Synopsis

Filters relevant inbound transactions for a given network and attempts to find an associated cctx from zetacore. If a 
cctx is not found, the associated transaction hash and amount is added to a list and displayed.

```
zetatool filterdeposit [command]
```
### Options

```
Available Commands:
btc         Filter inbound btc deposits
eth         Filter inbound eth deposits
```

### Flags
```
--btc-chain-id string   chain id used on zetachain to identify bitcoin - default: 8332 (default "8332")
```

### Options inherited from parent commands
```
--config string   custom config file: --config filename.json
```